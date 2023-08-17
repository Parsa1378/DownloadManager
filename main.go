package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func downloadInRange(url string, start, end int64, wg *sync.WaitGroup, filename string) {
	defer wg.Done()
	client := *&http.Client{}
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

	//create or open the file
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
	var filename string

	flag.StringVar(&filename, "-filename", "", "Name of the file to save as (optional)")
	flag.Parse()

	args := flag.Args()
	// if len(args) != 1 {

	// }
	url := args[0]
	fileSize, err := getFileSize(url)
	if err != nil {
		fmt.Print(err)
		return
	}

	const chunkSize = 10 * 1024 * 1024 //10 MB
	nChunk := chunkSize / fileSize
	var wg *sync.WaitGroup
	for i := int64(0); i < int64(nChunk); i++ {
		start := i * chunkSize
		end := chunkSize * (i + 1)
		if end > fileSize {
			end = fileSize
		}
		wg.Add(1)
		go downloadInRange(url, start, end, wg, filename)
	}
	wg.Wait()
	fmt.Printf("Download Completed")
}
