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
		log.Fatal("ListenAndServe: ", err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("index.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	tpl.Execute(w, nil)
}

func SendHandler(w http.ResponseWriter, r *http.Request) {
	captcha := r.PostFormValue("h-captcha-response")
	fmt.Println(captcha)

	valid := CheckHCaptcha(captcha)
	fmt.Println(valid)
	if valid {
		fmt.Fprintf(w, "The captcha was correct!")
	} else {
		fmt.Fprintf(w, "This captcha was NOT correct")
	}
}

func CheckHCaptcha(response string) bool {
	var hCaptcha string = "0x055a7D2c94f486D36e2d0B3110F72BB50d4cf52B"
	req, _ := http.NewRequest("POST", "https://hcaptcha.com/siteverify", nil)
	q := req.URL.Query()
	q.Add("secret", hCaptcha)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	var hCaptchaResponse map[string]interface{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &hCaptchaResponse)
	return hCaptchaResponse["success"].(bool)

}
