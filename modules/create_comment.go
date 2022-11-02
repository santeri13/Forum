package module

import (
	"database/sql"
	"log"
)

func Create_comment(username string,post_text string,ID string,post_type string){
	var usernameID string
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	userID_query, _ := db.Prepare(`SELECT ID FROM "users" WHERE USERNAME == ?`)
	err2 := userID_query.QueryRow(username).Scan(&usernameID)
	if err2 != nil {
		log.Fatal(err2)
	}
	post, err := db.Prepare(`
	INSERT INTO "`+post_type+`"(
		"POST_ID",
		"POST_COMENTATOR_ID",
		"TEXT") VALUES (?,?,?);
	`)
	if err != nil {
		log.Fatal(err)
	}
	post.Exec(ID,usernameID,post_text)
}