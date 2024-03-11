package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getSecretValue(secretName string) (string, error) {
	sess, _ := session.NewSession()
	svc := secretsmanager.New(sess)

	result, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}

func main() {
	// .envファイルから環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// 環境変数を取得
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	// 環境変数を使用してAWSの設定を行う
	// ...

	// Ginアプリケーションの設定
	r := gin.Default()

	// フロントエンドのルーティング
	r.Static("/", "./frontend/public")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/.next/server/pages" + c.Request.URL.Path)
	})

	// バックエンドのルーティング
	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from backend!",
		})
	})

	r.Run(":8080")
}
