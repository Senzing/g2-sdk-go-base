//go:build linux

package szconfigmanager_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/helper"
	"github.com/senzing-garage/sz-sdk-go-core/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-core/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

const (
	instanceName   = "SzConfigManager Test"
	observerOrigin = "SzConfigManager observer"
	verboseLogging = senzing.SzNoLogging
)

var (
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigManagerSingleton *szconfigmanager.Szconfigmanager
	// szConfigSingleton        *szconfig.Szconfig
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzconfigmanager_AddConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
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
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configComment := "Example configuration"
	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		handleError(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_CreateNewConfigAddDataSources() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configComment := "Example configuration"
	configID, err := szConfigManager.CreateNewConfigAddDataSources(ctx, 0, configComment, "TEST_DATASOURCE", "ANOTHER_DATASOURCE")
	if err != nil {
		handleError(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_GetConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}
	configDefinition, err := szConfigManager.GetConfig(ctx, configID)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(configDefinition, 6))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_CLASS":"ADDRESS","ATTR_CODE":"ADDR_CITY",...
}

func ExampleSzconfigmanager_GetConfigs() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configList, err := szConfigManager.GetConfigs(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(configList, 3))
	// Output: {"CONFIGS":[{...
}

func ExampleSzconfigmanager_GetDataSources() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	dataSources, err := szConfigManager.GetDataSources(ctx, 0)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(dataSources, 5))
	// Output: {"DATA_SOURCES":[{"DSRC_CODE":"CUSTOMERS","DSRC_ID":1001...
}

func ExampleSzconfigmanager_GetDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_GetTemplateConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configID, err := szConfigManager.GetTemplateConfigID(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_ReplaceDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
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
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}
	configComment := "Example configuration"
	newDefaultConfigID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		handleError(err)
	}
	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)
	if err != nil {
		_ = err
	}
	// Output:
}

func ExampleSzconfigmanager_SetDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configID, err := szConfigManager.GetDefaultConfigID(ctx) // For example purposes only. Normally would use output from GetConfigList()
	if err != nil {
		handleError(err)
	}
	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		handleError(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzconfigmanager_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager, err := getSzConfigManagerCore(ctx)
	if err != nil {
		handleError(err)
	}
	err = szConfigManager.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzconfigmanager_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager, err := getSzConfigManagerCore(ctx)
	if err != nil {
		handleError(err)
	}
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzconfigmanager_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager, err := getSzConfigManagerCore(ctx)
	if err != nil {
		handleError(err)
	}
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	result := szConfigManager.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func getTestDirectoryPath() string {
	return filepath.FromSlash("../target/test/szconfigmanager")
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

func getSzAbstractFactoryCore(ctx context.Context) (*szabstractfactory.Szabstractfactory, error) {
	var err error
	var result szabstractfactory.Szabstractfactory
	_ = ctx
	settings, err := getSettings()
	if err != nil {
		return &result, err
	}
	result = szabstractfactory.Szabstractfactory{
		ConfigID:       senzing.SzInitializeWithDefaultConfiguration,
		InstanceName:   instanceName,
		Settings:       settings,
		VerboseLogging: verboseLogging,
	}
	return &result, err
}

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	result, err := getSzAbstractFactoryCore(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

// func getSzConfigCore(ctx context.Context) (senzing.SzConfig, error) {
// 	var err error
// 	if szConfigSingleton == nil {
// 		settings, err := getSettings()
// 		if err != nil {
// 			return szConfigSingleton, fmt.Errorf("getSettings() Error: %w", err)
// 		}
// 		szConfigSingleton = &szconfig.Szconfig{}
// 		err = szConfigSingleton.SetLogLevel(ctx, logLevel)
// 		if err != nil {
// 			return szConfigSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
// 		}
// 		if logLevel == "TRACE" {
// 			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)
// 			err = szConfigSingleton.RegisterObserver(ctx, observerSingleton)
// 			if err != nil {
// 				return szConfigSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
// 			}
// 			err = szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
// 			if err != nil {
// 				return szConfigSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
// 			}
// 		}
// 		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
// 		if err != nil {
// 			return szConfigSingleton, fmt.Errorf("Initialize() Error: %w", err)
// 		}
// 	}
// 	return szConfigSingleton, err
// }

func getSzConfigManagerCore(ctx context.Context) (*szconfigmanager.Szconfigmanager, error) {
	var err error
	if szConfigManagerSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szConfigManagerSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{}
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

func handleError(err error) {
	fmt.Println(err)
}
