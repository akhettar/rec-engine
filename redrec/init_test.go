package redrec

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

var rr *Redrec

func TestMain(m *testing.M) {

	// 1. Initialise redis server
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	addr := s.Addr()

	rr, err = New(fmt.Sprintf("redis://%s", addr))
	if err != nil {
		log.Fatal("Redis init Error", err)
	}

	// 2. Run all the tests
	code := m.Run()

	// 3. terminate everything
	os.Exit(code)
}
