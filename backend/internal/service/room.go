package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
)

type RoomService struct {
	roomRepo *repository.RoomRepository
}

func NewRoomService(roomRepo *repository.RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) ListBuildings() ([]model.Building, error) {
	return s.roomRepo.ListBuildings()
}

func (s *RoomService) GetBuildingWithRooms(id uint) (*model.Building, error) {
	return s.roomRepo.GetBuildingWithRooms(id)
}

func (s *RoomService) ListRooms(floorID uint) ([]model.Room, error) {
	return s.roomRepo.ListRooms(floorID)
}

func (s *RoomService) GetRoom(id uint) (*model.Room, error) {
	return s.roomRepo.GetRoom(id)
}

func (s *RoomService) GetBedStats() (map[string]int64, error) {
	return s.roomRepo.GetBedStats()
}

func (s *RoomService) AssignBed(bedID, elderlyID uint) error {
	return s.roomRepo.AssignBed(bedID, elderlyID)
}

func (s *RoomService) ReleaseBed(bedID uint) error {
	return s.roomRepo.ReleaseBed(bedID)
}
