// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package universal_testing

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DataDog/test-infra-definitions/common/utils"
	"github.com/stretchr/testify/suite"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner/parameters"
	"github.com/DataDog/datadog-agent/test/universal-testing/client"
	"github.com/DataDog/datadog-agent/test/universal-testing/params"
)

const (
	deleteTimeout = 30 * time.Minute
)

// Suite manages the environment creation and runs tests.
type Suite[Env any] struct {
	suite.Suite

	params          params.Params
	defaultInfraDef InfraProvider[Env]
	currentInfraDef InfraProvider[Env]
	firstFailTest   string

	env *Env

	isUpdateEnvCalledInThisTest bool
}

type suiteConstraint[Env any] interface {
	suite.TestingSuite
	initSuite(infraName string, infraDef InfraProvider[Env], options ...params.Option)
}

// Run runs the tests defined in testSuite.
//
// t is an instance of type [*testing.T].
//
// testSuite is a pointer to a structure with a [universal_testing.Suite] embedded struct.
//
// infraDef defines the infrastructure which the tests will run on.
//
// options is an optional list of options like [DevMode] or [SkipDeleteOnFailure].
func Run[Env any, T suiteConstraint[Env]](t *testing.T, testSuite T, infraDef InfraProvider[Env], options ...params.Option) {
	suiteType := reflect.TypeOf(testSuite).Elem()
	name := suiteType.Name()
	pkgPaths := suiteType.PkgPath()
	pkgs := strings.Split(pkgPaths, "/")

	// Use the hash of PkgPath in order to generate a unique name
	hash := utils.StrHash(pkgs...)

	// Example: "e2e-e2eSuite-cbb731954db42b"
	defaultName := fmt.Sprintf("%v-%v-%v", pkgs[len(pkgs)-1], name, hash)

	testSuite.initSuite(defaultName, infraDef, options...)
	suite.Run(t, testSuite)
}

func (suite *Suite[Env]) initSuite(name string, infraDef InfraProvider[Env], options ...params.Option) {
	suite.params.Name = name
	suite.defaultInfraDef = infraDef
	for _, o := range options {
		o(&suite.params)
	}
}

// Env returns the current environment.
// In order to improve the efficiency, this function behaves as follow:
//   - It creates the default environment if no environment exists.
//   - It restores the default environment if [universal_testing.Suite.UpdateEnv] was not already called during this test.
//     This avoids having to restore the default environment for each test even if [universal_testing.Suite.UpdateEnv] immediately
//     overrides the environment.
func (suite *Suite[Env]) Env() *Env {
	if suite.env == nil || !suite.isUpdateEnvCalledInThisTest {
		suite.UpdateEnv(suite.defaultInfraDef)
	}
	return suite.env
}

// BeforeTest is executed right before the test starts and receives the suite and test names as input.
// This function is called by [testify Suite].
//
// If you override BeforeTest in your custom test suite type, the function must call [universal_testing.Suite.BeforeTest].
//
// [testify Suite]: https://pkg.go.dev/github.com/stretchr/testify/suite
func (suite *Suite[Env]) BeforeTest(suiteName, testName string) {
	_ = suiteName
	_ = testName
	suite.isUpdateEnvCalledInThisTest = false
}

// AfterTest is executed right after the test finishes and receives the suite and test names as input.
// This function is called by [testify Suite].
//
// If you override AfterTest in your custom test suite type, the function must call [universal_testing.Suite.AfterTest].
//
// [testify Suite]: https://pkg.go.dev/github.com/stretchr/testify/suite
func (suite *Suite[Env]) AfterTest(suiteName, testName string) {
	if suite.T().Failed() && suite.firstFailTest == "" {
		// As far as I know, there is no way to prevent other tests from being
		// run when a test fail. Even calling panic doesn't work.
		// Instead, this code stores the name of the first fail test and prevents
		// the environment to be updated.
		// Note: using os.Exit(1) prevents other tests from being run but at the
		// price of having no test output at all.
		suite.firstFailTest = fmt.Sprintf("%v.%v", suiteName, testName)
	}
}

// SetupSuite method will run before the tests in the suite are run.
// This function is called by [testify Suite].
//
// If you override SetupSuite in your custom test suite type, the function must call [universal_testing.Suite.SetupSuite].
//
// [testify Suite]: https://pkg.go.dev/github.com/stretchr/testify/suite
func (suite *Suite[Env]) SetupSuite() {
	skipDelete, _ := runner.GetProfile().ParamStore().GetBoolWithDefault(parameters.SkipDeleteOnFailure, false)
	if skipDelete {
		suite.params.SkipDeleteOnFailure = true
	}

	// Check if the Env type is correct otherwise raises an error before creating the env.
	err := client.CheckEnvStructValid[Env]()
	suite.Require().NoError(err)
}

// TearDownSuite run after all the tests in the suite have been run.
// This function is called by [testify Suite].
//
// If you override TearDownSuite in your custom test suite type, the function must call [universal_testing.Suite.TearDownSuite].
//
// [testify Suite]: https://pkg.go.dev/github.com/stretchr/testify/suite
func (suite *Suite[Env]) TearDownSuite() {
	if runner.GetProfile().AllowDevMode() && suite.params.DevMode {
		return
	}

	if suite.firstFailTest != "" && suite.params.SkipDeleteOnFailure {
		suite.Require().FailNow(fmt.Sprintf("%v failed. As SkipDeleteOnFailure feature is enabled the tests after %v were skipped. "+
			"The environment of %v was kept.", suite.firstFailTest, suite.firstFailTest, suite.firstFailTest))
		return
	}

	// TODO: Implement retry on delete
	ctx, cancel := context.WithTimeout(context.Background(), deleteTimeout)
	defer cancel()

	if suite.currentInfraDef == nil {
		// this happens if there's an error when instantiating the env
		return
	}

	err := suite.currentInfraDef.DeleteInfra(ctx, suite.params.Name)
	if err != nil {
		suite.T().Errorf("unable to delete infrastructure: err :%v", err)
		suite.T().Fail()
	}
}

// UpdateEnv updates the environment.
// This affects only the test that calls this function.
// Test functions that don't call UpdateEnv have the environment defined by [universal_testing.Run].
func (suite *Suite[Env]) UpdateEnv(infraDef InfraProvider[Env]) {
	if infraDef != suite.currentInfraDef {
		if (suite.firstFailTest != "" || suite.T().Failed()) && suite.params.SkipDeleteOnFailure {
			// In case of failure, do not override the environment
			suite.T().SkipNow()
		}
		env, err := infraDef.ProvisionInfraAndInitializeEnv(suite.T(), context.Background(), suite.params.Name, false)
		suite.Require().NoError(err)
		suite.env = env
		suite.currentInfraDef = infraDef
	}
	suite.isUpdateEnvCalledInThisTest = true
}
