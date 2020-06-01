load("@io_bazel_rules_go//go:def.bzl", "go_test", "go_library")

go_library(
    name = "tcn",
    srcs = [
        "serialize.go",
        "report.go",
        "signedreport.go",
        "keys.go",
    ],
    deps = [
        "@com_github_stretchr_testify//assert:go_default_library",
    ],
    importpath = "github.com/openmined/tcn-psi/tcn",
    visibility = ["//visibility:public"],
)

go_test(
    name = "tcn_test",
    srcs = ["parse_test.go"],
    embed = [":tcn"],
)
