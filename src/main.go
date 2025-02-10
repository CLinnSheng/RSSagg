package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CLinnSheng/RSSagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq" // use _ to tell compiler is ok not using this library
)

// In go, every variable & library that imported & declared must be use
// or else the compiler will give error

type apiConfig struct {
	DB *database.Queries
}

func main() {

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT IS NOT FOUND IN DONTENV")
	}

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	connection, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	defer connection.Close()

	db := database.New(connection)
	apiCfg := apiConfig{
		DB: db,
	}

	go startScrapping(db, 10, time.Minute)
	// router object
	router := chi.NewRouter()

	// options to give the client how can they access the server
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300}))

	v1Router := chi.NewRouter()

	// only respond to get request
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handleErr)

	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feed", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	router.Mount("/v1", v1Router)

	// server object
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server is running on port %s", portString)
	log.Fatal(server.ListenAndServe())

	fmt.Println(portString)
}
