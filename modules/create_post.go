package module

import (
	"database/sql"
	"fmt"
	"time"
)
type session struct {
	username string
	expiry   time.Time
}
var sessions = map[string]session{}

func Create_post(username string, postname string, post_text string, category []string)bool{
	var ID string
	db, _ := sql.Open("sqlite3", "./Database/database.db")
		post, err3 := db.Prepare(`
			INSERT INTO "forum"(
				"NAME",
				"TEXT",
				"AUTHOR_ID") VALUES (?,?,?);
			`)
			if err3 != nil {
				return false
			}
		stmt, err := db.Prepare(`
		SELECT ID FROM "users" WHERE USERNAME == ?`)
		if err != nil {
			return false
		}
		err2 := stmt.QueryRow(username).Scan(&ID)
		if err2 != nil {
			return false
		}
		post.Exec(postname,post_text,ID)
		// Create catgory
		var name string
		var post_id string
		post_ID, err := db.Prepare(`
			SELECT POST_ID FROM "forum" WHERE NAME == ?`)
		if err != nil {
			fmt.Println("Error running command: POST_ID FROM forum",)
			return false
		}
		err2 = post_ID.QueryRow(postname).Scan(&post_id)
		if err2 != nil {
			fmt.Println("Error running command: post_ID")
			return false
		}
		for i := 0; i < len(category); i++{
			categorys, err := db.Prepare(`
				SELECT ID FROM "categories" WHERE NAME == ?`)
			if err != nil {
				fmt.Println("Error running command: Problem with category")
				return false
			}
			err2 = categorys.QueryRow(category[i]).Scan(&name)
			if err2 != nil {
				fmt.Println("No Category")
				new_category, err := db.Prepare(`
					INSERT INTO "categories"(
					"NAME") VALUES (?);
				`)
				if err != nil {
					fmt.Println("Error running command: Not insert into category")
					return false
				}
				new_category.Exec(category[i])
			}
			new_category, _ := db.Prepare(`
				INSERT INTO "forum_categories"(
				"POST_ID",
				"CATEGORY_ID") VALUES (?,?);
			`)
			var category_id string
			category_ID, err := db.Prepare(`
			SELECT ID FROM "categories" WHERE NAME == ?`)
			if err != nil {
				fmt.Println("Error running command: ID")
				return false
			}
			err2 = category_ID.QueryRow(category[i]).Scan(&category_id)
			if err2 != nil {
				fmt.Println("Error running command:category_ID")
				return false
			}
			new_category.Exec(post_id,category_id)
		}
		return true
}