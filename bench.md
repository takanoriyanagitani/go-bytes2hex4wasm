# Simple Benchmark

- CPU: Apple M3 Max
- RAM: 48 GB

| input data | binary type            | output speed  | ratio  | CPU% | user |
|:----------:|:----------------------:|:-------------:|:------:|:----:|:----:|
| random gen | wasi-wasm              |     4.5 MiB/s |    0%  |  99% |  98% |
| random gen | native-wasm            |   541.7 MiB/s |   37%  |  15% |   9% |
| big file   | pure go(encoding/hex)  | 1,456.4 MiB/s | (100%) |  95% |  57% |
| big file   | native-wasm(optimized) | 3,927.5 MiB/s |  270%  |  89% |  63% |
| big file   | native-wasm            | 4,965.4 MiB/s |  341%  |  85% |  51% |
