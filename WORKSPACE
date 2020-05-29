workspace(name = "org_openmined_tcn_psi")
        
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
git_repository(
    name = "org_openmined_psi",
    remote = "https://github.com/OpenMined/PSI",
    branch = "master",
    init_submodules = True,
)
load("@org_openmined_psi//private_set_intersection:preload.bzl", "psi_preload")
psi_preload()

load("@org_openmined_psi//private_set_intersection:deps.bzl", "psi_deps")
psi_deps()


load("@org_openmined_tcn_psi//tcn_psi:preload.bzl", "tcn_psi_preload")

tcn_psi_preload()

load("@org_openmined_tcn_psi//tcn_psi:deps.bzl", "tcn_psi_deps")

tcn_psi_deps()
