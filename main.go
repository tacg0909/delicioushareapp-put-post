package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
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
}

type PutItem struct {
    UserId string `json:"userId"`
    Base64Image string `json:"base64Image"`
    EatingPlace EatingPlace `json:"eatingPlace"`
    // Base64Image string `json:"image"`
}

func MeshiteroPutPost(putItem PutItem) error {
    postId := uuid.NewString()

    db := dynamo.New(session.New(), &aws.Config{
        Region: aws.String(os.Getenv("DB_REGION")),
    })
    err := putUserPostToDetailTable(
        db,
        postId,
        putItem.EatingPlace,
        fmt.Sprintf("%s/large/%s.jpg", os.Getenv("S3_BUCKET_OBJECT_DOMAIN"), postId),
    )
    if err != nil {
        return err
    }
    err = putUserPostToOutlineTable(
        db,
        putItem.UserId,
        fmt.Sprintf("%s/small/%s.jpg", os.Getenv("S3_BUCKET_OBJECT_DOMAIN"), postId),
        postId,
    )
    if err != nil {
        return err
    }

    err = putImage(putItem.Base64Image, postId)
    return err
}

func putImage(base64Image string, postId string) error {
    imageBinary, err := base64.StdEncoding.DecodeString(base64Image)
    if err != nil {
        return err
    }
    largeImageBuf := bytes.NewBuffer(imageBinary)
    smallImageBuf := *largeImageBuf
    largeImage, err := resize.Resize(largeImageBuf, 1000)
    if err != nil {
        return err
    }
    err = putImageToBucket("large/" + postId + ".jpg", largeImage)
    if err != nil {
        return err
    }
    smallImage, err := resize.Resize(&smallImageBuf, 300)
    if err != nil {
        return err
    }
    err = putImageToBucket("small/" + postId + ".jpg", smallImage)
    return err
}

func putImageToBucket(key string, image bytes.Buffer) error {
    client := s3.New(session.New(), &aws.Config{
        Region: aws.String(os.Getenv("S3_BUCKET_REGION")),
    })
    _, err := client.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
        Key: aws.String(key),
        Body: bytes.NewReader(image.Bytes()),
        ContentType: aws.String("image/jpeg"),
    })
    return err
}

type UserPostOutline struct {
    UserId string
    PostedTime string
    SmallImageUrl string
    // PostId string
}

func putUserPostToOutlineTable(db *dynamo.DB, userId string, smallImageUrl string, postId string) error {
    table := db.Table(os.Getenv("USER_POST_OUTLINE_TABLE_NAME"))
    e := UserPostOutline{
        UserId: userId,
        PostedTime: time.Now().UTC().Format("2006-01-02-15-04-05-0700"),
        SmallImageUrl: smallImageUrl,
        // PostId: postId,
    }
    return table.Put(e).If("attribute_not_exists(PostedTime)").Run()
}

type UserPostDetail struct {
    PostId string `json:"postId"`
    EatingPlace EatingPlace `json:"eatingPlace"`
    LargeImageUrl string `json:"largeImageUrl"`
}

func putUserPostToDetailTable(db *dynamo.DB, postId string, eatingPlace EatingPlace, largeImageUrl string) error {
    table := db.Table(os.Getenv("USER_POST_DETAIL_TABLE_NAME"))
    item := UserPostDetail{
        PostId: postId,
        EatingPlace: eatingPlace,
        LargeImageUrl: largeImageUrl,
    }
    return table.Put(item).If("attribute_not_exists(PostId)").Run()
}
