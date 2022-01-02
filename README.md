# routesum-runner

`routesum-runner` - create CSVs of routesum performance data over multiple runs

# Project Description

This is a quick-and-dirty program that will run a routesum binary against an
input file multiple times with `-show-mem-stats` enabled and timing each run,
collecting and outputting that data in CSV format. The data can then be pasted
into a spreadsheet for analysis.

# Usage

```
$ ./routesum-runner -h
Usage of ./routesum-runner:
  -input string
        Path to routesum input. Required.
  -num-runs int
        Number of times to run the input. (default 1)
  -routesum string
        Path to the routesum binary. Defaults to first found in $PATH. (default "routesum")
  -run-label string
        Label for the run(s). Required.
  -time string
        Path to the time binary. (default "/usr/bin/time")

$ ./routesum-runner -input networks.txt -num-runs 10 -run-label 'binary search tree'
Input,Label,Metric,Amount
networks.txt,binary search tree,To Store Routes - Δ total bytes allocated,17897192
networks.txt,binary search tree,To Store Routes - Δ mallocs,583904
networks.txt,binary search tree,To Store Routes - Δ frees,310725
...
```

# Reporting Bugs and Issues

Bugs and other issues can be reported by filing an issue on our [GitHub issue
tracker](https://github.com/PatrickCronin/routesum-runner/issues).

# Copyright and License

This software is Copyright (c) 2021-2022 by Patrick Cronin.

This is free software, licensed under the terms of the [MIT
License](https://github.com/PatrickCronin/routesum-runner/LICENSE.md).
