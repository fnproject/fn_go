# FnProject SDK

This is a golang SDK for accessing the [Fn API](https://github.com/fnproject/fn/) - it allows you to create and modify applications and functions programmatically. 

Most of the code in the repository is automatically generated from the latest [Fn project API Swagger](https://github.com/fnproject/fn/blob/master/docs/swagger.yml)

Example:

```go
package main

import (
	"github.com/fnproject/fn_go/provider"
	"github.com/fnproject/fn_go"
	"github.com/fnproject/fn_go/client/apps"
	"context"
	"fmt"
)

func main() {

	config := provider.NewConfigSourceFromMap(map[string]string{
		"api-url": "http://localhost:8080",
	})

	currentProvider, err := fn_go.DefaultProviders.ProviderFromConfig("default", config, &provider.NopPassPhraseSource{})

	if err != nil {
		panic(err.Error())
	}

	appClient := currentProvider.APIClient().Apps

	ctx := context.Background()

	var cursor string
	for {
		params := &apps.GetAppsParams{
			Context: ctx,
		}
		if cursor != "" {
			params.Cursor = &cursor
		}

		gotApps, err := appClient.GetApps(params)
		if err != nil {
			panic(err.Error())
		}

		for _, app := range gotApps.Payload.Apps {
			fmt.Printf("App %s\n", app.Name)
		}

		if gotApps.Payload.NextCursor == "" {
			break
		}
		cursor = gotApps.Payload.NextCursor
	}
}

```