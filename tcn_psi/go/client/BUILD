load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "client",
    srcs = ["client.go"],
    importpath = "github.com/openmined/tcn-psi/client",
    visibility = ["//visibility:public"],
    deps = [
        "@org_openmined_psi//private_set_intersection/go/client",
    ],
)
