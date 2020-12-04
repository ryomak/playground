package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(3 * time.Second)
		fmt.Fprintf(w, "Welcome to the home page!")
	}) // ハンドラを登録してウェブページを表示させる
	http.ListenAndServe(":8081", nil)
}
