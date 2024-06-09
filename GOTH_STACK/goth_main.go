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
	Sections [][]Section
}


type Section struct{
	Section_Title string
	Products []MyDatabase.Product

}



var MyData Data
var SearchData Data
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


func SearchHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()
		searchQuery := r.FormValue("searchtext")

		// Now you can use the searchQuery to fetch data from the database or perform any other logic.
		db := MyDatabase.OpenConn()
		Generate_Row(db, searchQuery, "Secao1", &SearchData)

		tmplPath := filepath.Join("C:/Users/User/Desktop/Code/GOTH_STACK/templates", "SearchPage.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, SearchData)
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		SearchData = Data{
			Sections: [][]Section{},
		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}




func main() {
	
	db := MyDatabase.OpenConn()

	MyData = Data{
		Sections: [][]Section{},
	}

	SearchData = Data{
		Sections: [][]Section{},
	}

	Generate_Row(db, "livro", "Best Seller", &MyData)
	Generate_Row(db, "esgrima", "Recomendacoes", &MyData)
	Generate_Row(db, "Computer", "Ficcao  Fantasia", &MyData)
	Generate_Row(db, "esgrima", "Promocao", &MyData)
	
	
	http.HandleFunc("/", handler)
	http.HandleFunc("/SearchPage", SearchHandler)

	staticPath := "C:/Users/User/Desktop/Code/GOTH_STACK/static"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	http.ListenAndServe(":8222", nil)
}


func Generate_Row(db *sql.DB, SearchProduct string, RowTitle string, PageDS *Data){

	rows, _ := MyDatabase.DB_Search_and_Update(db, SearchProduct)
	
	PageSection = Section{
		Section_Title: RowTitle,
		Products: []MyDatabase.Product{},
	}
	

	var ElementsPerRow int
	if PageDS == &SearchData{
		ElementsPerRow = 5
	}else{
		ElementsPerRow = 5
	}
	
	var rowSections []Section
	
	for rows.Next(){
		var p MyDatabase.Product
		err := rows.Scan(&p.Title, &p.Price, &p.Reviews, &p.Imgurl, &p.Purl, &p.Lupdate, &p.Seller)
		if err != nil {
			fmt.Println("Error scanning row:", err)
				continue
		}
		PageSection.Products = append(PageSection.Products, p)
		
		if len(PageSection.Products) >= ElementsPerRow {
            rowSections = append(rowSections, PageSection)
			PageDS.Sections = append(PageDS.Sections, rowSections)
			rowSections = []Section{}

            PageSection = Section{
                Section_Title: RowTitle,
                Products:      []MyDatabase.Product{},
            }
        }
	}

	
	if PageDS == &SearchData{
		if len(PageSection.Products) > 0 {
		rowSections = append(rowSections, PageSection)
			PageDS.Sections = append(PageDS.Sections, rowSections)
		}	
	}

}



