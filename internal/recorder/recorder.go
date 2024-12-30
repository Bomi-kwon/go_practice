package recorder

import (
	"context"
	"go_project/internal/model"

	"gorm.io/gorm"
)

// Recorder는 DB와 직접 상호작용하는 인터페이스
type Recorder interface {
	Insert(ctx context.Context, model *model.Base) error
	Get(ctx context.Context, id uint) (*model.Base, error)
	GetAll(ctx context.Context) ([]*model.Base, error)
	Modify(ctx context.Context, model *model.Base) error
	Remove(ctx context.Context, model *model.Base) error
}

type recorder struct {
	db *gorm.DB
}

func NewRecorder(db *gorm.DB) Recorder {
	return &recorder{
		db: db,
	}
}

func (r *recorder) Insert(ctx context.Context, model *model.Base) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *recorder) Get(ctx context.Context, id uint) (*model.Base, error) {
	var base model.Base
	if err := r.db.WithContext(ctx).First(&base, id).Error; err != nil {
		return nil, err
	}
	return &base, nil
}

func (r *recorder) GetAll(ctx context.Context) ([]*model.Base, error) {
	var bases []*model.Base
	if err := r.db.WithContext(ctx).Find(&bases).Error; err != nil {
		return nil, err
	}
	return bases, nil
}

func (r *recorder) Modify(ctx context.Context, model *model.Base) error {
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *recorder) Remove(ctx context.Context, model *model.Base) error {
	return r.db.WithContext(ctx).Delete(model).Error
}
