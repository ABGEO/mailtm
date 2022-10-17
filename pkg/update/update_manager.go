package update

import (
	"io"

	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/hashicorp/go-version"
	"github.com/pterm/pterm"
)

type Manager struct {
	Version       types.Version
	GitHubService service.GitHubServiceInterface
	Writer        io.Writer
}

func NewManager(version types.Version, gitHubService service.GitHubServiceInterface, writer io.Writer) *Manager {
	return &Manager{
		Version:       version,
		GitHubService: gitHubService,
		Writer:        writer,
	}
}

func (manager *Manager) CheckUpdate() {
	const (
		repoOwner = "ABGEO"
		repoName  = "mailtm"
	)

	if manager.Version.Number == "dev" || manager.Version.Number == "test" {
		return
	}

	if release, err := manager.GitHubService.GetLatestRelease(repoOwner, repoName); err == nil {
		currentVersion, _ := version.NewVersion(manager.Version.Number)
		latestVersion, _ := version.NewVersion(release.TagName)

		if latestVersion.GreaterThan(currentVersion) {
			paragraph := pterm.DefaultParagraph.WithWriter(manager.Writer)

			paragraph.Println(pterm.Yellow("New version is available!"))
			paragraph.Printfln("Version %s was released on %s", latestVersion, release.PublishedAt.Format("02 January 2006"))
			paragraph.Printfln("You can download it from %s", release.HTMLURL)
		}
	}
}
