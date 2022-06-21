package repository

import "github.com/fkcs/gateway/internal/interfaces/command"

type RegisterRepoInterface interface {
	GetServers() ([]string, error)
	AddServer(info *command.ServerInfo) error
	DelServer(info *command.ServerInfo) error
}
