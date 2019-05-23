package main

import (
	"github.com/jykwon91/project/task"
	"github.com/jykwon91/project/rest"
)

func main() {
	task.InitTasks()
	rest.InitRestClient()
}
