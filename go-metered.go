package gometered

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/Wifx/gonetworkmanager"
)

/*
Quirks:
Does not preserve errors (see IsMeteredFull(), returns false apon error)
Only works on systemd based machines running network manager (PRs welcome)
- On systemd, it checks ALL current interfaces if they're metered and returns true if any are
*/
func IsMetered() bool {
	x, _ := IsMeteredFull()
	return x
}

//IsMeteredFull takes the same input as IsMetered(), however, it returns errors,
//so if you don't want to use the function inline and/or want to troubleshoot, use this.
func IsMeteredFull() (bool, error) {

	osType, err := checkOS()
	if err != nil {
		return false, err
	}

	switch osType {
	case "systemd":

		/* Create new instance of gonetworkmanager */
		nm, err := gonetworkmanager.NewNetworkManager()
		if err != nil {
			return false, err
		}

		resp, err := nm.GetPropertyMetered()
		if err != nil {
			return false, err
		}

		if resp == 1 || resp == 3 {
			return true, nil
		}

	}

	return false, fmt.Errorf("unknownos")

}

func checkOS() (string, error) {
	// systemd check

	// move to os.ReadDir once golang 1.16 is adopted by go-ipfs
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		return "", err
	}

	for _, f := range dirs {
		if f.IsDir() {
			_, err := strconv.Atoi(f.Name())
			if err == nil {
				exe, _ := os.Readlink(f.Name())
				if exe == "/usr/lib/systemd/systemd" {
					return "systemd", nil

				}
			}
		}
	}

	//fallback
	return "unknown", nil
}
