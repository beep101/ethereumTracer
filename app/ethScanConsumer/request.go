package ethScanConsumer

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var apiUrl string = "https://api.etherscan.io/api?apikey=E54VNGNEK3GFQAKXVXY8VP3DYJ7HK768C2"

func getEthTransactions(address string, from int, to int, internal bool) ([]byte, error) {
	params := make(map[string]string)

	params["module"] = "account"
	if internal {
		params["action"] = "txlistinternal"
	} else {
		params["action"] = "txlist"
	}
	params["address"] = address
	params["startblock"] = strconv.Itoa(from)
	if to > 0 {
		params["endblock"] = strconv.Itoa(to)
	}
	params["sort"] = "desc"

	return makeRequest(buildUrl(params))
}

func getCrrEthBalance(address string) ([]byte, error) {
	params := make(map[string]string)

	params["module"] = "account"
	params["action"] = "balance"
	params["address"] = address
	params["tag"] = "latest"

	return makeRequest(buildUrl(params))
}

func getLastBlockBefore(timestamp int64) ([]byte, error) {
	return getBlockAtTimestamp(timestamp, true)
}

func getBlockAtTimestamp(timestamp int64, before bool) ([]byte, error) {
	params := make(map[string]string)

	params["module"] = "block"
	params["action"] = "getblocknobytime"
	if before {
		params["closest"] = "before"
	} else {
		params["closest"] = "after"
	}
	params["timestamp"] = strconv.FormatInt(timestamp, 10)

	return makeRequest(buildUrl(params))
}

func getMinedBlocks(address string, uncles bool) ([]byte, error) {
	params := make(map[string]string)

	params["module"] = "account"
	params["action"] = "getminedblocks"
	params["address"] = address
	if uncles {
		params["blocktype"] = "uncles"
	} else {
		params["blocktype"] = "blocks"
	}

	return makeRequest(buildUrl(params))
}

func getERC20Transactions(address string, from int, to int) ([]byte, error) {
	params := make(map[string]string)

	params["module"] = "account"
	params["action"] = "tokentx"
	params["sort"] = "desc"
	params["address"] = address
	params["startblock"] = strconv.Itoa(from)
	if to > 0 {
		params["endblock"] = strconv.Itoa(to)
	}

	return makeRequest(buildUrl(params))
}

func getERC721Transactions(address string, from int, to int) ([]byte, error) {
	params := make(map[string]string)

	params["module"] = "account"
	params["action"] = "tokennfttx"
	params["sort"] = "desc"
	params["address"] = address
	params["startblock"] = strconv.Itoa(from)
	if to > 0 {
		params["endblock"] = strconv.Itoa(to)
	}

	return makeRequest(buildUrl(params))
}

func buildUrl(params map[string]string) string {
	url, _ := url.Parse(apiUrl)
	query := url.Query()

	for key, val := range params {
		query.Set(key, val)
	}

	url.RawQuery = query.Encode()
	return url.String()
}

func makeRequest(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New("API error")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("API response error")
	}
	return body, nil
}
