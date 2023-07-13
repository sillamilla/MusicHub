package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/sillamilla/user_microservice/handler"
	"github.com/sillamilla/user_microservice/internal/config"
	"github.com/sillamilla/user_microservice/internal/users/Mongo_storage"
	"github.com/sillamilla/user_microservice/internal/users/Redis_storage"
	"github.com/sillamilla/user_microservice/internal/users/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()

	//REDIS
	dbRedis := redis.NewClient(&redis.Options{
		Network:  cfg.Redis.Network,
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
	})

	defer func(dbRedis *redis.Client) {
		err := dbRedis.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dbRedis)

	err := dbRedis.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("Connect error Redis:", err)
	}

	//MONGO
	dbMongo, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Mongo.Address))
	if err != nil {
		log.Fatal(err)
	}
	defer func(dbMongo *mongo.Client, ctx context.Context) {
		err = dbMongo.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(dbMongo, context.Background())

	err = dbMongo.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Connect error Mongo:", err)
	}

	re := Redis_storage.New(dbRedis)
	mo := Mongo_storage.New(dbMongo)
	s := service.New(re, mo)
	h := handler.NewHandler(s)

	router := mux.NewRouter()
	router.HandleFunc("/signup", h.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/signin", h.SignIn).Methods(http.MethodPost)
	router.HandleFunc("/logout", h.Logout).Methods(http.MethodPost)
	router.HandleFunc("/editprofile", h.EditProfile).Methods(http.MethodPut)
	router.HandleFunc("/editpassword", h.ChangePassword).Methods(http.MethodPut)
	router.HandleFunc("/searchbyusername", h.SearchByUsername).Methods(http.MethodGet)

	router.HandleFunc("/setsession", h.UpsertSessions).Methods(http.MethodPost)
	router.HandleFunc("/getbyusername", h.GetByUsername).Methods(http.MethodGet)
	router.HandleFunc("/getbyid", h.GetById).Methods(http.MethodGet)
	router.HandleFunc("/getsession", h.GetSession).Methods(http.MethodGet)
	router.HandleFunc("/getbysession", h.GetBySession).Methods(http.MethodGet)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
