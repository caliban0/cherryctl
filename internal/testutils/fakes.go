package testutils

import (
	"errors"
	"github.com/cherryservers/cherryctl/internal/projects"
	"github.com/cherryservers/cherryctl/internal/servers"
	"github.com/cherryservers/cherryctl/internal/sshkeys"
	"github.com/cherryservers/cherryctl/internal/storages"
	"github.com/cherryservers/cherryctl/internal/teams"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

type FakeSSHKeysService struct {
	GetResp    func(sshKeyID int, opts *cherrygo.GetOptions) (cherrygo.SSHKey, *cherrygo.Response, error)
	CreateResp func(request *cherrygo.CreateSSHKey) (cherrygo.SSHKey, *cherrygo.Response, error)
	DeleteResp func(sshKeyID int) (cherrygo.SSHKey, *cherrygo.Response, error)
	ListResp   func(opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error)
	UpdateResp func(sshKeyID int, request *cherrygo.UpdateSSHKey) (cherrygo.SSHKey, *cherrygo.Response, error)
}

func (f FakeSSHKeysService) List(opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error) {
	if f.ListResp == nil {
		return []cherrygo.SSHKey{}, nil, errors.New("list called without injected ListResp mock")
	}
	return f.ListResp(opts)
}

func (f FakeSSHKeysService) Get(sshKeyID int, opts *cherrygo.GetOptions) (cherrygo.SSHKey, *cherrygo.Response, error) {
	if f.GetResp == nil {
		return cherrygo.SSHKey{}, nil, errors.New("get called without injected GetResp mock")
	}
	return f.GetResp(sshKeyID, opts)
}

func (f FakeSSHKeysService) Create(request *cherrygo.CreateSSHKey) (cherrygo.SSHKey, *cherrygo.Response, error) {
	if f.CreateResp == nil {
		return cherrygo.SSHKey{}, nil, errors.New("create called without injected CreateResp mock")
	}
	return f.CreateResp(request)
}

func (f FakeSSHKeysService) Delete(sshKeyID int) (cherrygo.SSHKey, *cherrygo.Response, error) {
	if f.DeleteResp == nil {
		return cherrygo.SSHKey{}, nil, errors.New("delete called without injected DeleteResp mock")
	}
	return f.DeleteResp(sshKeyID)
}

func (f FakeSSHKeysService) Update(sshKeyID int, request *cherrygo.UpdateSSHKey) (cherrygo.SSHKey, *cherrygo.Response, error) {
	if f.UpdateResp == nil {
		return cherrygo.SSHKey{}, nil, errors.New("update called without injected UpdateResp mock")
	}
	return f.UpdateResp(sshKeyID, request)
}

type FakeProjectService struct {
	ListSSHKeysResp func(projectID int, opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error)
}

func (f FakeProjectService) List(teamID int, opts *cherrygo.GetOptions) ([]cherrygo.Project, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeProjectService) Get(projectID int, opts *cherrygo.GetOptions) (cherrygo.Project, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeProjectService) Create(teamID int, request *cherrygo.CreateProject) (cherrygo.Project, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeProjectService) Update(projectID int, request *cherrygo.UpdateProject) (cherrygo.Project, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeProjectService) ListSSHKeys(projectID int, opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error) {
	if f.ListSSHKeysResp == nil {
		return []cherrygo.SSHKey{}, nil, errors.New("list ssh keys called without injected ListSSHKeysResp mock")
	}
	return f.ListSSHKeysResp(projectID, opts)
}

func (f FakeProjectService) Delete(projectID int) (*cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

type FakeStorageService struct {
	AttachResp func(request *cherrygo.AttachTo) (cherrygo.BlockStorage, *cherrygo.Response, error)
	CreateResp func(request *cherrygo.CreateStorage) (cherrygo.BlockStorage, *cherrygo.Response, error)
	DeleteResp func(storageID int) (*cherrygo.Response, error)
	DetachResp func(storageID int) (*cherrygo.Response, error)
	GetResp    func(storageID int, opts *cherrygo.GetOptions) (cherrygo.BlockStorage, *cherrygo.Response, error)
	ListResp   func(projectID int, opts *cherrygo.GetOptions) ([]cherrygo.BlockStorage, *cherrygo.Response, error)
	UpdateResp func(request *cherrygo.UpdateStorage) (cherrygo.BlockStorage, *cherrygo.Response, error)
}

func (f FakeStorageService) List(projectID int, opts *cherrygo.GetOptions) ([]cherrygo.BlockStorage, *cherrygo.Response, error) {
	if f.ListResp == nil {
		return []cherrygo.BlockStorage{}, nil, errors.New("list storages called without injected ListResp mock")
	}
	return f.ListResp(projectID, opts)
}

func (f FakeStorageService) Get(storageID int, opts *cherrygo.GetOptions) (cherrygo.BlockStorage, *cherrygo.Response, error) {
	if f.GetResp == nil {
		return cherrygo.BlockStorage{}, nil, errors.New("get storage called without injected GetResp mock")
	}
	return f.GetResp(storageID, opts)
}

func (f FakeStorageService) Create(request *cherrygo.CreateStorage) (cherrygo.BlockStorage, *cherrygo.Response, error) {
	if f.CreateResp == nil {
		return cherrygo.BlockStorage{}, nil, errors.New("create storage called without injected CreateResp mock")
	}
	return f.CreateResp(request)
}

func (f FakeStorageService) Delete(storageID int) (*cherrygo.Response, error) {
	if f.DeleteResp == nil {
		return nil, errors.New("delete storage called without injected DeleteResp mock")
	}
	return f.DeleteResp(storageID)
}

func (f FakeStorageService) Attach(request *cherrygo.AttachTo) (cherrygo.BlockStorage, *cherrygo.Response, error) {
	if f.AttachResp == nil {
		return cherrygo.BlockStorage{}, nil, errors.New("attach storage called without injected AttachResp mock")
	}
	return f.AttachResp(request)
}

func (f FakeStorageService) Detach(storageID int) (*cherrygo.Response, error) {
	if f.DetachResp == nil {
		return nil, errors.New("detach storage called without injected DetachResp mock")
	}
	return f.DetachResp(storageID)
}

func (f FakeStorageService) Update(request *cherrygo.UpdateStorage) (cherrygo.BlockStorage, *cherrygo.Response, error) {
	if f.UpdateResp == nil {
		return cherrygo.BlockStorage{}, nil, errors.New("update storage called without injected UpdateResp mock")
	}
	return f.UpdateResp(request)
}

type FakeServersService struct {
	ListResp   func(projectID int, opts *cherrygo.GetOptions) ([]cherrygo.Server, *cherrygo.Response, error)
	CreateResp func(request *cherrygo.CreateServer) (cherrygo.Server, *cherrygo.Response, error)
	DeleteResp func(serverID int) (cherrygo.Server, *cherrygo.Response, error)
}

func (f FakeServersService) List(projectID int, opts *cherrygo.GetOptions) ([]cherrygo.Server, *cherrygo.Response, error) {
	if f.ListResp == nil {
		return []cherrygo.Server{}, nil, errors.New("list servers called without injected ListResp mock")
	}
	return f.ListResp(projectID, opts)
}

func (f FakeServersService) Get(serverID int, opts *cherrygo.GetOptions) (cherrygo.Server, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) PowerOff(serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) PowerOn(serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) Create(request *cherrygo.CreateServer) (cherrygo.Server, *cherrygo.Response, error) {
	if f.CreateResp == nil {
		return cherrygo.Server{}, nil, errors.New("create server called without injected CreateResp mock")
	}
	return f.CreateResp(request)
}

func (f FakeServersService) Delete(serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	if f.DeleteResp == nil {
		return cherrygo.Server{}, nil, errors.New("delete server called without injected DeleteResp mock")
	}
	return f.DeleteResp(serverID)
}

func (f FakeServersService) PowerState(serverID int) (cherrygo.PowerState, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) Reboot(serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) EnterRescueMode(serverID int, fields *cherrygo.RescueServerFields) (cherrygo.Server, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) ExitRescueMode(serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) Update(serverID int, request *cherrygo.UpdateServer) (cherrygo.Server, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) Reinstall(serverID int, fields *cherrygo.ReinstallServerFields) (cherrygo.Server, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) ListSSHKeys(serverID int, opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeServersService) ResetBMCPassword(serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	//TODO implement me
	panic("implement me")
}

type FakeTeamsService struct {
	GetResp    func(teamID int, opts *cherrygo.GetOptions) (cherrygo.Team, *cherrygo.Response, error)
	CreateResp func(request *cherrygo.CreateTeam) (cherrygo.Team, *cherrygo.Response, error)
	ListResp   func(opts *cherrygo.GetOptions) ([]cherrygo.Team, *cherrygo.Response, error)
	UpdateResp func(teamID int, request *cherrygo.UpdateTeam) (cherrygo.Team, *cherrygo.Response, error)
	DeleteResp func(teamID int) (*cherrygo.Response, error)
}

func (f FakeTeamsService) List(opts *cherrygo.GetOptions) ([]cherrygo.Team, *cherrygo.Response, error) {
	if f.ListResp == nil {
		return []cherrygo.Team{}, nil, errors.New("list team called without injected ListResp mock")
	}
	return f.ListResp(opts)
}

func (f FakeTeamsService) Get(teamID int, opts *cherrygo.GetOptions) (cherrygo.Team, *cherrygo.Response, error) {
	if f.GetResp == nil {
		return cherrygo.Team{}, nil, errors.New("get team called without injected GetResp mock")
	}
	return f.GetResp(teamID, opts)
}

func (f FakeTeamsService) Create(request *cherrygo.CreateTeam) (cherrygo.Team, *cherrygo.Response, error) {
	if f.CreateResp == nil {
		return cherrygo.Team{}, nil, errors.New("create team called without injected CreateResp mock")
	}
	return f.CreateResp(request)
}

func (f FakeTeamsService) Update(teamID int, request *cherrygo.UpdateTeam) (cherrygo.Team, *cherrygo.Response, error) {
	if f.UpdateResp == nil {
		return cherrygo.Team{}, nil, errors.New("update team called without injected UpdateResp mock")
	}
	return f.UpdateResp(teamID, request)
}

func (f FakeTeamsService) Delete(teamID int) (*cherrygo.Response, error) {
	if f.DeleteResp == nil {
		return nil, errors.New("delete team called without injected DeleteResp mock")
	}
	return f.DeleteResp(teamID)
}

// FakeServicer is a fake implementation of the root command`s client. It uses fake services,
// that can be injected with mocked responses.
// Note: names like `service` and `servicer` might be confusing here, since they don't represent web services.
// They were chosen to align with `cherrygo` terminology.
type FakeServicer struct {
	FakeSSHKeysService *FakeSSHKeysService
	FakeProjectService *FakeProjectService
	FakeStorageService *FakeStorageService
	FakeServersService *FakeServersService
	FakeTeamsService   *FakeTeamsService
}

func (f *FakeServicer) API(command *cobra.Command) *cherrygo.Client {
	return &cherrygo.Client{SSHKeys: f.FakeSSHKeysService,
		Projects: f.FakeProjectService,
		Storages: f.FakeStorageService,
		Servers:  f.FakeServersService,
		Teams:    f.FakeTeamsService,
	}
}

func (f *FakeServicer) GetOptions() *cherrygo.GetOptions {
	return &cherrygo.GetOptions{}
}

func (f *FakeServicer) Config(cmd *cobra.Command) *viper.Viper {
	//TODO implement me
	panic("implement me")
}

type ExperimentalFixture struct {
	// RootCmd is the main command on which the tests should be executed.
	RootCmd  *cobra.Command
	outputer *SpyOutputer
	t        *testing.T
}

// NewExperimentalFixture sets up a testing fixture. Use `setupFunc` to inject mocks for resource client responses.
func NewExperimentalFixture(t *testing.T, setupFunc func(fakeServicer *FakeServicer)) ExperimentalFixture {
	// Set up fake clients.
	fakeSSHKeysClient := FakeSSHKeysService{}
	fakeProjectsClient := FakeProjectService{}
	fakeStorageClient := FakeStorageService{}
	fakeServersClient := FakeServersService{}
	fakeTeamsClient := FakeTeamsService{}

	// Set up a fake servicer.
	fakeServicer := FakeServicer{&fakeSSHKeysClient,
		&fakeProjectsClient,
		&fakeStorageClient,
		&fakeServersClient,
		&fakeTeamsClient}
	setupFunc(&fakeServicer)

	rootCmd := cobra.Command{}

	// Don't pollute test output.
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	outputer := SpyOutputer{}

	// Set up root command and add the necessary child commands.
	rootCmd.AddCommand(sshkeys.NewClient(&fakeServicer, &outputer).NewCommand(),
		projects.NewClient(&fakeServicer, &outputer).NewCommand(),
		storages.NewClient(&fakeServicer, &outputer).NewCommand(),
		servers.NewClient(&fakeServicer, &outputer).NewCommand(),
		teams.NewClient(&fakeServicer, &outputer).NewCommand())

	return ExperimentalFixture{
		RootCmd:  &rootCmd,
		outputer: &outputer,
		t:        t,
	}
}

// CheckOutput compares the actual output for the full resource and its tabled versions
// with what was expected and logs an error if there's a mismatch.
func (f ExperimentalFixture) CheckOutput(wantR interface{}, wantTH []string, wantTD [][]string) {
	if !reflect.DeepEqual(wantR, f.outputer.Resource) {
		f.t.Errorf("want resource: %v, got resource: %v", wantR, f.outputer.Resource)
	}
	if !reflect.DeepEqual(wantTH, f.outputer.TableHeader) {
		f.t.Errorf("want table header: %v, got table header: %v", wantTH, f.outputer.TableHeader)
	}
	if !reflect.DeepEqual(wantTD, f.outputer.TableData) {
		f.t.Errorf("want table data: %v, got table data: %v", wantTD, f.outputer.TableData)
	}
}

// CheckOutputLine compares the actual output line with what was expected and logs an error if there's a mismatch.
func (f ExperimentalFixture) CheckOutputLine(wantLine string) {
	if !reflect.DeepEqual(wantLine, f.outputer.Line) {
		f.t.Errorf("want: %v, got: %v", wantLine, f.outputer.Line)
	}
}
