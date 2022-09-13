package model

import (
	"bytes"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
)

const SonicNotFoundErr = "value not exists"

type RpcRequest map[string]interface{}

func NewJsonRpc(id string) *RpcRequest {
	return &RpcRequest{
		"jsonrpc": "2.0",
		"id":      id,
		"params":  map[string]interface{}{},
	}
}

func (r *RpcRequest) SetMethod(method string) *RpcRequest {
	(*r)["method"] = method
	return r
}

func (r *RpcRequest) Serialize() ([]byte, error) {
	return sonic.Marshal(r)
}

func (r *RpcRequest) SerializeChainable() []byte {
	p, _ := r.Serialize()
	return p
}

func CheckRespErr(h *http.Response) (bool, io.Reader) {
	length := h.ContentLength
	if length > 0 && length < 10*1024 {
		payload, err := io.ReadAll(h.Body)
		if err != nil {
			return false, bytes.NewReader(payload)
		}
		cont, err := sonic.Get(payload, "error")
		v, _ := cont.String()
		if err.Error() == SonicNotFoundErr || v != "" {
			return true, bytes.NewReader(payload)
		}
	}
	return true, h.Body
}
