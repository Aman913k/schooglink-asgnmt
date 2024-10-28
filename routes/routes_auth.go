package routes

import (
	"net/http"

	controller "github.com/Aman913k/controllers"
	"github.com/Aman913k/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/register", controller.Register).Methods("POST")
	router.HandleFunc("/login", controller.Login).Methods("POST")
	router.Handle("/profile/view", middleware.JWTAuth(http.HandlerFunc(controller.ViewProfile))).Methods("GET")
	router.Handle("/posts/create", middleware.JWTAuth(http.HandlerFunc(controller.CreatePost))).Methods("POST")
	router.Handle("/post/delete", middleware.JWTAuth(http.HandlerFunc(controller.DeletePost))).Methods("DELETE")
	router.HandleFunc("/posts", controller.GetAllPosts).Methods("GET")
	router.HandleFunc("/posts/{post_id}", controller.GetPostByID).Methods("GET")
	router.Handle("/profile/{id}", middleware.JWTAuth(http.HandlerFunc(controller.UpdateProfile))).Methods("PUT")

	router.HandleFunc("/posts/{post_id}", func(w http.ResponseWriter, r *http.Request) {
		middleware.JWTAuth(http.HandlerFunc(controller.UpdatePost)).ServeHTTP(w, r)
	}).Methods("PUT")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return router
}
