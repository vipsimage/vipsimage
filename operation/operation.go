package operation

import (
	"fmt"
	"net/url"

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
		err = fmt.Errorf("storage load functin not found. name: %s", th.storage.key)
		return
	}

	return fn(filePath, th.storage)
}

type keyValue struct {
	url.Values

	key   string
	value string
}

// Operation rule
type Operation struct {
	baseOperation

	handlerFunc []keyValue // image handler
	computeFunc []keyValue // image computed
}

// Execute image process
func (th Operation) Execute(img *vips.Image) (err error) {
	for _, params := range th.handlerFunc {
		fn, ok := GetHandler(params.key)
		if !ok {
			return fmt.Errorf("handler function: %s, not found", params.key)
		}

		err = fn(img, params, th.baseOperation)
		if err != nil {
			logrus.Errorln(err.Error())
			continue
		}
	}

	return
}

// Compute content use compute func
func (th Operation) Compute(content []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})

	for _, ic := range th.computeFunc {
		fn, ok := GetComputeFunc(ic.key)
		if !ok {
			err = fmt.Errorf("compute function: %s, not found", ic.key)
			return
		}

		result, err := fn(content, ic, th.baseOperation)
		if err != nil {
			logrus.Errorln(err.Error())
			continue
		}

		res[ic.key] = result
	}

	return
}

// HasCompute return true if operation has compute function
func (th Operation) HasCompute() bool {
	return len(th.computeFunc) != 0
}
