//go:build linux

package szconfig_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szconfig"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

const (
	instanceName   = "SzConfig Test"
	verboseLogging = senzing.SzNoLogging
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzconfig_AddDataSource() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "GO_TEST"
	result, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"DSRC_ID":1001}
}

func ExampleSzconfig_CloseConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzconfig_CreateConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfig_DeleteDataSource() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "TEST"
	err = szConfig.DeleteDataSource(ctx, configHandle, dataSourceCode)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzconfig_ExportConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(configDefinition, 10))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_CLASS":"ADDRESS","ATTR_CODE":"ADDR_CITY","ATTR_ID":1608,"DEFAULT_VALUE":null,"FELEM_CODE":"CITY","FELEM_REQ":"Any",...
}

func ExampleSzconfig_GetDataSources() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	result, err := szConfig.GetDataSources(ctx, configHandle)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleSzconfig_ImportConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	mockConfigHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configDefinition, err := szConfig.ExportConfig(ctx, mockConfigHandle)
	if err != nil {
		handleError(err)
	}
	configHandle, err := szConfig.ImportConfig(ctx, configDefinition)
	if err != nil {
		handleError(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzconfig_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	err := szConfig.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzconfig_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzconfig_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfig/szconfig_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	result := szConfig.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

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

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	var err error
	var result senzing.SzAbstractFactory
	_ = ctx
	settings, err := getSettings()
	if err != nil {
		panic(err)
	}
	result = &szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}
	return result
}

func getSzConfig(ctx context.Context) *szconfig.Szconfig {
	_ = ctx
	settings, err := getSettings()
	if err != nil {
		panic(err)
	}
	result := &szconfig.Szconfig{}
	err = result.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		panic(err)
	}
	return result
}

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfig")
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}
