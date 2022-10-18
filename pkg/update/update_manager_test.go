package update

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ManagerSuite struct {
	suite.Suite

	Buffer            *bytes.Buffer
	GitHubServiceMock *mocks.GitHubServiceInterface
}

func (suite *ManagerSuite) SetupSuite() {
	suite.Buffer = bytes.NewBufferString("")
	suite.GitHubServiceMock = mocks.NewGitHubServiceInterface(suite.T())
}

func (suite *ManagerSuite) TearDownTest() {
	suite.Buffer.Reset()
}

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func (suite *ManagerSuite) TestCheckUpdate_IsDevVersion() {
	version := types.Version{
		Number: "dev",
	}
	manager := NewManager(version, suite.GitHubServiceMock, suite.Buffer)

	manager.CheckUpdate()

	assert.Empty(suite.T(), suite.Buffer.String())
}

func (suite *ManagerSuite) TestCheckUpdate_IsUpToDate() {
	fixture := &dto.Release{
		TagName: "v1.0.5",
	}
	suite.GitHubServiceMock.On("GetLatestRelease", "ABGEO", "mailtm").Return(fixture, nil).Times(1)

	version := types.Version{
		Number: "1.0.5",
	}
	manager := NewManager(version, suite.GitHubServiceMock, suite.Buffer)

	manager.CheckUpdate()

	assert.Empty(suite.T(), suite.Buffer.String())
}

func (suite *ManagerSuite) TestCheckUpdate_IsOutdated() {
	fixture := &dto.Release{
		HTMLURL:     "https://foo.bar/v1.1.0",
		TagName:     "v1.1.0",
		PublishedAt: time.Now(),
	}
	suite.GitHubServiceMock.On("GetLatestRelease", "ABGEO", "mailtm").Return(fixture, nil).Times(1)

	version := types.Version{
		Number: "1.0.5",
	}
	manager := NewManager(version, suite.GitHubServiceMock, suite.Buffer)

	manager.CheckUpdate()

	assert.Contains(suite.T(), suite.Buffer.String(), "New version is available!")
	assert.Contains(
		suite.T(),
		suite.Buffer.String(),
		fmt.Sprintf("Version 1.1.0 was released on %s", fixture.PublishedAt.Format("02 January 2006")),
	)
	assert.Contains(suite.T(), suite.Buffer.String(), "You can download it from https://foo.bar/v1.1.0")
}
