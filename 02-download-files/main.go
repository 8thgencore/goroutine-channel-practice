package main

import (
	"errors"
	"fmt"
	"sync"
)

// fakeDownload - функция для имитации загрузки (пример)
func fakeDownload(url string) error {
	if url == "http://error.com" {
		return fmt.Errorf("failed to download: %s", url)
	}
	fmt.Printf("Downloaded: %s\n", url)
	return nil
}

func download(urls []string) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if err := fakeDownload(url); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}

		}(url)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func main() {
	urls := []string{
		"http://example.com/file1",
		"http://example.com/file2",
		"http://error.com",
		"http://example.com/file3",
	}

	if err := download(urls); err != nil {
		fmt.Println("Errors occurred:")
		fmt.Println(err)
	} else {
		fmt.Println("All downloads completed successfully!")
	}
}
