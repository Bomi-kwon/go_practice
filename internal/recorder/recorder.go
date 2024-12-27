package recorder

import (
	"context"
	"go_project/internal/model"

	"gorm.io/gorm"
)

// Recorder는 DB와 직접 상호작용하는 인터페이스
type Recorder interface {
	Create(ctx context.Context, entity interface{}) error
	Get(ctx context.Context, id uint) (interface{}, error)
	List(ctx context.Context) ([]interface{}, error)
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, entity interface{}) error
}

type recorder struct {
	db *gorm.DB
}

func NewRecorder(db *gorm.DB) Recorder {
	return &recorder{db: db}
}

func (r *recorder) Create(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *recorder) Get(ctx context.Context, id uint) (interface{}, error) {
	var base model.Base
	if err := r.db.WithContext(ctx).First(&base, id).Error; err != nil {
		return nil, err
	}
	return base, nil
}

func (r *recorder) List(ctx context.Context) ([]interface{}, error) {
	var bases []model.Base
	if err := r.db.WithContext(ctx).Find(&bases).Error; err != nil {
		return nil, err
	}

	// []model.Base를 []interface{}로 변환
	result := make([]interface{}, len(bases))
	for i, v := range bases {
		result[i] = v
	}
	return result, nil
}

func (r *recorder) Update(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *recorder) Delete(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Delete(entity).Error
}
