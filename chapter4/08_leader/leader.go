package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
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

	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("usage: ./leader NAME")
	}
	name := flag.Arg(0)
	s, err := concurrency.NewSession(client, concurrency.WithTTL(30))
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	e := concurrency.NewElection(s, "/chapter4/leader")

	err = e.Campaign(context.TODO(), name)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		fmt.Println(name + " is a leader.")
		time.Sleep(3 * time.Second)
	}
	err = e.Resign(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
