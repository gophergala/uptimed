package main

/*

#cgo LDFLAGS: -framework ApplicationServices
#include "ApplicationServices/ApplicationServices.h"

int64_t SystemIdleTime(void) {
  CFTimeInterval timeSinceLastEvent;
  timeSinceLastEvent = CGEventSourceSecondsSinceLastEventType(kCGEventSourceStateCombinedSessionState, kCGAnyInputEventType);
  return timeSinceLastEvent * 1000000000;
}

*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
	"time"
)

const (
	Freq = time.Second
	Min  = 30 * time.Second
)

func main() {
	b, _ := boottime()
	s, _ := sleeptime()
	w, _ := waketime()
	fmt.Println(b, s, w)
	for idleTime := range sysIdleTimeTicker(Freq, Min) {
		fmt.Println(idleTime)
	}
}

func boottime() (*time.Time, error) {
	return sysCtlTimeByName("kern.boottime")
}

func sleeptime() (*time.Time, error) {
	return sysCtlTimeByName("kern.sleeptime")
}

func waketime() (*time.Time, error) {
	return sysCtlTimeByName("kern.waketime")
}

func sysCtlTimeByName(name string) (*time.Time, error) {
	v, err := syscall.Sysctl(name)
	if err != nil {
		return nil, fmt.Errorf("%s error: %q", name, err)
	}

	var secs int64
	buf := bytes.NewBufferString(v)
	if err = binary.Read(buf, binary.LittleEndian, &secs); err != nil {
		return nil, fmt.Errorf("binary.Read error: %q", err)
	}
	u := time.Unix(int64(secs), 0)
	return &u, nil
}

func sysIdleTimeTicker(freq time.Duration, min time.Duration) <-chan time.Duration {
	c := make(chan time.Duration)
	go func() {
		var prev, curr time.Duration
		for _ = range time.NewTicker(freq).C {
			curr = time.Duration(C.SystemIdleTime())
			if curr < prev && prev >= min {
				c <- (prev + freq)
			}
			prev = curr
		}
	}()
	return c
}
