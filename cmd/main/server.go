package main

import (
	"curiosity/graph"
	"curiosity/internal/pkg/content/repo"
	"curiosity/internal/pkg/content/usecase"
	"curiosity/internal/utils"
	"database/sql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout", "/tmp/curiosity_logs"},
		ErrorOutputPaths: []string{"stderr", "/tmp/curiosity_err_logs"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey:   "level",
			TimeKey:    "ts",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}

	logger, err := cfg.Build()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("failed to sync logger", zap.Error(err))
		}
	}()

	connUrl, err := utils.GetConnectUrl()
	if err != nil {
		log.Fatalf("bad get connUrl: %v", err)
	}
	db, err := sql.Open("postgres", connUrl)
	if err != nil {
		log.Fatalf("bad db open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("bad ping to DB: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	postRepo := repo.NewPostRepository(db)
	commentRepo := repo.NewCommentRepository(db)

	postUseCase := usecase.NewPostUseCase(postRepo)
	commentUseCase := usecase.NewCommentUseCase(commentRepo)

	resolver := graph.NewResolver(postUseCase, commentUseCase, logger)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
