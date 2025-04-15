package main

import "todo_list_go/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
