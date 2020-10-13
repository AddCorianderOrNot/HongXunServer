package controllers

import (
	"HongXunServer/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sync"
)

type ChatController struct {
	Ctx     iris.Context
}

type Node struct {
	Conn *websocket.Conn
	//并行转串行,
	DataQueue chan []byte
}

//定义命令行格式
const (
	CmdSingleMsg = 10
	CmdHeart     = 0
)



//后端调度逻辑处理
func dispatch(data []byte) {
	msg := models.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	switch msg.Type {
	case CmdSingleMsg:
		sendMsg(msg.UserTo, data)
	case CmdHeart:
		//检测客户端的心跳
	}
}

//userid和Node映射关系表
var clientMap map[primitive.ObjectID]*Node = make(map[primitive.ObjectID]*Node, 0)
//读写锁
var rwlocker sync.RWMutex
//实现聊天的功能
func (c *ChatController) Get() {
	log.Println("Chat")
	var	userClaims models.UserClaims
	_ = jwt.ReadClaims(c.Ctx, &userClaims)
	userId := userClaims.UserId
	log.Println("UserId:", userId)
	conn, err := (&websocket.Upgrader{}).Upgrade(c.Ctx.ResponseWriter(), c.Ctx.Request(), nil)

	if err != nil {
		log.Println(err.Error())
		return
	}
	//获得websocket链接conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
	}

	rwlocker.Lock()
	clientMap[userId] = node
	log.Println(clientMap)
	rwlocker.Unlock()

	//开启协程处理发送逻辑
	go sendproc(node)

	//开启协程完成接收逻辑
	go recvproc(node)

	sendMsg(userId, []byte("welcome!"))
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

		dispatch(data)
		//todo对data进一步处理
		log.Printf("recv<=%s", data)
	}
}

//发送消息,发送到消息的管道
func sendMsg(userId primitive.ObjectID, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

