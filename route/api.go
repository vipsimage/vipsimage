package route

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/vipsimage/vipsimage/operation"
	"github.com/vipsimage/vipsimage/rule"
)

// bad is common error gin's json return, abort other middleware
func bad(c *gin.Context, code int, msg string) {
	logrus.Errorln(msg)

	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
	})

	c.Abort()
}

// getRules is GET: operation-rules API, return all operation rule
func getRules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": rule.GetAll(),
	})
}

type ruleParams struct {
	Alias         string `json:"alias"`
	OperationRule string `json:"operation_rule"`
}

// addRule is PUT: operation-rule API, add a operation rule.
func addRule(c *gin.Context) {
	var rp ruleParams
	err := c.ShouldBindJSON(&rp)
	if err != nil {
		bad(c, http.StatusBadRequest, err.Error())
		return
	}

	err = rule.Add(rp.Alias, rp.OperationRule)
	if err != nil {
		bad(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

// delRule is DELETE: operation-rule API, delete a operation rule.
func delRule(c *gin.Context) {
	alias, ok := c.GetQuery("alias")
	if !ok {
		bad(c, http.StatusBadRequest, "alias not found")
		return
	}

	rule.Del(alias)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

// updateRule is POST: operation-rule API, update one operation rule.
func updateRule(c *gin.Context) {
	var rp ruleParams
	err := c.ShouldBindJSON(&rp)
	if err != nil {
		bad(c, http.StatusBadRequest, err.Error())
		return
	}

	err = rule.Set(rp.Alias, rp.OperationRule)
	if err != nil {
		bad(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

// defaultHandle handle NoRoute
func defaultHandle(c *gin.Context) {
	// default get method
	switch c.Request.URL.Path {
	case "/operation-rules":
		getRules(c)
		return
	}

	// handle image
	c.Params = append(c.Params, gin.Param{
		Key:   "operation-rule",
		Value: "default",
	}, gin.Param{
		Key:   "original-path",
		Value: c.Request.URL.Path,
	})

	HandleImages(c)
}

// parseRule is POST: parse-rule API, return a parsed operation rule.
func parseRule(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		bad(c, http.StatusBadRequest, err.Error())
		return
	}

	op, err := operation.Parse(string(b))
	if err != nil {
		bad(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"msg":    "ok",
		"data":   op,
		"encode": base64.StdEncoding.EncodeToString(b),
	})
}

// patchRule test HandleImages
func patchRule(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		bad(c, http.StatusBadRequest, err.Error())
		return
	}

	// handle image
	c.Params = append(c.Params, gin.Param{
		Key:   "operation-rule",
		Value: base64.StdEncoding.EncodeToString(b),
	}, gin.Param{
		Key:   "original-path",
		Value: c.Request.URL.Path,
	})

	HandleImages(c)
}

// HandleImages process images according to operation rule.
func HandleImages(c *gin.Context) {
	operationRule := c.Param("operation-rule")
	originalPath := c.Param("original-path")

	var err error
	op, ok := rule.Get(operationRule)
	if !ok {
		// enable full rule
		if !viper.GetBool("vipsimage.enable-full-rule") {
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

	// use storage config, load image
	img, err := op.Storage.Load(originalPath)
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
			"msg":  "ok",
			"data": res,
			"info": gin.H{
				"size":         ib.Size,
				"content-type": ib.ContentType,
			},
		})
		return
	}

	// return image
	c.Header("Server", fmt.Sprintf("vipsimage/%s", rule.Version))
	c.Header("Content-Length", fmt.Sprint(ib.Size))
	c.Data(http.StatusOK, ib.ContentType, ib.Content)
}
