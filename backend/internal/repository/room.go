package repository

import (
	"elderly-care-system/internal/model"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// 获取所有楼栋
func (r *RoomRepository) ListBuildings() ([]model.Building, error) {
	var buildings []model.Building
	err := r.db.Preload("Floors", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Order("id ASC").Find(&buildings).Error
	return buildings, err
}

// 获取楼栋详情（包含楼层和房间）
func (r *RoomRepository) GetBuildingWithRooms(id uint) (*model.Building, error) {
	var building model.Building
	err := r.db.Preload("Floors.Rooms.Beds.Elderly", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&building, id).Error
	if err != nil {
		return nil, err
	}
	return &building, nil
}

// 获取所有房间及其床位
func (r *RoomRepository) ListRooms(floorID uint) ([]model.Room, error) {
	var rooms []model.Room
	query := r.db.Preload("Beds.Elderly")
	if floorID > 0 {
		query = query.Where("floor_id = ?", floorID)
	}
	err := query.Order("sort_order ASC").Find(&rooms).Error
	return rooms, err
}

// 获取房间详情
func (r *RoomRepository) GetRoom(id uint) (*model.Room, error) {
	var room model.Room
	err := r.db.Preload("Beds.Elderly").First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// 获取床位统计
func (r *RoomRepository) GetBedStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	var total, occupied, empty int64
	r.db.Model(&model.Bed{}).Count(&total)
	r.db.Model(&model.Bed{}).Where("status = ?", "occupied").Count(&occupied)
	r.db.Model(&model.Bed{}).Where("status = ?", "empty").Count(&empty)

	stats["total"] = total
	stats["occupied"] = occupied
	stats["empty"] = empty

	return stats, nil
}

// 更新床位状态
func (r *RoomRepository) UpdateBedStatus(bedID uint, status string) error {
	return r.db.Model(&model.Bed{}).Where("id = ?", bedID).Update("status", status).Error
}

// 分配床位
func (r *RoomRepository) AssignBed(bedID, elderlyID uint) error {
	return r.db.Model(&model.Bed{}).Where("id = ?", bedID).Updates(map[string]interface{}{
		"status":     "occupied",
		"elderly_id": elderlyID,
	}).Error
}

// 释放床位
func (r *RoomRepository) ReleaseBed(bedID uint) error {
	return r.db.Model(&model.Bed{}).Where("id = ?", bedID).Updates(map[string]interface{}{
		"status":     "empty",
		"elderly_id": nil,
	}).Error
}

// GetOccupancyRate 获取总入住率
func (r *RoomRepository) GetOccupancyRate() (float64, error) {
	var total, occupied int64
	r.db.Model(&model.Bed{}).Count(&total)
	r.db.Model(&model.Bed{}).Where("status = ?", "occupied").Count(&occupied)

	if total == 0 {
		return 0, nil
	}
	return float64(occupied) / float64(total) * 100, nil
}

// BuildingOccupancy 楼栋入住统计
type BuildingOccupancy struct {
	BuildingName  string
	TotalBeds     int
	OccupiedBeds  int
	OccupancyRate float64
}

// GetBuildingOccupancy 获取各楼栋入住率
func (r *RoomRepository) GetBuildingOccupancy() ([]BuildingOccupancy, error) {
	var results []BuildingOccupancy

	query := `
		SELECT
			b.name as building_name,
			COUNT(bed.id) as total_beds,
			SUM(CASE WHEN bed.status = 'occupied' THEN 1 ELSE 0 END) as occupied_beds,
			CASE
				WHEN COUNT(bed.id) > 0
				THEN ROUND(SUM(CASE WHEN bed.status = 'occupied' THEN 1 ELSE 0 END) * 100.0 / COUNT(bed.id), 2)
				ELSE 0
			END as occupancy_rate
		FROM buildings b
		LEFT JOIN floors f ON b.id = f.building_id
		LEFT JOIN rooms r ON f.id = r.floor_id
		LEFT JOIN beds bed ON r.id = bed.room_id
		GROUP BY b.id, b.name
		ORDER BY b.id
	`

	r.db.Raw(query).Scan(&results)
	return results, nil
}
