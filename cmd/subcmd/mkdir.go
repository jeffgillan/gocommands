package subcmd

import (
	"fmt"
	"strconv"

	irodsclient_fs "github.com/cyverse/go-irodsclient/fs"
	irodsclient_irodsfs "github.com/cyverse/go-irodsclient/irods/fs"
	"github.com/cyverse/gocommands/commons"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var mkdirCmd = &cobra.Command{
	Use:   "mkdir [collection1] [collection2] ...",
	Short: "Make iRODS collections",
	Long:  `This makes iRODS collections.`,
	RunE:  processMkdirCommand,
}

func AddMkdirCommand(rootCmd *cobra.Command) {
	// attach common flags
	commons.SetCommonFlags(mkdirCmd)
	mkdirCmd.Flags().Bool("parents", false, "Make parent collections")

	rootCmd.AddCommand(mkdirCmd)
}

func processMkdirCommand(command *cobra.Command, args []string) error {
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

	parent := false
	parentFlag := command.Flags().Lookup("parent")
	if parentFlag != nil {
		parent, err = strconv.ParseBool(parentFlag.Value.String())
		if err != nil {
			parent = false
		}
	}

	// Create a connection
	account := commons.GetAccount()
	filesystem, err := commons.GetIRODSFSClient(account)
	if err != nil {
		return err
	}

	defer filesystem.Release()

	if len(args) == 0 {
		return fmt.Errorf("not enough input arguments")
	}

	for _, targetPath := range args {
		err = makeOne(filesystem, targetPath, parent)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeOne(fs *irodsclient_fs.FileSystem, targetPath string, parent bool) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "makeOne",
	})

	cwd := commons.GetCWD()
	home := commons.GetHomeDir()
	zone := commons.GetZone()
	targetPath = commons.MakeIRODSPath(cwd, home, zone, targetPath)

	connection, err := fs.GetConnection()
	if err != nil {
		return err
	}
	defer fs.ReturnConnection(connection)

	logger.Debugf("making a collection %s", targetPath)

	err = irodsclient_irodsfs.CreateCollection(connection, targetPath, parent)
	if err != nil {
		return err
	}
	return nil
}
