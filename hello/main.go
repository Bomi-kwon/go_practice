package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Student struct {
	Id    int
	Name  string
	Age   int
	Score int
}

var students map[int]Student

var lastId int

type student_arr []Student

// Sort 인터페이스 구현을 위한 3개의 메서드
func (s student_arr) Len() int {
	return len(s)
}

func (s student_arr) Less(i, j int) bool {
	return s[i].Name < s[j].Name // 이름 기준으로 오름차순 정렬
}

func (s student_arr) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func SetupHandler(g *gin.Engine) { // 핸들러 인스턴스(라우터) 생성
	g.GET("/students", GetStudentListHandler)
	g.GET("/student/:id", GetStudentHandler)
	g.POST("/student", PostStudentHandler)
	g.DELETE("/student/:id", DeleteStudentHandler)

	students = make(map[int]Student) // 임시 데이터 생성
	students[0] = Student{1, "짱구", 5, 100}
	students[1] = Student{2, "짱아", 2, 20}
	students[2] = Student{3, "맹구", 5, 50}
	students[3] = Student{4, "훈이", 5, 60}
	students[4] = Student{5, "철수", 6, 90}

	lastId = 5
}

func GetStudentListHandler(c *gin.Context) {
	list := make(student_arr, 0)
	for _, student := range students {
		list = append(list, student)
	}
	c.JSON(http.StatusOK, list) // 학생 목록을 JSON 형식으로 반환
}

func GetStudentHandler(c *gin.Context) {
	id_string := c.Param("id") // 파라미터 추출
	if id_string == "" {
		c.AbortWithStatus(http.StatusBadRequest) // 파라미터가 없으면 400 상태 코드 반환
		return
	}
	id, err := strconv.Atoi(id_string) // 파라미터를 정수로 변환
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	student, ok := students[id] // 해당 id를 가진 학생 찾기
	if !ok {
		c.AbortWithStatus(http.StatusNotFound) // 학생이 없으면 404 상태 코드 반환
		return
	}
	c.JSON(http.StatusOK, student) // 학생 정보를 JSON 형식으로 반환
}

func PostStudentHandler(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	lastId++
	student.Id = lastId
	students[lastId] = student
	c.String(http.StatusCreated, "%d번 학생이 추가되었습니다.", lastId)
}

func DeleteStudentHandler(c *gin.Context) {
	id_string := c.Param("id") // 파라미터 추출
	if id_string == "" {
		c.AbortWithStatus(http.StatusBadRequest) // 파라미터가 없으면 400 상태 코드 반환
		return
	}
	id, err := strconv.Atoi(id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error()) // 파라미터가 정수가 아니면 400 상태 코드 반환
		return
	}
	delete(students, id)
	c.String(http.StatusOK, "%d번 학생이 삭제되었습니다.", id) // 학생 삭제 메시지 반환
}

func main() {
	r := gin.Default()
	SetupHandler(r)
	r.Run(":8080")
}
