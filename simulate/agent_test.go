package simulate

import (
	"fmt"
	"testing"
)

func Test_ChangePosition(t *testing.T) {
	ag := NewAgent("1", "token")
	for i := 0; i < 100; i++ {
		ag.ChangePosition()
		fmt.Println(ag.X, ag.Z)
		//t.Log(ag)
	}
}
