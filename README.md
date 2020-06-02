
# TCN-PSI Protocol

TCN protocol based on Private Set Intersection Cardinality.

## Description

1. The server receives a list of TCN reports and expands their Temporary Contact Numbers.
2. The server loads the encoded elements in the PSI server.
3. The client encodes a list of TCNs and loads them in the PSI client logic.
4. The client and server follow the PSI protocol, and the client receives the intersection size.

See the full PSI description [here](https://github.com/OpenMined/PSI/blob/master/private_set_intersection/cpp/psi_client.h)

## Requirements

There are requirements for the entire project which each language shares. There also could be requirements for each target language:

### Global Requirements

These are the common requirements across all target languages of this project.

- A compiler such as clang, gcc, or msvc
- [Bazel](https://bazel.build)

## Installation

The repository uses a folder structure to isolate the supported targets from one another:

```
tcn_psi/<target language>/<sources>
```

### Go

See the [Go README.md](tcn_psi/go/README.md)


## Usage

To use this library in another Bazel project, add the following in your WORKSPACE file:

```
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
   name = "org_openmined_tcn_psi",
   remote = "https://github.com/OpenMined/TCN-PSI",
   branch = "master",
   init_submodules = True,
)

load("@org_openmined_psi//tcn_psi:preload.bzl", "psi_preload")

psi_preload()

load("@org_openmined_psi//tcn_psi:deps.bzl", "psi_deps")

psi_deps()

```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Contributors

See [CONTRIBUTORS.md](CONTRIBUTORS.md).

## License
[Apache License 2.0](https://choosealicense.com/licenses/apache-2.0/)
