package main

import (
	"fmt"
	"net/http"
	"sync"
)

func getFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}
	return resp.ContentLength, nil
}

func downloadInRange(url string, start, end int64, wg *sync.WaitGroup) {
	defer wg.Done()
	
}
func main() {

}
