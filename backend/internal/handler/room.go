package handler

import (
	"elderly-care-system/internal/service"
	"elderly-care-system/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

// 获取楼栋列表
func (h *RoomHandler) ListBuildings(c *gin.Context) {
	buildings, err := h.roomService.ListBuildings()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, buildings)
}

// 获取楼栋详情（包含房间和床位）
func (h *RoomHandler) GetBuilding(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	building, err := h.roomService.GetBuildingWithRooms(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, building)
}

// 获取房间列表
func (h *RoomHandler) ListRooms(c *gin.Context) {
	floorID, _ := strconv.ParseUint(c.Query("floor_id"), 10, 32)
	rooms, err := h.roomService.ListRooms(uint(floorID))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, rooms)
}

// 获取床位统计
func (h *RoomHandler) GetBedStats(c *gin.Context) {
	stats, err := h.roomService.GetBedStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, stats)
}

// 分配床位
func (h *RoomHandler) AssignBed(c *gin.Context) {
	bedID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		ElderlyID uint `json:"elderly_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	err := h.roomService.AssignBed(uint(bedID), req.ElderlyID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

// 释放床位
func (h *RoomHandler) ReleaseBed(c *gin.Context) {
	bedID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	err := h.roomService.ReleaseBed(uint(bedID))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}
