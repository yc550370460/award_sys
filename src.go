package main

import (
	"./settings/jsonParse"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"html/template"
)

var settingVar = settings.JsonLoad("settings/config.json")
var staticPath = settingVar.Static
var templatePath = settingVar.Template
var bubbleStatus bool

func bubbleTest(ws *websocket.Conn) {
	var err error
	var count_msg int
	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		bubbleStatus = true
		fmt.Println(bubbleStatus)
		fmt.Println("Received back from client: " + reply)
		var number string
		var numbern int

		numbern = rand.Intn(10)
		fmt.Println(numbern)
		number = strconv.Itoa(numbern)
		fmt.Println(number)
		msg := reply + number

		fmt.Println("Sending to client: " + msg + number)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
		if count_msg%10 == 0{
			if err = websocket.Message.Send(ws, "too many messages"); err != nil {
				fmt.Println("Can't send")
				break
			}
		}
		count_msg = count_msg+1
	}
}

var staticReg = regexp.MustCompile("static")
var indexReg = regexp.MustCompile("index")

func route(w http.ResponseWriter, r *http.Request) {
	switch {
	case staticReg.MatchString(r.URL.Path):
		http.FileServer(http.Dir(staticPath))
	case indexReg.MatchString(r.URL.Path):
		http.FileServer(http.Dir(templatePath))
	default:
		w.Write([]byte("位置匹配项"))
	}
}

func award(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles(templatePath + "/award.gtpl")
		log.Println(t.Execute(w, nil))
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles(templatePath + "/login.gtpl")
		//employeeJsonPath := staticPath + "employee.json"
		//employVlue := settings.JsonLoad(employeeJsonPath)


		log.Println(t.Execute(w, nil))
	}
}

func bubble(ws *websocket.Conn) {
	var err error
	var count int
	for {
		if bubbleStatus{
			var msg string
			if count%2 == 0{
				msg = "不知道这个可以输入多长呢呢呢呢呢呢呢呢呢呢呢呢呢呢呢呢呢呢呢"
			}else{
				msg = "2018"
			}
			count = count + 1
			fmt.Println("Sending to client: " + msg)

			if err = websocket.Message.Send(ws, msg); err != nil {
				fmt.Println("Can't send")
				break
			}
			bubbleStatus = false
		}
	}
}


func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	http.HandleFunc("/award/", award)
	http.HandleFunc("/login/", login)
	http.Handle("/web_socket", websocket.Handler(bubbleTest))
	http.Handle("/bubble", websocket.Handler(bubble))

	http.ListenAndServe(":1234", nil)
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
