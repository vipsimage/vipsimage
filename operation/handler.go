package operation

import (
	"strconv"

	"github.com/vipsimage/pimg"
	"github.com/vipsimage/vips"
)

// watermark image
func watermark(img *vips.Image, params keyValue, op baseOperation) (err error) {
	pmg := pimg.Pimg{Image: img}
	wmVips, err := op.Load(params.Get("img"))
	if err != nil {
		return
	}
	wm := &pimg.Pimg{Image: wmVips}
	defer wm.Free()

	option := pimg.NewWatermarkOption().
		Direction(vips.DirectionNorthWest).
		WatermarkRepeat().
		WatermarkAutoScale(true)

	err = pmg.Watermark(wm, option)
	return
}

// thumbnail image
func thumbnail(img *vips.Image, params keyValue, _ baseOperation) (err error) {
	width := params.Get("width")
	if width == "" {
		width = params.value
	}

	w, err := strconv.Atoi(width)
	if err != nil {
		return
	}
	err = img.ThumbnailImage(w)
	return
}

// resize image
func resize(img *vips.Image, params keyValue, _ baseOperation) (err error) {
	scale := params.Get("scale")
	s, err := strconv.ParseFloat(scale, 10)
	if err != nil {
		return
	}

	err = img.Resize(s)
	return
}

// crop image
func crop(img *vips.Image, params keyValue, _ baseOperation) (err error) {
	left, top, width, height := params.Get("left"), params.Get("top"), params.Get("width"), params.Get("height")
	l, err := strconv.Atoi(left)
	if err != nil {
		return
	}
	t, err := strconv.Atoi(top)
	if err != nil {
		return
	}
	w, err := strconv.Atoi(width)
	if err != nil {
		return
	}
	h, err := strconv.Atoi(height)
	if err != nil {
		return
	}
	err = img.Crop(l, t, w, h)
	return
}

// smartCrop smart crop image
func smartCrop(img *vips.Image, params keyValue, _ baseOperation) (err error) {
	width, height := params.Get("width"), params.Get("height")

	w, err := strconv.Atoi(width)
	if err != nil {
		return
	}
	h, err := strconv.Atoi(height)
	if err != nil {
		return
	}

	err = img.SmartCrop(w, h)
	return
}

// HandlerFunc handle image function
type HandlerFunc func(img *vips.Image, params keyValue, op baseOperation) error

// GetHandler return handler function by name.
func GetHandler(name string) (HandlerFunc, bool) {
	switch name {
	case "wm", "watermark":
		return watermark, true
	case "thumb", "thumbnail":
		return thumbnail, true
	case "crop":
		return crop, true
	case "sc", "smart-crop":
		return smartCrop, true
	case "resize":
		return resize, true
	default:
		return nil, false
	}
}
