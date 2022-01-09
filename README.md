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
  -input value
        Path to an input file. Can be specified multiple times. At least once is required.
  -num-runs int
        Number of times to run each input against each routesum binary (default 5)
  -routesum value
        Path to routesum binary. Can be specified multiple times. At least once is required.
  -time string
        Path to the time binary. (default "/usr/bin/time")

$ ./routesum-runner -input ~/routesum-performance/dshield-intelfeed-ips.txt -input ~/routesum-performance/GeoLite2-Country-Networks.txt -routesum ./routesum-bst -routesum ./routesum-radix -num-runs 10
Input,Binary,Metric,Amount
dshield-intelfeed-ips.txt,routesum-bst,To Store Routes - Δ total allocated bytes,44914264
dshield-intelfeed-ips.txt,routesum-bst,To Store Routes - Δ mallocs,1594243
dshield-intelfeed-ips.txt,routesum-bst,To Store Routes - Δ frees,468337
...
dshield-intelfeed-ips.txt,routesum-radix,To Store Routes - Δ total allocated bytes,17896936
dshield-intelfeed-ips.txt,routesum-radix,To Store Routes - Δ mallocs,583901
dshield-intelfeed-ips.txt,routesum-radix,To Store Routes - Δ frees,312098
...
GeoLite2-Country-Networks.txt,routesum-bst,To Store Routes - Δ total allocated bytes,2355023304
GeoLite2-Country-Networks.txt,routesum-bst,To Store Routes - Δ mallocs,24400781
GeoLite2-Country-Networks.txt,routesum-bst,To Store Routes - Δ frees,23528854
...
GeoLite2-Country-Networks.txt,routesum-radix,To Store Routes - Δ total allocated bytes,852038168
GeoLite2-Country-Networks.txt,routesum-radix,To Store Routes - Δ mallocs,17451592
GeoLite2-Country-Networks.txt,routesum-radix,To Store Routes - Δ frees,16970186
...
```

# Reporting Bugs and Issues

Bugs and other issues can be reported by filing an issue on our [GitHub issue
tracker](https://github.com/PatrickCronin/routesum-runner/issues).

# Copyright and License

This software is Copyright (c) 2021-2022 by Patrick Cronin.

This is free software, licensed under the terms of the [MIT
License](https://github.com/PatrickCronin/routesum-runner/LICENSE.md).
