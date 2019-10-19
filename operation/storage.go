package operation

import (
	"net/url"

	"github.com/vipsimage/vips"

	"github.com/vipsimage/vipsimage/storage/local"
)

// parseStorage double parse params
func parseStorage(params url.Values) (storageOption keyValue, err error) {
	storageOption.key = params.Get("platform")
	if storageOption.key == "" {
		// set default storage
		storageOption.key = "local"
	}

	params.Del("platform")
	storageOption.Values = params
	return
}

// GetStorageLoadFunc return load function by name.
func GetStorageLoadFunc(name string) (StorageLoadFunc, bool) {
	switch name {
	case "local", "":
		return localLoad, true
	default:
		return nil, false
	}
}

// StorageLoadFunc storage setting
type StorageLoadFunc func(filePath string, params keyValue) (img *vips.Image, err error)

// localLoad load image from local file
func localLoad(filePath string, _ keyValue) (img *vips.Image, err error) {
	return local.Load(filePath)
}
