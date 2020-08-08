package main

import (
	"fmt"
	"icfs/adapters/shell"
	"icfs/app"

	"github.com/pkg/errors"
	"gopkg.in/abiosoft/ishell.v2"
)

func run() error {
	cancel, ipfsService, err := app.NewService()
	defer cancel()
	if err != nil {
		return errors.Wrap(err, "failed to create new ipfs service")
	}

	sh := &shell.Shell{Ish: ishell.New()}

	if !ipfsService.RepoExists() {
		sh.Ish.Println("enter bootstrap address")
		bootStr, err := sh.Ish.ReadLineErr()
		if err != nil {
			return errors.Wrap(err, "failed to readline from console")
		}
		err = ipfsService.SetupRepo(bootStr)
		if err != nil {
			return errors.Wrap(err, "failed to setup repo on default path")
		}
	}

	if err = ipfsService.StartService(); err != nil {
		return errors.Wrap(err, "failed to start ipfs service")
	}

	sh.Init(ipfsService)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %+v", err)
	}
}
