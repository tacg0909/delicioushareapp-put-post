package resize

import (
	"bytes"
	"image"
	"os"
	"testing"
)

func TestResize(t *testing.T) {
    targetWidth := 500
    targetHeight := 500
    img, err := os.Open("sample_image.jpg")
    if err != nil {
        t.Fatalf("error: %v", err)
    }
    i, err := Resize(img, targetWidth, targetHeight)
    c, format, err := image.DecodeConfig(bytes.NewReader(i.Bytes()))
    if format != "jpeg" {
        t.Fatalf("Resized image format is %s, But want jpeg", format)
    }
    if c.Width != targetWidth || c.Height != targetHeight {
        t.Fatalf("Resized image size is %dx%d, but want %dx%d", c.Width, c.Height, targetWidth, targetHeight)
    }
}
