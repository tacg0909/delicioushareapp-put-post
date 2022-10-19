package resize

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"io"

	"github.com/disintegration/imaging"
	"github.com/tacg0909/meshitero-put-post/calctargetsize"
	"golang.org/x/image/draw"
)

func Resize(imageBuf io.Reader, maxLength int) (resizedImage bytes.Buffer, err error) {
    decodedImage, err := imaging.Decode(imageBuf, imaging.AutoOrientation(true))
    if err != nil {
        return
    }
    r := decodedImage.Bounds()
    targetWidth, targetHeight := calctargetsize.CalcTargetSize(r.Dx(), r.Dy(), maxLength)
    resizedImage, err = resize(decodedImage, targetWidth, targetHeight)
    return
}

func resize(img image.Image, targetWidth int, targetHeight int) (resizedImage bytes.Buffer, err error) {
    destination := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
    draw.CatmullRom.Scale(
        destination,
        destination.Bounds(),
        img,
        img.Bounds(),
        draw.Over,
        nil,
    )
    err = jpeg.Encode(&resizedImage, destination, &jpeg.Options{Quality: 100})
    return
}

func resizeOld(imageBuf *bytes.Buffer, maxLength int) (err error) {
    decordedImage, _, err := image.Decode(imageBuf)
    if err != nil {
        return
    }
    rectangle := decordedImage.Bounds()
    width := rectangle.Dx()
    height := rectangle.Dy()
    targetWidth, targetHeight := calctargetsize.CalcTargetSize(width, height, maxLength)
    dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
    draw.CatmullRom.Scale(
        dst,
        dst.Bounds(),
        decordedImage,
        rectangle,
        draw.Over,
        nil,
    )
    err = jpeg.Encode(imageBuf, dst, &jpeg.Options{Quality: 100})
    return
}
