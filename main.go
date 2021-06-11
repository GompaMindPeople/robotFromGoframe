package main

import (
	_ "robotFromGoframe/boot"
	_ "robotFromGoframe/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
