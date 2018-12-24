package main

import (
	"./lib/redis"
	"./settings/jsonParse"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

// static path
var staticPath string
// template path
var templatePath string
// define Employee struct instance
var Employees settings.EmployeeList
// define redis interface instance
var redisClient redis.RedisR
// max messages for redis save
const MaxLengthBubble int = 100
// max count for websocket conn
const WebsocketMax = 10
// error code for apply websocket conn apply
const GetWebsocketFailed = -1

var pool WebsocketPool

type WebsocketUtil struct {
	connId uuid.UUID
	conn *websocket.Conn
	status bool
	activeTime int64
}

type WebsocketPool struct {
	sync.Mutex
	util []WebsocketUtil
	poolId uuid.UUID
}

func (cr *WebsocketPool)GetWebsocketConn(ws *websocket.Conn) bool{
	cr.Lock()
	activeTime := time.Now().Unix()
	defer cr.Unlock()
	if len(cr.util) >= WebsocketMax{
		fmt.Println("Reach the max and  will pop the previous conn")
		cr.util = cr.util[1:]
		IdTmp , errUuid := uuid.NewV4()
		if errUuid != nil{
			panic("Generate uuid failed")
		}
		cr.util = append(cr.util, WebsocketUtil{IdTmp, ws, true, activeTime})
	}else{
		IdTmp , errUuid := uuid.NewV4()
		if errUuid != nil{
			panic("Generate uuid failed")
		}
		cr.util = append(cr.util, WebsocketUtil{IdTmp, ws, true, activeTime})
	}
	return true
}

func (cr *WebsocketPool)FindWebsocketByConn(ws *websocket.Conn) (conn *WebsocketUtil){
	cr.Lock()
	defer cr.Unlock()

	var tmp int
	for index, value :=range cr.util{
		if value.conn == ws{
			tmp = index
		}
	}
	return &cr.util[tmp]
}


func (cr *WebsocketPool)UpdateWebsocketConnStatus(id uuid.UUID, status bool) {
	cr.Lock()
	defer cr.Unlock()
	var tmp int
	for index, value :=range cr.util{
		if value.connId == id{
			tmp = index
		}
	}
	cr.util[tmp].status = status
}

func (cr *WebsocketPool)ReleaseWebsocketConn(id uuid.UUID){
	//if _, isPresent := cr.util[id];isPresent{
	//	delete(cr.util, id)
	//}
	cr.Lock()
	defer cr.Unlock()

	var tmp int
	for index, value :=range cr.util{
		if value.connId == id{
			tmp = index
		}
	}
	for i,_ := range cr.util[tmp+1:]{
		cr.util = append(cr.util[:tmp], cr.util[i])
	}
}

func (cr *WebsocketPool)DeletePool(){
	pool = WebsocketPool{}
}


func init(){
	var Config settings.Config
	var ConfigData settings.JsonParse = &Config

	var Employee settings.EmployeeList
	var EmployeeData settings.JsonParse = &Employee

	// parse config.json
	errConfig := ConfigData.Load("settings/config.json")
	if errConfig != nil {
		panic("Parse config json failed")
	}
	staticPath = Config.Static
	templatePath = Config.Template

	// parse employee.json
	employeeConfig := EmployeeData.Load("static/employee.json")
	if employeeConfig != nil {
		panic("Parse config json failed")
	}
	Employees = Employee

	redisClient = &redis.GlobalRedisClient
}



func bubbleTest(ws *websocket.Conn) {
	var err error
	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println(err)
			panic(err)
		}

		RedisBubbleLength, errLength := redisClient.RedisLLen("bubble")
		if errLength != nil {
			panic(errLength)
		}
		for RedisBubbleLength > MaxLengthBubble {
			_, errGetData := redisClient.RedisLpop("bubble")
			if errGetData != nil {
				panic(errGetData)
			}
			RedisBubbleLength -= 1

		}
		errPush := redisClient.RedisRpush("bubble", reply)
		if errPush != nil {
			panic(errPush)
		} else {
			errStatus := redisClient.RedisSet("bubbleStatus", "true")
			if errStatus != nil {
				panic(errStatus)

			}
			msg := reply
			if err = websocket.Message.Send(ws, msg); err != nil {
				panic(err)
			}
		}
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

func _bubble(ws *websocket.Conn) {
	if status := pool.GetWebsocketConn(ws); status == false{
		if errSend := websocket.Message.Send(ws, "status:获取连接失败");errSend != nil{
			fmt.Println(errSend)
			panic("Send msg to web socket Error")
		}
		return
	}
	fmt.Print(pool)
	defer pool.DeletePool()
	for {
		bubbleLength, errLength := redisClient.RedisLLen("bubble")
		fmt.Println("12")
		if errLength != nil{
			panic("Get bubble length from redis Error")
		}
		if bubbleLength < 1{
			break
		}
		fmt.Println(bubbleLength)
		data, errGetData := redisClient.RedisRpop("bubble")
		fmt.Println(string(data))
		if errGetData != nil{
			fmt.Println(errGetData)
			panic("Get bubble data from redis Error")
		}
		for _, _ws := range pool.util{
			fmt.Println(_ws)
			if err := websocket.Message.Send(_ws.conn, string(data)); err != nil {
				panic("Send msg to web socket Error")
			}
		}

	}
}

func bubble(ws *websocket.Conn) {
	if status := pool.GetWebsocketConn(ws); status == false{
		if errSend := websocket.Message.Send(ws, "status:获取连接失败");errSend != nil{
			fmt.Println(errSend)
			panic("Send msg to web socket Error, Get new conn phase")
		}
		return
	}
	//for _, value := range pool.util{
	//	fmt.Println("---------")
	//	fmt.Println(value.connId)
	//	fmt.Println(value.conn)
	//	fmt.Println("---------")
	//}


	for {
		START:
		var reply string

		if errRecv := websocket.Message.Receive(ws, &reply); errRecv != nil {
			fmt.Println(errRecv)
			panic(errRecv)
		}
		wsRecv :=pool.FindWebsocketByConn(ws)
		wsRecv.activeTime = time.Now().Unix()


		time.Sleep(1*time.Second)
		bubbleLength, errLength := redisClient.RedisLLen("bubble")
		if errLength != nil{
			fmt.Println(errLength)
			panic("Get bubble length from redis Error")
		}
		if bubbleLength < 1{
			time.Sleep(10*time.Second)
			goto START
		}

		data, errGetData := redisClient.RedisRpop("bubble")
		fmt.Println(string(data))
		if errGetData != nil{
			fmt.Println(errGetData)
			panic("Get bubble data from redis Error")
		}

		for _, _ws := range pool.util {
			if time.Now().Unix() - _ws.activeTime < 10{
				fmt.Println(pool.util)
				if err := websocket.Message.Send(_ws.conn, string(data)); err != nil {
					fmt.Println(err)
					fmt.Println(_ws.connId)
					pool.ReleaseWebsocketConn(_ws.connId)
				}
				//_ws.status = false
				//pool.UpdateWebsocketConnStatus(_ws.connId, false)
			}else{
				pool.ReleaseWebsocketConn(_ws.connId)
			}
		}
	}
}


func SetTimer(){
	time.Sleep(3* time.Second)
}


func main() {
	var errUUID error
	pool = WebsocketPool{}
	pool.poolId , errUUID= uuid.NewV4()
	if errUUID != nil{
		panic("Apply Pool failed")
	}

	// static files url
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	// award page, just show the page
	http.HandleFunc("/award/", award)
	// login page
	http.HandleFunc("/login/", login)
	// Get the bubble text from user and store in redis
	http.Handle("/web_socket", websocket.Handler(bubbleTest))
	// get the bubble text from redis and send to award page
	http.Handle("/bubble", websocket.Handler(bubble))

	http.ListenAndServe(":1234", nil)
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
