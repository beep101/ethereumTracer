package main

import (
	eth "ethTracer/app/ethScanConsumer"
	"fmt"
)

func main() {
	str, _ := eth.GetBalanceAtDate("0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a", "2015-08-09")
	fmt.Println(str)
}
