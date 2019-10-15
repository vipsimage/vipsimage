package operation

import (
	"net/url"

	"github.com/vipsimage/vips"

	"github.com/vipsimage/vipsimage/storage/local"
)

func parseStorage(rule string) (storageOption keyValue, err error) {
	v, err := url.ParseQuery(rule)
	if err != nil {
		return
	}

	storageOption.key = v.Get("platform")
	if storageOption.key == "" {
		// set default storage
		storageOption.key = "local"
	}

	v.Del("platform")
	storageOption.value = v.Encode()
	return
}

func GetStorageLoadFunc(name string) (StorageLoadFunc, bool) {
	switch name {
	case "local", "":
		return localLoad, true
	default:
		return nil, false
	}
}

type StorageLoadFunc func(filePath string) (img *vips.Image, err error)

func localLoad(filePath string) (img *vips.Image, err error) {
	return local.Load(filePath)
}
