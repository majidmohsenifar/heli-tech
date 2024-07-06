package api

import (
	"net/http"

	"github.com/majidmohsenifar/heli-tech/gateway-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	swaggerURL = "0.0.0.0:8081"
)

func InitiateSwagger() {
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	r := gin.Default()
	url := ginSwagger.URL("./swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	go http.ListenAndServe(swaggerURL, r)
}
