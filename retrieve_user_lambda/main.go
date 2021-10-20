package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/go-sql-driver/mysql"
)

type FileInfo struct {
	Id              int    `json:"id" db:"id"`
	UserId          string `json:"userId" db:"user_id"`
	Username        string `json:"userName" db:"user_name"`
	FileName        string `json:"fileName" db:"file_name"`
	InsertTs        string `json:"insTs" db:"ins_ts"`
	UpdateTs        string `json:"updTs" db:"upd_ts"`
	FileDescription string `json:"fileDesc" db:"file_desc"`
	FilePath string `json:"filePath" db:"file_path"`
}

func HandleLambdaEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rds_host := os.Getenv("rds_host")
	name := os.Getenv("rds_user_name")
	password := os.Getenv("rds_password")
	db_name := os.Getenv("rds_db_name")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", name, password, rds_host, db_name))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)

	query := "SELECT id ,user_id,user_name,file_name, ins_ts, upd_ts,file_desc, file_path FROM pouch_db.USER_UPLOADS"


	id, ok := request.MultiValueQueryStringParameters["id"]
	if ok  {
			query = fmt.Sprintf("%s where user_id='%s'",query , id[0])

	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{Body: string(err.Error()), StatusCode: 500}, nil
	}
	defer rows.Close()

	var resp []FileInfo

	for rows.Next() {
		var a int
		var b string
		var c string
		var d string
		var e string
		var f string
		var g string
		var h string
		err := rows.Scan(&a, &b, &c, &d, &e, &f, &g, &h)
		if err != nil {
			log.Fatal(err)
			return events.APIGatewayProxyResponse{Body: string(err.Error()), StatusCode: 500}, nil
		}

		resp = append(resp, FileInfo{a, b, c, d, e, f, g, h})
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{Body: string(err.Error()), StatusCode: 500}, nil
	}

	response, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{Body: string(err.Error()), StatusCode: 500}, nil
	}



	return events.APIGatewayProxyResponse{Body: string(response),Headers: map[string]string{
		"Access-Control-Allow-Headers" : "Content-Type",
		"Access-Control-Allow-Origin": "*",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET",
	}, StatusCode: 200}, nil
}
func main() {
	lambda.Start(HandleLambdaEvent)
}
