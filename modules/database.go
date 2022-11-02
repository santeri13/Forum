package module

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
func CreateDatabase(){
	db, err := sql.Open("sqlite3", "./Database/database.db")
		if err != nil {
			fmt.Println("Error running command:", err)
			return
		}
	users, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "users"(
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"USERNAME" VARCHAR(20) NOT NULL,
			"PASSWORD" VARCHAR(200) NOT NULL,
			"EMAIL" VARCHAR(200) NOT NULL,
			"CREATION_DATE" TEXT,
			"UUID" TEXT
		);
	`)
	users.Exec()
	forum, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "forum"(
			"POST_ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"NAME" TEXT NOT NULL,
			"TEXT" TEXT NOT NULL,
			"AUTHOR_ID" INTEGER NOT NULL,
			"CREATION_DATE"    TEXT,
			FOREIGN KEY("AUTHOR_ID") REFERENCES users("ID")
		);
	`)
	forum.Exec()
	forum_categories, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "forum_categories"(
			"POST_ID" INTEGER NOT NULL,
			"CATEGORY_ID" INTEGER NOT NULL,
			FOREIGN KEY("POST_ID") REFERENCES forum("POST_ID"),
			FOREIGN KEY("CATEGORY_ID") REFERENCES categories("ID")
		);
	`)
	forum_categories.Exec()
	categories, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "categories" (
			"ID"    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"NAME"    TEXT NOT NULL
		);
	`)
	categories.Exec()
	categories_create, _ := db.Prepare(`
		INSERT INTO categories(NAME)
		SELECT("") WHERE NOT EXISTS (SELECT * FROM categories);
	`)
	categories_create.Exec()
	post_reactoion, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "post_reactoion"(
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"POST_ID" INTEGER NOT NULL,
			"POST_COMENTATOR_ID" INTEGER NOT NULL,
			"LIKE" INTEGER,
			"DISLIKE" INTEGER,
			FOREIGN KEY("POST_ID") REFERENCES forum("POST_ID"),
			FOREIGN KEY("POST_COMENTATOR_ID") REFERENCES users("ID")
		);
	`)
	post_reactoion.Exec()
	comments, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "comments"(
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"POST_ID" INTEGER NOT NULL,
			"POST_COMENTATOR_ID" INTEGER NOT NULL,
			"TEXT" TEXT NOT NULL,
			"CREATION_DATE" TEXT,
			FOREIGN KEY("POST_ID") REFERENCES forum("POST_ID"),
			FOREIGN KEY("POST_COMENTATOR_ID") REFERENCES users("ID")
		);
	`)
	comments.Exec()
	comment_reactoion, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "comment_reactoion"(
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"POST_ID" INTEGER NOT NULL,
			"POST_COMENTATOR_ID" INTEGER NOT NULL,
			"LIKE" INTEGER,
			"DISLIKE" INTEGER,
			FOREIGN KEY("POST_ID") REFERENCES comments("ID"),
			FOREIGN KEY("POST_COMENTATOR_ID") REFERENCES users("ID")
		);
	`)
	comment_reactoion.Exec()
	answers_for_forum, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "answers_for_forum"(
		"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"POST_ID" INTEGER NOT NULL,
		"TEXT" TEXT NOT NULL,
		"AUTHOR_ID" INTEGER NOT NULL,
		"CREATION_DATE" TEXT,
		FOREIGN KEY("AUTHOR_ID") REFERENCES users("ID"),
		FOREIGN KEY("POST_ID") REFERENCES forum("POST_ID")
		);
	`)
	answers_for_forum.Exec()
	answers_for_forum_reactoion, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "answers_for_forum_reactoion"(
		"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"POST_ID" INTEGER NOT NULL,
		"POST_COMENTATOR_ID" INTEGER NOT NULL,
		"LIKE" INTEGER,
		"DISLIKE" INTEGER,
		FOREIGN KEY("POST_ID") REFERENCES answers_for_forum("ID"),
		FOREIGN KEY("POST_COMENTATOR_ID") REFERENCES users("ID")
	);
	`)
	answers_for_forum_reactoion.Exec()
	answers_for_forum_comments, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "answers_for_forum_comments"(
		"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"POST_ID" INTEGER NOT NULL,
		"POST_COMENTATOR_ID" INTEGER NOT NULL,
		"TEXT" TEXT NOT NULL,
		"CREATION_DATE" TEXT,
		FOREIGN KEY("POST_ID") REFERENCES answers_for_forum("POST_ID"),
		FOREIGN KEY("POST_COMENTATOR_ID") REFERENCES users("ID")
		);
	`)
	answers_for_forum_comments.Exec()
	answers_for_forum_comments_reactoion, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "answers_for_forum_comments_reactoion"(
		"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"POST_ID" INTEGER NOT NULL,
		"POST_COMENTATOR_ID" INTEGER NOT NULL,
		"LIKE" INTEGER,
		"DISLIKE" INTEGER,
		FOREIGN KEY("POST_ID") REFERENCES answers_for_forum_comments("ID"),
		FOREIGN KEY("POST_COMENTATOR_ID") REFERENCES users("ID")
	);
	`)
	answers_for_forum_comments_reactoion.Exec()
}