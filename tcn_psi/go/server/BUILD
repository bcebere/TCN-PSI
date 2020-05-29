load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "server",
    srcs = ["server.go"],
    importpath = "github.com/bcebere/tcn-psi/server",
    visibility = ["//visibility:public"],
    deps = [
        "@org_openmined_psi//private_set_intersection/go/server",
    ],
)
