package usecase

import (
	"context"
	"errors"
	"go_project/internal/model"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Repository 모의 객체 정의
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Insert(ctx context.Context, model *model.Base) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *mockRepository) Get(ctx context.Context, id uint) (*model.Base, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Base), args.Error(1)
}

func (m *mockRepository) GetAll(ctx context.Context) ([]*model.Base, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Base), args.Error(1)
}

func (m *mockRepository) Modify(ctx context.Context, model *model.Base) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *mockRepository) Remove(ctx context.Context, model *model.Base) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

// 관련된 테스트를 하나의 Suite로 묶어서 관리
type UsecaseTestSuite struct {
	suite.Suite
	mockRepo *mockRepository
	uc       Usecase
}

func (s *UsecaseTestSuite) SetupTest() {
	// 테스트 초기화
	s.mockRepo = new(mockRepository)
	s.uc = NewUsecase(s.mockRepo)
}

func (s *UsecaseTestSuite) TearDownTest() {
	// 테스트 종료 시 mock 객체의 함수를 실제로 호출했는지 검증
	s.mockRepo.AssertExpectations(s.T())
}

func (s *UsecaseTestSuite) TestInsert() {
	tests := []struct {
		name    string
		model   *model.Base
		mockFn  func(*mockRepository)
		wantErr bool
	}{
		{
			name: "성공_케이스",
			model: &model.Base{
				Name: "추가할_데이터",
			},
			mockFn: func(m *mockRepository) {
				m.On("Insert", mock.Anything, mock.AnythingOfType("*model.Base")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "실패_케이스",
			model: &model.Base{
				Name: "추가할_데이터",
			},
			mockFn: func(m *mockRepository) {
				m.On("Insert", mock.Anything, mock.AnythingOfType("*model.Base")).
					Return(errors.New("생성 오류"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		// 반복문을 순회하며 테스트 케이스 실행
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockRepo)

			err := s.uc.Insert(context.Background(), tt.model) // 테스트할 메서드 호출

			if tt.wantErr {
				// 실패 케이스
				s.Error(err)
			} else {
				// 성공 케이스
				s.NoError(err)
				s.mockRepo.AssertCalled(s.T(), "Insert", mock.Anything, mock.MatchedBy(func(model *model.Base) bool {
					return model.Name == "추가할_데이터"
				}))
			}
			s.TearDownTest()
		})
	}
}

func (s *UsecaseTestSuite) TestGet() {
	tests := []struct {
		name    string
		id      uint
		mockFn  func(*mockRepository)
		want    *model.Base
		wantErr bool
	}{
		{
			name: "성공_케이스",
			id:   1,
			mockFn: func(m *mockRepository) {
				m.On("Get", mock.Anything, uint(1)).
					Return(&model.Base{ID: 1, Name: "조회할_데이터"}, nil)
			},
			want:    &model.Base{ID: 1, Name: "조회할_데이터"},
			wantErr: false,
		},
		{
			name: "실패_케이스",
			id:   2,
			mockFn: func(m *mockRepository) {
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
			tt.mockFn(s.mockRepo)

			got, err := s.uc.Get(context.Background(), tt.id) // 테스트할 메서드 호출

			if tt.wantErr {
				// 실패 케이스
				s.Error(err)
				s.Nil(got)
			} else {
				// 성공 케이스
				s.NoError(err)
				s.Equal(tt.want, got) // 원하는 값과 실제로 리턴된 값이 같은지 검증
			}
			s.TearDownTest()
		})
	}
}

func (s *UsecaseTestSuite) TestModify() {
	tests := []struct {
		name    string
		id      uint
		model   *model.Base
		mockFn  func(*mockRepository)
		wantErr bool
	}{
		{
			name: "성공_케이스",
			id:   1,
			model: &model.Base{
				ID:   1,
				Name: "수정할_데이터",
			},
			mockFn: func(m *mockRepository) {
				m.On("Get", mock.Anything, uint(1)).
					Return(&model.Base{ID: 1, Name: "기존_데이터"}, nil)
				m.On("Modify", mock.Anything, mock.AnythingOfType("*model.Base")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "실패_케이스_엔티티_없음",
			id:   2,
			model: &model.Base{
				ID:   2,
				Name: "수정할_데이터",
			},
			mockFn: func(m *mockRepository) {
				m.On("Get", mock.Anything, uint(2)).
					Return((*model.Base)(nil), errors.New("엔티티 없음"))
			},
			wantErr: true,
		},
		{
			name: "실패_케이스_업데이트_오류",
			id:   3,
			model: &model.Base{
				ID:   3,
				Name: "수정할_데이터",
			},
			mockFn: func(m *mockRepository) {
				m.On("Get", mock.Anything, uint(3)).
					Return(&model.Base{ID: 3, Name: "기존_데이터"}, nil)
				m.On("Modify", mock.Anything, mock.AnythingOfType("*model.Base")).
					Return(errors.New("업데이트 실패"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockRepo)

			err := s.uc.Modify(context.Background(), tt.id, tt.model) // 테스트할 메서드 호출

			if tt.wantErr {
				// 실패 케이스
				s.Error(err)
			} else {
				// 성공 케이스
				s.NoError(err)

				// Get 호출 검증
				s.mockRepo.AssertCalled(s.T(), "Get", mock.Anything, uint(1))

				// Modify 호출 검증
				s.mockRepo.AssertCalled(s.T(), "Modify", mock.Anything, mock.MatchedBy(func(model *model.Base) bool {
					return model.ID == uint(1) && model.Name == "수정할_데이터"
				}))
			}
			s.TearDownTest()
		})
	}
}

func (s *UsecaseTestSuite) TestRemove() {
	tests := []struct {
		name    string
		id      uint
		mockFn  func(*mockRepository)
		wantErr bool
	}{
		{
			name: "성공_케이스",
			id:   1,
			mockFn: func(m *mockRepository) {
				model := &model.Base{ID: 1, Name: "삭제할_데이터"}
				m.On("Get", mock.Anything, uint(1)).Return(model, nil)
				m.On("Remove", mock.Anything, mock.AnythingOfType("*model.Base")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "실패_케이스_엔티티_없음",
			id:   2,
			mockFn: func(m *mockRepository) {
				m.On("Get", mock.Anything, uint(2)).
					Return((*model.Base)(nil), errors.New("엔티티 없음"))
			},
			wantErr: true,
		},
		{
			name: "실패_케이스_삭제_오류",
			id:   3,
			mockFn: func(m *mockRepository) {
				model := &model.Base{ID: 3, Name: "삭제할_데이터"}
				m.On("Get", mock.Anything, uint(3)).Return(model, nil)
				m.On("Remove", mock.Anything, mock.AnythingOfType("*model.Base")).
					Return(errors.New("삭제 실패"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockRepo)

			err := s.uc.Remove(context.Background(), tt.id) // 테스트할 메서드 호출

			if tt.wantErr {
				// 실패 케이스
				s.Error(err)
			} else {
				// 성공 케이스
				s.NoError(err)

				// Get 호출 검증
				s.mockRepo.AssertCalled(s.T(), "Get", mock.Anything, uint(1))

				// Remove 호출 검증
				s.mockRepo.AssertCalled(s.T(), "Remove", mock.Anything, mock.MatchedBy(func(model *model.Base) bool {
					return model.ID == uint(1) && model.Name == "삭제할_데이터"
				}))
			}
			s.TearDownTest()
		})
	}
}

func (s *UsecaseTestSuite) TestGetAll() {
	tests := []struct {
		name    string
		mockFn  func(*mockRepository)
		want    []*model.Base
		wantErr bool
	}{
		{
			name: "성공_케이스",
			mockFn: func(m *mockRepository) {
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
			mockFn: func(m *mockRepository) {
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
			tt.mockFn(s.mockRepo)

			got, err := s.uc.GetAll(context.Background()) // 테스트할 메서드 호출

			if tt.wantErr {
				// 실패 케이스
				s.Error(err)
				s.Nil(got)
			} else {
				// 성공 케이스
				s.NoError(err)
				s.Equal(tt.want, got) // 원하는 값과 실제로 리턴된 값이 같은지 검증
			}
			s.TearDownTest()
		})
	}
}

// 테스트 실행을 위한 엔트리 포인트
func TestUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}
