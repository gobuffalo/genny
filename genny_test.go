package genny

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	PWD     string
	realPWD string
	cancels []context.CancelFunc
	*require.Assertions
}

func (s *Suite) SetupTest() {
	s.Assertions = require.New(s.T())

	pwd, err := os.Getwd()
	s.NoError(err)
	s.realPWD = pwd

	dir, err := ioutil.TempDir("", "")
	s.NoError(err)
	s.PWD = dir

	s.NoError(os.Chdir(s.PWD))
}

func (s *Suite) TearDownTest() {
	s.NoError(os.Chdir(s.realPWD))
	s.NoError(os.RemoveAll(s.PWD))
}

func (s *Suite) Command(name string, args ...string) *exec.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	s.cancels = append(s.cancels, cancel)
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd
}

func Test_Suite(t *testing.T) {
	s := &Suite{}
	suite.Run(t, s)
}
