package usecase

import (
	"context"
	"fmt"
	"go_project/internal/model"
	"go_project/internal/repository"
)

type Usecase interface {
	Insert(ctx context.Context, model *model.Base) error
	Get(ctx context.Context, id uint) (*model.Base, error)
	GetAll(ctx context.Context) ([]*model.Base, error)
	Modify(ctx context.Context, id uint, model *model.Base) error
	Remove(ctx context.Context, id uint) error
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

// 기본 CRUD 구현
func (u *usecase) Insert(ctx context.Context, model *model.Base) error {
	if err := u.repo.Insert(ctx, model); err != nil {
		return fmt.Errorf("생성 실패: %v", err)
	}
	return nil
}

func (u *usecase) Get(ctx context.Context, id uint) (*model.Base, error) {
	result, err := u.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("조회 실패: %v", err)
	}
	return result, nil
}

func (u *usecase) GetAll(ctx context.Context) ([]*model.Base, error) {
	results, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("목록 조회 실패: %v", err)
	}
	return results, nil
}

func (u *usecase) Modify(ctx context.Context, id uint, model *model.Base) error {
	// 먼저 존재하는지 확인
	_, err := u.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("업데이트할 모델을 찾을 수 없습니다: %v", err)
	}

	if err := u.repo.Modify(ctx, model); err != nil {
		return fmt.Errorf("업데이트 실패: %v", err)
	}
	return nil
}

func (u *usecase) Remove(ctx context.Context, id uint) error {
	// 먼저 존재하는지 확인
	model, err := u.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("삭제할 모델을 찾을 수 없습니다: %v", err)
	}

	if err := u.repo.Remove(ctx, model); err != nil {
		return fmt.Errorf("삭제 실패: %v", err)
	}
	return nil
}
