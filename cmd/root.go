package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"hujm.net/terminal_note/note"
)

/*
* @CreateTime: 2021/1/12 20:59
* @Author: hujiaming
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(0)
	}
}

// GetUserName ...
func GetUserName() string {
	if userName == "" {
		u, _ := user.Current()
		userName = u.Username
	}
	return userName
}

func CheckNotebook(userName string) {
	if !note.CheckDbFileExist(userName) {
		fmt.Printf("No existing notebook found for %s. Use `create` to create one.\n", userName)
		os.Exit(0)
	}
}

func GetInputPassword() string {
	password := ""
	prompt := &survey.Password{Message: "Please type your password:"}
	err := survey.AskOne(prompt, &password, nil)
	if err != nil {
		log.Panic(err)
	}
	return password
}
