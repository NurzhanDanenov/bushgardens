package handler

import (
	"bush/entity"
	"fmt"
	"mime/multipart"
	"net/http"
	//"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createAva(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}
	defer file.Close()

	// pass the file and its name to the controller
	ctx.Set("filePath", header.Filename)
	ctx.Set("file", file)

	// another func
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	strUserId, ok := userId.(string)
	if !ok {
		fmt.Println("value is not a string")
		return
	}

	//filename, ok := ctx.Get("filePath")
	//if !ok {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "filename not found"})
	//}

	fileGet, ok := ctx.Get("file")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "couldnt find file in request"})
		return
	}

	// upload file
	imageId, err := h.services.UploadImage.Upload(strUserId, fileGet.(multipart.File))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": imageId,
	})
}

func (h *Handler) getAllImages(ctx *gin.Context) {
	var input entity.DateOfIMages
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	strUserId, ok := userId.(string)
	if !ok {
		fmt.Println("value is not a string")
		return
	}

	err = h.services.UploadImage.Download(strUserId, input.Date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": userId,
	})
}
