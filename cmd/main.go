package main

import (
	"log"

	"go_project/internal/database"
	"go_project/internal/handler"
	"go_project/internal/recorder"
	"go_project/internal/repository"
	"go_project/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	// DB 초기화
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("데이터베이스 초기화 실패: %v", err)
	}

	// Recorder, Repository, Usecase, Handler 초기화
	rec := recorder.NewRecorder(db)
	repo := repository.NewRepository(rec)
	uc := usecase.NewUsecase(repo)
	h := handler.NewHandler(uc)

	// Router 설정
	r := gin.Default()

	// 라우트 설정
	h.RegisterRoutes(r)

	// 서버 시작
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
