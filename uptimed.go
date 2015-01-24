package main

/*
#cgo LDFLAGS: -framework Cocoa -framework ApplicationServices
#cgo CFLAGS: -x objective-c
#include "ApplicationServices/ApplicationServices.h"
#import "header.h"

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
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

const (
	Freq = time.Second
	Min  = 5 * time.Minute
)

func main() {
	go runMainThread()
	mainThread <- func() { C.StartApp() }

	totalIdleTime := time.Duration(0)
	for idleChange := range sysIdleTimeTicker(Freq, Min) {
		totalIdleTime += idleChange
		fmt.Println(idleChange.String())
		setMenuLabel(totalIdleTime.String())
	}
}

func setMenuLabel(l string) {
	mainThread <- func() {
		cs := C.CString(l)
		C.SetLabelText(cs)
		C.free(unsafe.Pointer(cs))
	}
}

var mainThread = make(chan func())

func runMainThread() {
	for f := range mainThread {
		go func() {
			runtime.LockOSThread()
			f()
		}()
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
