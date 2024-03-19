package goseeder

import (
	"database/sql"
	"flag"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type mockStatus struct {
	callCount  int
	callParams []interface{}
}

func newFnMock() (func(), *mockStatus) {
	mockStatus := mockStatus{
		callCount:  0,
		callParams: []interface{}{},
	}
	fn := func() {
		mockStatus.callCount++
	}

	return fn, &mockStatus
}

func newSeedFnMock() (func(s Seeder), *mockStatus) {
	mockStatus := mockStatus{
		callCount:  0,
		callParams: []interface{}{},
	}
	fnSeeder := func(s Seeder) {
		mockStatus.callCount++
	}

	return fnSeeder, &mockStatus
}

var conProvider = func() *sql.DB {
	db, _, _ := sqlmock.New()
	return db
}

func TestWithSeeder_not_seeding(t *testing.T) {
	clientMain, mockStatus := newFnMock()
	WithSeeder(conProvider, clientMain)
	require.Equal(t, 1, mockStatus.callCount)

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestWithSeeder_seed_all(t *testing.T) {
	clientMain, mockFnStatus := newFnMock()
	mockSeedFn, mockSeedStatus := newSeedFnMock()

	Registration{
		Name: "test_seed",
		Env:  "",
	}.Complete(mockSeedFn)

	os.Args = []string{
		"",
		"--gseed",
	}

	WithSeeder(conProvider, clientMain)
	require.Equal(t, 0, mockFnStatus.callCount)
	require.Equal(t, 1, mockSeedStatus.callCount)

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestWithSeeder_seed_named(t *testing.T) {
	clientMain, mockFnStatus := newFnMock()
	mockSeedFn, mockSeedStatus := newSeedFnMock()
	mockSeed2Fn, mockSeed2Status := newSeedFnMock()

	Registration{
		Name: "test_seed",
		Env:  "",
	}.Complete(mockSeedFn)
	Registration{
		Name: "another_weird_seed",
		Env:  "",
	}.Complete(mockSeed2Fn)

	os.Args = []string{
		"",
		"--gseed",
		"--gsnames=another_weird_seed",
	}

	WithSeeder(conProvider, clientMain)
	require.Equal(t, 0, mockFnStatus.callCount)
	require.Equal(t, 0, mockSeedStatus.callCount)
	require.Equal(t, 1, mockSeed2Status.callCount)

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestWithSeeder_seed_for_env(t *testing.T) {
	clientMain, mockFnStatus := newFnMock()
	mockCommonSeedFn, mockCommonSeedStatus := newSeedFnMock()
	mockSeedSecretFn, mockSeedSecretStatus := newSeedFnMock()
	mockSeedStageFn, mockSeedStageStatus := newSeedFnMock()

	Registration{
		Name: "common_seed",
		Env:  "",
	}.Complete(mockCommonSeedFn)
	Registration{
		Name: "secret_env_seed",
		Env:  "secret",
	}.Complete(mockSeedSecretFn)
	Registration{
		Name: "stage_env_seed",
		Env:  "stage",
	}.Complete(mockSeedStageFn)

	os.Args = []string{
		"",
		"--gseed",
		"--gsenv=secret",
	}

	WithSeeder(conProvider, clientMain)
	require.Equal(t, 0, mockFnStatus.callCount)
	require.Equal(t, 1, mockCommonSeedStatus.callCount)
	require.Equal(t, 1, mockSeedSecretStatus.callCount)
	require.Equal(t, 0, mockSeedStageStatus.callCount)

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestRegister(t *testing.T) {
	Register(dummySeeder)
	a := seeders[len(seeders)-1]
	require.Equal(t, "", a.env)
	require.Equal(t, "dummySeeder", a.name)
	require.NotEmpty(t, a.cb)
}

func TestRegisterForTest(t *testing.T) {
	RegisterForTest(dummySeeder)
	a := seeders[len(seeders)-1]
	require.Equal(t, "test", a.env)
	require.Equal(t, "dummySeeder", a.name)
	require.NotEmpty(t, a.cb)
}

func TestRegisterForEnv(t *testing.T) {
	RegisterForEnv("mySuperEnv", dummySeeder)
	a := seeders[len(seeders)-1]
	require.Equal(t, "mySuperEnv", a.env)
	require.Equal(t, "dummySeeder", a.name)
	require.NotEmpty(t, a.cb)
}

func TestRegisterForEnvNamed(t *testing.T) {
	RegisterForEnvNamed("mySuperEnv", dummySeeder, "customName")
	a := seeders[len(seeders)-1]
	require.Equal(t, "mySuperEnv", a.env)
	require.Equal(t, "customName", a.name)
	require.NotEmpty(t, a.cb)
}

func dummySeeder(s Seeder) {}
