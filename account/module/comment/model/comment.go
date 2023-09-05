package commentmodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionComment = "comments"

// type JsonPostedDate time.Time

// // Implement Marshaler and Unmarshaler interface
// func (j *JsonPostedDate) UnmarshalJSON(b []byte) error {
// 	fmt.Println("UnmarshalJSON")
// 	s := strings.Trim(string(b), "\"")
// 	t, err := time.Parse("2006-01-02", s)
// 	if err != nil {
// 		return err
// 	}
// 	*j = JsonPostedDate(t)
// 	return nil
// }

// func (j JsonPostedDate) MarshalJSON() ([]byte, error) {
// 	fmt.Println("MarshalJSON")
// 	return json.Marshal(time.Time(j))
// }

// // Maybe a Format function for printing your date
// func (j JsonPostedDate) Format(s string) string {
// 	t := time.Time(j)
// 	return t.Format(s)
// }

type CommentCreate struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Text string             `json:"text" bson:"text"`
}

func (data *CommentCreate) Fullfill() {

}

func (CommentCreate) CollectionsName() string {
	return "comments"
}

// func (data *CommentCreate) ToBson() primitive.D {
// 	return bson.D{
// 		{Key: "author", Value: data.Author},
// 		{Key: "discuss_id", Value: data.DiscussId},
// 		{Key: "posted", Value: data.Posted},
// 		{Key: "text", Value: data.Text},
// 		{Key: "parent_slug", Value: data.ParentSlug},
// 		{Key: "score", Value: data.Score},
// 		{Key: "slug", Value: data.Slug},
// 		{Key: "comment_replies_num", Value: data.CommentRepliesNum},
// 		{Key: "comment_likes", Value: data.CommentLikes},
// 		{Key: "comment_like_num", Value: data.CommentLikeNum},
// 		{Key: "full_slug", Value: data.FullSlug},
// 	}
// }
