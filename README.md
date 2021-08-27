# Frame IO

[![Go Reference](https://pkg.go.dev/badge/github.com/lemon-mint/frameio.svg)](https://pkg.go.dev/github.com/lemon-mint/frameio)

DoS Safe Frame based IO.

Convert Stream to Frame!

- [X] Payload Size Verification (Safe from attacks that falsify payload size)
- [X] Read CallBack (Reduce Memory Allocation)
- [X] Read To Buffer (Reduce Memory Allocation)

## Protocol

![frame](https://raw.githubusercontent.com/lemon-mint/frameio/master/img/01.png)
"N" is set to 1 if it is the last block of that frame If there are additional blocks after that, N is set to 0
