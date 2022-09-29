package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

func main() {
    lambda.Start(MeshiteroPutPost)
}

type PutItem struct {
    UserId string `json:"userId"`
    ImageUrl string `json:"imageUrl"`
}

type Event struct {
    UserId string
    PostedTime string
    ImageUrl string
}

func MeshiteroPutPost(c context.Context, putItem PutItem) (string, error) {
    db := dynamo.New(session.New(), &aws.Config{
        Region: aws.String("ap-northeast-1"),
    })
    table := db.Table("MeshiteroPostsTable")

    event := Event{
        UserId: putItem.UserId,
        PostedTime: time.Now().UTC().Format("2006-01-02-15-04-05-0700"),
        ImageUrl: putItem.ImageUrl,
    }
    err := table.Put(event).If("attribute_not_exists(PostedTime)").Run()
    if err != nil {
        return "", err
    }

    return "{success}", nil
}
