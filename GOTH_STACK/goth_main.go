package main

import (
	"GOTH_STACK/MyDatabase"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)



type Data struct {
	Sections []Section
}

type Section struct{
	Section_Title string
	Products []MyDatabase.Product

}

var MyData Data
var PageSection Section

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

	MyData = Data{
		Sections: []Section{},
	}

	db := MyDatabase.OpenConn()
	Generate_Row(db, "esgrima", "Birds Of A Feather")
	Generate_Row(db, "esgrima", "Mitski")
	Generate_Row(db, "esgrima", "Compiladores")

	
	

	http.HandleFunc("/", handler)
	staticPath := "C:/Users/User/Desktop/Code/GOTH_STACK/static"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	http.ListenAndServe(":8102", nil)
}


func Generate_Row(db *sql.DB, SearchProduct string, RowTitle string){

	rows, _ := MyDatabase.DB_Search_and_Update(db, SearchProduct)

	
	PageSection = Section{
		Section_Title: RowTitle,
		Products: []MyDatabase.Product{},
	}
	
	
	for rows.Next(){
		var p MyDatabase.Product
		err := rows.Scan(&p.Title, &p.Price, &p.Reviews, &p.Imgurl, &p.Purl, &p.Lupdate, &p.Seller)
		if err != nil {
			fmt.Println("Error scanning row:", err)
				continue
		}

		PageSection.Products = append(PageSection.Products, p)
		
	  }
	  MyData.Sections = append(MyData.Sections, PageSection)
	  
}
