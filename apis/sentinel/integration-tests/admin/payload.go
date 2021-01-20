package admin

import (
	"fmt"
	"os"
)

var loginPayload = fmt.Sprintf(
	`{
		"client_id": "%s",
		"client_secret": "%s"
	}`,
	os.Getenv("BITHIPPIE_CLIENT_ID"),
	os.Getenv("BITHIPPIE_CLIENT_SECRET"),
)

func createResourcePayload(name, sourceID string) string {
	return fmt.Sprintf(
		`{
			"data": {
				"attributes": {
					"context_id": "",
					"name": "%s",
					"source_id": "%s"
				},
				"type": "resource"
			}
		}`,
		name,
		sourceID,
	)
}

func createContextPayload(name string) string {
	return fmt.Sprintf(`{
  	"data": {
    	"attributes": {
      	"name": "%s"
    	},
    	"type": "context"
  	}
	}`,
		name)
}

const createPermissionPayload = `
{
  "data": {
    "attributes": {
      "name": "sentinel:read",
      "permitted": "allow"
    },
    "type": "permission"
  }
}`
