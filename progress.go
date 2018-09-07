// progress.go
package main

import (
	"fmt"
	"sync"
	"time"
)

type ProgressKeeper interface {
	Register(id string, progress float64, donemsg string)
	String() string
	Get() float64
	Message(id string) string
	SetMessage(msg string)
}

type SimpleProgressKeeper struct {
	progress  map[string]float64
	msgs      map[string]string
	mut       *sync.RWMutex
	currentID string
}

func NewSimpleProgressKeeper() *SimpleProgressKeeper {
	spk := &SimpleProgressKeeper{}
	spk.mut = &sync.RWMutex{}
	return spk
}

func (spk *SimpleProgressKeeper) Register(id string, progress float64, donemsg string) {
	spk.mut.Lock()
	if id != spk.currentID {
		defer spk.Complete(id, donemsg)
	}
	spk.progress[spk.currentID] = progress
	spk.mut.Unlock()
}

func (spk *SimpleProgressKeeper) String() string {
	spk.mut.RLock()
	percentage := int(100.0*spk.progress[spk.currentID] + 0.5)
	spk.mut.RUnlock()
	return fmt.Sprintf("The current progress is %d%%", percentage)
}

func (spk *SimpleProgressKeeper) Get() float64 {
	spk.mut.RLock()
	progress := spk.progress[spk.currentID]
	spk.mut.RUnlock()
	return progress
}

func (spk *SimpleProgressKeeper) Complete(id string, msg string) {
	// id is done
	spk.mut.Lock()
	spk.progress[id] = 1.0
	spk.msgs[id] = msg
	spk.mut.Unlock()
}

func (spk *SimpleProgressKeeper) Message(id string) string {
	spk.mut.RLock()
	msg := spk.msgs[id]
	spk.mut.RUnlock()
	return msg
}

func (spk *SimpleProgressKeeper) SetMessage(msg string) {
	spk.mut.Lock()
	spk.msgs[spk.currentID] = msg
	spk.mut.Unlock()
}

// Write can report the progress of an implementation of the ProgressKeeper.
// delay is how long to wait beteween each status update of the progress indicator.
func Write(wg *sync.WaitGroup, p ProgressKeeper, delay time.Duration) {
	// endless loop without a timeout, on purpose
	for {
		if p.Get() >= 1.0 {
			fmt.Println(p)
			break
		} else {
			fmt.Println(p)
		}
		time.Sleep(delay)
	}
	wg.Done()
}

func back() {
	fmt.Print("\b")
}

func startOfLine() {
	fmt.Print("\r")
}

// TODO: Think this through to allow multiple file copying operations at once,
//       that all report progress.
// Spin can report the progress of an implementation of the ProgressKeeper.
// delay is how long to wait beteween each status update of the progress indicator.
func Spin(wg *sync.WaitGroup, p ProgressKeeper, delay time.Duration) {
	spinner := []string{"|", "/", "-", "\\"}
	i := 0
	// endless loop without a timeout, on purpose
	for {
		startOfLine()
		//if p.Skipped() {
		//	fmt.Printf("[x] %s", spinner[i], p.Message())
		//} else if p.Complete() {
		//	fmt.Printf("[✓] %s\n", spinner[i], p.Message())
		//} else {
		id := "?"
		fmt.Printf("[%s] %s", spinner[i], p.Message(id))
		//}
		i++
		if i == len(spinner) {
			i = 0
		}
		if p.Get() >= 1.0 {
			break
		}
		time.Sleep(delay)
	}
	fmt.Println()
	fmt.Printf("[✓] done")
	fmt.Println("                                                                        ")
	wg.Done()
}
