package module

import (
	"fmt"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

func Login(username string, password string)bool{
	if(username=="" ||password=="" ){
		return false
	}
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	stmt, err := db.Prepare(`SELECT PASSWORD FROM "users" WHERE USERNAME == ?`)
	if err != nil {
		fmt.Println("Error running command:", err)
		return false
	}
	var passwords string
	err2 := stmt.QueryRow(username).Scan(&passwords)
	if err2 != nil {
		return false
	}
	switch bcrypt.CompareHashAndPassword([]byte(passwords), []byte(password))==nil{
	case false:
		return false
	}
	return true
}