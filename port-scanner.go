package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type (
	scanner struct {
		ip   string
		lock *semaphore.Weighted
	}
)

var (
	all    bool
	open   bool
	closed bool
	first  int
	last   int
	ip     string
)

func init() {
	cmd := Command{
		Name:        "port-scanner",
		Description: "Scans open/closed ports on an IP basis",
		Action:      portScanner,
	}

	cmd.FlagSet = flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	cmd.FlagSet.BoolVar(&all, "all", false, "list all available open and closed ports")
	cmd.FlagSet.BoolVar(&all, "a", false, "list all available open and closed ports (shorthand)")
	cmd.FlagSet.BoolVar(&open, "open", true, "list all available open ports")
	cmd.FlagSet.BoolVar(&open, "o", true, "list all available open ports (shorthand)")
	cmd.FlagSet.BoolVar(&closed, "closed", false, "list all available closed ports")
	cmd.FlagSet.BoolVar(&closed, "c", false, "list all available closed ports (shorthand)")
	cmd.FlagSet.IntVar(&first, "first", 1, "first port to listen on")
	cmd.FlagSet.IntVar(&first, "f", 1, "first port to listen on (shorthand)")
	cmd.FlagSet.IntVar(&last, "last", 65535, "last port enumerate to")
	cmd.FlagSet.IntVar(&last, "l", 65535, "last port enumerate to (shorthand)")
	cmd.FlagSet.StringVar(&ip, "ip", "127.0.0.1", "IP to test")
	Commands[cmd.Name] = cmd
}

func portScanner() {
	ps := &scanner{
		ip:   ip,
		lock: semaphore.NewWeighted(ulimit()),
	}
	ps.Start(500 * time.Millisecond)
}

func ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
	if err != nil {
		panic(err)
	}

	s := strings.TrimSpace(string(out))

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func scan(ip string, port int, timeout time.Duration) (int, string) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			scan(ip, port, timeout)
		} else {
			return port, "closed"
		}
		return 0, ""
	}

	conn.Close()
	return port, "open"
}

func (ps *scanner) Start(timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := first; port <= last; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			p, ptype := scan(ps.ip, port, timeout)

			if all {
				fmt.Println(fmt.Sprintf("%v\t%v", p, ptype))
			} else if open && ptype == "open" {
				fmt.Println(fmt.Sprintf("%v\t%v", p, ptype))
			} else if closed && ptype == "closed" {
				fmt.Println(fmt.Sprintf("%v\t%v", p, ptype))
			}

		}(port)
	}
}
