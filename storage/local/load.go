package local

import (
	"fmt"
	"os"
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

func Load(filePath string) (img *vips.Image, err error) {
	fp := FullPath(filePath)

	_, err = os.Stat(fp)
	if err != nil {
		return
	}

	return vips.NewFromFile(fp)
}
