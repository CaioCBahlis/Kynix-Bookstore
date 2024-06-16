package main

import (
	"GOTH_STACK/MyDatabase"
	"GOTH_STACK/Scrappers"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)



type Data struct {
	Sections [][]Section
}

type SData struct{
	Section_Title string
	Products []Scrappers.Product
}

type Exhibtion struct{
	Section_Title string
	Products []Scrappers.Book
}


type Section struct{
	Section_Title string
	Products []Scrappers.Product

}


var MyData Data
var SearchData SData
var OpenBooksData Data
var PageSection Section

func main() {

	//db := MyDatabase.OpenConn()

	staticPath := "C:/Users/User/Desktop/Code/GOTH_STACK/static"

	
	http.HandleFunc("/", handler)
	http.HandleFunc("/SearchPage", SearchHandler)
	http.HandleFunc("/OpenBooks", OpenBooks)
	http.HandleFunc("/OpenSearchHandler", OpenSearchHandler)
	http.HandleFunc("/ProductDownload", ProductHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	

	
	
	http.ListenAndServe(":8322", nil)
	
}


func handler(w http.ResponseWriter, r *http.Request) {

	
	MyData = Data{}
	db := MyDatabase.OpenConn()
	//Top 10 maiores BottleNecks dos Animes
	//Usar Singleton
	

	Generate_HUB_Row(db, "Harry Potter", "Best Sellers", &MyData)
	Generate_HUB_Row(db, "Harry Potter", "Recommended", &MyData)
	Generate_HUB_Row(db, "Computer", "Fiction & Fantasy", &MyData)
	Generate_HUB_Row(db, "esgrima", "On Sale", &MyData)

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

	db := MyDatabase.OpenConn()
	SearchData = SData{}
	
	

	if r.Method == http.MethodPost {
		r.ParseForm()
		searchQuery := r.FormValue("searchtext")

		
		Generate_Row(db, searchQuery, "Secao1", "Product")

		}else{
			fmt.Println("No Pages Found")
		}

		tmplPath := filepath.Join("C:/Users/User/Desktop/Code/GOTH_STACK/templates", "SearchStore.html", "")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, SearchData)
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			return
		}else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
}


func OpenBooks(w http.ResponseWriter, r *http.Request){

		db := MyDatabase.OBConn()
		Generate_HUB_Row(db, "Rust", "Rusted", &OpenBooksData)
			

		tmplPath := filepath.Join("C:/Users/User/Desktop/Code/GOTH_STACK/templates", "OpenBooks.html", "")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, OpenBooksData)
		OpenBooksData.Sections = [][]Section{}
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			return
		}
}

func OpenSearchHandler(w http.ResponseWriter, r *http.Request){

	db := MyDatabase.OBConn()
	
	MyExhibition := Exhibtion{}
	if r.Method == http.MethodPost{
		p := r.FormValue("OpenCategories")
		KeyWords := strings.Split(p, ",")
		Content := Generate_Open_Categories(db, KeyWords, "OpenBooks")
		

		tmplPath := filepath.Join("C:/Users/User/Desktop/Code/GOTH_STACK/templates", "SearchPage.html", "")
		tmpl, _ := template.ParseFiles(tmplPath)


		for value := range Content{
			MyExhibition.Products = append(MyExhibition.Products, value)
			
		}
		MyExhibition.Section_Title = MyExhibition.Products[0].Title


	
		tmpl.Execute(w, MyExhibition)
	}
}

func ProductHandler(w http.ResponseWriter, r *http.Request){

	db := MyDatabase.OBConn()
	ObjectTitle := r.FormValue("ProductObject")
	var SearchedBook Scrappers.Book
	
	

	rows, _ := MyDatabase.DB_Search_and_Update(db, "OpenBooks", ObjectTitle)
	
	rows.Scan(&SearchedBook.Title, &SearchedBook.Imgurl, &SearchedBook.Link1, &SearchedBook.Link2, &SearchedBook.Link3, &SearchedBook.Link4, &SearchedBook.Link5, &SearchedBook.Link6, &SearchedBook.Link7, time.Now())
	
	tmplPath := filepath.Join("C:/Users/User/Desktop/Code/GOTH_STACK/templates", "ProductPage.html", "")
	tmpl, _ := template.ParseFiles(tmplPath)
	tmpl.Execute(w,SearchedBook)
}


func Generate_Open_Categories(db *sql.DB, KeyWords []string, table string) map[Scrappers.Book]bool{

	Content := make(map[Scrappers.Book]bool)
	
	for _, value := range(KeyWords){
	
		rows, _ := MyDatabase.DB_Search_and_Update(db, table, value)
		for rows.Next(){
			MyBook := Scrappers.Book{}
			rows.Scan(&MyBook.Title, &MyBook.Imgurl, &MyBook.Link1, &MyBook.Link2, &MyBook.Link3, &MyBook.Link4, &MyBook.Link5, &MyBook.Link6, &MyBook.Link7, time.Now())
			if !Content[MyBook]{
				Content[MyBook] = true
			}			
		}
	}
	
	return Content

}

func Generate_Row(db *sql.DB, SearchProduct string, RowTitle string, table string){

	
	rows, _ := MyDatabase.DB_Search_and_Update(db, table, SearchProduct)
	
	
	for rows.Next(){
		var p Scrappers.Product
		err := rows.Scan(&p.Title, &p.Price, &p.Reviews, &p.Imgurl, &p.Purl, &p.Lupdate, &p.Seller)
		if err != nil {
			fmt.Println("Error scanning row:", err)
				continue
		}
	

		SearchData.Products = append(SearchData.Products, p)
    }


}
	

	func Generate_HUB_Row(db *sql.DB, RowTopic string, SectionTitle string, PageDS *Data){

	var table string
	if PageDS == &OpenBooksData{
		table = "OpenBooks"
	}else{
		table = "MyData" 
	}

	rows, _ :=  MyDatabase.DB_Search_and_Update(db, table, RowTopic)

	PageSection = Section{
		Section_Title: SectionTitle,
		Products: []Scrappers.Product{},
	}


	var HUBrow []Section

	
	for i := 0; i < 5; i++{
		rows.Next()

		var p Scrappers.Product
		err := rows.Scan(&p.Title, &p.Price, &p.Reviews, &p.Imgurl, &p.Purl, &p.Lupdate, &p.Seller)
		if err != nil{
			fmt.Println("Error scanning values for main page")
		}

		PageSection.Products = append(PageSection.Products, p)
	}

	HUBrow = append(HUBrow, PageSection)

	PageDS.Sections = append(PageDS.Sections, HUBrow)
	
}



