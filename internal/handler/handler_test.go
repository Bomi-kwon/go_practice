package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go_project/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// 응답 구조체 추가
type response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Usecase 모의 객체 정의
type mockUsecase struct {
	mock.Mock
}

func (m *mockUsecase) Insert(ctx context.Context, model *model.Base) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *mockUsecase) Get(ctx context.Context, id uint) (*model.Base, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Base), args.Error(1)
}

func (m *mockUsecase) GetAll(ctx context.Context) ([]*model.Base, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Base), args.Error(1)
}

func (m *mockUsecase) Modify(ctx context.Context, id uint, model *model.Base) error {
	args := m.Called(ctx, id, model)
	return args.Error(0)
}

func (m *mockUsecase) Remove(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type HandlerTestSuite struct {
	suite.Suite
	mockUc  *mockUsecase
	handler *Handler
}

func (s *HandlerTestSuite) SetupTest() {
	s.mockUc = new(mockUsecase)
	s.handler = NewHandler(s.mockUc)
}

func (s *HandlerTestSuite) setupRouter() *gin.Engine {
	// gin을 테스트 모드로 설정
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// API 라우트만 등록
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/resources", s.handler.GetAll)
			v1.GET("/resources/:id", s.handler.Get)
			v1.POST("/resources", s.handler.Insert)
			v1.PUT("/resources/:id", s.handler.Modify)
			v1.DELETE("/resources/:id", s.handler.Remove)
		}
	}

	return router
}

func (s *HandlerTestSuite) TestInsert() {
	tests := []struct {
		name   string
		input  *model.Base
		mockFn func(*mockUsecase)
		want   *response
	}{
		{
			name: "성공_케이스",
			input: &model.Base{
				Name: "테스트_데이터",
			},
			mockFn: func(m *mockUsecase) {
				m.On("Insert", mock.Anything, mock.AnythingOfType("*model.Base")).Return(nil)
			},
			want: &response{
				Status:  http.StatusCreated,
				Message: "성공",
				Data:    &model.Base{Name: "테스트_데이터"},
			},
		},
		{
			name: "실패_케이스",
			input: &model.Base{
				Name: "테스트_데이터",
			},
			mockFn: func(m *mockUsecase) {
				m.On("Insert", mock.Anything, mock.AnythingOfType("*model.Base")).
					Return(errors.New("생성 오류"))
			},
			want: &response{
				Status:  http.StatusInternalServerError,
				Message: "리소스 생성 실패",
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockUc)

			router := s.setupRouter()

			// HTTP 요청 생성
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/resources", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			s.Equal(tt.want.Status, w.Code)
		})
	}
}

func (s *HandlerTestSuite) TestGet() {
	tests := []struct {
		name   string
		id     string
		mockFn func(*mockUsecase)
		want   *response
	}{
		{
			name: "성공_케이스",
			id:   "1",
			mockFn: func(m *mockUsecase) {
				m.On("Get", mock.Anything, uint(1)).
					Return(&model.Base{ID: 1, Name: "테스트_데이터"}, nil)
			},
			want: &response{
				Status:  http.StatusOK,
				Message: "성공",
				Data:    &model.Base{ID: 1, Name: "테스트_데이터"},
			},
		},
		{
			name: "실패_케이스_없는_데이터",
			id:   "999",
			mockFn: func(m *mockUsecase) {
				m.On("Get", mock.Anything, uint(999)).
					Return((*model.Base)(nil), errors.New("데이터 없음"))
			},
			want: &response{
				Status:  http.StatusInternalServerError,
				Message: "리소스 조회 실패",
				Data:    nil,
			},
		},
		{
			name: "실패_케이스_잘못된_ID",
			id:   "invalid",
			mockFn: func(m *mockUsecase) {
				// mock 호출 없음 - ID 파싱 실패
			},
			want: &response{
				Status:  http.StatusBadRequest,
				Message: "잘못된 ID 형식",
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockUc)

			router := s.setupRouter()

			// HTTP 요청 생성
			req := httptest.NewRequest(http.MethodGet, "/api/v1/resources/"+tt.id, nil)
			w := httptest.NewRecorder()

			// 라우터를 통한 요청 처리
			router.ServeHTTP(w, req)

			// 응답 검증
			s.Equal(tt.want.Status, w.Code)
			if tt.want != nil {
				var got response
				err := json.NewDecoder(w.Body).Decode(&got)
				s.NoError(err)
				s.Equal(tt.want.Status, got.Status)
				s.Equal(tt.want.Message, got.Message)
				if got.Data != nil {
					gotData := &model.Base{}
					dataBytes, _ := json.Marshal(got.Data)
					json.Unmarshal(dataBytes, gotData)
					wantData := tt.want.Data.(*model.Base)
					s.Equal(wantData.ID, gotData.ID)
					s.Equal(wantData.Name, gotData.Name)
				}
			}
		})
	}
}

func (s *HandlerTestSuite) TestGetAll() {
	tests := []struct {
		name   string
		mockFn func(*mockUsecase)
		want   *response
	}{
		{
			name: "성공_케이스",
			mockFn: func(m *mockUsecase) {
				results := []*model.Base{
					{ID: 1, Name: "데이터1"},
					{ID: 2, Name: "데이터2"},
				}
				m.On("GetAll", mock.Anything).Return(results, nil)
			},
			want: &response{
				Status:  http.StatusOK,
				Message: "성공",
				Data: []*model.Base{
					{ID: 1, Name: "데이터1"},
					{ID: 2, Name: "데이터2"},
				},
			},
		},
		{
			name: "실패_케이스",
			mockFn: func(m *mockUsecase) {
				m.On("GetAll", mock.Anything).
					Return([]*model.Base{}, errors.New("조회 오류"))
			},
			want: &response{
				Status:  http.StatusInternalServerError,
				Message: "리소스 목록 조회 실패",
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockUc)

			router := s.setupRouter()

			// HTTP 요청 생성
			req := httptest.NewRequest(http.MethodGet, "/api/v1/resources", nil)
			w := httptest.NewRecorder()

			// 라우터를 통한 요청 처리
			router.ServeHTTP(w, req)

			// 응답 검증
			s.Equal(tt.want.Status, w.Code)
			if tt.want != nil {
				var got response
				err := json.NewDecoder(w.Body).Decode(&got)
				s.NoError(err)
				s.Equal(tt.want.Status, got.Status)
				s.Equal(tt.want.Message, got.Message)

				if got.Data != nil {
					var gotData []*model.Base
					dataBytes, _ := json.Marshal(got.Data)
					json.Unmarshal(dataBytes, &gotData)
					wantData := tt.want.Data.([]*model.Base)
					s.Equal(len(wantData), len(gotData))
					for i := range wantData {
						s.Equal(wantData[i].ID, gotData[i].ID)
						s.Equal(wantData[i].Name, gotData[i].Name)
					}
				}
			}
		})
	}
}

func (s *HandlerTestSuite) TestModify() {
	tests := []struct {
		name   string
		id     string
		input  *model.Base
		mockFn func(*mockUsecase)
		want   *response
	}{
		{
			name: "성공_케이스",
			id:   "1",
			input: &model.Base{
				ID:   1,
				Name: "수정된_데이터",
			},
			mockFn: func(m *mockUsecase) {
				m.On("Modify", mock.Anything, uint(1), mock.Anything).Return(nil)
			},
			want: &response{
				Status:  http.StatusOK,
				Message: "성공",
				Data:    &model.Base{ID: 1, Name: "수정된_데이터"},
			},
		},
		{
			name: "실패_케이스_없는_데이터",
			id:   "999",
			input: &model.Base{
				ID:   999,
				Name: "수정된_데이터",
			},
			mockFn: func(m *mockUsecase) {
				m.On("Modify", mock.Anything, uint(999), mock.Anything).
					Return(errors.New("데이터 없음"))
			},
			want: &response{
				Status:  http.StatusInternalServerError,
				Message: "리소스 수정 실패",
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockUc)

			router := s.setupRouter()

			// HTTP 요청 생성
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/resources/"+tt.id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 라우터를 통한 요청 처리
			router.ServeHTTP(w, req)

			// 응답 검증
			s.Equal(tt.want.Status, w.Code)
		})
	}
}

func (s *HandlerTestSuite) TestRemove() {
	tests := []struct {
		name   string
		id     string
		mockFn func(*mockUsecase)
		want   *response
	}{
		{
			name: "성공_케이스",
			id:   "1",
			mockFn: func(m *mockUsecase) {
				m.On("Remove", mock.Anything, uint(1)).Return(nil)
			},
			want: &response{
				Status:  http.StatusOK,
				Message: "성공",
				Data:    nil,
			},
		},
		{
			name: "실패_케이스_없는_데이터",
			id:   "999",
			mockFn: func(m *mockUsecase) {
				m.On("Remove", mock.Anything, uint(999)).
					Return(errors.New("데이터 없음"))
			},
			want: &response{
				Status:  http.StatusInternalServerError,
				Message: "리소스 삭제 실패",
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			tt.mockFn(s.mockUc)

			router := s.setupRouter()

			// HTTP 요청 생성
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/resources/"+tt.id, nil)
			w := httptest.NewRecorder()

			// 라우터를 통한 요청 처리
			router.ServeHTTP(w, req)

			// 응답 검증
			s.Equal(tt.want.Status, w.Code)
		})
	}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
