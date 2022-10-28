package main

import (
	"os"

	"github.com/guregu/dynamo"
)

type UserPostOutline struct {
    UserId string
    PostedTime string
    SmallImageUrl string
    PostId string
}

func putUserPostToOutlineTable(db *dynamo.DB, userId string, smallImageUrl string, postId string, postedTime string) error {
    table := db.Table(os.Getenv("USER_POST_OUTLINE_TABLE_NAME"))
    e := UserPostOutline{
        UserId: userId,
        PostedTime: postedTime,
        SmallImageUrl: smallImageUrl,
        PostId: postId,
    }
    return table.Put(e).If("attribute_not_exists(PostedTime)").Run()
}
