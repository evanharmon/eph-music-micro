package storage

import (
	"fmt"
	"os"
	"testing"
)

// BeforeTestGetEnv unsets the environment variable for a Google Cloud Project
// ID
func BeforeTestGetEnv() error {
	key := "GOOGLE_PROJECT_ID"
	if err := os.Unsetenv(key); err != nil {
		return fmt.Errorf("BeforeTestGetEnv: failed to unset ENV %s: %v", key, err)
	}

	return nil
}

// TestGetEnv tests the safe lookup of environment variables
func TestGetEnv(t *testing.T) {
	// Test Empty String
	if err := BeforeTestGetEnv(); err != nil {
		t.Errorf("Before failed: %v", err)
	}

	val, err := getEnv("")
	if len(val) != 0 {
		t.Error("Empty string as argument should return empty string as default")
	}
	if err == nil {
		t.Error("Empty string as argument string should throw error")
	}

	// Test ENV missing
	if err := BeforeTestGetEnv(); err != nil {
		t.Errorf("Before failed: %v", err)
	}

	key := "GOOGLE_PROJECT_ID"
	val2, err2 := getEnv(key)
	if len(val2) != 0 {
		t.Error("Missing ENV var should return empty string as default")
	}
	if err2 == nil {
		t.Error("Missing ENV var should throw error")
	}

	// VALID ENV should return
	if err := BeforeTestGetEnv(); err != nil {
		t.Errorf("Before failed: %v", err)
	}

	os.Setenv(key, "")
	val4, err4 := getEnv(key)
	if err4 != nil {
		t.Error("Valid ENV set should not throw error")
	}
	if len(val4) != 0 {
		t.Error("Valid ENV set should get value")
	}

	// AFTER
	os.Unsetenv(key)
}

// TestConfigureStoraage tests creating a client for the Google Cloud Storage
// API
func TestConfigureStorage(t *testing.T) {
	key := "eph-music"
	storageBucket, err := configureStorage(key)
	if storageBucket == nil {
		t.Errorf("should return valid storageBucket for bucketId %s\n", key)
	}

	if err != nil {
		t.Errorf("should not throw error on valid bucketId %s\n", key)
	}
}
