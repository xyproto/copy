package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"encoding/json"
)

type Configuration struct {
	AliasExpressions []string
}

// interpret an alias expression on the form: alias=username@hostname:port
func interpret(aliasExpression string) (string, string, string, string) {
	fields := strings.SplitN(aliasExpression, "=", 2)
	alias := fields[0]
	usernameHostPort := fields[1]
	fields = strings.SplitN(usernameHostPort, "@", 2)
	username := fields[0]
	hostPort := fields[1]
	fields = strings.SplitN(hostPort, ":", 2)
	hostname := fields[0]
	port := fields[1]
	return alias, username, hostname, port
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
	file, err := os.OpenFile(filename, os.O_CREATE | os.O_WRONLY, 0644)
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

func main() {
	conf, err := ReadConfig("~/.config/copy.conf")
	if err != nil {
		err = WriteConfig("~/.config/copy.conf", &Configuration{})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	var addHost, removeHost string
	flag.StringVar(&addHost, "a", "", "Add a host alias")
	flag.StringVar(&removeHost, "r", "", "Remove a host alias")
	flag.Parse()

	if addHost != "" {
		if !(strings.Contains(addHost, "=") && strings.Contains(addHost, "@") && strings.Contains(addHost, ":")) {
			fmt.Fprintln(os.Stderr, "INVALID ALIAS EXPRESSION", addHost)
			os.Exit(1)
		}
		alias, username, hostname, port := interpret(addHost)
		conf.AliasExpressions = append(conf.AliasExpressions, addHost)
		err = WriteConfig("~/.config/copy.conf", conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("SAVED %s@%s:%s as %s\n", username, hostname, port, alias)
		os.Exit(0)
	}

	if removeHost != "" {
		fmt.Printf("REMOVED %s\n", removeHost)
		os.Exit(0)
	}

	fmt.Println(flag.Args())
	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "TWO ARGUMENTS ARE REQUIRED")
		os.Exit(1)
	}
	from := flag.Args()[0]
	to := flag.Args()[1]

	fmt.Println("COPY", from, to)
}
