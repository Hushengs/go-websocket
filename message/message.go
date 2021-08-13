package message

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func Send(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	dataMessage := make(map[string]interface{}, 0)
	dataMessage["message_id"] = "8633"
	dataMessage["title"] = message
	messages, _ := json.Marshal(dataMessage)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:63791",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rdb.Publish(ctx, "chan_message", messages)
	rdb.Close()
}
