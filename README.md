# UCode SDK
UCode SDK is a Go package that provides a simple and efficient way to interact with the UCode API. This SDK offers various methods to perform CRUD operations, manage relationships, and handle data retrieval from the UCode platform.

## Table of Contents

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Usage](#usage)
   - [Creating Objects](#creating-objects)
   - [Retrieving Objects](#retrieving-objects)
   - [Updating Objects](#updating-objects)
   - [Deleting Objects](#deleting-objects)
4. [Error Handling](#error-handling)
5. [Examples](#examples)

## Installation

To install the UCode SDK, use the following command:

```bash
go get github.com/ucode-io/ucode_sdk
```

## Configuration

Before using the SDK, you need to configure it with your UCode API credentials and settings.

```go
import "github.com/ucode-io/ucode_sdk"

config := &ucodesdk.Config{
    BaseURL:        "https://api.admin.u-code.io",
    FunctionName:   "your-function-name",
    AppId: "your_app_id"
}

// Create a new UCode API client
ucodeApi := ucodesdk.NewSDK(config)
```

Make sure to set the `APP_ID` environment variable before running your application.

## Usage

### Creating Objects

To create a new object in a specific table:

```go
createRequest :=  map[string]any{
        "name":  "Example Object",
        "price": 100,
}

createdObject, response, err := ucodeApi.Items("your_table_slug").Create(createRequest).Exec()
if err != nil {
    log.Fatalf("Error creating object: %v", err)
}

fmt.Printf("Created object: %+v\n", createdObject)
```

### Retrieving Objects

#### Get List Slim

```go

objectList, response, err := ucodeApi.Items("your_table_slug").
    GetList().
    Page(1).
    Limit(10).
    Filter(map[string]any{}). // add any filters here
    WithRelation(true).
    Exec()
if err != nil {
    log.Fatalf("Error retrieving object list: %v", err)
}

fmt.Printf("Retrieved objects: %+v\n", objectList)
```

#### Get Single Slim

To retrieve a single object with selected relations:

```go

singleSlimObject, response, err := ucodeApi.Items("your_table_slug").
    GetSingle("object_guid").
    ExecSlim() // Use Exec() for simple
if err != nil {
    log.Fatalf("Error retrieving single slim object: %v", err)
}

fmt.Printf("Retrieved slim object: %+v\n", singleSlimObject)
```

#### Get List Aggregation

To perform an aggregation query (MongoDB only):

```go
aggregationPipeline := []map[string]any{
    {
        "$match": map[string]any{
            "field": map[string]any{
                "$exists": true,
                "$eq":     "value",
            },
        },
    },
    {
        "$group": map[string]any{
            "_id": "$group_field",
            "count": map[string]any{
                "$sum": 1,
            },
        },
    },
}

aggregationRequest :=  map[string]any{
    "pipelines": aggregationPipeline,
}

aggregationResult, response, err := ucodeApi.Items("your_table_slug").
    GetList().
    Pipelines(aggregationRequest).
    ExecAggregation()
if err != nil {
    log.Fatalf("Error performing aggregation: %v", err)
}

fmt.Printf("Aggregation result: %+v\n", aggregationResult)
```

### Updating Objects

#### Update Single Object

```go
updateRequest := map[string]any{
    "guid":  "object_guid",
    "name":  "Updated Name",
    "price": 150,
}

updatedObject, response, err := ucodeApi.Items("your_table_slug").
    Update(updatedObject).
    DisableFaas(false). //default true
    ExecSingle()
if err != nil {
    log.Fatalf("Error updating object: %v", err)
}

fmt.Printf("Updated object: %+v\n", updatedObject)
```

#### Update Multiple Objects

```go
multiUpdateRequest := map[string]any{
    "objects": []map[string]any{
        {"guid": "object1_guid", "name": "Updated Name 1"},
        {"guid": "object2_guid", "name": "Updated Name 2"},
    },
}

updatedObjects, response, err := ucodeApi.Items("your_table_slug").
    Update(updatedObject).
    ExecMultiple()
if err != nil {
    log.Fatalf("Error updating multiple objects: %v", err)
}

fmt.Printf("Updated objects: %+v\n", updatedObjects)
```

### Deleting Objects

#### Delete Single Object

```go

response, err := ucodeApi.Items("your_table_slug").
    Delete().
    DisableFaas(false).
    Single("object_guid").
    Exec()
if err != nil {
    log.Fatalf("Error deleting object: %v", err)
}

fmt.Printf("Delete response: %+v\n", response)
```

#### Delete Multiple Objects

```go

response, err := ucodeApi.Items("your_table_slug").
    Delete().
    Multiple([]string{"object1_guid", "object2_guid"}).
    Exec()
if err != nil {
    log.Fatalf("Error deleting multiple objects: %v", err)
}

fmt.Printf("Multiple delete response: %+v\n", response)
```

## Error Handling

All methods in the SDK return an error as the last return value. Always check for errors and handle them appropriately in your application.

```go
if err != nil {
    log.Printf("An error occurred: %v", err)
    // Handle the error (e.g., retry the operation, log it, or return it to the user)
}
```

## Examples

For more detailed examples and use cases, please refer to the `function_test.go` file in the SDK repository. This file contains comprehensive test cases that demonstrate how to use various features of the SDK.

---

For any issues, feature requests, or questions, please open an issue in the GitHub repository or contact the maintainers.