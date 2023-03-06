package v1

import (
	"grpc-research/bff-rest/forum/auth"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const thrityDaysInSecond = 30 * 24 * 60 * 60

func NewAuthHandler(authClient auth.AuthClient) authHandler {
	return authHandler{
		auth: authClient,
	}
}

type authHandler struct {
	auth auth.AuthClient
}

func (a *authHandler) PostLoginHandler(ctx *gin.Context) {
	var request struct {
		Username      string `json:"username" binding:"required_without=Email,omitempty,min=5,max=50"`
		Email         string `json:"email" binding:"required_without=Username,omitempty,email"`
		PlainPassword string `json:"password" binding:"required,min=8"`
	}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}
	if request.Email != "" {
		resp, err := a.auth.Login(ctx, &auth.LoginRequest{
			Auth:     &auth.LoginRequest_Email{Email: request.Email},
			Password: request.PlainPassword,
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.SetCookie("refreshToken", resp.GetRefreshKey(), thrityDaysInSecond, "/", os.Getenv("HOST"), true, true)
		ctx.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"data": gin.H{
				"accessToken": resp.GetAccessKey(),
			},
		})
		return
	}
	resp, err := a.auth.Login(ctx, &auth.LoginRequest{
		Auth:     &auth.LoginRequest_Username{Username: request.Username},
		Password: request.PlainPassword,
	})
	if err != nil {
		status, ok := status.FromError(err)
		if ok && status.Code() == codes.PermissionDenied {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "fail",
				"message": status.Err(),
			})
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return

	}

	ctx.SetCookie("refreshToken", resp.GetRefreshKey(), thrityDaysInSecond, "/", os.Getenv("HOST"), true, true)
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"accessToken": resp.GetAccessKey(),
		},
	})
}

func (a *authHandler) PutAccessTokenHandler(ctx *gin.Context) {
	refreshTokenCookie, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "cookies not found in your browser, must be login",
		})
		return
	}
	resp, err := a.auth.RefreshAccess(ctx, &auth.RefreshRequest{RefreshKey: refreshTokenCookie})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"accessToken": resp.GetAccessKey(),
		},
	})
}
func (a *authHandler) DeleteLogoutHandler(ctx *gin.Context) {
	refreshTokenCookie, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "cookies not found in your browser, must be login",
		})
		return
	}
	_, err = a.auth.Logout(ctx, &auth.LogoutRequest{RefreshKey: refreshTokenCookie})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.SetCookie("refreshToken", "", -1, "/", os.Getenv("HOST"), true, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "log out success",
	})
}
