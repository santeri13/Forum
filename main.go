package main

import (
	"database/sql"
	"fmt"
	modules "forum/modules"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type session struct {
	username string
	expiry   time.Time
}

var sessions = map[string]session{}

type comments struct {
	ID   string
	Text string
}
type Forum struct {
	ID       string
	Name     string
	Text     string
	Like     string
	Dislike  string
	Comments []comments
}
type Options struct {
	ID   string
	Name string
}

func main() {
	modules.CreateDatabase()
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	http.HandleFunc("/", ForumPage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/register", RegisterPage)
	http.HandleFunc("/original_post/", Original_post)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/create_comment", Create_comment)
	http.HandleFunc("/like", Like)
	http.HandleFunc("/dislike", Dislike)
	http.HandleFunc("/post_creation", Forum_post)
	http.HandleFunc("/answer_for_forum", Answer_for_forum)
	http.HandleFunc("/option", Option)
	http.ListenAndServe(":8080", nil)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("HTTP status 500 - Internal server error: %s", err)
	}
}
func ForumPage(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			data := map[string]interface{}{"postdata": Forums_data("", "", "", "", ""), "listdata": Options_data()}
			tmpl := template.Must(template.ParseFiles("templates/forum.html"))
			tmpl.Execute(w, data)
			return
		}
		data := map[string]interface{}{"postdata": Forums_data("", "", "", "", ""), "listdata": Options_data()}
		tmpl := template.Must(template.ParseFiles("templates/forum.html"))
		tmpl.Execute(w, data)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		data := map[string]interface{}{"postdata": Forums_data("", "", "", "", ""), "listdata": Options_data()}
		tmpl := template.Must(template.ParseFiles("templates/forum.html"))
		tmpl.Execute(w, data)
		return
	}
	if userSession.isExpired(){
		delete(sessions, sessionToken)
		data := map[string]interface{}{"postdata": Forums_data("", "", "", "", ""), "listdata": Options_data()}
		tmpl := template.Must(template.ParseFiles("templates/forum.html"))
		tmpl.Execute(w, data)
		return
	}
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	query_id := db.QueryRow(`SELECT UUID FROM "users" WHERE username == ?`, userSession.username)
	var uuid string
	query_id.Scan(&uuid)
	if sessionToken != uuid{
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}
	data := map[string]interface{}{"postdata": Forums_data("", "", "", "", ""), "username": userSession.username, "listdata": Options_data()}
	tmpl := template.Must(template.ParseFiles("templates/forum.html"))
	tmpl.Execute(w, data)
}
func LoginPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
	case "POST":
		switch modules.Login(r.FormValue("username"), r.FormValue("psw")) {
		case true:
			create_session(r.FormValue("username"), w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		case false:
			data := map[string]interface{}{"error": "Wrong username or password"}
			tmpl := template.Must(template.ParseFiles("templates/login.html"))
			tmpl.Execute(w, data)
			return
		}
	}
}
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
	case "POST":
		switch r.FormValue("psw") == r.FormValue("psw-repeat") && r.FormValue("email") != "" {
		case true:
			db, _ := sql.Open("sqlite3", "./Database/database.db")
			usernmae := db.QueryRow(`SELECT ID FROM "users" WHERE USERNAME == ?`, r.FormValue("username"))
			email := db.QueryRow(`SELECT ID FROM "users" WHERE  EMAIL == ?`, r.FormValue("email"))
			name := ""
			mail := ""
			usernmae.Scan(&name)
			email.Scan(&mail)
			switch true {
			case name == "" && mail == "":
				modules.Register_user(r.FormValue("username"), r.FormValue("email"), r.FormValue("psw"))
				create_session(r.FormValue("username"), w)
				http.Redirect(w, r, "/", http.StatusFound)
			case name != "":
				data := map[string]interface{}{"error": "Name already used"}
				tmpl := template.Must(template.ParseFiles("templates/register.html"))
				tmpl.Execute(w, data)
			case mail != "":
				data := map[string]interface{}{"error": "Email already used"}
				tmpl := template.Must(template.ParseFiles("templates/register.html"))
				tmpl.Execute(w, data)
			}
		case false:
			data := map[string]interface{}{"error": "Passwords are not same"}
			tmpl := template.Must(template.ParseFiles("templates/register.html"))
			tmpl.Execute(w, data)
		}
	}
}
func Forum_post(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		sessionToken := c.Value
		userSession, exists := sessions[sessionToken]
		if !exists {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if userSession.isExpired() {
			delete(sessions, sessionToken)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		tmpl := template.Must(template.ParseFiles("templates/post.html"))
		tmpl.Execute(w, nil)
	case "POST":
		c, _ := r.Cookie("session_token")
		sessionToken := c.Value
		userSession, _ := sessions[sessionToken]
		r.ParseForm()
		switch modules.Create_post(userSession.username, r.FormValue("post_name"), r.FormValue("post_text"), strings.Split(r.FormValue("category"), ",")) {
		case true:
			http.Redirect(w, r, "/", http.StatusFound)
			return
		case false:
			http.Redirect(w, r, "/post_creation", http.StatusFound)
			return
		}
	}
}
func Original_post(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := r.URL.Path[15:]
		db, _ := sql.Open("sqlite3", "./Database/database.db")
		name_query, _ := db.Prepare(`SELECT NAME FROM "forum" WHERE POST_ID == ? `)
		text_query, _ := db.Prepare(`SELECT TEXT FROM "forum" WHERE POST_ID == ?`)
		like_query,_ := db.Prepare(`SELECT count(NULLIF(LIKE, '')) FROM post_reactoion WHERE POST_ID == ?`)
		dislike_query,_ := db.Prepare(`SELECT count(NULLIF(Dislike, '')) FROM post_reactoion WHERE POST_ID == ?`)
		var name string
		var text string
		var like string
		var dislike string
		err2 := name_query.QueryRow(id).Scan(&name)
		if err2 != nil {
			panic(err2)
		}
		err2 = text_query.QueryRow(id).Scan(&text)
		if err2 != nil {
			panic(err2)
		}
		err2 = like_query.QueryRow(id).Scan(&like)
		if err2 != nil {
			panic(err2)
		}
		err2 = dislike_query.QueryRow(id).Scan(&dislike)
		if err2 != nil {
			panic(err2)
		}
		data := map[string]interface{}{"postid": id, "postname": name, "posttext": text, "answerdata": Answers_data(id), "comentsdata": Comments_data(id),"Like":like, "Dislike":dislike}
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				tmpl := template.Must(template.ParseFiles("templates/post_page.html"))
				tmpl.Execute(w, data)
				return
			}
			tmpl := template.Must(template.ParseFiles("templates/post_page.html"))
			tmpl.Execute(w, data)
			return
		}
		sessionToken := c.Value
		userSession, exists := sessions[sessionToken]
		if !exists {
			tmpl := template.Must(template.ParseFiles("templates/post_page.html"))
			tmpl.Execute(w, data)
			return
		}
		if userSession.isExpired() {
			delete(sessions, sessionToken)
			tmpl := template.Must(template.ParseFiles("templates/post_page.html"))
			tmpl.Execute(w, data)
			return
		}
		query_id := db.QueryRow(`SELECT UUID FROM "users" WHERE username == ?`, userSession.username)
		var uuid string
		query_id.Scan(&uuid)
		if sessionToken != uuid{
			http.Redirect(w, r, "/logout", http.StatusFound)
			return
		}
		data = map[string]interface{}{"postid": id, "postname": name, "posttext": text, "answerdata": Answers_data(id), "comentsdata": Comments_data(id), "username": userSession.username,"Like":like, "Dislike":dislike}
		tmpl := template.Must(template.ParseFiles("templates/post_page.html"))
		tmpl.Execute(w, data)
	}
}

func Answer_for_forum(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	query_id := db.QueryRow(`SELECT UUID FROM "users" WHERE username == ?`, userSession.username)
	var uuid string
	query_id.Scan(&uuid)
	if sessionToken != uuid{
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}
	ID := r.FormValue("postid")
	text := r.FormValue("answer")
	modules.Create_Answer(userSession.username, text, ID)
	http.Redirect(w, r, "/original_post/"+ID, http.StatusFound)
}
func Create_comment(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	query_id := db.QueryRow(`SELECT UUID FROM "users" WHERE username == ?`, userSession.username)
	var uuid string
	query_id.Scan(&uuid)
	if sessionToken != uuid{
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}
	ID := r.FormValue("postid")
	text := r.FormValue("comment")
	post_type := r.FormValue("link")
	modules.Create_comment(userSession.username, text, ID, post_type)
	http.Redirect(w, r, "/original_post/"+ID, http.StatusFound)
}
func Like(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	query_id := db.QueryRow(`SELECT UUID FROM "users" WHERE username == ?`, userSession.username)
	var uuid string
	query_id.Scan(&uuid)
	if sessionToken != uuid{
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}
	post_type := r.FormValue("link")
	id := r.FormValue("commentid")
	if post_type == "post_reactoion"{
		modules.Like(userSession.username, id, post_type)
		http.Redirect(w, r, r.FormValue("postlink"), http.StatusFound)
	}else if post_type == "answers_for_forum_reactoion"{
		db, _ := sql.Open("sqlite3", "./Database/database.db")
		usernmae := db.QueryRow(`SELECT POST_ID FROM "answers_for_forum" WHERE ID == ?`, r.FormValue("commentid"))
		var redirect_id string
		usernmae.Scan(&redirect_id)
		modules.Like(userSession.username, id, post_type)
		http.Redirect(w, r, r.FormValue("postlink")+redirect_id, http.StatusFound)
	}

}
func Dislike(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	query_id := db.QueryRow(`SELECT UUID FROM "users" WHERE username == ?`, userSession.username)
	var uuid string
	query_id.Scan(&uuid)
	if sessionToken != uuid{
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}
	post_type := r.FormValue("link")
	id := r.FormValue("commentid")
	if post_type == "post_reactoion"{
		modules.Dislike(userSession.username, id, post_type)
		http.Redirect(w, r, r.FormValue("postlink"), http.StatusFound)
	}else if post_type == "answers_for_forum_reactoion"{
		db, _ := sql.Open("sqlite3", "./Database/database.db")
		usernmae := db.QueryRow(`SELECT POST_ID FROM "answers_for_forum" WHERE ID == ?`, r.FormValue("commentid"))
		var redirect_id string
		usernmae.Scan(&redirect_id)
		modules.Dislike(userSession.username, id, post_type)
		http.Redirect(w, r, r.FormValue("postlink")+redirect_id, http.StatusFound)
	}
}
func Option(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	query_id := db.QueryRow(`SELECT UUID FROM "users" WHERE username == ?`, userSession.username)
	var uuid string
	query_id.Scan(&uuid)
	if sessionToken != uuid{
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}
	usernmae := db.QueryRow(`SELECT ID FROM "users" WHERE USERNAME == ?`, userSession.username)
	var id string
	usernmae.Scan(&id)
	data := map[string]interface{}{"postdata": Forums_data(r.FormValue("option"), r.FormValue("like"), r.FormValue("dislike"), id, r.FormValue("created_posts")), "username": userSession.username, "listdata": Options_data()}
	tmpl := template.Must(template.ParseFiles("templates/forum.html"))
	tmpl.Execute(w, data)
}
func create_session(username string, w http.ResponseWriter) {
	u1 := uuid.NewV4().String()
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	stmt, _ := db.Prepare(`
	UPDATE users
		set UUID= ? where username==?;
	`)
	stmt.Exec(u1,username)
	expiresAt := time.Now().Add(3600 * time.Second)
	sessions[u1] = session{
		username: username,
		expiry:   expiresAt,
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   u1,
		Expires: expiresAt,
	})
}
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}
func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	sessionToken := c.Value
	delete(sessions, sessionToken)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusFound)
	return
}
func Forums_data(option string, like string, dislike string, userid string, created_posts string) []Forum {
	var forums []Forum
	var option_id string
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	switch option {
	case "":
		switch created_posts {
		case "":
			switch true {
			case like == "" && dislike == "":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME,count(NULLIF(post_reactoion.LIKE, '')),count(NULLIF(post_reactoion.DISLIKE, '')) FROM forum Left JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID GROUP BY forum.POST_ID`)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name, &foru.Like, &foru.Dislike); err != nil {
						panic(err)
					}
					forums = append(forums, foru)
				}
			case like == "true" && dislike == "":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where post_reactoion.POST_COMENTATOR_ID==`+userid+` and post_reactoion.Like==1`)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						log.Fatal(err)
					}
					forums = append(forums, foru)
				}
			case like == "" && dislike == "true":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where post_reactoion.POST_COMENTATOR_ID==`+userid+` and post_reactoion.Dislike==1`)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						panic(err)
					}
					forums = append(forums, foru)
				}
			case like == "true" && dislike == "true":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where (post_reactoion.DISLIKE== 1 or post_reactoion.LIKE== 1) and post_reactoion.POST_COMENTATOR_ID==`+ userid)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						panic(err)
					}
					forums = append(forums, foru)
				}
			}
		default:
			switch true {
			case like == "" && dislike == "":
				forum, err := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum Left JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where forum.AUTHOR_ID==`+userid+` group by post_reactoion.POST_ID`)
				if(err != nil){
					log.Fatal(err)
				}
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						log.Fatal(err)
					}
					forums = append(forums, foru)
				}
			case like == "true" && dislike == "":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where post_reactoion.Like==1 and post_reactoion.POST_COMENTATOR_ID==` + userid + ` and forum.AUTHOR_ID==` + userid)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						panic(err)
					}
					forums = append(forums, foru)
				}
			case like == "" && dislike == "true":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where post_reactoion.Dislike==1 and post_reactoion.POST_COMENTATOR_ID==` + userid + ` and forum.AUTHOR_ID==` + userid)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						panic(err)
					}
					forums = append(forums, foru)
				}
			case like == "true" && dislike == "true":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where post_reactoion.POST_COMENTATOR_ID==` + userid + ` and forum.AUTHOR_ID==` + userid + ` and (post_reactoion.Dislike==1 or post_reactoion.Like==1)`)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						log.Fatal(err)
					}
					forums = append(forums, foru)
				}
			}
		}
	default:
		id, _ := db.Prepare(`SELECT ID FROM "categories" WHERE NAME == ?`)
		err2 := id.QueryRow(option).Scan(&option_id)
		if err2 != nil {
			fmt.Println("Error running command: post_ID")
		}
		switch created_posts {
		case "":
			switch true {
			case like == "" && dislike == "":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where forum_categories.CATEGORY_ID==` + option_id +` GROUP BY forum.POST_ID`)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						log.Fatal(err)
					}
					forums = append(forums, foru)
				}
			case like == "true" && dislike == "":
				forum, err := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where forum_categories.CATEGORY_ID==` + option_id + ` and post_reactoion.Like==1 and post_reactoion.POST_COMENTATOR_ID==` + userid)
				if err != nil {
					log.Fatal(err)
				}
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
					}
					forums = append(forums, foru)
				}
			case like == "" && dislike == "true":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where forum_categories.CATEGORY_ID==` + option_id + ` and post_reactoion.Dislike==1 and post_reactoion.POST_COMENTATOR_ID==` + userid)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						panic(err)
					}
					forums = append(forums, foru)
				}
			case like == "true" && dislike == "true":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where forum_categories.CATEGORY_ID==` + option_id + ` and post_reactoion.POST_COMENTATOR_ID==` + userid + ` and (post_reactoion.Dislike==1 or post_reactoion.Like==1)`)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						log.Fatal(err)
					}
					forums = append(forums, foru)
				}
			}
		default:
			switch true {
			case like == "" && dislike == "":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID==forum_categories.POST_ID Left JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where forum.AUTHOR_ID==`+userid+` and forum_categories.CATEGORY_ID==`+option_id+` GROUP BY forum.POST_ID`)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						log.Fatal(err)
					}
					forums = append(forums, foru)
				}
			case like == "true" && dislike == "":
				forum, err := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where post_reactoion.Like==1 and post_reactoion.POST_COMENTATOR_ID==` + userid + ` and forum.AUTHOR_ID==` + userid + ` and forum_categories.CATEGORY_ID==` + option_id)
				if(err != nil){
					log.Fatal(err)
				}
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						panic(err)
					}
					forums = append(forums, foru)
				}
			case like == "" && dislike == "true":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where post_reactoion.Dislike==1 and post_reactoion.POST_COMENTATOR_ID==` + userid + ` and forum.AUTHOR_ID==` + userid + ` and forum_categories.CATEGORY_ID==` + option_id)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						panic(err)
					}
					forums = append(forums, foru)
				}
			case like == "true" && dislike == "true":
				forum, _ := db.Query(`SELECT forum.POST_ID, forum.NAME FROM forum LEFT JOIN forum_categories on forum.POST_ID == forum_categories.POST_ID LEFT JOIN post_reactoion on forum.POST_ID == post_reactoion.POST_ID where post_reactoion.POST_COMENTATOR_ID==` + userid + ` and forum.AUTHOR_ID==` + userid + ` and forum_categories.CATEGORY_ID==` + option_id + ` and (post_reactoion.Dislike==1 or post_reactoion.Like==1)`)
				for forum.Next() {
					var foru Forum
					if err := forum.Scan(&foru.ID, &foru.Name); err != nil {
						log.Fatal(err)
					}
					forums = append(forums, foru)
				}
			}
		}
	}
	return forums
}
func Answers_data(id string) []Forum {
	var forums []Forum
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	forum, _ := db.Query(`SELECT answers_for_forum.ID, answers_for_forum.TEXT,count(NULLIF(answers_for_forum_reactoion.LIKE, '')),count(NULLIF(answers_for_forum_reactoion.DISLIKE, '')) FROM answers_for_forum LEFT JOIN answers_for_forum_reactoion on answers_for_forum.ID==answers_for_forum_reactoion.POST_ID where answers_for_forum.POST_ID==`+id+` group by answers_for_forum.ID`)
	for forum.Next() {
		var foru Forum
		if err := forum.Scan(&foru.ID, &foru.Text, &foru.Like, &foru.Dislike); err != nil {
			panic(err)
		}
		forums = append(forums, foru)
	}
	if err := forum.Err(); err != nil {
		panic(err)
	}
	return forums
}
func Comments_data(id string) []Forum {
	var forums []Forum
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	forum, _ := db.Query(`SELECT ID, TEXT FROM "comments" WHERE POST_ID == ` + id)
	for forum.Next() {
		var foru Forum
		if err := forum.Scan(&foru.ID, &foru.Text); err != nil {
			panic(err)
		}
		forums = append(forums, foru)
	}
	if err := forum.Err(); err != nil {
		panic(err)
	}
	return forums
}
func Options_data() []Options {
	var options []Options
	db, _ := sql.Open("sqlite3", "./Database/database.db")
	option, _ := db.Query(`SELECT ID, NAME FROM "categories"`)
	for option.Next() {
		var op Options
		if err := option.Scan(&op.ID, &op.Name); err != nil {
			panic(err)
		}
		options = append(options, op)
	}
	if err := option.Err(); err != nil {
		panic(err)
	}
	return options
}