package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func addValue(client *clientv3.Client, key string, d int) {
	_, err := concurrency.NewSTM(client, func(stm concurrency.STM) error {
		v := stm.Get(key)
		value, _ := strconv.Atoi(v)
		value += d
		stm.Put(key, strconv.Itoa(value))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	key := "/chapter4/stm"
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
