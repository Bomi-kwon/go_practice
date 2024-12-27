package repository

import (
	"context"
	"go_project/internal/recorder"
)

// Repository 인터페이스는 비즈니스 로직을 위한 데이터 접근 계층
type Repository interface {
	Create(ctx context.Context, entity interface{}) error
	Get(ctx context.Context, id uint) (interface{}, error)
	List(ctx context.Context) ([]interface{}, error)
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, entity interface{}) error
}

type repository struct {
	recorder recorder.Recorder
}

func NewRepository(recorder recorder.Recorder) Repository {
	return &repository{
		recorder: recorder,
	}
}

func (r *repository) Create(ctx context.Context, entity interface{}) error {
	if err := r.recorder.Create(ctx, entity); err != nil {
		return err
	}
	return nil
}

func (r *repository) Get(ctx context.Context, id uint) (interface{}, error) {
	result, err := r.recorder.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) List(ctx context.Context) ([]interface{}, error) {
	results, err := r.recorder.List(ctx)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *repository) Update(ctx context.Context, entity interface{}) error {
	if err := r.recorder.Update(ctx, entity); err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, entity interface{}) error {
	if err := r.recorder.Delete(ctx, entity); err != nil {
		return err
	}
	return nil
}
