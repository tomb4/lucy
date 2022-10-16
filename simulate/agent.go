package simulate

import (
	"log"
	pb "lucy/proto/MetaGateway"
	"strings"
	"time"

	"github.com/duke-git/lancet/convertor"
	"github.com/duke-git/lancet/random"
	"github.com/gorilla/websocket"
)

type Direction int

type DtAct struct {
	OffsetX float32
	OffSetZ float32
}

const (
	up    Direction = 1
	down  Direction = 2
	left  Direction = 3
	right Direction = 4
)

var (
	DirectionMap = map[Direction]DtAct{
		up:    {OffsetX: 0, OffSetZ: 2},
		down:  {OffsetX: 0, OffSetZ: -2},
		left:  {OffsetX: -2, OffSetZ: 0},
		right: {OffsetX: 2, OffSetZ: 0},
	}

	SquareSteps = []Direction{left, up, right, right, down, down, left, left, up, right}
)

type Agent struct {
	UserId string
	Uid    int32
	Token  string
	X      float32
	Y      float32
	Z      float32
	LastDt Direction
}

func NewAgent(uid string, token string) *Agent {
	id, _ := convertor.ToInt(uid)
	return &Agent{
		UserId: uid,
		Uid:    int32(id),
		Token:  token,
		X:      MetaSceneInitialX,
		Y:      MetaSceneInitialY,
		Z:      MetaSceneInitialZ,
	}
}

func (a *Agent) SetPosition(x, y, z float32) {
	a.X = x
	a.Y = y
	a.Z = z
}

func (a *Agent) ChangePosition() {
	act, ok := DirectionMap[SquareSteps[a.LastDt]]
	if !ok {
		return
	}

	a.X += act.OffsetX
	a.Z += act.OffSetZ

	a.LastDt++
	if int(a.LastDt) >= len(SquareSteps) {
		a.LastDt = 0
	}
}

func (a *Agent) SendMessage(c *websocket.Conn, pack Packet) error {
	bts, err := pack.encode()
	if err != nil {
		return err
	}
	err = c.WriteMessage(websocket.BinaryMessage, bts)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) StartMove(url string, stop chan bool, callback func(uid string)) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read_err:", err)
				return
			}
			//log.Printf("recv: %s", message)
			if strings.Contains(string(message), "登录失败") {
				log.Println("登录失败", a.Uid)
				return
			}
		}
	}()

	// login
	err = a.SendMessage(c, NewPacket(int32(pb.CmdId_LoginReqCmdId), &pb.LoginReq{
		UserId:    a.Uid,
		Token:     a.Token,
		Type:      1,
		LoginType: 1,
		ReqId:     time.Now().UnixNano() + int64(random.RandInt(1, 1000)),
	}))
	if err != nil {
		log.Println("login err:", err)
		return
	}

	// move
	ticker := time.NewTicker(time.Second)
	pingTk := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			log.Println("done -> ", a.Uid)
			callback(a.UserId)
			_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		case <-pingTk.C:
			log.Println("ping...", a.Uid)
			err = a.SendMessage(c, NewPacket(int32(pb.CmdId_PingCmdId), &pb.Ping{
				ReqId: time.Now().UnixNano() + int64(random.RandInt(1, 1000)),
			}))
			if err != nil {
				log.Println("ping write:", err)
				return
			}
		case <-ticker.C:
			log.Println("moving...", a.Uid)
			err = a.SendMessage(c, NewPacket(int32(pb.CmdId_ClientStateEventReqCmdId), &pb.ClientStateEventReq{
				User: &pb.UserStateEvent{
					X:            a.X,
					Y:            a.Y,
					Z:            a.Z,
					UserId:       a.Uid,
					From:         a.Uid,
					OnlineStatus: 1,
				},
				ReqId: time.Now().UnixNano() + int64(random.RandInt(1, 1000)),
			}))
			if err != nil {
				log.Println("move write:", err)
				return
			}
			a.ChangePosition()
		case <-stop:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}
