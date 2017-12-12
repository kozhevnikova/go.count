package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getUrls() []string {
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		url := scanner.Text()
		urls = append(urls, url)
	}
	return urls
}

func calculateCount(urls []string) int {
	var count int
	need := "Go"
	needBytes := []byte(need)
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Errorf("ERROR", err)
		}
		text, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("ERROR", err)
		}
		count = bytes.Count(text, needBytes)
		fmt.Printf("Count for %s is %d \n", url, count)
	}
	return count
}

func main() {
	urls := getUrls()
	count := calculateCount(urls)
	fmt.Println(count)
}
