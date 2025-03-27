package main

/*
#include <stdio.h>
#include <time.h>
#include <errno.h>
#cgo LDFLAGS:
// simple C code
int do_some_c_call(void)
{

	struct timespec timeout;
	timeout.tv_sec = 0;
	timeout.tv_nsec = 1000;

	if (0!=nanosleep(&timeout, &timeout) && errno == EINTR)
	{
		printf("INTERRUPTED\n");
		return -1;
	}

	return 0;
}
*/
import "C"


import (
        "sync"
	"os"
	"atmi"
	"time"
)

func main() {
    
	atmi.RuntimeInit()
	var wg sync.WaitGroup

	go func() {
		time.Sleep(60 * time.Second)
		//OK in this time..
		os.Exit(0)
	}()

	for i :=0; i<100000; i++  {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if C.do_some_c_call()!=0 {
					os.Exit(-1)
				}
			}
		}()

		/* check for run-time */
	}
	wg.Wait()
}
