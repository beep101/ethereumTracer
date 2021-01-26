package ethScanConsumer

import (
	"fmt"
	"time"
)

type Event interface {
	Report() string
	Effect() float64
	Block() int
	Time() int64
	IsValid() bool
}

type Transaction struct {
	Acc         string `json:"-"`
	Type        string `json:"-"`
	From        string
	To          string
	BlockNumber int   `json:",string"`
	TimeStamp   int64 `json:",string"`
	Value       string
	ValueEth    float64 `json:"-"`
	IsError     int     `json:",string"`
	GasPrice    string
	GasPriceEth float64 `json:"-"`
	GasUsed     int     `json:",string"`
}

func (t Transaction) Report() string {
	resp := time.Unix(t.TimeStamp, 0).String()[:16]
	resp += " : "
	if t.Acc == t.To {
		resp += " RECEIVED FROM "
		resp += t.From
	} else {
		resp += " SENT TO "
		resp += t.To
	}
	resp += " = "
	resp += fmt.Sprint(t.ValueEth)
	resp += " ETH ; "
	resp += t.Type
	resp += " TRANSACTION"
	return resp
}

func (t Transaction) Block() int {
	return t.BlockNumber
}

func (t Transaction) Effect() float64 {
	if t.Acc == t.To {
		if !t.IsValid() {
			return 0
		}
		return t.ValueEth
	} else {
		if !t.IsValid() {
			return 0 - t.GasPriceEth*float64(t.GasUsed)
		}
		return 0 - t.ValueEth - t.GasPriceEth*float64(t.GasUsed)
	}
}

func (t Transaction) Time() int64 {
	return t.TimeStamp
}

func (t Transaction) IsValid() bool {
	if t.IsError == 1 || t.ValueEth == 0 || t.To == "" || t.From == "" {
		return false
	}
	return true
}

type Mine struct {
	Acc         string `json:"-"`
	Type        string `json:"-"`
	BlockNumber int    `json:",string"`
	TimeStamp   int64  `json:",string"`
	BlockReward string
	RewardEth   float64 `json:"-"`
}

func (m Mine) Report() string {
	resp := time.Unix(m.TimeStamp, 0).String()[:16]
	resp += " : "
	resp += fmt.Sprint(m.RewardEth)
	resp += " ETH AS "
	resp += m.Type
	resp += " MINING AWARD"
	return resp
}

func (m Mine) Block() int {
	return m.BlockNumber
}

func (m Mine) Effect() float64 {
	return m.RewardEth
}

func (m Mine) Time() int64 {
	return m.TimeStamp
}

func (m Mine) IsValid() bool {
	return true
}

//

type Token interface {
	GetContract() string
	GetName() string
	GetSymbol() string
	GetEffect() float64
	GetReport() string
	GetType() string
	Time() int64
	GetFee() float64
}

type ECR20 struct {
	Acc             string `json:"-"`
	From            string
	To              string
	ContractAddress string
	BlockNumber     int   `json:",string"`
	TimeStamp       int64 `json:",string"`
	TokenName       string
	TokenSymbol     string
	Value           string
	TokenDecimal    int     `json:",string"`
	ValueNum        float64 `json:"-"`
	GasUsed         int     `json:",string"`
	GasPrice        string
	GasPriceEth     float64 `json:"-"`
}

type ECR271 struct {
	Acc             string `json:"-"`
	From            string
	To              string
	ContractAddress string
	BlockNumber     int   `json:",string"`
	TimeStamp       int64 `json:",string"`
	TokenName       string
	TokenSymbol     string
	TokenID         string
	GasUsed         int `json:",string"`
	GasPrice        string
	GasPriceEth     float64 `json:"-"`
}

func (tkn ECR20) GetFee() float64 {
	if tkn.Acc == tkn.From {
		return float64(0) - tkn.GasPriceEth*float64(tkn.GasUsed)
	}
	return float64(0)
}

func (tkn ECR271) GetFee() float64 {
	if tkn.Acc == tkn.From {
		return float64(0) - tkn.GasPriceEth*float64(tkn.GasUsed)
	}
	return float64(0)
}

func (tkn ECR20) GetEffect() float64 {
	if tkn.Acc == tkn.To {
		return tkn.ValueNum
	}
	return 0 - tkn.ValueNum
}

func (tkn ECR271) GetEffect() float64 {
	if tkn.Acc == tkn.To {
		return 1
	}
	return -1
}

func (tkn ECR20) GetReport() string {
	resp := time.Unix(tkn.TimeStamp, 0).String()[:16]
	resp += " : "
	if tkn.Acc == tkn.To {
		resp += "RECEIVED FROM "
		resp += tkn.From
	} else {
		resp += "SENT TO "
		resp += tkn.To
	}
	resp += " = "
	resp += fmt.Sprint(fmt.Sprint(tkn.ValueNum))
	resp += " "
	if tkn.TokenName == "" {
		resp += tkn.ContractAddress
	} else {
		resp += tkn.TokenName
	}
	resp += " ECR20 TOKENS"
	return resp
}

func (tkn ECR271) GetReport() string {
	resp := time.Unix(tkn.TimeStamp, 0).String()[:16]
	resp += " : "
	if tkn.Acc == tkn.To {
		resp += " RECEIVED FROM : "
		resp += tkn.From
	} else {
		resp += " SENT TO : "
		resp += tkn.To
	}
	resp += " : "
	resp += tkn.TokenID
	resp += " ID "
	if tkn.TokenName == "" {
		resp += tkn.ContractAddress
	} else {
		resp += tkn.TokenName
	}
	resp += " ECR721 TOKEN"

	return resp
}

func (tkn ECR20) GetAddress() string {
	return tkn.ContractAddress
}
func (tkn ECR271) GetAddress() string {
	return tkn.ContractAddress
}

func (tkn ECR20) GetName() string {
	return tkn.TokenName
}
func (tkn ECR271) GetName() string {
	return tkn.TokenName
}

func (tkn ECR20) GetSymbol() string {
	return tkn.TokenSymbol
}
func (tkn ECR271) GetSymbol() string {
	return tkn.TokenSymbol
}

func (tkn ECR20) GetContract() string {
	return tkn.ContractAddress
}
func (tkn ECR271) GetContract() string {
	return tkn.ContractAddress
}

func (tkn ECR20) GetType() string {
	return "ECR20"
}
func (tkn ECR271) GetType() string {
	return "ECR721"
}

func (tkn ECR20) Time() int64 {
	return tkn.TimeStamp
}
func (tkn ECR271) Time() int64 {
	return tkn.TimeStamp
}
