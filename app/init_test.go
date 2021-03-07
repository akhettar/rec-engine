package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

var app *App

func TestMain(m *testing.M) {

	// 1. Initialise redis server
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	app = InitialiseApp(fmt.Sprintf("redis://%s", s.Addr()))

	// 2. Run all the tests
	code := m.Run()

	// 3. terminate everything
	os.Exit(code)
}
