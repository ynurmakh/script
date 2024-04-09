package gitconfigs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	gouseful "github.com/p0n41k/goUseful"
)

func GitConfigs() {
	/*

		parse gitZeroCommit
			replace it

			parce .git
				rename url
				create gitzerocommit


	*/

	ok := parceGitZeroCommit()
	if ok == false {
		parceGitOriginal()
		ok = parceGitZeroCommit()
		if ok == false {
			log.Fatalln("Can`t parse .git folder")
		}
	}

	gouseful.Bash("git init")
}

func parceGitZeroCommit() bool {
	zeroCommit := ".gitZeroCommit"
	_, err := os.ReadDir(zeroCommit)
	if err != nil {
		return false
	}

	_, err = gouseful.Bash([]string{
		"rm -r -f .git",
		"cp -r -f .gitZeroCommit .git",
	})
	if err != nil {
		log.Println("\n" + err.Error() + "\n")
	}

	return true
}

func parceGitOriginal() error {
	originalGitRepoName := ".git"
	_, err := os.ReadDir(originalGitRepoName)
	if err != nil {
		return errors.New("can`t pase .git" + err.Error())
	}

	gitsConfigName := ".git/config"
	file, err := os.ReadFile(gitsConfigName)
	if err != nil {
		return errors.New("can`t parse .git/config" + err.Error())
	}

	fileLines := strings.Split(string(file), "\n")

	for i := 0; i < len(fileLines); i++ {
		startIndex := strings.Index(fileLines[i], "url = ")
		if startIndex != -1 {

			link := fileLines[i][startIndex+6:]
			if len(link) > 3 && link[:4] == "git@" {
				_, err = gouseful.Bash("cp -r .git .gitZeroCommit")
				if err != nil {
					return errors.New("can`t create .gitZeroCommit" + err.Error())
				}
				return nil
			}

			httpDomainOwnerRepo := fileLines[i][startIndex+6:]

			domain, owner, repo := httpsToGit(httpDomainOwnerRepo)

			gitLink := "git@" + domain + ":" + owner + "/" + repo

			fileLines[i] = fileLines[i][:startIndex+6] + gitLink
			fmt.Println(fileLines[i], "<<<")
		}
	}

	forWrite := strings.Join(fileLines, "\n")

	err = os.WriteFile(gitsConfigName, []byte(forWrite), 0662)
	if err != nil {
		return errors.New("can`t write new .git/config" + err.Error())
	}

	_, err = gouseful.Bash("cp -r .git .gitZeroCommit")
	if err != nil {
		return errors.New("can`t create .gitZeroCommit" + err.Error())
	}

	return nil
}

func httpsToGit(httpDomainOwnerRepo string) (string, string, string) {
	slashIndex := strings.Index(httpDomainOwnerRepo, "/")
	domainOwnerRepo := httpDomainOwnerRepo[slashIndex+2:]

	slashIndex = strings.Index(domainOwnerRepo, "/")
	domain := domainOwnerRepo[:slashIndex]
	ownerRepo := domainOwnerRepo[slashIndex+1:]

	slashIndex = strings.Index(ownerRepo, "/")
	owner := ownerRepo[:slashIndex]
	repo := ownerRepo[slashIndex+1:]

	return domain, owner, repo
}
