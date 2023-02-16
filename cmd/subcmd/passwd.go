package subcmd

import (
	"fmt"
	"syscall"

	irodsclient_fs "github.com/cyverse/go-irodsclient/fs"
	irodsclient_irodsfs "github.com/cyverse/go-irodsclient/irods/fs"
	"github.com/cyverse/gocommands/commons"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var passwdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "Change iRODS user password",
	Long:  `This changes iRODS user password.`,
	RunE:  processPasswdCommand,
}

func AddPasswdCommand(rootCmd *cobra.Command) {
	// attach common flags
	commons.SetCommonFlags(passwdCmd)

	rootCmd.AddCommand(passwdCmd)
}

func processPasswdCommand(command *cobra.Command, args []string) error {
	cont, err := commons.ProcessCommonFlags(command)
	if err != nil {
		return err
	}

	if !cont {
		return nil
	}

	// handle local flags
	_, err = commons.InputMissingFields()
	if err != nil {
		return err
	}

	// Create a connection
	account := commons.GetAccount()
	filesystem, err := commons.GetIRODSFSClient(account)
	if err != nil {
		return err
	}

	err = changePassword(filesystem)
	if err != nil {
		return err
	}
	return nil
}

func changePassword(fs *irodsclient_fs.FileSystem) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "changePassword",
	})

	account := commons.GetAccount()

	connection, err := fs.GetConnection()
	if err != nil {
		return err
	}
	defer fs.ReturnConnection(connection)

	logger.Debugf("changing password for user %s", account.ClientUser)

	pass := false
	for i := 0; i < 3; i++ {
		fmt.Print("Current iRODS Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		fmt.Print("\n")
		currentPassword := string(bytePassword)

		if len(currentPassword) == 0 {
			fmt.Println("Please provide password")
			fmt.Println("")
			continue
		}

		if currentPassword == account.Password {
			pass = true
			break
		}

		fmt.Println("Wrong password")
		fmt.Println("")
	}

	if !pass {
		return fmt.Errorf("password mismatched")
	}

	pass = false
	newPassword := ""
	for i := 0; i < 3; i++ {
		fmt.Print("New iRODS Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		fmt.Print("\n")
		newPassword = string(bytePassword)

		if len(newPassword) == 0 {
			fmt.Println("Please provide password")
			fmt.Println("")
			continue
		}

		if newPassword != account.Password {
			pass = true
			break
		}

		fmt.Println("Please provide new password")
		fmt.Println("")
	}

	if !pass {
		return fmt.Errorf("invalid password provided")
	}

	newPasswordConfirm := ""
	fmt.Print("Confirm New iRODS Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	fmt.Print("\n")
	newPasswordConfirm = string(bytePassword)

	if newPassword != newPasswordConfirm {
		return fmt.Errorf("password mismatched")
	}

	err = irodsclient_irodsfs.ChangeUserPassword(connection, account.ClientUser, account.ClientZone, newPassword)
	if err != nil {
		return err
	}
	return nil

}
