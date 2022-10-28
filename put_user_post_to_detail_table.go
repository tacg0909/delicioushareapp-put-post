package main

import (
	"os"

	"github.com/guregu/dynamo"
)

type Contributor struct {
    UserId string `json:"userId"`
}

type UserPostDetail struct {
    PostId string `json:"postId"`
    EatingPlace EatingPlace `json:"eatingPlace"`
    LargeImageUrl string `json:"largeImageUrl"`
    Contributor Contributor `json:"contributor"`
    PostedTime string `json:"postedTime"`
}

func putUserPostToDetailTable(db *dynamo.DB, postId string, eatingPlace EatingPlace, largeImageUrl string, userId string, postedTime string) error {
    table := db.Table(os.Getenv("USER_POST_DETAIL_TABLE_NAME"))
    item := UserPostDetail{
        PostId: postId,
        EatingPlace: eatingPlace,
        LargeImageUrl: largeImageUrl,
        Contributor: Contributor {
            UserId: userId,
        },
        PostedTime: postedTime,
    }
    return table.Put(item).If("attribute_not_exists(PostId)").Run()
}
