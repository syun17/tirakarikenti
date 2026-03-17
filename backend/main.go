package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	rekognitionTypes "github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type ScoreHistory struct {
	ID        string `dynamodbav:"id" json:"id"`
	Score     int    `dynamodbav:"score" json:"score"`
	ImageURL  string `dynamodbav:"image_url" json:"image_url"`
	CreatedAt string `dynamodbav:"created_at" json:"created_at"`
}

var (
	s3Client          *s3.Client
	rekognitionClient *rekognition.Client
	dynamoClient      *dynamodb.Client
	s3Bucket          string
	dynamoTable       = "RoomScoreHistory"
)

func initAWS() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error reading it, using environment variables")
	}

	s3Bucket = os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Println("S3_BUCKET environment variable is not set. Assuming running in container with env_file applied.")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1" // Fallback to us-east-1 if not specified
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	token := os.Getenv("AWS_SESSION_TOKEN")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, token)),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	rekognitionClient = rekognition.NewFromConfig(cfg)
	dynamoClient = dynamodb.NewFromConfig(cfg)
}

func main() {
	initAWS()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	r.POST("/analyze", func(c *gin.Context) {
		fileHeader, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open image"})
			return
		}
		defer file.Close()

		// 1. S3にアップロード
		objectKey := uuid.New().String() + "-" + fileHeader.Filename
		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("failed to upload image to s3: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image to s3"})
			return
		}

		imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s3Bucket, objectKey)

		// 2. Rekognitionで解析
		rekOutput, err := rekognitionClient.DetectLabels(context.TODO(), &rekognition.DetectLabelsInput{
			Image: &rekognitionTypes.Image{
				S3Object: &rekognitionTypes.S3Object{
					Bucket: aws.String(s3Bucket),
					Name:   aws.String(objectKey),
				},
			},
			MaxLabels:     aws.Int32(10),
			MinConfidence: aws.Float32(70),
		})
		if err != nil {
			log.Printf("failed to detect labels: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze image"})
			return
		}

		for _, label := range rekOutput.Labels {
			fmt.Println("Label:", aws.ToString(label.Name), aws.ToFloat32(label.Confidence))
		}

		for _, label := range rekOutput.Labels {
			for _, instance := range label.Instances {
				box := instance.BoundingBox
				fmt.Println(
					aws.ToString(label.Name),
					aws.ToFloat32(box.Left),
					aws.ToFloat32(box.Top),
				)
			}
		}

		// 3. ラベル数で簡易スコア生成
		labelCount := len(rekOutput.Labels)
		score := labelCount * 10
		if score > 100 {
			score = 100
		}

		// 4. DynamoDBに保存
		record := ScoreHistory{
			ID:        uuid.New().String(),
			Score:     score,
			ImageURL:  imageURL,
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		av, err := attributevalue.MarshalMap(record)
		if err != nil {
			log.Printf("failed to marshal dynamodb record: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare database record"})
			return
		}

		_, err = dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(dynamoTable),
			Item:      av,
		})
		if err != nil {
			log.Printf("failed to put item to dynamodb: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save history"})
			return
		}

		// 5. 返却
		c.JSON(http.StatusOK, gin.H{
			"score":     score,
			"image_url": imageURL,
		})
	})

	r.GET("/history", func(c *gin.Context) {
		out, err := dynamoClient.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String(dynamoTable),
		})
		if err != nil {
			log.Printf("failed to scan dynamodb: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get history"})
			return
		}

		var history []ScoreHistory
		err = attributevalue.UnmarshalListOfMaps(out.Items, &history)
		if err != nil {
			log.Printf("failed to unmarshal dynamodb results: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse history"})
			return
		}

		// 要素がない場合は空配列を返す
		if history == nil {
			history = []ScoreHistory{}
		}

		c.JSON(http.StatusOK, history)
	})

	r.Run(":8080")
}
