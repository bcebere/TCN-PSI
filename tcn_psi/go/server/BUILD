load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "server",
    srcs = ["server.go"],
    importpath = "github.com/openmined/tcn-psi/server",
    visibility = ["//visibility:public"],
    deps = [
        "@org_openmined_psi//private_set_intersection/go/server",
        "@org_openmined_tcn_psi//tcn_psi/go/tcn",
        ]
)

go_test(
    name = "server_test",
    srcs = ["server_test.go"],
    embed = [":server"],
    deps = ["@org_openmined_tcn_psi//tcn_psi/go/client"],
)
