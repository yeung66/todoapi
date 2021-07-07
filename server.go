package main

import (
	"github.com/yeung66/todoapi/internal/auth"
	database "github.com/yeung66/todoapi/internal/db"
	"github.com/yeung66/todoapi/internal/todos"
	"github.com/yeung66/todoapi/internal/users"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/yeung66/todoapi/graph"
	"github.com/yeung66/todoapi/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	err := database.Init()
	if err != nil {
		panic("failed to connect database")
	}

	database.Db.AutoMigrate(&todos.TodoItem{}, &users.User{})

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
