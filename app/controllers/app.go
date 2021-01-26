package controllers

import (
	ethScan "ethTracer/app/ethScanConsumer"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Transactions(eth bool, from int, address string) revel.Result {
	res := make([]string, 0)
	var err error = nil
	if eth {
		res, err = ethScan.AllTransactions(address, from)
		if err != nil {
			return c.RenderError(err)
		}
	} else {
		res, err = ethScan.AllTokenTransactions(address, from)
		if err != nil {
			return c.RenderError(err)
		}
	}
	return c.RenderJSON(res)
}

func (c App) Balance(eth bool, address, date string) revel.Result {
	if eth {
		res, err := ethScan.BalanceAtDate(address, date)
		if err != nil {
			return c.RenderError(err)
		}
		return c.RenderText(res)
	} else {
		res, err := ethScan.TokenBalanceAtDate(address, date)
		if err != nil {
			return c.RenderError(err)
		}
		return c.RenderJSON(res)
	}

}
