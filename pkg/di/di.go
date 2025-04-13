package di

import (
	"log"
	"net/http"

	"github.com/ratheeshkumar25/chatApp/pkg/api"
	"github.com/ratheeshkumar25/chatApp/pkg/chat"
	"github.com/ratheeshkumar25/chatApp/pkg/server"
)

func Init() {
	chatRoom := chat.NewChatRoom()

	handler := api.NewHandler(chatRoom)
	server := server.NewServer(handler)

	log.Println("Server staring :3000")

	log.Fatal(http.ListenAndServe(":3000", server.Router))

}
