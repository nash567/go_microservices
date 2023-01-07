package log

import (
	"fmt"
	"net/http"

	"github.com/loger-service/internal/service/log/model"
	"github.com/loger-service/internal/service/log/repo"
	"github.com/loger-service/internal/utils"
)

type Service struct {
	repo repo.Repo
}

func NewService(repo repo.Repo) *Service {
	return &Service{repo: repo}
}
func (s *Service) WriteLog(w http.ResponseWriter, r *http.Request) {
	var reqPayload model.JsonPayload

	_ = utils.ReadJSON(w, r, &reqPayload)

	fmt.Println("request payload is", reqPayload)
	event := &model.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}

	err := s.repo.Insert(event)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JsonResponse{
		Error:   false,
		Message: "logged",
	}
	fmt.Println("sending respone to request")
	utils.WriteJSON(w, http.StatusAccepted, resp)

}
