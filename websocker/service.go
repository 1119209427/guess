package websocker

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"guess/logic"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	writeWait = 120 * time.Second

	pongWait = 120 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

//提供读取信息和写入信息的功能
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// GuessConn 玩家的链接
type GuessConn struct {
	ws         *websocket.Conn
	send       chan []byte
	number     int   //猜题次数
	numberKick int   //踢出的记录次数2
	timeLog    int64 //等待时间
	flag       bool  //判读是否在游戏状态
	difficulty bool  //是否选择难度
	legitimate bool  //是否合法
	secrete    int   //随机的生成数
	//first bool   //是否为第一次计数
}

//var secreteMap map[string]int//存储随机数

// Message 处理消息的结构体
type Message struct {
	hub    *Hub
	data   []byte
	roomid string
	conn   *GuessConn
}

// ReadPump 将消息读入处理管道
func (m *Message) ReadPump() {
	c := m.conn
	//防止浪费资源
	defer func() {
		m.hub.unlogin <- m
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	//SetReadDeadline设置基础网络连接的读取截止时间。
	err := c.ws.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fmt.Println(err.Error())
	}
	//SetPongHandler为从对等方接收的pong消息设置处理程序。
	c.ws.SetPongHandler(func(string) error {
		err = c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return err
	})
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println(err.Error())
			}
			break
		}
		go m.GuessCheck(msg)
	}
}

func (m *Message) GuessCheck(msg []byte) {
	c := m.conn
	//判读消息是否合法
	baseStr := "准备简单中等困难"
	if c.numberKick < 3 {
		testStr := string(msg[:])
		for _, i := range testStr {
			flag := strings.Contains(baseStr, string(i))
			if !flag || !IsNum(string(i)) {
				c.legitimate = false
				c.numberKick += 1
				m.hub.warnings <- m
			}
			break
		}
	}
	if c.numberKick > 3 {
		m.hub.kickOut <- m
		fmt.Println("要被踢出群聊了")
		c.ws.Close()
	}
	//判读是否准备，没有提醒，有的话进行难度选择
	if !c.flag && c.legitimate {
		m.hub.unprepare <- m //将消息传入提醒管道
	}
	if string(msg[:]) == "准备" {
		c.flag = true

		m.hub.unprepare <- m
	}
	//准备
	if c.flag && !c.difficulty && c.legitimate {
		//判读难度
		if string(msg[:]) == "简单" {
			c.number = 10 //设置猜数次数
			c.difficulty = true
			c.secrete = logic.Easy()
			m.hub.easy <- m
		}
		if string(msg[:]) == "中等" {
			c.number = 20 //设置猜数次数
			c.difficulty = true
			c.secrete = logic.Medium()
			m.hub.medium <- m
		}
		if string(msg[:]) == "困难" {
			c.number = 30 //设置猜数次数
			c.difficulty = true
			c.secrete = logic.Hard()
			m.hub.hard <- m
		}

	}
	if c.flag && c.difficulty && c.legitimate {
		str := logic.CheckGuess(string(msg[:]), c.secrete)
		//如果不正确
		if !strings.Contains(str, "Correct, you Legend!") {
			m.hub.remindCnt <- m
			c.number--
		}
		M := &Message{data: []byte(str), roomid: m.roomid, conn: c}
		if c.number == 0 {
			m.hub.remind <- m
			c.flag = false
			c.difficulty = false
		}
		m.hub.broadcast <- M
	}
}

func (c *GuessConn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (m *Message) WriteGuess() {
	c := m.conn
	ticker := time.NewTimer(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.write(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Println(err)
					return
				}

			}
			err := c.write(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}
}

func ServerWs(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	roomid := ctx.Param("room_id")
	hub := NewHub()

	c := &GuessConn{send: make(chan []byte, 256), ws: conn}
	m := &Message{hub, nil, roomid, c}
	m.hub.login <- m
	go m.ReadPump()
	go m.WriteGuess()

}
