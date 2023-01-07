package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/authentication-service/internal/service/auth/repo"
	"github.com/authentication-service/internal/utils"
)

type Service struct {
	repo repo.Repo
}

func NewService(repo repo.Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authenticating i m called")
	var reqPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := utils.ReadJSON(w, r, &reqPayload)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := s.repo.GetByEmail(reqPayload.Email)
	if err != nil {
		utils.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := s.repo.PasswordMatches(reqPayload.Password, *user)
	if err != nil || !valid {
		utils.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}
	// log authenticate

	err = s.LogRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	payload := utils.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	utils.WriteJSON(w, http.StatusAccepted, payload)
}

func (s *Service) LogRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
