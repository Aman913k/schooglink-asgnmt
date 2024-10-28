package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Aman913k/database"
	"github.com/Aman913k/middleware"
	"github.com/Aman913k/models"
	"github.com/Aman913k/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// EncryptUserPassword hashes the user's password and inserts the user record into the MongoDB database.
func EncryptUserPassword(user *models.User) (*mongo.InsertOneResult, error) {
	hashedPassword, err := models.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	const dbName = "schooguser"
	collection := database.GetCollection(dbName)
	user.Password = hashedPassword
	return collection.InsertOne(context.TODO(), user)
}

// Register registers a new user.
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Email already in use"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	const colName = "schooguser"
	collection := database.GetCollection(colName)

	// Check if email is already in use
	err = collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&user)
	if err == nil {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	} else if err != mongo.ErrNoDocuments {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Validate email and password
	if !utils.IsValidGmail(user.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}
	if !utils.StrongPassword(user.Password) {
		http.Error(w, "Password is weak", http.StatusForbidden)
	}

	// Register user by hashing password and saving to the database
	result, err := EncryptUserPassword(&user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Login authenticates a user and generates a JWT token.
// @Summary Login user
// @Description Login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "Login details"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Invalid email or password"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	const colName = "schooguser"
	collection := database.GetCollection(colName)
	var foundUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(foundUser.Email, foundUser.Name)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// ViewProfile retrieves and displays the profile of the authenticated user.
// @Summary View user profile
// @Description View profile details of the logged-in user
// @Tags Profile
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Router /profile/view [get]
func ViewProfile(w http.ResponseWriter, r *http.Request) {
	userEmail, ok := r.Context().Value(middleware.EmailContextKey).(string)
	if !ok || userEmail == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	const colName = "schooguser"
	collection := database.GetCollection(colName)

	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}





// UpdateProfile updates a user's profile.
// @Summary Update user profile
// @Description Update the profile details of a logged-in user.
// @Tags Profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {string} string "Profile updated successfully"
// @Failure 400 {string} string "Invalid user data"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Failed to update profile"
// @Router /profile/{id} [put]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Retrieve email from context
	email, ok := r.Context().Value(middleware.EmailContextKey).(string)
	if !ok || email == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}


	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("schooguser")
	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": userID, "email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found or unauthorized", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	var updatedUser models.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	user.Name = updatedUser.Name
	// You can also update other fields like Email if needed, e.g.:
	// user.Email = updatedUser.Email

	// Update the user in the database
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": userID}, bson.M{"$set": user})
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Send the response back
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User updated successfully",
		"user":    user,
	})
}
