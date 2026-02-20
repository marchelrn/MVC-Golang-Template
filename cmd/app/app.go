package app

import (
	"mini_jira/config"
	"mini_jira/internal/server"
)

func RunApp() {
	config.Load()
	server.Run()
}
