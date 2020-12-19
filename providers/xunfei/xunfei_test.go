package xunfei

import (
	"fmt"
	"testing"

	"github.com/zairl23/config"
)


func init() {
	// init config
	if err := config.Init("config.yaml"); err != nil {
		panic(err)
	}
}

func TestTtsOnline(t *testing.T) {
	xunfei := NewXunfei(config.Get("HOSTURL"), config.Get("HOST"), config.Get("APPID"), config.Get("APISECRET"), config.Get("APIKEY"))
	err := xunfei.TtsOnline("当当是个聪明的小男孩", "test.mp3")

	if err != nil {
		t.Fatal(err)
	}

}