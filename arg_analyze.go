package main

import "os"

func (o *options) arg_analyze() {

	args := os.Args[1:]
	if len(args) > 0 {

		for _, a := range args {
			switch a {

				case "help", "-help", "--help":
					help()

				case "c", "C":
					o.lang = c
				case "objc", "objective-c", "Objective-C":
					o.lang = c
					o.objc = true
				case "cpp", "c++", "C++":
					o.lang = cpp
				case "objc++", "objcpp", "objective-c++", "Objective-C++":
					o.lang = cpp
					o.objc = true
				case "swift", "Swift":
					o.lang = swift
				case "go", "Go", "golang":
					o.lang = golang
				case "gc":
					o.lang = golang
					o.cc_go = gc
				case "gccgo":
					o.lang = golang
					o.cc_go = gccgo
				case "rustc", "rust":
					o.lang = rustc
				case "cargo":
					o.lang = cargo

				case "xcode":
					o.cc = clang
					o.clang = clang_xcode
				case "clt":
					o.cc = clang
					o.clang = clang_clt
				case "brew","llvm":
					o.cc = clang
					o.clang = clang_brew
				case "macports":
					o.cc = clang
					o.clang = clang_macports
				case "graalvm":
					o.cc = clang
					o.clang = clang_graalvm
				case "clang":
					o.cc = clang
					o.lang = c
				case "clang++":
					o.cc = clang
					o.lang = cpp
				case "gcc":
					o.cc = gcc
					o.lang = c
				case "g++":
					o.cc = gcc
					o.lang = cpp
				case "msvc","cl":
					o.plat = windows
					o.cc = msvc
				case "vc":
					o.plat = windows
					o.cc = msvc
					o.lang = c
				case "vc++":
					o.plat = windows
					o.cc = msvc
					o.lang = cpp
				case "mingw":
					o.plat = windows
					o.cc = gcc

				case "c90":
					o.lang = c
					o.std_c = c90
				case "gnu90":
					o.lang = c
					o.std_c = gnu90
				case "c99":
					o.lang = c
					o.std_c = c99
				case "gnu99":
					o.lang = c
					o.std_c = gnu99
				case "c11":
					o.lang = c
					o.std_c = c11
				case "gnu11":
					o.lang = c
					o.std_c = gnu11
				case "c17":
					o.lang = c
					o.std_c = c17
				case "gnu17":
					o.lang = c
					o.std_c = gnu17
				case "c2x", "c20":
					o.lang = c
					o.std_c = c2x
				case "gnu2x":
					o.lang = c
					o.std_c = gnu2x

				case "c++98":
					o.lang = cpp
					o.std_cpp = cpp98
				case "gnu++98":
					o.lang = cpp
					o.std_cpp = gnupp98
				case "c++11", "c++0x":
					o.lang = cpp
					o.std_cpp = cpp11
				case "gnu++11", "gnu++0x":
					o.lang = cpp
					o.std_cpp = gnupp11
				case "c++14", "c++1y":
					o.lang = cpp
					o.std_cpp = cpp14
				case "gnu++14", "gnu++1y":
					o.lang = cpp
					o.std_cpp = gnupp14
				case "c++17", "c++1z":
					o.lang = cpp
					o.std_cpp = cpp17
				case "gnu++17", "gnu++1z":
					o.lang = cpp
					o.std_cpp = gnupp17
				case "c++20", "c++2a":
					o.lang = cpp
					o.std_cpp = cpp20
				case "gnu++20", "gnu++2a":
					o.lang = cpp
					o.std_cpp = gnupp20
				case "c++23", "c++2b":
					o.lang = cpp
					o.std_cpp = cpp23
				case "gnu++23", "gnu++2b":
					o.lang = cpp
					o.std_cpp = gnupp23

				case "shared", "so":
					o.output = shared
				case "dylib", "dynamiclib", "dll":
					o.output = dylib
				case "bundle":
					o.output = bundle
				case "object", "obj", "compile", "o":
					o.output = object
				case "assembly", "assemble", "asm", "s", "S":
					o.output = assembly_att
				case "assembly_intel", "asm_intel":
					o.output = assembly_intel
				case "preprocess":
					o.output = preprocess
				case "precompile":
					o.output = precompile
				case "llvm-byte", "llvm-bytecode", "llvm-bit", "llvm-bitcode", "bc":
					o.output = llvm_bytecode
				case "ir", "llvm-ir", "emit-llvm", "ll":
					o.output = llvm_ir
				case "syntax", "syntax-check", "syntax-check-only", "syntax-only":
					o.output = syntax_check
				case "show-macros", "macros":
					o.output = show_macros

				case "universal":
					o.arch = universal
				case "x86_64", "amd64":
					o.arch = x86_64
				case "x86_64h":
					o.arch = x86_64h
				case "x32", "x32abi":
					o.arch = x32
				case "i686", "i386", "x86":
					o.arch = i686
				case "i586":
					o.arch = i586
				case "arm64", "aarch64":
					o.arch = arm64
				case "arm64_be", "aarch64_be":
					o.arch = arm64_be
				case "arm64e":
					o.arch = arm64e
				case "arm64_32":
					o.arch = arm64_32
				case "armv7", "arm":
					o.arch = armv7
				case "armv7hf", "armhf":
					o.arch = armv7hf
				case "armv7el", "armv7le", "armel", "armle":
					o.arch = armv7el
				case "armv7s":
					o.arch = armv7s
				case "armv7k":
					o.arch = armv7k
				case "armv7m":
					o.arch = armv7m
				case "armv7em":
					o.arch = armv7em
				case "armv6":
					o.arch = armv6
				case "armv6hf":
					o.arch = armv6hf
				case "armv6el":
					o.arch = armv6el
				case "armv6m":
					o.arch = armv6m
				case "armv5","armv5el":
					o.arch = armv5
				case "mips":
					o.arch = mips
				case "mipsel", "mipsle":
					o.arch = mipsel
				case "mips64":
					o.arch = mips64
				case "mips64el", "mips64le":
					o.arch = mips64el
				case "ppc", "powerpc":
					o.arch = ppc
				case "ppc64", "powerpc64":
					o.arch = ppc64
				case "ppc64el", "powerpc64el":
					o.arch = ppc64el
				case "s390x":
					o.arch = s390x
				case "riscv64":
					o.arch = riscv64
				case "sparc64":
					o.arch = sparc64

				case "musl":
					o.musl = true

				case "native":
					o.native = true

				case "darwin":
					o.plat = darwin
				case "macos", "mac", "osx":
					o.plat = macos
				case "maccatalyst", "catalyst":
					o.plat = maccatalyst
				case "ios", "ipados", "iphoneos":
					o.plat = ios
				case "ios-simulator", "ios_simulator", "ios simulator":
					o.plat = ios_simulator
				case "watchos", "applewatch", "watch":
					o.plat = watchos
				case "watchos-simulator", "watchos_simulator", "watchos simulator":
					o.plat = watchos_simulator
				case "tvos", "appletv", "tv":
					o.plat = tvos
				case "tvos-simulator", "tvos_simulator", "tvos simulator":
					o.plat = tvos_simulator
				case "driverkit":
					o.plat = driverkit
				case "linux":
					o.plat = linux
				case "windows":
					o.plat = windows
				case "wasm":
					o.plat = wasm

				case "libc++", "libcpp":
					o.stdlib = libcpp
				case "libstdc++", "libstdcpp":
					o.stdlib = libstdcpp

				case "ld":
					o.linker = ld
				case "bfd":
					o.linker = bfd
				case "gold":
					o.linker = gold
				case "lld":
					o.linker = lld

				case "strict":
					o.flags_level = strict
				case "loose":
					o.flags_level = loose

				case "static", "static-link":
					o.link_type = statically
				case "dynamic", "dynamic-link":
					o.link_type = dynamically

				case "dry-run", "dryrun":
					o.dry_run = true
				case "stack-usage":
					o.stack_usage = true
				case "debug", "g":
					o.optimize = debug
				case "unoptimized", "O0":
					o.optimize = unoptimized
				case "release", "optimized", "optimize", "O2":
					o.optimize = normal
				case "debug-optimized", "Og":
					o.optimize = debug_optimized
				case "lto":
					o.optimize = lto
				case "faster", "fast", "Ofast", "O3":
					o.optimize = faster
				case "smaller", "small", "Os", "Oz":
					o.optimize = smaller
				case "protect","protection","sanitize":
					o.protection = true
				case "math":
					o.math = true
				case "openmp", "omp":
					o.openmp = true
				case "boost-random":
					o.boost_random = true
				case "thread", "pthread":
					o.thread = true
				case "posix":
					o.posix = true
				case "filesystem", "c++fs", "fs":
					o.cppfs = true
				case "opencl":
					o.opencl = true
				case "userlib":
					o.userlib = true
				case "no-unused-result":
					o.no_unused_result = true

				default:
					err("不明なオプション: " + a)

			}
		}
	} else {
		err("引数がありません")
	}

}
