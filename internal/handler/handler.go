package handler

import (
	"go_project/internal/model"
	"go_project/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc usecase.Usecase
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// 정적 파일 제공 설정
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// 메인 페이지 라우트
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// API 라우트
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// 최종 엔드포인트 URL들:
			// GET    /api/v1/resources     - 전체 리소스 목록 조회
			// GET    /api/v1/resources/:id - 특정 ID의 리소스 조회 (예: /api/v1/resources/1)
			// POST   /api/v1/resources     - 새로운 리소스 생성
			// PUT    /api/v1/resources/:id - 특정 ID의 리소스 수정 (예: /api/v1/resources/1)
			// DELETE /api/v1/resources/:id - 특정 ID의 리소스 삭제 (예: /api/v1/resources/1)
			v1.GET("/resources", h.GetAll)
			v1.GET("/resources/:id", h.Get)
			v1.POST("/resources", h.Insert)
			v1.PUT("/resources/:id", h.Modify)
			v1.DELETE("/resources/:id", h.Remove)
		}
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	results, err := h.uc.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "리소스 목록 조회 실패",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "성공",
		"data":    results,
	})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "잘못된 ID 형식",
			"data":    nil,
		})
		return
	}

	result, err := h.uc.Get(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "리소스 조회 실패",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "성공",
		"data":    result,
	})
}

func (h *Handler) Insert(c *gin.Context) {
	var resource model.Base
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "잘못된 요청 데이터",
			"data":    nil,
		})
		return
	}

	if err := h.uc.Insert(c, &resource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "리소스 생성 실패",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "성공",
		"data":    resource,
	})
}

func (h *Handler) Modify(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "잘못된 ID 형식",
			"data":    nil,
		})
		return
	}

	var resource model.Base
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "잘못된 요청 데이터",
			"data":    nil,
		})
		return
	}

	if err := h.uc.Modify(c, uint(id), &resource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "리소스 수정 실패",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "성공",
		"data":    resource,
	})
}

func (h *Handler) Remove(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "잘못된 ID 형식",
			"data":    nil,
		})
		return
	}

	if err := h.uc.Remove(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "리소스 삭제 실패",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "성공",
		"data":    nil,
	})
}
