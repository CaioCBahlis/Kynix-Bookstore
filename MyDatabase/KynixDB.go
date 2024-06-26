package MyDatabase

import (
	"GOTH_STACK/Scrappers"
	"database/sql"
	"fmt"
	"time"
	"os"
	_ "github.com/lib/pq" 
)


type MyDatabase interface{
	DB_Insert(item Scrappers.Item) error
	DB_Search(string) (*sql.Rows, error)
}

type OpenBookDB struct{
	Connection *sql.DB
	
}

type KynixDB struct{
	Connection *sql.DB
}

var OpenConnection MyDatabase
var KynixConnection MyDatabase


func Get_Kynix_DB() (MyDatabase, error){

	if KynixConnection != nil{
		fmt.Println("There's already an instance of this connection open")
		return KynixConnection, nil
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))


	CONNECTION, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	KynixConnection := &KynixDB{Connection: CONNECTION}
	return KynixConnection, nil
	
}


func Get_Open_DB() (MyDatabase, error){

	if OpenConnection != nil{
		fmt.Println("There's already an instance of this connection open")
		return OpenConnection, nil
	}

	DB_DATA := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	CONNECTION, err := sql.Open("postgres", DB_DATA)
	if err != nil {
		return nil, err
	}

	OpenConnection := &OpenBookDB{Connection: CONNECTION}
	return OpenConnection, nil
}


func (ky *KynixDB) DB_Insert(Item Scrappers.Item) error{

	
	MyProduct := Item.(Scrappers.Product)
	InsertSql := `INSERT INTO product (title, price, reviews, imgurl, purl, lupdate, seller) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	
	_, err := ky.Connection.Exec(InsertSql, MyProduct.Title, MyProduct.Price, MyProduct.Reviews, MyProduct.Imgurl, MyProduct.Purl, time.Now(), MyProduct.Seller)
	if err != nil {
		return err
	}
	

	return nil	
}


func (ky *KynixDB) DB_Search(search string) (*sql.Rows, error){
	
	QueryInput := "SELECT * FROM product WHERE product.title LIKE '%' || $1 || '%'"
	

	//fmt.Println(QueryInput, search)
	DBSearch, err := ky.Connection.Query(QueryInput, search)
	if err != nil{
		return nil, nil
	}
	
	if !DBSearch.Next(){ //Logica da Shopee
		fmt.Println("No Items Found, Womp Womp")
		return &sql.Rows{}, nil
	} 

	return DBSearch, nil
	
}


func (ob *OpenBookDB) DB_Insert(Item Scrappers.Item) error{

	MyBook := Item.(Scrappers.Book)

	QUERY := "INSERT INTO openbooks (title, imgurl, link1, link2, link3, link4, link5, link6, link7, lupdate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"

	_, err := ob.Connection.Exec(QUERY, MyBook.Title, MyBook.Imgurl, MyBook.Link1, MyBook.Link2, MyBook.Link3, MyBook.Link4, MyBook.Link5, MyBook.Link6, MyBook.Link7, time.Now())
	if err != nil{
		fmt.Println("Error Adding Book to the Database")
		return err
	}

	return nil
}

func (ob *OpenBookDB) DB_Search(search string) (*sql.Rows, error){
	

	QueryInput := "SELECT * FROM OpenBooks WHERE OpenBooks.title LIKE '%' || $1 || '%'"
	
	


	fmt.Println(QueryInput, search)
	DBSearch, err := ob.Connection.Query(QueryInput, search)
	if err != nil{
		fmt.Println("Error Adding to the DataBase")
	}
	
	if !DBSearch.Next(){ //Logica da Shopee
		fmt.Println("No Items Found, Womp Womp")
		return &sql.Rows{}, nil
	} 

	return DBSearch, nil
	
}








