package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ElderlyHandler struct {
	elderlyService *service.ElderlyService
}

func NewElderlyHandler(elderlyService *service.ElderlyService) *ElderlyHandler {
	return &ElderlyHandler{elderlyService: elderlyService}
}

func (h *ElderlyHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	elderly, total, err := h.elderlyService.List(page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  elderly,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func (h *ElderlyHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	elderly, err := h.elderlyService.Get(uint(id))
	if err != nil {
		response.Error(c, 404, "Elderly not found")
		return
	}

	response.Success(c, elderly)
}

func (h *ElderlyHandler) Create(c *gin.Context) {
	var req service.CreateElderlyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	elderly, err := h.elderlyService.Create(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, elderly)
}

func (h *ElderlyHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.UpdateElderlyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	if err := h.elderlyService.Update(uint(id), &req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}
