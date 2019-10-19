package route

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Route return vipsimage define route
func Route() *gin.Engine {
	r := gin.Default()

	if viper.GetBool("auth.without-get-image") {
		r.GET("/:operation-rule/*original-path", HandleImages).Use(authMiddleWare())
	} else {
		r.Use(authMiddleWare()).GET("/:operation-rule/*original-path", HandleImages)
	}

	r.PUT("operation-rule", addRule)
	r.DELETE("operation-rule", delRule)
	r.POST("operation-rule", updateRule)
	r.POST("operation-rule-parse", parseRule)
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
