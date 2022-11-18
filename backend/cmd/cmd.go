package main

import "ctigroupjsc.com/phuocnn/gps-management/app/cli"

func RunCmd() {
	cli.NewCommander()
	commander := cli.NewCommander()
	commander.AddHandler("admin",
		&cli.AccountCommander{
			UserRepository:    container.UserRepository,
			AccountRepository: container.AccountRepository,
		},
	)
	commander.Handle()
}
