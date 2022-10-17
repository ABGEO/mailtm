package service

import "github.com/abgeo/mailtm/pkg/dto"

type GitHubServiceInterface interface {
	GetLatestRelease(owner, repo string) (*dto.Release, error)
}
