package main

import (
	"fmt"
	"im/im"
)

func chat() {

}

func main() {

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
