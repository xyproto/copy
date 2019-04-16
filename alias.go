package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RemoteHostAlias struct {
	Alias    string `alias`
	Username string `username`
	Hostname string `hostname`
	Port     int    `port`
}

type RemoteFile struct {
	Username string `username`
	Hostname string `hostname`
	Port     int    `port`
	Path     string
}

// interpret an expression on the form: username@hostname:port_path
func NewRemoteFile(unparsedRemoteFile string, hostAliases []RemoteHostAlias) (*RemoteFile, error) {
	if (!strings.Contains(unparsedRemoteFile, "@") && !strings.Contains(unparsedRemoteFile, ":")) && strings.Contains(unparsedRemoteFile, "_") {
		fields := strings.SplitN(unparsedRemoteFile, "_", 2)
		if fields[0] != "" {
			// Check if the first part of the expression is a remote host alias
			for _, rh := range hostAliases {
				if rh.Alias == fields[0] {
					// Yes, return that
					rf := &RemoteFile{}
					rf.Username = rh.Username
					rf.Hostname = rh.Hostname
					rf.Port = rh.Port
					rf.Path = fields[1]
					return rf, nil
				}
			}
		}
	}

	// The first part is not an existing remote host alias

	usernameHostPortPath := unparsedRemoteFile
	var (
		username     string
		hostPortPath string
	)
	if !strings.Contains(usernameHostPortPath, "@") {
		username = os.Getenv("LOGNAME")
		hostPortPath = usernameHostPortPath
	} else {
		fields := strings.SplitN(usernameHostPortPath, "@", 2)
		username = fields[0]
		hostPortPath = fields[1]
	}
	// TODO: Also support file expressions on the form hostname:filename
	if !strings.Contains(hostPortPath, "_") {
		return nil, errors.New("A remote file must have a path: username@hostname:port_path")
	}
	fields := strings.SplitN(hostPortPath, "_", 2)
	hostPort := fields[0]
	path := fields[1]
	var (
		hostname string
		port     int
		err      error
	)
	if !strings.Contains(hostPort, ":") {
		hostname = hostPort
		port = 22
	} else {
		fields = strings.SplitN(hostPort, ":", 2)
		hostname = fields[0]
		port, err = strconv.Atoi(fields[1])
		if err != nil {
			// The port is not a number, assume that the path follows the colon instead
			port = 22
			if path == "" {
				path = fields[1]
			}
		}
	}
	if hostname == "" {
		return nil, errors.New("A remote file must have a hostname: username@hostname:port_path")
	}
	rf := &RemoteFile{}
	rf.Username = username
	rf.Hostname = hostname
	rf.Port = port
	rf.Path = path
	return rf, nil
}

// interpret an alias expression on the form: alias=username@hostname:port
func NewRemoteHostAlias(unparsedRemoteHostAlias string) (*RemoteHostAlias, error) {
	if !strings.Contains(unparsedRemoteHostAlias, "=") {
		return nil, errors.New("Alias expressions must contain an equal sign!")
	}
	fields := strings.SplitN(unparsedRemoteHostAlias, "=", 2)
	alias := fields[0]
	usernameHostPort := fields[1]
	var (
		username string
		hostPort string
	)
	if !strings.Contains(usernameHostPort, "@") {
		username = os.Getenv("LOGNAME")
		hostPort = usernameHostPort
	} else {
		fields = strings.SplitN(usernameHostPort, "@", 2)
		username = fields[0]
		hostPort = fields[1]
	}
	var (
		hostname string
		port     int
		err      error
	)
	if !strings.Contains(hostPort, ":") {
		hostname = hostPort
		port = 22
	} else {
		fields = strings.SplitN(hostPort, ":", 2)
		hostname = fields[0]
		port, err = strconv.Atoi(fields[1])
		if err != nil {
			port = 22
		}
	}
	ae := &RemoteHostAlias{}
	ae.Alias = alias
	ae.Username = username
	ae.Hostname = hostname
	ae.Port = port
	return ae, nil
}

// Check if the given list of remote host aliases contains the given alias
func ahas(ras []RemoteHostAlias, ra *RemoteHostAlias) bool {
	for _, e := range ras {
		if e.Alias == ra.Alias {
			return true
		}
	}
	return false
}

// Return the RemoteHostAlias for the given alias name, if it exists
func aget(ras []RemoteHostAlias, alias string) (*RemoteHostAlias, error) {
	for _, e := range ras {
		if e.Alias == alias {
			return &e, nil
		}
	}
	return nil, errors.New("Could not find remote host alias: " + alias)
}

func (ae *RemoteHostAlias) String() string {
	s := ""
	if ae.Alias != "" {
		s += ae.Alias + "="
	}
	if ae.Username != "" {
		s += ae.Username + "@"
	}
	if ae.Hostname != "" {
		s += ae.Hostname
	}
	if ae.Port != 0 {
		s += fmt.Sprintf(":%d", ae.Port)
	}
	return s
}

func (ae *RemoteFile) String() string {
	s := ""
	if ae.Username != "" {
		s += ae.Username + "@"
	}
	if ae.Hostname != "" {
		s += ae.Hostname
	}
	if ae.Port != 0 {
		s += fmt.Sprintf(":%d", ae.Port)
	}
	if ae.Path != "" {
		s += "_" + ae.Path
	}
	return s
}
