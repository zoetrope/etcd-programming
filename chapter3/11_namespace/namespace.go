package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/namespace"
)

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	newClient := namespace.NewKV(client.KV, "/chapter3")

	_, err = newClient.Put(context.TODO(), "/ns/1", "hoge")
	if err != nil {
		log.Fatal(err)
	}
	resp, _ := client.Get(context.TODO(), "/chapter3/ns/1")
	fmt.Printf("%s: %s\n", resp.Kvs[0].Key, resp.Kvs[0].Value)

	_, err = client.Put(context.TODO(), "/chapter3/ns/2", "test")
	if err != nil {
		log.Fatal(err)
	}
	resp, _ = newClient.Get(context.TODO(), "/ns/2")
	fmt.Printf("%s: %s\n", resp.Kvs[0].Key, resp.Kvs[0].Value)

	client.KV = namespace.NewKV(client.KV, "/chapter3")
	client.Watcher = namespace.NewWatcher(client.Watcher, "/chapter3")
	client.Lease = namespace.NewLease(client.Lease, "/chapter3")
}
