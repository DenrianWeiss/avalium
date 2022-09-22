package ws

import (
	"encoding/json"
	"github.com/DenrianWeiss/avalium/service/forward"
	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"strings"
	"sync"
)

var StatefulCall = map[string]bool{
	"eth_subscribe":   true,
	"eth_unsubscribe": true,
}

type RpcMsg struct {
	JsonRPC string        `json:"jsonrpc"`
	Id      interface{}   `json:"id,omitempty"`
	Method  string        `json:"method,omitempty"`
	Params  []interface{} `json:"params"`
	Result  interface{}   `json:"result,omitempty"`
	Error   interface{}   `json:"error,omitempty"`
}

type RpcResp struct {
	JsonRPC string      `json:"jsonrpc"`
	Id      interface{} `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type Upstream struct {
	conn             *websocket.Conn
	idMapping        map[string]string // Map downstream sub id to upstream
	reverseIdMapping map[string]string // Map upstream id to downstream
	HistorySub       map[string]string
	idLock           sync.RWMutex
	recoverTime      int
}

func ForwardStatelessReq(req interface{}) []byte {
	r, _ := sonic.Marshal(req)
	resp := forward.RpcPullUntilAvailable(r)
	respStr, err := io.ReadAll(resp)
	if err != nil {
		log.Println("Websocket state forward error", err)
	}
	return respStr
}

func ReceiveFromSocket(conn *websocket.Conn, ch chan RpcResp, errCh chan error) {
	for {
		v := RpcResp{}
		err := conn.ReadJSON(&v)
		if err != nil {
			errCh <- err
			return
		}
		ch <- v
	}
}

func (u *Upstream) ReceiveRoutine(downstream *websocket.Conn, closeCh chan bool) {
	receiveCh := make(chan RpcResp)
	defer close(receiveCh)
	errorCh := make(chan error)
	defer close(errorCh)
	go ReceiveFromSocket(u.conn, receiveCh, errorCh)
	for {
		select {
		case <-closeCh:
			{
				u.conn.Close()
				return
			}
		case <-errorCh:
			{
				/// Todo Handle error
				u.conn.Close()
				return
			}
		case m := <-receiveCh:
			{
				u.idLock.RLock()
				// Replace subId
				// Find if event is subscription
				if m.Method == "eth_subscription" {
					// Replace subscription id
					r := m.Params.(map[string]interface{})
					v, ok := r["subscription"]
					if ok {
						r["subscription"] = u.reverseIdMapping[v.(string)]
					}
				} else {
					if m.Id != nil {
						switch m.Id.(type) {
						case string:
							{
								id := m.Id.(string)
								if strings.HasPrefix(id, "subscribe_") {

								} else if strings.HasPrefix(id, "unsubscribe_") {

								}
							}
						default:

						}
					}
				}
				// Send Message to Deduplicate
				str, err := json.Marshal(m)
				if err != nil {
					log.Println(err.Error())
				} else {
					if !IsDup(string(str)) {
						err := downstream.WriteMessage(websocket.TextMessage, str)
						if err != nil {
							log.Println("send message error")
						}
					}
				}

				u.idLock.RUnlock()
			}
		}
	}
}

func (u *Upstream) CloseSubscribe(id string) {
	u.idLock.Lock()
	// todo
	defer u.idLock.RUnlock()
}

func (u *Upstream) New(conn *websocket.Conn) {
	u.conn = conn
	u.HistorySub = make(map[string]string)
	u.idMapping = make(map[string]string)
	u.reverseIdMapping = make(map[string]string)
	u.idLock = sync.RWMutex{}
}

func (u *Upstream) NewSubscribe(req RpcMsg, subId string) {
	// Todo Convert req id back
	u.idLock.Lock()
	defer u.idLock.Unlock()
	req.Id = "subscribe_" + subId
	err := u.conn.WriteJSON(req)
	if err != nil {
		u.conn.Close()
		return
	}
}

func (u *Upstream) Unsubscribe(req RpcMsg, subId string) {
	// Todo Convert req id back
	u.idLock.Lock()
	defer u.idLock.Unlock()
	req.Id = "unsubscribe_" + subId
	// Convert subId to local sub id
	// Remove subId mapping
	// Remove subId in history
}

func HandleWs(conn *websocket.Conn, lock *sync.Mutex) {
	// Connect to upstream
	for {
		req := RpcMsg{}
		err := conn.ReadJSON(req)
		if err != nil {
			log.Println(err)
			if err != io.ErrUnexpectedEOF {
				continue
			}
			_ = conn.Close()
			return
		}
		// Check If Request is stateful
		method := req.Method
		_, ok := StatefulCall[method]
		if ok {
			// Stateful Handler

		} else {
			// Forward to random upstream
			lock.Lock()
			err := conn.WriteMessage(websocket.TextMessage, ForwardStatelessReq(req))
			if err != nil {
				log.Println("websocket write err", err)
			}
			lock.Unlock()
		}
	}
}
