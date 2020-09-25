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

	lease, err := client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Put(context.TODO(), "/chapter3/lease", "value",
		clientv3.WithLease(lease.ID))
	if err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := client.Get(context.TODO(), "/chapter3/lease")
		if err != nil {
			log.Fatal(err)
		}
		if resp.Count == 0 {
			fmt.Println("'/chapter3/lease' disappeared")
			break
		}
		fmt.Printf("[%v] %s\n",
			time.Now().Format("15:04:05"), resp.Kvs[0].Value)
		time.Sleep(1 * time.Second)
	}
}
