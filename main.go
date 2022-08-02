package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {

	createNAddress(10)
	
	mainAddress, receiptsAddress := getAddress()

	fmt.Printf("Main address : %s\n", mainAddress)
	
	generateTransation(mainAddress, receiptsAddress)
	
	fmt.Printf("Fin\n")
}

func getAddress() (string, []string) {
	body := EthRequestBody{
    Jsonrpc: 2.0,
    Id: 3,
    Method: "eth_accounts",
    Params: []interface{}{},
	}

	receiptBody := sendEthRequest(body)

	mainAddress := receiptBody.Result[0].(string)
	receiptsAddress := []string{}

	for i, address := range(receiptBody.Result) {
		if i == 0 {
			continue
		}

		receiptsAddress = append(receiptsAddress, address.(string))
	}

	return mainAddress, receiptsAddress
}

type EthRequestBody struct {
	Jsonrpc float64 			`json:"jsonrpc"`
	Id 			int 					`json:"id"`
	Method	string				`json:"method"`
	Params	[]interface{}	`json:"params"`
}

type EthReceiptBody struct {
	Jsonrpc string 				`json:"jsonrpc"`
	Id 			int 					`json:"id"`
	Result	[]interface{}	`json:"result"`
}

func createNAddress(x int) {
	var wg sync.WaitGroup

	for i:=0; i<x; i++ {
		body := EthRequestBody{
			Jsonrpc: 2.0,
			Id: 5,
			Method: "personal_newAccount",
			Params: []interface{}{""},
		}

		go func() {
			defer wg.Done()
			wg.Add(1)
			sendEthRequest(body)
		}()
	}
	wg.Wait()
}

func sendEthRequest(body EthRequestBody) EthReceiptBody {

	bufferRequestBody, _ := json.Marshal(body)
	
	// fmt.Printf("%v : %v\n", bufferRequestBody, err)
	// fmt.Printf("%v\n", strings.NewReader(string(bufferRequestBody)))

	resp, err := http.Post("http://localhost:8545", "application/json", strings.NewReader(string(bufferRequestBody)))

	// fmt.Printf("%v : %v\n", resp, err)
	
	defer resp.Body.Close()
	buffferReceiptBody, err := io.ReadAll(resp.Body)
	
	// fmt.Printf("%v : %v\n", buffferReceiptBody, err)
		
	var ReceiptBody EthReceiptBody
	err = json.Unmarshal(buffferReceiptBody, &ReceiptBody)
	
	// fmt.Printf("%v : %v\n", ReceiptBody, err)

	if err != nil {

	}
	
	return ReceiptBody
}

func generateTransation(mainAddress string, receiptAddress []string) {
	
	nb := 10000

	fmt.Printf("Nombre de transaction à envoyer : %v\n", nb*len(receiptAddress))
	
	for i:=0; i<nb; i++{

		unlockAccount(mainAddress, "")
		for _, address := range(receiptAddress) {
			s, _ := randomHex(10000)
			body := EthRequestBody{
				Jsonrpc: 2.0,
				Id: 6,
				Method: "eth_sendTransaction",
				Params: []interface{}{EthSendTransactionParams{
					From: mainAddress,
					To: address,
					Value: "0xf4240",
					Data: fmt.Sprintf("0x%x", s),
				}},
			}
			sendEthRequest(body)
			// 	time.Sleep(1000000000)
					time.Sleep(700000000)
		}
	}
}

func randomHex(n int) (string, error) {
  bytes := make([]byte, n)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }
  return hex.EncodeToString(bytes), nil
}

type EthSendTransactionParams struct {
	From 	string	`json:"from"`
	To 		string	`json:"to"`
	Value string	`json:"value"`
	Data 	string	`json:"data"`
}

func unlockAccount(mainAddress string, mpd string) {
	body := EthRequestBody{
    Jsonrpc: 2.0,
    Id: 6,
    Method: "personal_unlockAccount",
    Params: []interface{}{mainAddress, mpd},
	}

	sendEthRequest(body)
}