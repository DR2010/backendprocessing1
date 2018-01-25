// Main web application program for restauranteweb
// -----------------------------------------------
// .../src/restaurantetick/restaurantetick.go
// -----------------------------------------------
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	btcmarketshandler "restaurantetick/areas/btcmarketshandler"
	helper "restauranteweb/areas/helper"
	"time"

	"github.com/go-redis/redis"
	// _ "github.com/go-sql-driver/mysql"
)

var mongodbvar helper.DatabaseX
var credentials helper.Credentials

var db *sql.DB
var err error
var redisclient *redis.Client

// Looks after the main routing
//
func main() {

	redisclient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	loadreferencedatainredis()

	// Read variables from server
	//
	envirvar := new(helper.RestEnvVariables)
	envirvar.APIAPIServerIPAddress, _ = redisclient.Get("Web.APIServer.IPAddress").Result()
	envirvar.APIAPIServerPort, _ = redisclient.Get("Web.APIServer.Port").Result()
	envirvar.WEBServerPort, _ = redisclient.Get("WEBServerPort").Result()
	envirvar.WEBDebug, _ = redisclient.Get("Web.Debug").Result()
	envirvar.RecordCurrencyTick, _ = redisclient.Get("RecordCurrencyTick").Result()
	envirvar.RunningFromServer, _ = redisclient.Get("RunningFromServer").Result()

	// btcmarketshandler.SendEmail(redisclient, "StartingSystemNow"+envirvar.RunningFromServer)

	fmt.Println(">>> Web Server: restauranteweb.exe running.")
	fmt.Println("Loading reference data in cache - Redis")

	mongodbvar.Location = "localhost"
	mongodbvar.Database = "restaurante"

	credentials.UserID = "No User"

	for {

		fmt.Println(">>> Web Server: restauranteweb.exe running. ")
		t := time.Now()
		fmt.Println(t.Format("2006-01-02 15:04:05"))

		strjson, btccoin := emailtick()

		fmt.Println(strjson)

		fmt.Println("Recording...")
		btcmarketshandler.RecordTick(redisclient, btccoin, "GoLangTick")

		// for x := 0; x < len(btccoin); x++ {
		// 	fmt.Println(btccoin[x].Currency + " " + btccoin[x].CotacaoAtual + " " + btccoin[x].DateTime)
		// }
		fmt.Println("Sleeping 120,000...")
		time.Sleep(100000 * time.Millisecond)

	}

}

func loadreferencedatainredis() {

	variable := helper.Readfileintostruct()
	err = redisclient.Set("Web.MongoDB.Database", variable.APIMongoDBDatabase, 0).Err()
	err = redisclient.Set("Web.APIServer.Port", variable.APIAPIServerPort, 0).Err()
	err = redisclient.Set("WEBServerPort", variable.WEBServerPort, 0).Err()
	err = redisclient.Set("Web.MongoDB.Location", variable.APIMongoDBLocation, 0).Err()
	err = redisclient.Set("Web.APIServer.IPAddress", variable.APIAPIServerIPAddress, 0).Err()
	err = redisclient.Set("Web.Debug", variable.WEBDebug, 0).Err()
	err = redisclient.Set("RecordCurrencyTick", variable.RecordCurrencyTick, 0).Err()
	err = redisclient.Set("RunningFromServer", variable.RunningFromServer, 0).Err()

}

func root(httpwriter http.ResponseWriter, r *http.Request) {

	// create new template
	var listtemplate = `
		{{define "listtemplate"}}
		This is my web site, Daniel - aka D#.
		<p/>
		<p/>
		<picture>
			<img src="images/avatar.png" alt="Avatar" width="400" height="400">
		</picture>
		{{end}}
		`

	t, _ := template.ParseFiles("templates/indextemplate.html")
	t, _ = t.Parse(listtemplate)

	t.Execute(httpwriter, listtemplate)
	return
}

func root2(httpwriter http.ResponseWriter, r *http.Request) {
	http.ServeFile(httpwriter, r, "index.html")

	return
}

func btcrecordtick(httpwriter http.ResponseWriter, req *http.Request) {

	// Fazer o record tick aceitar um parametro para gravar de onde a rotina foi chamada
	// .... btcrecordtick?rotina=CURLubuntuAUTO
	// .... btcrecordtick?rotina=WindowsPC
	// .... btcrecordtick?rotina=WindowsPCCURL

	// params := req.URL.Query()
	// var rotina = params.Get("rotina")
	// if rotina == "" {
	// 	rotina = "Not sure - web test most likely"
	// }

	rotina := "restaurantetick.service"

	var listofbit = btcmarketshandler.GetBalance(redisclient)
	btcmarketshandler.RecordTick(redisclient, listofbit, rotina)

	jsonval, _ := json.Marshal(listofbit)
	jsonstring := string(jsonval)

	http.Error(httpwriter, jsonstring, 200)
}

func emailtick() (string, []btcmarketshandler.BalanceCrypto) {

	var listofbit = btcmarketshandler.GetBalance(redisclient)

	jsonval, _ := json.Marshal(listofbit)
	jsonstring := string(jsonval)

	return jsonstring, listofbit
}
