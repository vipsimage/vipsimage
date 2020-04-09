package operation

import (
	"github.com/vipsimage/vips"

	"github.com/vipsimage/vipsimage/storage/local"
)

// StorageLoader is storage interface
type StorageLoader interface {
	Load(filePath string) (img *vips.Image, err error)
}

// Storage option
type Storage struct {
}

func (th Storage) Load(filePath string) (img *vips.Image, err error) {
	return Load(filePath)
}

// Load local image
func Load(filePath string) (img *vips.Image, err error) {
	return local.Load(filePath)
}
