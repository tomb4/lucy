package simulate

import "testing"

func Test_LoadUser(t *testing.T) {
	t.Log(len(loadUser("./testdata/user.json")))
}

func Test_loadConfig(t *testing.T) {
	c := loadConfig("./config.yml")
	t.Log(c)
}
