package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

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

func main() {
	var addHost, removeHost string
	flag.StringVar(&addHost, "a", "", "Add a host alias")
	flag.StringVar(&removeHost, "r", "", "Remove a host alias")
	flag.Parse()

	if addHost != "" {
		if !(strings.Contains(addHost, "=") && strings.Contains(addHost, "@") && strings.Contains(addHost, ":")) {
			fmt.Println("INVALID ALIAS EXPRESSION", addHost)
			os.Exit(1)
		}
		alias, username, hostname, port := interpret(addHost)
		fmt.Printf("SAVED %s@%s:%s as %s\n", username, hostname, port, alias)
		os.Exit(0)
	}

	if removeHost != "" {
		fmt.Printf("REMOVED %s\n", removeHost)
		os.Exit(0)
	}

	fmt.Println(flag.Args())
	if len(flag.Args()) < 2 {
		fmt.Println("TWO ARGUMENTS ARE REQUIRED")
		os.Exit(1)
	}
	from := flag.Args()[0]
	to := flag.Args()[1]

	fmt.Println("COYP", from, to)
}
