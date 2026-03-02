package database

import (
	"fmt"
	"log"

	"elderly-care-system/internal/config"
	"elderly-care-system/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate - GORM will add missing columns automatically
	// This won't delete existing columns or change column types
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Elderly{},
		&model.ElderlyFamily{},
		&model.Building{},
		&model.Floor{},
		&model.Room{},
		&model.Bed{},
		&model.CareItem{},
		&model.CareStandard{},
		&model.CareRecord{},
		&model.ServiceRequest{},
		&model.FeeItem{},
		&model.Bill{},
		&model.BillItem{},
		&model.Payment{},
		&model.HealthRecord{},
		&model.Schedule{},
		&model.Medication{},
		&model.MedicationRecord{},
		&model.MedicationLog{},
		&model.VisitAppointment{},
		&model.Alert{},
		&model.AlertRule{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
