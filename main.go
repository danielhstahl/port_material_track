package main
import (
	"database/sql"
	_ "github.com/lib/pq"
    //"fmt"
    "log"
    "net/http"
    "os"
    "io/ioutil"
    "encoding/json"
)
var db *sql.DB
type env struct{ //environmental variables
    PGUsername string
    PGPassword string
    Statichtml string
}
type portType struct{
    Port string
}
type materialType struct{
    Material string
}
type asOf struct{
    Date string
}
type transaction struct{
    Port string 
    Date string
    Amount int
    Material string
}
type getResults struct{
    Amount int
    Material string
    Port string
}
var currEnv env
type success struct{
    Success bool
}
type failure struct{
    Failure error
}

func writePortType(w http.ResponseWriter, r *http.Request){
    /*w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");*/
    decoder := json.NewDecoder(r.Body)
    var t portType   
    err := decoder.Decode(&t)
    if err != nil {
        log.Println(err) 
    }
    _, err1:=db.Exec(`INSERT INTO main.possibleports VALUES($1)`, t.Port)
    if err1!=nil{
        log.Println(err1)
    }
    results:=new(success)
    results.Success=true
    json.NewEncoder(w).Encode(results)
}
func writeMaterialType(w http.ResponseWriter, r *http.Request){
    /*w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");*/
    decoder := json.NewDecoder(r.Body)
    var t materialType   
    err := decoder.Decode(&t)
    if err != nil {
        log.Println(err)
    }
    _, err1:=db.Exec(`INSERT INTO main.possiblematerials VALUES($1)`, t.Material)
    if err1!=nil{
        log.Println(err1)
    }
    results:=new(success)
    results.Success=true
    json.NewEncoder(w).Encode(results)
}
func writeTransaction(w http.ResponseWriter, r *http.Request){
    /*w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");*/
    decoder := json.NewDecoder(r.Body)
    var t transaction   
    err := decoder.Decode(&t)
    if err != nil {
        log.Println(err)
    }
    _, err1:=db.Exec(`INSERT INTO main.materialtransactions VALUES($1, $2, $3, $4)`, t.Port, t.Material, t.Date, t.Amount)
    
    
    if err1!=nil{
        errors:=new(failure)
        log.Println(err1)
        errors.Failure=err
        json.NewEncoder(w).Encode(errors)
    }else{
        results:=new(success)
        results.Success=true
        json.NewEncoder(w).Encode(results)
    }
    
}
func getAsOfMaterials(w http.ResponseWriter, r *http.Request){
    /*w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");*/
    decoder := json.NewDecoder(r.Body)
    var t asOf
    err := decoder.Decode(&t)
    log.Println("got here 110")
    if err != nil {
        log.Println(err)
    }
    var results []getResults
	rows, err1 := db.Query(`SELECT SUM(amount) as amount, material, port FROM main.materialtransactions WHERE transactiondate <= $1 GROUP BY material, port ORDER BY material, port`, "'"+t.Date+"'")
    if err1!=nil{
        errors:=new(failure)
        log.Println(err1)
        json.NewEncoder(w).Encode(errors)
    }else{
        defer rows.Close()
        for rows.Next(){
            var portRow getResults 
            err:=rows.Scan(&portRow.Amount, &portRow.Material, &portRow.Port)
            if err!=nil{
                log.Println(err)
            }
            results=append(results, portRow)
        }
        json.NewEncoder(w).Encode(results)
    }
    
}
func getPort(w http.ResponseWriter, r *http.Request){
    /*w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");*/
    var results []string
	rows, err1 := db.Query(`SELECT port FROM main.possibleports ORDER BY port`)
    if err1!=nil{
        log.Fatal(err1)
    }
    defer rows.Close()
    for rows.Next(){
        var portRow string 
        err:=rows.Scan(&portRow)
        if err!=nil{
            log.Fatal(err)
        }
        results=append(results, portRow)
    }
    json.NewEncoder(w).Encode(results)
}
func getMaterial(w http.ResponseWriter, r *http.Request){
    /*w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");*/
    var results []string
	rows, err1 := db.Query(`SELECT material FROM main.possiblematerials ORDER BY material`)
    if err1!=nil{
        log.Fatal(err1)
    }
    defer rows.Close()
    for rows.Next(){
        var matRow string 
        err:=rows.Scan(&matRow)
        if err!=nil{
            log.Fatal(err)
        }
        results=append(results, matRow)
    }
    json.NewEncoder(w).Encode(results)
}

func init(){
    f, err := os.OpenFile("log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Couldn't open log")
    }
    log.SetOutput(f)
    pwd, _ := os.Getwd()
    file, err1 := ioutil.ReadFile(pwd+"/environment.json") // For read access
    if err1!=nil{
        log.Fatal(err1)
    }
    err2:=json.Unmarshal(file,&currEnv)
    if err2!=nil{
        log.Fatal(err2)
    }
    var err3 error
    db, err3 = sql.Open("postgres", "user="+currEnv.PGUsername+" dbname=portserver password="+currEnv.PGPassword)
	if err3 != nil {
		log.Fatal(err3)
	}
}
func main(){
    defer db.Close()
    fs := http.FileServer(http.Dir(currEnv.Statichtml))
    http.Handle("/", fs)
    http.HandleFunc("/writePort", writePortType)
    http.HandleFunc("/writeMaterial", writeMaterialType)
    http.HandleFunc("/writeTransaction", writeTransaction)
    http.HandleFunc("/getResults", getAsOfMaterials)
    http.HandleFunc("/getPorts", getPort)
    http.HandleFunc("/getMaterials", getMaterial)
    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)
}
