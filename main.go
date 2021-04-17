package main

import (
	"log"
	"sync"
	"time"

	"gocnnews/channel"
	"gocnnews/model"
	"gocnnews/server"
)

func main() {
	go server.Start()
	model.Init()
	for {
		runChannels()
		time.Sleep(time.Minute * 10)
	}
}

func runChannels() {
	st := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(channel.Channels))
	for _, fn := range channel.Channels {
		go func(fn channel.ChannelFunc) {
			model.SaveArticles(fn())
			wg.Done()
		}(fn)
	}
	wg.Wait()
	log.Println("[info]\t channel,", time.Since(st))
}
