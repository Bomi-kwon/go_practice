package main

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
)

type Student struct {
	Id    int
	Name  string
	Age   int
	Score int
}

var students map[int]Student

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

func MakeWebHandler() http.Handler { // 핸들러 인스턴스(라우터) 생성
	mux := mux.NewRouter()
	mux.HandleFunc("/students", GetStudentListHandler).Methods("GET")
	// "/students" 로 들어오는 GET 요청을 받을 때만 GetStudentListHandler 핸들러 동작

	students = make(map[int]Student) // 임시 데이터 생성
	students[0] = Student{1, "짱구", 5, 100}
	students[1] = Student{2, "짱아", 2, 20}
	students[2] = Student{3, "맹구", 5, 50}
	students[3] = Student{4, "훈이", 5, 60}
	students[4] = Student{5, "철수", 6, 90}

	return mux
}

func GetStudentListHandler(w http.ResponseWriter, r *http.Request) {
	list := make(student_arr, 0)
	for _, student := range students {
		list = append(list, student) // 학생 목록을 슬라이스에 복사
	}
	sort.Sort(list) // map은 정렬 불가능하고 슬라이스만 가능해서..
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func main() {
	http.ListenAndServe(":8080", MakeWebHandler())
}
