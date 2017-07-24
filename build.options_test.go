package main

import (
	"strings"
	"testing"
)

func TestBuildOptions1(t *testing.T) {

	cmd := BuildOptions(BuildOptionsConfig{
		Src:  "a/x.js",
		Dist: "b/x.js",
	})
	cmdStr := strings.Join(cmd, " ")
	excepted := "bundle a/x.js b/x.js"

	if cmdStr != excepted {
		t.Error(
			"For", cmd,
			"expected", excepted,
			"got", cmdStr,
		)
	}

}

func TestBuildOptions2(t *testing.T) {

	cmd := BuildOptions(BuildOptionsConfig{
		Src:     "a/x.js",
		Dist:    "b/x.js",
		Exclude: "module",
	})
	cmdStr := strings.Join(cmd, " ")
	excepted := "bundle a/x.js - module b/x.js"

	if cmdStr != excepted {
		t.Error(
			"For", cmd,
			"expected", excepted,
			"got", cmdStr,
		)
	}
}

func TestBuildOptions3(t *testing.T) {

	cmd := BuildOptions(BuildOptionsConfig{
		Src:            "a/x.js",
		Dist:           "b/x.js",
		Format:         "abc",
		AllowedFormats: "default|abc",
	})
	cmdStr := strings.Join(cmd, " ")
	excepted := "build a/x.js b/x.js --format abc"

	if cmdStr != excepted {
		t.Error(
			"For", cmd,
			"expected", excepted,
			"got", cmdStr,
		)
	}
}

func TestBuildOptions4(t *testing.T) {

	cmd := BuildOptions(BuildOptionsConfig{
		Src:            "a/x.js",
		Dist:           "b/x.js",
		Exclude:        "module",
		Format:         "abc",
		AllowedFormats: "default|abc",
	})
	cmdStr := strings.Join(cmd, " ")
	excepted := "build a/x.js - module b/x.js --format abc"

	if cmdStr != excepted {
		t.Error(
			"For", cmd,
			"expected", excepted,
			"got", cmdStr,
		)
	}
}

func TestBuildOptions5(t *testing.T) {

	cmd := BuildOptions(BuildOptionsConfig{
		Src:            "a/x.js",
		Dist:           "b/x.js",
		Exclude:        "module",
		Format:         "abc",
		NoMangle:       true,
		AllowedFormats: "default|abc",
	})
	cmdStr := strings.Join(cmd, " ")
	excepted := "build a/x.js - module b/x.js --no-mangle --format abc"

	if cmdStr != excepted {
		t.Error(
			"For", cmd,
			"expected", excepted,
			"got", cmdStr,
		)
	}
}

func TestBuildOptions6(t *testing.T) {

	cmd := BuildOptions(BuildOptionsConfig{
		Src:            "a/x.js",
		Dist:           "b/x.js",
		Exclude:        "module",
		Format:         "abc",
		NoMangle:       true,
		Minify:         true,
		AllowedFormats: "default|abc",
	})
	cmdStr := strings.Join(cmd, " ")
	excepted := "build a/x.js - module b/x.js --minify --no-mangle --format abc"

	if cmdStr != excepted {
		t.Error(
			"For", cmd,
			"expected", excepted,
			"got", cmdStr,
		)
	}
}
