package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/jemmyh/tnote/note"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

/*
* @CreateTime: 2021/1/12 20:59
* @Author: JemmyHu <hujm20151021@gmail.com>
* @Description:
 */

var (
	userName string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&userName, "owner", "o", "", "owner of notebook")
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  `Terminal Notebook is CLI App, which is implemented by Golang.`,
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(0)
	}
}

// GetUserName returns the current logined user, the same as `whoami`
func GetUserName() string {
	if userName == "" {
		u, _ := user.Current()
		userName = u.Username
	}
	return userName
}

// CheckNotebook make app exit if user's notebook does not exist.
func CheckNotebook(userName string) {
	if !note.CheckDbFileExist(userName) {
		fmt.Printf("No existing notebook found for %s. Use `create` to create one.\n", userName)
		os.Exit(0)
	}
}

// GetInputPassword gets password from stdin.
func GetInputPassword() string {
	password := ""
	prompt := &survey.Password{Message: "Please type your password:"}
	err := survey.AskOne(prompt, &password, nil)
	if err != nil {
		log.Panic(err)
	}
	return password
}
