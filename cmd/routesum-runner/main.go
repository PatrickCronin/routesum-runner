// Package main provides a program that creates CSVs of routesum performance data over multiple runs
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	a, err := parseArgs()
	if err != nil {
		fatalf(fmt.Errorf("parse args: %w", err))
	}

	if err := runAllInputsAndBinaries(a); err != nil {
		fatalf(fmt.Errorf("run and interpret: %w", err))
	}
}

func fatalf(err error) {
	fmt.Fprintf(os.Stderr, "%+v\n", err)
	os.Exit(-1)
}

func runAllInputsAndBinaries(a *args) error {
	csvOut := csv.NewWriter(os.Stdout)
	if err := csvOut.Write([]string{"Input", "Binary", "Metric", "Amount"}); err != nil {
		return errors.Wrap(err, "write csv header")
	}

	for _, inputPath := range a.inputPaths {
		for _, rsBinPath := range a.rsBinPaths {
			err := runNTimesAndInterpret(a.timeBinPath, rsBinPath, inputPath, a.numRuns, csvOut)
			if err != nil {
				return fmt.Errorf("processing %s with %s: %w", inputPath, rsBinPath, err)
			}
		}
	}

	csvOut.Flush()
	if err := csvOut.Error(); err != nil {
		return errors.Wrapf(err, "flush csv buffer")
	}

	return nil
}

func runNTimesAndInterpret(
	timeBinPath, rsBinPath, inputPath string,
	numRuns int,
	csvOut *csv.Writer,
) error {
	inputFile, err := os.Open(filepath.Clean(inputPath))
	if err != nil {
		return errors.Wrapf(err, "open %s for reading", inputPath)
	}
	defer func() { //nolint: gosec // we're just reading from the file
		if err := inputFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close input file: %+v\n", err)
		}
	}()

	inputBase := filepath.Base(inputPath)
	rsBinBase := filepath.Base(rsBinPath)

	for i := 0; i < numRuns; i++ {
		if i != 0 {
			if _, err := inputFile.Seek(0, 0); err != nil {
				return errors.Wrap(err, "rewind input file")
			}
		}

		var b bytes.Buffer
		cmd := exec.Command(filepath.Clean(timeBinPath), filepath.Clean(rsBinPath), "-show-mem-stats") //nolint: gosec
		cmd.Stdin = inputFile
		cmd.Stdout = nil
		cmd.Stderr = &b

		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, "run")
		}

		measurements, err := interpret(b.String())
		if err != nil {
			return fmt.Errorf("interpret mem stat output: %w", err)
		}
		for _, m := range measurements {
			if err := csvOut.Write([]string{inputBase, rsBinBase, m.metric, m.amount}); err != nil {
				return errors.Wrap(err, "write csv data line")
			}
		}
	}

	return nil
}

type measurement struct {
	metric string
	amount string
}

var (
	sectionLineRE = regexp.MustCompile(`^\S`)
	timeLineRE    = regexp.MustCompile(`^ *(?:\d+[.]\d+) real *(\d+[.]\d+) user *(\d+[.]\d+) sys\s*$`)
)

func interpret(memStats string) ([]measurement, error) {
	var measurements []measurement

	var section string
	s := bufio.NewScanner(strings.NewReader(memStats))
	for s.Scan() {
		line := s.Text()
		if sectionLineRE.MatchString(line) {
			// Starting a new section
			section = line
			continue
		}

		matches := timeLineRE.FindStringSubmatch(line)
		if len(matches) > 0 {
			// Parsing time output
			measurements = append(measurements, []measurement{
				{
					metric: "User-space Time",
					amount: matches[1],
				},
				{
					metric: "Kernel Time",
					amount: matches[2],
				},
			}...,
			)
			continue
		}

		// routesum memory metric
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			panic(line)
		}
		measurements = append(measurements, measurement{
			metric: fmt.Sprintf(
				"%s - %s",
				section,
				strings.TrimSpace(strings.SplitN(parts[0], "(", 2)[0]),
			),
			amount: strings.TrimSpace(parts[1]),
		})
	}

	if err := s.Err(); err != nil {
		return nil, errors.Wrap(err, "scan memstat output")
	}

	return measurements, nil
}
