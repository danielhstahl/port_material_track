package main
import (
	"database/sql"
	_ "github.com/lib/pq"
    "github.com/abbot/go-http-auth"
    "log"
    "net/http"
    "os"
    "io/ioutil"
    "encoding/json"
)
const isDev bool = false
var db *sql.DB
func Secret(user string, realm string) string {
    if user == "admin" {
        return currEnv.StandardPassword
    }
    return ""
}
type env struct{ //environmental variables
    PGUsername string
    PGPassword string
    StandardPassword string
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
    Comment string
}
type viewTransaction struct{
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

func writePortType(w http.ResponseWriter, r *auth.AuthenticatedRequest){
    if(isDev){
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    }
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
func writeMaterialType(w http.ResponseWriter, r *auth.AuthenticatedRequest){
    if(isDev){
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    }
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
func writeTransaction(w http.ResponseWriter, r *auth.AuthenticatedRequest){
    if(isDev){
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    }
    decoder := json.NewDecoder(r.Body)
    var t transaction   
    err := decoder.Decode(&t)
    if err != nil {
        log.Println(err)
        return;
    }
    _, err1:=db.Exec(`INSERT INTO main.materialtransactions (port, material, transactiondate, amount, comment) VALUES($1, $2, $3, $4, $5)`, t.Port, t.Material, t.Date, t.Amount, t.Comment)
    
    
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
func getAsOfMaterials(w http.ResponseWriter, r *auth.AuthenticatedRequest){
    decoder := json.NewDecoder(r.Body)
    var t asOf
    err := decoder.Decode(&t)
    if err != nil {
        log.Println(err)
    }
    var results []viewTransaction
	rows, err1 := db.Query(`SELECT SUM(amount) as amount, material, port FROM main.materialtransactions WHERE transactiondate <= $1 GROUP BY material, port ORDER BY material, port`, "'"+t.Date+"'")
    if err1!=nil{
        errors:=new(failure)
        log.Println(err1)
        errors.Failure=err1
        json.NewEncoder(w).Encode(errors)
    }else{
        defer rows.Close()
        for rows.Next(){
            var portRow viewTransaction 
            err:=rows.Scan(&portRow.Amount, &portRow.Material, &portRow.Port)
            if err!=nil{
                log.Println(err)
            }
            results=append(results, portRow)
        }
        json.NewEncoder(w).Encode(results)
    }
    
}
func getAllResults(w http.ResponseWriter, r *auth.AuthenticatedRequest){
    if(isDev){
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    }
    var results []transaction
	rows, err1 := db.Query(`SELECT  port, CAST(transactiondate as char(10)) as transactiondate, amount, material, comment FROM main.materialtransactions ORDER BY material, port, transactiondate`)
    if err1!=nil{
        errors:=new(failure)
        errors.Failure=err1
        log.Println(err1)
        json.NewEncoder(w).Encode(errors)
    }else{
        defer rows.Close()
        for rows.Next(){
            var getRow transaction 
            err:=rows.Scan(&getRow.Port, &getRow.Date, &getRow.Amount, &getRow.Material, &getRow.Comment)
            if err!=nil{
                log.Println(err)
            }
            results=append(results, getRow)
        }
        json.NewEncoder(w).Encode(results)
    }
    
}
func getPort(w http.ResponseWriter, r *auth.AuthenticatedRequest){
    if(isDev){
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    }
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
func getMaterial(w http.ResponseWriter, r *auth.AuthenticatedRequest){
    if(isDev){
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    }
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
    
    ath := auth.NewBasicAuthenticator("example.com", Secret)
    http.HandleFunc("/", ath.Wrap(func(res http.ResponseWriter, req *auth.AuthenticatedRequest) {
        http.FileServer(http.Dir(currEnv.Statichtml)).ServeHTTP(res, &req.Request)
    }))
    http.HandleFunc("/writePort", ath.Wrap(writePortType))
    http.HandleFunc("/writeMaterial", ath.Wrap(writeMaterialType))
    http.HandleFunc("/writeTransaction", ath.Wrap(writeTransaction))
    http.HandleFunc("/getResults", ath.Wrap(getAsOfMaterials))
    http.HandleFunc("/getPorts", ath.Wrap(getPort))
    http.HandleFunc("/getMaterials", ath.Wrap(getMaterial))
    http.HandleFunc("/getAllResults", ath.Wrap(getAllResults))
    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)
}
