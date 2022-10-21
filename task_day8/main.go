package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	route.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./asset/"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/contactMe", contactMe).Methods("GET")
	route.HandleFunc("/project", addProject).Methods("GET")
	route.HandleFunc("/projectDetail/{index}", projectDetail).Methods("GET")
	route.HandleFunc("/addProject", addProjectPost).Methods("POST")
	route.HandleFunc("/contact", AddContact).Methods("POST")
	route.HandleFunc("/deleteProject/{index}", deleteProject).Methods("GET")
	route.HandleFunc("/editProject/{index}", editProject).Methods("GET")
	route.HandleFunc("/update-project/{Id}", submitEdit).Methods("POST")

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

	response := map[string]interface{}{
		"Projects": dataProject,
	}

	// w.Write([]byte("home"))
	// w.WriteHeader(http.StatusOK)

	tmpl.Execute(w, response)
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

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("projectDetail.html")

	if err != nil {
		w.Write([]byte("massage : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if index == i {
			ProjectDetail = Project{
				Name:        data.Name,
				Duration:    data.Duration,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Description: data.Description,
				NodeJs:      data.NodeJs,
				Javascript:  data.Javascript,
				ReactJs:     data.ReactJs,
				Html5:       data.Html5,
				Image:       data.Image,
			}
		}
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpl.Execute(w, data)
}

var dataProject = []Project{}

type Project struct {
	Name        string
	Duration    string
	StartDate   string
	EndDate     string
	Description string
	NodeJs      string
	Javascript  string
	ReactJs     string
	Html5       string
	Image       string
	Id          int
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
	var nodeJs = r.PostForm.Get("nodeJs")
	var javascript = r.PostForm.Get("javascript")
	var reactJs = r.PostForm.Get("reactJs")
	var html5 = r.PostForm.Get("html5")
	var image = r.PostForm.Get("image")

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

	var newPoject = Project{
		Name:        name,
		Duration:    duration,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		NodeJs:      nodeJs,
		Javascript:  javascript,
		ReactJs:     reactJs,
		Html5:       html5,
		Image:       image,
		Id:          len(dataProject),
	}
	fmt.Println(newPoject)

	dataProject = append(dataProject, newPoject)
	// fmt.Println(dataProject)

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
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject = append(dataProject[:index], dataProject[index+1:]...)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("addproject-update.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	var projectDetail = Project{}
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	for i, data := range dataProject {
		if index == i {
			projectDetail = Project{
				Name:        data.Name,
				Duration:    data.Duration,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Description: data.Description,
				NodeJs:      data.NodeJs,
				Javascript:  data.Javascript,
				ReactJs:     data.ReactJs,
				Html5:       data.Html5,
				Image:       data.Image,
				Id:          data.Id,
			}

		}
	}
	data := map[string]interface{}{
		"editProject": projectDetail,
	}
	tmpl.Execute(w, data)
}

func submitEdit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("name")
	var description = r.PostForm.Get("description")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")
	var nodeJs = r.PostForm.Get("nodeJs")
	var javascript = r.PostForm.Get("javascript")
	var reactJs = r.PostForm.Get("reactJs")
	var html5 = r.PostForm.Get("html5")
	var image = r.PostForm.Get("image")

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

	var newPoject = Project{
		Name:        name,
		Duration:    duration,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		NodeJs:      nodeJs,
		Javascript:  javascript,
		ReactJs:     reactJs,
		Html5:       html5,
		Image:       image,
		Id:          len(dataProject),
	}
	fmt.Println(newPoject)

	dataProject = append(dataProject, newPoject)
	// fmt.Println(dataProject)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

// go run main.go
