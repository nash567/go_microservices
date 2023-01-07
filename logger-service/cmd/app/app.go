package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/loger-service/api/proto/logs"
	logRPC "github.com/loger-service/api/rpcs/logs"
	"github.com/loger-service/config"
	"github.com/loger-service/internal/db"
	logger "github.com/loger-service/internal/service/log"
	"github.com/loger-service/internal/service/log/repo"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type Application struct {
	DB         *mongo.Client
	httpServer *http.Server
	grpcServer *grpc.Server
	cfg        *config.AppConfig
	svc        *logger.Service
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
	a.svc = buildLogService(a.DB)
	fmt.Println("url port is", a.cfg.URLPort)
	a.httpServer = &http.Server{Addr: fmt.Sprintf(":%v", a.cfg.URLPort), Handler: registerHTTPEndpoints(a)}
	repo := repo.NewRepository(a.DB)
	a.grpcServer = registerGRPCEndpoints(*repo)
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

	// start grpc server
	go func() {
		listenOn := fmt.Sprintf("%v:%v", a.cfg.URLHost, a.cfg.GRPCPort)
		listener, listenerErr := net.Listen("tcp", listenOn)
		if listenerErr != nil {
			log.Fatalf("failed to listen on %s", listenOn)
		}

		log.Printf("grpc server started")
		if err := a.grpcServer.Serve(listener); err != nil {
			log.Fatal("failed to serve gRPC server")
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

func registerGRPCEndpoints(repo repo.Repo) *grpc.Server {
	grpcServer := grpc.NewServer()
	logs.RegisterLogServiceServer(grpcServer, logRPC.NewServer(repo))
	return grpcServer
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

	mux.HandleFunc("/log", a.svc.WriteLog).Methods("POST")
	handler := c.Handler(mux)
	return handler
}

func buildLogService(db *mongo.Client) *logger.Service {
	repo := repo.NewRepository(db)
	svc := logger.NewService(*repo)
	return svc

}
func setupDB(cfg *db.Config) (*mongo.Client, error) {
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
