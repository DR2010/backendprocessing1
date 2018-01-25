// Package btcmarketshandler Handler for dishes web
// -----------------------------------------------------------
// .../src/restauranteweb/areas/disherhandler/ordershandler.go
// -----------------------------------------------------------
package btcmarketshandler

import (
	"fmt"
	"html/template"
	"net/http"
	helper "restauranteweb/areas/helper"

	"github.com/go-redis/redis"
)

// This is the template to display as part of the html template
//

// ControllerInfo is
type ControllerInfo struct {
	Name     string
	Message  string
	Currency string
	FromDate string
	ToDate   string
}

// Row is
type Row struct {
	Description []string
}

// Coin is
type Coin struct {
	Short string
	Name  string
}

// DisplayTemplate is
type DisplayTemplate struct {
	Info       ControllerInfo
	FieldNames []string
	Rows       []Row
	Btccoin    []BalanceCrypto
	Coins      []Coin
}

var mongodbvar helper.DatabaseX

// List = assemble results of API call to dish list
//
func List(httpwriter http.ResponseWriter, redisclient *redis.Client) {

	// create new template
	t, _ := template.ParseFiles("templates/basictemplate.html", "templates/btcmarkets/btcmarketstemplate.html")

	// Get list of orders (api call)
	//
	var list = ListCoinsIhave(redisclient)

	// Assembly display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Daniel Investment List"
	items.Info.Currency = "NA"

	var numberoffields = 4

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Currency"
	items.FieldNames[1] = "Balance"
	items.FieldNames[2] = "CotacaoAtual"
	items.FieldNames[3] = "ValueInCashAUD"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Btccoin = make([]BalanceCrypto, len(list))

	for i := 0; i < len(list); i++ {

		items.Btccoin[i] = BalanceCrypto{}
		items.Btccoin[i].Balance = list[i].Balance
		items.Btccoin[i].Currency = list[i].Currency
		items.Btccoin[i].CotacaoAtual = list[i].CotacaoAtual
		items.Btccoin[i].ValueInCashAUD = list[i].ValueInCashAUD

	}

	t.Execute(httpwriter, items)

}

// ListV2 = assemble results of API call to dish list
//
func ListV2(httpwriter http.ResponseWriter, redisclient *redis.Client) []BalanceCrypto {

	// create new template
	t, _ := template.ParseFiles("templates/basictemplate.html", "templates/btcmarkets/btcmarketstemplate.html")

	// Get list of orders (api call)
	//
	var list = ListCoinsIhave(redisclient)

	// Assembly display structure for html template
	//
	items := DisplayTemplate{}
	items.Info.Name = "Coins"
	items.Info.Currency = "SUMMARY"

	var numberofcoins = 8
	items.Coins = make([]Coin, numberofcoins)
	items.Coins[0].Short = "AUD"
	items.Coins[0].Name = "Australian Dollar"
	items.Coins[1].Short = "BTC"
	items.Coins[1].Name = "Bitcoin"
	items.Coins[2].Short = "BTC"
	items.Coins[2].Name = "Bitcoin"
	items.Coins[3].Short = "LTC"
	items.Coins[3].Name = "Litecoin"
	items.Coins[4].Short = "ETH"
	items.Coins[4].Name = "Ethereum"
	items.Coins[5].Short = "XRP"
	items.Coins[5].Name = "Ripple"
	items.Coins[6].Short = "BCH"
	items.Coins[6].Name = "Bitcash"
	items.Coins[7].Short = "ALL"
	items.Coins[7].Name = "All Coins"

	var numberoffields = 4

	// Set colum names
	items.FieldNames = make([]string, numberoffields)
	items.FieldNames[0] = "Currency"
	items.FieldNames[1] = "Balance"
	items.FieldNames[2] = "CotacaoAtual"
	items.FieldNames[3] = "ValueInCashAUD"

	// Set rows to be displayed
	items.Rows = make([]Row, len(list))
	items.Btccoin = make([]BalanceCrypto, len(list))

	var RetBtccoin []BalanceCrypto
	RetBtccoin = make([]BalanceCrypto, len(list))

	for i := 0; i < len(list); i++ {

		items.Btccoin[i] = BalanceCrypto{}
		items.Btccoin[i].Balance = list[i].Balance
		items.Btccoin[i].Currency = list[i].Currency
		items.Btccoin[i].CotacaoAtual = list[i].CotacaoAtual
		items.Btccoin[i].ValueInCashAUD = list[i].ValueInCashAUD

		// New code to return values to write to mongo every minute or every call
		// 31/12/2017
		//
		RetBtccoin[i] = BalanceCrypto{}
		RetBtccoin[i].Balance = list[i].Balance
		RetBtccoin[i].Currency = list[i].Currency
		RetBtccoin[i].CotacaoAtual = list[i].CotacaoAtual
		RetBtccoin[i].ValueInCashAUD = list[i].ValueInCashAUD

	}

	t.Execute(httpwriter, items)

	return RetBtccoin
}

// GetBalance = assemble results of API call to dish list
//
func GetBalance(redisclient *redis.Client) []BalanceCrypto {

	// Get list of orders (api call)
	//
	var list = ListCoinsIhave(redisclient)

	var RetBtccoin []BalanceCrypto
	RetBtccoin = make([]BalanceCrypto, len(list))

	for i := 0; i < len(list); i++ {

		// New code to return values to write to mongo every minute or every call
		// 31/12/2017
		//
		RetBtccoin[i] = BalanceCrypto{}
		RetBtccoin[i].Balance = list[i].Balance
		RetBtccoin[i].Currency = list[i].Currency
		RetBtccoin[i].CotacaoAtual = list[i].CotacaoAtual
		RetBtccoin[i].ValueInCashAUD = list[i].ValueInCashAUD
		RetBtccoin[i].Volume24 = list[i].Volume24
		RetBtccoin[i].BestBid = list[i].BestBid
		RetBtccoin[i].BestAsk = list[i].BestAsk

	}

	return RetBtccoin
}

// RecordTick is xxx
func RecordTick(redisclient *redis.Client, balcrypto []BalanceCrypto, rotina string) {

	for i := 0; i < len(balcrypto); i++ {

		bcc := BalanceCrypto{}
		bcc.Balance = balcrypto[i].Balance
		bcc.Currency = balcrypto[i].Currency
		bcc.CotacaoAtual = balcrypto[i].CotacaoAtual
		bcc.ValueInCashAUD = balcrypto[i].ValueInCashAUD
		bcc.BestAsk = balcrypto[i].BestAsk
		bcc.BestBid = balcrypto[i].BestBid
		bcc.Volume24 = balcrypto[i].Volume24
		balcrypto[i].Rotina = rotina
		bcc.Rotina = rotina
		fmt.Println("RecordTick: " + bcc.Currency + " <--")

		APIcallAdd(redisclient, bcc, rotina)
	}

	return
}

// Lpad is left pad
func Lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}
