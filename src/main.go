package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	workers "github.com/jrallison/go-workers"
	phantomjs "github.com/nladuo/go-phantomjs-fetcher"
)

// Contents :
type Contents struct {
	// {"args":[1,2],"at":1502818316.2218475,"class":"Add","enqueued_at":1502818316.2218478,"jid":"704ea5578be22bb1743c2777","queue":"myqueue3"}
	Args []string `json:"args"`
	// At         float64 `json:"at"`
	// Class      string  `json:"class"`
	// EnqueuedAt float64 `json:"enqueued_at"`
	// Jid        string  `json:"jid"`
	// Queue      string  `json:"queue"`
}

// getRandomNumber:
func getRandomNumber() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(3000) + 21000
}

// phantomWorker :
func phantomWorker(message *workers.Msg) {
	fmt.Println(message.Jid())
	// message.Args() is a wrapper around go-simplejson (http://godoc.org/github.com/bitly/go-simplejson)
	a := message.ToJson()

	res := Contents{}
	json.Unmarshal([]byte(a), &res)

	port := getRandomNumber()
	loadURL(&res.Args[0], &res.Args[1], &port, message.Jid())

}

// loadURL :
func loadURL(url *string, checkContent *string, port *int, jid string) {
	fmt.Printf("JID(%s) - Starting PhantomJS\n", jid)
	fetcher, err := phantomjs.NewFetcher(*port, nil)
	fmt.Printf("JID(%s) -Started PhantomJS\n", jid)
	checkErr(err)
	defer fetcher.ShutDownPhantomJSServer()

	jsRunAt := phantomjs.RUN_AT_DOC_END
	resp, err := fetcher.GetWithJS(*url, "", jsRunAt)
	checkErr(err)
	fmt.Printf("JID(%s) - Success: %t\n", jid, strings.Contains(resp.Content, *checkContent))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type myMiddleware struct{}

func (r *myMiddleware) Call(queue string, message *workers.Msg, next func() bool) (acknowledge bool) {
	// do something before each message is processed
	acknowledge = next()
	// do something after each message is processed
	return
}

func main() {
	workers.Configure(map[string]string{
		// location of redis instance
		"server": os.Getenv("REDIS_URL"),
		// instance of the database
		"database": "0",
		// number of connections to keep open with redis
		"pool": "30",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "1",
	})

	workers.Middleware.Append(&myMiddleware{})

	if !(os.Getenv("IS_WORKER") == "true") {
		// Add a job to a queue with retry
		for i := 0; i < 9; i++ {
			// Add a job to a queue
			go workers.EnqueueWithOptions("phantomQueue", "phantomWorker", []string{"https://www.google.com/", "Feeling Lucky"}, workers.EnqueueOptions{Retry: true})
		}
		// stats will be available at http://localhost:8081/stats
		workers.StatsServer(8081)
	} else {
		// pull messages from "phantomQueue" with concurrency of 20
		numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
		checkErr(err)
		workers.Process("phantomQueue", phantomWorker, numWorkers)
	}
	// Blocks until process is told to exit via unix signal
	workers.Run()
}
