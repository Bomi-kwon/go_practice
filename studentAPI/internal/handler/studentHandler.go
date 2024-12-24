package handler

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"studentAPI/internal/model"

	"github.com/gin-gonic/gin"
)

var students map[int]model.Student
var lastId int

func SetupHandler(g *gin.Engine) {
	g.Static("/static", "static")
	g.LoadHTMLGlob("templates/*")

	g.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	g.GET("/students", GetStudentListHandler)
	g.GET("/student/:id", GetStudentHandler)
	g.POST("/student", PostStudentHandler)
	g.DELETE("/student/:id", DeleteStudentHandler)
	g.PUT("/student/:id", UpdateStudentHandler)

	students = make(map[int]model.Student)
	students[0] = model.Student{Id: 0, Name: "짱구", Age: 5, Score: 100}
	students[1] = model.Student{Id: 1, Name: "짱아", Age: 2, Score: 20}
	students[2] = model.Student{Id: 2, Name: "맹구", Age: 5, Score: 50}
	students[3] = model.Student{Id: 3, Name: "훈이", Age: 5, Score: 60}
	students[4] = model.Student{Id: 4, Name: "철수", Age: 6, Score: 90}

	lastId = 4
}

func GetStudentListHandler(c *gin.Context) {
	acceptHeader := c.GetHeader("Accept")

	list := make(model.StudentArr, 0)
	for _, student := range students {
		list = append(list, student)
	}
	sort.Sort(list)

	if strings.Contains(acceptHeader, "text/html") || strings.Contains(acceptHeader, "*/*") {
		c.HTML(http.StatusOK, "students.html", list)
		return
	}

	c.JSON(http.StatusOK, list)
}

func GetStudentHandler(c *gin.Context) {
	id_string := c.Param("id")
	if id_string == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	student, ok := students[id]
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, student)
}

func PostStudentHandler(c *gin.Context) {
	var student model.Student
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
	id_string := c.Param("id")
	if id_string == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	delete(students, id)
	c.String(http.StatusOK, "%d번 학생이 삭제되었습니다.", id)
}

func UpdateStudentHandler(c *gin.Context) {
	id_string := c.Param("id")
	if id_string == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if _, exists := students[id]; !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	student.Id = id
	students[id] = student
	c.String(http.StatusOK, "%d번 학생 정보가 수정되었습니다.", id)
}
