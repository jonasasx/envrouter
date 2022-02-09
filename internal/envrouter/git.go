package envrouter

import (
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	cryptossh "golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"time"
)

type GitClient interface {
	GetCommitByHash(applicationName string, hash string) (*api.Commit, error)
	GetLatestCommit(repositoryName string, ref string) (*api.Commit, error)
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

func (g *gitClient) GetCommitByHash(applicationName string, hash string) (*api.Commit, error) {
	application, err := g.applicationService.FindByName(applicationName)
	if err != nil {
		return nil, err
	}
	options, fs, storer, err := g.getGitOptions(*application.RepositoryName)
	if err != nil {
		return nil, err
	}
	r, err := git.Clone(storer, fs, options)
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
	options, fs, storer, err := g.getGitOptions(repositoryName)
	if err != nil {
		return nil, err
	}
	r, err := git.Clone(storer, fs, options)
	if err != nil {
		panic(err)
	}
	h, err := r.ResolveRevision(plumbing.Revision("origin/" + ref))
	if err != nil {
		panic(err)
	}
	return g.getCommitByHash(r, *h)
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

func (g gitStorage) GetLatestCommit(repositoryName string, ref string, force bool) (*api.Commit, error) {
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
