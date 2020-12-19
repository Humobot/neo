package neo

import (
	"fmt"
	"testing"

	"humobot/neo/providers/xunfei"

	"github.com/zairl23/config"
)


func init() {
	// init config
	if err := config.Init("config.yaml"); err != nil {
		panic(err)
	}
}

func TestSpeak(t *testing.T) {
	xunfei := xunfei.NewXunfei(config.Get("providers.xunfei.HOSTURL"), config.Get("providers.xunfei.HOST"), config.Get("providers.xunfei.APPID"), config.Get("providers.xunfei.APISECRET"), config.Get("providers.xunfei.APIKEY"))
	
	AddProvider("xunfei", xunfei)

	neo := NewNeo()

	err := neo.Use("xunfei").Speak("当当的妈是个好母亲", "online", "neo_mom.mp3")

	if err != nil {
		t.Fatal(err)
	}

}