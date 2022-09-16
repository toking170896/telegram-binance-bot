package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Validation struct {
	Address string `json:"address,omitempty"`
	Valid   bool   `json:"valid,omitempty"`
}

type Wallet struct {
	Address       string  `json:"address,omitempty"`
	Confirmations int     `json:"confirmations,omitempty"`
	Balance       float64 `json:"balance,omitempty"`
}

type GetLastTransactionRes [][]Transaction

type Transaction struct {
	Type            int      `json:"type"`
	Version         int      `json:"version"`
	ID              string   `json:"id"`
	Sender          string   `json:"sender"`
	SenderKeyType   string   `json:"senderKeyType"`
	SenderPublicKey string   `json:"senderPublicKey"`
	Fee             int      `json:"fee"`
	Timestamp       int64    `json:"timestamp"`
	Recipient       string   `json:"recipient"`
	Amount          int64    `json:"amount"`
	Proofs          []string `json:"proofs"`
	Status          string   `json:"status"`
	Height          int      `json:"height"`
	EffectiveFee    int      `json:"effectiveFee"`
}

const ltoAPI = "http://144.76.119.169:6869"

func ValidateAddress(address string) (*Validation, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", ltoAPI+"/addresses/validate/"+address, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res *Validation
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetWalletBalance(address string) (*Wallet, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", ltoAPI+"/addresses/balance/"+address, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res *Wallet
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetLastTransaction(address string) (*Transaction, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", ltoAPI+"/transactions/address/"+address+"/limit/1", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res GetLastTransactionRes
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res[0][0], nil
}
