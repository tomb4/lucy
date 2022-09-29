package simulate

type Agent struct {
	UserId string
}

func NewAgent(uid string) *Agent {
	return &Agent{
		UserId: uid,
	}
}

func (a Agent) StartMove() {

}
