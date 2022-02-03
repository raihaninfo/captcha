package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/send", SendHandler)

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("index.gohtml")
	if err != nil {
		log.Println(err)
	}
	tpl.Execute(w, nil)
}

func SendHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	captcha := r.PostFormValue("g-recaptcha-response")
	fmt.Println(captcha)

	valid := CheckGoogleCaptcha(captcha)
	fmt.Println(valid)
	if valid {
		fmt.Fprintf(w, "The captcha was correct!")
	} else {
		fmt.Fprintf(w, "This captcha was NOT correct")
	}
}

func CheckGoogleCaptcha(response string) bool {
	var googleCaptchav3 string = "6LcYhFUeAAAAAPNsgv5l9reKZj8qNT-zB9nqrNzY"
	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		log.Println(err)
	}
	q := req.URL.Query()
	q.Add("secret", googleCaptchav3)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	var googleResponse map[string]interface{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &googleResponse)
	return googleResponse["success"].(bool)

}
