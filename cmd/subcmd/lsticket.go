package subcmd

import (
	"fmt"

	irodsclient_fs "github.com/cyverse/go-irodsclient/fs"
	"github.com/cyverse/gocommands/commons"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var lsticketCmd = &cobra.Command{
	Use:   "lsticket [ticket_string1] [ticket_string2] ...",
	Short: "List tickets for the user",
	Long:  `This lists tickets for the user.`,
	RunE:  processLsticketCommand,
}

func AddLsticketCommand(rootCmd *cobra.Command) {
	// attach common flags
	commons.SetCommonFlags(lsticketCmd)

	rootCmd.AddCommand(lsticketCmd)
}

func processLsticketCommand(command *cobra.Command, args []string) error {
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

	// Create a file system
	account := commons.GetAccount()
	filesystem, err := commons.GetIRODSFSClient(account)
	if err != nil {
		return err
	}

	defer filesystem.Release()

	if len(args) == 0 {
		return fmt.Errorf("not enough input arguments")
	}

	for _, ticket := range args {
		err = getTicket(filesystem, ticket)
		if err != nil {
			return err
		}
	}

	return nil
}

func getTicket(filesystem *irodsclient_fs.FileSystem, ticket string) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "getTicket",
	})

	logger.Debugf("get ticket: %s", ticket)

	ticketInfo, err := filesystem.GetTicketForAnonymousAccess(ticket)
	if err != nil {
		return err
	}

	fmt.Printf("[%s]\n", ticketInfo.Name)
	fmt.Printf("  id: %d\n", ticketInfo.ID)
	fmt.Printf("  string: %s\n", ticketInfo.Name)

	if ticketInfo.ExpireTime.IsZero() {
		fmt.Print("  expiry time: none\n")
	} else {
		fmt.Printf("  expiry time: %s\n", ticketInfo.ExpireTime)
	}

	fmt.Printf("  ticket type: %s\n", ticketInfo.Type)
	fmt.Printf("  collection name: %s\n", ticketInfo.Path)

	return nil
}
