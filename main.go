package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	branchCmd := "git remote -v | awk '{print $1}' | sort -u"
	cmd := exec.Command("bash","-c", branchCmd )
	stdout, err := cmd.Output()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "exec find git branch cmd get exception, cmd> %s, err: %v\n", branchCmd, err)
		os.Exit(1)
	}

	cmdOut := string(stdout)

	if cmdOut == "" {
		_, _ = fmt.Fprintf(os.Stdout, "cmd output is empty, cmd> %s\n", branchCmd)
		os.Exit(2)
	}

	branchBox := strings.Split(cmdOut, "\n")
	if len(branchBox) <= 0 {
		_, _ = fmt.Fprintf(os.Stdout, "There is no git source in the current working directory.\n")
		os.Exit(3)
	}

	mt := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, branch := range branchBox{
		if branch == "" {
			continue
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, name string) {
			defer wg.Done()

			pushCmd := exec.Command("git", "push", name, "master")
			stdout, err := pushCmd.CombinedOutput()
			if err != nil {
				mt.Lock()
				_, _ = fmt.Fprintf(os.Stdout, "exec cmd stdout pipe get exception, err: %v", err)
				mt.Unlock()
				return
			}

			mt.Lock()
			_, _ = fmt.Fprintf(os.Stdout, "---> %s begain <---\n", name)
			_, _ = fmt.Fprintf(os.Stdout, "%s", string(stdout))
			_, _ = fmt.Fprintf(os.Stdout, "--- %s end ---\n", name)
			mt.Unlock()
		}(&wg, branch)
	}

	wg.Wait()
}
