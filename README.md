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
Input,Metric,Binary,Amount
GeoLite2-ASN-Networks.txt,Num internal nodes,routesum-bst,4622053
GeoLite2-ASN-Networks.txt,Num leaf nodes,routesum-bst,4622055
GeoLite2-ASN-Networks.txt,Size of all internal nodes,routesum-bst,463039380
...
GeoLite2-ASN-Networks.txt,Num internal nodes,routesum-radix,622053
GeoLite2-ASN-Networks.txt,Num leaf nodes,routesum-radix,622055
GeoLite2-ASN-Networks.txt,Size of all internal nodes,routesum-radix,63039380
...
GeoLite2-Country-Networks.txt,Num internal nodes,routesum-bst,496323
GeoLite2-Country-Networks.txt,Num leaf nodes,routesum-bst,496325
GeoLite2-Country-Networks.txt,Size of all internal nodes,routesum-bst,412266326
...
GeoLite2-Country-Networks.txt,Num internal nodes,routesum-radix,96323
GeoLite2-Country-Networks.txt,Num leaf nodes,routesum-radix,96325
GeoLite2-Country-Networks.txt,Size of all internal nodes,routesum-radix,12266326
...
```

# Reporting Bugs and Issues

Bugs and other issues can be reported by filing an issue on our [GitHub issue
tracker](https://github.com/PatrickCronin/routesum-runner/issues).

# Copyright and License

This software is Copyright (c) 2021-2022 by Patrick Cronin.

This is free software, licensed under the terms of the [MIT
License](https://github.com/PatrickCronin/routesum-runner/LICENSE.md).
