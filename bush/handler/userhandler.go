package handler

import (
	"bush/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Register(ctx *gin.Context) {
	var input *entity.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	//err := h.services.Authorization.CreateFirebaseUser(input)
	//if err != nil {
	//	newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	//	return
	//}

	err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	//
	//	ctx.JSON(http.StatusOK, map[string]interface{}{
	//		"id": id,
	//	})
}

func (h *Handler) Login(ctx *gin.Context) {
	var input entity.Token

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
