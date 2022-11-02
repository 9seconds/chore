//go:build !(aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris)

package commands

var SignalsToRelay = []os.Signal{}
