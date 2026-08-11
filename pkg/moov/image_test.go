package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageMetadataMarshal(t *testing.T) {
	input := []byte(`{
		"imageID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		"publicID": "qJRAaAwwF5hmfeAFdHjIb",
		"altText": "Test image",
		"link": "https://api.moov.io/images/qJRAaAwwF5hmfeAFdHjIb",
		"createdOn": "2024-01-15T10:30:00Z",
		"updatedOn": "2024-01-15T10:30:00Z",
		"disabledOn": "2024-01-15T10:30:00Z"
	}`)

	time, _ := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")

	metadata := new(moov.ImageMetadata)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&metadata)
	require.NoError(t, err)

	want := moov.ImageMetadata{
		ImageID:    "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		PublicID:   "qJRAaAwwF5hmfeAFdHjIb",
		AltText:    moov.PtrOf("Test image"),
		Link:       "https://api.moov.io/images/qJRAaAwwF5hmfeAFdHjIb",
		CreatedOn:  time,
		UpdatedOn:  time,
		DisabledOn: &time,
	}
	require.Equal(t, want, *metadata)
}

func TestImageConflictError(t *testing.T) {
	existingImage := moov.ImageMetadata{
		ImageID:   "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		PublicID:  "qJRAaAwwF5hmfeAFdHjIb",
		AltText:   moov.PtrOf("Test image"),
		Link:      "https://api.moov.io/images/qJRAaAwwF5hmfeAFdHjIb",
		CreatedOn: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		UpdatedOn: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	t.Run("implements error interface", func(t *testing.T) {
		err := &moov.ImageConflictError{
			ExistingImage: existingImage,
		}

		// Verify it implements error interface
		var asError error = err
		require.NotNil(t, asError)
		require.Contains(t, err.Error(), "duplicate image detected")
		require.Contains(t, err.Error(), existingImage.ImageID)
	})

	t.Run("can unwrap via errors.As", func(t *testing.T) {
		conflictErr := &moov.ImageConflictError{
			ExistingImage: existingImage,
		}

		var target *moov.ImageConflictError
		require.True(t, errors.As(conflictErr, &target))
		require.Equal(t, existingImage.ImageID, target.ExistingImage.ImageID)
		require.Equal(t, existingImage.PublicID, target.ExistingImage.PublicID)
		require.Equal(t, existingImage.AltText, target.ExistingImage.AltText)
		require.Equal(t, existingImage.Link, target.ExistingImage.Link)
		require.Equal(t, existingImage.CreatedOn, target.ExistingImage.CreatedOn)
		require.Equal(t, existingImage.UpdatedOn, target.ExistingImage.UpdatedOn)
	})
}

func TestUploadImageV2026_10(t *testing.T) {
	existingImage := moov.ImageMetadata{
		ImageID:   "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		PublicID:  "qJRAaAwwF5hmfeAFdHjIb",
		AltText:   moov.PtrOf("Test image"),
		Link:      "https://api.moov.io/images/qJRAaAwwF5hmfeAFdHjIb",
		CreatedOn: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		UpdatedOn: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	t.Run("sends v2026.10.00 version header", func(t *testing.T) {
		var capturedVersion string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedVersion = r.Header.Get(moov.VersionHeader)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"imageID": "new-image-id",
				"publicID": "newPublicId",
				"link": "https://api.moov.io/images/newPublicId",
				"createdOn": "2024-01-15T10:30:00Z",
				"updatedOn": "2024-01-15T10:30:00Z"
			}`))
		}))
		t.Cleanup(srv.Close)

		client, err := moov.NewClient(
			moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
			moov.WithMoovURLScheme("http"),
		)
		require.NoError(t, err)
		client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

		_, imgReader := randomImage(t, 100, 100, encodePNG)
		uploaded, err := client.UploadImageV2026_10(context.Background(), "account-123", imgReader, nil)

		require.NoError(t, err)
		require.NotNil(t, uploaded)
		assert.Equal(t, moov.Version2026_10.String(), capturedVersion)
		assert.Equal(t, "new-image-id", uploaded.ImageID)
	})

	t.Run("returns ImageConflictError on 409 with ImageMetadata body", func(t *testing.T) {
		imageJSON, _ := json.Marshal(existingImage)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write(imageJSON)
		}))
		t.Cleanup(srv.Close)

		client, err := moov.NewClient(
			moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
			moov.WithMoovURLScheme("http"),
		)
		require.NoError(t, err)
		client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

		_, imgReader := randomImage(t, 100, 100, encodePNG)
		uploaded, err := client.UploadImageV2026_10(context.Background(), "account-123", imgReader, nil)

		require.Nil(t, uploaded)
		require.Error(t, err)

		var conflictErr *moov.ImageConflictError
		require.True(t, errors.As(err, &conflictErr), "expected ImageConflictError, got %T: %v", err, err)
		assert.Equal(t, existingImage.ImageID, conflictErr.ExistingImage.ImageID)
		assert.Equal(t, existingImage.PublicID, conflictErr.ExistingImage.PublicID)
		assert.Equal(t, existingImage.AltText, conflictErr.ExistingImage.AltText)
		assert.Equal(t, existingImage.Link, conflictErr.ExistingImage.Link)
	})

	t.Run("returns generic error on 409 without ImageMetadata body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"error": "some conflict"}`))
		}))
		t.Cleanup(srv.Close)

		client, err := moov.NewClient(
			moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
			moov.WithMoovURLScheme("http"),
		)
		require.NoError(t, err)
		client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

		_, imgReader := randomImage(t, 100, 100, encodePNG)
		uploaded, err := client.UploadImageV2026_10(context.Background(), "account-123", imgReader, nil)

		require.Nil(t, uploaded)
		require.Error(t, err)

		// Should be a generic error, not ImageConflictError
		var conflictErr *moov.ImageConflictError
		require.False(t, errors.As(err, &conflictErr), "did not expect ImageConflictError, got one")
	})
}

func Test_Images(t *testing.T) {
	mc := NewTestClient(t)
	ctx := context.Background()
	accountID := MERCHANT_ID

	uploadedImageID := ""

	t.Run("upload image", func(t *testing.T) {
		_, imgReader := randomImage(t, 100, 100, encodePNG)
		metadata := &moov.ImageMetadataRequest{
			AltText: moov.PtrOf("Test image from moov-go SDK"),
		}

		uploaded, err := mc.UploadImage(ctx, accountID, imgReader, metadata)
		require.NoError(t, err)
		require.NotNil(t, uploaded)
		uploadedImageID = uploaded.ImageID
		require.NotEmpty(t, uploaded.ImageID)
		require.NotEmpty(t, uploaded.PublicID)
		require.NotEmpty(t, uploaded.Link)
		require.Equal(t, metadata.AltText, uploaded.AltText)
	})
	// cleanup image even if a below test fails
	t.Cleanup(func() {
		if uploadedImageID != "" {
			_ = mc.DeleteImage(ctx, accountID, uploadedImageID)
		}
	})

	t.Run("list images", func(t *testing.T) {
		gotImages, err := mc.ListImageMetadata(ctx, accountID)
		require.NoError(t, err)
		require.Greater(t, len(gotImages), 0)
	})

	t.Run("get image", func(t *testing.T) {
		got, err := mc.GetImageMetadata(ctx, accountID, uploadedImageID)
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("update image", func(t *testing.T) {
		_, imgReader := randomImage(t, 100, 100, encodeJPEG)
		metadata := &moov.ImageMetadataRequest{
			AltText: moov.PtrOf("Updated test image"),
		}

		updated, err := mc.UpdateImage(ctx, accountID, uploadedImageID, imgReader, metadata)
		require.NoError(t, err)
		require.NotNil(t, updated)
		require.Equal(t, metadata.AltText, updated.AltText)
	})

	t.Run("update image metadata", func(t *testing.T) {
		metadata := moov.ImageMetadataRequest{
			AltText: moov.PtrOf("Updated metadata only"),
		}

		updated, err := mc.UpdateImageMetadata(ctx, accountID, uploadedImageID, metadata)
		require.NoError(t, err)
		require.NotNil(t, updated)
		require.Equal(t, metadata.AltText, updated.AltText)
	})

	t.Run("delete image", func(t *testing.T) {
		err := mc.DeleteImage(ctx, accountID, uploadedImageID)
		require.NoError(t, err)
	})
}

type encoderFunc func(img image.Image) (io.Reader, error)

func encodePNG(img image.Image) (io.Reader, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	return &buf, err
}

func encodeJPEG(img image.Image) (io.Reader, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	return &buf, err
}

func randomImage(t *testing.T, w, h int, enc encoderFunc) (image.Image, io.Reader) {
	t.Helper()

	if w <= 0 || h <= 0 {
		t.Fatalf("invalid dimensions: %dx%d", w, h)
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(rnd.Intn(256)),
				G: uint8(rnd.Intn(256)),
				B: uint8(rnd.Intn(256)),
				A: 255,
			})
		}
	}

	var out io.Reader
	if enc != nil {
		var err error
		out, err = enc(img)
		require.NoError(t, err)
	}

	return img, out
}
