package g2config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go/g2api"
	g2configapi "github.com/senzing/g2-sdk-go/g2config"
	"github.com/senzing/g2-sdk-go/g2error"
	futil "github.com/senzing/go-common/fileutil"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-logging/logging"
	"github.com/stretchr/testify/assert"
)

const (
	defaultTruncation = 76
	printResults      = false
	moduleName        = "Config Test Module"
	verboseLogging    = 0
)

var (
	configInitialized bool     = false
	globalG2config    G2config = G2config{}
	logger            logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return g2error.Cast(logger.NewError(errorId, err), err)
}

func getTestObject(ctx context.Context, test *testing.T) g2api.G2config {
	return &globalG2config
}

func getG2Config(ctx context.Context) g2api.G2config {
	return &globalG2config
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, g2config g2api.G2config, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func baseDirectoryPath() string {
	return filepath.FromSlash("../target/test/g2config")
}

func dbTemplatePath() string {
	return filepath.FromSlash("../testdata/sqlite/G2C.db")
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
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
	var err error = nil
	ctx := context.TODO()
	logger, err = logging.NewSenzingSdkLogger(ComponentId, g2configapi.IdMessages)
	if err != nil {
		return createError(5901, err)
	}

	baseDir := baseDirectoryPath()
	err = os.RemoveAll(filepath.Clean(baseDir)) // cleanup any previous test run
	if err != nil {
		return fmt.Errorf("Failed to remove target test directory (%v): %w", baseDir, err)
	}
	err = os.MkdirAll(filepath.Clean(baseDir), 0770) // recreate the test target directory
	if err != nil {
		return fmt.Errorf("Failed to recreate target test directory (%v): %w", baseDir, err)
	}

	// get the database URL and determine if external or a local file just created
	dbUrl, _, err := setupDB(false)
	if err != nil {
		return err
	}

	// get the INI params
	iniParams, err := setupIniParams(dbUrl)
	if err != nil {
		return err
	}

	// setup the config
	err = setupG2config(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	// get the tem
	return err
}

func teardownG2config(ctx context.Context) error {
	// check if not initialized
	if !configInitialized {
		return nil
	}

	// destroy the G2config
	err := globalG2config.Destroy(ctx)
	if err != nil {
		return err
	}
	configInitialized = false

	return nil
}

func teardown() error {
	ctx := context.TODO()
	err := teardownG2config(ctx)
	return err
}

func restoreG2config(ctx context.Context) error {
	iniParams, err := getIniParams()
	if err != nil {
		return err
	}

	err = setupG2config(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return err
	}

	return nil
}

func setupDB(preserveDB bool) (string, bool, error) {
	var err error = nil

	// get the base directory
	baseDir := baseDirectoryPath()

	// get the template database file path
	dbFilePath := dbTemplatePath()

	dbFilePath, err = filepath.Abs(dbFilePath)
	if err != nil {
		err = fmt.Errorf("failed to obtain absolute path to database file (%s): %s",
			dbFilePath, err.Error())
		return "", false, err
	}

	// check the environment for a database URL
	dbUrl, envUrlExists := os.LookupEnv("SENZING_TOOLS_DATABASE_URL")

	dbTargetPath := filepath.Join(baseDirectoryPath(), "G2C.db")

	dbTargetPath, err = filepath.Abs(dbTargetPath)
	if err != nil {
		err = fmt.Errorf("failed to make target database path (%s) absolute: %w",
			dbTargetPath, err)
		return "", false, err
	}

	dbDefaultUrl := fmt.Sprintf("sqlite3://na:na@%s", dbTargetPath)

	dbExternal := envUrlExists && dbDefaultUrl != dbUrl

	if !dbExternal {
		// set the database URL
		dbUrl = dbDefaultUrl

		if !preserveDB {
			// copy the SQLite database file
			_, _, err := futil.CopyFile(dbFilePath, baseDir, true)

			if err != nil {
				err = fmt.Errorf("setup failed to copy template database (%v) to target path (%v): %w",
					dbFilePath, baseDir, err)
				// fall through to return the error
			}
		}
	}

	return dbUrl, dbExternal, err
}

func setupIniParams(dbUrl string) (string, error) {
	configAttrMap := map[string]string{"databaseUrl": dbUrl}

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(configAttrMap)

	if err != nil {
		err = createError(5902, err)
	}

	return iniParams, err
}

func getIniParams() (string, error) {
	dbUrl, _, err := setupDB(true)
	if err != nil {
		return "", err
	}
	iniParams, err := setupIniParams(dbUrl)
	if err != nil {
		return "", err
	}
	return iniParams, nil
}

func setupG2config(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	if configInitialized {
		return fmt.Errorf("G2config is already setup and has not been torn down")
	}
	globalG2config.SetLogLevel(ctx, logging.LevelInfoName)
	log.SetFlags(0)
	err := globalG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	configInitialized = true
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2config_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2config.SetObserverOrigin(ctx, origin)
}

func TestG2config_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2config.SetObserverOrigin(ctx, origin)
	actual := g2config.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestG2config_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	actual, err := g2config.AddDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2config_AddDataSource_WithLoad(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	jsonConfig, err := g2config.Save(ctx, configHandle)
	testError(test, ctx, g2config, err)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
	configHandle2, err := g2config.Load(ctx, jsonConfig)
	testError(test, ctx, g2config, err)
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	actual, err := g2config.AddDataSource(ctx, configHandle2, inputJson)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
	err = g2config.Close(ctx, configHandle2)
	testError(test, ctx, g2config, err)
}

func TestG2config_Close(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2config_Create(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	actual, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
}

func TestG2config_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	actual, err := g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "Original", actual)
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	_, err = g2config.AddDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)
	actual, err = g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "     Add", actual)
	err = g2config.DeleteDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)
	actual, err = g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "  Delete", actual)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2config_DeleteDataSource_WithLoad(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	actual, err := g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "Original", actual)
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	_, err = g2config.AddDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)
	actual, err = g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "     Add", actual)
	jsonConfig, err := g2config.Save(ctx, configHandle)
	testError(test, ctx, g2config, err)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
	configHandle2, err := g2config.Load(ctx, jsonConfig)
	testError(test, ctx, g2config, err)
	err = g2config.DeleteDataSource(ctx, configHandle2, inputJson)
	testError(test, ctx, g2config, err)
	actual, err = g2config.ListDataSources(ctx, configHandle2)
	testError(test, ctx, g2config, err)
	printResult(test, "  Delete", actual)
	err = g2config.Close(ctx, configHandle2)
	testError(test, ctx, g2config, err)
}

func TestG2config_ListDataSources(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	actual, err := g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2config_Load(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	jsonConfig, err := g2config.Save(ctx, configHandle)
	testError(test, ctx, g2config, err)
	actual, err := g2config.Load(ctx, jsonConfig)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
}

func TestG2config_Save(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	actual, err := g2config.Save(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
}

func TestG2config_Init(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, err := getIniParams()
	testError(test, ctx, g2config, err)
	err = g2config.Init(ctx, moduleName, iniParams, verboseLogging)
	testError(test, ctx, g2config, err)
}

func TestG2config_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	err := g2config.Destroy(ctx)
	testError(test, ctx, g2config, err)

	// restore the state that existed prior to this test
	configInitialized = false
	restoreG2config(ctx)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2config_SetObserverOrigin() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2config.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2config_GetObserverOrigin() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2config.SetObserverOrigin(ctx, origin)
	result := g2config.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2config_AddDataSource() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	result, err := g2config.AddDataSource(ctx, configHandle, inputJson)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DSRC_ID":1001}
}

func ExampleG2config_Close() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = g2config.Close(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_Create() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2config_DeleteDataSource() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	inputJson := `{"DSRC_CODE": "TEST"}`
	err = g2config.DeleteDataSource(ctx, configHandle, inputJson)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_ListDataSources() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	result, err := g2config.ListDataSources(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleG2config_Load() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	mockConfigHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := g2config.Save(ctx, mockConfigHandle)
	if err != nil {
		fmt.Println(err)
	}
	configHandle, err := g2config.Load(ctx, jsonConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2config_Save() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := g2config.Save(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(jsonConfig, 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"ADVANCED":"Yes","INTERNAL":"No"},...
}

func ExampleG2config_SetLogLevel() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	err := g2config.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	moduleName := "Test module name"
	iniParams, err := getIniParams()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := 0
	err = g2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2config_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	err := g2config.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
