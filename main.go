package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func main() {
	client := http.Client{}
	resp, err := client.Get("http://192.168.1.184")
	if err != nil {
		log.Println(err)
	}
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		log.Println(err)
	}
	go deleteOld()
	for {
		mr := multipart.NewReader(resp.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Println(err)
			}
			slurp, err := io.ReadAll(p)
			if err != nil {
				log.Println(err)
			}
			f, err := os.Create(fmt.Sprintf("data/%s.jpg", time.Now().Round(time.Millisecond*100)))
			if err != nil {
				log.Println(err)
			}
			_, err = io.Copy(f, bytes.NewReader(slurp))
			if err != nil {
				log.Println(err)
			}
			f.Close()
		}
	}
}

func deleteOld() {
	for {
		files, err := ioutil.ReadDir("data")
		if err != nil {
			log.Println(err)
		}
		now := time.Now()
		for _, fileInfo := range files {
			if diff := now.Sub(fileInfo.ModTime()); diff > time.Hour*48 {
				err := os.Remove(fmt.Sprintf("data/%s", fileInfo.Name()))
				if err != nil {
					log.Println(err)
				}
				fmt.Printf("Deleting %s which is %s old\n", fileInfo.Name(), diff)
			}
		}
	}
}
