package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
)

const mpath string = "F:\\myProjects\\goserver\\backend\\main.go"

type reqCap struct {
	timer    *time.Timer
	callback func(o *reqCap)
	duration time.Duration
	state    bool
	cmd      *exec.Cmd
	name     string
}

func New_reqCap(call func(o *reqCap), dur time.Duration) *reqCap {
	aux := reqCap{
		timer:    time.NewTimer(time.Second * 20),
		callback: call,
		duration: dur,
		state:    true,
	}
	aux.timer.Stop()

	go func(r *reqCap) {
		for r.state {
			<-r.timer.C
			r.callback(r)
		}
	}(&aux)

	return &aux
}

func (o *reqCap) Close_reqCap() {
	o.timer.Stop()
	o.state = false
}

func (o *reqCap) Capture(name string) bool {
	o.name = name
	return o.timer.Reset(o.duration)
}

func kill(pid string) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", pid)
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}

func restartProgram(o *reqCap) {

	fmt.Println(o.name)

	if o.cmd != nil {
		//fmt.Println("deleting ", strconv.Itoa(o.cmd.Process.Pid))
		kill(strconv.Itoa(o.cmd.Process.Pid))
	}

	o.cmd = exec.Command("F:\\golang\\bin\\go.exe", "run", mpath)

	o.cmd.Stdout = os.Stdout
	o.cmd.Stderr = os.Stderr

	if err := o.cmd.Start(); err != nil {
		panic(err)
	}
	fmt.Println("--------------------------------------------------")
	fmt.Println("----------------------CHILD PROGRAM---------------")
	/* 	str := fmt.Sprintf("(ParentProcessId=%d)", o.cmd.Process.Pid)

	   	aux, _ := exec.Command("wmic", "process", "where", str, "get", "Caption,ProcessId").Output()
	   	aux1 := string(aux)
	   	fmt.Println(aux1)

	   	re := regexp.MustCompile("\\d+")
	   	m := string(re.Find(aux))

	   	o.pid = m
	*/

}

func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	capt := New_reqCap(restartProgram, 2000*time.Millisecond)
	capt.name = "Inner Porg Initializing"
	restartProgram(capt)

	go func() {

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {

					if filepath.Ext(event.Name) == ".go" {
						fmt.Println("request " + event.Name)
						capt.Capture("File Modified: " + event.Name)
					}
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					fmt.Println(" Error on watcher Errors")
				}

			}
		}
	}()

	err = watcher.Add("F:\\myProjects\\goserver\\backend")
	if err != nil {
		panic("watcher add error")
	}

	<-make(chan int)
}
