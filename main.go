package main

import (
	"log"

	"github.com/antklim/go-dynamodb/dyno"
	"github.com/antklim/go-dynamodb/invoice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("ap-southeast-2")}))
	client := dynamodb.New(sess)
	repo := dyno.NewRepository(client)
	service := invoice.NewService(repo)

	inv, err := service.GetInvoice("123")
	log.Println(inv, err)
}
