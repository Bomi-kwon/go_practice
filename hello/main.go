package main

import (
	"encoding/json"
	"net/http"
)

type Student struct {
	Id    int
	Name  string
	Age   int
	Score int
}

func MakeWebHandler() http.Handler { // 핸들러 인스턴스(라우터) 생성
	mux := http.NewServeMux()
	mux.HandleFunc("/student", StudentHandler) // 핸들러 함수 등록
	return mux
}

func StudentHandler(w http.ResponseWriter, r *http.Request) {
	var student = Student{1668002, "권보미", 29, 85}
	json, _ := json.Marshal(student) // 구조체를 []byte 타입으로 변환
	// json은 이제 다음과 같은 형태의 []byte가 됩니다:
	// {"Id":1668002,"Name":"권보미","Age":29,"Score":85}
	w.Header().Add("content-type", "application/json") // 응답이 JSON 형식임을 표시
	w.WriteHeader(http.StatusOK)                       // 200 OK 상태 코드 전송
	w.Write(json)
}

func main() {
	http.ListenAndServe(":8080", MakeWebHandler())
}
