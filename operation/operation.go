package operation

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vipsimage/vips"
)

type baseOperation struct {
	rule    string   // operation rule
	storage keyValue // storage option
	format  string   // target format
}

// FormatImage format image
func (th baseOperation) FormatImage(img *vips.Image) (ib ImageBuffer, err error) {
	return ImageFormat(img, th.format)
}

func (th baseOperation) Load(filePath string) (img *vips.Image, err error) {
	fn, ok := GetStorageLoadFunc(th.storage.key)
	if !ok {
		err = errors.New(fmt.Sprintf("Storage load functin not found."))
		return
	}

	return fn(filePath)
}

type keyValue struct {
	key   string
	value string
}

type Operation struct {
	baseOperation

	handlerFunc []keyValue // image handler
	computeFunc []keyValue // image computed
}

func (th Operation) Execute(img *vips.Image) (err error) {
	for _, handle := range th.handlerFunc {
		fn, ok := GetHandler(handle.key)
		if !ok {
			return errors.New(fmt.Sprintf("function: %s, not found", handle.key))
		}

		err = fn(img, handle.value, th.baseOperation)
		if err != nil {
			logrus.Errorln(err.Error())
			continue
		}
	}

	return
}

func (th Operation) Compute(content []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})

	for _, ic := range th.computeFunc {
		fn, ok := GetComputeFunc(ic.key)
		if !ok {
			err = errors.New(fmt.Sprintf("function: %s, not found", ic.key))
			return
		}

		result, err := fn(content, th.baseOperation)
		if err != nil {
			logrus.Errorln(err.Error())
			continue
		}

		res[ic.key] = result
	}

	return
}

func (th Operation) HasCompute() bool {
	return len(th.computeFunc) != 0
}
