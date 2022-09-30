package simulate

import (
	"github.com/duke-git/lancet/convertor"
	"github.com/duke-git/lancet/random"
	"github.com/gorilla/websocket"
	"log"
	pb "lucy/proto/MetaGateway"
	"net/url"
	"time"
)

const (
	MetaSceneInitialX = 2
	MetaSceneInitialY = 3.4
	MetaSceneInitialZ = 1
)

type Agent struct {
	UserId string
	Uid    int32
	Token  string
	X      float32
	Y      float32
	Z      float32
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

func (a *Agent) ChangePosition() {
	a.X += 0.3
	a.Z += 0.3
	if a.X > 30 {
		a.X = MetaSceneInitialX
	}
	if a.Z > 20 {
		a.Z = MetaSceneInitialY
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

func (a *Agent) StartMove(stop chan bool) {
	u := url.URL{Scheme: "ws", Host: WsServer, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
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
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
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
			a.ChangePosition()
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
