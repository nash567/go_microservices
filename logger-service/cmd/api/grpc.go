package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/loger-service/data"
	"github.com/loger-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	//write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(&logEntry)

	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	res := &logs.LogResponse{Result: "Logged!"}

	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen for grpc:%v", err)
	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("Starting grpc server on port: %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for grpc:%v", err)

	}
}
