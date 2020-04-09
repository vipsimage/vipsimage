package route

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // pprof
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Route return vipsimage define route
func Route() *gin.Engine {
	go func() {
		fmt.Println("bind pprof on :6060")
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	r := gin.Default()

	if viper.GetBool("auth.without-get-image") {
		r.Handle("VIEW", "*original-path", patchRule)
		r.GET("/:operation-rule/*original-path", HandleImages).Use(authMiddleWare())
	} else {
		r.Use(authMiddleWare()).GET("/:operation-rule/*original-path", HandleImages)
		r.Handle("VIEW", "*original-path", patchRule)
	}

	r.PUT("operation-rule", addRule)
	r.DELETE("operation-rule", delRule)
	r.POST("operation-rule", updateRule)
	r.POST("parse-rule", parseRule)
	r.NoRoute(defaultHandle)

	return r
}

// authMiddleWare auth any request
func authMiddleWare() gin.HandlerFunc {
	authType := viper.GetString("auth.type")
	if strings.ToLower(authType) == "basic" &&
		viper.GetBool("auth.enable") {

		name := viper.GetString("auth.name")
		pass := viper.GetString("auth.pass")

		return gin.BasicAuth(map[string]string{
			name: pass,
		})
	}

	return func(c *gin.Context) {}
}
