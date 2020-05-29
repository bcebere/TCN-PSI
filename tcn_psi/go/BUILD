load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/openmined/tcn-psi
gazelle(name = "gazelle")

go_binary(
    name = "tcn_psi_test",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/openmined/tcn-psi",
    visibility = ["//visibility:private"],
    deps = [
        "//tcn_psi/go/client",
        "//tcn_psi/go/server",
    ],
)
