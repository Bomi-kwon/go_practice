package main

import (
	"fmt"
	"net/http"
	"strconv" // 타입 변환
)

func barHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()    // 쿼리 파라미터를 가져오는 함수
	name := values.Get("name") // key가 name인 값을 가져오는 함수
	if name == "" {
		name = "World"
	}
	id, _ := strconv.Atoi(values.Get("id")) // key가 id인 값(string)을 가져와서 정수로 변환
	fmt.Fprintf(w, "Hello, %s! id: %d", name, id)
}

func main() {
	mux := http.NewServeMux() // 라우터 생성
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	mux.HandleFunc("/bar", barHandler)

	http.ListenAndServe(":8080", mux) // mux 인스턴스
}
