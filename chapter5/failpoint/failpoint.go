package failpoint

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func nextRev() int64 {
	p := "./last_revision"
	f, err := os.Open(p)
	if err != nil {
		os.Remove(p)
		return 0
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		os.Remove(p)
		return 0
	}
	rev, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		os.Remove(p)
		return 0
	}
	return rev + 1
}

func saveRev(rev int64) error {
	p := "./last_revision"
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(strconv.FormatInt(rev, 10))
	return err
}

// Run runs sample code for failpoint
func Run() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	rev := nextRev()
	fmt.Printf("loaded revision: %d\n", rev)
	ch := client.Watch(context.TODO(), "/chapter5/watch_file", clientv3.WithRev(rev))
	for resp := range ch {
		if resp.Err() != nil {
			log.Fatal(resp.Err())
		}
		for _, ev := range resp.Events {
			doSomething(ev)
			if vExampleOneLine, __fpErr := __fp_ExampleOneLine.Acquire(); __fpErr == nil {
				defer __fp_ExampleOneLine.Release()
				_, __fpTypeOK := vExampleOneLine.(struct{})
				if !__fpTypeOK {
					goto __badTypeExampleOneLine
				}
			__badTypeExampleOneLine:
				__fp_ExampleOneLine.BadType(vExampleOneLine, "struct{}")
			}
			err := saveRev(ev.Kv.ModRevision)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("saved: %d\n", ev.Kv.ModRevision)
		}
	}
}

func doSomething(ev *clientv3.Event) {
	fmt.Printf("[%d] %s %q : %q\n", ev.Kv.ModRevision, ev.Type, ev.Kv.Key, ev.Kv.Value)
}
