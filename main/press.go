package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var counter int = 0

type httpConnect interface {
	onetest()
}

func onetest(url string, verify string, lock *sync.Mutex) int64 {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	_ = client

	t1 := time.Now()

	if strings.Index(url, "https") > -1 {
		var resp *http.Response
		if verify == "on" {
			resp, _ = http.Get(url)
		} else {
			resp, _ = client.Get(url)
		}

		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			_, _ = ioutil.ReadAll(resp.Body)
		}
	} else {
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			_, _ = ioutil.ReadAll(resp.Body)
		}
	}

	//resp, err := http.Get("https://meituan.51hutui.com/hello")
	//resp, err := http.Get("http://meituan.51hutui.com/hello")

	//resp, err := client.Get("https://120.197.235.100:6243/mtmid/axb/binding")
	//resp, err := http.Get("http://120.197.235.100:6280/mtmid/axb/binding")

	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()

	//if resp.StatusCode == http.StatusOK {
	//	_, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}

	t2 := time.Now()
	d := t2.Sub(t1)
	lock.Lock()
	counter++
	fmt.Printf("这是第%d次", counter)
	lock.Unlock()
	//fmt.Println(d)
	//fmt.Println(int64(d))
	return int64(d) / 1000000
}

const capcatity = 1000
const len = capcatity / 100

func main() {
	args := os.Args
	verify := args[1]
	threads, _ := strconv.Atoi(args[2])
	lock := &sync.Mutex{}
	url := args[3]
	fmt.Println("verify https ca:", verify)
	fmt.Println("concurrent threads:", threads)
	fmt.Println("url:", url)

	c := make(chan int64, capcatity)
	for k := 0; k < threads; k++ {
		go func(c chan int64) {
			for true {
				c <- onetest(url, verify, lock)
			}
		}(c)
	}

	step := 0
	for true {
		step = step + 1
		t1 := time.Now()
		var total int64 = 0

		var store []int64 = make([]int64, len)

		for i := 0; i < capcatity; i++ {
			t := <-c
			total = total + t

			if t > store[0] {
				store[0] = t
				for i := 0; i < len-1; i++ {
					if store[i] < store[i+1] {
						break
					} else {
						temp := store[i+1]
						store[i+1] = store[i]
						store[i] = temp
					}
				}
			}

		}
		t2 := time.Now()
		d := t2.Sub(t1)
		//if store[0] < 60{
		//	fmt.Print("*")
		//}else if store[0] < 100{
		//	fmt.Print("#")
		//}else{
		//	fmt.Print("\n")
		//fmt.Println("the max 1%=", store)
		if step > 1 {
			fmt.Print("avg: ", float64(total)/capcatity, "ms")
			fmt.Println(" qps:", capcatity/(float64(d)/1000000000.0), " tp99:", store[0], "ms")
		} else {
			fmt.Print("没有走这条路")
		}

		//}

	}

}
