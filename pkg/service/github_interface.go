package service

import dto "github.com/abgeo/mailtm/pkg/dto/github"

type GitHubServiceInterface interface {
	GetLatestRelease(owner, repo string) (*dto.Release, error)
}
