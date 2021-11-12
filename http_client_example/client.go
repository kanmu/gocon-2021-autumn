package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	req, err := http.NewRequest("GET", "http://localhost:8000/flag.txt", nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bb := &bytes.Buffer{}
	io.Copy(bb, resp.Body)
	fmt.Printf("%s\n", bb.String())
}
