package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/abbot/go-http-auth"
	_ "github.com/lib/pq"
)

const isDev bool = false

var db *sql.DB

func Secret(user string, realm string) string {
	if user == "admin" {
		return currEnv.StandardPassword
	}
	return ""
}

type env struct { //environmental variables
	PGUsername       string
	PGPassword       string
	StandardPassword string
	Statichtml       string
}
type portType struct {
	Port string
}
type materialType struct {
	Material string
}
type asOf struct {
	Date string
}
type transaction struct {
	Port     string
	Date     string
	Amount   int
	Material string
	Comment  string
}
type allTransactions struct {
	FirstPort  transaction
	SecondPort transaction
}
type viewTransaction struct {
	Amount   int
	Material string
	Port     string
}

var currEnv env

type success struct {
	Success bool
}
type failure struct {
	Failure error
}

func port(w http.ResponseWriter, r *auth.AuthenticatedRequest){
	switch r.Method {
		case "GET":     
			getPort(w, r)
		case "POST":
			writePortType(w, r)
	}
}
func writePortType(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
	decoder := json.NewDecoder(r.Body)
	var t portType
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
	}
	_, err1 := db.Exec(`INSERT INTO main.possibleports VALUES($1)`, t.Port)
	if err1 != nil {
		log.Println(err1)
	}
	results := new(success)
	results.Success = true
	json.NewEncoder(w).Encode(results)
}
func material(w http.ResponseWriter, r *auth.AuthenticatedRequest){
	switch r.Method {
		case "GET":     
			getMaterial(w, r)
		case "POST":
			writeMaterialType(w, r)
	}
}
func writeMaterialType(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
	decoder := json.NewDecoder(r.Body)
	var t materialType
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
	}
	_, err1 := db.Exec(`INSERT INTO main.possiblematerials VALUES($1)`, t.Material)
	if err1 != nil {
		log.Println(err1)
	}
	results := new(success)
	results.Success = true
	json.NewEncoder(w).Encode(results)
}
func transaction_req(w http.ResponseWriter, r *auth.AuthenticatedRequest){
	switch r.Method {
		case "GET":     
			getAllResults(w, r)
		case "POST":
			writeTransaction(w, r)
	}
}
func writeTransaction(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
	decoder := json.NewDecoder(r.Body)
	var t allTransactions
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		return
	}

	tx, errTx := db.Begin()
	if errTx != nil {
		return
	}

	_, errFirstInput := tx.Exec(`INSERT INTO main.materialtransactions (port, material, transactiondate, amount, comment) VALUES($1, $2, $3, $4, $5)`, t.FirstPort.Port, t.FirstPort.Material, t.FirstPort.Date, t.FirstPort.Amount, t.FirstPort.Comment)

	if t.SecondPort.Port != "" {
		_, errSecondInput := tx.Exec(`INSERT INTO main.materialtransactions (port, material, transactiondate, amount, comment) VALUES($1, $2, $3, $4, $5)`, t.SecondPort.Port, t.SecondPort.Material, t.SecondPort.Date, t.SecondPort.Amount, t.SecondPort.Comment)
		err = errSecondInput
	}

	defer func() {
		if errFirstInput != nil || err != nil {
			tx.Rollback()
			errors := new(failure)
			log.Println(errFirstInput)
			log.Println(err)
			errors.Failure = err
			json.NewEncoder(w).Encode(errors)
			return
		}
		err = tx.Commit()
		if err != nil {
			errors := new(failure)
			log.Println(err)
			errors.Failure = err
			json.NewEncoder(w).Encode(errors)
			return
		}
		results := new(success)
		results.Success = true
		json.NewEncoder(w).Encode(results)

	}()

}
func all(w http.ResponseWriter, r *auth.AuthenticatedRequest){
	switch r.Method {
		case "GET":     
			getAsOfMaterials(w, r)
		case "DELETE":
			deleteAll(w, r)
	}
}
func getAsOfMaterials(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	asOf := r.URL.Query().Get("report_date")
	log.Println("This is as of date:", asOf)
	var results []viewTransaction
	rows, err := db.Query(`SELECT SUM(amount) as amount, material, port FROM main.materialtransactions WHERE transactiondate <= $1 GROUP BY material, port ORDER BY material, port`, "'"+asOf+"'")
	if err != nil {
		errors := new(failure)
		log.Println(err)
		errors.Failure = err
		json.NewEncoder(w).Encode(errors)
	} else {
		defer rows.Close()
		for rows.Next() {
			var portRow viewTransaction
			err := rows.Scan(&portRow.Amount, &portRow.Material, &portRow.Port)
			if err != nil {
				log.Println(err)
			}
			results = append(results, portRow)
		}
		json.NewEncoder(w).Encode(results)
	}

}
func getAllResults(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
	var results []transaction
	rows, err := db.Query(`SELECT  port, CAST(transactiondate as char(10)) as transactiondate, amount, material, CASE WHEN comment IS NULL THEN '' ELSE comment END as comment FROM main.materialtransactions ORDER BY material, port, transactiondate`)
	if err != nil {
		errors := new(failure)
		errors.Failure = err
		log.Println(err)
		json.NewEncoder(w).Encode(errors)
	} else {
		defer rows.Close()
		for rows.Next() {
			var getRow transaction
			err := rows.Scan(&getRow.Port, &getRow.Date, &getRow.Amount, &getRow.Material, &getRow.Comment)
			if err != nil {
				log.Println(err)
			}
			results = append(results, getRow)
		}
		json.NewEncoder(w).Encode(results)
	}

}
func getPort(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
	var results []string
	rows, err := db.Query(`SELECT port FROM main.possibleports ORDER BY port`)
	if err != nil {
		errors := new(failure)
		errors.Failure = err
		log.Println(err)
		json.NewEncoder(w).Encode(errors)
	} else {
		defer rows.Close()
		for rows.Next() {
			var portRow string
			err := rows.Scan(&portRow)
			if err != nil {
				log.Println(err)
			}
			results = append(results, portRow)
		}
		json.NewEncoder(w).Encode(results)
	}
}
func getMaterial(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
	var results []string
	rows, err := db.Query(`SELECT material FROM main.possiblematerials ORDER BY material`)
	if err != nil {
		errors := new(failure)
		errors.Failure = err
		log.Println(err)
		json.NewEncoder(w).Encode(errors)
	} else {
		defer rows.Close()
		for rows.Next() {
			var matRow string
			err := rows.Scan(&matRow)
			if err != nil {
				log.Println(err)
			}
			results = append(results, matRow)
		}
		json.NewEncoder(w).Encode(results)
	}
}
func deleteAll(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
	tx, errTx := db.Begin()
	if errTx != nil {
		return
	}
	_, err := tx.Exec(`TRUNCATE main.possibleports, main.possiblematerials, main.materialtransactions`)

	defer func() {
		if err != nil {
			tx.Rollback()
			errors := new(failure)
			log.Println(err)
			errors.Failure = err
			json.NewEncoder(w).Encode(errors)
			return
		}
		err = tx.Commit()
		if err != nil {
			errors := new(failure)
			log.Println(err)
			errors.Failure = err
			json.NewEncoder(w).Encode(errors)
			return
		}
		results := new(success)
		results.Success = true
		json.NewEncoder(w).Encode(results)

	}()
}

func init() {
	f, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Couldn't open log")
	}
	log.SetOutput(f)
	pwd, _ := os.Getwd()
	file, err1 := ioutil.ReadFile(pwd + "/environment.json") // For read access
	if err1 != nil {
		log.Println(err1)
	}
	err2 := json.Unmarshal(file, &currEnv)
	if err2 != nil {
		log.Println(err2)
	}
	var err3 error
	db, err3 = sql.Open("postgres", "user="+currEnv.PGUsername+" dbname=portserver password="+currEnv.PGPassword)
	if err3 != nil {
		log.Fatal(err3)
	}
	log.Println(currEnv.Statichtml)
}
func main() {
	defer db.Close()

	ath := auth.NewBasicAuthenticator("example.com", Secret)
	http.HandleFunc("/", ath.Wrap(func(res http.ResponseWriter, req *auth.AuthenticatedRequest) {
		http.FileServer(http.Dir(currEnv.Statichtml)).ServeHTTP(res, &req.Request)
	}))
	
	
	http.HandleFunc("/port", ath.Wrap(port))
	http.HandleFunc("/material", ath.Wrap(material))
	http.HandleFunc("/transaction", ath.Wrap(transaction_req))
	http.HandleFunc("/all", ath.Wrap(all))
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
