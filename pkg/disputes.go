package moov

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Dispute struct {
	DisputeID                string              `json:"disputeID,omitempty"`
	CreatedOn                time.Time           `json:"createdOn,omitempty"`
	Amount                   Amount              `json:"amount,omitempty"`
	NetworkReasonCode        string              `json:"networkReasonCode,omitempty"`
	NetworkReasonDescription string              `json:"networkReasonDescription,omitempty"`
	RespondBy                time.Time           `json:"respondBy,omitempty"`
	Status                   string              `json:"status,omitempty"`
	Transfer                 SynchronousTransfer `json:"transfer,omitempty"`
}

// ListDisputes lists of Disputes that are associated with a Moov account
// https://docs.moov.io/api/money-movement/disputes/list/
func (c Client) ListDisputes(count int,
	skip int,
	respondStartDateTime time.Time,
	respondEndDateTime time.Time,
	status string,
	merchantAccountID string,
	cardholderAccountID string,
	startDateTime time.Time,
	endDateTime time.Time,
	orderBy string) ([]Dispute, error) {
	var resp []Dispute

	values := url.Values{}
	// Convert time values to ISO8601 format
	respondStartDateTimeStr := respondStartDateTime.Format(time.RFC3339)
	respondEndDateTimeStr := respondEndDateTime.Format(time.RFC3339)
	startDateTimeStr := startDateTime.Format(time.RFC3339)
	endDateTimeStr := endDateTime.Format(time.RFC3339)

	// Add non-empty fields to the query string
	if count > 0 {
		values.Add("count", fmt.Sprint(count))
	}
	if skip > 0 {
		values.Add("skip", fmt.Sprint(skip))
	}
	if !respondStartDateTime.IsZero() {
		values.Add("respondStartDateTime", respondStartDateTimeStr)
	}
	if !respondEndDateTime.IsZero() {
		values.Add("respondEndDateTime", respondEndDateTimeStr)
	}
	if status != "" {
		values.Add("status", status)
	}
	if merchantAccountID != "" {
		values.Add("merchantAccountID", merchantAccountID)
	}
	if cardholderAccountID != "" {
		values.Add("cardholderAccountID", cardholderAccountID)
	}
	if !startDateTime.IsZero() {
		values.Add("startDateTime", startDateTimeStr)
	}
	if !endDateTime.IsZero() {
		values.Add("endDateTime", endDateTimeStr)
	}
	if orderBy != "" {
		values.Add("orderBy", orderBy)
	}

	urlStr := fmt.Sprintf("%s/%s?%s", baseURL, pathDisputes, values.Encode())

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, urlStr, nil, nil)
	if err != nil {
		return resp, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resp)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resp, nil
	case http.StatusTooManyRequests:
		return resp, ErrRateLimit
	}
	return resp, ErrDefault
}

// GetDispute retrieves a dispute for the given dispute id
// https://docs.moov.io/api/money-movement/disputes/get/
func (c Client) GetDispute(disputeID string) (Dispute, error) {
	resp := Dispute{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, pathDisputes, disputeID)

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resp, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resp)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resp, nil
	case http.StatusNotFound:
		return resp, ErrNoAccount
	case http.StatusTooManyRequests:
		return resp, ErrRateLimit
	}
	return resp, ErrDefault
}
