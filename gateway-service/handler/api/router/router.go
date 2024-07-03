package router

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/majidmohsenifar/heli-tech/gateway-service/docs"
	"github.com/majidmohsenifar/heli-tech/gateway-service/handler/api"
	"github.com/majidmohsenifar/heli-tech/gateway-service/handler/api/middleware"
	"github.com/majidmohsenifar/heli-tech/gateway-service/service/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	Service = "gateway-service"
)

type Router struct {
	Handler *gin.Engine
	Routes  []Route
}

type Route struct {
	Path        string
	Name        string
	Method      string
	RouterGroup *gin.RouterGroup
	Handler     gin.HandlerFunc
}

func (r *Router) AddRoute(
	rg *gin.RouterGroup,
	method string,
	path string,
	name string,
	handler gin.HandlerFunc,
) error {
	switch method {
	case http.MethodPost:
		rg.POST(path, handler)
	case http.MethodGet:
		rg.GET(path, handler)
	case http.MethodPut:
		rg.PUT(path, handler)
	case http.MethodDelete:
		rg.DELETE(path, handler)
	default:
		return errors.New("invalid method")
	}
	route := Route{
		Path:        path,
		Name:        name,
		Method:      method,
		RouterGroup: rg,
		Handler:     handler,
	}
	r.Routes = append(r.Routes, route)
	return nil
}

// @license.name				Apache 2.0
// @license.url					http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey  ApiKeyAuth
// @in 							header
// @name 						Authorization
// @query.collection.format 	multi
// @externalDocs.description  	OpenAPI
// @externalDocs.url          	https://swagger.io/resources/open-api/
func New(
	userHandler *api.UserHandler,
	userService *user.Service,
	logger *slog.Logger,
) *Router {
	gin.SetMode(gin.ReleaseMode)

	docs.SwaggerInfo.Title = "Heli tech Gateway"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Heli tech API documentation"
	docs.SwaggerInfo.BasePath = "/"

	router := &Router{}
	r := gin.New()
	r.Use(cors.Default())
	r.Use(globalRecover(logger))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(
			http.StatusNotFound,
			api.ResponseFailure{
				Success: false,
				Error: api.ErrorCode{
					Code:    http.StatusNotFound,
					Message: "URL not found",
				},
			})
	})

	v1 := r.Group("/api/v1")
	//these route do not need token
	router.AddRoute(v1, http.MethodPost, "/auth/register", "register", userHandler.Register)
	router.AddRoute(v1, http.MethodPost, "/auth/login", "login", userHandler.Register)

	securedV1 := r.Group("/api/v1")
	securedV1.Use(middleware.JwtMiddleware(userService))
	{
		paymentRoutes := securedV1.Group("/payment")
		{
			//TODO: handle these routes later
			router.AddRoute(paymentRoutes, http.MethodPost, "/withdraw", "withdraw", userHandler.Register)
			router.AddRoute(paymentRoutes, http.MethodPost, "/deposit", "deposit", userHandler.Register)
		}
	}
	router.Handler = r
	return router
}

func globalRecover(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if rec := recover(); rec != nil {
				err := errors.New("error 500")
				logger.Error(fmt.Sprintf("error  500 in global recover %v", rec), err)
				api.MakeErrorResponseWithCode(c.Writer, http.StatusInternalServerError, "error 500")
			}
		}(c)
		c.Next()
	}
}
