package shim

import (
	"github.com/fnproject/fn_go/clientv2/apps"
	"github.com/fnproject/fn_go/modelsv2"
	"github.com/fnproject/fn_go/provider/oracle/shim/client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	compartmentId := "CreateAppCompartment"
	shim := NewAppsShim(c, compartmentId)

	syslogUrl := "CreateAppSyslogUrl"
	app := modelsv2.App{
		Name:      "CreateAppName",
		SyslogURL: &syslogUrl,
		Annotations: map[string]interface{}{
			annotationSubnet: []interface{}{"CreateAppSubnet"},
		},
		Config: map[string]string{
			"CreateAppKey": "CreateAppValue",
		},
	}

	createAppOK, err := shim.CreateApp(&apps.CreateAppParams{
		Body: &app,
	})
	assert.NoError(t, err)

	expectedAnnotations := app.Annotations
	expectedAnnotations[annotationCompartmentId] = compartmentId

	result := createAppOK.GetPayload()
	assert.Equal(t, app.Name, result.Name)
	assert.Equal(t, app.SyslogURL, result.SyslogURL)
	assert.Equal(t, expectedAnnotations, result.Annotations)
	assert.Equal(t, app.Config, result.Config)
	assert.NotEmpty(t, result.ID)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
}

func TestDeleteApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	compartmentId := "DeleteAppCompartment"
	shim := NewAppsShim(c, compartmentId)

	_, err := shim.DeleteApp(&apps.DeleteAppParams{
		AppID: "DeleteAppId",
	})
	assert.NoError(t, err)
}

func TestGetApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	compartmentId := "GetAppCompartment"
	shim := NewAppsShim(c, compartmentId)

	appId := "GetAppId"
	getAppOK, err := shim.GetApp(&apps.GetAppParams{
		AppID: appId,
	})
	assert.NoError(t, err)

	result := getAppOK.GetPayload()
	assert.Equal(t, appId, result.ID)
	assert.NotEmpty(t, result.Name)
	assert.NotEmpty(t, result.SyslogURL)
	assert.NotEmpty(t, result.Annotations[annotationSubnet])
	assert.NotEmpty(t, result.Annotations[annotationCompartmentId])
	assert.NotEmpty(t, result.Config)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)
}

func TestListAppsMultiplePages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	compartmentId := "ListAppsCompartment"
	shim := NewAppsShim(c, compartmentId)

	params := &apps.ListAppsParams{}
	var results []*modelsv2.App
	for {
		listAppsOK, err := shim.ListApps(params)
		assert.NoError(t, err)

		result := listAppsOK.GetPayload()
		results = append(results, result.Items...)

		if result.NextCursor == "" {
			break
		}

		params.Cursor = &result.NextCursor
	}
	assert.Len(t, results, 9)
	app := results[0]
	assert.Equal(t, compartmentId, app.Annotations[annotationCompartmentId])
	assert.NotEmpty(t, app.ID)
	assert.NotEmpty(t, app.Name)
	assert.NotEmpty(t, app.Annotations[annotationSubnet])
	assert.NotEmpty(t, app.CreatedAt)
	assert.NotEmpty(t, app.UpdatedAt)
}

// Test get-by-name full model behaviour
func TestListAppsGetByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	compartmentId := "ListAppsCompartment"
	shim := NewAppsShim(c, compartmentId)

	name := "ListAppsName"
	listAppsOK, err := shim.ListApps(&apps.ListAppsParams{
		Name: &name,
	})
	assert.NoError(t, err)

	result := listAppsOK.GetPayload().Items
	assert.Len(t, result, 1)
	app := result[0]
	assert.NotEmpty(t, app.Config)
	assert.NotEmpty(t, app.SyslogURL)
}

func TestUpdateAppConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	compartmentId := "UpdateAppCompartment"
	shim := NewAppsShim(c, compartmentId)

	// Test config update (incl. merge behaviour)
	app := modelsv2.App{
		Config: map[string]string{
			"GetApplicationKey2":    "",
			"UpdatedApplicationKey": "UpdatedApplicationValue",
		},
	}

	appId := "UpdateAppId"
	updateAppOK, err := shim.UpdateApp(&apps.UpdateAppParams{
		AppID: appId,
		Body:  &app,
	})
	assert.NoError(t, err)

	result := updateAppOK.GetPayload()
	expectedConfig := map[string]string{
		"GetApplicationKey1":    "GetApplicationValue1",
		"UpdatedApplicationKey": "UpdatedApplicationValue",
	}
	assert.Equal(t, expectedConfig, result.Config)
	// Check we haven't inadvertently changed syslogUrl
	assert.Equal(t, "OriginalApplicationSyslogUrl", *result.SyslogURL)
}

func TestUpdateAppSyslogUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewMockFunctionsManagementClientBasic(ctrl)
	compartmentId := "UpdateAppCompartment"
	shim := NewAppsShim(c, compartmentId)

	// Test syslogUrl update
	syslogUrl := "UpdatedApplicationSyslogUrl"
	app := modelsv2.App{
		SyslogURL: &syslogUrl,
	}
	appId := "UpdateAppId"
	updateAppOK, err := shim.UpdateApp(&apps.UpdateAppParams{
		AppID: appId,
		Body:  &app,
	})
	assert.NoError(t, err)

	result := updateAppOK.GetPayload()
	// Check we haven't inadvertently changed config
	expectedConfig := map[string]string{
		"UpdateApplicationKey1": "UpdateApplicationValue1",
		"UpdateApplicationKey2": "UpdateApplicationValue2",
	}
	assert.Equal(t, app.SyslogURL, result.SyslogURL)
	assert.Equal(t, expectedConfig, result.Config)
}
