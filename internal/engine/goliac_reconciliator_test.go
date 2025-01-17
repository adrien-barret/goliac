package engine

import (
	"context"
	"fmt"
	"testing"

	"github.com/Alayacare/goliac/internal/config"
	"github.com/Alayacare/goliac/internal/entity"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gosimple/slug"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type GoliacLocalMock struct {
	users     map[string]*entity.User
	externals map[string]*entity.User
	teams     map[string]*entity.Team
	repos     map[string]*entity.Repository
	rulesets  map[string]*entity.RuleSet
}

func (m *GoliacLocalMock) Clone(accesstoken, repositoryUrl, branch string) error {
	return nil
}
func (m *GoliacLocalMock) ListCommitsFromTag(tagname string) ([]*object.Commit, error) {
	return nil, fmt.Errorf("not tag %s found", tagname)
}
func (m *GoliacLocalMock) GetHeadCommit() (*object.Commit, error) {
	return nil, nil
}
func (m *GoliacLocalMock) CheckoutCommit(commit *object.Commit) error {
	return nil
}
func (m *GoliacLocalMock) PushTag(tagname string, hash plumbing.Hash, accesstoken string) error {
	return nil
}
func (m *GoliacLocalMock) LoadRepoConfig() (error, *config.RepositoryConfig) {
	return nil, &config.RepositoryConfig{}
}
func (m *GoliacLocalMock) LoadAndValidate() ([]error, []entity.Warning) {
	return nil, nil
}
func (m *GoliacLocalMock) LoadAndValidateLocal(fs afero.Fs, path string) ([]error, []entity.Warning) {
	return nil, nil
}
func (m *GoliacLocalMock) Teams() map[string]*entity.Team {
	return m.teams
}
func (m *GoliacLocalMock) Repositories() map[string]*entity.Repository {
	return m.repos
}
func (m *GoliacLocalMock) Users() map[string]*entity.User {
	return m.users
}
func (m *GoliacLocalMock) ExternalUsers() map[string]*entity.User {
	return m.externals
}
func (m *GoliacLocalMock) RuleSets() map[string]*entity.RuleSet {
	return m.rulesets
}
func (m *GoliacLocalMock) UpdateAndCommitCodeOwners(repoconfig *config.RepositoryConfig, dryrun bool, accesstoken string, branch string, tagname string) error {
	return nil
}
func (m *GoliacLocalMock) SyncUsersAndTeams(repoconfig *config.RepositoryConfig, plugin UserSyncPlugin, dryrun bool) error {
	return nil
}
func (m *GoliacLocalMock) Close() {

}

type GoliacRemoteMock struct {
	users      map[string]string
	teams      map[string]*GithubTeam // key is the slug team
	repos      map[string]*GithubRepository
	teamsrepos map[string]map[string]*GithubTeamRepo // key is the slug team
	rulesets   map[string]*GithubRuleSet
	appids     map[string]int
}

func (m *GoliacRemoteMock) Load() error {
	return nil
}
func (m *GoliacRemoteMock) IsEnterprise() bool {
	return true
}
func (m *GoliacRemoteMock) FlushCache() {

}
func (m *GoliacRemoteMock) RuleSets() map[string]*GithubRuleSet {
	return m.rulesets
}
func (m *GoliacRemoteMock) Users() map[string]string {
	return m.users
}

func (m *GoliacRemoteMock) TeamSlugByName() map[string]string {
	slugs := make(map[string]string)
	for _, v := range m.teams {
		slugs[v.Name] = slug.Make(v.Name)
	}
	return slugs
}
func (m *GoliacRemoteMock) Teams() map[string]*GithubTeam {
	return m.teams
}
func (m *GoliacRemoteMock) Repositories() map[string]*GithubRepository {
	return m.repos
}
func (m *GoliacRemoteMock) RepositoriesByRefId() map[string]*GithubRepository {
	return make(map[string]*GithubRepository)
}
func (m *GoliacRemoteMock) TeamRepositories() map[string]map[string]*GithubTeamRepo {
	return m.teamsrepos
}
func (m *GoliacRemoteMock) AppIds() map[string]int {
	return m.appids
}

type ReconciliatorListenerRecorder struct {
	UsersCreated map[string]string
	UsersRemoved map[string]string

	TeamsCreated      map[string][]string
	TeamMemberAdded   map[string][]string
	TeamMemberRemoved map[string][]string
	TeamDeleted       map[string]bool

	RepositoryCreated              map[string]bool
	RepositoryTeamAdded            map[string][]string
	RepositoryTeamUpdated          map[string][]string
	RepositoryTeamRemoved          map[string][]string
	RepositoriesDeleted            map[string]bool
	RepositoriesUpdatePrivate      map[string]bool
	RepositoriesUpdateArchived     map[string]bool
	RepositoriesSetExternalUser    map[string]string
	RepositoriesRemoveExternalUser map[string]bool

	RuleSetCreated map[string]*GithubRuleSet
	RuleSetUpdated map[string]*GithubRuleSet
	RuleSetDeleted []int
}

func NewReconciliatorListenerRecorder() *ReconciliatorListenerRecorder {
	r := ReconciliatorListenerRecorder{
		UsersCreated:                   make(map[string]string),
		UsersRemoved:                   make(map[string]string),
		TeamsCreated:                   make(map[string][]string),
		TeamMemberAdded:                make(map[string][]string),
		TeamMemberRemoved:              make(map[string][]string),
		TeamDeleted:                    make(map[string]bool),
		RepositoryCreated:              make(map[string]bool),
		RepositoryTeamAdded:            make(map[string][]string),
		RepositoryTeamUpdated:          make(map[string][]string),
		RepositoryTeamRemoved:          make(map[string][]string),
		RepositoriesDeleted:            make(map[string]bool),
		RepositoriesUpdatePrivate:      make(map[string]bool),
		RepositoriesUpdateArchived:     make(map[string]bool),
		RepositoriesSetExternalUser:    make(map[string]string),
		RepositoriesRemoveExternalUser: make(map[string]bool),
		RuleSetCreated:                 make(map[string]*GithubRuleSet),
		RuleSetUpdated:                 make(map[string]*GithubRuleSet),
		RuleSetDeleted:                 make([]int, 0),
	}
	return &r
}
func (r *ReconciliatorListenerRecorder) AddUserToOrg(dryrun bool, ghuserid string) {
	r.UsersCreated[ghuserid] = ghuserid
}
func (r *ReconciliatorListenerRecorder) RemoveUserFromOrg(dryrun bool, ghuserid string) {
	r.UsersRemoved[ghuserid] = ghuserid
}
func (r *ReconciliatorListenerRecorder) CreateTeam(dryrun bool, teamname string, description string, members []string) {
	r.TeamsCreated[teamname] = append(r.TeamsCreated[teamname], members...)
}
func (r *ReconciliatorListenerRecorder) UpdateTeamAddMember(dryrun bool, teamslug string, username string, role string) {
	r.TeamMemberAdded[teamslug] = append(r.TeamMemberAdded[teamslug], username)
}
func (r *ReconciliatorListenerRecorder) UpdateTeamRemoveMember(dryrun bool, teamslug string, username string) {
	r.TeamMemberRemoved[teamslug] = append(r.TeamMemberRemoved[teamslug], username)
}
func (r *ReconciliatorListenerRecorder) DeleteTeam(dryrun bool, teamslug string) {
	r.TeamDeleted[teamslug] = true
}
func (r *ReconciliatorListenerRecorder) CreateRepository(dryrun bool, reponame string, descrition string, writers []string, readers []string, public bool) {
	r.RepositoryCreated[reponame] = true
}
func (r *ReconciliatorListenerRecorder) UpdateRepositoryAddTeamAccess(dryrun bool, reponame string, teamslug string, permission string) {
	r.RepositoryTeamAdded[reponame] = append(r.RepositoryTeamAdded[reponame], teamslug)
}
func (r *ReconciliatorListenerRecorder) UpdateRepositoryUpdateTeamAccess(dryrun bool, reponame string, teamslug string, permission string) {
	r.RepositoryTeamUpdated[reponame] = append(r.RepositoryTeamUpdated[reponame], teamslug)
}
func (r *ReconciliatorListenerRecorder) UpdateRepositoryRemoveTeamAccess(dryrun bool, reponame string, teamslug string) {
	r.RepositoryTeamRemoved[reponame] = append(r.RepositoryTeamRemoved[reponame], teamslug)
}
func (r *ReconciliatorListenerRecorder) DeleteRepository(dryrun bool, reponame string) {
	r.RepositoriesDeleted[reponame] = true
}
func (r *ReconciliatorListenerRecorder) UpdateRepositoryUpdatePrivate(dryrun bool, reponame string, private bool) {
	r.RepositoriesUpdatePrivate[reponame] = true
}
func (r *ReconciliatorListenerRecorder) UpdateRepositoryUpdateArchived(dryrun bool, reponame string, archived bool) {
	r.RepositoriesUpdateArchived[reponame] = true
}
func (r *ReconciliatorListenerRecorder) UpdateRepositorySetExternalUser(dryrun bool, reponame string, githubid string, permission string) {
	r.RepositoriesSetExternalUser[githubid] = permission
}
func (r *ReconciliatorListenerRecorder) UpdateRepositoryRemoveExternalUser(dryrun bool, reponame string, githubid string) {
	r.RepositoriesRemoveExternalUser[githubid] = true
}
func (r *ReconciliatorListenerRecorder) AddRuleset(dryrun bool, ruleset *GithubRuleSet) {
	r.RuleSetCreated[ruleset.Name] = ruleset
}
func (r *ReconciliatorListenerRecorder) UpdateRuleset(dryrun bool, ruleset *GithubRuleSet) {
	r.RuleSetUpdated[ruleset.Name] = ruleset
}
func (r *ReconciliatorListenerRecorder) DeleteRuleset(dryrun bool, rulesetid int) {
	r.RuleSetDeleted = append(r.RuleSetDeleted, rulesetid)
}
func (r *ReconciliatorListenerRecorder) Begin(dryrun bool) {
}
func (r *ReconciliatorListenerRecorder) Rollback(dryrun bool, err error) {
}
func (r *ReconciliatorListenerRecorder) Commit(dryrun bool) {
}

func TestReconciliation(t *testing.T) {

	t.Run("happy path: new team", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		newTeam := &entity.Team{}
		newTeam.Name = "new"
		newTeam.Spec.Owners = []string{"new.owner"}
		newTeam.Spec.Members = []string{"new.member"}
		local.teams["new"] = newTeam

		newOwner := entity.User{}
		newOwner.Name = "new.owner"
		newOwner.Spec.GithubID = "new_owner"
		local.users["new.owner"] = &newOwner
		newMember := entity.User{}
		newMember.Name = "new.member"
		newMember.Spec.GithubID = "new_member"
		local.users["new.member"] = &newMember

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 2 members created
		assert.Equal(t, 2, len(recorder.TeamsCreated["new"]))
		assert.Equal(t, 1, len(recorder.TeamsCreated["new-owners"]))
	})

	t.Run("happy path: new team with non english slug", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		newTeam := &entity.Team{}
		newTeam.Name = "nouveauté"
		newTeam.Spec.Owners = []string{"new.owner"}
		newTeam.Spec.Members = []string{"new.member"}
		local.teams["nouveauté"] = newTeam

		newOwner := entity.User{}
		newOwner.Name = "new.owner"
		newOwner.Spec.GithubID = "new_owner"
		local.users["new.owner"] = &newOwner
		newMember := entity.User{}
		newMember.Name = "new.member"
		newMember.Spec.GithubID = "new_member"
		local.users["new.member"] = &newMember

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 2 members created
		assert.Equal(t, 2, len(recorder.TeamsCreated["nouveaute"]))
		assert.Equal(t, 1, len(recorder.TeamsCreated["nouveaute-owners"]))
	})

	t.Run("happy path: existing team with new members", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing.owner", "existing.owner2"}
		existingTeam.Spec.Members = []string{"existing.member"}
		local.teams["existing"] = existingTeam

		existing_owner := entity.User{}
		existing_owner.Name = "existing.owner"
		existing_owner.Spec.GithubID = "existing_owner"
		local.users["existing.owner"] = &existing_owner

		existing_owner2 := entity.User{}
		existing_owner2.Name = "existing.owner2"
		existing_owner2.Spec.GithubID = "existing_owner2"
		local.users["existing.owner2"] = &existing_owner2

		existing_member := entity.User{}
		existing_member.Name = "existing.member"
		existing_member.Spec.GithubID = "existing_member"
		local.users["existing.member"] = &existing_member

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["existing"] = existing
		existingowners := &GithubTeam{
			Name:    "existing-owners",
			Slug:    "existing-owners",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["existing-owners"] = existingowners

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 members added
		assert.Equal(t, 0, len(recorder.TeamsCreated))
		assert.Equal(t, 1, len(recorder.TeamMemberAdded["existing"]))
	})

	t.Run("happy path: existing team with non english slug with new members", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		existingTeam := &entity.Team{}
		existingTeam.Name = "exist ing"
		existingTeam.Spec.Owners = []string{"existing.owner", "existing.owner2"}
		existingTeam.Spec.Members = []string{"existing.member"}
		local.teams["exist ing"] = existingTeam

		existing_owner := entity.User{}
		existing_owner.Name = "existing.owner"
		existing_owner.Spec.GithubID = "existing_owner"
		local.users["existing.owner"] = &existing_owner

		existing_owner2 := entity.User{}
		existing_owner2.Name = "existing.owner2"
		existing_owner2.Spec.GithubID = "existing_owner2"
		local.users["existing.owner2"] = &existing_owner2

		existing_member := entity.User{}
		existing_member.Name = "existing.member"
		existing_member.Spec.GithubID = "existing_member"
		local.users["existing.member"] = &existing_member

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "exist ing",
			Slug:    "exist-ing",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["exist-ing"] = existing

		existingowners := &GithubTeam{
			Name:    "exist ing-owners",
			Slug:    "exist-ing-owners",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["exist-ing-owners"] = existingowners

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 members added
		assert.Equal(t, "exist-ing", remote.TeamSlugByName()["exist ing"])
		assert.Equal(t, 0, len(recorder.TeamsCreated))
		assert.Equal(t, 1, len(recorder.TeamMemberAdded["exist-ing"]))
	})

	t.Run("happy path: new team + adding everyone team", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{
			EveryoneTeamEnabled: true,
		}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		newTeam := &entity.Team{}
		newTeam.Name = "new"
		newTeam.Spec.Owners = []string{"new.owner"}
		newTeam.Spec.Members = []string{"new.member"}
		local.teams["new"] = newTeam

		newOwner := entity.User{}
		newOwner.Name = "new.owner"
		newOwner.Spec.GithubID = "new_owner"
		local.users["new.owner"] = &newOwner
		newMember := entity.User{}
		newMember.Name = "new.member"
		newMember.Spec.GithubID = "new_member"
		local.users["new.member"] = &newMember

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 2 members created
		assert.Equal(t, 2, len(recorder.TeamsCreated["new"]))
		assert.Equal(t, 1, len(recorder.TeamsCreated["new-owners"]))
		// and the everyone team
		assert.Equal(t, 2, len(recorder.TeamsCreated["everyone"]))
	})

	t.Run("happy path: removed team without destructive operation", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		removing := &GithubTeam{
			Name:    "removing",
			Slug:    "removing",
			Members: []string{"existing_owner", "existing_owner"},
		}
		remote.teams["removing"] = removing

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 team deleted
		assert.Equal(t, 0, len(recorder.TeamDeleted))
	})

	t.Run("happy path: removed team", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()
		repoconfig := &config.RepositoryConfig{}
		repoconfig.DestructiveOperations.AllowDestructiveTeams = true
		r := NewGoliacReconciliatorImpl(recorder, repoconfig)
		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		removing := &GithubTeam{
			Name:    "removing",
			Slug:    "removing",
			Members: []string{"existing_owner", "existing_owner"},
		}
		remote.teams["removing"] = removing

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 team deleted
		assert.Equal(t, 1, len(recorder.TeamDeleted))
	})

	t.Run("happy path: new repo without owner", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()
		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		newRepo := &entity.Repository{}
		newRepo.Name = "new"
		newRepo.Spec.Readers = []string{}
		newRepo.Spec.Writers = []string{}
		local.repos["new"] = newRepo

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 repo created
		assert.Equal(t, 1, len(recorder.RepositoryCreated))
	})

	t.Run("happy path: new repo with owner", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		newRepo := &entity.Repository{}
		newRepo.Name = "new"
		newRepo.Spec.Readers = []string{}
		newRepo.Spec.Writers = []string{}
		owner := "existing"
		newRepo.Owner = &owner
		local.repos["new"] = newRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{"existing_member"}
		local.teams["existing"] = existingTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["existing"] = existing

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 repo created
		assert.Equal(t, 1, len(recorder.RepositoryCreated))
	})

	t.Run("happy path: existing repo with new owner (from read to write)", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		lRepo := &entity.Repository{}
		lRepo.Name = "myrepo"
		lRepo.Spec.Readers = []string{}
		lRepo.Spec.Writers = []string{}
		lowner := "existing"
		lRepo.Owner = &lowner
		local.repos["myrepo"] = lRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{"existing_member"}
		local.teams["existing"] = existingTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["existing"] = existing
		rRepo := GithubRepository{
			Name: "myrepo",
		}
		remote.repos["myrepo"] = &rRepo

		remote.teamsrepos["existing"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["existing"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "READ",
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 team updated
		assert.Equal(t, 0, len(recorder.RepositoryCreated))
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
		assert.Equal(t, 1, len(recorder.RepositoryTeamRemoved))
		assert.Equal(t, 1, len(recorder.RepositoryTeamAdded))
		assert.Equal(t, 0, len(recorder.RepositoryTeamUpdated))
	})

	t.Run("happy path: existing repo without new owner but with everyone team", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{
			EveryoneTeamEnabled: true,
		}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		lRepo := &entity.Repository{}
		lRepo.Name = "myrepo"
		lRepo.Spec.Readers = []string{}
		lRepo.Spec.Writers = []string{}
		lowner := "existing"
		lRepo.Owner = &lowner
		local.repos["myrepo"] = lRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{"existing_member"}
		local.teams["existing"] = existingTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["existing"] = existing
		rRepo := GithubRepository{
			Name: "myrepo",
		}
		remote.repos["myrepo"] = &rRepo

		remote.teamsrepos["existing"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["existing"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "WRITE",
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 team updated
		assert.Equal(t, 0, len(recorder.RepositoryCreated))
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
		assert.Equal(t, 0, len(recorder.RepositoryTeamRemoved))
		// we have a new "everyone" team for the repository
		assert.Equal(t, 1, len(recorder.RepositoryTeamAdded))
		assert.Equal(t, 0, len(recorder.RepositoryTeamUpdated))
	})

	t.Run("happy path: add a team to an existing repo", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		lRepo := &entity.Repository{}
		lRepo.Name = "myrepo"
		lRepo.Spec.Readers = []string{"reader"}
		lRepo.Spec.Writers = []string{}
		lowner := "existing"
		lRepo.Owner = &lowner
		local.repos["myrepo"] = lRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{"existing_member"}
		local.teams["existing"] = existingTeam

		readerTeam := &entity.Team{}
		readerTeam.Name = "reader"
		readerTeam.Spec.Owners = []string{"existing_owner"}
		readerTeam.Spec.Members = []string{"existing_member"}
		local.teams["reader"] = readerTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner", "existing_member"},
		}
		reader := &GithubTeam{
			Name:    "reader",
			Slug:    "reader",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["existing"] = existing
		remote.teams["reader"] = reader
		rRepo := GithubRepository{
			Name: "myrepo",
		}
		remote.repos["myrepo"] = &rRepo

		remote.teamsrepos["existing"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["existing"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "ADMIN",
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 team added
		assert.Equal(t, 0, len(recorder.RepositoryCreated))
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
		assert.Equal(t, 0, len(recorder.RepositoryTeamRemoved))
		assert.Equal(t, 1, len(recorder.RepositoryTeamAdded))
		assert.Equal(t, 0, len(recorder.RepositoryTeamUpdated))
	})

	t.Run("happy path: remove a team from an existing repo", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		lRepo := &entity.Repository{}
		lRepo.Name = "myrepo"
		lRepo.Spec.Readers = []string{}
		lRepo.Spec.Writers = []string{}
		lowner := "existing"
		lRepo.Owner = &lowner
		local.repos["myrepo"] = lRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{"existing_member"}
		local.teams["existing"] = existingTeam

		readerTeam := &entity.Team{}
		readerTeam.Name = "reader"
		readerTeam.Spec.Owners = []string{"existing_owner"}
		readerTeam.Spec.Members = []string{"existing_member"}
		local.teams["reader"] = readerTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner", "existing_member"},
		}
		reader := &GithubTeam{
			Name:    "reader",
			Slug:    "reader",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["existing"] = existing
		remote.teams["reader"] = reader
		rRepo := GithubRepository{
			Name: "myrepo",
		}
		remote.repos["myrepo"] = &rRepo

		remote.teamsrepos["existing"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["existing"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "WRITE",
		}
		remote.teamsrepos["reader"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["reader"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "WRITE",
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 team removed
		assert.Equal(t, 0, len(recorder.RepositoryCreated))
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
		assert.Equal(t, 1, len(recorder.RepositoryTeamRemoved))
		assert.Equal(t, 0, len(recorder.RepositoryTeamAdded))
		assert.Equal(t, 0, len(recorder.RepositoryTeamUpdated))
	})

	t.Run("happy path: remove a team member", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		lRepo := &entity.Repository{}
		lRepo.Name = "myrepo"
		lRepo.Spec.Readers = []string{}
		lRepo.Spec.Writers = []string{}
		lowner := "existing"
		lRepo.Owner = &lowner
		local.repos["myrepo"] = lRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{}
		local.teams["existing"] = existingTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["existing"] = existing
		rRepo := GithubRepository{
			Name: "myrepo",
		}
		remote.repos["myrepo"] = &rRepo

		remote.teamsrepos["existing"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["existing"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "WRITE",
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 member removed
		assert.Equal(t, 0, len(recorder.RepositoryCreated))
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
		assert.Equal(t, 0, len(recorder.RepositoryTeamRemoved))
		assert.Equal(t, 0, len(recorder.RepositoryTeamAdded))
		assert.Equal(t, 0, len(recorder.RepositoryTeamUpdated))
		assert.Equal(t, 1, len(recorder.TeamMemberRemoved))
	})

	t.Run("happy path: add a team AND add it to an existing repo", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()
		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}
		lRepo := &entity.Repository{}
		lRepo.Name = "myrepo"
		lRepo.Spec.Readers = []string{"reader"}
		lRepo.Spec.Writers = []string{}
		lowner := "existing"
		lRepo.Owner = &lowner
		local.repos["myrepo"] = lRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{"existing_member"}
		local.teams["existing"] = existingTeam

		readerTeam := &entity.Team{}
		readerTeam.Name = "reader"
		readerTeam.Spec.Owners = []string{"existing_owner"}
		readerTeam.Spec.Members = []string{"existing_member"}
		local.teams["reader"] = readerTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner", "existing_member"},
		}
		remote.teams["existing"] = existing
		rRepo := GithubRepository{
			Name: "myrepo",
		}
		remote.repos["myrepo"] = &rRepo

		remote.teamsrepos["existing"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["existing"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "WRITE",
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 repo updated
		assert.Equal(t, 0, len(recorder.RepositoryCreated))
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
		assert.Equal(t, 0, len(recorder.RepositoryTeamRemoved))
		assert.Equal(t, 1, len(recorder.RepositoryTeamAdded))
	})

	t.Run("happy path: existing repo with new external write collaborator", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users:     make(map[string]*entity.User),
			externals: make(map[string]*entity.User),
			teams:     make(map[string]*entity.Team),
			repos:     make(map[string]*entity.Repository),
		}
		outside1 := entity.User{}
		outside1.Name = "outside1"
		outside1.Spec.GithubID = "outside1-githubid"
		local.externals["outside1"] = &outside1

		lRepo := &entity.Repository{}
		lRepo.Name = "myrepo"
		lRepo.Spec.Readers = []string{}
		lRepo.Spec.Writers = []string{}
		lRepo.Spec.ExternalUserWriters = []string{"outside1"}
		lowner := "existing"
		lRepo.Owner = &lowner
		local.repos["myrepo"] = lRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{}
		local.teams["existing"] = existingTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner"},
		}
		remote.teams["existing"] = existing
		rRepo := GithubRepository{
			Name:          "myrepo",
			ExternalUsers: make(map[string]string),
		}
		remote.repos["myrepo"] = &rRepo

		remote.teamsrepos["existing"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["existing"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "WRITE",
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 team updated
		assert.Equal(t, 0, len(recorder.RepositoryCreated))
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
		assert.Equal(t, 0, len(recorder.RepositoryTeamRemoved))
		assert.Equal(t, 0, len(recorder.RepositoryTeamAdded))
		assert.Equal(t, 0, len(recorder.RepositoryTeamUpdated))
		assert.Equal(t, 1, len(recorder.RepositoriesSetExternalUser))
		assert.Equal(t, 0, len(recorder.RepositoriesRemoveExternalUser))
	})

	t.Run("happy path: existing repo with deleted external write collaborator", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users:     make(map[string]*entity.User),
			externals: make(map[string]*entity.User),
			teams:     make(map[string]*entity.Team),
			repos:     make(map[string]*entity.Repository),
		}

		lRepo := &entity.Repository{}
		lRepo.Name = "myrepo"
		lRepo.Spec.Readers = []string{}
		lRepo.Spec.Writers = []string{}
		lRepo.Spec.ExternalUserWriters = []string{}
		lowner := "existing"
		lRepo.Owner = &lowner
		local.repos["myrepo"] = lRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{}
		local.teams["existing"] = existingTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner"},
		}
		remote.teams["existing"] = existing
		rRepo := GithubRepository{
			Name:          "myrepo",
			ExternalUsers: make(map[string]string),
		}
		rRepo.ExternalUsers["outside1-githubid"] = "WRITE"
		remote.repos["myrepo"] = &rRepo

		remote.teamsrepos["existing"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["existing"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "WRITE",
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 team updated
		assert.Equal(t, 0, len(recorder.RepositoryCreated))
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
		assert.Equal(t, 0, len(recorder.RepositoryTeamRemoved))
		assert.Equal(t, 0, len(recorder.RepositoryTeamAdded))
		assert.Equal(t, 0, len(recorder.RepositoryTeamUpdated))
		assert.Equal(t, 0, len(recorder.RepositoriesSetExternalUser))
		assert.Equal(t, 1, len(recorder.RepositoriesRemoveExternalUser))
	})

	t.Run("happy path: existing repo with changed external write collaborator (from read to write)", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users:     make(map[string]*entity.User),
			externals: make(map[string]*entity.User),
			teams:     make(map[string]*entity.Team),
			repos:     make(map[string]*entity.Repository),
		}

		outside1 := entity.User{}
		outside1.Name = "outside1"
		outside1.Spec.GithubID = "outside1-githubid"
		local.externals["outside1"] = &outside1

		lRepo := &entity.Repository{}
		lRepo.Name = "myrepo"
		lRepo.Spec.Readers = []string{}
		lRepo.Spec.Writers = []string{}
		lRepo.Spec.ExternalUserWriters = []string{}
		lRepo.Spec.ExternalUserReaders = []string{"outside1"}
		lowner := "existing"
		lRepo.Owner = &lowner
		local.repos["myrepo"] = lRepo

		existingTeam := &entity.Team{}
		existingTeam.Name = "existing"
		existingTeam.Spec.Owners = []string{"existing_owner"}
		existingTeam.Spec.Members = []string{}
		local.teams["existing"] = existingTeam

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		existing := &GithubTeam{
			Name:    "existing",
			Slug:    "existing",
			Members: []string{"existing_owner"},
		}
		remote.teams["existing"] = existing
		rRepo := GithubRepository{
			Name:          "myrepo",
			ExternalUsers: make(map[string]string),
		}
		rRepo.ExternalUsers["outside1-githubid"] = "WRITE"
		remote.repos["myrepo"] = &rRepo

		remote.teamsrepos["existing"] = make(map[string]*GithubTeamRepo)
		remote.teamsrepos["existing"]["myrepo"] = &GithubTeamRepo{
			Name:       "myrepo",
			Permission: "WRITE",
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 team updated
		assert.Equal(t, 0, len(recorder.RepositoryCreated))
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
		assert.Equal(t, 0, len(recorder.RepositoryTeamRemoved))
		assert.Equal(t, 0, len(recorder.RepositoryTeamAdded))
		assert.Equal(t, 0, len(recorder.RepositoryTeamUpdated))
		assert.Equal(t, 1, len(recorder.RepositoriesSetExternalUser))
		assert.Equal(t, 0, len(recorder.RepositoriesRemoveExternalUser))
	})

	t.Run("happy path: removed repo without destructive operation", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		removing := &GithubRepository{
			Name: "removing",
		}
		remote.repos["removing"] = removing

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 repo deleted
		assert.Equal(t, 0, len(recorder.RepositoriesDeleted))
	})

	t.Run("happy path: removed repo", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()
		repoconfig := &config.RepositoryConfig{}
		repoconfig.DestructiveOperations.AllowDestructiveRepositories = true
		r := NewGoliacReconciliatorImpl(recorder, repoconfig)

		local := GoliacLocalMock{
			users: make(map[string]*entity.User),
			teams: make(map[string]*entity.Team),
			repos: make(map[string]*entity.Repository),
		}

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}
		removing := &GithubRepository{
			Name: "removing",
		}
		remote.repos["removing"] = removing

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 repo deleted
		assert.Equal(t, 1, len(recorder.RepositoriesDeleted))
	})
}

func TestReconciliationRulesets(t *testing.T) {

	t.Run("happy path: no new ruleset in goliac conf", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()
		repoconf := config.RepositoryConfig{}

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users:    make(map[string]*entity.User),
			teams:    make(map[string]*entity.Team),
			repos:    make(map[string]*entity.Repository),
			rulesets: make(map[string]*entity.RuleSet),
		}

		newRuleset := &entity.RuleSet{}
		newRuleset.Name = "new"
		newRuleset.Spec.Enforcement = "evaluate"
		newRuleset.Spec.Rules = append(newRuleset.Spec.Rules, struct {
			Ruletype   string
			Parameters entity.RuleSetParameters
		}{
			"required_signatures", entity.RuleSetParameters{},
		})
		local.rulesets["new"] = newRuleset

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 ruleset created
		assert.Equal(t, 0, len(recorder.RuleSetCreated))
		assert.Equal(t, 0, len(recorder.RuleSetUpdated))
		assert.Equal(t, 0, len(recorder.RuleSetDeleted))
	})

	t.Run("happy path: new ruleset", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{
			Rulesets: make([]struct {
				Pattern string
				Ruleset string
			}, 0),
		}
		repoconf.Rulesets = append(repoconf.Rulesets, struct {
			Pattern string
			Ruleset string
		}{
			Pattern: ".*",
			Ruleset: "new",
		})

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users:    make(map[string]*entity.User),
			teams:    make(map[string]*entity.Team),
			repos:    make(map[string]*entity.Repository),
			rulesets: make(map[string]*entity.RuleSet),
		}

		newRuleset := &entity.RuleSet{}
		newRuleset.Name = "new"
		newRuleset.Spec.Enforcement = "evaluate"
		newRuleset.Spec.Rules = append(newRuleset.Spec.Rules, struct {
			Ruletype   string
			Parameters entity.RuleSetParameters
		}{
			"required_signatures", entity.RuleSetParameters{},
		})
		local.rulesets["new"] = newRuleset

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 ruleset created
		assert.Equal(t, 1, len(recorder.RuleSetCreated))
		assert.Equal(t, 0, len(recorder.RuleSetUpdated))
		assert.Equal(t, 0, len(recorder.RuleSetDeleted))
	})

	t.Run("happy path: update ruleset (enforcement)", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{
			Rulesets: make([]struct {
				Pattern string
				Ruleset string
			}, 0),
		}
		repoconf.Rulesets = append(repoconf.Rulesets, struct {
			Pattern string
			Ruleset string
		}{
			Pattern: ".*",
			Ruleset: "update",
		})

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users:    make(map[string]*entity.User),
			teams:    make(map[string]*entity.Team),
			repos:    make(map[string]*entity.Repository),
			rulesets: make(map[string]*entity.RuleSet),
		}

		lRuleset := &entity.RuleSet{}
		lRuleset.Name = "update"
		lRuleset.Spec.Enforcement = "evaluate"
		lRuleset.Spec.Rules = append(lRuleset.Spec.Rules, struct {
			Ruletype   string
			Parameters entity.RuleSetParameters
		}{
			"required_signatures", entity.RuleSetParameters{},
		})
		local.rulesets["update"] = lRuleset

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}

		rRuleset := &GithubRuleSet{
			Name:        "update",
			Enforcement: "active",
			Rules:       make(map[string]entity.RuleSetParameters),
		}
		rRuleset.Rules["required_signatures"] = entity.RuleSetParameters{}
		remote.rulesets["update"] = rRuleset

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 ruleset created
		assert.Equal(t, 0, len(recorder.RuleSetCreated))
		assert.Equal(t, 1, len(recorder.RuleSetUpdated))
		assert.Equal(t, 0, len(recorder.RuleSetDeleted))
	})

	t.Run("happy path: delete ruleset", func(t *testing.T) {
		recorder := NewReconciliatorListenerRecorder()

		repoconf := config.RepositoryConfig{
			Rulesets: make([]struct {
				Pattern string
				Ruleset string
			}, 0),
		}
		repoconf.DestructiveOperations.AllowDestructiveRulesets = true

		r := NewGoliacReconciliatorImpl(recorder, &repoconf)

		local := GoliacLocalMock{
			users:    make(map[string]*entity.User),
			teams:    make(map[string]*entity.Team),
			repos:    make(map[string]*entity.Repository),
			rulesets: make(map[string]*entity.RuleSet),
		}

		remote := GoliacRemoteMock{
			users:      make(map[string]string),
			teams:      make(map[string]*GithubTeam),
			repos:      make(map[string]*GithubRepository),
			teamsrepos: make(map[string]map[string]*GithubTeamRepo),
			rulesets:   make(map[string]*GithubRuleSet),
			appids:     make(map[string]int),
		}

		rRuleset := &GithubRuleSet{
			Name:        "delete",
			Enforcement: "active",
			Rules:       make(map[string]entity.RuleSetParameters),
		}
		rRuleset.Rules["required_signatures"] = entity.RuleSetParameters{}
		remote.rulesets["delete"] = rRuleset

		r.Reconciliate(context.TODO(), &local, &remote, "teams", false)

		// 1 ruleset created
		assert.Equal(t, 0, len(recorder.RuleSetCreated))
		assert.Equal(t, 0, len(recorder.RuleSetUpdated))
		assert.Equal(t, 1, len(recorder.RuleSetDeleted))
	})
}
