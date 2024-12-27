package handler

import (
	"net/http"
	"strconv"

	"go_project/internal/model"
	"go_project/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/resources", h.List)
			v1.GET("/resources/:id", h.Get)
			v1.POST("/resources", h.Create)
			v1.PUT("/resources/:id", h.Update)
			v1.DELETE("/resources/:id", h.Delete)
		}
	}
}

func (h *Handler) List(c *gin.Context) {
	results, err := h.svc.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "리소스 목록 조회 실패",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "성공",
		Data:    results,
	})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "잘못된 ID 형식",
			Data:    nil,
		})
		return
	}

	result, err := h.svc.Get(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "리소스 조회 실패",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "성공",
		Data:    result,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var resource interface{}
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "잘못된 요청 데이터",
			Data:    nil,
		})
		return
	}

	if err := h.svc.Create(c, resource); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "리소스 생성 실패",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, model.Response{
		Status:  http.StatusCreated,
		Message: "성공",
		Data:    resource,
	})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "잘못된 ID 형식",
			Data:    nil,
		})
		return
	}

	var resource interface{}
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "잘못된 요청 데이터",
			Data:    nil,
		})
		return
	}

	if err := h.svc.Update(c, uint(id), resource); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "리소스 수정 실패",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "성공",
		Data:    resource,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: "잘못된 ID 형식",
			Data:    nil,
		})
		return
	}

	if err := h.svc.Delete(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "리소스 삭제 실패",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "성공",
		Data:    nil,
	})
}
