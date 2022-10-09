package websocker

import "fmt"

type Hub struct {
	rooms      map[string]map[*GuessConn]bool
	difficulty chan *Message
	broadcast  chan *Message //广播猜数信息
	warnings   chan *Message //警告信息
	remind     chan *Message //提醒猜数已次数以完
	remindCnt  chan *Message //提醒剩余猜数次数和
	unprepare  chan *Message
	easy       chan *Message
	medium     chan *Message
	hard       chan *Message
	kickOut    chan *Message
	login      chan *Message
	unlogin    chan *Message
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*GuessConn]bool),
		difficulty: make(chan *Message),
		broadcast:  make(chan *Message),
		warnings:   make(chan *Message),
		unprepare:  make(chan *Message),
		kickOut:    make(chan *Message),
		login:      make(chan *Message),
		unlogin:    make(chan *Message),
		remind:     make(chan *Message),
		remindCnt:  make(chan *Message),
		easy:       make(chan *Message),
		medium:     make(chan *Message),
		hard:       make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case m := <-h.login:
			conns := h.rooms[m.roomid]
			if conns == nil {
				conn := make(map[*GuessConn]bool)
				h.rooms[m.roomid] = conn
				fmt.Println("注册人数", len(conns))
				fmt.Println("room==", h.rooms)
			}
			h.rooms[m.roomid][m.conn] = true
			fmt.Println("在线人数 == ", len(conns))
			fmt.Println("rooms ==", h.rooms)
			for con := range conns {
				delMsg := "系统信息:欢迎新玩家加入" + m.roomid
				data := []byte(delMsg)
				select {
				case con.send <- data:
				}
			}
		case m := <-h.unlogin:
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					delete(conns, m.conn)
					close(m.conn.send)
					for conn := range conns {
						delMsg := "系统消息:有玩家离开了" + m.roomid
						data := []byte(delMsg)
						select {
						case conn.send <- data:
						}
					}
					if len(conns) == 0 {
						delete(h.rooms, m.roomid)
					}
				}
			}
		case m := <-h.kickOut:
			conns := h.rooms[m.roomid]
			notice := "系统公告:由于您多次发送不合法信息,已被踢出！！！"
			select {
			case m.conn.send <- []byte(notice):
			}
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					delete(conns, m.conn)
					close(m.conn.send)
					if len(conns) == 0 {
						delete(h.rooms, m.roomid)
					}
				}
			}
		case m := <-h.warnings:
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := fmt.Sprintf("系统公告:请输入合法字符！！！你还有%d发言资格", 3-m.conn.numberKick)
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}
		case m := <-h.unprepare:
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := "系统公共:请准备后再发言，准备游戏请在聊天框输入准备"
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}

		case m := <-h.easy:
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := "系统公告:你选择了简单难度，游戏规则如下:电脑随机生成两个数字，用户输入两个数字。电脑反馈XAYB，x代表数字完全猜对的数量，y代表数字猜对但位置不对的数量。如电脑为12，用户输入12，返回2A0B."
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}
		case m := <-h.medium:
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := "系统公告:你选择了中等难度，游戏规则如下:电脑随机生成三个数字，用户输入三个数字。电脑反馈XAYB，x代表数字完全猜对的数量，y代表数字猜对但位置不对的数量.用户输入124，返回2A1B."
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}
		case m := <-h.hard:
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := "系统公告:你选择了困难难度，游戏规则如下:电脑随机生成四个数字，用户输入四个数字。电脑反馈XAYB，x代表数字完全猜对的数量，y代表数字猜对但位置不对的数量。如电脑为1234，用户输入1243，返回2A2B."
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}
		case m := <-h.remindCnt:
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := fmt.Sprintf("系统公告:你猜错了还有%d次机会", m.conn.number)
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}
		case m := <-h.remind:
			conns := h.rooms[m.roomid]
			if conns != nil {
				if _, ok := conns[m.conn]; ok {
					notice := "系统公告:游戏失败，如果重新开始请依次在聊天框输入准备和难度选择"
					select {
					case m.conn.send <- []byte(notice):
					}
				}
			}
		case m := <-h.broadcast:
			conns := h.rooms[m.roomid]
			for con := range conns {
				if con == m.conn { //自己发送的信息，不用再发给自己
					continue
				}
				select {
				case con.send <- m.data:
				default:
					close(con.send)
					delete(conns, con)
					if len(conns) == 0 {
						delete(h.rooms, m.roomid)
					}
				}
			}
		}
	}
}
