package simulate

import "testing"

func Test_LoadUser(t *testing.T) {
	t.Log(len(loadUser("./testdata/user.json")))
}
