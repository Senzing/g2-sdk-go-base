# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], [markdownlint],
and this project adheres to [Semantic Versioning].

## [Unreleased]

-

## [0.8.8] - 2025-01-31

### Changed in 0.8.8

- Update dependencies
- Updated test cases and examples

## [0.8.7] - 2025-01-24

### Changed in 0.8.7

- Return empty string for non-withInfo methods

## [0.8.6] - 2024-12-10

### Changed in 0.8.6

- Update dependencies

## [0.8.5] - 2024-11-07

### Changed in 0.8.5

- Update dependencies

## [0.8.4] - 2024-10-30

### Changed in 0.8.4

- Migrate from `gohelpers` to `szhelpers`

## [0.8.3] - 2024-10-01

### Changed in 0.8.3

- Update dependencies
- Add `PreprocessRecord()`

## [0.8.2] - 2024-09-11

### Changed in 0.8.2

- Update dependencies
- Added test cases.

## [0.8.1] - 2024-08-27

### Changed in 0.8.1

- Modify method calls to match Senzing API 4.0.0-24237

## [0.8.0] - 2024-08-22

### Changed in 0.8.0

- Change from `g2` to `sz`/`er`

## [0.7.5] - 2024-08-02

### Changed in 0.7.5

- Update to `template-go`

## [0.7.4] - 2024-06-26

### Changed in 0.7.4

- Updated dependencies
- Included `parametertests` spike
- Synchronized with [sz-sdk-go-grpc](https://github.com/senzing-garage/sz-sdk-go-grpc) and [sz-sdk-go-mock](https://github.com/senzing-garage/sz-sdk-go-mock)

## [0.7.3] - 2024-06-13

### Changed in 0.7.3

- Updated methods to Senzing 4.0.0-24162

## [0.7.2] - 2024-05-31

### Changed in 0.7.2

- Message ID changed from `senzing-` to `SZSDK`

## [0.7.1] - 2024-05-08

### Added in 0.7.1

- `SzDiagnostic.GetFeature`
- `SzEngine.FindInterestingEntitiesByEntityId`
- `SzEngine.FindInterestingEntitiesByRecordId`

### Deleted in 0.7.1

- `SzEngine.GetRepositoryLastModifiedTime`

## [0.7.0] - 2024-04-26

### Changed in 0.7.0

- Move from `g2-sdk-go-base` to `sz-sdk-go-core`
- Updated dependencies
  - github.com/stretchr/testify v1.9.0
  - google.golang.org/grpc v1.63.2

## [0.6.1] - 2024-02-29

### Changed in 0.6.1

- Added `G2Diagnostic.PurgeRepository()`

## [0.6.0] - 2024-02-27

### Changed in 0.6.0

- Updated dependencies
- Deleted methods not used in V4

## [0.5.0] - 2024-01-26

### Changed in 0.5.0

- Renamed module to `github.com/senzing-garage/sz-sdk-go-core`
- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Update dependencies
  - github.com/senzing-garage/g2-sdk-go v0.9.0
  - google.golang.org/grpc v1.61.0

## [0.4.0] - 2024-01-02

### Changed in 0.4.0

- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Update dependencies
  - github.com/aquilax/truncate v1.0.0
  - github.com/senzing-garage/go-common v0.4.0
  - github.com/senzing-garage/go-logging v1.4.0
  - github.com/senzing-garage/go-observing v0.3.0
  - github.com/senzing/g2-sdk-go v0.8.0
  - github.com/stretchr/testify v1.8.4
  - google.golang.org/grpc v1.60.1

## [0.3.4] - 2023-12-12

### Added in 0.3.4

- `ExportCSVEntityReportIterator` and `ExportJSONEntityReportIterator`

### Changed in 0.3.4

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.7.6
  - google.golang.org/grpc v1.60.0

## [0.3.3] - 2023-10-31

### Changed in 0.3.3

- Support for changed method signatures in Senzing G2Config API
- Update dependencies
  - github.com/senzing-garage/go-common v0.3.2-0.20231018174900-c1895fb44c30

## [0.3.2] - 2023-10-18

### Changed in 0.3.2

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.7.4
  - github.com/senzing-garage/go-common v0.3.1
  - github.com/senzing-garage/go-logging v1.3.3
  - github.com/senzing-garage/go-observing v0.2.8
  - google.golang.org/grpc v1.59.0

## [0.3.1] - 2023-10-12

### Changed in 0.3.1

- Changed from `int` to `int64` where required by the SenzingAPI
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.7.3
  - google.golang.org/grpc v1.58.3

### Deleted in 0.3.1

- `g2product.ValidateLicenseFile`
- `g2product.ValidateLicenseStringBase64`

## [0.3.0] - 2023-10-01

### Changed in 0.3.0

- Support for SenzingAPI 3.8.0
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.7.0
  - google.golang.org/grpc v1.58.2

### Removed in 0.3.0

- In `g2diagnostic.`
  - CloseEntityListBySize
  - FetchNextEntityBySize
  - FindEntitiesByFeatureIDs
  - GetDataSourceCounts
  - GetEntityDetails
  - GetEntityListBySize
  - GetEntityResume
  - GetEntitySizeBreakdown
  - GetFeature
  - GetGenericFeatures
  - GetMappingStatistics
  - GetRelationshipDetails
  - GetResolutionStatistics
- In `g2engine.`
  - AddRecordWithInfoWithReturnedRecordID
  - AddRecordWithReturnedRecordID
  - CheckRecord
  - ProcessRedoRecord
  - ProcessRedoRecordWithInfo
  - ProcessWithResponse
  - ProcessWithResponseResize

## [0.2.7] - 2023-09-01

### Changed in 0.2.7

- Last version before SenzingAPI 3.8.0

## [0.2.6] - 2023-08-28

### Changed in 0.2.6

- Enablement for Windows

## [0.2.5] - 2023-08-24

### Changed in 0.2.5

- Changes to automated tests to isolate test suites from each other
- Changes to automated tests to isolate individual tests from each other

## [0.2.4] - 2023-08-08

### Changed in 0.2.4

- Changes to accomodate macOS builds -- cleanup tests for multiplatform differences
- Switched Linux github action workflow test to use Senzing Staging Repository

## [0.2.3] - 2023-08-07

### Changed in 0.2.3

- Refactor to `template-go`
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.8
  - github.com/senzing-garage/go-common v0.2.11
  - github.com/senzing-garage/go-logging v1.3.2
  - github.com/senzing-garage/go-observing v0.2.7
  - google.golang.org/grpc v1.57.0

## [0.2.2] - 2023-07-07

### Changed in 0.2.2

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.7
  - github.com/senzing-garage/go-common v0.1.4
  - github.com/senzing-garage/go-logging v1.2.6
  - google.golang.org/grpc v1.56.2

## [0.2.1] - 2023-06-16

### Changed in 0.2.1

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.5
  - github.com/senzing-garage/go-common v0.1.4
  - github.com/senzing-garage/go-logging v1.2.6
  - github.com/senzing-garage/go-observing v0.2.6
  - github.com/stretchr/testify v1.8.4
  - google.golang.org/grpc v1.56.0

## [0.2.0] - 2023-05-26

### Changed in 0.2.0

- Fixed method signature for g2config.Load()
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.4

## [0.1.11] - 2023-05-19

### Changed in 0.1.11

- Fixed go/CGO memory issues
- Update dependencies
  - github.com/senzing-garage/go-observing v0.2.5

## [0.1.10] - 2023-05-11

### Changed in 0.1.10

- Update dependencies
  - github.com/senzing-garage/go-common v0.1.3
  - github.com/senzing-garage/go-logging v1.2.3

## [0.1.9] - 2023-05-10

### Changed in 0.1.9

- Added GetObserverOrigin() and SetObserverOrigin() to g2* packages
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.2
  - github.com/senzing-garage/go-observing v0.2.2

## [0.1.8] - 2023-04-21

### Changed in 0.1.8

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.1

## [0.1.7] - 2023-04-20

### Changed in 0.1.7

- Updated dependencies
- Migrated from `github.com/senzing-garage/go-logging/logger` to `github.com/senzing-garage/go-logging/logging`

## [0.1.6] - 2023-04-18

### Fixed in 0.1.6

- Fixed concurrency issue with unregistering observer

## [0.1.5] - 2023-04-14

### Changed in 0.1.5

- Improved underlying CGO for g2engine

## [0.1.4] - 2023-03-27

### Changed in 0.1.4

- Fix copy/paste error in getRepositoryLastModifiedTime

## [0.1.3] - 2023-03-22

### Changed in 0.1.3

- Migrated to `github.com/senzing/g2-sdk-go v0.5.0`
- Refactored documentation

## [0.1.2] - 2023-03-22

### Changed in 0.1.2

- Fixed internal logic

## [0.1.1] - 2023-02-21

### Added to 0.1.1

- Modified `GetSdkId()` signature

## [0.1.0] - 2023-02-14

### Added to 0.1.0

- Initial functionality
-

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[markdownlint]: https://dlaa.me/markdownlint/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
