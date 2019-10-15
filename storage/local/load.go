package local

import (
	"fmt"
	"strings"

	"github.com/vipsimage/vips"
)

func Load(filePath string) (*vips.Image, error) {
	fp := FullPath(filePath)

	return vips.NewFromFile(fp)
}

func FullPath(filePath string) string {
	filePath = "/" + strings.Trim(filePath, "/")
	return fmt.Sprintf("data/images%s", filePath)
}
