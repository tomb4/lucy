package simulate

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"log"
	"net/url"
	"time"
)

const (
	MetaSceneInitialX = 2
	MetaSceneInitialY = 3.4
	MetaSceneInitialZ = 1
)

type AgentManager struct {
	agentMap cmap.ConcurrentMap[*Agent]
}

func NewAgentManager() *AgentManager {
	agentMap := cmap.New[*Agent]()

	mgr := new(AgentManager)
	mgr.agentMap = agentMap
	return mgr
}

func (a AgentManager) AddAgent(ag *Agent) {
	a.agentMap.Set(ag.UserId, ag)
}

func (a AgentManager) Count() int {
	return a.agentMap.Count()
}

func (a AgentManager) AgentMove(ws string, stop chan bool) {
	a.Arrange()

	u := url.URL{Scheme: "ws", Host: ws, Path: "/echo"}
	for _, v := range a.agentMap.Items() {
		go v.StartMove(u.String(), stop, func(uid string) {
			log.Println("AgentMove uid", uid)
			a.agentMap.Remove(uid)
		})
	}

	time.Sleep(time.Second * 3)
	log.Println("agent count:", a.Count())
}

func (a AgentManager) Arrange() {
	row, col, maxCol, interval := 0, 0, 10, 5
	for _, v := range a.agentMap.Items() {
		x := MetaSceneInitialX + col*interval
		z := MetaSceneInitialZ + row*interval
		v.SetPosition(float32(x), MetaSceneInitialY, float32(z))
		col++
		if col >= maxCol {
			row++
			col = 0
		}
	}
}
