package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	route := mux.NewRouter()

	connection.DataBaseConnect()

	route.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./asset/"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/contact-me", contactMe).Methods("GET")
	route.HandleFunc("/contact", AddContact).Methods("POST")
	route.HandleFunc("/project", addProject).Methods("GET")
	route.HandleFunc("/add-project", middleware.UploadFile(addProjectPost)).Methods("POST")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	route.HandleFunc("/add_project-update/{id}", submitEdit).Methods("POST")
	route.HandleFunc("/register", register).Methods("GET")
	route.HandleFunc("/submit-register", submitRegister).Methods("POST")
	route.HandleFunc("/login", login).Methods("GET")
	route.HandleFunc("/submit-login", submitLogin).Methods("POST")
	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

type SessionData struct {
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = SessionData{}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("index.html")
	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, f1 := range fm {

			flashes = append(flashes, f1.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, " ")
	println(flashes)

	if session.Values["IsLogin"] != true {
		println("login")
		data, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, tb_project.name, description, duratin,technologies,image FROM tb_project ORDER BY id DESC")
		var result []Project
		for data.Next() {
			var each = Project{}
			err := data.Scan(&each.Id, &each.NameProject, &each.Description, &each.Duration, &each.Technologies, &each.Image)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(each.Technologies)
			result = append(result, each)
		}

		resData := map[string]interface{}{
			"DataSession": Data,
			"Projects":    result,
		}
		w.WriteHeader(http.StatusOK)

		tmpl.Execute(w, resData)
	} else {

		sessionID := session.Values["ID"].(int)
		fmt.Println(sessionID)
		data, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, tb_project.name, description, duration,technologies,image FROM tb_project WHERE tb_project.author_id = $1 ORDER BY id DESC", sessionID)
		var result []Project
		for data.Next() {
			var each = Project{}
			err := data.Scan(&each.Id, &each.NameProject, &each.Description, &each.Duration, &each.Technologies, &each.Image)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(each.Technologies)
			result = append(result, each)
		}

		resData := map[string]interface{}{
			"DataSession": Data,
			"Projects":    result,
		}
		w.WriteHeader(http.StatusOK)

		tmpl.Execute(w, resData)
	}
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("public/contactMe.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("public/addProject.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

type Project struct {
	NameProject      string
	Duration         string
	StartDate        time.Time
	EndDate          time.Time
	Description      string
	Technologies     []string
	Image            string
	Author           string
	Id               int
	Format_StartDate string
	Format_EndDate   string
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("public/projectDetail.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(), "SELECT tb_project.id, tb_project.name, description,start_date,end_date,duration,image, tb_user.name as author FROM tb_project LEFT JOIN tb_user ON tb_project.author_id = tb_user.id WHERE tb_project.id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.NameProject, &ProjectDetail.Description, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Image, &ProjectDetail.Author)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	ProjectDetail.Format_StartDate = ProjectDetail.StartDate.Format("2 January 2006")
	ProjectDetail.Format_EndDate = ProjectDetail.EndDate.Format("2 January 2006")

	data := map[string]interface{}{
		"ProjectDetail": ProjectDetail,
	}

	tmpl.Execute(w, data)
}

func addProjectPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var nameProject = r.PostForm.Get("name")
	var description = r.PostForm.Get("description")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")

	var technologies []string
	technologies = r.Form["technologies"]

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	println(session)

	author := session.Values["ID"].(int)
	fmt.Println(author)

	timePost, _ := time.Parse("2006-01-02", startDate)
	timeNow, _ := time.Parse("2006-01-02", endDate)
	println(timeNow.String())
	println(timePost.String())

	hours := timeNow.Sub(timePost).Hours()
	days := hours / 24
	weeks := math.Floor(days / 7)
	months := math.Floor(days / 30)
	years := math.Floor(days / 365)

	var duration string

	if years > 0 {
		duration = strconv.FormatFloat(years, 'f', 0, 64) + " Years"
	} else if months > 0 {
		duration = strconv.FormatFloat(months, 'f', 0, 64) + " Months"
	} else if weeks > 0 {
		duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " Weeks"
	} else if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " Days"
	} else if hours > 0 {
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hours"
	} else {
		duration = "0 Days"
	}
	println(hours)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (name, description,start_date, end_date,duration,technologies,author_id,image) VALUES ($1, $2, $3, $4, $5,$6,$7,$8)", nameProject, description, timeNow, timePost, duration, technologies, author, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func AddContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Nama : " + r.PostForm.Get("nama"))
	fmt.Println("Email : " + r.PostForm.Get("email"))
	fmt.Println("Phone Number : " + r.PostForm.Get("phoneNumber"))
	fmt.Println("Subject : " + r.PostForm.Get("subject"))
	fmt.Println("Message : " + r.PostForm.Get("message"))
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("public/addproject-update.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	var ProjectDetail = Project{}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description, start_date, end_date FROM tb_project WHERE id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.NameProject, &ProjectDetail.Description, &ProjectDetail.StartDate, &ProjectDetail.EndDate)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data := map[string]interface{}{
		"editProject": ProjectDetail,
	}
	tmpl.Execute(w, data)
}

func submitEdit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var nameProject = r.PostForm.Get("name")
	var description = r.PostForm.Get("description")

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_project SET name = $1, description = $2, image = $3 WHERE id = $4", nameProject, description, image, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("public/register.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func submitRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("public/login.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func submitLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	user := User{}

	// mengambil data email, dan melakukan pengecekan email
	err = connection.Conn.QueryRow(context.Background(),
		"SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		fmt.Println("Email belum terdaftar")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : Email belum terdaftar " + err.Error()))
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)

		return
	}

	// melakukan pengecekan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Password salah")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : Password salah " + err.Error()))
		return
	}

	//berfungsi untuk menyimpan data kedalam sessions browser
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["ID"] = user.Id
	session.Values["IsLogin"] = true
	session.Options.MaxAge = 10800 // 3 JAM

	session.AddFlash("succesfull login", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout")
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// go run main.go
