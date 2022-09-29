package simulate

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "lucy/proto/party"
)

const (
	StressMoveMethod = 1 // 多人移动压测
)

type (
	Input struct {
		Method int
		Count  int
	}
)

type Service struct {
	partyCli pb.PartyServiceClient
}

func NewService() *Service {
	svc := &Service{}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	svc.partyCli = pb.NewPartyServiceClient(conn)

	return svc
}

func (s *Service) Handle(in Input) {
	switch in.Method {
	case StressMoveMethod:
		s.StressMove(in.Count)
	default:
		log.Println("unknown method")
	}
}

func (s *Service) StressMove(c int) {
	for i := 0; i < c; i++ {

		party, err := s.partyCli.JoinMetaParty(context.TODO(), &pb.JoinMetaPartyReq{
			PartyId: "1",
			UserId:  "1",
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(party.PartyId)

		go NewAgent("").StartMove()
	}
}
