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

/*
func SearchHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()
		searchQuery := r.FormValue("searchtext")

		// Now you can use the searchQuery to fetch data from the database or perform any other logic.
		db := MyDatabase.OpenConn()
		Generate_Row(db, searchQuery, "Search Results")

		tmplPath := filepath.Join("C:/Users/User/Desktop/Code/GOTH_STACK/templates", "SearchPage.html")
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

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
*/



func main() {
	
	db := MyDatabase.OpenConn()

	MyData = Data{
		Sections: []Section{},
	}

	Generate_Row(db, "esgrima", "Secao 1")
	Generate_Row(db, "livro", "secao 2")
	Generate_Row(db, "esgrima", "secao 3")
	Generate_Row(db, "livro", "secao 4")
	
	

	http.HandleFunc("/", handler)
	//http.HandleFunc("/SearchPage", SearchHandler)

	staticPath := "C:/Users/User/Desktop/Code/GOTH_STACK/static"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	http.ListenAndServe(":8103", nil)
}


func Generate_Row(db *sql.DB, SearchProduct string, RowTitle string){

	rows, _ := MyDatabase.DB_Search_and_Update(db, SearchProduct)

	
	PageSection = Section{
		Section_Title: RowTitle,
		Products: []MyDatabase.Product{},
	}
	
	ElementsPerRow := 6

	for i := 0; i < ElementsPerRow; i++{
		rows.Next()
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
