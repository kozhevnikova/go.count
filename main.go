package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	urls := make(chan string)
	threadpool := make(chan struct{}, 4)
	var totalCount int64

	workers := &sync.WaitGroup{}

	go func() {
		for {
			url := <-urls

			workers.Add(1)

			threadpool <- struct{}{}
			go func(url string) {
				defer workers.Done()
				var countForOneUrl int64
				countForOneUrl = calculateCount(url)
				atomic.AddInt64(&totalCount, countForOneUrl)
				time.Sleep(time.Microsecond)
				<-threadpool
				fmt.Printf("URL: %s, Count: %d \n", url, countForOneUrl)
			}(url)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		match, err := regexp.MatchString(
			`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`, url)
		if match == false {
			fmt.Errorf("ERROR. Check URL", url)
		}
		if err != nil {
			fmt.Errorf("ERROR", err)
		}
		urls <- url
	}

	workers.Wait()

	time.Sleep(time.Second)

	totalCountFinal := atomic.LoadInt64(&totalCount)
	fmt.Printf("Total count %d \n", totalCountFinal)
	fmt.Println("End of the program\n")
}

func calculateCount(url string) int64 {
	var countInInt64 int64
	needToFind := "Go"
	needToFindInBytes := []byte(needToFind)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("ERROR", err)
	}
	textFromPage, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("ERROR", err)
	}
	count := bytes.Count(textFromPage, needToFindInBytes)
	countInInt64 = int64(count)
	resp.Body.Close()
	return countInInt64
}
