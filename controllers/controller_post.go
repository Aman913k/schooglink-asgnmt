package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Aman913k/database"
	"github.com/Aman913k/middleware"
	"github.com/Aman913k/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



// CreatePost godoc
// @Summary Create a new post
// @Description Allows a logged-in user to create a new blog post. The user's email and name are retrieved from the request context.
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post content"
// @Success 201 {object} map[string]interface{} "Post created successfully"
// @Failure 400 {string} string "Invalid post data"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Failed to create post"
// @Router /posts/create [post]
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Retrieving user email and name from context
	email, ok := r.Context().Value(middleware.EmailContextKey).(string)
	if !ok || email == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	name, ok := r.Context().Value(middleware.NameContextKey).(string)
	if !ok || name == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid post data", http.StatusBadRequest)
		return
	}

	post.Email = email
	post.Author = name
	post.Created_AT = time.Now()

	collection := database.GetCollection("schoogpost")
	insertResult, err := collection.InsertOne(context.TODO(), post)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created successfully",
		"postID":  insertResult.InsertedID,
	})
}

// DeletePost deletes a blog post.
// @Summary Delete a post
// @Description Delete a post by post ID.
// @Tags Posts
// @Accept json
// @Produce json
// @Param post_id query string true "Post ID"
// @Success 200 {string} string "Post deleted successfully"
// @Failure 400 {string} string "Post ID is required"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Post not found"
// @Failure 500 {string} string "Failed to delete post"
// @Router /posts [delete]
func DeletePost(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(middleware.EmailContextKey).(string)

	if !ok || email == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := r.URL.Query().Get("post_id")

	if postIDStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("schoogpost")

	filter := bson.M{"_id": postID, "email": email}
	var post models.Post

	err = collection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Deleting post if found
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Post deleted successfully"}`))
}

// GetAllPosts retrieves all blog posts.
// @Summary Get all posts
// @Description Get all blog posts.
// @Tags Posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post "List of posts"
// @Failure 500 {string} string "Failed to fetch posts"
// @Router /posts [get]
func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection("schoogpost")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var posts []primitive.M
	for cursor.Next(context.TODO()) {
		var post bson.M
		if err := cursor.Decode(&post); err != nil {
			http.Error(w, "Failed to decode post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error while reading posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetPostByID retrieves a single blog post by ID.
// @Summary Get a post by ID
// @Description Get a blog post by post ID.
// @Tags Posts
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} models.Post "Post details"
// @Failure 400 {string} string "Invalid Post ID format"
// @Failure 404 {string} string "Post not found"
// @Failure 500 {string} string "Database error"
// @Router /posts/{post_id} [get]
func GetPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDStr := vars["post_id"]

	// Converting post ID to ObjectID
	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("schoogpost")
	var post models.Post
	err = collection.FindOne(context.TODO(), bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// UpdatePost updates a post.
// @Summary Update a post
// @Description Update the post details of a logged-in user.
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param user body models.Post true "Updated post data"
// @Success 200 {string} string "Post updated successfully"
// @Failure 400 {string} string "Invalid post data"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Post not found"
// @Failure 500 {string} string "Failed to update post"
// @Router /posts/{post_id} [put]
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(middleware.EmailContextKey).(string)
	if !ok || email == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	postIDStr := vars["post_id"]

	// Converting post ID to ObjectID
	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("schoogpost")
	var post models.Post
	err = collection.FindOne(context.TODO(), bson.M{"_id": postID, "email": email}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Post not found or unauthorized", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	var updatedPost models.Post
	err = json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, "Invalid post data", http.StatusBadRequest)
		return
	}

	post.Title = updatedPost.Title
	post.Content = updatedPost.Content
	post.Updated_AT = time.Now()

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": postID}, bson.M{"$set": post})
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post updated successfully",
		"post":    post,
	})
}
