package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
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

	_, err = client.Delete(context.TODO(), "/chapter3/compaction")
	if err != nil {
		log.Fatal(err)
	}

	// prepare data
	_, err = client.Put(context.TODO(), "/chapter3/compaction", "hoge")
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Put(context.TODO(), "/chapter3/compaction", "fuga")
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Put(context.TODO(), "/chapter3/compaction", "fuga")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Get(context.TODO(), "/chapter3/compaction")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("--- prepared data: ")
	for i := resp.Kvs[0].CreateRevision; i <= resp.Kvs[0].ModRevision; i++ {
		r, err := client.Get(context.TODO(), "/chapter3/compaction", clientv3.WithRev(i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("rev: %d, value: %s\n", i, r.Kvs[0].Value)
	}

	// compaction
	_, err = client.Compact(context.TODO(), resp.Kvs[0].ModRevision)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n--- compacted: %d\n", resp.Kvs[0].ModRevision)
	for i := resp.Kvs[0].CreateRevision; i <= resp.Kvs[0].ModRevision; i++ {
		r, err := client.Get(context.TODO(), "/chapter3/compaction", clientv3.WithRev(i))
		if err != nil {
			fmt.Printf("failed to get: %v\n", err)
			continue
		}
		fmt.Printf("rev: %d, value: %s\n", i, r.Kvs[0].Value)
	}
}
