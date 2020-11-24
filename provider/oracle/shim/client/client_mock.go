package client

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/oracle/oci-go-sdk/v28/common"
	"github.com/oracle/oci-go-sdk/v28/functions"
	"time"
)

func NewMockFunctionsManagementClientBasic(ctrl *gomock.Controller) FunctionsManagementClient {
	m := NewMockFunctionsManagementClient(ctrl)

	// Create
	m.EXPECT().
		CreateApplication(
			gomock.Any(),
			gomock.AssignableToTypeOf(functions.CreateApplicationRequest{}),
		).
		DoAndReturn(
			func(ctx context.Context, request functions.CreateApplicationRequest) (functions.CreateApplicationResponse, error) {
				id := "CreateApplicationId"
				return functions.CreateApplicationResponse{
					Application: functions.Application{
						Id:             &id,
						CompartmentId:  request.CompartmentId,
						DisplayName:    request.DisplayName,
						LifecycleState: functions.ApplicationLifecycleStateActive,
						Config:         request.Config,
						SubnetIds:      request.SubnetIds,
						SyslogUrl:      request.SyslogUrl,
						FreeformTags:   request.FreeformTags,
						DefinedTags:    request.DefinedTags,
						TimeCreated:    &common.SDKTime{Time: time.Now()},
						TimeUpdated:    &common.SDKTime{Time: time.Now()},
					},
				}, nil
			},
		).
		AnyTimes()

	// Delete
	m.EXPECT().
		DeleteApplication(
			gomock.Any(),
			gomock.AssignableToTypeOf(functions.DeleteApplicationRequest{}),
		).
		Return(functions.DeleteApplicationResponse{}, nil).
		AnyTimes()

	// Get
	m.EXPECT().
		GetApplication(
			gomock.Any(),
			gomock.AssignableToTypeOf(functions.GetApplicationRequest{}),
		).
		DoAndReturn(
			func(ctx context.Context, request functions.GetApplicationRequest) (functions.GetApplicationResponse, error) {
				compartment := "GetApplicationCompartment"
				displayName := "GetApplicationDisplayName"
				syslogUrl := "GetApplicationSyslogUrl"
				return functions.GetApplicationResponse{
					Application: functions.Application{
						Id:             request.ApplicationId,
						CompartmentId:  &compartment,
						DisplayName:    &displayName,
						LifecycleState: functions.ApplicationLifecycleStateActive,
						Config: map[string]string{
							"GetApplicationKey1": "GetApplicationValue1",
							"GetApplicationKey2": "GetApplicationValue2",
						},
						SubnetIds:    []string{"GetApplicationSubnet"},
						SyslogUrl:    &syslogUrl,
						FreeformTags: nil,
						DefinedTags:  nil,
						TimeCreated:  &common.SDKTime{Time: time.Now()},
						TimeUpdated:  &common.SDKTime{Time: time.Now()},
					},
				}, nil
			},
		).
		AnyTimes()

	// List
	m.EXPECT().
		ListApplications(
			gomock.Any(),
			gomock.AssignableToTypeOf(functions.ListApplicationsRequest{}),
		).
		DoAndReturn(
			func(ctx context.Context, request functions.ListApplicationsRequest) (functions.ListApplicationsResponse, error) {
				if request.DisplayName != nil && *request.DisplayName != "" {
					return functions.ListApplicationsResponse{
						Items: []functions.ApplicationSummary{
							newBasicApplicationSummary(0),
						},
					}, nil
				}

				page0 := []functions.ApplicationSummary{
					newBasicApplicationSummary(0),
					newBasicApplicationSummary(1),
					newBasicApplicationSummary(2),
				}
				page1 := []functions.ApplicationSummary{
					newBasicApplicationSummary(3),
					newBasicApplicationSummary(4),
					newBasicApplicationSummary(5),
				}
				page2 := []functions.ApplicationSummary{
					newBasicApplicationSummary(6),
					newBasicApplicationSummary(7),
					newBasicApplicationSummary(8),
				}

				var response functions.ListApplicationsResponse
				if request.Page == nil {
					opcNextPage := "1"
					response = functions.ListApplicationsResponse{Items: page0, OpcNextPage: &opcNextPage}
				} else if *request.Page == "1" {
					opcNextPage := "2"
					response = functions.ListApplicationsResponse{Items: page1, OpcNextPage: &opcNextPage}
				} else if *request.Page == "2" {
					response = functions.ListApplicationsResponse{Items: page2}
				}
				return response, nil
			},
		).
		AnyTimes()

	// Update
	m.EXPECT().
		UpdateApplication(
			gomock.Any(),
			gomock.AssignableToTypeOf(functions.UpdateApplicationRequest{}),
		).
		DoAndReturn(
			func(ctx context.Context, request functions.UpdateApplicationRequest) (functions.UpdateApplicationResponse, error) {
				id := "UpdateApplicationId"
				compartment := "UpdateApplicationCompartment"
				displayName := "UpdateApplicationDisplayName"
				config := map[string]string{
					"UpdateApplicationKey1": "UpdateApplicationValue1",
					"UpdateApplicationKey2": "UpdateApplicationValue2",
				}
				if request.Config != nil {
					config = request.Config
				}
				syslogUrl := "OriginalApplicationSyslogUrl"
				if request.SyslogUrl != nil {
					syslogUrl = *request.SyslogUrl
				}
				return functions.UpdateApplicationResponse{
					Application: functions.Application{
						Id:             &id,
						CompartmentId:  &compartment,
						DisplayName:    &displayName,
						LifecycleState: functions.ApplicationLifecycleStateActive,
						Config:         config,
						SubnetIds:      []string{"UpdateApplicationSubnet"},
						SyslogUrl:      &syslogUrl,
						FreeformTags:   nil,
						DefinedTags:    nil,
						TimeCreated:    &common.SDKTime{Time: time.Now()},
						TimeUpdated:    &common.SDKTime{Time: time.Now()},
					},
				}, nil
			},
		).
		AnyTimes()

	return m
}

func newBasicApplicationSummary(n int) functions.ApplicationSummary {
	id := fmt.Sprintf("ApplicationSummaryId%d", n)
	compartment := "ApplicationSummaryCompartment"
	displayName := fmt.Sprintf("ApplicationSummaryDisplayName%d", n)
	return functions.ApplicationSummary{
		Id:             &id,
		CompartmentId:  &compartment,
		DisplayName:    &displayName,
		LifecycleState: functions.ApplicationLifecycleStateActive,
		SubnetIds:      []string{"ApplicationSummarySubnet"},
		FreeformTags:   nil,
		DefinedTags:    nil,
		TimeCreated:    &common.SDKTime{Time: time.Now()},
		TimeUpdated:    &common.SDKTime{Time: time.Now()},
	}
}
