package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch user",
	Long:  `Switch the active user`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			user := args[0]
			for _, profile := range configuration.Profiles {
				// Look for the user by email
				if profile.Email == user {
					home, err := homedir.Dir()
					checkErrors(err)

					hostPrivateKey := home + "/.ssh/id_rsa"
					hostPublicKey := home + "/.ssh/id_rsa.pub"
					userPrivateKey := home + "/.ssh/" + profile.Key
					userPublicKey := home + "/.ssh/" + profile.Key + ".pub"

					fmt.Printf("Moving ssh key ~/.ssh/%s_rsa, ~/.ssh/%s_rsa.pub\n", profile.Key, profile.Key)

					copyFile(userPrivateKey, hostPrivateKey)
					copyFile(userPublicKey, hostPublicKey)

					fmt.Printf("Setting Git name to %s and email to %s\n", profile.Name, profile.Email)

					cmd := "git"
					argsName := []string{"config", "--global", "user.name", profile.Name}
					argsEmail := []string{"config", "--global", "user.email", profile.Email}
					// Set Git Config name
					if err := exec.Command(cmd, argsName...).Run(); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					// Set Git Config Email
					if err := exec.Command(cmd, argsEmail...).Run(); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					os.Exit(0)
				}
			}

			fmt.Printf("No profile found for %s.\n", user)
			os.Exit(1)
		} else { // If a user is not provided
			fmt.Println("Error: you must provide a user")
			fmt.Println("Usage: switchr switch [user]")
			os.Exit(2)
		}
	},
}

func init() {
	RootCmd.AddCommand(switchCmd)
}

func copyFile(source string, destination string) {
	srcFile, err := os.Open(source)
	checkErrors(err)
	defer srcFile.Close()

	destFile, err := os.Create(destination) // Create if the file doesn't exist
	checkErrors(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	checkErrors(err)

	err = destFile.Sync()
	checkErrors(err)
}

func checkErrors(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
