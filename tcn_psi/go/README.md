# TCN-PSI Cardinality - Go [![Go Report Card](https://goreportcard.com/badge/github.com/OpenMined/TCN-PSI)](https://goreportcard.com/report/github.com/OpenMined/TCN-PSI)

TCN protocol based on Private Set Intersection Cardinality.


## TCN client [![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/OpenMined/TCN-PSI/tcn_psi/go/client)
```
import "github.com/openmined/tcn-psi/client"
```

## TCN server [![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/OpenMined/TCN-PSI/tcn_psi/go/server)
```
import "github.com/openmined/tcn-psi/client"
```

## Tests
```
bazel test //tcn_psi/go/... --test_output=all
```

## Benchmarks
```
bazel test //tcn_psi/go/... --test_arg=-test.bench=. --test_output=all
```

## Integration

* Add Bazel depends to your WORKSPACE, as indicated in the [Usage](https://github.com/OpenMined/PSI#Usage) section.
* Add the server or the client to your deps in the BUILD file


```
go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    deps = [
            "@org_openmined_tcn_psi//tcn_psi/go/server",
            "@org_openmined_tcn_psi//tcn_psi/go/client",
            ],
)
```


* Import and use the library

```go
package main
import (
    "fmt"
    "github.com/openmined/tcn-psi/client"
    "github.com/openmined/tcn-psi/server"
)

func main(){
    tcnServer, err := server.CreateWithNewKey()
    if err == nil {
        fmt.Println("server loaded")
    }

    tcnClient, err := client.Create()
    if err == nil  {
        fmt.Println("client loaded")
    }
}
```

