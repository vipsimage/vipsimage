package operation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vipsimage/pimg"
	"github.com/vipsimage/vips"
)

// watermark image
// params:
//
// img: required, watermark image, load by operation rule.
//      if the image path contains slashes, you need to escape it twice.
//      e.g. test/test/vipsimage.png => test%2Ftest%2Fvipsimage.png
//      Golang: url.PathEscape(url.PathEscape("test/test/vipsimage.png"))
//      Javascript: encodeURI(encodeURIComponent('test/test/vipsimage.png'))
//
// direction: optional, default south-east, e,s,w,t,se etc
// repeat: optional, default false, watermark repeat, direction usually use north-west
// scale: optional, float, default scale to 1/7 of the target image width
// auto-scale: optional, default scale to 1/7 of the target image width
// offset-x, offset-y: optional, int, offset of watermark
// angle: optional, the Angle of clockwise rotation
func watermark(img *vips.Image, params keyValue, op baseOperation) (err error) {
	pmg := pimg.Pimg{Image: img}
	// defer pmg.Free()

	wmPath := params.Get("img")
	if wmPath == "" {
		err = fmt.Errorf("watermark image not found in this operation rule")
		return
	}

	wmVips, err := op.Load(wmPath)
	if err != nil {
		return
	}
	wm := &pimg.Pimg{Image: wmVips}
	// defer wm.Free()

	err = pmg.Watermark(wm, createWatermarkOption(params))
	return
}

func createWatermarkOption(params keyValue) *pimg.WatermarkOption { // nolint
	option := pimg.NewWatermarkOption()

	option.Direction(watermarkDirection(params.Get("direction")))
	// repeat
	if params.Get("repeat") != "" {
		option.WatermarkRepeat()
	}

	// watermark scale
	if scale := params.Get("scale"); scale != "" {
		s, err := strconv.ParseFloat(scale, 64)
		if err == nil && s != 0 {
			option.WatermarkScale(s)
		}
	}

	// auto scale
	if params.Get("auto-scale") != "" {
		option.WatermarkAutoScale(true)
	}

	// rotate
	if angle := params.Get("angle"); angle != "" {
		s, err := strconv.ParseFloat(angle, 64)
		if err == nil && s != 0 {
			option.Rotation(s)
		}
	}

	// set offset
	offsetX, offsetY, err := watermarkOffset(params)
	if err == nil && (offsetX != 0 || offsetY != 0) {
		option.SetOffset(offsetX, offsetY)
	}

	return option
}

func watermarkOffset(params keyValue) (offsetX, offsetY int, err error) {
	if ox := params.Get("offset-x"); ox != "" {
		offsetX, err = strconv.Atoi(ox)
		if err != nil {
			return
		}
	}

	if oy := params.Get("offset-y"); oy != "" {
		offsetY, err = strconv.Atoi(oy)
		if err != nil {
			return
		}
	}

	return
}

// watermarkDirection parse watermark direction
func watermarkDirection(direction string) vips.CompassDirection {
	direction = strings.ToLower(direction)

	switch direction {
	case "e", "east":
		return vips.DirectionEast
	case "s", "south":
		return vips.DirectionSouth
	case "w", "west":
		return vips.DirectionWest
	case "n", "north":
		return vips.DirectionNorth
	case "se", "south-east":
		return vips.DirectionSouthEast
	case "sw", "south-west":
		return vips.DirectionSouthWest
	case "ne", "north-east":
		return vips.DirectionNorthEast
	case "nw", "north-west":
		return vips.DirectionNorthWest
	case "c", "center":
		return vips.DirectionCentre
	}

	// default direction
	return vips.DirectionSouthEast
}
