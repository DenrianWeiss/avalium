package forward

import (
	"bytes"
	"github.com/DenrianWeiss/avalium/config"
	"github.com/DenrianWeiss/avalium/model"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"net/http"
)

func ReceiveHttp(r *fasthttp.RequestCtx) {
	b := r.Request.Body()
	// Check RPC Handlers
	rStream := RpcPullUntilAvailable(b)
	if rStream == nil {
		r.SetStatusCode(http.StatusGatewayTimeout)
		r.SetContentType("application/json")
	} else {
		r.SetBodyStream(rStream, -1)
		r.SetStatusCode(http.StatusOK)
		r.SetContentType("application/json")
	}
	return
}

func RpcPullUntilAvailable(req []byte) io.Reader {
	layer := model.GetMaxLayer()
	for i := 0; i < layer; i++ {
		ep := model.GetRandomRPC(i)
		if config.IsDebug() {
			log.Println("Forwarding to", ep.RpcUrl)
		}
		post, err := http.Post(ep.RpcUrl, "application/json", bytes.NewReader(req))
		if err != nil {
			continue
		}
		v, rStream := model.CheckRespErr(post)
		if v {
			return rStream
		}
	}
	return nil
}
