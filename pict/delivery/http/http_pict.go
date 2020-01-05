package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/fidellr/edu_malay/model"
	"github.com/fidellr/edu_malay/pict"
	"github.com/fidellr/edu_malay/utils"
	"github.com/labstack/echo"
)

type Handler struct {
	service pict.PictUsecase
}

type Response struct {
	FileName string `json:"file_name"`
}

func NewPictHandler(e *echo.Echo, service pict.PictUsecase) {
	handler := &Handler{
		service,
	}

	e.POST("/picture/upload/:to/:name/:on_date", handler.Upload)
}

func (h *Handler) Upload(c echo.Context) error {
	var s3PictBucketInfo map[string]string
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	s3PictBucketInfo, err := determineS3PictBucketInfo(c.Param("to"), c.Param("name"), c.Param("on_date"))
	if s3PictBucketInfo == nil {
		return c.JSON(http.StatusInternalServerError, model.ResponseError{Message: err.Error()})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	fileName, err := h.service.Upload(ctx, file, s3PictBucketInfo)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, Response{FileName: fileName})
}

// func (h *Handler) Download()

func determineS3PictBucketInfo(toFolderName, prefixName, suffixName string) (map[string]string, error) {
	var s3PictBucketInfo map[string]string
	fileName := fmt.Sprintf("%s-%s", prefixName, suffixName)

	switch toFolderName {
	case "clc":
		s3PictBucketInfo = map[string]string{
			"folder_name": "clc",
			"file_name":   fileName,
			"aws_bucket":  "picture",
		}
		break
	case "teacher":
		s3PictBucketInfo = map[string]string{
			"folder_name": "teacher",
			"file_name":   fileName,
			"aws_bucket":  "picture",
		}
		break
	default:
		break
	}

	if s3PictBucketInfo == nil {
		return nil, errors.New("unsupported destination folder name")
	}

	return s3PictBucketInfo, nil
}
