package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
)

type Data struct {
	Duration string
}

func runChecker(address string, h *hub) {
	for _ = range time.Tick(10 * time.Second) {
		now := time.Now()
		resp, err := http.Get(address)
		defer resp.Body.Close()

		if err != nil {
			glog.Errorf("Cant get content from %s, err: %v", address, err)
		} else {
			elapsed := time.Since(now)
			glog.Infof("Request ok, took: %s", elapsed)
			duration := Data{fmt.Sprintf("%d", int64(elapsed/time.Millisecond))}

			resp, err := json.Marshal(duration)
			if err != nil {
				glog.Error("Cant marshal json: %v", err)
				return
			}

			h.broadcast <- resp
		}
	}
}
