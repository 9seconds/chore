package testlib

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

const pipesToFlushTimeout = time.Second

type RedirectStreamsTestSuite struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (suite *RedirectStreamsTestSuite) Setup(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	outR, outW, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	errR, errW, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer wg.Done()
		io.Copy(&suite.stdout, outR)
	}()

	go func() {
		defer wg.Done()
		io.Copy(&suite.stderr, errR)
	}()

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	t.Cleanup(func() {
		time.Sleep(pipesToFlushTimeout)

		outW.Close()
		errW.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr

		wg.Wait()
	})

	os.Stdout = outW
	os.Stderr = errW
}

func (suite *RedirectStreamsTestSuite) Stdout() string {
	return suite.stdout.String()
}

func (suite *RedirectStreamsTestSuite) Stderr() string {
	return suite.stderr.String()
}
