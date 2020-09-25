package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func addValue(client *clientv3.Client, key string, d int) {
	resp, _ := client.Get(context.TODO(), key)
	value, _ := strconv.Atoi(string(resp.Kvs[0].Value))
	value += d
	_, err := client.Put(context.TODO(), key, strconv.Itoa(value))
	if err != nil {
		log.Fatal(err)
	}
}

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

	key := "/chapter4/conflict"
	_, err = client.Put(context.TODO(), key, "10")
	if err != nil {
		log.Fatal(err)
	}
	go addValue(client, key, 5)
	go addValue(client, key, -3)
	time.Sleep(1 * time.Second)
	resp, _ := client.Get(context.TODO(), key)
	fmt.Println(string(resp.Kvs[0].Value))
}
