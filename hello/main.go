package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	// http://localhost:8080/파일이름.png 입력시 화면에 이미지 출력됨
	http.ListenAndServe(":8080", nil) // mux 인스턴스
}
