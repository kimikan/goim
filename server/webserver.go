package server

import "net/http"

func test(w http.ResponseWriter, r *http.Request) {

}

//webserver entrance
func StartWebServer() error {
	//http.HandleFunc("/users", im.ViewUserList)
	//http.HandleFunc("/users/view")
	http.HandleFunc("/", test)

	return http.ListenAndServe(":9998", nil)
}
