package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type PutArgs struct {
	Key   string
	Value string
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Value string
}

type PutReply struct {
}

type KV struct {
	mu    sync.Mutex
	store map[string]string
}

func (store *KV) Get(args *GetArgs, replyArgs *GetReply) error {

	store.mu.Lock()
	value, ok := store.store[args.Key]

	replyArgs.Value = value

	store.mu.Unlock()

	if !ok {
		return errors.New("Some error occured in Getting from the key value store")
	}

	return nil
}

func (store *KV) Put(args *PutArgs, replyArgs *PutReply) error {

	store.mu.Lock()
	store.store[args.Key] = args.Value
	store.mu.Unlock()

	return nil
}

func server() {

	keyValue := new(KV)
	keyValue.store = make(map[string]string)

	rpc.Register(keyValue)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	go http.Serve(l, nil)
}

func get(args *GetArgs) {
	serverAddress := "localhost:1234"
	client, err := rpc.DialHTTP("tcp", serverAddress)

	if err != nil {
		log.Fatal("dialing:", err)
	}
	reply := GetReply{}

	err = client.Call("KV.Get", args, &reply)
	fmt.Printf("Get %s : %s \n", "name", reply.Value)
}

func put(args *PutArgs) {

	serverAddress := "localhost:1234"
	client, err := rpc.DialHTTP("tcp", serverAddress)

	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply *PutReply

	err = client.Call("KV.Put", &PutArgs{Key: args.Key, Value: args.Value}, &reply)

}

func main() {
	server()

	put(&PutArgs{Key: "name", Value: "sagar"})
	get(&GetArgs{Key: "name"})
	put(&PutArgs{Key: "name", Value: "diptilal"})
	get(&GetArgs{Key: "name"})
}
