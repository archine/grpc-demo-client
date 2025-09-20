package controller

import (
	"context"
	"grpc-demo-client/model/dto"
	"strconv"
	"time"

	"github.com/archine/gin-plus/v4/component/ioc"
	"github.com/archine/gin-plus/v4/resp"
	"github.com/archine/grpc-demo-proto/user"
	"github.com/gin-gonic/gin"
)

func init() {
	ioc.RegisterBeanDef(&UserController{})
}

type UserController struct {
	userClient user.UserServiceClient `autowire:""`
}

func (u *UserController) SetRoutes(router *gin.RouterGroup) {
	router.Group("/user").
		POST("/", u.create).
		GET("/", u.getById).
		GET("/list", u.list)
}

// create 创建用户
func (u *UserController) create(ctx *gin.Context) {
	var arg dto.CreateUserArg
	if !resp.ParamValidation(ctx, &arg) {
		return
	}

	apiCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	response, err := u.userClient.CreateUser(apiCtx, &user.CreateUserRequest{
		Name:  arg.Name,
		Email: arg.Email,
		Age:   int32(arg.Age),
	})
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	resp.Json(ctx, response.GetId())
}

// getById 根据ID获取用户
func (u *UserController) getById(ctx *gin.Context) {
	userIdStr := ctx.Query("id")
	if userIdStr == "" {
		resp.BadRequest(ctx, "id不能为空")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		resp.BadRequest(ctx, "无效的ID")
		return
	}

	apiCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	response, err := u.userClient.GetUser(apiCtx, &user.GetUserRequest{Id: int32(userId)})
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	resp.Json(ctx, response.GetUser())
}

// list 获取用户列表
func (u *UserController) list(ctx *gin.Context) {
	apiCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	response, err := u.userClient.FindUserList(apiCtx, &user.FindUserListRequest{})
	if err != nil {
		resp.Error(ctx, err)
		return
	}

	resp.Json(ctx, response.GetUsers())
}
