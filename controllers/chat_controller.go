package controllers

import (
	"HongXunServer/middleware"
	"HongXunServer/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"log"
	"net/http"
	"sync"
)

type ChatController struct {
	Ctx     iris.Context
}

type Node struct {
	Conn *websocket.Conn
	//并行转串行,
	DataQueue chan []byte
	UserEmail string
}

//定义命令行格式
const (
	CmdSingleMsg = 0
	CmdHeart     = 1
)



//后端调度逻辑处理
func dispatch(userFrom string,data []byte) {
	message := models.Message{}
	err := json.Unmarshal(data, &message)
	message.UserFrom = userFrom
	message.UserName = userFrom
	log.Println(message)
	if err != nil {
		log.Println(err.Error())
	}
	msg, _ := json.Marshal(message)
	switch message.Type {
	case CmdSingleMsg:
		sendMsg(message.UserTo, msg)
	case CmdHeart:
		//检测客户端的心跳
	}
}

//userid和Node映射关系表
var clientMap  = make(map[string]*Node, 0)
//读写锁
var rwlocker sync.RWMutex
//实现聊天的功能
func Chat(ctx iris.Context) {
	log.Println("Chat")
	var	userClaims models.UserClaims
	token := jwt.FromQuery(ctx)
	err := middleware.J.VerifyTokenString(ctx, token, &userClaims)
	isLegal := true
	if err != nil {
		isLegal = false
		log.Println(err)
	}
	userEmail := userClaims.UserEmail
	log.Println("UserEmail:", userEmail)
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isLegal
		}}).Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)

	if err != nil {
		log.Println(err)
		return
	}
	//获得websocket链接conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		UserEmail: userEmail,
	}

	rwlocker.Lock()
	clientMap[userEmail] = node
	log.Println(clientMap)
	rwlocker.Unlock()

	//开启协程处理发送逻辑
	go sendproc(node)

	//开启协程完成接收逻辑
	go recvproc(node)

	sendMsg(userEmail, []byte("welcome!"))
}

//发送逻辑
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

//接收逻辑
func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		dispatch(node.UserEmail, data)
		//todo对data进一步处理
		log.Printf("recv<=%s", data)
	}
}

//发送消息,发送到消息的管道
func sendMsg(UserEmail string, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[UserEmail]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

