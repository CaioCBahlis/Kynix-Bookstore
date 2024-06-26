package main

import (
	"GOTH_STACK/MyDatabase"
	"GOTH_STACK/Scrappers"
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

	

	staticPath := "C:/Users/User/Desktop/Code/GOTH_STACK/static"

	
	http.HandleFunc("/", handler)
	http.HandleFunc("/SearchPage", SearchHandler)
	http.HandleFunc("/OpenBooks", OpenBooks)
	http.HandleFunc("/OpenSearchHandler", OpenSearchHandler)
	http.HandleFunc("/ProductDownload", ProductHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	

	
	fmt.Println("Servin")
	http.ListenAndServe(":8322", nil)
	fmt.Println("Not Serving")
	
	
}


func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("handler")
	MyData = Data{}
	
	Generate_HUB_Row("Harry Potter", "Best Sellers", &MyData)
	Generate_HUB_Row("Harry Potter", "Recommended", &MyData)
	Generate_HUB_Row("Computer", "Fiction & Fantasy", &MyData)
	Generate_HUB_Row("esgrima", "On Sale", &MyData)

	fmt.Println("Content Served")

	tmplPath := filepath.Join("C:/Users/User/Desktop/Code/GOTH_STACK/templates", "myhtml.html", "")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	err = tmpl.Execute(w, MyData)
	fmt.Println("Content Executed")
	
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}	


func SearchHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Search handler")
	SearchData = SData{}
	
	

	if r.Method == http.MethodPost {
		r.ParseForm()
		searchQuery := r.FormValue("searchtext")

		
		Generate_Row(searchQuery, "Secao1")

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

	fmt.Println("Open Handler")
		Generate_HUB_Row("Rust", "Rusted", &OpenBooksData)
			

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

	fmt.Println("Open Search handler")
	MyExhibition := Exhibtion{}
	if r.Method == http.MethodPost{
		p := r.FormValue("OpenCategories")
		KeyWords := strings.Split(p, ",")
		Content := Generate_Open_Categories(KeyWords)
		

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

	fmt.Println("Product handler")
 	OpenDB, err := MyDatabase.Get_Open_DB()
	if err != nil{
		panic(err)
	}
	
	ObjectTitle := r.FormValue("ProductObject")
	var SearchedBook Scrappers.Book
	
	

	rows, _ := OpenDB.DB_Search(ObjectTitle)
	
	rows.Scan(&SearchedBook.Title, &SearchedBook.Imgurl, &SearchedBook.Link1, &SearchedBook.Link2, &SearchedBook.Link3, &SearchedBook.Link4, &SearchedBook.Link5, &SearchedBook.Link6, &SearchedBook.Link7, time.Now())
	
	tmplPath := filepath.Join("C:/Users/User/Desktop/Code/GOTH_STACK/templates", "ProductPage.html", "")
	tmpl, _ := template.ParseFiles(tmplPath)
	tmpl.Execute(w,SearchedBook)
}


func Generate_Open_Categories(KeyWords []string) map[Scrappers.Book]bool{

	OpenDB, err := MyDatabase.Get_Open_DB()
	if err != nil{
		panic(err)
	}

	Content := make(map[Scrappers.Book]bool)
	
	for _, value := range(KeyWords){
	
		rows, _ := OpenDB.DB_Search(value)
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

func Generate_Row(SearchProduct string, RowTitle string){

	KynixDB, err := MyDatabase.Get_Kynix_DB()
	if err != nil{
		panic(err)
	}
	rows, _ := KynixDB.DB_Search(SearchProduct)
	
	
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
	

	func Generate_HUB_Row(RowTopic string, SectionTitle string, PageDS *Data){


	var db MyDatabase.MyDatabase
	var err error

	if PageDS == &OpenBooksData{
		db, err = MyDatabase.Get_Open_DB()
	}else{
		db, err = MyDatabase.Get_Kynix_DB()
	}

	if err != nil{
		panic(err)
	}
	

	rows, _ :=  db.DB_Search(RowTopic)

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



