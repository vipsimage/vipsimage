package local

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vipsimage/vips"
)

var dataDir = getDataDir()

func getDataDir() string {
	if gin.IsDebugging() {
		return "data"
	}

	return "/data"
}

func FullPath(filePath string) string {
	filePath = "/" + strings.Trim(filePath, "/")
	return fmt.Sprintf("%s/images%s", dataDir, filePath)
}

func Load(filePath string) (*vips.Image, error) {
	fp := FullPath(filePath)

	return vips.NewFromFile(fp)
}
