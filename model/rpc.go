package model

import (
	"math/rand"
	"sync"
)

type RPC struct {
	RpcUrl    string `json:"rpc_url" mapstructure:"rpc_url"`
	RpcWeight int    `json:"rpc_weight" mapstructure:"rpc_weight"`
}

type RPCList struct {
	TireN [][]RPC `json:"tire_n" mapstructure:"tire_n"`
}

var (
	cached    = RPCList{}
	cacheLock sync.RWMutex
)

func GenerateCache(r RPCList) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cached = RPCList{}
	for _, layer := range r.TireN {
		var tire []RPC
		for _, rpc := range layer {
			for i := 0; i < rpc.RpcWeight; i++ {
				tire = append(tire, rpc)
			}
		}
		cached.TireN = append(cached.TireN, tire)
	}

}

func GetRandomRPC(layerN int) RPC {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	layer := cached.TireN[layerN]
	rpc := layer[rand.Intn(len(layer))]
	return rpc
}

func GetMaxLayer() int {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	return len(cached.TireN)
}

func KickoutRpc(rpc string) {

}
