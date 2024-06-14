package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prithuadhikary/user-service/helper"
	"github.com/prithuadhikary/user-service/model"
	"github.com/prithuadhikary/user-service/service"
	"github.com/prithuadhikary/user-service/util"
)

type UserController interface {
	Signup(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	Whoami(ctx *gin.Context)
}

type userController struct {
	service service.UserService
}

func (controller userController) Signup(ctx *gin.Context) {
	request := &model.SignupRequest{}
	if err := ctx.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
		util.RenderBindingErrors(ctx, err.(validator.ValidationErrors))
		return
	}
	err := controller.service.Signup(request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": err.Error(),
		})
	}

	ctx.JSON(200, gin.H{
		"message": "Successfully created user",
	})
}

func (controller userController) Signin(ctx *gin.Context) {
	request := &model.SigninRequest{}
	fmt.Println("ini request", request.Username)
	if err := ctx.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
		util.RenderBindingErrors(ctx, err.(validator.ValidationErrors))
		return
	}
	id, session, err := controller.service.Signin(request)
	fmt.Println("ini error signin", err)
	if err != nil {
		var status int
		var message string

		// Check the type of error returned by userService.Signin
		switch {
		case errors.Is(err, errors.New("username or password might be wrong")):
			status = http.StatusUnauthorized
			message = "Username or password might be wrong"
		case errors.Is(err, errors.New("database error occured")):
			status = http.StatusInternalServerError
			message = "Database error occurred"
		default:
			status = http.StatusNotAcceptable
			message = "Something went wrong"
		}

		ctx.AbortWithStatusJSON(status, gin.H{
			"message": message,
		})
		return
	}
	fmt.Println("ini session", session.ID)
	ctx.SetCookie("session_id", session.ID.String(), int(time.Until(session.ExpiresAt)), "/", "", false, true)

	ctx.JSON(http.StatusAccepted, gin.H{
		"message":    "Login successful",
		"session_id": session.ID,
		"id":         id,
	})
}

func (controller userController) ServiceChecking(ctx *gin.Context) {
	connection, load, err := helper.GetCurrentConnectionAndLoad("https://jsonplaceholder.typicode.com/posts")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"message": err.Error(),
		})

		return
	}
	fmt.Println("ini data server", connection, load)
	ctx.JSON(http.StatusAccepted, gin.H{
		"connection": connection,
		"load":       load,
	})
}

func (controller userController) Whoami(ctx *gin.Context) {
	cookie := ctx.Request.Header.Get("cookie")
	fmt.Println("ini data cookie", cookie)

	user, err := controller.service.Whoami(&model.Whoami{
		SessionID: cookie,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "success",
		"user":    user,
	})
}

func (controller userController) Signout(ctx *gin.Context) {
	request := &model.Signout{}

	if err := ctx.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
		util.RenderBindingErrors(ctx, err.(validator.ValidationErrors))
		return
	}

	err := controller.service.Signout(request)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"messasge": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully logout",
	})
}

func NewUserController(engine *gin.Engine, userService service.UserService) {
	controller := &userController{
		service: userService,
	}
	api := engine.Group("api")
	{
		api.POST("users", controller.Signup)
		api.POST("users/login", controller.Signin)
		api.POST("users/logout", controller.Signout)
		api.GET("users/whoami", controller.Whoami)
		api.GET("/service", controller.ServiceChecking)
	}
}
