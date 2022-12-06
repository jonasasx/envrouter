package envrouter

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/utils"
	"os"
	"reflect"

	cryptossh "golang.org/x/crypto/ssh"
	"strings"
	"time"
)

type GitClient interface {
	GetCommitByHash(repositoryName string, hash string) (*api.Commit, error)
	SupplyRefsHeads(repositoryName string, supplier func(ref string, commit *api.Commit)) error
}

type gitClient struct {
	repositoryService        RepositoryService
	credentialsSecretService CredentialsSecretService
}

func NewGitClient(
	repositoryService RepositoryService,
	credentialsSecretService CredentialsSecretService,
) GitClient {
	return &gitClient{
		repositoryService,
		credentialsSecretService,
	}
}

func (g *gitClient) getRepository(repositoryName string) (*git.Repository, error) {
	options, err := g.getGitOptions(repositoryName)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/tmp/git/%s", repositoryName)
	var r *git.Repository
	if _, err := os.Stat(path); os.IsNotExist(err) {
		r, err = git.PlainClone(path, true, options)
		if err != nil {
			return nil, err
		}
	} else {
		log.Debugf("Start fetching %s", repositoryName)
		r, err = git.PlainOpenWithOptions(path, &git.PlainOpenOptions{})
		if err != nil {
			return nil, err
		}
		err = r.Fetch(&git.FetchOptions{RemoteName: "origin", Auth: options.Auth, Progress: options.Progress})
		if err != nil && fmt.Sprint(err) != "already up-to-date" {
			return nil, err
		}
		log.Debugf("Fetched %s", repositoryName)
	}
	return r, err
}

func (g *gitClient) GetCommitByHash(repositoryName string, hash string) (*api.Commit, error) {
	r, err := g.getRepository(repositoryName)
	if err != nil {
		return nil, err
	}
	return g.getCommitByHash(r, plumbing.NewHash(hash))
}

func (g *gitClient) getCommitByHash(repository *git.Repository, hash plumbing.Hash) (*api.Commit, error) {
	commit, err := repository.CommitObject(hash)
	if err != nil {
		return nil, err
	}
	when := commit.Author.When.Format(time.RFC3339)
	return &api.Commit{
		Author:    &commit.Author.Name,
		Message:   &commit.Message,
		Sha:       hash.String(),
		Timestamp: &when,
	}, nil
}

func (g *gitClient) SupplyRefsHeads(repositoryName string, supplier func(ref string, commit *api.Commit)) error {
	r, err := g.getRepository(repositoryName)
	if err != nil {
		return err
	}
	iter, err := r.References()
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		if strings.HasPrefix(string(ref.Name()), "refs/remotes/origin/") {
			refName := strings.Replace(ref.Name().Short(), "origin/", "", 1)
			log.Debugf("ref: %v", refName)
			commit, err := g.getCommitByHash(r, ref.Hash())
			if err != nil {
				return err
			}
			supplier(refName, commit)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *gitClient) getGitOptions(repositoryName string) (*git.CloneOptions, error) {
	repository, err := g.repositoryService.FindByName(repositoryName)
	if err != nil {
		return nil, err
	}
	credentials, err := g.credentialsSecretService.FindByName(repository.CredentialsSecret)
	if err != nil {
		return nil, err
	}

	options := &git.CloneOptions{
		NoCheckout: true,
		URL:        repository.Url,
		Progress:   nil,
	}

	if credentials != nil {
		if len(credentials.Username) > 0 && len(credentials.Password) > 0 {
			options.Auth = &http.BasicAuth{
				Username: credentials.Username,
				Password: credentials.Password,
			}
		} else if len(credentials.Key) > 0 {
			key, err := ssh.NewPublicKeys("git", []byte(credentials.Key), "")
			if err != nil {
				return nil, err
			}
			key.HostKeyCallbackHelper = ssh.HostKeyCallbackHelper{
				HostKeyCallback: cryptossh.InsecureIgnoreHostKey(),
			}
			options.Auth = key
		}
	}
	return options, nil
}

type GitStorage interface {
	GetAllRefsHeads() ([]*api.Ref, error)
	GetCommitByHash(applicationName string, hash string) (*api.Commit, error)
	ScanRefsHeads(repositoryName string) error
}

type gitStorage struct {
	gitClient      GitClient
	eventsObserver utils.Observer
	commits        map[string]*api.Commit
	branches       map[string]map[string]*api.Commit
}

func NewGitStorage(
	gitClient GitClient,
	eventsObserver utils.Observer,
) GitStorage {
	return &gitStorage{
		gitClient:      gitClient,
		eventsObserver: eventsObserver,
		commits:        map[string]*api.Commit{},
		branches:       map[string]map[string]*api.Commit{},
	}
}

func (g *gitStorage) GetAllRefsHeads() ([]*api.Ref, error) {
	var result []*api.Ref
	for repositoryName, v := range g.branches {
		for ref, commit := range v {
			result = append(result, &api.Ref{
				Repository: repositoryName,
				Commit:     *commit,
				Ref:        ref,
			})
		}
	}
	return result, nil
}

func (g *gitStorage) GetCommitByHash(repositoryName string, hash string) (*api.Commit, error) {
	if commit, ok := g.commits[hash]; ok {
		return commit, nil
	}
	commit, err := g.gitClient.GetCommitByHash(repositoryName, hash)
	if err != nil {
		return nil, err
	}
	if commit != nil {
		g.commits[hash] = commit
	}
	return commit, nil
}

func (g *gitStorage) addRefHeadCommit(repositoryName string, ref string, commit *api.Commit) {
	g.commits[commit.Sha] = commit
	if _, ok := g.branches[repositoryName]; !ok {
		g.branches[repositoryName] = map[string]*api.Commit{}
	}
	if oldCommit, ok := g.branches[repositoryName][ref]; !ok || !reflect.DeepEqual(oldCommit, commit) {
		g.branches[repositoryName][ref] = commit
		g.eventsObserver.Publish(nil, api.SSEvent{
			ItemType: "RefHead",
			Item: api.Ref{
				Repository: repositoryName,
				Commit:     *commit,
				Ref:        ref,
			},
			Event: "UPDATED",
		})
	}
}

func (g *gitStorage) ScanRefsHeads(repositoryName string) error {
	err := g.gitClient.SupplyRefsHeads(repositoryName, func(ref string, commit *api.Commit) {
		g.addRefHeadCommit(repositoryName, ref, commit)
	})
	if err != nil {
		return err
	}
	return nil
}

type GitScanJob interface {
	Scan()
}

type gitScanJob struct {
	repositoryService RepositoryService
	gitStorage        GitStorage
}

func NewGitScanJob(repositoryService RepositoryService, gitStorage GitStorage) GitScanJob {
	return &gitScanJob{
		repositoryService: repositoryService,
		gitStorage:        gitStorage,
	}
}

func (g *gitScanJob) Scan() {
	for {
		rs, err := g.repositoryService.FindAll()
		if err != nil {
			log.Errorf("Error on git scan %v", err)
		}
		for _, v := range rs {
			err := g.gitStorage.ScanRefsHeads(v.Name)
			if err != nil {
				log.Errorf("Error on git scan %v", err)
			}
		}
		log.Info("ScanRefsHeads finished")
		time.Sleep(30 * time.Second)
	}
}
