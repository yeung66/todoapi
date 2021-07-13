package main

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
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
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		err.Message = e.Error()

		var wrongUserErr *users.WrongUsernameOrPasswordError
		var permissionErr *users.PermissionDeniedError
		var noUserErr *users.HasNoTodoItemError
		if errors.As(e, &wrongUserErr) {
			err.Extensions = map[string]interface{}{
				"code": "UNAUTHENTICATED",
			}
		} else if errors.As(e, &permissionErr) {
			err.Extensions = map[string]interface{}{
				"code": "FORBIDDEN",
			}
		} else if errors.As(e, &noUserErr) {
			err.Extensions = map[string]interface{}{
				"code": "FORBIDDEN",
			}
		}

		return err
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
