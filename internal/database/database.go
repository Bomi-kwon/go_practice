package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// PostgreSQL 데이터베이스 연결 정보
	host     = "localhost"   // 데이터베이스 서버 주소 (예: localhost 또는 127.0.0.1)
	user     = "bbomi"       // PostgreSQL 사용자 이름
	password = "1204"        // PostgreSQL 사용자 비밀번호
	dbname   = "go_practice" // 연결할 데이터베이스 이름
	port     = "5432"        // PostgreSQL 포트 번호 (기본값: 5432)
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 실패: %v", err)
	}

	// 데이터베이스 마이그레이션
	// 이 부분 활성화시 데이터베이스에 테이블이 없으면 Base 테이블 자동으로 생성
	// 테이블 있으면 스키마 변경사항 자동으로 반영
	// if err := db.AutoMigrate(&model.Base{}); err != nil {
	// 	return nil, fmt.Errorf("마이그레이션 실패: %v", err)
	// }

	fmt.Println("데이터베이스 연결 성공")

	return db, nil
}
