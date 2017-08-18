package main

import (
	"fmt"
	"im/im"
	"net/http"
)

func chat() {

}

func main() {
	http.HandleFunc("/users", im.ViewUserList)
	//http.HandleFunc("/users/view")
	http.HandleFunc("/", im.ViewIndex)

	http.ListenAndServe(":9999", nil)
}

func main2() {

	u := im.NewUser("kan", "xxxxx")

	if u != nil {
		bs, err := u.ToBytes()

		if err == nil {
			fmt.Println(bs)

			fmt.Println(im.ParseUser(bs))
		}
	}
	fmt.Println("Service started")
}
