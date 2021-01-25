package main

type currentHost interface {
	copyTeamFiles() error
	install() error
}
