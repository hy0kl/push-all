package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	upstreamCmd := "git remote -v | awk '{print $1}' | sort -u"
	cmd := exec.Command("bash","-c", upstreamCmd )
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "exec find git upstream cmd get exception, cmd> %s, err: %v\n", upstreamCmd, err)
		os.Exit(1)
	}

	cmdOut := stdout.String()
	if cmdOut == "" {
		_, _ = fmt.Fprintf(os.Stdout, "cmd output is empty\ncmd> %s\nstderr: %s", upstreamCmd, stderr.String())
		os.Exit(2)
	}

	upstreamBox := strings.Split(cmdOut, "\n")
	if len(upstreamBox) <= 0 {
		_, _ = fmt.Fprintf(os.Stdout, "There is no git source in the current working directory.\n")
		os.Exit(3)
	}

	mt := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, upstream := range upstreamBox{
		if upstream == "" {
			continue
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, up string) {
			defer wg.Done()

			pushCmd := exec.Command("git", "push", up, "master")
			stdout, err := pushCmd.CombinedOutput()
			if err != nil {
				mt.Lock()
				_, _ = fmt.Fprintf(os.Stdout, "exec cmd stdout pipe get exception.\nupstream: %s, stderr: %s, err: %v\n",
					up, string(stdout), err)
				mt.Unlock()
				return
			}

			mt.Lock()
			_, _ = fmt.Fprintf(os.Stdout, "---> %s begain <---\n", up)
			_, _ = fmt.Fprintf(os.Stdout, "%s", string(stdout))
			_, _ = fmt.Fprintf(os.Stdout, "--- %s end ---\n", up)
			mt.Unlock()
		}(&wg, upstream)
	}

	wg.Wait()
}
