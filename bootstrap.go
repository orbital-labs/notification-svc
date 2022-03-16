package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
)

func initEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}
	return nil
}

func initAWS() *session.Session {
	key := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(
		&aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, ""),
			Region:      aws.String("ap-southeast-1"),
		})
	if err != nil {
		panic(err)
	}

	return sess
}
