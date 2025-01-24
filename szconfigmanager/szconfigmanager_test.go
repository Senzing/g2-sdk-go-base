package szconfigmanager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/fileutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultTruncation = 2500
	instanceName      = "SzConfigManager Test"
	observerOrigin    = "SzConfigManager observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badConfigDefinition       = "\n\t"
	badConfigID               = int64(0)
	badCurrentDefaultConfigID = int64(0)
	badLogLevelName           = "BadLogLevelName"
	badNewDefaultConfigID     = int64(0)
)

// Nil/empty parameters

var (
	nilConfigComment          string
	nilConfigDefinition       string
	nilConfigID               int64
	nilCurrentDefaultConfigID int64
	nilNewDefaultConfigID     int64
)

var (
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigManagerSingleton *Szconfigmanager
	szConfigSingleton        *szconfig.Szconfig
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzconfigmanager_AddConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	now := time.Now()
	szConfig, err := getSzConfig(ctx)
	require.NoError(test, err)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), baseTen)
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_AddConfig_badConfigDefinition(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	now := time.Now()
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	_, err := szConfigManager.AddConfig(ctx, badConfigDefinition, configComment)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_AddConfig_nilConfigDefinition(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	now := time.Now()
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	_, err := szConfigManager.AddConfig(ctx, nilConfigDefinition, configComment)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_AddConfig_nilConfigComment(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	now := time.Now()
	szConfig, err := getSzConfig(ctx)
	require.NoError(test, err)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), baseTen)
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	actual, err := szConfigManager.AddConfig(ctx, configDefinition, nilConfigComment)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_CreateNewConfigAddDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	timestamp := time.Now().UTC().Format(time.RFC3339)
	configComment := fmt.Sprintf("TestSzconfigmanager_CreateNewConfigAddDataSources: %s", timestamp)
	dataSource1 := "TEST_1"
	dataSource2 := "TEST_2"
	dataSource3 := "TEST_3"
	actual, err := szConfigManager.CreateNewConfigAddDataSources(ctx, 0, configComment, dataSource1, dataSource2, dataSource3)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_CreateNewConfigAddDataSources_Again(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	timestamp := time.Now().UTC().Format(time.RFC3339)
	configComment := fmt.Sprintf("TestSzconfigmanager_CreateNewConfigAddDataSources_Again: %s", timestamp)
	dataSource1 := "TEST_1"
	dataSource2 := "TEST_2"
	dataSource3 := "TEST_3"
	actual, err := szConfigManager.CreateNewConfigAddDataSources(ctx, 0, configComment, dataSource1, dataSource2, dataSource3)
	require.NoError(test, err)
	printActual(test, actual)
	datasources, err := szConfigManager.GetDataSources(ctx, actual)
	require.NoError(test, err)
	printActual(test, datasources)
	configs, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, configs)
}

func TestSzconfigmanager_CreateNewConfigAddDataSources_Multi(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	now := time.Now().UTC()
	timestamp := now.Format(time.RFC3339)
	datasourceSuffix := strconv.FormatInt(now.Unix(), baseTen)
	configComment := fmt.Sprintf("TestSzconfigmanager_CreateNewConfigAddDataSources_Multi: %s", timestamp)

	datasources, err := szConfigManager.GetDataSources(ctx, 0)
	require.NoError(test, err)
	printActual(test, datasources)

	defaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	require.NoError(test, err)
	printActual(test, defaultConfigID)

	datasources, err = szConfigManager.GetDataSources(ctx, defaultConfigID)
	require.NoError(test, err)
	printActual(test, datasources)

	newConfig1, err := szConfigManager.CreateNewConfigAddDataSources(ctx, 0, configComment, "THIS"+datasourceSuffix, "WAS"+datasourceSuffix, "HERE")
	require.NoError(test, err)
	printActual(test, newConfig1)
	datasources, err = szConfigManager.GetDataSources(ctx, newConfig1)
	require.NoError(test, err)
	printActual(test, datasources)

	newConfig2, err := szConfigManager.CreateNewConfigAddDataSources(ctx, newConfig1, configComment, "THIS"+datasourceSuffix, "WAS"+datasourceSuffix, "HERE")
	require.NoError(test, err)
	printActual(test, newConfig2)
	datasources, err = szConfigManager.GetDataSources(ctx, newConfig2)
	require.NoError(test, err)
	printActual(test, datasources)

	newConfig3, err := szConfigManager.CreateNewConfigAddDataSources(ctx, newConfig2, configComment, "THIS"+datasourceSuffix, "WAS"+datasourceSuffix, "HERE")
	require.NoError(test, err)
	printActual(test, newConfig3)
	datasources, err = szConfigManager.GetDataSources(ctx, newConfig3)
	require.NoError(test, err)
	printActual(test, datasources)

	newConfig4, err := szConfigManager.CreateNewConfigAddDataSources(ctx, newConfig2, "Alternate comment", "THIS"+datasourceSuffix, "WAS"+datasourceSuffix, "HERE")
	require.NoError(test, err)
	printActual(test, newConfig4)
	datasources, err = szConfigManager.GetDataSources(ctx, newConfig4)
	require.NoError(test, err)
	printActual(test, datasources)

	configs, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, configs)
}

func TestSzconfigmanager_GetConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	configID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}
	actual, err := szConfigManager.GetConfig(ctx, configID)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_GetConfig_badConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetConfig(ctx, badConfigID)
	assert.Equal(test, "", actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_GetConfig_nilConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetConfig(ctx, nilConfigID)
	assert.Equal(test, "", actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_GetConfigs(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetDataSources(ctx, 0)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_GetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetDefaultConfigID(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_GetTemplateConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetDataSources(ctx, 0)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_ReplaceDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	currentDefaultConfigID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}

	// TODO: This is kind of a cheater.

	newDefaultConfigID, err2 := szConfigManager.GetDefaultConfigID(ctx)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()-2")
	}

	err := szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)
	require.NoError(test, err)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badCurrentDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	newDefaultConfigID, err2 := szConfigManager.GetDefaultConfigID(ctx)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}
	err := szConfigManager.ReplaceDefaultConfigID(ctx, badCurrentDefaultConfigID, newDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzReplaceConflict)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badNewDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	currentDefaultConfigID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}
	err := szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, badNewDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_nilCurrentDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	newDefaultConfigID, err2 := szConfigManager.GetDefaultConfigID(ctx)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}
	err := szConfigManager.ReplaceDefaultConfigID(ctx, nilCurrentDefaultConfigID, newDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzReplaceConflict)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_nilNewDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	currentDefaultConfigID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}
	err := szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, nilNewDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_SetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	configID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}
	err := szConfigManager.SetDefaultConfigID(ctx, configID)
	require.NoError(test, err)
}

func TestSzconfigmanager_SetDefaultConfigID_badConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.SetDefaultConfigID(ctx, badConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_SetDefaultConfigID_nilConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.SetDefaultConfigID(ctx, nilConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func TestSzconfigmanager_getByteArray(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	szProduct.getByteArray(10)
}

func TestSzconfigmanager_getByteArrayC(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	szProduct.getByteArrayC(10)
}

func TestSzconfigmanager_newError(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	err := szProduct.newError(ctx, 1)
	require.Error(test, err)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfigmanager_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	_ = szConfigManager.SetLogLevel(ctx, badLogLevelName)
}

func TestSzconfigmanager_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
}

func TestSzconfigmanager_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	actual := szConfigManager.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzconfigmanager_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfigmanager_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getSzConfigManagerAsInterface(ctx)
	actual, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_Initialize(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	settings, err := getSettings()
	require.NoError(test, err)
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzconfigmanager_Initialize_error
// func TestSzconfigmanager_Initialize_error(test *testing.T) {}

func TestSzconfigmanager_Destroy(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.Destroy(ctx)
	require.NoError(test, err)
}

func TestSzconfigmanager_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerSingleton = nil
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzconfigmanager_Destroy_error
// func TestSzconfigmanager_Destroy_error(test *testing.T) {}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getDatabaseTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

func getSettings() (string, error) {
	var result string

	// Determine Database URL.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		return result, fmt.Errorf("failed to make target database path (%s) absolute. Error: %w", dbTargetPath, err)
	}
	databaseURL := fmt.Sprintf("sqlite3://na:na@nowhere/%s", dbTargetPath)

	// Create Senzing engine configuration JSON.

	configAttrMap := map[string]string{"databaseUrl": databaseURL}
	result, err = settings.BuildSimpleSettingsUsingMap(configAttrMap)
	if err != nil {
		return result, fmt.Errorf("failed to BuildSimpleSettingsUsingMap(%s) Error: %w", configAttrMap, err)
	}
	return result, err
}

func getSzConfig(ctx context.Context) (senzing.SzConfig, error) {
	var err error
	if szConfigSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szConfigSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		szConfigSingleton = &szconfig.Szconfig{}
		err = szConfigSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			return szConfigSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
		}
		if logLevel == "TRACE" {
			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szConfigSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				return szConfigSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
			}
			err = szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				return szConfigSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
			}
		}
		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			return szConfigSingleton, fmt.Errorf("Initialize() Error: %w", err)
		}
	}
	return szConfigSingleton, err
}

func getSzConfigManager(ctx context.Context) (*Szconfigmanager, error) {
	var err error
	if szConfigManagerSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szConfigManagerSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		szConfigManagerSingleton = &Szconfigmanager{}
		err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			return szConfigManagerSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
		}
		if logLevel == "TRACE" {
			szConfigManagerSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szConfigManagerSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				return szConfigManagerSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
			}
			err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				return szConfigManagerSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
			}
		}
		err = szConfigManagerSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			return szConfigManagerSingleton, fmt.Errorf("Initialize() Error: %w", err)
		}
	}
	return szConfigManagerSingleton, err
}

func getSzConfigManagerAsInterface(ctx context.Context) senzing.SzConfigManager {
	result, err := getSzConfigManager(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfigmanager")
}

func getTestObject(ctx context.Context, test *testing.T) *Szconfigmanager {
	result, err := getSzConfigManager(ctx)
	require.NoError(test, err)
	return result
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		if errors.Is(err, szerror.ErrSzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzBadInput) {
			fmt.Printf("\nBad user input error detected. \n\n")
		}
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error
	err = setupDirectories()
	if err != nil {
		return fmt.Errorf("Failed to set up directories. Error: %w", err)
	}
	err = setupDatabase()
	if err != nil {
		return fmt.Errorf("Failed to set up database. Error: %w", err)
	}
	err = setupSenzingConfiguration()
	if err != nil {
		return fmt.Errorf("Failed to set up Senzing configuration. Error: %w", err)
	}
	return err
}

func setupDatabase() error {
	var err error

	// Locate source and target paths.

	testDirectoryPath := getTestDirectoryPath()
	dbTargetPath, err := filepath.Abs(filepath.Join(testDirectoryPath, "G2C.db"))
	if err != nil {
		return fmt.Errorf("failed to make target database path (%s) absolute. Error: %w",
			dbTargetPath, err)
	}
	databaseTemplatePath, err := filepath.Abs(getDatabaseTemplatePath())
	if err != nil {
		return fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			databaseTemplatePath, err.Error())
	}

	// Copy template file to test directory.

	_, _, err = fileutil.CopyFile(databaseTemplatePath, testDirectoryPath, true) // Copy the SQLite database file.
	if err != nil {
		return fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
			databaseTemplatePath, testDirectoryPath, err)
	}
	return err
}

func setupDirectories() error {
	var err error
	testDirectoryPath := getTestDirectoryPath()
	err = os.RemoveAll(filepath.Clean(testDirectoryPath)) // cleanup any previous test run
	if err != nil {
		return fmt.Errorf("failed to remove target test directory (%v): %w", testDirectoryPath, err)
	}
	err = os.MkdirAll(filepath.Clean(testDirectoryPath), 0750) // recreate the test target directory
	if err != nil {
		return fmt.Errorf("failed to recreate target test directory (%v): %w", testDirectoryPath, err)
	}
	return err
}

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()

	// Create sz objects.

	settings, err := getSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings. Error: %w", err)
	}
	szConfig := &szconfig.Szconfig{}
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return fmt.Errorf("failed to szConfig.Initialize(). Error: %w", err)
	}
	defer func() { handleError(szConfig.Destroy(ctx)) }()

	szConfigManager := &Szconfigmanager{}
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		return fmt.Errorf("failed to szConfigManager.Initialize(). Error: %w", err)
	}
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()

	// Create an in memory Senzing configuration.

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to szConfig.CreateConfig(). Error: %w", err)
	}

	// Add data sources to in-memory Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return fmt.Errorf("failed to szConfig.AddDataSource(). Error: %w", err)
		}
	}

	// Create a string representation of the in-memory configuration.

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return fmt.Errorf("failed to szConfig.ExportConfig(). Error: %w", err)
	}

	// Close szConfig in-memory object.

	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		return fmt.Errorf("failed to szConfig.CloseConfig(). Error: %w", err)
	}

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szconfigmanager_test at %s", now.UTC())
	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return fmt.Errorf("failed to szConfigManager.AddConfig(). Error: %w", err)
	}

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return fmt.Errorf("failed to szConfigManager.SetDefaultConfigID(). Error: %w", err)
	}

	return err
}

func teardown() error {
	ctx := context.TODO()
	err := teardownSzConfig(ctx)
	if err != nil {
		return err
	}
	err = teardownSzConfigManager(ctx)
	if err != nil {
		return err
	}
	return err
}

func teardownSzConfig(ctx context.Context) error {
	err := szConfigSingleton.UnregisterObserver(ctx, observerSingleton)
	if err != nil {
		return err
	}
	err = szConfigSingleton.Destroy(ctx)
	if err != nil {
		return err
	}
	szConfigSingleton = nil
	return nil
}

func teardownSzConfigManager(ctx context.Context) error {
	err := szConfigManagerSingleton.UnregisterObserver(ctx, observerSingleton)
	if err != nil {
		return err
	}
	err = szConfigManagerSingleton.Destroy(ctx)
	if err != nil {
		return err
	}
	szConfigManagerSingleton = nil
	return nil
}
