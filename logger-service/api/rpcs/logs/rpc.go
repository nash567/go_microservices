package logs

import (
	"context"

	"github.com/loger-service/api/proto/logs"
	"github.com/loger-service/internal/service/log/model"
	"github.com/loger-service/internal/service/log/repo"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	repo repo.Repo
}

func NewServer(repo repo.Repo) *LogServer {
	return &LogServer{repo: repo}
}
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	//write the log
	logEntry := model.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.repo.Insert(&logEntry)

	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	res := &logs.LogResponse{Result: "Logged!"}

	return res, nil
}
