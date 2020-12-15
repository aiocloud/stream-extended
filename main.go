package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	verCode = "1.0.4"

	flags struct {
		Path    string
		VerCode bool
	}

	Data []struct {
		Domain []string `json:"domains"`
		Remote []string `json:"remotes"`
	}
)

func main() {
	flag.StringVar(&flags.Path, "c", "", "Path")
	flag.BoolVar(&flags.VerCode, "v", false, "VerCode")
	flag.Parse()

	if flags.VerCode {
		fmt.Println(verCode)
		return
	}

	{
		data, err := ioutil.ReadFile(flags.Path)
		if err != nil {
			log.Fatalf("[APP][ioutil.ReadFile] %v", err)
		}

		if err := json.Unmarshal(data, &Data); err != nil {
			log.Fatalf("[APP][json.Unmarshal] %v", err)
		}
	}

	go startHTTP()
	go startTLS()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
	<-channel
}
