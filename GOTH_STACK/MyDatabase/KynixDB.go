package MyDatabase

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	_ "github.com/lib/pq" 
)

type Product struct {
	Title   string `json:"title"`
	Price   string `json:"price"`
	Reviews string `json:"reviews"`
	Imgurl  string `json:"imgurl"`
	Purl    string `json:"purl"`
	Lupdate string `json:"lupdate"`
	Seller  string `json:"seller "`
}





func OpenConn() *sql.DB{

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))


	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return db

}




func DBinsert(db *sql.DB, MyProduct Product) {

	InsertSql := `INSERT INTO product (title, price, reviews, imgurl, purl, lupdate, seller) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	
	_, err := db.Exec(InsertSql, MyProduct.Title, MyProduct.Price, MyProduct.Reviews, MyProduct.Imgurl, MyProduct.Purl, time.Now(), MyProduct.Seller)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(db, InsertSql)
	
	
}


func DB_Search_and_Update(db *sql.DB, search string) (sql.Rows, bool){
	

	QueryInput := "SELECT * FROM product WHERE product.title LIKE '%' || $1 || '%'"
	DBSearch, _ := db.Query(QueryInput, search)
	
	if !DBSearch.Next(){ //Logica da Shopee
		fmt.Println("No Items Found, Womp Womp")
		return sql.Rows{}, false
	} 

	return *DBSearch, true
	
}







