package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/fidellr/edu_malay/clc"
	"github.com/fidellr/edu_malay/model"
	clcModel "github.com/fidellr/edu_malay/model/clc"
	"github.com/fidellr/edu_malay/utils"
)

type Handler struct {
	service clc.ProfileUsecase
}

func NewClcProfileHandler(e *echo.Echo, service clc.ProfileUsecase) {
	handler := &Handler{
		service,
	}

	e.GET("/clcs", handler.FindAll)
	e.GET("/clc/:id", handler.GetByID)
	e.POST("/clc", handler.Create)
	e.PUT("/clc/:id", handler.Update)
	e.POST("/clc/:id", handler.Remove)
	e.POST("/clc/assemble-profile/:clc_id/:teacher_id/:start_date", handler.AssembleProfile)
	e.PUT("/clc/assemble-profile/:clc_id/:teacher_id", handler.UpdateAssembledProfile)
}

func (h *Handler) FindAll(c echo.Context) error {
	var num int
	var err error
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if c.QueryParam("num") != "" {
		num, err = strconv.Atoi(c.QueryParam("num"))
		if err != nil {
			return utils.ConstraintErrorf("%s", err.Error())
		}
	}

	filter := &model.Filter{
		Num:    num,
		Cursor: c.QueryParam("cursor"),
		Search: c.QueryParam("search"),
	}

	clcs, nextCursor, err := h.service.FindAll(ctx, filter)
	if err != nil {
		return utils.ConstraintErrorf("%s", err.Error())
	}

	c.Response().Header().Set("X-Cursor", nextCursor)
	return c.JSON(http.StatusOK, clcs)
}

func (h *Handler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	clc := new(clcModel.ProfileEntity)
	if err := c.Bind(clc); err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	err := h.service.Create(ctx, clc)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	clc, err := h.service.GetByID(ctx, c.Param("id"))
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(utils.GetStatusCode(err), clc)
}

func (h *Handler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	clc := new(clcModel.ProfileEntity)
	if err := c.Bind(clc); err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	err := h.service.Update(ctx, c.Param("id"), clc)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) AssembleProfile(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.service.AssembleProfile(ctx, c.Param("clc_id"), c.Param("teacher_id"), c.Param("start_date"))
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) UpdateAssembledProfile(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	isEditingS := c.QueryParam("is_editing")
	var err error
	var isEditing bool
	if isEditingS != "" {
		isEditing, err = strconv.ParseBool(isEditingS)
		if err != nil {
			return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
		}
	}

	startDate := c.QueryParam("start_date")
	if isEditing && startDate == "" {
		return c.JSON(http.StatusBadRequest, model.ResponseError{Message: "editing required start_date queryparam :D"})
	}

	err = h.service.UpdateAssembledProfile(ctx, c.Param("clc_id"), c.Param("teacher_id"), startDate, isEditing)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) Remove(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	err := h.service.Remove(ctx, c.Param("id"))
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
