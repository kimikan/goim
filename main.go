package main

import (
	"fmt"
	"goim/im"
	"goim/server"
	"net/http"
)

func maine() {
	http.HandleFunc("/users", im.ViewUserList)
	//http.HandleFunc("/users/view")
	http.HandleFunc("/", im.ViewIndex)

	http.ListenAndServe(":9999", nil)
}

func main() {
	go server.StartTCPServer()

	err := server.StartWebServer()
	if err != nil {
		fmt.Println(err)
	}
}
