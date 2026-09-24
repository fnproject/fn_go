package shim

import (
	"fmt"
	"github.com/fnproject/fn_go/clientv2/fns"
	"github.com/fnproject/fn_go/modelsv2"
	"github.com/fnproject/fn_go/provider/oracle/shim/client"
	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/oracle/oci-go-sdk/v65/functions"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateCodeOnlyFunctionSourceDetailsDirect(t *testing.T) {
	archiveBytes := []byte("zip-bytes")
	fn := &modelsv2.Fn{
		CodeOnly:          true,
		SourceType:        "direct",
		SourceArchive:     strfmt.Base64(archiveBytes),
		RuntimeConfigType: "FUNCTION_UPDATE",
		RuntimeName:       "python311.ol9",
		Handler:           "hello_world.handler",
	}

	details, err := createCodeOnlyFunctionSourceDetails(fn)
	assert.NoError(t, err)

	archiveDetails, ok := details.(functions.CreateArchiveFunctionSourceDetails)
	assert.True(t, ok)

	directDetails, ok := archiveDetails.ArchiveSourceDetails.(functions.CreateDirectArchiveSourceDetails)
	assert.True(t, ok)
	assert.Equal(t, archiveBytes, directDetails.ArchiveFile)

	runtimeConfig, ok := archiveDetails.RuntimeConfig.(functions.CreateFunctionUpdateRuntimeConfig)
	assert.True(t, ok)
	if assert.NotNil(t, runtimeConfig.FunctionsRuntimeName) {
		assert.Equal(t, "python311.ol9", *runtimeConfig.FunctionsRuntimeName)
	}
	if assert.NotNil(t, archiveDetails.Handler) {
		assert.Equal(t, "hello_world.handler", *archiveDetails.Handler)
	}
}

func TestCreateCodeOnlyFunctionSourceDetailsObjectStorageManual(t *testing.T) {
	fn := &modelsv2.Fn{
		CodeOnly:              true,
		SourceType:            "object-storage",
		SourceBucketName:      "code-only-test-files",
		SourceNamespace:       "oraclefunctionsdevelopm",
		SourceObjectName:      "hello.0.0.1.zip",
		SourceObjectVersionID: "object-version-id",
		RuntimeConfigType:     "MANUAL",
		RuntimeName:           "python311.ol9",
		RuntimeVersionID:      "ocid1.functionsruntimeversion.oc1..exampleuniqueID",
		Handler:               "hello_world.handler",
	}

	details, err := createCodeOnlyFunctionSourceDetails(fn)
	assert.NoError(t, err)

	archiveDetails, ok := details.(functions.CreateArchiveFunctionSourceDetails)
	assert.True(t, ok)

	objectStorageDetails, ok := archiveDetails.ArchiveSourceDetails.(functions.CreateObjectStorageArchiveSourceDetails)
	assert.True(t, ok)
	if assert.NotNil(t, objectStorageDetails.BucketName) {
		assert.Equal(t, "code-only-test-files", *objectStorageDetails.BucketName)
	}
	if assert.NotNil(t, objectStorageDetails.Namespace) {
		assert.Equal(t, "oraclefunctionsdevelopm", *objectStorageDetails.Namespace)
	}
	if assert.NotNil(t, objectStorageDetails.ObjectName) {
		assert.Equal(t, "hello.0.0.1.zip", *objectStorageDetails.ObjectName)
	}
	if assert.NotNil(t, objectStorageDetails.ObjectVersionId) {
		assert.Equal(t, "object-version-id", *objectStorageDetails.ObjectVersionId)
	}

	runtimeConfig, ok := archiveDetails.RuntimeConfig.(functions.CreateManualRuntimeConfig)
	assert.True(t, ok)
	if assert.NotNil(t, runtimeConfig.FunctionsRuntimeName) {
		assert.Equal(t, "python311.ol9", *runtimeConfig.FunctionsRuntimeName)
	}
	if assert.NotNil(t, runtimeConfig.FunctionsRuntimeVersionId) {
		assert.Equal(t, "ocid1.functionsruntimeversion.oc1..exampleuniqueID", *runtimeConfig.FunctionsRuntimeVersionId)
	}
	if assert.NotNil(t, archiveDetails.Handler) {
		assert.Equal(t, "hello_world.handler", *archiveDetails.Handler)
	}
}

func TestCreateArchiveSourceDetailsValidation(t *testing.T) {
	_, err := createArchiveSourceDetails(&modelsv2.Fn{SourceType: "direct"})
	assert.EqualError(t, err, "direct source requires archive bytes")

	_, err = createArchiveSourceDetails(&modelsv2.Fn{SourceType: "object-storage", SourceBucketName: "bucket"})
	assert.EqualError(t, err, "object-storage source requires bucket, namespace, and object name")

	_, err = createArchiveSourceDetails(&modelsv2.Fn{SourceType: "something-else"})
	assert.EqualError(t, err, `unsupported code-only source type "something-else"`)
}

func TestCreateArchiveSourceDetailsDirectFile(t *testing.T) {
	archivePath := t.TempDir() + "/function.zip"
	archiveBytes := []byte("zip-bytes-from-file")
	assert.NoError(t, os.WriteFile(archivePath, archiveBytes, 0600))

	details, err := createArchiveSourceDetails(&modelsv2.Fn{
		SourceType: "direct",
		SourceFile: archivePath,
	})
	assert.NoError(t, err)

	directDetails, ok := details.(functions.CreateDirectArchiveSourceDetails)
	assert.True(t, ok)
	assert.Equal(t, archiveBytes, directDetails.ArchiveFile)
}

func TestCreateCodeOnlyUpdateSourceDetailsDirect(t *testing.T) {
	archiveBytes := []byte("updated-zip-bytes")
	fn := &modelsv2.Fn{
		CodeOnly:          true,
		SourceType:        "direct",
		SourceArchive:     strfmt.Base64(archiveBytes),
		RuntimeConfigType: "FUNCTION_UPDATE",
		RuntimeName:       "python311.ol9",
		Handler:           "hello_world.handler",
	}

	details, err := createCodeOnlyUpdateSourceDetails(fn)
	assert.NoError(t, err)

	archiveDetails, ok := details.(functions.UpdateArchiveFunctionSourceDetails)
	assert.True(t, ok)

	directDetails, ok := archiveDetails.ArchiveSourceDetails.(functions.UpdateDirectArchiveSourceDetails)
	assert.True(t, ok)
	assert.Equal(t, archiveBytes, directDetails.ArchiveFile)

	runtimeConfig, ok := archiveDetails.RuntimeConfig.(functions.UpdateFunctionUpdateRuntimeConfig)
	assert.True(t, ok)
	if assert.NotNil(t, runtimeConfig.FunctionsRuntimeName) {
		assert.Equal(t, "python311.ol9", *runtimeConfig.FunctionsRuntimeName)
	}
	if assert.NotNil(t, archiveDetails.Handler) {
		assert.Equal(t, "hello_world.handler", *archiveDetails.Handler)
	}
}

func TestCreateUpdateRuntimeConfigManualAllowsOmittedRuntimeName(t *testing.T) {
	config, err := createUpdateRuntimeConfig(&modelsv2.Fn{
		RuntimeConfigType: "MANUAL",
		RuntimeVersionID:  "ocid1.functionsruntimeversion.oc1..exampleuniqueID",
	})
	assert.NoError(t, err)

	manualConfig, ok := config.(functions.UpdateManualRuntimeConfig)
	assert.True(t, ok)
	assert.Nil(t, manualConfig.FunctionsRuntimeName)
	if assert.NotNil(t, manualConfig.FunctionsRuntimeVersionId) {
		assert.Equal(t, "ocid1.functionsruntimeversion.oc1..exampleuniqueID", *manualConfig.FunctionsRuntimeVersionId)
	}
}

func TestCreateFn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	timeout := int32(30)
	fn := modelsv2.Fn{
		Name:    "CreateFnName",
		AppID:   "CreateFnAppId",
		Memory:  128,
		Timeout: &timeout,
		Image:   "CreateFnImage",
		Annotations: map[string]interface{}{
			annotationImageDigest: "CreateFnDigest",
		},
		Config: map[string]string{
			"CreateFnKey": "CreateFnValue",
		},
	}

	createFnOK, err := shim.CreateFn(&fns.CreateFnParams{
		Body: &fn,
	})
	assert.NoError(t, err)

	result := createFnOK.GetPayload()

	expectedAnnotations := fn.Annotations
	expectedAnnotations[annotationCompartmentId] = "CreateFunctionCompartment"
	expectedAnnotations[annotationInvokeEndpoint] = fmt.Sprintf("CreateFunctionInvokeEndpoint/20181201/functions/%s/actions/invoke", result.ID)
	expectedAnnotations[annotationPCStrategy] = "NONE"

	assert.Equal(t, fn.Name, result.Name)
	assert.Equal(t, fn.AppID, result.AppID)
	assert.Equal(t, fn.Memory, result.Memory)
	assert.Equal(t, fn.Timeout, result.Timeout)
	assert.Equal(t, fn.Image, result.Image)
	assert.Equal(t, expectedAnnotations, result.Annotations)
	assert.Equal(t, fn.Config, result.Config)
	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
}

func TestDeleteFn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	_, err := shim.DeleteFn(&fns.DeleteFnParams{
		FnID: "DeleteFnId",
	})
	assert.NoError(t, err)
}

func TestGetFn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	fnId := "GetFnId"
	getFnOK, err := shim.GetFn(&fns.GetFnParams{
		FnID: fnId,
	})
	assert.NoError(t, err)

	result := getFnOK.GetPayload()
	assert.Equal(t, fnId, result.ID)
	assert.NotEmpty(t, result.AppID)
	assert.NotEmpty(t, result.Name)
	assert.NotEmpty(t, result.Memory)
	assert.NotEmpty(t, result.Timeout)
	assert.NotEmpty(t, result.Image)
	assert.NotEmpty(t, result.Annotations[annotationImageDigest])
	assert.NotEmpty(t, result.Annotations[annotationInvokeEndpoint])
	assert.NotEmpty(t, result.Annotations[annotationCompartmentId])
	assert.NotEmpty(t, result.Config)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
}

func TestListFnsMultiplePages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	appId := "ListFnsAppId"
	params := &fns.ListFnsParams{AppID: &appId}
	var results []*modelsv2.Fn
	for {
		listFnsOK, err := shim.ListFns(params)
		assert.NoError(t, err)

		result := listFnsOK.GetPayload()
		results = append(results, result.Items...)

		if result.NextCursor == "" {
			break
		}

		params.Cursor = &result.NextCursor
	}
	assert.Len(t, results, 9)
	fn := results[0]
	assert.Equal(t, appId, fn.AppID)
	assert.NotEmpty(t, fn.ID)
	assert.NotEmpty(t, fn.Name)
	assert.NotEmpty(t, fn.Memory)
	assert.NotEmpty(t, fn.Timeout)
	assert.NotEmpty(t, fn.Image)
	assert.NotEmpty(t, fn.Annotations[annotationImageDigest])
	assert.NotEmpty(t, fn.Annotations[annotationInvokeEndpoint])
	assert.NotEmpty(t, fn.Annotations[annotationCompartmentId])
	assert.NotEmpty(t, fn.CreatedAt)
	assert.NotEmpty(t, fn.UpdatedAt)
}

// Test get-by-name full model behaviour
func TestListFnsGetByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	appId := "ListFnsAppId"
	name := "ListFnsName"
	listFnsOK, err := shim.ListFns(&fns.ListFnsParams{
		AppID: &appId,
		Name:  &name,
	})
	assert.NoError(t, err)

	result := listFnsOK.GetPayload().Items
	assert.Len(t, result, 1)
	fn := result[0]
	assert.NotEmpty(t, fn.Config)
}

func TestUpdateFnConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	// Test config update (incl. merge behaviour)
	fn := modelsv2.Fn{
		Config: map[string]string{
			"GetFunctionKey2":    "",
			"UpdatedFunctionKey": "UpdatedFunctionValue",
		},
	}

	fnId := "UpdateFnId"
	updateFnOK, err := shim.UpdateFn(&fns.UpdateFnParams{
		FnID: fnId,
		Body: &fn,
	})
	assert.NoError(t, err)

	result := updateFnOK.GetPayload()
	expectedConfig := map[string]string{
		"GetFunctionKey1":    "GetFunctionValue1",
		"UpdatedFunctionKey": "UpdatedFunctionValue",
	}
	assert.Equal(t, expectedConfig, result.Config)
	// Check we haven't inadvertently updated other values
	assert.Equal(t, "GetFunctionImage", result.Image)
	assert.Equal(t, "GetFunctionDigest", result.Annotations[annotationImageDigest])
	assert.Equal(t, uint64(128), result.Memory)
	assert.Equal(t, int32(30), *result.Timeout)
}

func TestUpdateFnMemory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	fn := modelsv2.Fn{
		Memory: 256,
	}

	fnId := "UpdateFnId"
	updateFnOK, err := shim.UpdateFn(&fns.UpdateFnParams{
		FnID: fnId,
		Body: &fn,
	})
	assert.NoError(t, err)

	result := updateFnOK.GetPayload()
	assert.Equal(t, fn.Memory, result.Memory)
	// Check we haven't inadvertently updated other values
	expectedConfig := map[string]string{
		"UpdateFunctionKey1": "UpdateFunctionValue1",
		"UpdateFunctionKey2": "UpdateFunctionValue2",
	}
	assert.Equal(t, expectedConfig, result.Config)
	assert.Equal(t, "OriginalFunctionImage", result.Image)
	assert.Equal(t, "OriginalFunctionDigest", result.Annotations[annotationImageDigest])
	assert.Equal(t, int32(30), *result.Timeout)
}

func TestUpdateFnTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	timeout := int32(60)
	fn := modelsv2.Fn{
		Timeout: &timeout,
	}

	fnId := "UpdateFnId"
	updateFnOK, err := shim.UpdateFn(&fns.UpdateFnParams{
		FnID: fnId,
		Body: &fn,
	})
	assert.NoError(t, err)

	result := updateFnOK.GetPayload()
	assert.Equal(t, fn.Timeout, result.Timeout)
	// Check we haven't inadvertently updated other values
	expectedConfig := map[string]string{
		"UpdateFunctionKey1": "UpdateFunctionValue1",
		"UpdateFunctionKey2": "UpdateFunctionValue2",
	}
	assert.Equal(t, expectedConfig, result.Config)
	assert.Equal(t, "OriginalFunctionImage", result.Image)
	assert.Equal(t, "OriginalFunctionDigest", result.Annotations[annotationImageDigest])
	assert.Equal(t, uint64(128), result.Memory)
}

func TestUpdateFnImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	fn := modelsv2.Fn{
		Image: "UpdateFnImage",
		Annotations: map[string]interface{}{
			annotationImageDigest: "UpdateFnDigest",
		},
	}

	fnId := "UpdateFnId"
	updateFnOK, err := shim.UpdateFn(&fns.UpdateFnParams{
		FnID: fnId,
		Body: &fn,
	})
	assert.NoError(t, err)

	result := updateFnOK.GetPayload()
	assert.Equal(t, fn.Image, result.Image)
	assert.Equal(t, fn.Annotations[annotationImageDigest], result.Annotations[annotationImageDigest])
	// Check we haven't inadvertently updated other values
	expectedConfig := map[string]string{
		"UpdateFunctionKey1": "UpdateFunctionValue1",
		"UpdateFunctionKey2": "UpdateFunctionValue2",
	}
	assert.Equal(t, expectedConfig, result.Config)
	assert.Equal(t, uint64(128), result.Memory)
	assert.Equal(t, int32(30), *result.Timeout)
}

func TestUpdateFnBlankDigest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	shim := NewFnsShim(c)

	fn := modelsv2.Fn{
		Annotations: map[string]interface{}{
			annotationImageDigest: "",
		},
	}

	fnId := "UpdateFnId"
	updateFnOK, err := shim.UpdateFn(&fns.UpdateFnParams{
		FnID: fnId,
		Body: &fn,
	})
	assert.NoError(t, err)
	result := updateFnOK.GetPayload()
	assert.NotEmpty(t, result.Annotations[annotationImageDigest])
}
