package service

import (
	"errors"
	"fmt"
	"github.com/alt-golang/config"
	"github.com/alt-golang/logger"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
)

type ConfigService struct {
	Logger      logger.Logger
	Dir         string
	GitUrl      string
	GitBranch   string
	GitUsername string
	GitToken    string
}

func (configService *ConfigService) Init() {
	configService.GitCloneAndCheckout()
}

func (configService *ConfigService) GitCloneAndCheckout() {
	if configService.GitUrl != "" {
		configService.Logger.Info("Cloning git repo at URL: " + configService.GitUrl)
		repo, err := git.PlainClone(configService.Dir, false, &git.CloneOptions{
			URL:      configService.GitUrl,
			Progress: os.Stdout,
			Auth: &http.BasicAuth{
				Username: configService.GitUsername,
				Password: configService.GitToken,
			},
		})
		if err != nil {
			fmt.Println(err)
			if fmt.Sprint(err) != "repository already exists" {
				configService.Logger.Error("Git Clone Failed: system startup exiting:" + fmt.Sprint(err))
				os.Exit(1)
			}

			configService.Logger.Info("Opening git repo at path: " + configService.Dir)
			repo, err = git.PlainOpen(configService.Dir)
			fmt.Println(err)
		}
		w, err := repo.Worktree()
		branchName := plumbing.NewBranchReferenceName(configService.GitBranch)
		headRef, err := repo.Head()

		configService.Logger.Info("Checking out branch: " + configService.GitBranch)
		branchRef := plumbing.NewHashReference(branchName, headRef.Hash())
		err = repo.Storer.SetReference(branchRef)
		err = w.Checkout(&git.CheckoutOptions{
			Force:  true,
			Branch: branchRef.Name(),
		})
		fmt.Println(err)
	}
}

func (configService *ConfigService) GitPull() (bool, error) {

	if configService.GitUrl != "" {

		configService.Logger.Info("Opening git repo at path: " + configService.Dir)
		repo, err := git.PlainOpen(configService.Dir)

		w, err := repo.Worktree()

		branchName := plumbing.NewBranchReferenceName(configService.GitBranch)

		err = w.Pull(&git.PullOptions{RemoteName: "origin",
			ReferenceName: branchName,
			Auth: &http.BasicAuth{
				Username: configService.GitUsername,
				Password: configService.GitToken,
			}})

		if err != nil && fmt.Sprint(err) != "already up-to-date" {
			message := "Git Pull Failed on branch (" + configService.GitBranch + "):" + fmt.Sprint(err)
			configService.Logger.Error(message)
			return false, errors.New(message)
		}
	}
	return true, nil
}

func (configService ConfigService) Get(environment string, instance string, profile string, path string) (interface{}, error) {

	configService.Logger.Info(fmt.Sprintf("Fetching config for environment:%s, instance:%s, profile:%s , path: %s", fmt.Sprint(environment), fmt.Sprint(instance), fmt.Sprint(profile), fmt.Sprint(path)))

	os.Setenv("GO_ENV", environment)
	os.Setenv("GO_APP_INSTANCE", instance)
	os.Setenv("GO_PROFILES_ACTIVE", profile)

	_, err := configService.GitPull()

	if err != nil {
		return "", err
	}
	conf := config.GetServiceConfigFromDir(configService.Dir)
	var result interface{}
	result, _ = conf.Get("")
	if path != "" {
		result, _ = conf.Get(path)
	}

	configService.Logger.Info(fmt.Sprintf("Sending config for environment:%s, instance:%s, profile:%s , path: %s as:", environment, instance, profile, path) + fmt.Sprint(result))

	return result, nil
}
