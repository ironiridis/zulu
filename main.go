package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var flagOneShot bool
var flagPermissive bool
var flagSequential bool
var flagRate time.Duration
var flagDriver string

const rateMinimum = time.Minute * 5
const rateDefault = time.Hour * 1

type driverInterface interface {
	doSync(connstr string, tod time.Time) error
	description() string
}

var drivers map[string]driverInterface

// connResolve breaks a connection string into an optional username, an optional password,
// a required host ip (optionally resolved from a hostname), and an optional port.
func connResolve(conn, defaultUsername, defaultPassword string, defaultPort int) (*net.TCPAddr, string, string, error) {
	var auth, host, username, password string
	parts := strings.SplitN(conn, "@", 2)
	if len(parts) == 2 {
		auth, host = parts[0], parts[1]
		authparts := strings.SplitN(auth, ":", 2)
		if len(authparts) == 2 {
			if authparts[0] != "" {
				username = authparts[0]
			} else {
				username = defaultUsername
			}
			password = authparts[1]
		} else {
			username = auth
			password = defaultPassword
		}
	} else {
		username = defaultUsername
		password = defaultPassword
		host = conn
	}

	var addr *net.TCPAddr
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		ip, err := net.ResolveIPAddr("ip", host)
		if err != nil {
			return nil, "", "", err
		}
		addr = &net.TCPAddr{IP: ip.IP, Port: defaultPort}
	}
	return addr, username, password, nil
}

func registerDriver(n string, d driverInterface) {
	if n == "list" {
		log.Printf("driver attempting to register using the reserved keyword %q\n", n)
		panic("internal error")
	}
	if drivers == nil {
		drivers = make(map[string]driverInterface)
	}
	drivers[n] = d
}

func invoke(d driverInterface, c string) {
	err := d.doSync(c, time.Now())
	if err == nil {
		return
	}
	if flagPermissive {
		log.Printf("warning: sync with %q failed: %v\n", c, err)
	} else {
		log.Fatalf("error: sync with %q failed: %v\n", c, err)
	}
}

func invokeSequential(d driverInterface, conns []string) {
	for _, c := range conns {
		invoke(d, c)
	}
}

func invokeThreads(d driverInterface, conns []string) {
	var wg sync.WaitGroup
	for i := range conns {
		wg.Add(1)
		go func(c string, wg *sync.WaitGroup) {
			invoke(d, c)
			wg.Done()
		}(conns[i], &wg)
	}
	wg.Wait()
}

func main() {
	flag.BoolVar(&flagOneShot, "oneshot", false, "Set time and exit")
	flag.BoolVar(&flagPermissive, "permissive", false, "Don't exit on any failure")
	flag.BoolVar(&flagSequential, "sequential", false, "Process each argument sequentially")
	flag.DurationVar(&flagRate, "rate", rateDefault, "Frequency with which to set time")
	flag.StringVar(&flagDriver, "driver", "list", "Type of device to sync")
	flag.Parse()

	if flagRate < rateMinimum {
		log.Fatalf("Specified rate %v is too short; minimum is %v.\n", flagRate, rateMinimum)
	}
	if flagRate != rateDefault && flagOneShot {
		log.Println("Flags oneshot and rate both specified, but oneshot ignores rate.")
	}

	d, ok := drivers[flagDriver]
	if !ok {
		if flagDriver == "list" {
			log.Println("List of available drivers in this build:")
		} else {
			log.Printf("Unknown driver %q specified; please choose from the following list:\n", flagDriver)
		}
		var k string
		for k = range drivers {
			log.Printf("%q: %s\n", k, drivers[k].description())
		}
		log.Fatalf("Try: %s -driver %q 10.0.0.2\n", os.Args[0], k)
	}

	for {
		if flagSequential {
			invokeSequential(d, flag.Args())
		} else {
			invokeThreads(d, flag.Args())
		}
		if flagOneShot {
			log.Println("Oneshot specified, exiting")
			return
		}
		time.Sleep(flagRate)
	}
}
