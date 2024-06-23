package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Proxy struct {
	store   *Store
	env     Env
	servers map[string]*websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewProxy(env Env) *Proxy {
	return &Proxy{
		store:   NewStore(env),
		env:     env,
		servers: make(map[string]*websocket.Conn),
	}
}

func (p *Proxy) handleNew(c *gin.Context) {
	var req ConnectionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	conn := NewConnection(req)

	if err := p.store.SetConnection(conn); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, conn.Iden)
}

func verifyToken(bearer string, actual string) bool {
	token := strings.Split(bearer, " ")[1]
	return token == actual
}

func (p *Proxy) handleListen(c *gin.Context) {
	iden := c.Param("iden")

	connInfo, err := p.store.GetConnection(iden)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	// Check if the Authorization header is present
	token := c.GetHeader("Authorization")
	if token == "" || !verifyToken(token, connInfo.Token) {
		c.JSON(401, gin.H{"error": "Authorization header required"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	p.servers[iden] = ws
	ws.WriteMessage(websocket.TextMessage, []byte("Connected to the proxy server"))
}

// func (p *Proxy) handleSock(conn *websocket.Conn) {
// 	for {
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}

// 		// send back same msg
// 		//TODO: Write the client websocket connection here
// 		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
// 			break
// 		}
// 	}
// }

func (p *Proxy) handleClient(c *gin.Context) {
	iden := c.Param("iden")

	ws := p.servers[iden]
	if ws == nil {
		c.JSON(404, gin.H{"error": "Connection not found"})
		return
	}

	msg := []byte(c.Query("route"))

	if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	gotRes := false
	// read the message from websocket
	for {
		if gotRes {
			break
		}
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		c.JSON(200, gin.H{"message": string(msg)})
		gotRes = true
	}
}

func (p *Proxy) Run() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "RevLocal Proxy is running!",
		})
	})

	r.POST("/new", p.handleNew)
	r.GET("/get/:iden", p.handleListen)
	r.GET("/client/:iden", p.handleClient)

	r.Run(p.env.Addr())
}
