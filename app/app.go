package app

import (
	"github.com/bayudha2/go-test-0/config"
	"github.com/bayudha2/go-test-0/controllers/authcontroller"
	"github.com/bayudha2/go-test-0/controllers/postcontroller"
	"github.com/bayudha2/go-test-0/controllers/productcontroller"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var R *mux.Router

func Initialize() {
	R = mux.NewRouter().StrictSlash(true)
	InitializeRoutes(R)
}

func InitializeRoutes(R *mux.Router) {
	R.HandleFunc("/signup", authcontroller.Register).Methods("POST")
	R.HandleFunc("/signin", authcontroller.Login).Methods("POST")
	R.HandleFunc("/refresh", authcontroller.Refresh).Methods("POST")

	secure := R.PathPrefix("/v1").Subrouter()
	secure.Use(config.IsAuthorized)
	secure.HandleFunc("/signout", authcontroller.Logout).Methods("POST")

	secure.HandleFunc("/products", productcontroller.GetProducts).Methods("GET")
	secure.HandleFunc("/product", productcontroller.CreateProduct).Methods("POST")
	secure.HandleFunc("/product/{id}", productcontroller.GetProduct).Methods("GET")
	secure.HandleFunc("/product/{id}", productcontroller.UpdateProduct).Methods("PUT")
	secure.HandleFunc("/product/{id}", productcontroller.DeleteProduct).Methods("DELETE")

	secure.HandleFunc("/posts", postcontroller.GetPosts).Methods("GET")
	secure.HandleFunc("/post", postcontroller.CreatePost).Methods("POST")
	secure.HandleFunc("/post/{id}", postcontroller.GetPost).Methods("GET")
	secure.HandleFunc("/post/{id}", postcontroller.UpdatePost).Methods("PUT")
	secure.HandleFunc("/post/{id}", postcontroller.DeletePost).Methods("DELETE")
}
