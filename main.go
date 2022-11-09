package main

import (
	"autoreload/reqcap"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
)

const mpath string = "F:\\myProjects\\goserver\\backend\\"
const mfile string = "main.go"

func kill_pgm(pid string) error {
	kill_pgm := exec.Command("TASKKILL", "/T", "/F", "/PID", pid)
	kill_pgm.Stderr = os.Stderr
	kill_pgm.Stdout = os.Stdout
	return kill_pgm.Run()
}

var count int = 0

func restartProgram(o *reqcap.ReqCap) {

	fmt.Println(o.Name)

	if o.Cmd != nil {
		kill_pgm(strconv.Itoa(o.Pid))
	}
	//o.Cmd = exec.Command("F:\\golang\\bin\\go.exe", "run", mpath)
	o.Cmd = exec.Command("F:\\myProjects\\goserver\\autoreaload\\start.bat", mpath, mfile)

	if err := o.Cmd.Start(); err != nil {
		panic(err)
	}

	fmt.Println(o.Cmd.Process.Pid)
	u := fmt.Sprintf("(ParentProcessId=%d)", o.Cmd.Process.Pid)
	pidStr, _ := exec.Command("wmic", "process", "where", u, "get", "Caption,ProcessId").Output()
	reg := regexp.MustCompile("\\d+")
	pid, _ := strconv.Atoi(string(reg.Find(pidStr)))
	o.Pid = pid

	fmt.Println(fmt.Sprintf("%d  -------------------------------", count))
	count++

}

func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	capt := reqcap.New_reqCap(restartProgram, 2000*time.Millisecond)
	capt.Name = "Inner Porg Initializing"
	restartProgram(capt)

	go func() {

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) {
					fmt.Println("created ", event.Name)
					fileInfo, _ := os.Stat(event.Name)
					if fileInfo.IsDir() {
						fmt.Println("add file")
						watcher.Add(event.Name)
					}
				}
				if event.Has(fsnotify.Remove) {
					fmt.Println("Delete ", event.Name)
					ff, _ := os.Stat(event.Name)
					if ff.IsDir() {
						fmt.Println("remove dir")
						//watcher.Lis(event.Name)
					}
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

	watchAllFolders(mpath, func(p string) {
		err := watcher.Add(p)
		if err != nil {
			panic(err)
		}
	})
	<-make(chan int)
}

func watchAllFolders(p string, evt func(string)) {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}
	evt(p)
	for _, f := range files {
		if f.IsDir() {
			np := p + "\\" + f.Name()

			watchAllFolders(np, evt)
		}

	}
}
