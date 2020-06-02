load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "client",
    srcs = ["client.go"],
    importpath = "github.com/openmined/tcn-psi/client",
    visibility = ["//visibility:public"],
    deps = [
        "@org_openmined_psi//private_set_intersection/go/client",
        "@org_openmined_tcn_psi//tcn_psi/go/tcn",
        ],
)

go_test(
    name = "client_test",
    srcs = ["client_test.go"],
    embed = [":client"],
    deps = [
        "@org_openmined_tcn_psi//tcn_psi/go/server",
    ],
)
