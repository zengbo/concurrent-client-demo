package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var apiPrefix string

type Response struct {
	Result string `json:"result"`
}

func createCacheTag(keyName string) {
	url := apiPrefix + "/cache_tag"
	jsonData := map[string]string{"key": keyName}
	jsonValue, _ := json.Marshal(jsonData)
	_, _ = http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
}

func deleteCacheTag() {
	url := apiPrefix + "/cache_tag"
	req, _ := http.NewRequest("DELETE", url, nil)
	client := &http.Client{}
	_, _ = client.Do(req)
}

func checkCacheTag() {
	url := apiPrefix + "/cache_tag/check"
	res, err := http.Get(url)
	if err != nil {
		// Handle error according to your needs
		fmt.Println("Error in GET request:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		// Handle error here
		fmt.Println("Error reading body:", err)
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		// Handle error here
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	if response.Result == "inconsistent" {
		fmt.Println("Inconsistent result found")
		os.Exit(1)
	}

}

func main() {
	var iterations int
	flag.StringVar(&apiPrefix, "apiPrefix", "http://localhost/api", "API prefix URL")
	flag.IntVar(&iterations, "iterations", 50, "The number of iterations of the test")
	flag.Parse()

	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		wg.Add(21)
		for j := 0; j < 20; j++ {
			go func(j int) {
				createCacheTag("key_" + strconv.Itoa(j))
				wg.Done()
			}(j)
		}
		go func() {
			deleteCacheTag()
			wg.Done()
		}()
		wg.Wait()
		checkCacheTag()
		time.Sleep(time.Second * 5)
	}
}
