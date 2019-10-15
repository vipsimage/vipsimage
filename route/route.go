package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/vipsimage/vipsimage/operation"
	"github.com/vipsimage/vipsimage/utils/conf"
)

var short = make(map[string]operation.Operation)
var enableShortRule = conf.Getenv("ENABLE_SHORT_RULE", "true") == "true"

func init() {
	p, err := operation.Parse("wm:img=car-watermark.png&auto=true&test=true;f:webp;thumb:400;compute-md5;")
	if err != nil {
		panic(err)
	}
	short["default"] = p
}

func Route() *gin.Engine {
	r := gin.Default()

	r.GET("/:operation-rule/*original-path", HandleImages)
	r.NoRoute(defaultHandle)

	return r
}

func bad(c *gin.Context, code int, msg string) {
	logrus.Errorln(msg)

	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
	})

	c.Abort()
}

func defaultHandle(c *gin.Context) {
	c.Params = append(c.Params, gin.Param{
		Key:   "operation-rule",
		Value: "default",
	}, gin.Param{
		Key:   "original-path",
		Value: c.Request.URL.Path,
	})

	HandleImages(c)
}

func HandleImages(c *gin.Context) {
	operationRule := c.Param("operation-rule")
	originalPath := c.Param("original-path")

	var err error
	op, ok := short[operationRule]
	if !ok {
		if !enableShortRule {
			bad(c, http.StatusForbidden, "short rule disable")
			return
		}

		// parse operation
		op, err = operation.Parse(operationRule)
		if err != nil {
			bad(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	img, err := op.Load(originalPath)
	if err != nil {
		bad(c, http.StatusBadRequest, err.Error())
		return
	}
	defer img.Free()

	// image handle
	err = op.Execute(img)
	if err != nil {
		bad(c, http.StatusInternalServerError, err.Error())
		return
	}

	// image format
	ib, err := op.FormatImage(img)
	if err != nil {
		bad(c, http.StatusInternalServerError, err.Error())
		return
	}

	// image compute
	if op.HasCompute() {
		res, err := op.Compute(ib.Content)
		if err != nil {
			bad(c, http.StatusInternalServerError, err.Error())
			return
		}

		// return computed result
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "",
			"data": res,
			"info": gin.H{
				"size":         ib.Size,
				"content-type": ib.ContentType,
			},
		})
		return
	}

	// return image
	c.Header("Server", "vipsimage/1.0.0")
	c.Header("Content-Length", fmt.Sprint(ib.Size))
	c.Data(http.StatusOK, ib.ContentType, ib.Content)
}
