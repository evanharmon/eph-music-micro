package helper

import (
	"os"
	"testing"

	testhelper "github.com/evanharmon/eph-music-micro/helper/testhelper"
)

//
var (
	envVarProjectID = "GOOGLE_PROJECT_ID"
	projectID       = ""
)

// setup sets the necessary vars from ENV exports
func setup(t *testing.T) {
	id, err := GetEnv(envVarProjectID)
	testhelper.Ok(t, err)
	testhelper.Assert(t, len(id) != 0, "Export necessary ENV vars before running tests")
	projectID = id
}

// TestGetEnv tests the safe lookup and retrieval of environment variable values
func TestGetEnv(t *testing.T) {
	tests := map[string]struct {
		export string
		result string
	}{
		"argument error":    {"", ""},
		"env missing error": {envVarProjectID, ""},
		"valid env empty":   {envVarProjectID, ""},
		"valid env":         {envVarProjectID, projectID},
	}

	for k, test := range tests {
		if k == "argument error" {
			val, err := GetEnv(test.export)
			testhelper.Assert(t, len(val) == 0, "should return empty string as default")
			testhelper.Throws(t, err)
		}
		if k == "env missing error" {
			setup(t)
			testhelper.Ok(t, os.Unsetenv(envVarProjectID))
			_, err := GetEnv(test.export)
			testhelper.Throws(t, err)
		}
		if k == "valid env empty" {
			os.Setenv(envVarProjectID, "")
			val, err := GetEnv(test.export)
			testhelper.Ok(t, err)
			testhelper.Assert(t, len(val) == 0, "Valid ENV set should get value")
		}
		if k == "valid env" {
			os.Setenv(envVarProjectID, projectID)
			val, err := GetEnv(test.export)
			testhelper.Ok(t, err)
			testhelper.Assert(t, len(val) != 0, "Valid ENV set should get value")
		}
	}
}
