package initialize

import (
	"fmt"
	"time"

	"git.selly.red/Selly-Modules/natsio"

	"git.selly.red/Selly-Server/warehouse/pkg/app/config"
)

func nats() {
	cfg := config.GetENV().Nats
	err := natsio.Connect(natsio.Config{
		URL:            cfg.URL,
		User:           cfg.Username,
		Password:       cfg.Password,
		RequestTimeout: 2 * time.Minute,
	})
	if err != nil {
		panic(err)
	}

	s := natsio.GetServer()
	c, err := s.NewJSONEncodedConn()
	if err != nil {
		panic(err)
	}

	fmt.Println(c)
}
