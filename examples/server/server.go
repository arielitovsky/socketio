package main

import (
	"log"
	"net/http"
	"time"

	"github.com/arielitovsky/ariesocketio"
	"github.com/arielitovsky/ariesocketio/websocket"
)

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
}

type Desc struct {
	Text string `json:"text"`
}

func main() {
	server := ariesocketio.NewServer(*websocket.GetDefaultWebsocketTransport())

	server.On(ariesocketio.OnConnection, func(c *ariesocketio.Channel) {
		log.Println("[server] connected! id:", c.Id())
		log.Println("[server]", c.RemoteAddr().Network()+" "+c.RemoteAddr().String()+
			" --> "+c.LocalAddr().Network()+" "+c.LocalAddr().String())

		c.Join("room")
		server.BroadcastTo("room", "/admin", Message{1, "new members!"})
		time.Sleep(100 * time.Millisecond)
		c.BroadcastTo("room", "/admin", Message{2, "hello everyone!"})

		_ = c.Emit("message", Message{10, "{\"chinese\":\"中文才是最屌的\"}"})

		// return [][]byte
		result, err := c.Ack("/ackFromServer", time.Second*5, "go", 3)
		if err != nil {
			log.Println("[server] ack cb err:", err)
			return
		}
		log.Println("[server] ack cb:", result)

		time.Sleep(3 * time.Second)
		log.Println("ReadBytes", c.ReadBytes())
		log.Println("WriteBytes", c.WriteBytes())
	})
	server.On(ariesocketio.OnDisconnection, func(c *ariesocketio.Channel, reason websocket.CloseError) {
		log.Println("[server] received disconnect", c.Id(), "code:", reason.Code, "text:", reason.Text)
	})

	server.On("message", func(c *ariesocketio.Channel, arg1 string, arg2 Message, arg3 int, arg4 bool) {
		log.Println("[server] received message:", "arg1:", arg1, "arg2:", arg2, "arg3:", arg3, "arg4:", arg4)
	})

	// listen ack event
	server.On("/ackFromClient", func(c *ariesocketio.Channel, msg Message, num int) (int, Desc, string) {
		log.Println("[server] received ack:", msg, num)
		return 1, Desc{Text: "resp"}, "server"
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	log.Println("[server] starting server...")
	log.Panic(http.ListenAndServe(":2233", serveMux))
}
