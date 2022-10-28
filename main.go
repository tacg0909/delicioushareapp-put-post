package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

func main() {
    lambda.Start(MeshiteroPutPost)
}

type EatingPlace struct {
    Name string `json:"name"`
    Address string `json:"address"`
    Website string `json:"website"`
    Id string `json:"id"`
}

type PutItem struct {
    UserId string `json:"userId"`
    Base64Image string `json:"base64Image"`
    EatingPlace EatingPlace `json:"eatingPlace"`
}

func MeshiteroPutPost(putItem PutItem) error {
    postId := uuid.NewString()

    postedTime := time.Now().UTC().Format("2006-01-02-15-04-05-0700")

    db := dynamo.New(session.New(), &aws.Config{
        Region: aws.String(os.Getenv("DB_REGION")),
    })
    err := putUserPostToDetailTable(
        db,
        postId,
        putItem.EatingPlace,
        fmt.Sprintf("%s/large/%s.jpg", os.Getenv("S3_BUCKET_OBJECT_DOMAIN"), postId),
        putItem.UserId,
        postedTime,
    )
    if err != nil {
        return err
    }
    err = putUserPostToOutlineTable(
        db,
        putItem.UserId,
        fmt.Sprintf("%s/small/%s.jpg", os.Getenv("S3_BUCKET_OBJECT_DOMAIN"), postId),
        postId,
        postedTime,
    )
    if err != nil {
        return err
    }

    err = putImage(putItem.Base64Image, postId)
    return err
}
