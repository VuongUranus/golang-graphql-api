package main

import (
	"golang-graphql-api/database"
	"golang-graphql-api/graph"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

const defaultPort = "8080"

func init(){
	database.Connect("mongodb://localhost:27017/")
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", c.Handler(srv))

	log.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Panic(http.ListenAndServe(":"+port, nil))
}
