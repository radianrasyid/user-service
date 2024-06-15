package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prithuadhikary/user-service/model"
	"github.com/prithuadhikary/user-service/service"
)

type AuthMiddleware interface {
	Auth(ctx *gin.Context)
}

type authController struct {
	service service.UserService
}

func (controller *authController) Auth(ctx *gin.Context) {
	cookie := ctx.Request.Header.Get("cookie")

	if cookie == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	user, err := controller.service.Whoami(&model.Whoami{
		SessionID: cookie,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	ctx.Set("user", user)

	ctx.Next()
}

func NewAuthMiddleware(service service.UserService) AuthMiddleware {
	return &authController{
		service: service,
	}
}
