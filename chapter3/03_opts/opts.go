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

	_, err = client.Put(context.TODO(), "/chapter3/option/key3", "val2")
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Put(context.TODO(), "/chapter3/option/key1", "val3")
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Put(context.TODO(), "/chapter3/option/key2", "val1")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Get(context.TODO(), "/chapter3/option/",
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByValue, clientv3.SortAscend),
		clientv3.WithKeysOnly(),
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, kv := range resp.Kvs {
		fmt.Printf("%s: %s\n", kv.Key, kv.Value)
	}
}
