package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"personal-web/connection"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	connection.DataBaseConnect()

	route.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./asset/"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/contactMe", contactMe).Methods("GET")
	route.HandleFunc("/project", addProject).Methods("GET")
	route.HandleFunc("/projectDetail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/addProject", addProjectPost).Methods("POST")
	route.HandleFunc("/contact", AddContact).Methods("POST")
	route.HandleFunc("/deleteProject/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/editProject/{id}", editProject).Methods("GET")
	route.HandleFunc("/addproject-update/{id}", submitEdit).Methods("POST")

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("index.html")

	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id,name,description,duration FROM tb_project")

	var result []Project
	for data.Next() {
		var each = Project{}
		err := data.Scan(&each.Id, &each.Name, &each.Description, &each.Duration)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}
	resData := map[string]interface{}{
		"Projects": result,
	}

	tmpl.Execute(w, resData)
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("contactMe.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("addProject.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

type Project struct {
	Name             string
	Duration         string
	StartDate        time.Time
	EndDate          time.Time
	Description      string
	NodeJs           string
	Javascript       string
	ReactJs          string
	Html5            string
	Image            string
	Id               int
	Format_StartDate string
	Format_EndDate   string
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("projectDetail.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description, start_date, end_date, duration FROM tb_project WHERE id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Name, &ProjectDetail.Description, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration)
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

	var name = r.PostForm.Get("name")
	var description = r.PostForm.Get("description")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")

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

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (name, description,start_date, end_date,duration) VALUES ($1, $2, $3, $4, $5)", name, description, timeNow, timePost, duration)
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
	var tmpl, err = template.ParseFiles("addproject-update.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	var ProjectDetail = Project{}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, description, start_date, end_date FROM tb_project WHERE id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Name, &ProjectDetail.Description, &ProjectDetail.StartDate, &ProjectDetail.EndDate)

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

	var name = r.PostForm.Get("name")
	var description = r.PostForm.Get("description")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")

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

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_project SET name = $1, description = $2, timeNow = $3, timePost = $4, duration = $5 WHERE id = $4", name, description, timeNow, timePost, duration, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

// go run main.go
