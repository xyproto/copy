package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const versionString = "copy 0.1"

// Check if the given list of strings has the given string
func has(sl []string, s string) bool {
	for _, e := range sl {
		if e == s {
			return true
		}
	}
	return false
}

// Check if the given path is a directory
func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	return fi.IsDir(), err
}

func main() {
	// Start out by reading in the configuration, if available
	conf, err := ReadConfig(configFile)
	if err != nil {
		conf = &Configuration{}
		err = WriteConfig(configFile, conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Parse flags, if any
	var (
		unparsedRemoteHostAlias, removeHost string
		listHosts, verbose, version, help   bool
	)
	flag.StringVar(&unparsedRemoteHostAlias, "a", "", "Add a host alias")
	flag.StringVar(&removeHost, "r", "", "Remove a host alias")
	flag.BoolVar(&listHosts, "l", false, "List host aliases")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&version, "version", false, "Version info")
	flag.BoolVar(&help, "help", false, "Help output")
	flag.Parse()

	if version {
		fmt.Println(versionString)
		os.Exit(1)
	}

	if help {
		fmt.Println(versionString)
		fmt.Println("TO IMPLEMENT: USAGE INFORMATION") // TODO
	}

	// Perform actions, based on the given flags and arguments

	if listHosts {
		for _, rh := range conf.RemoteHosts {
			// &rh instead of rh to print with .String() instead of printing the structure
			fmt.Println(&rh)
		}
		os.Exit(0)
	}

	if unparsedRemoteHostAlias != "" {
		newAE, err := NewRemoteHostAlias(unparsedRemoteHostAlias)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if ahas(conf.RemoteHosts, newAE) {
			var newRemoteHosts []RemoteHostAlias
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
		var filteredRemoteHosts []RemoteHostAlias
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

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "At least one arguments is required:")
		fmt.Fprintln(os.Stderr, "copy FROM")
		os.Exit(1)
	}

	to := ""
	froms := args
	if len(args) > 1 {
		froms = args[0 : len(args)-1]
		to = args[len(args)-1]
	}

	// Set up the source to copy files from
	so := &Source{}
	var fromFiles []string
	for _, globExpression := range froms {
		globResult, err := filepath.Glob(globExpression)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fromFiles = append(fromFiles, globResult...)
	}
	if len(fromFiles) == 0 {
		// Assume that the rest of the arguments are remote files
		for _, fromExpression := range froms {
			rf, err := NewRemoteFile(fromExpression, conf.RemoteHosts)
			if err != nil {
				fmt.Fprintln(os.Stderr, fromExpression+" does not exist.")
				os.Exit(1)
			}
			so.RemoteFiles = append(so.RemoteFiles, rf)
		}
		if len(so.RemoteFiles) == 0 {
			// TODO: This never happens?
			fmt.Fprintln(os.Stderr, "None of these are valid remote file expressions:")
			for _, fromExpression := range froms {
				fmt.Fprintf(os.Stderr, "\t%s\n", fromExpression)
			}
			os.Exit(1)
		}
	} else {
		so.Files = fromFiles
	}

	// Set up the target to copy files to
	ta := &Target{}
	if to == "" {
		var err error
		if ta.Directory, err = os.Getwd(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		if ok, err := IsDir(to); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else if ok {
			// A directory
			ta.Directory = to
		} else {
			// Not a directory, check if it's a remote host alias
			ae, err := aget(conf.RemoteHosts, to)
			if err != nil {
				fmt.Fprintln(os.Stderr, to+" is not a file and not a remote host alias")
				os.Exit(1)
			}
			ta.RemoteHost = ae
		}
	}

	if verbose {
		fmt.Println(so)
		fmt.Println(ta)
	}

	pk := NewSimpleProgressKeeper()

	// Copy the files and report the progress
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := so.Copy(wg, ta, pk)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	wg.Add(1)
	//go Write(wg, pk, 100*time.Millisecond)
	go Spin(wg, pk, 100*time.Millisecond)
	wg.Wait()
}
