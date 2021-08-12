package main

import (
	"fmt"
	"go-websocket/ws"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	fmt.Println("websockct连接开始：" + username)
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)
	if error != nil {
		fmt.Println("websockct连接错误:", error)
		return
	}
	client := &ws.Client{ID: username, Socket: conn, Send: make(chan []byte)}
	fmt.Println("websockct创建客户端：", client)
	ws.Manager.Register <- client
	go client.Read()
	go client.Write()
}

func Index(w http.ResponseWriter, r *http.Request) {
	t1, err := template.ParseFiles("./view/index.html")
	if err != nil {
		panic(err)
	}
	t1.Execute(w, "")
}

func main() {
	go ws.Manager.Start()
	http.HandleFunc("/admin/index", Index)
	http.HandleFunc("/interface/ws", IndexHandler)
	http.ListenAndServe("127.0.0.1:8081", nil)
}
