package pimg

import (
	"image"

	"github.com/vipsimage/vips"
)

type WatermarkOption struct {
	compassDirection vips.CompassDirection
	blendMode        vips.BlendMode
	wmReplicate      bool
	wmScale          float64
	wmAutoScale      bool

	offsetX, offsetY int
}

func NewWatermarkOption() *WatermarkOption {
	return &WatermarkOption{
		compassDirection: vips.DirectionCentre,
		blendMode:        vips.BlendModeOver,
		wmReplicate:      false,
	}
}

func (th *WatermarkOption) Direction(direction vips.CompassDirection) *WatermarkOption {
	th.compassDirection = direction
	return th
}

func (th *WatermarkOption) BlendMode(mode vips.BlendMode) *WatermarkOption {
	th.blendMode = mode
	return th
}

func (th *WatermarkOption) WatermarkRepeat() *WatermarkOption {
	th.wmReplicate = true
	return th
}

func (th *WatermarkOption) SetOffset(x, y int) *WatermarkOption {
	th.offsetX, th.offsetY = x, y
	return th
}

func (th *WatermarkOption) WatermarkScale(scale float64) *WatermarkOption {
	th.wmScale = scale
	return th
}

func (th *WatermarkOption) WatermarkAutoScale(auto bool) *WatermarkOption {
	th.wmAutoScale = auto
	return th
}

func (th *Pimg) Watermark(watermark *Pimg, op *WatermarkOption) (err error) {
	// copy watermark
	wm, err := watermark.Copy()
	if err != nil {
		return
	}
	defer wm.Free()

	var scale float64
	if op.wmScale <= 0 || op.wmAutoScale {
		scale = float64(th.Width()) / float64(wm.Width()) / 7
	} else {
		scale = op.wmScale
	}

	err = wm.Resize(scale)
	if err != nil {
		return
	}

	if op.wmReplicate {
		err = wm.Replicate(th.Width()/wm.Width()+1, th.Height()/wm.Height()+1)
		if err != nil {
			return
		}
	}

	// Calculation of position
	point := watermarkPosition(op.compassDirection, rect{
		width:  th.Width(),
		height: th.Height(),
	}, rect{
		width:  wm.Width(),
		height: wm.Height(),
	})

	// excursion
	point.X += op.offsetX
	point.Y += op.offsetY

	err = th.Composite2(wm, op.blendMode, point)
	if err != nil {
		return
	}

	return
}

type rect struct {
	width, height int
}

func watermarkPosition(direction vips.CompassDirection, base, wm rect) (point image.Point) {
	switch direction {
	case vips.DirectionCentre:
		point.X = (base.width - wm.width) / 2
		point.Y = (base.height - wm.height) / 2
	case vips.DirectionEast:
		point.X = base.width - wm.width
		point.Y = (base.height - wm.height) / 2
	case vips.DirectionSouth:
		point.X = (base.width - wm.width) / 2
		point.Y = base.height - wm.height
	case vips.DirectionWest:
		point.Y = (base.height - wm.height) / 2
	case vips.DirectionNorth:
		point.X = (base.width - wm.width) / 2
	case vips.DirectionNorthEast:
		point.X = base.width - wm.width
	case vips.DirectionSouthEast:
		point.X = base.width - wm.width
		point.Y = base.height - wm.height
	case vips.DirectionSouthWest:
		point.Y = base.height - wm.height
	case vips.DirectionNorthWest:
	}

	return
}
