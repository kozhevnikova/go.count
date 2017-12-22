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
	"time"
)

func main() {
	urls := make(chan string)
	threadpool := make(chan struct{}, 4)
	var totalCount int

	workers := &sync.WaitGroup{}

	go func() {
		for {
			url := <-urls

			workers.Add(1)

			threadpool <- struct{}{}
			go func(url string) {
				defer workers.Done()
				countForOneUrl := calculateCount(url)
				totalCount += countForOneUrl
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

	fmt.Println("Total", totalCount)
	fmt.Println("End of the program\n")
}

func calculateCount(url string) int {
	need := "Go"
	needBytes := []byte(need)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("ERROR", err)
	}
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("ERROR", err)
	}
	count := bytes.Count(text, needBytes)
	return count
}
