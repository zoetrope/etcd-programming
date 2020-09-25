package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func printResponse(resp *clientv3.GetResponse) {
	fmt.Println("header: " + resp.Header.String())
	for i, kv := range resp.Kvs {
		fmt.Printf("kv[%d]: %s\n", i, kv.String())
	}
	fmt.Println()
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

	_, err = client.Delete(context.TODO(), "/chapter3/rev/", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Put(context.TODO(), "/chapter3/rev/1", "123")
	if err != nil {
		log.Fatal(err)
	}
	resp, _ := client.Get(context.TODO(), "/chapter3/rev", clientv3.WithPrefix())
	printResponse(resp)

	_, err = client.Put(context.TODO(), "/chapter3/rev/1", "456")
	if err != nil {
		log.Fatal(err)
	}
	resp, _ = client.Get(context.TODO(), "/chapter3/rev", clientv3.WithPrefix())
	printResponse(resp)

	_, err = client.Put(context.TODO(), "/chapter3/rev/2", "999")
	if err != nil {
		log.Fatal(err)
	}
	resp, _ = client.Get(context.TODO(), "/chapter3/rev", clientv3.WithPrefix())
	printResponse(resp)

	resp, _ = client.Get(context.TODO(), "/chapter3/rev",
		clientv3.WithPrefix(), clientv3.WithRev(resp.Kvs[0].CreateRevision))
	printResponse(resp)
}
