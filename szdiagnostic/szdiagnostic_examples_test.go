//go:build linux

package szdiagnostic_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

const (
	instanceName   = "SzDiagnostic Test"
	verboseLogging = senzing.SzNoLogging
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_CheckDatastorePerformance() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}
	secondsToRun := 1
	result, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(result, 2))
	// Output: {"insertTime":1000,...
}

func ExampleSzdiagnostic_GetDatastoreInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}
	result, err := szDiagnostic.GetDatastoreInfo(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"dataStores":[{"id":"CORE","type":"sqlite3","location":"nowhere"}]}
}

func ExampleSzdiagnostic_GetFeature() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}
	featureID := int64(1)
	result, err := szDiagnostic.GetFeature(ctx, featureID)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"LIB_FEAT_ID":1,"FTYPE_CODE":"NAME","ELEMENTS":[{"FELEM_CODE":"FULL_NAME","FELEM_VALUE":"Robert Smith"},{"FELEM_CODE":"SUR_NAME","FELEM_VALUE":"Smith"},{"FELEM_CODE":"GIVEN_NAME","FELEM_VALUE":"Robert"},{"FELEM_CODE":"CULTURE","FELEM_VALUE":"ANGLO"},{"FELEM_CODE":"CATEGORY","FELEM_VALUE":"PERSON"},{"FELEM_CODE":"TOKENIZED_NM","FELEM_VALUE":"ROBERT|SMITH"}]}
}

func ExampleSzdiagnostic_PurgeRepository() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}
	err = szDiagnostic.PurgeRepository(ctx)
	if err != nil {
		handleError(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzdiagnostic_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzdiagnostic_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	result := szDiagnostic.GetObserverOrigin(ctx)
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

func getSzDiagnostic(ctx context.Context) *szdiagnostic.Szdiagnostic {
	_ = ctx
	settings, err := getSettings()
	if err != nil {
		panic(err)
	}
	result := &szdiagnostic.Szdiagnostic{}
	err = result.Initialize(ctx, instanceName, settings, 0, verboseLogging)
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
