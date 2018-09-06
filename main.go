package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type AliasExpression struct {
	Alias    string `alias`
	Username string `username`
	Hostname string `hostname`
	Port     int    `port`
}

var configFile = filepath.Join(os.Getenv("HOME"), ".config/copy.conf")

type Configuration struct {
	RemoteHosts []AliasExpression
}

// interpret an alias expression on the form: alias=username@hostname:port
func NewAliasExpression(unparsedAliasExpression string) (*AliasExpression, error) {
	if !strings.Contains(unparsedAliasExpression, "=") {
		return nil, errors.New("Alias expressions must contain an equal sign!")
	}
	fields := strings.SplitN(unparsedAliasExpression, "=", 2)
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
	ae := &AliasExpression{}
	ae.Alias = alias
	ae.Username = username
	ae.Hostname = hostname
	ae.Port = port
	return ae, nil
}

func ReadConfig(filename string) (*Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func WriteConfig(filename string, config *Configuration) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	return nil
}

// Check if the given list of strings has the given string
func has(sl []string, s string) bool {
	for _, e := range sl {
		if e == s {
			return true
		}
	}
	return false
}

// Check if the given list of aliasExpressions contains the given alias
func ahas(aes []AliasExpression, ae *AliasExpression) bool {
	for _, e := range aes {
		if e.Alias == ae.Alias {
			return true
		}
	}
	return false
}

func main() {
	conf, err := ReadConfig(configFile)
	if err != nil {
		conf = &Configuration{}
		err = WriteConfig(configFile, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	var (
		unparsedAliasExpression, removeHost string
		listHosts                           bool
	)
	flag.StringVar(&unparsedAliasExpression, "a", "", "Add a host alias")
	flag.StringVar(&removeHost, "r", "", "Remove a host alias")
	flag.BoolVar(&listHosts, "l", false, "List host aliases")
	flag.Parse()

	if listHosts {
		fmt.Println(conf.RemoteHosts)
		os.Exit(0)
	}

	if unparsedAliasExpression != "" {
		newAE, err := NewAliasExpression(unparsedAliasExpression)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if ahas(conf.RemoteHosts, newAE) {
			var newRemoteHosts []AliasExpression
			for _, aliasExpression := range conf.RemoteHosts {
				if aliasExpression.Alias == newAE.Alias {
					// Already exists
					newRemoteHosts = append(newRemoteHosts, *newAE)
				} else {
					newRemoteHosts = append(newRemoteHosts, aliasExpression)
				}
			}
			conf.RemoteHosts = newRemoteHosts
		} else {
			conf.RemoteHosts = append(conf.RemoteHosts, *newAE)
		}
		err = WriteConfig(configFile, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		//fmt.Printf("SAVED %s@%s:%d as %s\n", newAE.Username, newAE.Hostname, newAE.Port, newAE.Alias)
		os.Exit(0)
	}

	if removeHost != "" {
		var filteredRemoteHosts []AliasExpression
		for _, aliasExpression := range conf.RemoteHosts {
			if aliasExpression.Alias == removeHost {
				//fmt.Printf("REMOVED %s\n", removeHost)
				// Skip
				continue
			}
			filteredRemoteHosts = append(filteredRemoteHosts, aliasExpression)
		}
		conf.RemoteHosts = filteredRemoteHosts
		err = WriteConfig(configFile, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "Two arguments are required:")
		fmt.Fprintln(os.Stderr, "copy FROM TO")
		os.Exit(1)
	}

	from := flag.Args()[0]
	to := flag.Args()[1]

	// First assume that from and to are local files.
	// Look into aliases only if local files of the same names are not found.

	a, b := os.Stat(from)
	c, d := os.Stat(to)

	fmt.Println(a, b, c, d)

	fmt.Println("TO IMPLEMENT: COPY", from, to)
}
