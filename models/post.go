package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post represents the blog post model for the application.
// @Description Post model
// @Name Post
// @Property id objectId `json:"id,omitempty" bson:"_id,omitempty"`
// @Property email string `json:"email,omitempty"`
// @Property title string `json:"title"`
// @Property content string `json:"content"`
// @Property author string `json:"author"`
// @Property created_at string `json:"created_at"` // Date when the post was created
// @Property updated_at string `json:"updated_at"` // Date when the post was last updated
type Post struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email      string             `json:"email,omitempty"`
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Author     string             `json:"author"`
	Created_AT time.Time          `json:"created_at"`
	Updated_AT time.Time          `json:"updated_at"`
}


// type UpdatePostResponse struct {
//     Message string `json:"message"`
//     User    User   `json:"user"`
// }