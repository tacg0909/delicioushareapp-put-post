package main

import (
	"bytes"
	"encoding/base64"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/tacg0909/meshitero-put-post/resize"
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
    // UserId string `json:"userId"`
    // Base64Image string `json:"base64Image"`
    Base64Image string `json:"image"`
    // EatingPlace EatingPlace `json:"eatingPlace"`
}

type Event struct {
    UserId string
    PostedTime string
    ImageUrl string
}

func MeshiteroPutPost(putItem PutItem) (err error) {
    postId := uuid.NewString()

    imageBinary, err := base64.StdEncoding.DecodeString(putItem.Base64Image)
    if err != nil {
        return
    }
    largeImageBuf := bytes.NewBuffer(imageBinary)
    smallImageBuf := *largeImageBuf

    largeImage, err := resize.Resize(largeImageBuf, 1000)
    if err != nil {
        return
    }
    err = putImageToBucket("large/" + postId + ".jpg", largeImage)
    if err != nil {
        return
    }

    smallImage, err := resize.Resize(&smallImageBuf, 300)
    if err != nil {
        return
    }
    err = putImageToBucket("small/" + postId + ".jpg", smallImage)
    if err != nil {
        return
    }

    // db := dynamo.New(session.New(), &aws.Config{
    //     Region: aws.String(os.Getenv("DB_REGION")),
    // })
    // err = putUserPostToOutlineTable(db, putItem.UserId, putItem.ImageUrl, postId)
    return
}

func putImageToBucket(key string, image bytes.Buffer) (err error) {
    client := s3.New(session.New(), &aws.Config{
        Region: aws.String(os.Getenv("S3_BUCKET_REGION")),
    })
    _, err = client.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
        Key: aws.String(key),
        Body: bytes.NewReader(image.Bytes()),
        ContentType: aws.String("image/jpeg"),
    })
    return
}

type putUserPostOutlineEvent struct {
    UserId string
    PostedTime string
    SmallImageUrl string
    PostId string
}

func putUserPostToOutlineTable(db *dynamo.DB, userId string, smallImageUrl string, postId string) (err error) {
    table := db.Table(os.Getenv("USER_POSTS_OUTLINE_TABLE_NAME"))
    e := putUserPostOutlineEvent{
        UserId: userId,
        PostedTime: time.Now().UTC().Format("2006-01-02-15-04-05-0700"),
        SmallImageUrl: smallImageUrl,
        PostId: postId,
    }
    err = table.Put(e).If("attribute_not_exists(PostedTime)").Run()
    return
}
