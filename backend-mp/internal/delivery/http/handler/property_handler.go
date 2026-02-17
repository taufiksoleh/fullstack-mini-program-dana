package handler

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/delivery/http/request"
	"github.com/taufiksoleh/backend-mp/internal/delivery/http/response"
	"github.com/taufiksoleh/backend-mp/internal/domain/apperror"
	"github.com/taufiksoleh/backend-mp/internal/domain/repository"
	"github.com/taufiksoleh/backend-mp/internal/usecase"
)

type PropertyHandler struct {
	usecase usecase.PropertyUsecase
}

func NewPropertyHandler(uc usecase.PropertyUsecase) *PropertyHandler {
	return &PropertyHandler{usecase: uc}
}

func httpStatus(err error) int {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, apperror.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, apperror.ErrValidation):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func (h *PropertyHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filter := repository.PropertyFilter{
		Keyword:  c.Query("keyword"),
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.usecase.GetAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(httpStatus(err), response.NewError(err.Error()))
		return
	}

	totalPages := int(math.Ceil(float64(result.Total) / float64(pageSize)))
	items := make([]response.PropertyListItem, len(result.Properties))
	for i := range result.Properties {
		items[i] = response.ToPropertyListItem(&result.Properties[i])
	}
	c.JSON(http.StatusOK, response.List(items, response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      result.Total,
		TotalPages: totalPages,
	}))
}

func (h *PropertyHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(errInvalidID))
		return
	}
	property, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(httpStatus(err), response.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(response.ToPropertyDetail(property)))
}

func (h *PropertyHandler) GetByDocumentID(c *gin.Context) {
	property, err := h.usecase.GetByDocumentID(c.Request.Context(), c.Param("documentId"))
	if err != nil {
		c.JSON(httpStatus(err), response.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(response.ToPropertyDetail(property)))
}

func (h *PropertyHandler) Create(c *gin.Context) {
	var body request.PropertyBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		return
	}
	property, err := h.usecase.Create(c.Request.Context(), body.ToUsecaseRequest())
	if err != nil {
		c.JSON(httpStatus(err), response.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, response.Success(response.ToPropertyDetail(property)))
}

func (h *PropertyHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(errInvalidID))
		return
	}
	var body request.PropertyBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		return
	}
	property, err := h.usecase.Update(c.Request.Context(), id, body.ToUsecaseRequest())
	if err != nil {
		c.JSON(httpStatus(err), response.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(response.ToPropertyDetail(property)))
}

func (h *PropertyHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(errInvalidID))
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(httpStatus(err), response.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
