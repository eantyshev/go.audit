package main

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	flag.Parse()
	// we need to wait for all other components:
	// http may be on but mongo isn't
	// TODO: implement a more intelligent polling
	time.Sleep(10 * time.Second)
	os.Exit(m.Run())
}
