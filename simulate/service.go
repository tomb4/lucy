package simulate

import (
	"encoding/json"
	"log"
	pb "lucy/proto/party"
	"os"
	"os/signal"
	"time"

	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	StressMoveMethod = 1 // 多人移动压测
)

type Service struct {
	partyCli pb.PartyServiceClient
	userList []User
	count    int
	httpCli  *req.Client
	conf     *Config
	agentMgr *AgentManager
}

func NewService(cfg string, count int) *Service {
	svc := &Service{}
	svc.conf = loadConfig(cfg)

	svc.httpCli = req.C()
	svc.userList = loadUser("./simulate/testdata/user.json")
	svc.count = count
	conn, err := grpc.Dial(svc.conf.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	svc.partyCli = pb.NewPartyServiceClient(conn)
	svc.agentMgr = NewAgentManager()

	return svc
}

func loadConfig(path string) *Config {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	c := Config{
		PartyId:  viper.GetString("service.partyId"),
		HttpAddr: viper.GetString("service.httpAddr"),
		WsAddr:   viper.GetString("service.wsAddr"),
		GrpcAddr: viper.GetString("service.grpcAddr"),
	}
	return &c
}

func loadUser(path string) []User {
	bts, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	users := make([]struct {
		ID      string `json:"id"`
		Account struct {
			Token string `json:"token"`
		} `json:"account"`
	}, 0)

	_ = json.Unmarshal(bts, &users)

	userList := make([]User, 0, len(users))
	for _, v := range users {
		userList = append(userList, User{
			Id:    v.ID,
			Token: v.Account.Token,
		})
	}
	return userList
}

func (s *Service) Handle(method int) {
	switch method {
	case StressMoveMethod:
		s.StressMove()
	default:
		log.Println("unknown method")
	}
}

func (s *Service) StressMove() {
	interrupt := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(interrupt, os.Interrupt)

	if s.count > len(s.userList) {
		log.Println("count overflow")
		return
	}

	users := s.userList[:s.count]
	for _, user := range users {
		resp := Resp{}
		_, err := s.httpCli.R().
			SetHeader("Authorization", "token "+user.Token).
			SetBodyJsonMarshal(MetaJoinReq{
				PartyId: s.conf.PartyId,
			}).
			SetResult(&resp).
			Post(s.conf.HttpAddr + "/v1/party/meta/join")
		if err != nil {
			log.Println(err)
			continue
		}
		if resp.Code != 0 {
			log.Println(resp.Message)
			continue
		}

		ag := NewAgent(user.Id, user.Token)
		s.agentMgr.AddAgent(ag)
	}
	s.agentMgr.AgentMove(s.conf.WsAddr, stop)

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
				PartyId: s.conf.PartyId,
			}).
			SetResult(&resp).
			Post(s.conf.HttpAddr + "/v1/party/meta/exit")
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
