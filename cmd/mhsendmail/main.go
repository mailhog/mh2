package main

import (
	"github.com/mailhog/mh2/cmd"
	mhsendmail "github.com/mailhog/mhsendmail/cmd"
)

func main() {
	cmd.Main(mhsendmail.Go)
}
