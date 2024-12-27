package database

import (
	"fmt"

	"go_project/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	user     = "bbomi"
	password = "1204"
	dbname   = "go_practice"
	port     = "5432"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 실패: %v", err)
	}

	// 데이터베이스 마이그레이션
	if err := db.AutoMigrate(&model.Base{}); err != nil {
		return nil, fmt.Errorf("마이그레이션 실패: %v", err)
	}

	return db, nil
}
