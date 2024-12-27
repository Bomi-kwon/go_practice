package recorder

import (
	"context"

	"gorm.io/gorm"
)

// Recorder는 DB와 직접 상호작용하는 인터페이스
type Recorder interface {
	Create(ctx context.Context, entity interface{}) error
	FindByID(ctx context.Context, id uint) (interface{}, error)
	Find(ctx context.Context) ([]interface{}, error)
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, entity interface{}) error
}

type recorder struct {
	db *gorm.DB
}

func NewRecorder(db *gorm.DB) Recorder {
	return &recorder{
		db: db,
	}
}

func (r *recorder) Create(ctx context.Context, entity interface{}) error {
	if result := r.db.WithContext(ctx).Create(entity); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *recorder) FindByID(ctx context.Context, id uint) (interface{}, error) {
	var result interface{}
	if result := r.db.WithContext(ctx).First(&result, id); result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

func (r *recorder) Find(ctx context.Context) ([]interface{}, error) {
	var results []interface{}
	if result := r.db.WithContext(ctx).Find(&results); result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}

func (r *recorder) Update(ctx context.Context, entity interface{}) error {
	if result := r.db.WithContext(ctx).Save(entity); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *recorder) Delete(ctx context.Context, entity interface{}) error {
	if result := r.db.WithContext(ctx).Delete(entity); result.Error != nil {
		return result.Error
	}
	return nil
}
