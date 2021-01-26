package ethScanConsumer

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

func BalanceAtDate(address string, date string) (string, error) {
	t, err := time.Parse(time.RFC3339, date+"T00:00:00+00:00")
	if err != nil {
		return "", errors.New("Bad date")
	}
	ut := t.Unix()
	block, err := blockAfterDate(ut)
	if err != nil {
		return "", err
	}

	c := make(chan eventResult, 4)
	go ethTransactions(address, block, -1, false, c)
	go ethTransactions(address, block, -1, true, c)
	go minedBlocks(address, false, c)
	go minedBlocks(address, true, c)

	//token transactions, goroutine sleeps for 1sec, not to exceed max api calls
	tknRes := make(chan tokenResult, 2)
	balanceRes := make(chan balanceResult)
	go func(address string, block int, tknRes chan tokenResult, balanceRes chan balanceResult) {
		time.Sleep(1 * time.Second)
		go ecr20Transactions(address, block, -1, tknRes)
		go ecr271Transactions(address, block, -1, tknRes)
		go balanceOfAccount(address, balanceRes)
	}(address, block, tknRes, balanceRes)
	//...

	prefix := ""
	sum := float64(0)

	for i := 0; i < 4; i++ {
		res := <-c
		if res.err != nil {
			return "", res.err
		}
		if len(res.events) >= 10000 {
			prefix = "Might be incorrect : "
		}
		for _, e := range res.events {
			if e.Block() >= block {
				sum += e.Effect()
			}
		}
	}

	//count in token transaction fees
	for i := 0; i < 2; i++ {
		res := <-tknRes
		if res.err != nil {
			return "", res.err
		}
		if len(res.tokens) >= 10000 {
			prefix = "Might be incorrect : "
		}
		for _, tkn := range res.tokens {
			sum += tkn.GetFee()
		}
	}
	//...
	crrRes := <-balanceRes
	if crrRes.err != nil {
		return "", crrRes.err
	}
	sum = crrRes.balance - sum

	return prefix + "ETH = " + fmt.Sprint(sum), nil
}

func AllTransactions(address string, block int) ([]string, error) {
	c := make(chan eventResult, 4)
	go ethTransactions(address, block, -1, false, c)
	go ethTransactions(address, block, -1, true, c)
	go minedBlocks(address, false, c)
	go minedBlocks(address, true, c)
	events := make([]Event, 0)

	for i := 0; i < 4; i++ {
		res := <-c
		if res.err != nil {
			return nil, res.err
		}
		events = append(events, res.events...)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time() > events[j].Time()
	})

	ret := make([]string, 0)
	for _, e := range events {
		if e.IsValid() {
			ret = append(ret, e.Report())
		}
	}
	if len(ret) == 0 {
		ret = append(ret, "No transactions found")
	}
	return ret, nil
}

func AllTokenTransactions(address string, block int) ([]string, error) {
	c := make(chan tokenResult, 2)
	go ecr20Transactions(address, block, -1, c)
	go ecr271Transactions(address, block, -1, c)
	tokens := make([]Token, 0)

	for i := 0; i < 2; i++ {
		res := <-c
		if res.err != nil {
			return nil, res.err
		}
		tokens = append(tokens, res.tokens...)
	}

	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].Time() > tokens[j].Time()
	})

	ret := make([]string, 0)
	for _, t := range tokens {
		ret = append(ret, t.GetReport())
	}
	if len(ret) == 0 {
		ret = append(ret, "No transactions found")
	}
	return ret, nil
}

func TokenBalanceAtDate(address string, date string) ([]string, error) {
	t, err := time.Parse(time.RFC3339, date+"T00:00:00+00:00")
	if err != nil {
		return nil, errors.New("Bad date")
	}
	ut := t.Unix()
	block, err := blockAtDate(ut)
	if err != nil {
		return nil, err
	}

	c := make(chan tokenResult, 2)
	go ecr20Transactions(address, 0, block, c)
	go ecr271Transactions(address, 0, block, c)
	tokens := make([]Token, 0)

	for i := 0; i < 2; i++ {
		res := <-c
		if res.err != nil {
			return nil, res.err
		}
		tokens = append(tokens, res.tokens...)
	}

	result := make(map[string]float64)
	for _, tkn := range tokens {
		loc := ""
		if tkn.GetName() == "" {
			loc = tkn.GetType() + " " + tkn.GetContract()
		} else {
			loc = tkn.GetType() + " " + tkn.GetName()
		}
		result[loc] = result[loc] + tkn.GetEffect()
	}
	response := make([]string, 0)
	for key, val := range result {
		response = append(response, key+" = "+fmt.Sprint(val))
	}
	if len(response) == 0 {
		response = append(response, "No tokens found")
	}
	return response, nil
}
