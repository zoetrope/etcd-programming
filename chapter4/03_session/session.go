package main

import (
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
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

	session, err := concurrency.NewSession(client, concurrency.WithTTL(10))
	if err != nil {
		log.Fatal(err)
	}
	select {
	case <-session.Done():
		log.Fatal("session has been orphaned")
	}
}
