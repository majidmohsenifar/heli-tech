package api

import (
	"net/http"

	"git.energy/ting/gateway-service/config"
	"git.energy/ting/gateway-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitialSwagger() {
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	r := gin.Default()
	url := ginSwagger.URL("./swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	go http.ListenAndServe(config.SwaggerUrl(), r)
}
