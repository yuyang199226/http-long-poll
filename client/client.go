package main

import (
	"fmt"
	"io"
	"net/http"
)

const appId = "111"

func main() {
	for {
		fmt.Println("start http long poll")
		url := fmt.Sprintf("http://127.0.0.1:8081/get_config?app_id=%s", appId)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}

		if resp.StatusCode == 304 {
			fmt.Println("content not change")
		} else if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(body))
		}
		resp.Body.Close()

	}
}
