package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CareHandler struct {
	careService *service.CareService
}

func NewCareHandler(careService *service.CareService) *CareHandler {
	return &CareHandler{careService: careService}
}

func (h *CareHandler) ListRecords(c *gin.Context) {
	elderlyID, _ := strconv.ParseUint(c.Query("elderly_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	records, total, err := h.careService.ListRecords(uint(elderlyID), page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  records,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func (h *CareHandler) CreateRecord(c *gin.Context) {
	staffID, _ := c.Get("user_id")

	var req service.CreateCareRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	record, err := h.careService.CreateRecord(&req, staffID.(uint))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, record)
}

func (h *CareHandler) GetMyTasks(c *gin.Context) {
	staffID, _ := c.Get("user_id")

	tasks, err := h.careService.GetMyTasks(staffID.(uint))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, tasks)
}

func (h *CareHandler) ListCareItems(c *gin.Context) {
	items, err := h.careService.ListCareItems()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, items)
}

// ListServiceRequests 获取服务请求列表
func (h *CareHandler) ListServiceRequests(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	requests, total, err := h.careService.ListServiceRequests(status, page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       requests,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// CreateServiceRequest 创建服务请求（老人呼叫）
func (h *CareHandler) CreateServiceRequest(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req service.CreateServiceRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	request, err := h.careService.CreateServiceRequest(&req, userID.(uint))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, request)
}

// HandleServiceRequest 处理服务请求
func (h *CareHandler) HandleServiceRequest(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.HandleServiceRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	err := h.careService.HandleServiceRequest(uint(id), req.Status, userID.(uint))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// Health Record handlers
func (h *CareHandler) ListHealthRecords(c *gin.Context) {
	elderlyID, _ := strconv.ParseUint(c.Query("elderly_id"), 10, 32)
	recordType := c.Query("record_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	records, total, err := h.careService.ListHealthRecords(uint(elderlyID), recordType, page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       records,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

func (h *CareHandler) GetLatestHealthRecords(c *gin.Context) {
	elderlyID, _ := strconv.ParseUint(c.Param("elderly_id"), 10, 32)

	records, err := h.careService.GetLatestHealthRecords(uint(elderlyID))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, records)
}

func (h *CareHandler) CreateHealthRecord(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req service.CreateHealthRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	record, err := h.careService.CreateHealthRecord(&req, userID.(uint))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, record)
}

func (h *CareHandler) DeleteHealthRecord(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	err := h.careService.DeleteHealthRecord(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}
