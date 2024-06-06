package main

import (
	"GOTH_STACK/MyDatabase"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)



type data struct {
	Title  string
	Products []MyDatabase.Product
}

var MyData data

func handler(w http.ResponseWriter, r *http.Request) {
	

	tmplPath := filepath.Join("C:/Users/User/Desktop/Code/GOTH_STACK/templates", "myhtml.html", "")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, MyData)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}


func main() {
	db := MyDatabase.OpenConn()
	rows, _ := MyDatabase.DB_Search_and_Update(db, "esgrima")

	MyData = data{
		Title: "MyTitle",
		Products: []MyDatabase.Product{},
	}
	
	
	for rows.Next(){
		var p MyDatabase.Product
		err := rows.Scan(&p.Title, &p.Price, &p.Reviews, &p.Imgurl, &p.Purl, &p.Lupdate, &p.Seller)
		if err != nil {
			fmt.Println("Error scanning row:", err)
				continue
		}

		MyData.Products = append(MyData.Products, p)
	  }
	
	

	http.HandleFunc("/", handler)
	staticPath := "C:/Users/User/Desktop/Code/GOTH_STACK/static"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	http.ListenAndServe(":8102", nil)
}
