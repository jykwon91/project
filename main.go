package main

import (
	"github.com/jykwon91/project/rest"
	"github.com/jykwon91/project/task"
)

func main() {
	task.InitTasks()
	rest.InitRestClient()
}
