// Package shell is an interactive shell as a frontend for testing services
package shell

import (
	"icfs/app"

	"gopkg.in/abiosoft/ishell.v2"
)

type Shell struct {
	service *app.IpfsService
	i       *ishell.Shell
}

func New(s *app.IpfsService) *Shell {
	return &Shell{service: s}
}

func (sh *Shell) Init() {
	sh.i = ishell.New()
	sh.i.Println("Shell started")
	sh.i.AddCmd(&ishell.Cmd{Name: "add", Help: "add a file to ipfs", Func: func(c *ishell.Context) {
		filename := c.Args[0]
		filepath := "./" + filename
		cid, err := sh.service.AddFile(filepath)
		if err != nil {
			c.Err(err)
			return
		}
		c.Printf("filed added with cid: %s\n", cid)
	}})

	sh.i.AddCmd(&ishell.Cmd{Name: "get", Help: "get a file from ipfs", Func: func(c *ishell.Context) {
		cidStr := c.Args[0]
		err := sh.service.GetFile(cidStr)
		if err != nil {
			c.Err(err)
			return
		}
		c.Printf("file written successfully\n")
	}})

	sh.i.AddCmd(&ishell.Cmd{Name: "connect", Help: "connect to a peer", Func: func(c *ishell.Context) {
		addrStr := c.Args[0]
		err := sh.service.Connect(addrStr)
		if err != nil {
			c.Err(err)
			return
		}
		c.Printf("connected to %s\n", addrStr)
	}})
	sh.i.Run()
}
