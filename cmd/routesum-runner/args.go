package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type pathsFlag []string

func (paths *pathsFlag) String() string {
	return strings.Join(*paths, ", ")
}

func (paths *pathsFlag) Set(newPath string) error {
	*paths = append(*paths, newPath)
	return nil
}

type args struct {
	timeBinPath string
	numRuns     int
	rsBinPaths  []string
	inputPaths  []string
}

func parseArgs() (*args, error) {
	timeBinPath := flag.String("time", "/usr/bin/time", "Path to the time binary.")
	numRuns := flag.Int(
		"num-runs",
		5,
		"Number of times to run each input against each routesum binary",
	)

	var rsBinPaths pathsFlag
	flag.Var(
		&rsBinPaths,
		"routesum",
		"Path to routesum binary. Can be specified multiple times. At least once is required.",
	)

	var inputPaths pathsFlag
	flag.Var(
		&inputPaths,
		"input",
		"Path to an input file. Can be specified multiple times. At least once is required.",
	)

	flag.Parse()

	if err := assertPathExists("time", *timeBinPath); err != nil {
		return nil, err
	}

	if len(rsBinPaths) == 0 {
		argFatal("You must specify at least one routesum binary with the `routesum` arg.")
	}
	cleanedRSBinPaths, err := cleanAndAssertExistence("routesum binary", rsBinPaths)
	if err != nil {
		return nil, err
	}

	if len(inputPaths) == 0 {
		argFatal("You must specify at least one input file with the `input` arg.")
	}
	cleanedInputPaths, err := cleanAndAssertExistence("input file", inputPaths)
	if err != nil {
		return nil, err
	}

	return &args{
		timeBinPath: *timeBinPath,
		numRuns:     *numRuns,
		rsBinPaths:  cleanedRSBinPaths,
		inputPaths:  cleanedInputPaths,
	}, nil
}

func argFatal(msg string) {
	fmt.Fprintf(os.Stderr, msg+"\n")
	os.Exit(1)
}

func cleanAndAssertExistence(pathDescriptor string, paths pathsFlag) ([]string, error) {
	cleaned := make([]string, 0, len(paths))
	for _, path := range paths {
		if err := assertPathExists(pathDescriptor, path); err != nil {
			return nil, err
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, errors.Wrap(err, "confirm/compute absolute path")
		}
		cleaned = append(cleaned, absPath)
	}

	return cleaned, nil
}

func assertPathExists(pathDescriptor, path string) error {
	_, err := os.Stat(filepath.Clean(path))
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		argFatal(fmt.Sprintf("Failed to find a(n) %s at %s.", pathDescriptor, path))
	}
	return err //nolint: wrapcheck
}
