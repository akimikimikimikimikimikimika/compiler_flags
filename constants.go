package main

const xcode_dev_dir string = "/Applications/Xcode.app"
const clt_dev_dir string = "/Library/Developer/CommandLineTools"
const xcrun string = "/usr/bin/xcrun"

const c_macos_xcode string = "/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/clang"
const c_macos_clt string = "/Library/Developer/CommandLineTools/usr/bin/clang"
const c_macos_brew string = "/usr/local/opt/llvm/bin/clang"
const cpp_macos_brew string = "/usr/local/opt/llvm/bin/clang++"
const c_macos_macports string = "/opt/local/bin/clang-mp-13"
const cpp_macos_macports string = "/opt/local/bin/clang++-mp-13"
const c_macos_graalvm string = "/Library/Java/JavaVirtualMachines/graalvm.jdk/Contents/Home/lib/llvm/bin/clang"
const cpp_macos_graalvm string = "/Library/Java/JavaVirtualMachines/graalvm.jdk/Contents/Home/lib/llvm/bin/clang++"
const c_macos_gcc string = "/usr/local/opt/gcc/bin/gcc-11"
const cpp_macos_gcc string = "/usr/local/opt/gcc/bin/g++-11"
const emscripten_macos string = "/usr/local/opt/emscripten/bin/emcc"
const go_macos string = "/usr/local/opt/go/bin/go"
const go_linux string = "/usr/bin/go"

const c_linux_clang string = "/usr/bin/clang-11"
const cpp_linux_clang string = "/usr/bin/clang++-11"

const cargo_rustup string = "~/.cargo/bin/cargo"
const rustc_rustup string = "~/.cargo/bin/rustc"
