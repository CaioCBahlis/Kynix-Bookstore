package MyDatabase


import (
	"GOTH_STACK/Scrappers"
	"database/sql"
	"fmt"
	"os"
	"time"
	_ "github.com/lib/pq" 
)

func OBConn() *sql.DB{
	 
	DB_DATA := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", DB_DATA)
	if err != nil {
		panic(err)
	}
	
	return db

}



func OB_Insert(db *sql.DB, MyBook *Scrappers.Book){

	QUERY := "INSERT INTO openbooks (title, imgurl, link1, link2, link3, link4, link5, link6, link7, lupdate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"

	_, err := db.Exec(QUERY, MyBook.Title, MyBook.Imgurl, MyBook.Link1, MyBook.Link2, MyBook.Link3, MyBook.Link4, MyBook.Link5, MyBook.Link6, MyBook.Link7, time.Now())
	if err != nil{
		fmt.Println("error adding to the openbook database, womp womp")
	}
}

