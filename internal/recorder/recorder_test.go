package recorder

import (
	"context"
	"errors"
	"fmt"
	"go_project/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "your_host"          // 데이터베이스 서버 주소 (예: localhost 또는 127.0.0.1)
	user     = "your_username"      // PostgreSQL 사용자 이름
	password = "your_password"      // PostgreSQL 사용자 비밀번호
	dbname   = "your_database_name" // 테스트 데이터베이스 이름
	port     = "your_port"          // PostgreSQL 포트 번호 (기본값: 5432)
)

type RecorderTestSuite struct {
	suite.Suite
	db       *gorm.DB
	recorder Recorder
}

func (s *RecorderTestSuite) SetupSuite() {
	// 테스트용 DB 연결 설정
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)

	// 테이블 자동 생성
	err = db.AutoMigrate(&model.Base{})
	s.Require().NoError(err)

	s.db = db
	s.recorder = NewRecorder(db)
}

func (s *RecorderTestSuite) TearDownTest() {
	// 각 테스트 후 테이블 초기화
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Base{})
}

func (s *RecorderTestSuite) TestInsert() {
	// given
	m := &model.Base{
		Name: "테스트_데이터",
	}

	// when
	err := s.recorder.Insert(context.Background(), m)

	// then
	s.NoError(err)
	s.NotZero(m.ID)
	s.NotZero(m.CreatedAt)
	s.NotZero(m.UpdatedAt)

	// DB에서 직접 확인
	var saved model.Base
	err = s.db.First(&saved, m.ID).Error
	s.NoError(err)
	s.Equal(m.Name, saved.Name)
}

func (s *RecorderTestSuite) TestGet() {
	// given
	m := &model.Base{
		Name: "테스트_데이터",
	}
	s.db.Create(m)

	// when
	result, err := s.recorder.Get(context.Background(), m.ID)

	// then
	s.NoError(err)
	s.Equal(m.ID, result.ID)
	s.Equal(m.Name, result.Name)
}

func (s *RecorderTestSuite) TestGet_NotFound() {
	// when
	result, err := s.recorder.Get(context.Background(), 999)

	// then
	s.Error(err)
	s.Nil(result)
}

func (s *RecorderTestSuite) TestGetAll() {
	// given
	ms := []*model.Base{
		{Name: "테스트1"},
		{Name: "테스트2"},
		{Name: "테스트3"},
	}
	for _, e := range ms {
		s.db.Create(e)
	}

	// when
	results, err := s.recorder.GetAll(context.Background())

	// then
	s.NoError(err)
	s.Len(results, len(ms))
	for i, result := range results {
		s.Equal(ms[i].Name, result.Name)
	}
}

func (s *RecorderTestSuite) TestModify() {
	// given
	m := &model.Base{
		Name: "테스트_데이터",
	}
	s.db.Create(m)
	oldUpdatedAt := m.UpdatedAt

	time.Sleep(time.Second) // UpdatedAt 변경 확인을 위해 잠시 대기

	// when
	m.Name = "수정된_데이터"
	err := s.recorder.Modify(context.Background(), m)

	// then
	s.NoError(err)

	// DB에서 직접 확인
	var updated model.Base
	err = s.db.First(&updated, m.ID).Error
	s.NoError(err)
	s.Equal("수정된_데이터", updated.Name)
	s.True(updated.UpdatedAt.After(oldUpdatedAt))
}

func (s *RecorderTestSuite) TestRemove() {
	// given
	m := &model.Base{
		Name: "테스트_데이터",
	}
	s.db.Create(m)

	// when
	err := s.recorder.Remove(context.Background(), m)

	// then
	s.NoError(err)

	// DB에서 직접 확인
	var deleted model.Base
	err = s.db.First(&deleted, m.ID).Error
	s.Error(err) // record not found 에러 기대
	s.True(errors.Is(err, gorm.ErrRecordNotFound))
}

func TestRecorderSuite(t *testing.T) {
	suite.Run(t, new(RecorderTestSuite))
}
