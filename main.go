package main

import (
	"github.com/chnyangzhen/kago-fly/pkg/cmd"
	"github.com/chnyangzhen/kago-fly/pkg/constant"
	"strconv"
	"time"
)
import _ "github.com/chnyangzhen/kago-fly/src/user"

func main() {
	app := cmd.App{
		Name:        constant.AppName,
		Version:     "1.0",
		Copyright:   "(c) " + strconv.Itoa(time.Now().Year()),
		Description: "",
		Banner:      "asd",
	}

	cmd.Go(app)
}
