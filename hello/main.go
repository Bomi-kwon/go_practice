package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
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

func MakeWebHandler() http.Handler { // 핸들러 인스턴스(라우터) 생성
	mux := mux.NewRouter()

	// 1. 학생 리스트 조회
	mux.HandleFunc("/students", GetStudentListHandler).Methods("GET")
	// "/students" 로 들어오는 GET 요청을 받을 때만 GetStudentListHandler 핸들러 동작

	// 2. 학생 상세 조회
	mux.HandleFunc("/students/{id:[0-9]+}", GetStudentHandler).Methods("GET")
	// "/students" 아래 숫자로 된 경로 오면 GetStudentHandler 핸들러 동작
	// gorilla/mux에서 자동으로 id값을 내부 map에 저장

	// 3. 학생 추가
	mux.HandleFunc("/students", PostStudentHandler).Methods("POST")
	// "/students" 로 들어오는 POST 요청을 받을 때만 PostStudentHandler 핸들러 동작

	// 4. 학생 삭제
	mux.HandleFunc("/students/{id:[0-9]+}", DeleteStudentHandler).Methods("DELETE")
	// "/students" 아래 숫자로 된 경로 오면 DeleteStudentHandler 핸들러 동작

	students = make(map[int]Student) // 임시 데이터 생성
	students[0] = Student{1, "짱구", 5, 100}
	students[1] = Student{2, "짱아", 2, 20}
	students[2] = Student{3, "맹구", 5, 50}
	students[3] = Student{4, "훈이", 5, 60}
	students[4] = Student{5, "철수", 6, 90}

	lastId = 5

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

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)               // 저장했던 id값 읽어오기
	id, _ := strconv.Atoi(vars["id"]) // 문자열을 정수로 변환
	student, ok := students[id]       // 해당 id의 학생 정보 읽기
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 해당 id의 학생 정보가 없으면 404 상태 코드 반환
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
}

func PostStudentHandler(w http.ResponseWriter, r *http.Request) {
	student := Student{}
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 요청 본문이 유효하지 않으면 400 상태 코드 반환
		return
	}
	lastId++
	students[lastId] = student        // 맵에 학생 추가
	w.WriteHeader(http.StatusCreated) // 201 상태 코드 반환 (학생 추가 성공)
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if _, ok := students[id]; !ok {
		w.WriteHeader(http.StatusNotFound) // 해당 id의 학생 정보가 없으면 404 상태 코드 반환
		return
	}
	delete(students, id) // go안에 내장되어 있는 맵 속 데이터 삭제 함수
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.ListenAndServe(":8080", MakeWebHandler())
}
