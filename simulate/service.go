package simulate

import (
	"github.com/imroc/req/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "lucy/proto/party"
	"os"
	"os/signal"
	"time"
)

const (
	StressMoveMethod = 1 // 多人移动压测
	HttpServer       = "https://api-test.booyah.cc"
	//WsServer = "localhost:8888"
	WsServer = "meta-gateway-test.booyah.cc"
)

type (
	Input struct {
		Method int
		Count  int
	}

	User struct {
		Id    string
		Token string
	}

	Resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	MetaJoinReq struct {
		PartyId   string `json:"partyId"`
		FollowUid string `json:"followUid"`
	}

	MetaExitReq struct {
		PartyId string `json:"partyId"`
	}
)

type Service struct {
	partyCli pb.PartyServiceClient
	userList []User
	partyId  string
	count    int
	httpCli  *req.Client
}

func NewService() *Service {
	svc := &Service{}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	svc.httpCli = req.C()
	svc.partyCli = pb.NewPartyServiceClient(conn)
	svc.userList = []User{
		{Id: "25647", Token: "1a2fd588-4513-423c-bf79-b4b7783b72c0"},
	}
	svc.partyId = "31522535"

	return svc
}

func (s *Service) Handle(in Input) {
	s.count = in.Count

	switch in.Method {
	case StressMoveMethod:
		s.StressMove(in.Count)
	default:
		log.Println("unknown method")
	}
}

func (s *Service) StressMove(c int) {
	interrupt := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(interrupt, os.Interrupt)

	if c > len(s.userList) {
		log.Println("count overflow")
		return
	}

	users := s.userList[:c]
	for _, user := range users {
		resp := Resp{}
		_, err := s.httpCli.R().
			SetHeader("Authorization", "token "+user.Token).
			SetBodyJsonMarshal(MetaJoinReq{
				PartyId: "31522535",
			}).
			SetResult(&resp).
			Post(HttpServer + "/v1/party/meta/join")
		if err != nil {
			log.Println(err)
			return
		}
		if resp.Code != 0 {
			log.Println(resp.Message)
			return
		}

		go NewAgent(user.Id, user.Token).StartMove(stop)
	}

	<-interrupt

	log.Println("stop")
	close(stop)
	time.Sleep(time.Second * 1)

	s.Clear()
	log.Println("clear")
}

func (s *Service) Clear() {
	users := s.userList[:s.count]
	for _, user := range users {
		resp := Resp{}
		_, err := s.httpCli.R().
			SetHeader("Authorization", "token "+user.Token).
			SetBodyJsonMarshal(MetaExitReq{
				PartyId: "31522535",
			}).
			SetResult(&resp).
			Post(HttpServer + "/v1/party/meta/exit")
		if err != nil {
			log.Println(err)
			return
		}
		if resp.Code != 0 {
			log.Println(resp.Message)
			return
		}
		time.Sleep(time.Millisecond * 10)
	}
}
