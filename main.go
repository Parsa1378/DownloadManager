package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func getFileSize(url string) (int64, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}
	return resp.ContentLength, nil
}

func downloadInRange(url string, start, end int64, wg *sync.WaitGroup, filename string) {
	defer wg.Done()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error downloding range %d-%d : %s\n", start, end, err)
		return
	}
	defer res.Body.Close()
	var saveFilename string
	if filename == "" {
		saveFilename = filepath.Base(url)
	} else {
		saveFilename = filename
	}

	file, err := os.OpenFile(saveFilename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error in opening the file: %s", saveFilename)
		return
	}
	defer file.Close()

	file.Seek(start, io.SeekStart)
	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Printf("Error in writing in the file: %s\n", err)
		return
	}
}
func main() {
	start := time.Now()
	var filename string
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: dm <URL> [options]")
		return
	}
	if len(args) > 1 {
		filename = args[2]
	}
	url := args[0]
	fileSize, err := getFileSize(url)
	fmt.Println(fileSize)
	if err != nil {
		fmt.Print(err)
		return
	}

	const chunkSize = 50 * 1024 * 1024 //50 MB
	nChunk := fileSize / chunkSize
	fmt.Println(nChunk)
	var wg sync.WaitGroup
	for i := int64(0); i < nChunk; i++ {
		start := i * chunkSize
		end := chunkSize * (i + 1)
		if end > fileSize {
			end = fileSize
		}
		wg.Add(1)
		go downloadInRange(url, start, end-1, &wg, filename)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Download Completed in:%v\n", elapsed)
}
