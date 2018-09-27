package main

import (
	"encoding/hex"
	"fmt"
	"goim/im"
	"goim/server"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func maine() {
	http.HandleFunc("/users", im.ViewUserList)
	//http.HandleFunc("/users/view")
	http.HandleFunc("/", im.ViewIndex)

	http.ListenAndServe(":9999", nil)
}

func main() {
	fmt.Println("Ok, started!")
	go server.StartTCPServer()

	err := server.StartWebServer()
	if err != nil {
		fmt.Println(err)
	}
}

func main2() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	bs, e := ioutil.ReadFile("output.txt")
	if e != nil {
		log.Fatal(e)
	}
	//fmt.printl
	//str :=

	content, e2 := hex.DecodeString(string(bs))
	if e2 != nil {
		log.Fatal(e2)
	}

	fmt.Println(string(content))
}
