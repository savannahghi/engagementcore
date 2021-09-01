package presentation

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/savannahghi/engagement/pkg/engagement/presentation/graph"
	"github.com/savannahghi/engagement/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/engagement/pkg/engagement/usecases"

	"github.com/savannahghi/serverutils"

	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	mbBytes              = 1048576
	serverTimeoutSeconds = 120
)

// AllowedOrigins is list of CORS origins allowed to interact with
// this service
var AllowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.healthcloud.co.ke",
	"https://a.bewell.co.ke",
	"https://b.bewell.co.ke",
	"http://localhost:5000",
	"https://europe-west3-bewell-app.cloudfunctions.net",
}
var allowedHeaders = []string{
	"Authorization", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type", " X-Authorization", " Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers",
}

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {
	r := mux.NewRouter() // gorilla mux
	SharedUnauthenticatedRoutes(ctx, r)
	AuthenticatedGraphQLRoute(ctx, r)
	SharedAuthenticatedISCRoutes(ctx, r)
	// return the combined router
	return r, nil
}

// PrepareServer starts up a server
func PrepareServer(
	ctx context.Context,
	port int,
	allowedOrigins []string,
) *http.Server {
	// start up the router
	r, err := Router(ctx)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	// start the server
	addr := fmt.Sprintf(":%d", port)
	h := handlers.CompressHandlerLevel(r, gzip.BestCompression)

	h = handlers.CORS(
		handlers.AllowedHeaders(allowedHeaders),
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST"}),
	)(h)
	h = handlers.CombinedLoggingHandler(os.Stdout, h)
	h = handlers.ContentTypeHandler(
		h,
		"application/json",
		"application/x-www-form-urlencoded",
	)
	srv := &http.Server{
		Handler:      h,
		Addr:         addr,
		WriteTimeout: serverTimeoutSeconds * time.Second,
		ReadTimeout:  serverTimeoutSeconds * time.Second,
	}
	log.Infof("Server running at port %v", addr)
	return srv
}

//HealthStatusCheck endpoint to check if the server is working.
func HealthStatusCheck(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		log.Fatal(err)
	}
}

// GQLHandler sets up a GraphQL resolver
func GQLHandler(ctx context.Context,
	usecases usecases.Usecases,
) http.HandlerFunc {
	resolver, err := graph.NewResolver(ctx, usecases)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}
