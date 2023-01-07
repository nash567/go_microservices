package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/authentication-service/config"
	"github.com/authentication-service/internal/db"
	"github.com/authentication-service/internal/service/auth"
	"github.com/authentication-service/internal/service/auth/repo"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Application struct {
	DB         *sql.DB
	httpServer *http.Server
	cfg        *config.AppConfig
	svc        *auth.Service
}

func (a *Application) Init(ctx context.Context, configFiles string) {
	var err error
	a.cfg, err = config.LoadConfig(strings.Split(configFiles, ",")...)
	if err != nil {
		fmt.Println("failed to load config files")
	}
	a.DB, err = setupDB(a.cfg.Database)

	if err != nil {
		log.Fatal("error connecting to database", err)
	}
	a.svc = buildAuthService(a.DB)
	fmt.Println("url port is", a.cfg.URLPort)
	a.httpServer = &http.Server{Addr: fmt.Sprintf(":%v", a.cfg.URLPort), Handler: registerHTTPEndpoints(a)}

}

func (a *Application) Start() {

	go func() {
		log.Println("http server started at port", a.cfg.URLPort)
		fmt.Println("http server is", a.httpServer)
		if err := a.httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("Error running Api", err)
			} else {
				log.Fatal("Stopping API")

			}
		}
	}()
}

func (a *Application) Stop(ctx context.Context) {
	if a.httpServer != nil {
		log.Println("Shutting down the Server...")
		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			log.Println("Error shutting down server gracefully")
		}
	}
}

func registerHTTPEndpoints(a *Application) http.Handler {
	mux := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*", "http://*"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ey weasad"))
	}).Methods("GET")
	mux.HandleFunc("/authenticate", a.svc.Authenticate).Methods("POST")
	handler := c.Handler(mux)
	return handler
}

func buildAuthService(db *sql.DB) *auth.Service {
	repo := repo.NewRepository(db)
	svc := auth.NewService(*repo)
	return svc

}
func setupDB(cfg *db.Config) (*sql.DB, error) {
	conn, err := db.NewConnection(cfg)
	if err != nil {

		return nil, fmt.Errorf("failed to open connection: %w", err)
	}
	if err := db.Verify(conn); err != nil {
		return nil, fmt.Errorf("failed to verify database connection: %w", err)
	}

	fmt.Println("Connection to db succesfull..")

	return conn, nil
}
