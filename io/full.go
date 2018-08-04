package io

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// TryReadFull uses io.ReadFull to read from response body.
func TryReadFull() error {
	const reqURI = `https://google.com/`
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("get", reqURI, nil)
	if err != nil {
		return fmt.Errorf("can't create request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("can't perform request: %v", err)
	}
	defer resp.Body.Close()
	buf := make([]byte, resp.ContentLength)
	n, err := io.ReadFull(resp.Body, buf)
	if err != nil {
		return fmt.Errorf("can't read from response body: %v", err)
	}
	fmt.Printf("Read from resp.Body: %d bytes", n)
	fmt.Printf("%s", buf)
	return nil
}