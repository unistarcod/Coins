package api

import (
	"fmt"
	"net/http"
	"strings"
)

func createwallet() {
	fmt.Println("wallet")
	body := strings.NewReader("{\"jsonrpc\":\"1.0\",\"id\":\"curltext\",\"method\":\"getblockchaininfo\",\"params\":[]}")
	req, err := http.NewRequest("POST", "http://127.0.0.1:8332", body)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth("rpcuser", "rpcpassword")
	req.Header.Set("Content-Type", "text/plain;")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	fmt.Println(resp.Proto)
	defer resp.Body.Close()

}
