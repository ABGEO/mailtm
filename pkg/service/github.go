package service

import (
	"encoding/json"
	"time"

	dto "github.com/abgeo/mailtm/pkg/dto/github"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/go-resty/resty/v2"
)

type GitHubService struct {
	client *resty.Client
}

func NewGitHubService() *GitHubService {
	const timeout = 30 * time.Second

	client := resty.New()
	client.SetBaseURL("https://api.github.com/").
		SetTimeout(timeout).
		SetHeader("Accept", "application/vnd.github+json")

	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal

	return &GitHubService{
		client: client,
	}
}

func (svc *GitHubService) GetLatestRelease(owner, repo string) (release *dto.Release, err error) {
	_, err = svc.client.R().
		SetPathParams(types.StrMap{"owner": owner, "repo": repo}).
		SetResult(&release).
		Get("/repos/{owner}/{repo}/releases/latest")
	if err != nil {
		return release, err
	}

	return release, nil
}
