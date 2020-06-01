#
# Copyright 2020 the authors listed in CONTRIBUTORS.md
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")


def tcn_psi_deps():
    # Language-specific dependencies.

    # Protobuf.
    rules_proto_dependencies()

    rules_proto_toolchains()

    # Golang.
    go_rules_dependencies()

    go_register_toolchains()

    rules_pkg_dependencies()

    gazelle_dependencies()

    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        tag = "v1.3.0",
    )

    go_repository(
        name = "com_github_juliangruber_go_intersect",
        importpath = "github.com/juliangruber/go-intersect",
        tag = "master",
    )
