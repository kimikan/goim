package server

import (
	"goim/im"
	"html/template"
	"io/ioutil"
	"net/http"
)

func returnResult(w http.ResponseWriter, result string) {
	content, err := ioutil.ReadFile("static/html/result.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	t, _ := template.New("webpage").Parse(string(content))
	err = t.Execute(w, result)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	key := r.PostFormValue("key")
	displayname := r.PostFormValue("displayname")
	description := r.PostFormValue("description")
	if len(key) <= 0 || len(key) > 1000 ||
		len(displayname) <= 0 || len(displayname) > 10 {
		returnResult(w, "error input !")
		return
	}

	info := &im.UserInfo{
		Key:         key,
		DisplayName: displayname,
		Description: description,
	}
	err := im.SetUserInfo(info)
	if err != nil {
		returnResult(w, "Failed "+err.Error())
		return
	}
	returnResult(w, "Ok, register success!")
}

func index(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("static/html/reg.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	t, _ := template.New("webpage").Parse(string(content))
	err = t.Execute(w, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func checkDetails(w http.ResponseWriter, r *http.Request) {

}

//webserver entrance
func StartWebServer() error {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/register", register)
	http.HandleFunc("/check", checkDetails)
	http.HandleFunc("/", index)

	return http.ListenAndServe(":9998", nil)
}
