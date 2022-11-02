package module

import (
	"database/sql"
	"log"
)

func Dislike(username string,p string,post_type string){
	var ID string
    var Like string
    var Dislike string
	var usernameID string
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	userID_query, _ := db.Prepare(`SELECT ID FROM "users" WHERE USERNAME == ?`)
	err2 := userID_query.QueryRow(username).Scan(&usernameID)
	if err2 != nil {
		log.Fatal(err2)
	}
	check_on_dislike, _ := db.Query(`SELECT ID, LIKE, DISLIKE FROM "`+post_type+`" WHERE POST_COMENTATOR_ID == `+usernameID+` AND POST_ID==`+p)

    for check_on_dislike.Next() {
        if err := check_on_dislike.Scan(&ID,&Like,&Dislike); err != nil && err != sql.ErrNoRows{
            dislike, _ := db.Prepare(`
			INSERT INTO "`+post_type+`"(
				"POST_ID",
				"POST_COMENTATOR_ID",
				"LIKE",
				"DISLIKE") VALUES (?,?,?,?);
			`)
			dislike.Exec(p,usernameID,"",1)
			return
        }
    }
	check_on_dislike.Close()
	switch(Dislike){
	case "":
		if Like == ""{
			dislike, _ := db.Prepare(`
			INSERT INTO "`+post_type+`"(
				"POST_ID",
				"POST_COMENTATOR_ID",
				"LIKE",
				"DISLIKE") VALUES (?,?,?,?);
			`)
			dislike.Exec(p,usernameID,"",1)
			Dislike = ""
			Like = ""
			return
		}
		dislike, _ := db.Prepare(`
			UPDATE "`+post_type+`" SET "LIKE" = ?,
			"DISLIKE" = ? WHERE "POST_ID" == ? AND "POST_COMENTATOR_ID" == ?
		`)
		dislike.Exec("",1,p,usernameID)
		Dislike = ""
		Like = ""
		return
	case "1":
		delete_dislike, _ := db.Prepare(`
		DELETE from "`+post_type+`"
		WHERE "POST_COMENTATOR_ID" == ? AND "POST_ID" == ?
		`)
		delete_dislike.Exec(usernameID,p)
		Dislike = ""
		Like = ""
		return
	}
}