package main

import (
	"context"
	"fmt"
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

	session, err := concurrency.NewSession(client)
	if err != nil {
		log.Fatal(err)
	}

	mutex := concurrency.NewMutex(session, "/chapter4/mutex_txn")
RETRY:
	select {
	case <-session.Done():
		log.Fatal("session has been orphaned")
	default:
	}
	err = mutex.Lock(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Txn(context.TODO()).
		If(mutex.IsOwner()).
		Then(clientv3.OpPut("/chapter4/mutex_owner", "test")).
		Commit()
	if err != nil {
		log.Fatal(err)
	}
	if !resp.Succeeded {
		fmt.Println("the lock was not acquired")
		mutex.Unlock(context.TODO())
		goto RETRY
	}
	err = mutex.Unlock(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
