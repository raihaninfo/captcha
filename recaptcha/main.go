package main
 
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)
 
func main() {
 
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/send", SendHandler)
 
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
 
}
 
func IndexHandler(w http.ResponseWriter, r *http.Request) {
 
	form := `<script src='https://www.google.com/recaptcha/api.js?render=6LfyG1IeAAAAAPinKen3-8jPmvINU0ec1LdTl_tm'></script>
	<form action="/send" method="POST">
	<div class="g-recaptcha" data-sitekey="PUBLIC_KEY_HERE"></div>
	<input type="submit">`
 
	fmt.Fprintf(w, form)
}
 
func SendHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	captcha := r.PostFormValue("g-recaptcha-response")
 
	valid := CheckGoogleCaptcha(captcha)
 
	if valid {
		fmt.Fprintf(w, "The captcha was correct!")
	} else {
		fmt.Fprintf(w, "This captcha was NOT correct, check the public and secret keys.")
	}
}
 
func CheckGoogleCaptcha(response string) bool {
	var googleCaptcha string = "SECRET_KEY_HERE"
	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	q := req.URL.Query()
	q.Add("secret", googleCaptcha)
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