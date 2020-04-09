package operation

import (
	"strings"

	"github.com/vipsimage/pimg"
	"github.com/vipsimage/vips"
)

// Watermark option
type Watermark struct {
	Storage `json:"storage"`

	Img       string  `json:"img"`
	Direction string  `json:"direction"`
	Repeat    bool    `json:"repeat"`
	Scale     float64 `json:"scale"`
	AutoScale bool    `json:"auto_scale"`

	Angle   float64 `json:"angle"`
	OffsetX int     `json:"offset_x"`
	OffsetY int     `json:"offset_y"`
}

// Handle watermark image
// params:
//
// img:
//
// direction: optional, default south-east, e,s,w,t,se etc
// repeat: optional, default false, watermark repeat, direction usually use north-west
// scale: optional, float, default scale to 1/7 of the target image width
// auto-scale: optional, default scale to 1/7 of the target image width
// offset-x, offset-y: optional, int, offset of watermark
// angle: optional, the Angle of clockwise rotation
func (th Watermark) Handle(img *vips.Image) (err error) {
	pmg := pimg.Pimg{Image: img}

	wmVips, err := th.Load(th.Img)
	if err != nil {
		return
	}
	wm := &pimg.Pimg{Image: wmVips}

	option := pimg.NewWatermarkOption()
	option.Direction(watermarkDirection(th.Direction))

	if th.Repeat {
		option.WatermarkRepeat()
	}

	// watermark scale
	if th.Scale != 0 {
		option.WatermarkScale(th.Scale)
	}

	// auto scale
	option.WatermarkAutoScale(th.AutoScale)

	// rotation
	option.Rotation(th.Angle)

	// set offset
	option.SetOffset(th.OffsetX, th.OffsetY)

	return pmg.Watermark(wm, option)
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
