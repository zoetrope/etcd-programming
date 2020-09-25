package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cfg := clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 3 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// write value
	_, err = client.Put(context.TODO(), "/chapter3/kv", "my-value")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("wrote")

	// read value
	resp, err := client.Get(context.TODO(), "/chapter3/kv")
	if err != nil {
		log.Fatal(err)
	}
	if resp.Count == 0 {
		log.Fatal("/chapter3/kv not found")
	}
	fmt.Println(string(resp.Kvs[0].Value))

	// delete value
	_, err = client.Delete(context.TODO(), "/chapter3/kv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("deleted")
}
