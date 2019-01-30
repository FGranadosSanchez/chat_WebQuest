package main

import(
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var update = websocket.Upgrader{}

type Message struct {
	Username string`json:"username"`
	From string`json:"from"`
	Message string`json:"message"`  
}

func handleConnections(w http.ResponseWriter,r *http.Request) {
	ws, err := update.Upgrade(w,r,nil)
	if err !=nil{
		log.Fatal(err)
	}
	defer ws.Close()
	clients[ws] = true
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil{
			log.Printf("err: %v",err)
			delete(clients,ws)
			break
		}
		broadcast<-msg
	}
}
func handleMessage()  {
	for{
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("err: %v",err)
				client.Close()
				delete(clients,client)
			}
		}

	}
}

func main()  {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/",fs)
	http.HandleFunc("/ws",handleConnections)
	go handleMessage()
	log.Println("Server start correct in port: 1235")
	err := http.ListenAndServe(":1235",nil)
	if err != nil {
		log.Fatal("Error en la escucha del puerto ",err)
	}
}