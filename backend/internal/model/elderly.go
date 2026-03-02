package model

import (
	"time"
	"gorm.io/gorm"
)

type Elderly struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"size:50;not null"`
	Gender           string         `json:"gender" gorm:"size:10"`
	BirthDate        *time.Time     `json:"birth_date"`
	IDCard           string         `json:"id_card" gorm:"size:20"`
	Phone            string         `json:"phone" gorm:"size:20"`
	EmergencyContact string         `json:"emergency_contact" gorm:"size:50"`
	EmergencyPhone   string         `json:"emergency_phone" gorm:"size:20"`
	AdmissionDate    *time.Time     `json:"admission_date"`
	BedID            *uint          `json:"bed_id"`
	HealthStatus     string         `json:"health_status" gorm:"type:text"`
	CareLevel        int            `json:"care_level" gorm:"default:1"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	Bed        *Bed          `json:"bed,omitempty" gorm:"foreignKey:BedID"`
	Families   []ElderlyFamily `json:"families,omitempty" gorm:"foreignKey:ElderlyID"`
}

type Building struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:50;not null"`
	FloorsCount int       `json:"floors_count" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at"`
	Floors      []Floor   `json:"floors,omitempty" gorm:"foreignKey:BuildingID"`
}

type Floor struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	BuildingID uint   `json:"building_id" gorm:"not null"`
	Name      string  `json:"name" gorm:"size:50;not null"`
	SortOrder int     `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	Rooms     []Room  `json:"rooms,omitempty" gorm:"foreignKey:FloorID"`
}

type Room struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	FloorID     uint   `json:"floor_id" gorm:"not null"`
	Name        string `json:"name" gorm:"size:50;not null"`
	BedCapacity int    `json:"bed_capacity" gorm:"default:1"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	Beds        []Bed  `json:"beds,omitempty" gorm:"foreignKey:RoomID"`
}

type Bed struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	RoomID    uint       `json:"room_id" gorm:"not null"`
	Name      string     `json:"name" gorm:"size:50;not null"`
	Status    string     `json:"status" gorm:"size:20;default:empty"` // empty/occupied/reserved
	ElderlyID *uint      `json:"elderly_id"`
	SortOrder int        `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time  `json:"created_at"`
	Elderly   *Elderly   `json:"elderly,omitempty" gorm:"foreignKey:ElderlyID"`
}

type ElderlyFamily struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ElderlyID uint      `json:"elderly_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Relation  string    `json:"relation" gorm:"size:30"`
	IsPrimary bool      `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Elderly) TableName() string {
	return "elderly"
}

func (ElderlyFamily) TableName() string {
	return "elderly_family"
}
