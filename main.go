package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-websocket/message"
	"go-websocket/ws"
	"net/http"
	"text/template"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var ctx = context.Background()

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
	//如果是分布式，每台接口(API)服务器,都需要启动，每台几千管理自己的websocket连接池
	go ws.Manager.Start()

	manage1 := http.NewServeMux()
	manage1.HandleFunc("/admin/index", Index)
	manage1.HandleFunc("/admin/message/send", message.Send)
	go http.ListenAndServe("127.0.0.1:8081", manage1)

	interface1 := http.NewServeMux()
	interface1.HandleFunc("/interface/ws", IndexHandler)
	go http.ListenAndServe("127.0.0.1:8082", interface1)

	interface2 := http.NewServeMux()
	interface2.HandleFunc("/interface/ws1", IndexHandler)
	go http.ListenAndServe("127.0.0.2:8083", interface2)
	redisSubscription()
	fmt.Println("启动成功")
	select {}
}

func redisSubscription() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:63791",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.PubSubChannels(ctx, "chan_message")
	go func() {
		pb := rdb.Subscribe(ctx, "chan_message")
		ch := pb.Channel()
		for msg := range ch {
			params := make(map[string]interface{}, 0)
			json.Unmarshal([]byte(msg.Payload), &params)
			go ws.Manager.Send([]byte(params["title"].(string)), &ws.Client{})
		}
	}()
}
