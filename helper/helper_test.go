package helper

import (
	"os"
	"testing"

	test "github.com/evanharmon/eph-music-micro/helper/testhelper"
)

//
var (
	envVarProjectID = "GOOGLE_PROJECT_ID"
	projectID       = ""
)

// setup sets the necessary vars from ENV exports
func setup(t *testing.T) {
	id, err := GetEnv(envVarProjectID)
	test.Ok(t, err)
	test.Assert(t, len(id) != 0, "Export necessary ENV vars before running tests")
	projectID = id
}

// TestGetEnv tests the safe lookup of environment variables
func TestGetEnv(t *testing.T) {
	// Test Empty String
	val, err := GetEnv("")
	test.Assert(t, len(val) == 0, "Empty string as argument should return empty string as default")
	test.Throws(t, err)

	// Test ENV missing
	setup(t)
	test.Ok(t, os.Unsetenv(envVarProjectID))
	_, err2 := GetEnv(envVarProjectID)
	test.Throws(t, err2)

	// VALID ENV empty string should return
	os.Setenv(envVarProjectID, "")
	val4, err4 := GetEnv(envVarProjectID)
	test.Ok(t, err4)
	test.Assert(t, len(val4) == 0, "Valid ENV set should get value")

	// VALID ENV value should return
	os.Setenv(envVarProjectID, projectID)
	val5, err5 := GetEnv(envVarProjectID)
	test.Assert(t, len(val5) != 0, "Valid ENV set should get value")
	test.Ok(t, err5)
}
