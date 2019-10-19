package operation

import (
	"github.com/sirupsen/logrus"
	"github.com/vipsimage/vips"
)

// ImageBuffer contain formatted image info
type ImageBuffer struct {
	Content     []byte
	ContentType string
	Size        int
}

// ImageFormat format vips.Image
func ImageFormat(img *vips.Image, format string) (ib ImageBuffer, err error) {
	switch format {
	case "jpg", "jpeg":
		ib.ContentType = "image/jpeg"
		ib.Content, ib.Size, err = img.JPEGSaveBuffer()
	case "webp":
		ib.ContentType = "image/webp"
		ib.Content, ib.Size, err = img.WEBPSaveBuffer()
	case "png":
		ib.ContentType = "image/png"
		ib.Content, ib.Size, err = img.PNGSaveBuffer()
	case "tif", "tiff":
		ib.ContentType = "image/tiff"
		ib.Content, ib.Size, err = img.TIFFSaveBuffer()
	case "heif":
		ib.ContentType = "image/heif"
		ib.Content, ib.Size, err = img.HEIFSaveBuffer()
	default:
		logrus.WithField("format", format).Warnln("format not found, use default format jpeg")

		ib.ContentType = "image/jpeg"
		ib.Content, ib.Size, err = img.JPEGSaveBuffer()
	}
	return
}
