package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	timePath := flag.String("time", "/usr/bin/time", "Path to the time binary.")
	rsPath := flag.String("routesum", "routesum", "Path to the routesum binary. Defaults to first found in $PATH.")
	inputPath := flag.String("input", "", "Path to routesum input. Required.")
	numRuns := flag.Int("num-runs", 1, "Number of times to run the input.")
	runLabel := flag.String("run-label", "", "Label for the run(s). Required.")

	flag.Parse()
	if len(*runLabel) == 0 {
		fatalf(errors.New("run-label cannot be empty"))
	}

	if err := runAndInterpret(*timePath, *rsPath, *inputPath, *numRuns, *runLabel); err != nil {
		fatalf(fmt.Errorf("run and interpret: %w", err))
	}
}

func fatalf(err error) {
	fmt.Fprintf(os.Stderr, "%+v\n", err)
	os.Exit(-1)
}

func runAndInterpret(timePath, rsPath, inputPath string, numRuns int, runLabel string) error {
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

	csvOut := csv.NewWriter(os.Stdout)
	if err := csvOut.Write([]string{"Input", "Label", "Metric", "Amount"}); err != nil {
		return errors.Wrap(err, "write csv header")
	}

	for i := 0; i < numRuns; i++ {
		if i != 0 {
			if _, err := inputFile.Seek(0, 0); err != nil {
				return errors.Wrap(err, "rewind input file")
			}
		}

		var b bytes.Buffer
		cmd := exec.Command(filepath.Clean(timePath), filepath.Clean(rsPath), "-show-mem-stats") //nolint: gosec
		cmd.Stdin = inputFile
		cmd.Stdout = nil
		cmd.Stderr = &b

		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, "run routesum")
		}

		measurements, err := interpret(b.String())
		if err != nil {
			return fmt.Errorf("interpret mem stat output: %w", err)
		}
		for _, m := range measurements {
			if err := csvOut.Write([]string{inputBase, runLabel, m.metric, m.amount}); err != nil {
				return errors.Wrap(err, "write csv data line")
			}
		}
	}

	csvOut.Flush()
	if err := csvOut.Error(); err != nil {
		return errors.Wrapf(err, "flush csv buffer")
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
