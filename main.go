package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"math/big"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"golang.org/x/image/draw"
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
    imageBuf := bytes.NewBuffer(imageBinary)

    err = resize(imageBuf, 1000)
    err = putImageToBucket("large/" + postId + ".jpg", *imageBuf)

    err = resize(imageBuf, 300)
    err = putImageToBucket("small/" + postId + ".jpg", *imageBuf)

    // db := dynamo.New(session.New(), &aws.Config{
    //     Region: aws.String(os.Getenv("DB_REGION")),
    // })
    // err = putUserPostToOutlineTable(db, putItem.UserId, putItem.ImageUrl, postId)
    return
}

func resize(imageBuf *bytes.Buffer, maxLength int) (err error) {
    decordedImage, _, err := image.Decode(imageBuf)
    if err != nil {
        return
    }
    rectangle := decordedImage.Bounds()
    width := rectangle.Dx()
    height := rectangle.Dy()
    targetWidth, targetHeight := targetSize(width, height, maxLength)
    dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
    draw.CatmullRom.Scale(
        dst,
        dst.Bounds(),
        decordedImage,
        rectangle,
        draw.Over,
        nil,
    )
    err = jpeg.Encode(imageBuf, dst, &jpeg.Options{Quality: 100})
    return
}

func targetSize(width int, height int, maxLength int) (targetWidth int, targetHeight int) {
    gcd := int(new(big.Int).GCD(nil, nil, big.NewInt(int64(width)), big.NewInt(int64(height))).Int64())
    widthRate := width / gcd
    heightRate := height / gcd
    if widthRate > heightRate {
        rate := maxLength / widthRate
        targetWidth = widthRate * rate
        targetHeight = heightRate * rate
        return
    }
    rate := maxLength / heightRate
    targetWidth = widthRate * rate
    targetHeight = heightRate * rate
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
