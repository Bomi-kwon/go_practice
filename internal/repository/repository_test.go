package repository

import (
	"context"
	"errors"
	"go_project/internal/model"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Recorder 모의 객체 정의
type mockRecorder struct {
	mock.Mock
}

func (m *mockRecorder) Insert(ctx context.Context, model *model.Base) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *mockRecorder) Get(ctx context.Context, id uint) (*model.Base, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Base), args.Error(1)
}

func (m *mockRecorder) GetAll(ctx context.Context) ([]*model.Base, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Base), args.Error(1)
}

func (m *mockRecorder) Modify(ctx context.Context, model *model.Base) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *mockRecorder) Remove(ctx context.Context, model *model.Base) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

type RepositoryTestSuite struct {
	suite.Suite
	mockRecorder *mockRecorder
	repo         Repository
}

func (s *RepositoryTestSuite) SetupTest() {
	s.mockRecorder = new(mockRecorder)
	s.repo = NewRepository(s.mockRecorder)
}

func (s *RepositoryTestSuite) TestInsert() {
	tests := []struct {
		name    string
		model   *model.Base
		mockFn  func(*mockRecorder)
		wantErr bool
	}{
		{
			name: "성공_케이스",
			model: &model.Base{
				Name: "테스트_데이터",
			},
			mockFn: func(m *mockRecorder) {
				m.On("Insert", mock.Anything, mock.AnythingOfType("*model.Base")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "실패_케이스",
			model: &model.Base{
				Name: "테스트_데이터",
			},
			mockFn: func(m *mockRecorder) {
				m.On("Insert", mock.Anything, mock.AnythingOfType("*model.Base")).
					Return(errors.New("생성 오류"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockRecorder)

			err := s.repo.Insert(context.Background(), tt.model)

			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *RepositoryTestSuite) TestGet() {
	tests := []struct {
		name    string
		id      uint
		mockFn  func(*mockRecorder)
		want    *model.Base
		wantErr bool
	}{
		{
			name: "성공_케이스",
			id:   1,
			mockFn: func(m *mockRecorder) {
				m.On("Get", mock.Anything, uint(1)).
					Return(&model.Base{ID: 1, Name: "테스트_데이터"}, nil)
			},
			want:    &model.Base{ID: 1, Name: "테스트_데이터"},
			wantErr: false,
		},
		{
			name: "실패_케이스",
			id:   2,
			mockFn: func(m *mockRecorder) {
				m.On("Get", mock.Anything, uint(2)).
					Return((*model.Base)(nil), errors.New("조회 오류"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockRecorder)

			got, err := s.repo.Get(context.Background(), tt.id)

			if tt.wantErr {
				s.Error(err)
				s.Nil(got)
			} else {
				s.NoError(err)
				s.Equal(tt.want, got)
			}
		})
	}
}

func (s *RepositoryTestSuite) TestGetAll() {
	tests := []struct {
		name    string
		mockFn  func(*mockRecorder)
		want    []*model.Base
		wantErr bool
	}{
		{
			name: "성공_케이스",
			mockFn: func(m *mockRecorder) {
				results := []*model.Base{
					{ID: 1, Name: "데이터1"},
					{ID: 2, Name: "데이터2"},
				}
				m.On("GetAll", mock.Anything).Return(results, nil)
			},
			want: []*model.Base{
				{ID: 1, Name: "데이터1"},
				{ID: 2, Name: "데이터2"},
			},
			wantErr: false,
		},
		{
			name: "실패_케이스",
			mockFn: func(m *mockRecorder) {
				m.On("GetAll", mock.Anything).
					Return([]*model.Base{}, errors.New("조회 오류"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockRecorder)

			got, err := s.repo.GetAll(context.Background())

			if tt.wantErr {
				s.Error(err)
				s.Nil(got)
			} else {
				s.NoError(err)
				s.Equal(tt.want, got)
			}
		})
	}
}

func (s *RepositoryTestSuite) TestModify() {
	tests := []struct {
		name    string
		model   *model.Base
		mockFn  func(*mockRecorder)
		wantErr bool
	}{
		{
			name: "성공_케이스",
			model: &model.Base{
				ID:   1,
				Name: "수정된_데이터",
			},
			mockFn: func(m *mockRecorder) {
				m.On("Modify", mock.Anything, mock.AnythingOfType("*model.Base")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "실패_케이스",
			model: &model.Base{
				ID:   2,
				Name: "수정된_데이터",
			},
			mockFn: func(m *mockRecorder) {
				m.On("Modify", mock.Anything, mock.AnythingOfType("*model.Base")).
					Return(errors.New("업데이트 오류"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockRecorder)

			err := s.repo.Modify(context.Background(), tt.model)

			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *RepositoryTestSuite) TestRemove() {
	tests := []struct {
		name    string
		model   *model.Base
		mockFn  func(*mockRecorder)
		wantErr bool
	}{
		{
			name: "성공_케이스",
			model: &model.Base{
				ID:   1,
				Name: "삭제할_데이터",
			},
			mockFn: func(m *mockRecorder) {
				m.On("Remove", mock.Anything, mock.AnythingOfType("*model.Base")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "실패_케이스",
			model: &model.Base{
				ID:   2,
				Name: "삭제할_데이터",
			},
			mockFn: func(m *mockRecorder) {
				m.On("Remove", mock.Anything, mock.AnythingOfType("*model.Base")).
					Return(errors.New("삭제 오류"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockRecorder)

			err := s.repo.Remove(context.Background(), tt.model)

			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
