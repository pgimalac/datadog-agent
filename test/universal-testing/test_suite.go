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

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/runner/parameters"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/params"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/infra"
	"github.com/DataDog/test-infra-definitions/common/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/suite"
)

const (
	deleteTimeout = 30 * time.Minute
)

// Suite manages the environment creation and runs E2E tests.
type Suite[Env any] struct {
	suite.Suite

	params          params.Params
	defaultStackDef *StackDefinition[Env]
	currentStackDef *StackDefinition[Env]
	firstFailTest   string

	// These fields are initialized in SetupSuite
	env *Env

	isUpdateEnvCalledInThisTest bool
}

type suiteConstraint[Env any] interface {
	suite.TestingSuite
	initSuite(stackName string, stackDef *StackDefinition[Env], options ...params.Option)
}

// Run runs the tests defined in e2eSuite
//
// t is an instance of type [*testing.T].
//
// e2eSuite is a pointer to a structure with a [e2e.Suite] embbeded struct.
//
// stackDef defines the stack definition.
//
// options is an optional list of options like [DevMode], [SkipDeleteOnFailure] or [WithStackName].
//
//	type vmSuite struct {
//		e2e.Suite[e2e.VMEnv]
//	}
//	// ...
//	e2e.Run(t, &vmSuite{}, e2e.EC2VMStackDef())
func Run[Env any, T suiteConstraint[Env]](t *testing.T, e2eSuite T, stackDef *StackDefinition[Env], options ...params.Option) {
	suiteType := reflect.TypeOf(e2eSuite).Elem()
	name := suiteType.Name()
	pkgPaths := suiteType.PkgPath()
	pkgs := strings.Split(pkgPaths, "/")

	// Use the hash of PkgPath in order to have a uniq stack name
	hash := utils.StrHash(pkgs...)

	// Example: "e2e-e2eSuite-cbb731954db42b"
	defaultStackName := fmt.Sprintf("%v-%v-%v", pkgs[len(pkgs)-1], name, hash)

	e2eSuite.initSuite(defaultStackName, stackDef, options...)
	suite.Run(t, e2eSuite)
}

func (suite *Suite[Env]) initSuite(stackName string, stackDef *StackDefinition[Env], options ...params.Option) {
	suite.params.StackName = stackName
	suite.defaultStackDef = stackDef
	for _, o := range options {
		o(&suite.params)
	}
}

// Env returns the current environment.
// In order to improve the efficiency, this function behaves as follow:
//   - It creates the default environment if no environment exists.
//   - It restores the default environment if [e2e.Suite.UpdateEnv] was not already called during this test.
//     This avoid having to restore the default environment for each test even if [suite.UpdateEnv] immedialy
//     overrides the environment.
func (suite *Suite[Env]) Env() *Env {
	if suite.env == nil || !suite.isUpdateEnvCalledInThisTest {
		suite.UpdateEnv(suite.defaultStackDef)
	}
	return suite.env
}

// BeforeTest is executed right before the test starts and receives the suite and test names as input.
// This function is called by [testify Suite].
//
// If you override BeforeTest in your custom test suite type, the function must call [e2e.Suite.BeforeTest].
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
// If you override AfterTest in your custom test suite type, the function must call [e2e.Suite.AfterTest].
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
// If you override SetupSuite in your custom test suite type, the function must call [e2e.Suite.SetupSuite].
//
// [testify Suite]: https://pkg.go.dev/github.com/stretchr/testify/suite
func (suite *Suite[Env]) SetupSuite() {
	skipDelete, _ := runner.GetProfile().ParamStore().GetBoolWithDefault(parameters.SkipDeleteOnFailure, false)
	if skipDelete {
		suite.params.SkipDeleteOnFailure = true
	}

	suite.Require().NotEmptyf(suite.params.StackName, "The stack name is empty. You must define it with WithName")
	// Check if the Env type is correct otherwise raises an error before creating the env.
	err := client.CheckEnvStructValid[Env]()
	suite.Require().NoError(err)
}

// TearDownSuite run after all the tests in the suite have been run.
// This function is called by [testify Suite].
//
// If you override TearDownSuite in your custom test suite type, the function must call [e2e.Suite.TearDownSuite].
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
	err := infra.GetStackManager().DeleteStack(ctx, suite.params.StackName)
	if err != nil {
		suite.T().Errorf("unable to delete stack: %s, err :%v", suite.params.StackName, err)
		suite.T().Fail()
	}
}

func createEnv[Env any](suite *Suite[Env], stackDef *StackDefinition[Env]) (*Env, auto.UpResult, error) {
	var env *Env
	ctx := context.Background()

	_, stackOutput, err := infra.GetStackManager().GetStackNoDeleteOnFailure(
		ctx,
		suite.params.StackName,
		stackDef.configMap,
		func(ctx *pulumi.Context) error {
			var err error
			env, err = stackDef.envFactory(ctx)
			return err
		}, false)

	return env, stackOutput, err
}

// UpdateEnv updates the environment.
// This affects only the test that calls this function.
// Test functions that don't call UpdateEnv have the environment defined by [e2e.Run].
func (suite *Suite[Env]) UpdateEnv(stackDef *StackDefinition[Env]) {
	if stackDef != suite.currentStackDef {
		if (suite.firstFailTest != "" || suite.T().Failed()) && suite.params.SkipDeleteOnFailure {
			// In case of failure, do not override the environment
			suite.T().SkipNow()
		}
		env, upResult, err := createEnv(suite, stackDef)
		suite.Require().NoError(err)
		err = client.CallStackInitializers(suite.T(), env, upResult)
		suite.Require().NoError(err)
		suite.env = env
		suite.currentStackDef = stackDef
	}
	suite.isUpdateEnvCalledInThisTest = true
}
