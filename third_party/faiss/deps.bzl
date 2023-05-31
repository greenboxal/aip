load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

_ALL_CONTENT = """\
load("@rules_foreign_cc//foreign_cc:defs.bzl", "cmake")

filegroup(
    name = "all_srcs",
    srcs = glob(["**"]),
    visibility = ["//visibility:public"],
)

cmake(
    name = "faiss",
    lib_source = ":all_srcs",
    cache_entries = {
        "FAISS_ENABLE_C_API": "ON",
        "FAISS_ENABLE_PYTHON": "OFF",
        "FAISS_ENABLE_GPU": "OFF",
        "BUILD_SHARED_LIBS": "ON",
        "BUILD_TESTING": "OFF",
        "CMAKE_BUILD_TYPE": "Release",
    },
    out_shared_libs = select({
        "@io_bazel_rules_go//go/platform:darwin": [
            "libfaiss_c.dylib",
            "libfaiss.dylib",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "libfaiss_c.dll",
            "libfaiss.dll",
        ],
        "//conditions:default": [
            "libfaiss_c.so",
            "libfaiss.so",
        ],
    }),
    targets = [
        "faiss_c",
        "faiss",
    ],
    visibility = ["//visibility:public"],
)
"""

def faiss_dependencies():
    http_archive(
        name = "com_github_facebookresearch_faiss",
        build_file_content = _ALL_CONTENT,
        strip_prefix = "faiss-1.7.4",
        sha256 = "d9a7b31bf7fd6eb32c10b7ea7ff918160eed5be04fe63bb7b4b4b5f2bbde01ad",
        patch_args = ["-p1"],
        patches = ["//third_party/faiss:faiss.patch"],
        urls = [
            "https://github.com/facebookresearch/faiss/archive/refs/tags/v1.7.4.tar.gz",
        ],
    )
