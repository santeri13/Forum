package module

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)
func Register_user(username string, email string, password string){
		db, _ := sql.Open("sqlite3", "./Database/database.db")
		stmt, _ := db.Prepare(`
		INSERT INTO "users"(
			"USERNAME",
			"PASSWORD",
			"EMAIL") VALUES (?,?,?);
		`)
		passwords := []byte(password)
		hashedPassword, err := bcrypt.GenerateFromPassword(passwords, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		stmt.Exec(username,hashedPassword,email)
}