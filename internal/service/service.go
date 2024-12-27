package service

import (
	"context"
	"fmt"
	"go_project/internal/repository"
)

type Service interface {
	Create(ctx context.Context, entity interface{}) error
	Get(ctx context.Context, id uint) (interface{}, error)
	List(ctx context.Context) ([]interface{}, error)
	Update(ctx context.Context, id uint, entity interface{}) error
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

// 기본 CRUD 구현
func (s *service) Create(ctx context.Context, entity interface{}) error {
	if err := s.repo.Create(ctx, entity); err != nil {
		return fmt.Errorf("생성 실패: %v", err)
	}
	return nil
}

func (s *service) Get(ctx context.Context, id uint) (interface{}, error) {
	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("조회 실패: %v", err)
	}
	return result, nil
}

func (s *service) List(ctx context.Context) ([]interface{}, error) {
	results, err := s.repo.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("목록 조회 실패: %v", err)
	}
	return results, nil
}

func (s *service) Update(ctx context.Context, id uint, entity interface{}) error {
	// 먼저 존재하는지 확인
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("업데이트할 엔티티를 찾을 수 없습니다: %v", err)
	}

	if err := s.repo.Update(ctx, entity); err != nil {
		return fmt.Errorf("업데이트 실패: %v", err)
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	// 먼저 존재하는지 확인
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("삭제할 엔티티를 찾을 수 없습니다: %v", err)
	}

	if err := s.repo.Delete(ctx, entity); err != nil {
		return fmt.Errorf("삭제 실패: %v", err)
	}
	return nil
}
