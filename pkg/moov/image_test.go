package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
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

func Test_Images(t *testing.T) {
	t.Skip("todo: enable after resolving image upload limit")

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
		require.NotEmpty(t, uploaded.ImageID)
		require.NotEmpty(t, uploaded.PublicID)
		require.NotEmpty(t, uploaded.Link)
		require.Equal(t, metadata.AltText, uploaded.AltText)

		uploadedImageID = uploaded.ImageID
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
