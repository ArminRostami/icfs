package main

import (
	"fmt"
	"icfs/adapters/shell"
	"icfs/app"

	"github.com/pkg/errors"
)

func run() error {
	cancel, ips, err := app.NewIpfsService()
	defer cancel()

	if err != nil {
		return errors.Wrap(err, "failed to start new ipfs service")
	}
	ish := shell.New(ips)
	ish.Init()
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("run failed: %+v", err)
	}
}
