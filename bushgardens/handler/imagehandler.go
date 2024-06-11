package handler

import (
	"fmt"
	"net/http"
	//"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createAva(ctx *gin.Context) {
	//file, header, err := ctx.Request.FormFile("file")
	//if err != nil {
	//	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	//		"error": "Bad request",
	//	})
	//	return
	//}
	//defer file.Close()
	//
	//// pass the file and its name to the controller
	//ctx.Set("filePath", header.Filename)
	//ctx.Set("file", file)

	// another func
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	//var image entity.Image
	//
	//filename, ok := ctx.Get("filePath")
	//if !ok {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "filename not found"})
	//}
	//
	//fileGet, ok := ctx.Get("file")
	//if !ok {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "couldnt find file in request"})
	//	return
	//}
	fmt.Println(userId)
}
