package envrouter

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"os"

	cryptossh "golang.org/x/crypto/ssh"
	"strings"
	"time"
)

type GitClient interface {
	GetCommitByHash(applicationName string, hash string) (*api.Commit, error)
	GetLatestCommit(repositoryName string, ref string) (*api.Commit, error)
	GetAllLatestCommits(repositoryName string, supplier func(ref string, commit *api.Commit)) error
}

type gitClient struct {
	applicationService       ApplicationService
	repositoryService        RepositoryService
	credentialsSecretService CredentialsSecretService
}

func NewGitClient(
	applicationService ApplicationService,
	repositoryService RepositoryService,
	credentialsSecretService CredentialsSecretService,
) GitClient {
	return &gitClient{
		applicationService,
		repositoryService,
		credentialsSecretService,
	}
}

func (g *gitClient) getRepository(repositoryName string) (*git.Repository, error) {
	options, _, _, err := g.getGitOptions(repositoryName)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/tmp/git/%s", repositoryName)
	var r *git.Repository
	if _, err := os.Stat(path); os.IsNotExist(err) {
		r, err = git.PlainClone(path, true, options)
	} else {
		r, err = git.PlainOpenWithOptions(path, &git.PlainOpenOptions{})
		err = r.Fetch(&git.FetchOptions{RemoteName: "origin", Depth: 1, Auth: options.Auth})
		if err != nil {
			return nil, err
		}
	}
	return r, err
}

func (g *gitClient) GetCommitByHash(applicationName string, hash string) (*api.Commit, error) {
	application, err := g.applicationService.FindByName(applicationName)
	if err != nil {
		return nil, err
	}
	r, err := g.getRepository(*application.RepositoryName)
	if err != nil {
		panic(err)
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

func (g *gitClient) GetLatestCommit(repositoryName string, ref string) (*api.Commit, error) {
	r, err := g.getRepository(repositoryName)
	if err != nil {
		panic(err)
	}
	h, err := r.ResolveRevision(plumbing.Revision("origin/" + ref))
	if err != nil {
		panic(err)
	}
	return g.getCommitByHash(r, *h)
}

func (g *gitClient) GetAllLatestCommits(repositoryName string, supplier func(ref string, commit *api.Commit)) error {
	r, err := g.getRepository(repositoryName)
	if err != nil {
		return err
	}
	iter, err := r.References()
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		if strings.HasPrefix(string(ref.Name()), "refs/remotes/origin/") {
			refName := strings.Replace(ref.Name().Short(), "origin/", "", 1)
			log.Infof("ref: %v", refName)
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

func (g *gitClient) getGitOptions(repositoryName string) (*git.CloneOptions, billy.Filesystem, *memory.Storage, error) {
	repository, err := g.repositoryService.FindByName(repositoryName)
	if err != nil {
		return nil, nil, nil, err
	}
	credentials, err := g.credentialsSecretService.FindByName(repository.CredentialsSecret)
	if err != nil {
		return nil, nil, nil, err
	}

	options := &git.CloneOptions{
		NoCheckout: true,
		URL:        repository.Url,
		Progress:   os.Stdout,
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
				panic(err)
			}
			key.HostKeyCallbackHelper = ssh.HostKeyCallbackHelper{
				HostKeyCallback: cryptossh.InsecureIgnoreHostKey(),
			}
			options.Auth = key
		}
	}
	return options, memfs.New(), memory.NewStorage(), nil
}

type GitStorage interface {
	GetCommitByHash(applicationName string, hash string) (*api.Commit, error)
	GetLatestCommit(repositoryName string, ref string, force bool) (*api.Commit, error)
	Scan(repositoryName string) error
}

type gitStorage struct {
	gitClient GitClient
	commits   map[string]*api.Commit
	branches  map[string]map[string]*api.Commit
}

func NewGitStorage(
	gitClient GitClient,
) GitStorage {
	return &gitStorage{
		gitClient: gitClient,
		commits:   map[string]*api.Commit{},
		branches:  map[string]map[string]*api.Commit{},
	}
}

func (g *gitStorage) GetCommitByHash(applicationName string, hash string) (*api.Commit, error) {
	if commit, ok := g.commits[hash]; ok {
		return commit, nil
	}
	commit, err := g.gitClient.GetCommitByHash(applicationName, hash)
	if err != nil {
		return nil, err
	}
	if commit != nil {
		g.commits[hash] = commit
	}
	return commit, nil
}

func (g *gitStorage) GetLatestCommit(repositoryName string, ref string, force bool) (*api.Commit, error) {
	if !force {
		if repository, ok := g.branches[repositoryName]; ok {
			if commit, ok := repository[ref]; ok {
				return commit, nil
			}
		}
	}
	commit, err := g.gitClient.GetLatestCommit(repositoryName, ref)
	if err != nil {
		return nil, err
	}
	if commit != nil {
		if _, ok := g.branches[repositoryName]; !ok {
			g.branches[repositoryName] = map[string]*api.Commit{}
		}
		g.branches[repositoryName][ref] = commit
	}
	return nil, err
}

func (g *gitStorage) Scan(repositoryName string) error {
	err := g.gitClient.GetAllLatestCommits(repositoryName, func(ref string, commit *api.Commit) {
		g.commits[commit.Sha] = commit
		if _, ok := g.branches[repositoryName]; !ok {
			g.branches[repositoryName] = map[string]*api.Commit{}
		}
		g.branches[repositoryName][ref] = commit
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
			err := g.gitStorage.Scan(v.Name)
			if err != nil {
				log.Errorf("Error on git scan %v", err)
			}
		}
		log.Info("Scan finished")
		time.Sleep(30 * time.Second)
	}
}
