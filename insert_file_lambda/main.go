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

type FileInfo struct {
	UserId          string    `json:"userId" db:"user_id"`
	Username          string    `json:"userName" db:"user_name"`
	FileName        string `json:"fileName" db:"file_name"`
	InsertTs        string `json:"insTs" db:"ins_ts"`
	UpdateTs        string `json:"updTs" db:"upd_ts"`
	FileDescription string `json:"fileDesc" db:"file_desc"`
	FilePath string `json:"filePath" db:"file_path"`
}

func HandleLambdaEvent(fileInfo FileInfo) (int64, error) {
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
	result, err := db.Exec("INSERT INTO pouch_db.USER_UPLOADS(user_id, user_name, file_name, file_desc, file_path, ins_ts, upd_ts) VALUES ( ?, ?,?, ?, ?, ?, ?)", fileInfo.UserId,fileInfo.Username, fileInfo.FileName, fileInfo.FileDescription, fileInfo.FilePath, time.Now().String(), time.Now().String())

	// if there is an error inserting, handle it
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return id, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
