package main

type language int
const (
	c language = iota
	cpp
	swift
	golang
	rustc
	cargo
)

type c_standard int
const (
	c90 c_standard = iota
	c99
	c11
	c17
	c2x
	gnu90
	gnu99
	gnu11
	gnu17
	gnu2x
)

type cpp_standard int
const (
	cpp98 cpp_standard = iota
	cpp11
	cpp14
	cpp17
	cpp20
	cpp23
	gnupp98
	gnupp11
	gnupp14
	gnupp17
	gnupp20
	gnupp23
)

type compiler int
const (
	default_compiler compiler = iota
	clang
	gcc
	msvc
)

type go_implementation int
const (
	gc go_implementation = iota
	gccgo
)

type clang_type int
const (
	clang_default clang_type = iota
	clang_xcode
	clang_clt
	clang_macports
	clang_brew
	clang_graalvm
)

type architecture int
const (
	default_arch architecture = iota
	universal
	x86_64
	x86_64h
	i686
	i586
	x32
	arm64
	arm64_be
	arm64e
	arm64_32
	armv7
	armv7hf
	armv7el
	armv7s
	armv7k
	armv7m
	armv7em
	armv6
	armv6hf
	armv6el
	armv6m
	armv5
	mips
	mipsel
	mips64
	mips64el
	s390x
	ppc64
	ppc64el
	ppc
	riscv64
	sparc64
)

type platform int
const (
	default_plat platform = iota
	darwin
	macos
	maccatalyst
	ios
	ios_simulator
	watchos
	watchos_simulator
	tvos
	tvos_simulator
	driverkit
	linux
	windows
	wasm
)

type output_type int
const (
	executable output_type = iota
	shared
	dylib
	bundle
	object
	assembly_att
	assembly_intel
	preprocess
	precompile
	llvm_ir
	llvm_bytecode
	syntax_check
	show_macros
)

type stdlib_type int
const (
	default_stdlib stdlib_type = iota
	libstdcpp
	libcpp
)

type linker_type int
const (
	default_linker linker_type = iota
	ld
	bfd
	lld
	gold
)

type flags_level_type int
const (
	loose = iota
	neutral
	strict
)

type link_lib_type int
const (
	default_link_type = iota
	statically
	dynamically
)

type optimize_type int
const (
	normal = iota
	unoptimized
	debug
	debug_optimized
	lto
	faster // -0fast, LTO
	smaller
)

type options struct {
	lang             language
	objc             bool
	cc               compiler
	cc_go            go_implementation
	clang            clang_type
	arch             architecture
	native           bool
	plat             platform
	std_c            c_standard
	std_cpp          cpp_standard
	stdlib           stdlib_type
	linker           linker_type
	output           output_type
	flags_level      flags_level_type
	link_type        link_lib_type
	optimize         optimize_type
	musl             bool
	dry_run          bool
	stack_usage      bool
	protection       bool
	math             bool
	openmp           bool
	boost_random     bool
	thread           bool
	posix            bool
	cppfs            bool
	opencl           bool
	userlib          bool
	no_unused_result bool
	cross            bool
	flags            []string
	flags_shell      bool
}

func init_options() options {
	return options{
		lang:             c,
		objc:             false,
		cc:               default_compiler,
		cc_go:            gc,
		clang:            clang_default,
		arch:             default_arch,
		native:           false,
		plat:             default_plat,
		std_c:            c17,
		std_cpp:          cpp20,
		stdlib:           default_stdlib,
		linker:           default_linker,
		output:           executable,
		flags_level:      neutral,
		link_type:        default_link_type,
		optimize:         normal,
		musl:             false,
		dry_run:          false,
		stack_usage:      false,
		protection:       false,
		math:             false,
		openmp:           false,
		boost_random:     false,
		thread:           false,
		posix:            false,
		cppfs:            false,
		opencl:           false,
		userlib:          false,
		no_unused_result: false,
		cross:            true,
		flags:            []string{},
		flags_shell:      false,
	}
}
