package main

import (
	"os"
	"runtime"
)

func (o *options) env_analyze() {

	var native bool = true

	if o.plat==default_plat {

		val,exist := os.LookupEnv("FLAGS_PLATFORM")

		if exist {
			native = false
			switch val {
				case "darwin":
					o.plat = darwin
				case "linux":
					o.plat = linux
				case "windows":
					o.plat = windows
				case "mingw":
					o.plat = windows
					o.cc = gcc
				case "wasm":
					o.plat = wasm
				default:
					exist = false
			}
		}
		if !exist {
			switch runtime.GOOS {
				case "darwin":
					o.plat = darwin
				case "linux":
					o.plat = linux
				case "windows":
					o.plat = windows
			}
		}

	} else { native = false }

	if o.cc==default_compiler {

		val,exist := os.LookupEnv("FLAGS_CC")

		if exist {
			switch val {
				case "clang":
					o.cc = clang
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
				case "gcc":
					o.cc = gcc
				case "msvc":
					o.cc = msvc
				default:
					exist = false
			}
		}
		if !exist {
			switch o.simplified_plat() {
				case darwin:
					o.cc = clang
				case linux:
					o.cc = gcc
				case windows:
					if runtime.GOOS=="windows" {
						o.cc = clang
					} else {
						o.cc = gcc
					}
			}
		}

	}

	{
		val,exist := os.LookupEnv("FLAGS_GO")
		if exist {
			switch val {
				case "gc": o.cc_go = gc
				case "gccgo": o.cc_go = gccgo
			}
		}
	}

	if o.simplified_plat()==darwin && o.clang==clang_default {
		if is_exist(xcode_dev_dir) {
			o.clang = clang_xcode
		} else
		if is_exist(clt_dev_dir) {
			o.clang = clang_clt
		} else
		if is_exist(c_macos_brew) {
			o.clang = clang_brew
		} else
		if is_exist(c_macos_macports) {
			o.clang = clang_macports
		} else
		if is_exist(c_macos_graalvm) {
			o.clang = clang_graalvm
		}
	}

	if o.stdlib==default_stdlib {
		if o.cc==gcc {
			o.stdlib = libstdcpp
		} else {
			if o.plat==linux && o.arch!=default_arch {
				o.stdlib = libstdcpp
			} else {
				o.stdlib = libcpp
			}
		}
	}

	if o.arch==default_arch {

		val,exist := os.LookupEnv("FLAGS_ARCH")

		if exist {
			native = false
			switch val {
				case "universal":
					o.arch = universal
				case "x86_64","amd64":
					o.arch = x86_64
				case "x86_64h":
					o.arch = x86_64h
				case "x32","x32abi":
					o.arch = x32
				case "i686","i386","x86":
					o.arch = i686
				case "i586":
					o.arch = i586
				case "arm64","aarch64":
					o.arch = arm64
				case "arm64e":
					o.arch = arm64e
				case "arm64_32":
					o.arch = arm64_32
				case "armv7","arm":
					o.arch = armv7
				case "armv7s":
					o.arch = armv7s
				case "armv7k":
					o.arch = armv7k
				case "armv7m":
					o.arch = armv7m
				case "armv7em":
					o.arch = armv7em
				case "armv7hf":
					o.arch = armv7hf
				case "armv7el","armv7le":
					o.arch = armv7el
				case "armv6":
					o.arch = armv6
				case "armv6m":
					o.arch = armv6m
				case "mips":
					o.arch = mips
				case "mipsel","mipsle":
					o.arch = mipsel
				case "mips64":
					o.arch = mips64
				case "mips64el","mips64le":
					o.arch = mips64el
				case "ppc","powerpc":
					o.arch = ppc
				case "ppc64","powerpc64":
					o.arch = ppc64
				case "ppc64el","powerpc64el":
					o.arch = ppc64el
				case "s390x":
					o.arch = s390x
				case "riscv64":
					o.arch = riscv64
				case "sparc64":
					o.arch = sparc64
				default:
					exist = false
			}
		}
		if !exist {
			switch o.plat {
				case linux,windows:
					switch runtime.GOARCH {
						case "amd64": o.arch = x86_64
						case "arm64": o.arch = arm64
						case "386":   o.arch = i686
						case "arm":   o.arch = armv7hf
					}
				case darwin,macos,maccatalyst,ios_simulator,watchos_simulator,tvos_simulator:
					switch runtime.GOARCH {
						case "amd64": o.arch = x86_64h
						case "arm64": o.arch = arm64e
						case "i686":  o.arch = x86_64
						case "arm":   o.arch = arm64
					}
				case ios,tvos:
					o.arch = arm64e
				case watchos:
					o.arch = arm64_32
			}
		}
	} else { native = false }

	if o.linker==default_linker {

		val,exist := os.LookupEnv("FLAGS_LINKER")
		if exist {
			switch val {
				case "ld":
					o.linker = ld
				case "bfd":
					o.linker = bfd
				case "gold":
					o.linker = gold
				case "lld":
					o.linker = lld
				default: exist = false
			}
		}

	}

	switch o.arch {
		case x86_64,x86_64h:
			if runtime.GOARCH=="amd64" { o.cross=false }
		case arm64,arm64e:
			if runtime.GOARCH=="arm64" { o.cross=false }
	}
	if native { o.native = true }
	if o.native && o.cross {
		err("クロスコンパイルにおける native 指定は無効です")
	}

}