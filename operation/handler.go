package operation

import (
	"net/url"
	"strconv"

	"github.com/vipsimage/pimg"
	"github.com/vipsimage/vips"
)

func watermark(img *vips.Image, params string, op baseOperation) (err error) {
	v, err := url.ParseQuery(params)
	if err != nil {
		return
	}

	pmg := pimg.Pimg{Image: img}
	wmVips, err := op.Load(v.Get("img"))
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

func thumbnail(img *vips.Image, params string, _ baseOperation) (err error) {
	v, err := url.ParseQuery(params)
	if err != nil {
		return
	}

	width := v.Get("width")
	if len(v) == 1 && width == "" {
		width = params
	}

	w, err := strconv.Atoi(width)
	if err != nil {
		return
	}
	err = img.ThumbnailImage(w)
	return
}

func resize(img *vips.Image, params string, _ baseOperation) (err error) {
	v, err := url.ParseQuery(params)
	if err != nil {
		return
	}

	scale := v.Get("scale")
	s, err := strconv.ParseFloat(scale, 10)
	if err != nil {
		return
	}

	err = img.Resize(s)
	return
}

func crop(img *vips.Image, params string, _ baseOperation) (err error) {
	v, err := url.ParseQuery(params)
	if err != nil {
		return
	}
	left, top, width, height := v.Get("left"), v.Get("top"), v.Get("width"), v.Get("height")
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

func smartCrop(img *vips.Image, params string, _ baseOperation) (err error) {
	v, err := url.ParseQuery(params)
	if err != nil {
		return
	}

	width, height := v.Get("width"), v.Get("height")

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

type HandlerFunc func(img *vips.Image, params string, op baseOperation) error

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
