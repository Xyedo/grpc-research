package v1

import (
	"grpc-research/bff-rest/forum/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewUserHandler(userClient user.UserClient) userHandler {
	return userHandler{
		user: userClient,
	}

}

type userHandler struct {
	user user.UserClient
}

func (h *userHandler) PostUserHandler(ctx *gin.Context) {
	var request struct {
		Username      string `json:"username" binding:"required,min=5,max=50"`
		Email         string `json:"email" binding:"required,email"`
		PlainPassword string `json:"password" binding:"required,min=8"`
	}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}
	resp, err := h.user.AddUser(ctx, &user.AddUserRequest{
		Username: request.Username,
		Email:    request.Email,
		Password: request.PlainPassword,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"id": resp.GetId(),
		},
	})
}
