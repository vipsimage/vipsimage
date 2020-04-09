package operation

import (
	"reflect"

	"github.com/vipsimage/vips"
)

// Operation operation
type Operation struct {
	*Thumbnail `json:"thumbnail,omitempty"`
	*Resize    `json:"resize,omitempty"`
	*Crop      `json:"crop,omitempty"`
	*SmartCrop `json:"smart-crop,omitempty"`
	*Watermark `json:"watermark,omitempty"`
	*Rotate    `json:"rotate,omitempty"`
}

// Execute image process
func (th Operation) Execute(img *vips.Image) (err error) {
	v := reflect.ValueOf(th)

	for i := 0; i < v.NumField(); i++ {
		filed := v.Field(i)
		if !filed.IsZero() {
			err = filed.Interface().(Handler).Handle(img)
			if err != nil {
				return
			}
		}
	}
	return
}

// Handlers is handle function interface
type Handler interface {
	Handle(img *vips.Image) (err error)
}

// Thumbnail handle option
type Thumbnail struct {
	Width int `json:"width,omitempty"`
}

// Handle image thumbnail
func (th Thumbnail) Handle(img *vips.Image) (err error) {
	err = img.ThumbnailImage(th.Width)
	return
}

// Resize option
type Resize struct {
	Scale float64 `json:"scale,omitempty"`
}

// Handle resize image
func (th Resize) Handle(img *vips.Image) (err error) {
	err = img.Resize(th.Scale)
	return
}

// Crop option
type Crop struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Handle crop image
func (th Crop) Handle(img *vips.Image) (err error) {
	err = img.Crop(th.Left, th.Top, th.Width, th.Height)
	return
}

// SmartCrop option
type SmartCrop struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Handle smart crop image
func (th SmartCrop) Handle(img *vips.Image) (err error) {
	err = img.SmartCrop(th.Width, th.Height)
	return
}

// rotate image option
type Rotate struct {
	Angle float64 `json:"angle,omitempty"`
}

// Handle rotate image
func (th Rotate) Handle(img *vips.Image) (err error) {
	err = img.Rotate(th.Angle)
	return
}
