package cli

import (
	"os"
	"testing"
	"time"

	"github.com/tpbowden/swarm-ingress-router/types"
)

var currentConfig types.Configuration
var appType string
var success bool

type fakeApp struct {
}

func (f fakeApp) Start() {
	return
}

func fakeAbort(string) {
	success = false
}

func fakeServer(config types.Configuration) types.Startable {
	currentConfig = config
	appType = "server"
	return types.Startable(fakeApp{})
}

func fakeCollector(config types.Configuration) types.Startable {
	currentConfig = config
	appType = "collector"
	return types.Startable(fakeApp{})
}

var subject = CLI{
	initServer:    fakeServer,
	initCollector: fakeCollector,
	abort:         fakeAbort,
}

type cliTest struct {
	description     string
	environment     map[string]string
	expectedConfig  types.Configuration
	expectedCommand string
	arguments       []string
	success         bool
}

var defaultConfig = types.Configuration{
	Redis:        "localhost:6379",
	Bind:         "0.0.0.0",
	PollInterval: 10 * time.Second,
}

var testCases = []cliTest{
	{
		description:     "Default server config with no environment",
		environment:     map[string]string{},
		expectedConfig:  defaultConfig,
		expectedCommand: "server",
		arguments:       []string{"anything", "server"},
		success:         true,
	},
	{
		description:     "Default collector config with no environment",
		environment:     map[string]string{},
		expectedConfig:  defaultConfig,
		expectedCommand: "collector",
		arguments:       []string{"anything", "collector"},
		success:         true,
	},
	{
		description: "Collector config with overrides",
		environment: map[string]string{
			"INGRESS_REDIS":         "some-redis:1234",
			"INGRESS_BIND":          "1.2.3.4",
			"INGRESS_POLL_INTERVAL": "10m",
		},
		expectedConfig: types.Configuration{
			Redis:        "some-redis:1234",
			Bind:         "1.2.3.4",
			PollInterval: 10 * time.Minute,
		},
		expectedCommand: "collector",
		arguments:       []string{"anything", "collector"},
		success:         true,
	},
	{
		description:     "Invalid arguments",
		environment:     map[string]string{},
		expectedConfig:  types.Configuration{},
		expectedCommand: "",
		arguments:       []string{"anything", "something"},
		success:         false,
	},
	{
		description:     "Too few arguments",
		environment:     map[string]string{},
		expectedConfig:  types.Configuration{},
		expectedCommand: "",
		arguments:       []string{"anything"},
		success:         false,
	},
}

func setup(testCase cliTest) {
	currentConfig = types.Configuration{}
	appType = ""
	success = true

	for k, v := range testCase.environment {
		os.Setenv(k, v)
	}
}

func teardown(testCase cliTest) {
	for k, _ := range testCase.environment {
		os.Unsetenv(k)
	}
}

func TestCLI(t *testing.T) {
	for _, testCase := range testCases {
		setup(testCase)

		subject.GetConfig(testCase.arguments)

		if currentConfig != testCase.expectedConfig {
			t.Errorf("'%s': Config did not match - expected %+v, got %+v", testCase.description, testCase.expectedConfig, currentConfig)
		}

		if appType != testCase.expectedCommand {
			t.Errorf("'%s': Expected app to be a server, got %v", testCase.description, appType)
		}

		if success != testCase.success {
			t.Errorf("'%s': Expected app success to be %v, got %v", testCase.description, testCase.success, success)
		}

		teardown(testCase)
	}
}
