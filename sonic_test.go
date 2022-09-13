package main

import (
	"github.com/bytedance/sonic"
	"log"
	"testing"
)

func TestSonicReadErr(t *testing.T) {
	json := `{"error":"true"}`
	json2 := `{"test": 1}`
	non_json := `3fnivbd`
	_, err := sonic.Get([]byte(json), "error")
	log.Println("test1", err)
	_, err = sonic.Get([]byte(json2), "error")
	log.Println("test2", err)
	_, err = sonic.Get([]byte(non_json), "error")
	log.Println("test3", err)
}
