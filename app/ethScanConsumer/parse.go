package ethScanConsumer

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type transactionResponse struct {
	Status int           `json:",string"`
	Result []Transaction `json:"result"`
}

type eventResult struct {
	events []Event
	err    error
}

type tokenResult struct {
	tokens []Token
	err    error
}

func ethTransactions(wallet string, from int, to int, internal bool, c chan eventResult) {
	body, err := getEthTransactions(wallet, from, to, internal)
	if err != nil {
		c <- eventResult{events: nil, err: err}
		return
	}
	response := &transactionResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		c <- eventResult{events: nil, err: err}
		return
	}
	transactions := make([]Event, 0)
	if response.Status == 1 {
		for _, t := range response.Result {
			t.Acc = wallet
			if internal {
				t.Type = "INTERNAL"
			} else {
				t.Type = "REGULAR"
			}
			t.ValueEth = stringToEth(t.Value)
			t.GasPriceEth = stringToEth(t.GasPrice)
			transactions = append(transactions, t)
		}

	}
	c <- eventResult{events: transactions, err: nil}
	return
}

//*****************************

type blockTimeResponse struct {
	Status int `json:",string"`
	Result int `json:",string"`
}

func blockAtDate(timestamp int64) (int, error) {
	body, err := getLastBlockBefore(timestamp)
	if err != nil {
		return -1, err
	}
	response := &blockTimeResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return -1, err
	}
	if response.Status == 0 {
		return -1, errors.New("Block not found")
	}
	return response.Result, nil
}

func blockAfterDate(timestamp int64) (int, error) {
	body, err := getBlockAtTimestamp(timestamp, false)
	if err != nil {
		return -1, err
	}
	response := &blockTimeResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return -1, err
	}
	if response.Status == 0 {
		return -1, errors.New("Block not found")
	}
	return response.Result, nil
}

type balanceResult struct {
	balance float64
	err     error
}

type balanceResponse struct {
	Status int `json:",string"`
	Result string
}

func balanceOfAccount(address string, c chan balanceResult) {
	body, err := getCrrEthBalance(address)
	if err != nil {
		c <- balanceResult{balance: 0, err: err}
		return
	}
	response := &balanceResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		c <- balanceResult{balance: 0, err: err}
		return
	}
	if response.Status == 0 {
		c <- balanceResult{balance: 0, err: errors.New("Account not found")}
		return
	}
	c <- balanceResult{balance: stringToEth(response.Result), err: nil}
	return
}

//********************************************************

type minedBlocksResponse struct {
	Status int    `json:",string"`
	Result []Mine `json:"result"`
}

func minedBlocks(address string, uncle bool, c chan eventResult) {
	body, err := getMinedBlocks(address, uncle)
	if err != nil {
		c <- eventResult{events: nil, err: err}
		return
	}
	response := &minedBlocksResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		c <- eventResult{events: nil, err: err}
		return
	}
	mined := make([]Event, 0)
	if response.Status == 1 {
		for _, m := range response.Result {
			m.Acc = address
			if uncle {
				m.Type = "UNCLE"
			} else {
				m.Type = "REGULAR"
			}
			m.RewardEth = stringToEth(m.BlockReward)
			mined = append(mined, m)
		}
	}
	c <- eventResult{events: mined, err: nil}
	return
}

//********************************************************

type ecr20Response struct {
	Status int     `json:",string"`
	Result []ECR20 `json:"result"`
}

func ecr20Transactions(address string, from int, to int, c chan tokenResult) {
	body, err := getERC20Transactions(address, from, to)
	if err != nil {
		c <- tokenResult{tokens: nil, err: err}
		return
	}
	response := &ecr20Response{}
	err = json.Unmarshal(body, response)
	if err != nil {
		c <- tokenResult{tokens: nil, err: err}
		return
	}
	tokens := make([]Token, 0)
	if response.Status == 1 {
		for _, tkn := range response.Result {
			tkn.Acc = address
			tkn.ValueNum = stringToVal(tkn.Value, tkn.TokenDecimal)
			tkn.GasPriceEth = stringToEth(tkn.GasPrice)
			tokens = append(tokens, tkn)
		}
	}
	c <- tokenResult{tokens: tokens, err: nil}
	return
}

type ecr271Response struct {
	Status int      `json:",string"`
	Result []ECR271 `json:"result"`
}

func ecr271Transactions(address string, from int, to int, c chan tokenResult) {
	body, err := getERC721Transactions(address, from, to)
	if err != nil {
		c <- tokenResult{tokens: nil, err: err}
		return
	}
	response := &ecr271Response{}
	err = json.Unmarshal(body, response)
	if err != nil {
		c <- tokenResult{tokens: nil, err: err}
		return
	}
	tokens := make([]Token, 0)
	if response.Status == 1 {
		for _, tkn := range response.Result {
			tkn.Acc = address
			tkn.GasPriceEth = stringToEth(tkn.GasPrice)
			tokens = append(tokens, tkn)
		}
	}
	c <- tokenResult{tokens: tokens, err: nil}
	return
}

//********************************************************

func stringToEth(val string) float64 {
	return stringToVal(val, 18)
}

func stringToVal(val string, dec int) float64 {
	if len(val) <= dec {
		val = "0." + strings.Repeat("0", dec-len(val)) + val
	} else {
		val = val[:len(val)-dec] + "." + val[len(val)-dec:]
	}
	res, _ := strconv.ParseFloat(val, 64)
	return res
}
