package e2e

import (
	"flag"
	"testing"
)

var useLocal = flag.Bool("local", false, "run tests on a local VM")

type testSuite struct {
	Suite[VMEnv]
}

func TestVMSuite(t *testing.T) {
	if *useLocal {
		Run(t, &testSuite{}, EC2VMStackDef())
	} else {
		Run(t, &testSuite{}, LocalVMInfraDef("win2016.json"))
	}
}

func (t *testSuite) Test() {
	t.Env().VM.Execute("cat /etc/os-release")
}
