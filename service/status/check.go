package status

import (
	"bytes"
	"github.com/DenrianWeiss/avalium/model"
	"net/http"
	"time"
)

const (
	checkAliveInterval = 5 * time.Second
	checkAliveCount    = 5
)

func CheckAlive(rpc string) bool {
	post, err := http.Post(rpc, "application/json", bytes.NewReader(model.NewJsonRpc("1").SetMethod("eth_blockNumber").SerializeChainable()))
	if err != nil {
		return false
	}
	v, _ := model.CheckRespErr(post)
	return v
}

func CheckUntilRecovered(rpc string) {
	tick := time.NewTicker(checkAliveInterval)
	count := 0
	failCount := 0
	for {
		select {
		case <-tick.C:
			if CheckAlive(rpc) {
				count += 1
			} else {
				count = 0
				failCount += 1
			}
			if count >= checkAliveCount {
				return
			}
		}
	}
}
