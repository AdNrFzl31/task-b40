package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	route.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./asset/"))))

	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("GET")
	route.HandleFunc("/contactMe", contactMe).Methods("GET")
	route.HandleFunc("/projectDetail", projectDetail).Methods("GET")

	route.HandleFunc("/addProject", addProjectPost).Methods("POST")

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

	// w.Write([]byte("home"))
	// w.WriteHeader(http.StatusOK)

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

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf8")
	var tmpl, err = template.ParseFiles("contactMe.html")

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

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	data := map[string]interface{}{
		"Title":   "DUMBWAYS WEB APP",
		"Content": "Lorem ipsum dolor sit amet consectetur, adipisicing elit. Nulla optio pariatur quos doloremque neque vitae aliquam voluptate perferendis? Eaque enim quisquam ipsam unde, expedita saepe aliquid a praesentium est fuga.",
		"Id":      id,
	}

	tmpl.Execute(w, data)
}

func addProjectPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("projectName : " + r.PostForm.Get("name"))
	fmt.Println("startDate : " + r.PostForm.Get("startDate"))
	fmt.Println("endDate : " + r.PostForm.Get("endDate"))
	fmt.Println("description : " + r.PostForm.Get("description"))
	fmt.Println("technologies : " + r.PostForm.Get("icon"))
	fmt.Println("image : " + r.PostForm.Get("image"))

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

// go run main.go
