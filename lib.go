package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"path/filepath"
)

func (o *options) simplified_plat() platform {
	switch o.plat {
		case darwin,macos,maccatalyst,ios,ios_simulator,watchos,watchos_simulator,tvos,tvos_simulator:
			return darwin
		default:
			return o.plat
	}
}

func expand(path string) string {

	if strings.HasPrefix(path,"~/") {
		u,_ := user.Current()
		home := u.HomeDir
		return filepath.Join(home,path[2:])
	} else { return path }

}

func is_exist(path string) bool {
	_,err := os.Lstat(path)
	return err==nil
}

func warn(message string) {
	fmt.Fprintln(os.Stderr,message)
}

func err(message string) {
	fmt.Fprintln(os.Stderr,message)
	os.Exit(1)
}