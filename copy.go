package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Source is a list of local and remote files
type Source struct {
	Files       []string
	RemoteFiles []*RemoteFile
}

// Target is either a directory or a remote host directory
type Target struct {
	Directory  string
	RemoteHost *RemoteHostAlias
}

func (so *Source) String() string {
	s := "FROM:\n"
	if len(so.Files) > 0 {
		for _, fn := range so.Files {
			s += "\t" + fn + "\n"
		}
	}
	if len(so.RemoteFiles) > 0 {
		for _, rf := range so.RemoteFiles {
			s += "\t" + rf.String() + "\n"
		}
	}
	if strings.HasSuffix(s, "\n") {
		s = s[:len(s)-1]
	}
	return s
}

func (ta *Target) String() string {
	s := "TO:\n"
	if ta.Directory != "" {
		s += "\t" + ta.Directory + "\n"
	}
	if ta.RemoteHost != nil {
		s += "\t" + ta.RemoteHost.String() + "\n"
	}
	if strings.HasSuffix(s, "\n") {
		s = s[:len(s)-1]
	}
	return s
}

// TODO: Rewrite this and think it through so that multiple goroutines can both copy files and also report progress!
// Copy will copy the sources to a target, while reporting the progress
func (so *Source) Copy(wg *sync.WaitGroup, ta *Target, pk ProgressKeeper) error {
	return nil
	var err error
	curdir, err := os.Getwd()
	if err != nil {
		return err
	}
	curdir = filepath.Clean(curdir)
	progress := 0.0
	elements := len(so.Files) + len(so.RemoteFiles)
	step := 1.0 / float64(elements)
	for _, localFile := range so.Files {
		pk.Register(localFile, progress, "done")
		filename := filepath.Base(localFile)
		localDir := filepath.Dir(localFile)
		from := filepath.Clean(localDir)
		to := ""
		if ta.Directory == "." {
			to = curdir
		} else if ta.Directory != "" {
			to = filepath.Clean(ta.Directory)
		} else {
			// TODO
			to = ta.RemoteHost.String() + "(REMOTE)"
		}

		if from == to {
			pk.SetMessage(fmt.Sprintf("Skipping %s, already in %s", filename, from))
		} else {
			if from == "." {
				pk.SetMessage(fmt.Sprintf("Copying %s to %s...", filename, to))
			} else {
				pk.SetMessage(fmt.Sprintf("Copying %s from %s to %s...", filename, from, to))
			}
			time.Sleep(1 * time.Second)
		}
		progress += step
		if progress < 1.0 {
			pk.Register(localFile, progress, "done")
		}
	}
	for _, rf := range so.RemoteFiles {
		pk.Register(rf.String(), progress, "done")
		filename := filepath.Base(rf.Path)
		// TODO
		from := rf.Path + "(REMOTE)"
		to := ""
		if ta.Directory != "" {
			localDir := filepath.Dir(ta.Directory)
			to = filepath.Clean(localDir)
		} else {
			to = ta.RemoteHost.String()
		}
		if from == "." {
			pk.SetMessage(fmt.Sprintf("Copying %s to %s...", filename, to))
		} else {
			pk.SetMessage(fmt.Sprintf("Copying %s from %s to %s...", filename, from, to))
		}
		time.Sleep(1 * time.Second)
		progress += step
		if progress < 1.0 {
			pk.Register(rf.String(), progress, "done")
		}
	}
	//pk.Register(1.0)
	wg.Done()
	return nil
}
