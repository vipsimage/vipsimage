package operation

import (
	"github.com/sirupsen/logrus"
	"github.com/vipsimage/vips"
)

// Format option
type Format struct {
	Target string `json:"target" validate:"required,eq=jpg|eq=jpeg|eq=webp|eq=png|eq=tif|eq=tiff|eq=heif"`
}

// ImageBuffer contain formatted image info
type ImageBuffer struct {
	Content     []byte
	ContentType string
	Size        int
}

// FormatImage format vips.Image
func (th Format) FormatImage(img *vips.Image) (ib ImageBuffer, err error) {
	switch th.Target {
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
		logrus.WithField("format", th.Target).Warnln("format not found, use default format jpeg")

		ib.ContentType = "image/jpeg"
		ib.Content, ib.Size, err = img.JPEGSaveBuffer()
	}
	return
}
