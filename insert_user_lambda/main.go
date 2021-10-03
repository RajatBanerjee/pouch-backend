package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/go-sql-driver/mysql"
)

type UserData struct {
	UserId          int    `json:"userId" db:"USER_ID"`
	FileName        string `json:"fileName" db:"FILE_NAME"`
	InsertTs        string `json:"insTs" db:"INS_TS"`
	UpdateTs        string `json:"updTs" db:"UPD_TS"`
	FileDescription string `json:"fileDesc" db:"FILE_DESC"`
}

func HandleLambdaEvent(userData UserData) (int64, error) {
	rds_host := os.Getenv("rds_host")
	name := os.Getenv("rds_user_name")
	password := os.Getenv("rds_password")
	db_name := os.Getenv("rds_db_name")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", name, password, rds_host, db_name))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)

	// perform a db.Query insert
	result, err := db.Exec("INSERT INTO pouch_db.USER_UPLOADS(USER_ID, FILE_NAME, FILE_DESC) VALUES ( ?,?,?)", userData.UserId, userData.FileName, userData.FileDescription)

	// if there is an error inserting, handle it
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	//
	//// Execute the query
	//err = db.QueryRow("SELECT  INS_TS, UPD_TS FROM pouch_db.USER_UPLOADS where USER_ID = ?", id).Scan(&userData.InsertTs, &userData.UpdateTs)
	//if err != nil {
	//	log.Fatal(err) // proper error handling instead of panic in your app
	//}

	return id, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
