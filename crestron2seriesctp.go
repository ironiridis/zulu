package main

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type driverCrestron2SeriesCTP struct{}

func (d *driverCrestron2SeriesCTP) doSync(connstr string, tod time.Time) error {
	var err error
	a, _, p, err := connResolve(connstr, "", "", 41795)
	if err != nil {
		return err
	}

	t, err := net.DialTCP("tcp", nil, a)
	if err != nil {
		return fmt.Errorf("cannot connect to %s: %w", a, err)
	}
	defer t.Close()

	// this should get us to a prompt, and possibly respond to a password prompt
	var setuptime int
	for {
		if setuptime > 10 {
			return fmt.Errorf("got stuck in connection setup, giving up")
		}
		setuptime++
		err = t.SetReadDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			return fmt.Errorf("cannot set read deadline: %w", err)
		}
		buf := make([]byte, 250)
		n, err := t.Read(buf)
		if n == 12 && bytes.Equal(buf[:n], []byte("Password? \x0d\x0a")) {
			pwtx := fmt.Sprintf("%s\x0d", p)
			_, err = t.Write([]byte(pwtx))
			if err != nil {
				return fmt.Errorf("cannot send password response: %w", err)
			}
			continue
		}
		if n == 0 {
			// we (probably) timed out on this read, so probe for a response
			_, err = t.Write([]byte{0x0d})
			if err != nil {
				// if the connection closed or some other fatal error has happened
				// we will likely catch it here anyway
				return fmt.Errorf("cannot complete connection setup: %w", err)
			}
			continue
		}
		if buf[n-1] == 0x3e {
			// we got something that ended in ">", assume it's a prompt
			break
		}
	}

	ctpTimeFmt := "15:04:05 01/02/2006"
	tdcmd := fmt.Sprintf("TIMEDATE %s\x0d", time.Now().Format(ctpTimeFmt))
	_, err = t.Write([]byte(tdcmd))
	if err != nil {
		return fmt.Errorf("failed sending time command: %w", err)
	}
	time.Sleep(1)

	byecmd := "bye\x0d"
	_, err = t.Write([]byte(byecmd))
	if err != nil {
		return fmt.Errorf("failed to cleanly end control session: %w", err)
	}
	time.Sleep(1)

	return nil
}

func (d *driverCrestron2SeriesCTP) description() string {
	return "Connects to a 2-series Crestron controller over CTP. Typical port is 41795."
}

func init() {
	registerDriver("crestron2seriesctp", &driverCrestron2SeriesCTP{})
}
