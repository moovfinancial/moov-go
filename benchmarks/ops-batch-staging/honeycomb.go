package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"
)

const defaultHoneycombMCPURL = "https://mcp.honeycomb.io/mcp"

type honeycombMCPClient struct {
	URL        string
	Token      string
	SessionID  string
	HTTPClient *http.Client
}

type honeycombQueryResult struct {
	RunPK string
	Rows  []map[string]any
}

type jsonRPCEnvelope struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *honeycombMCPClient) runQuery(ctx context.Context, arguments map[string]any) (honeycombQueryResult, error) {
	var toolResult struct {
		IsError bool `json:"isError"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
			URI  string `json:"uri"`
		} `json:"content"`
	}
	if err := c.call(ctx, "tools/call", map[string]any{
		"name":      "run_query",
		"arguments": arguments,
	}, &toolResult); err != nil {
		return honeycombQueryResult{}, err
	}
	if toolResult.IsError {
		return honeycombQueryResult{}, errors.New(firstToolText(toolResult.Content))
	}

	var resourceURI string
	for _, content := range toolResult.Content {
		if content.Type == "resource_link" && strings.HasSuffix(content.URI, "/json") {
			resourceURI = content.URI
			break
		}
	}
	if resourceURI == "" {
		return honeycombQueryResult{}, errors.New("Honeycomb query did not return a JSON result resource")
	}

	var resourceResult struct {
		Contents []struct {
			Text string `json:"text"`
		} `json:"contents"`
	}
	if err := c.readResource(ctx, resourceURI, &resourceResult); err != nil {
		return honeycombQueryResult{}, err
	}
	if len(resourceResult.Contents) != 1 {
		return honeycombQueryResult{}, fmt.Errorf("Honeycomb result resource returned %d documents", len(resourceResult.Contents))
	}
	var raw struct {
		Results []struct {
			Data map[string]any `json:"data"`
		} `json:"results"`
	}
	if err := json.Unmarshal([]byte(resourceResult.Contents[0].Text), &raw); err != nil {
		return honeycombQueryResult{}, fmt.Errorf("decoding Honeycomb query result: %w", err)
	}
	rows := make([]map[string]any, 0, len(raw.Results))
	for _, result := range raw.Results {
		rows = append(rows, result.Data)
	}
	return honeycombQueryResult{RunPK: queryRunPK(resourceURI), Rows: rows}, nil
}

func (c *honeycombMCPClient) readResource(ctx context.Context, uri string, target any) error {
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		lastErr = c.call(ctx, "resources/read", map[string]any{"uri": uri}, target)
		if lastErr == nil {
			return nil
		}
		if attempt == 4 {
			break
		}
		timer := time.NewTimer(time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return fmt.Errorf("reading Honeycomb result resource: %w", lastErr)
}

func (c *honeycombMCPClient) call(ctx context.Context, method string, params any, target any) error {
	if c.Token == "" {
		return errors.New("HONEYCOMB_MCP_API_KEY is required")
	}
	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	})
	if err != nil {
		return fmt.Errorf("encoding Honeycomb MCP request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("creating Honeycomb MCP request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if c.SessionID != "" {
		req.Header.Set("Mcp-Session-Id", c.SessionID)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling Honeycomb MCP: %w", err)
	}
	defer resp.Body.Close()
	if sessionID := resp.Header.Get("Mcp-Session-Id"); sessionID != "" {
		c.SessionID = sessionID
	}
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
		return fmt.Errorf("Honeycomb MCP returned HTTP %d", resp.StatusCode)
	}

	data, err := readMCPResponseData(resp.Header.Get("Content-Type"), resp.Body)
	if err != nil {
		return err
	}
	var envelope jsonRPCEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return fmt.Errorf("decoding Honeycomb MCP response: %w", err)
	}
	if envelope.Error != nil {
		return fmt.Errorf("Honeycomb MCP error %d: %s", envelope.Error.Code, envelope.Error.Message)
	}
	if err := json.Unmarshal(envelope.Result, target); err != nil {
		return fmt.Errorf("decoding Honeycomb MCP result: %w", err)
	}
	return nil
}

func readMCPResponseData(contentType string, reader io.Reader) ([]byte, error) {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err == nil && mediaType == "application/json" {
		data, readErr := io.ReadAll(io.LimitReader(reader, 2<<20))
		if readErr != nil {
			return nil, fmt.Errorf("reading Honeycomb MCP JSON response: %w", readErr)
		}
		return data, nil
	}
	return readSSEData(reader)
}

func readSSEData(reader io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64<<10), 2<<20)
	var data strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" && data.Len() > 0 {
			return []byte(data.String()), nil
		}
		if strings.HasPrefix(line, "data:") {
			if data.Len() > 0 {
				data.WriteByte('\n')
			}
			data.WriteString(strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading Honeycomb MCP response: %w", err)
	}
	if data.Len() > 0 {
		return []byte(data.String()), nil
	}
	return nil, errors.New("Honeycomb MCP response contained no data event")
}

func firstToolText(content []struct {
	Type string `json:"type"`
	Text string `json:"text"`
	URI  string `json:"uri"`
}) string {
	for _, current := range content {
		if current.Type == "text" && current.Text != "" {
			return current.Text
		}
	}
	return "Honeycomb tool call failed"
}

func queryRunPK(uri string) string {
	trimmed := strings.TrimSuffix(uri, "/json")
	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}

func number(row map[string]any, key string) (int64, error) {
	value, ok := row[key]
	if !ok {
		return 0, fmt.Errorf("Honeycomb result missing %q", key)
	}
	number, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("Honeycomb result %q has type %T", key, value)
	}
	return int64(number), nil
}

func newHoneycombMCPClient(url, token string) *honeycombMCPClient {
	if url == "" {
		url = defaultHoneycombMCPURL
	}
	return &honeycombMCPClient{
		URL:        url,
		Token:      token,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}
