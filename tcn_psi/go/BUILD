load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/openmined/tcn-psi
gazelle(name = "gazelle")

go_binary(
    name = "tcn_psi_test",
    srcs = ["main.go"],
    deps = ["//tcn_psi/go/client:client", "//tcn_psi/go/server:server"],
    visibility = ["//visibility:public"],
)
