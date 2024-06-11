package handler

import (
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type errorst struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(ctx *gin.Context, statuscode int, message string) {
	logrus.Error(message)
	ctx.AbortWithStatusJSON(statuscode, errorst{message})
}
