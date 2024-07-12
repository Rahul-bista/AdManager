package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Application struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	PackageName string `json:"package_name"`
	ImageURL    string `json:"image_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open(sqlite.Open("packages.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&Application{})

	router := http.NewServeMux()

	router.HandleFunc("/addPackage", addPackageHandler)
	router.HandleFunc("/packages", getAllPackagesHandler)
	router.HandleFunc("/deletePackage", deletePackageHandler)

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	router.HandleFunc("/", formHandler)

	fmt.Println("Server running at http://127.0.0.1:8000")
	server.ListenAndServe()
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "applicationForm.html")
}

func addPackageHandler(w http.ResponseWriter, r *http.Request) {
	var app Application
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	app.PackageName = r.FormValue("package_name")
	app.ImageURL = r.FormValue("image_url")
	app.Title = r.FormValue("title")
	app.Description = r.FormValue("description")
	app.URL = r.FormValue("url")

	db.Create(&app)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

func getAllPackagesHandler(w http.ResponseWriter, r *http.Request) {
	var apps []Application
	db.Find(&apps)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func deletePackageHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Fprintf(w, "id: %s",id)
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	var app Application
	result := db.Delete(&app, idInt)

	if result.Error != nil {
		http.Error(w, "Error deleting package", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
