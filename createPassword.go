package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	//args:=os.Args[1:]
	pswd := os.Args[1]
	bytePassword := []byte(pswd)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err)

	}
	fmt.Println(string(hashedPassword))
}
