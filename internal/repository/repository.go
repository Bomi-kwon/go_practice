package repository

import (
	"context"
	"go_project/internal/model"
	"go_project/internal/recorder"
)

// Repository 인터페이스는 비즈니스 로직을 위한 데이터 접근 계층
type Repository interface {
	Insert(ctx context.Context, model *model.Base) error
	Get(ctx context.Context, id uint) (*model.Base, error)
	GetAll(ctx context.Context) ([]*model.Base, error)
	Modify(ctx context.Context, model *model.Base) error
	Remove(ctx context.Context, model *model.Base) error
}

type repository struct {
	recorder recorder.Recorder
}

func NewRepository(recorder recorder.Recorder) Repository {
	return &repository{
		recorder: recorder,
	}
}

func (r *repository) Insert(ctx context.Context, model *model.Base) error {
	return r.recorder.Insert(ctx, model)
}

func (r *repository) Get(ctx context.Context, id uint) (*model.Base, error) {
	result, err := r.recorder.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) GetAll(ctx context.Context) ([]*model.Base, error) {
	results, err := r.recorder.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *repository) Modify(ctx context.Context, model *model.Base) error {
	return r.recorder.Modify(ctx, model)
}

func (r *repository) Remove(ctx context.Context, model *model.Base) error {
	return r.recorder.Remove(ctx, model)
}
