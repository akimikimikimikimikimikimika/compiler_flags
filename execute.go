package main

import (
	"fmt"
	"strings"
)

func (o *options) execute() {

	switch o.lang {
		case c,cpp:
			o.c_configure()
		case swift:
			o.swift_configure()
		case golang:
			switch o.cc_go {
				case gc:    o.go_configure()
				case gccgo: o.gccgo_configure()
			}
		case rustc,cargo:
			o.rust_configure()
	}

	args := strings.Join(o.flags," ")
	if o.flags_shell { args = fmt.Sprintf("sh -c '%s'",args) }
	fmt.Print(args)

}

func (o *options) c_configure() {

	// compiler/target specification
	switch o.simplified_plat() {
		case darwin:
			var cc string
			switch o.cc {
				case clang:
					switch o.clang {
						case clang_xcode,clang_clt:
							cc = xcrun
							o.apple_platforms()
						case clang_brew,clang_macports,clang_graalvm:
							switch o.plat {
								case darwin,macos:
								default:
									err("このプラットフォームはこのコンパイラでサポートされていません")
							}
							switch o.arch {
								case x86_64,x86_64h:
									o.arch = x86_64
								case arm64,arm64e:
									if o.clang==clang_graalvm {
										err("このアーキテクチャはこのコンパイラでサポートされていません")
									}
									o.arch = arm64
								default:
									err("このアーキテクチャはこのコンパイラでサポートされていません")
							}
							switch o.clang {
								case clang_brew:
									switch o.lang {
										case c:   cc = c_macos_brew
										case cpp: cc = cpp_macos_brew
									}
								case clang_macports:
									switch o.lang {
										case c:   cc = c_macos_macports
										case cpp: cc = cpp_macos_macports
									}
								case clang_graalvm:
									switch o.lang {
										case c:   cc = c_macos_graalvm
										case cpp: cc = cpp_macos_graalvm
									}
							}
							o.add(cc)
						case clang_default:
							err("使用する clang コンパイラが指定されていません")
					}
				case gcc:
					switch o.plat {
						case darwin,macos:
						default:
							err("このプラットフォームはこのコンパイラでサポートされていません")
					}
					switch o.arch {
						case x86_64,x86_64h:
						default:
							err("このアーキテクチャはこのコンパイラでサポートされていません")
					}
					switch o.lang {
						case c:   cc = c_macos_gcc
						case cpp: cc = cpp_macos_gcc
					}
					o.add(cc)
				default:
					err("サポートしていないコンパイラが指定されました")
			}
			if !is_exist(cc) { warn("この環境で該当するコンパイラが見つかりません") }
			if o.linker!=default_linker { err("リンカの指定には対応していません") }
		case linux:
			switch o.cc {
				case gcc:
					var cc,options = o.linux_gcc_target(o.musl,false)
					o.add(cc)
					if options!="" { o.add(options) }
					if o.linker==ld { err("このリンカはサポートされていません") }
				case clang:
					var target,options = o.linux_gcc_target(false,true)
					var cc string
					switch o.lang {
						case c:   cc = c_linux_clang
						case cpp: cc = cpp_linux_clang
					}
					if !is_exist(cc) { warn("この環境で該当するコンパイラが見つかりません") }
					o.add(cc)
					o.add("--target="+target)
					if options!="" { o.add(options) }
					if o.cross {
						o.add(
							"--gcc-toolchain=/usr",
							"--sysroot","/usr/"+target,
						)
						if o.lang==cpp {
							o.add(
								"-stdlib=libstdc++",
								"-I","/usr/"+target+"/include/c++/10",
								"-I","/usr/"+target+"/include/c++/10/"+target,
								"-I","/usr/"+target+"/include/c++/10/backward",
							)
						}
					}
					switch o.linker {
						case bfd,gold:
							err("このリンカはサポートされていません")
					}
				default:
					err("サポートしていないコンパイラが指定されました")
			}

		case windows:
			switch o.cc {
				case msvc:
					switch o.arch {
						case x86_64:
							o.add("cl-x64.bat")
						case arm64:
							o.add("cl-arm64.bat")
						case i686,i586:
							o.add("cl-x86.bat")
						default:
							err("サポートしていないアーキテクチャが指定されました")
					}
					if o.linker!=default_linker { err("リンカの指定には対応していません") }
				case clang:
					o.add("clang.exe")
					switch o.arch {
						case x86_64:
							o.add("--target=x86_64-pc-windows-msvc")
						case arm64:
							o.add("--target=arm64-pc-windows-msvc")
						case i686:
							o.add("--target=i686-pc-windows-msvc")
						case i586:
							o.add("--target=i586-pc-windows-msvc")
						default:
							err("サポートしていないアーキテクチャが指定されました")
					}
					if o.linker!=default_linker { err("リンカの指定には対応していません") }
				case gcc:
					var cc string
					switch o.arch {
						case x86_64:
							cc = "x86_64-w64-mingw32"
						case i686,i586:
							cc = "i686-w64-mingw32"
						default:
							err("サポートしていないアーキテクチャが指定されました")
					}
					if o.cc!=gcc { err("サポートしていないコンパイラが指定されました") }
					switch o.lang {
						case c:   cc+="-gcc"
						case cpp: cc+="-g++"
					}
					o.add(cc)
				default:
					err("サポートしていないコンパイラが指定されました")
			}
		case wasm:
			if is_exist(emscripten_macos) {
				o.add(emscripten_macos)
			} else {
				o.add("emcc")
				warn("この環境で該当するコンパイラが見つかりません")
			}
			if o.linker!=default_linker { err("リンカの指定には対応していません") }
		default:
			err("サポートしていないプラットフォームが指定されました")
	}

	if o.musl {
		if o.plat!=linux {
			err("musl はこのプラットフォームではサポートしていません")
		}
		if o.cc!=gcc {
			err("musl は GCC のみで使用できます")
		}
	}

	if o.native {
		switch o.plat {
			case darwin,macos,maccatalyst,linux:
				o.add("-march=native")
			case windows:
				if o.cc!=msvc { o.add("-march=native") }
			default: err("native 指定は無効です")
		}
	}

	// std / objective-c
	switch o.lang {

		case c:

			if o.cc!=msvc {

				if o.output!=show_macros {
					switch o.std_c {
						case   c90: o.add("-std=c90")
						case gnu90: o.add("-std=gnu90")
						case   c99: o.add("-std=c99")
						case gnu99: o.add("-std=gnu99")
						case   c11: o.add("-std=c11")
						case gnu11: o.add("-std=gnu11")
						case   c17: o.add("-std=c17")
						case gnu17: o.add("-std=gnu17")
						case   c2x: o.add("-std=c2x")
						case gnu2x: o.add("-std=gnu2x")
					}
				}

				if o.objc { o.add("-ObjC","-lobjc") }

			} else {
				o.add("/TC")
				switch o.std_c {
					case c11: o.add("/std:c11")
					case c17: o.add("/std:c17")
					default:
						err("未対応の C 標準が指定されました")
				}
			}

		case cpp:

			if o.cc!=msvc {

				if o.output!=show_macros {
					switch o.std_cpp {
						case   cpp98: o.add("-std=c++98")
						case gnupp98: o.add("-std=gnu++98")
						case   cpp11: o.add("-std=c++11")
						case gnupp11: o.add("-std=gnu++11")
						case   cpp14: o.add("-std=c++14")
						case gnupp14: o.add("-std=gnu++14")
						case   cpp17: o.add("-std=c++17")
						case gnupp17: o.add("-std=gnu++17")
						case   cpp20: o.add("-std=c++20")
						case gnupp20: o.add("-std=gnu++20")
						case   cpp23: o.add("-std=c++2b")
						case gnupp23: o.add("-std=gnu++2b")
					}
				}

				if o.cc!=gcc {
					switch o.stdlib {
						case libcpp: o.add("-stdlib=libc++")
						case libstdcpp: o.add("-stdlib=libstdc++")
					}
				}

				if o.objc { o.add("-ObjC++","-lobjc") }
				if o.cppfs { o.add("-lstdc++fs") }

			} else {
				o.add("/TP")
				switch o.std_cpp {
					case cpp14: o.add("/std:c++14")
					case cpp17: o.add("/std:c++17")
					case cpp20: o.add("/std:c++20")
					case cpp23: o.add("/std:c++latest")
					default:
						err("未対応の C++ 標準が指定されました")
				}
			}

	}

	switch o.output {
		case shared:
			if o.cc==msvc { err("このコンパイラは共有ライブラリの出力はできません") }
			o.add("-shared")
		case dylib:
			if o.cc==msvc {
				o.add("/LD")
			} else {
				o.add("-dynamiclib")
			}
		case bundle:
			if o.cc==msvc { err("このコンパイラはバンドルの出力はできません") }
			o.add("-bundle")
		case object: // -c (preprocess,compile,assemble)
			if o.cc==msvc {
				o.add("/c")
			} else {
				o.add("--compile")
			}
		case assembly_att: // -S (preprocess,compile)
			if o.cc==msvc { err("このコンパイラはアセンブリの出力はできません") }
			o.add("--assemble")
			if o.flags_level!=loose { o.add("-fverbose-asm") } // 解読に便利なコメントを付ける
		case assembly_intel:
			if o.cc==msvc { err("このコンパイラはアセンブラの出力はできません") }
			if o.arch!=x86_64 && o.arch!=x86_64h {
				err("サポートしていないアーキテクチャが指定されました")
			}
			if o.cc==gcc {
				if o.simplified_plat()==darwin {
					err("Darwin で gcc コンパイラの Intel 式アセンブリ出力はできません")
				}
				o.add("--assemble","-masm=intel")
			}
			o.add("--assemble","-mllvm","--x86-asm-syntax=att")
		case preprocess: // -E (preprocess)
			if o.cc==msvc {
				o.add("/E")
			} else {
				o.add("--preprocess")
			}
		case precompile:
			if o.cc==msvc { err("このコンパイラはプリコンパイルには非対応です") }
			o.add("--precompile")
		case llvm_bytecode:
			if o.cc==gcc || o.cc==msvc { err("このコンパイラで LLVM ビットコード の出力はできません") }
			o.add("-emit-llvm","-c")
		case llvm_ir:
			if o.cc==gcc || o.cc==msvc { err("このコンパイラで LLVM IR の出力はできません") }
			o.add("-emit-llvm","-S")
		case syntax_check:
			if o.cc==msvc {
				o.add("/Zs")
			} else {
				o.add("-fsyntax-only")
			}
		case show_macros:
			if o.cc==msvc { err("このコンパイラは定義済みマクロの出力はできません") }
			o.add("-dM","-E","-","<","/dev/null")
			o.flags_shell = true
			return
	}

	// flags
	if o.flags_level != loose {
		if o.cc!=msvc {
			o.add(
				"-Werror", // 警告を全てエラーとする
				"-Wall", // 基本的な警告
				"-Wextra", // 追加の警告
				"-Wshadow", // 隠れ変数を生じさせると警告
				"-Wunreachable-code", // 明らかに到達できないコードがあれば警告
				"-Wpointer-arith", // void* など怪しいポインタ演算を警告する
				"-Wwrite-strings", // const char* を char* に変換すると警告
				"-Winit-self", // int i = i; のような宣言を警告
				"-Wfloat-equal", // 浮動小数を == で比較すると警告
				"-Wimplicit-fallthrough", // case 文の break 忘れを警告
				"-Wsign-compare", // signed の値を unsigned の値と比較すると警告
				"-Wredundant-decls", // 同じ変数を複数回宣言すると警告
				"-Wdisabled-optimization", // 最適化に失敗した場合に警告
				"-Wchar-subscripts", // 配列添字が負値をとりうる char 型になる場合に警告
				"-fno-ident", // コンパイラの情報をバイナリに付加しない
				"-fno-common", // グローバル変数を .bss セクションに配置
				"-fno-strict-aliasing", // int を short* に代入するなど型名が不適切な代入を無効にする
				"-ffunction-sections","-fdata-sections",
				"-fPIC", // アドレス空間のランダム化 (or -fPIE)
				"-pipe", // 一時ファイルを介さず、パイプを使ってデータを渡す
			)
			if o.cc != gcc {
				o.add(
					"-Wcomma","-Wunreachable-code-return","-fcolor-diagnostics",
				)
			}
		} else {
			o.add(
				"/nologo", // ロゴを出力しない
				"/MP", // ソースコードの処理を並列に実行します
				"/WX", // 警告を全てエラーとする
				"/Wall", // 基本的な警告
				"/sdl", // 追加のセキュリティ機能と警告を有効にします
				"/guard:cf",
				"/guard:ehcont",
				"/analyze",
				"/EHsc", // C++ 例外処理を有効にします + extern "C" では、既定値が nothrow に設定されます
			)
		}
	}
	if o.flags_level == strict && o.cc!=msvc {
		o.add(
			"-Wundef", // 未定義のマクロ変数に対して警告
			"-Wunused", // 定義されていても未使用であれば警告
			"-Wpadded", // 構造体でパディングが発生した場合警告
			"-Wstrict-prototypes", // 引数の型を明示しない関数があれば警告
			"-Wconversion", // unsigned x = -1 のような値が変わってしまうものに警告
			"-Wdouble-promotion", // float 型が double 型に変換されているときに警告
			"-Wcast-align=strict", // char* -> int* のようにポインタのアライメントが変わるようなキャストを警告
			"-Waggregate-return", // 構造体を返す関数を呼び出せば警告
			"-Wswitch-default", // default 節のない switch に警告
			"-Wswitch-enum", // enum に対する switch で全ての値を尽くしていないものに警告
			"-Wfloat-overflow-conversion",
			"-Wfloat-zero-conversion",
		)
	}

	o.optimize_flags()

	if o.cc!=msvc {

		// linker
		switch o.output {
			case executable,shared,dylib,bundle:
				switch o.linker {
					case default_linker:
					case ld:   o.add("-fuse-ld=ld")
					case bfd:  o.add("-fuse-ld=bfd")
					case gold: o.add("-fuse-ld=gold")
					case lld:  o.add("-fuse-ld=lld")
				}
		}

		if o.userlib {
			if !is_exist("/usr/local/include") || !is_exist("/usr/local/lib") {
				o.add("-I/usr/local/include","-L/usr/local/lib")
			} else { warn("この環境でユーザーライブラリが見つかりませんでした") }
		}

		if o.protection {
			o.add(
				"-ftrapv", // 数値オーバーフローが発生すれば実行を終了する
				"-fstack-protector-all", // スタック領域範囲外への書き込みをチェックする
				"-D_FORTIFY_SOURCE=2", // バッファオーバーフローを検出する
				"-D_GLIBCXX_ASSERTIONS", // C++ コンテナ型の境界チェックを行う
				"-fsanitize=address,thread,undefined,integer", // メモリアドレス,スレッド,未定義動作,整数オーバーフローのサニタイザを有効にする
				"-fno-sanitize-recover=all", // サニタイザがエラーを見つけたら終了
			)
		}
		if o.dry_run { o.add("-###") }
		if o.stack_usage { o.add("-fstack-usage") }
		if o.math { o.add("-lm") }
		if o.boost_random { o.add("-lboost_random") }
		if o.thread {
			o.add("-pthread")
			if o.plat==linux { o.add("-lrt") }
		}
		if o.posix  { o.add("-D_POSIX_C_SOURCE=199309L") }
		if o.no_unused_result { o.add("-Wno-unused-result") }
		if o.opencl { o.add("-framework","OpenCL","-DCL_SILENCE_DEPRECATION") }


		if o.link_type==default_link_type && o.musl {
			o.link_type = statically
		}
		if o.link_type==statically { o.add("-static") }

	} else {

		o.add(
			"/source-charset:utf-8", // 入力ファイルのエンコードを UTF-8 とする
			"/execution-charset:utf-8", // 出力結果のエンコードを UTF-8 とする
		)

	}

	if o.openmp {
		switch o.cc {
			case clang:
				o.add("-Xpreprocessor","-fopenmp","-lomp")
			case gcc:
				o.add("-fopenmp","-lgomp")
			case msvc:
				o.add("/openmp:llvm") // or /openmp
		}
	}

}

func (o *options) swift_configure() {

	if !is_exist(xcrun) { warn("この環境で該当するコンパイラが見つかりません") }

	o.apple_platforms()
	o.optimize_flags()

}

func (o *options) go_configure() {

	var os string
	var arch string
	var option string = ""

	switch o.plat {
		case darwin,macos,ios:
			if o.plat!=ios { os = "darwin" } else { os = "ios" }
			switch o.arch {
				case x86_64,x86_64h:
					arch = "amd64"
				case arm64,arm64e:
					arch = "arm64"
				default:
					err("サポートしていないアーキテクチャが指定されました")
			}
		case linux:
			os = "linux"
			switch o.arch {
				case x86_64:
					arch = "amd64"
				case i686:
					arch = "386"
					option = "GO386=sse2"
				case i586:
					arch = "386"
					option = "GO386=softfloat"
				case arm64:
					arch = "arm64"
				case armv7,armv7hf,armv7el:
					arch = "arm"
					option = "GOARM=7"
				case armv6,armv6hf,armv6el:
					arch = "arm"
					option = "GOARM=6"
				case armv5:
					arch = "arm"
					option = "GOARM=5"
				case mips64:
					arch = "mips64"
				case mips64el:
					arch = "mips64le"
				case mips:
					arch = "mips"
				case mipsel:
					arch = "mipsle"
				case ppc64:
					arch = "ppc64"
				case ppc64el:
					arch = "ppc64le"
				case s390x:
					arch = "s390x"
				case riscv64:
					arch = "riscv64"
				default:
					err("サポートしていないアーキテクチャが指定されました")
			}
		case windows:
			os = "windows"
			switch o.arch {
				case x86_64:
					arch = "amd64"
				case i686,i586:
					arch = "386"
				case arm64:
					arch = "arm64"
				case armv7,armv7hf,armv7el:
					arch = "arm"
				default:
					err("サポートしていないアーキテクチャが指定されました")
			}
		case wasm:
			os = "js"
			arch = "wasm"
		default:
			err("サポートしていないプラットフォームが指定されました")
	}

	o.add("env","GOOS="+os,"GOARCH="+arch)
	if option!="" { o.add(option) }
	if is_exist(go_macos) {
		o.add(go_macos)
	} else if is_exist(go_linux) {
		o.add(go_linux)
	} else { warn("この環境で該当するコンパイラが見つかりません") }
	o.add("build","-ldflags=-s","-ldflags=-w")

}

func (o *options) gccgo_configure() {

	// compiler/target specification
	if o.plat!=linux { err("サポートしていないプラットフォームが指定されました") }
	if o.cc!=gcc { err("サポートしていないコンパイラが指定されました") }
	var cc,options = o.linux_gcc_target(o.musl,false)
	o.add(cc)
	if options!="" { o.add(options) }

	if o.native {
		switch o.plat {
			case darwin,macos,maccatalyst,linux:
				o.add("-march=native")
			default: err("native 指定は無効です")
		}
	}

	switch o.output {
		case shared:        o.add("-shared")
		case dylib:         o.add("-dynamiclib")
		case bundle:        o.add("-bundle")
		case object: // -c (preprocess,compile,assemble)
			o.add("--compile")
		case assembly_att: // -S (preprocess,compile)
			o.add("--assemble")
		case assembly_intel:
			err("Intel 式アセンブリ出力はできません")
		case preprocess: // -E (preprocess)
			o.add("--preprocess")
		case llvm_bytecode:
			err("このコンパイラで LLVM ビットコード の出力はできません")
		case llvm_ir:
			err("このコンパイラで LLVM IR の出力はできません")
		case precompile:    o.add("--precompile")
		case syntax_check:  o.add("-fsyntax-only")
		case show_macros:
			o.add("-dM","-E","-","<","/dev/null")
			o.flags_shell = true
			return
	}

	// flags
	if o.flags_level != loose {
		o.add(
			"-Werror", // 警告を全てエラーとする
			"-Wall", // 基本的な警告
			"-fPIC",
			"-pipe", // 一時ファイルを介さず、パイプを使ってデータを渡す
		)
	}
	if o.flags_level == strict {
		o.add(
			"-Wextra", // 追加の警告
		)
	}

	switch o.linker {
		case ld:   o.add("-fuse-ld=ld")
		case bfd:  o.add("-fuse-ld=bfd")
		case gold: o.add("-fuse-ld=gold")
		case lld:  o.add("-fuse-ld=lld")
	}

	if o.userlib {
		if !is_exist("/usr/local/include") || !is_exist("/usr/local/lib") {
			o.add("-I/usr/local/include","-L/usr/local/lib")
		} else { warn("この環境でユーザーライブラリが見つかりませんでした") }
	}

	o.optimize_flags()

	if o.dry_run { o.add("-###") }
	if o.stack_usage { o.add("-fstack-usage") }
	if o.no_unused_result { o.add("-Wno-unused-result") }
	if o.link_type==statically { o.add("-static") }

	if o.musl {
		err("Go で musl は現在のところサポートされていません")
	}

}

func (o *options) rust_configure() {

	var target string = ""
	var linker string = ""
	if o.output!=syntax_check {
		if o.musl && o.plat!=linux {
			err("musl はこのプラットフォームではサポートしていません")
		}
		switch o.plat {
			case darwin,macos:
				switch o.cc {
					case clang:
						switch o.arch {
							case x86_64,x86_64h:
								target = "x86_64-apple-darwin"
							case arm64,arm64e:
								if o.clang==clang_graalvm {
									err("サポートしていないアーキテクチャが指定されました")
								}
								target = "aarch64-apple-darwin"
							default:
								err("サポートしていないアーキテクチャが指定されました")
						}
						switch o.clang {
							case clang_xcode:
								o.add("env","DEVELOPER_DIR="+xcode_dev_dir)
								linker = "/usr/bin/clang"
							case clang_clt:
								o.add("env","DEVELOPER_DIR="+clt_dev_dir)
								linker = "/usr/bin/clang"
							case clang_brew:     linker = c_macos_brew
							case clang_macports: linker = c_macos_macports
							case clang_graalvm:  linker = c_macos_graalvm
							case clang_default:
								err("使用する clang コンパイラが指定されていません")
						}
					case gcc:
						switch o.arch {
							case x86_64,x86_64h:
								target="x86_64-apple-darwin"
							default:
								err("サポートしていないアーキテクチャが指定されました")
						}
						linker = c_macos_gcc
					default:
						err("サポートしていないコンパイラが指定されました")
				}
			case ios:
				switch o.arch {
					case x86_64:
						target = "x86_64-apple-ios"
					case arm64:
						target = "aarch64-apple-ios"
					default:
						err("サポートしていないアーキテクチャが指定されました")
				}
				if o.cc!=clang {
					err("このプラットフォームはこのコンパイラでサポートされていません")
				}
				switch o.clang {
					case clang_xcode: linker = c_macos_xcode
					case clang_clt:   linker = c_macos_clt
					case clang_default:
						err("使用する clang コンパイラが指定されていません")
					default:
						err("このプラットフォームはこのコンパイラでサポートされていません")
				}
			case linux:
				target = o.linux_rust_target()
				linker,_ = o.linux_gcc_target(false,false)
			case windows:
				switch o.arch {
					case x86_64:
						switch o.cc {
							case msvc:
								target = "x86_64-pc-windows-msvc"
								linker = "cl-x64.bat"
							case clang:
								target = "x86_64-pc-windows-msvc"
								linker = "clang.exe"
							case gcc:
								target = "x86_64-pc-windows-gnu"
								linker = "x86_64-w64-mingw32-gcc"
						}
					case i686:
						switch o.cc {
							case msvc:
								target = "i686-pc-windows-msvc"
								linker = "cl-x86.bat"
							case clang:
								target = "i686-pc-windows-msvc"
								linker = "clang.exe"
							case gcc:
								target = "i686-pc-windows-gnu"
								linker = "i686-w64-mingw32-gcc"
						}
					case i586:
						switch o.cc {
							case msvc:
								target = "i586-pc-windows-msvc"
								linker = "cl-x86.bat"
							case clang:
								target = "i586-pc-windows-msvc"
								linker = "clang.exe"
							case gcc:
								target = "i586-pc-windows-gnu"
								linker = "i686-w64-mingw32-gcc"
						}
					case arm64:
						switch o.cc {
							case msvc:
								target = "aarch64-pc-windows-msvc"
								linker = "cl-arm64.bat"
							case clang:
								target = "aarch64-pc-windows-msvc"
								linker = "clang.exe"
							case gcc:
								err("このアーキテクチャはこのコンパイラでサポートされていません")
						}
					default:
						err("サポートしていないアーキテクチャが指定されました")
				}
			case wasm:
				target="wasm32-unknown-unknown"
				linker = emscripten_macos
			default:
				err("サポートしていないプラットフォームが指定されました")
		}
	}

	switch o.lang {
		case rustc:
			rustc := expand(rustc_rustup)
			if !is_exist(rustc) { warn("この環境で rustc が見つかりません") }
			o.add(
				rustc,
				"--target="+target,
			)
		case cargo:
			cargo := expand(cargo_rustup)
			if !is_exist(cargo) { warn("この環境で cargo が見つかりません") }
			if o.output==syntax_check {
				o.add(
					cargo,"check",
					"--target-dir","target.nosync",
				)
			} else {
				o.add(
					cargo,"rustc","--bins",
					"--target-dir","target.nosync",
					"--target="+target,
				)
				switch o.optimize {
					case normal,lto,faster,smaller:
						o.add("--release")
				}
				o.add("--")
			}
	}

	if o.output!=syntax_check {

		switch o.output {
			case executable:
				if o.flags_level!=loose {
					o.add("--emit","link")
				}
			case object:        o.add("--emit","obj")
			case assembly_att:  o.add("--emit","asm")
			case llvm_bytecode: o.add("--emit","llvm-bc")
			case llvm_ir:       o.add("--emit","llvm-ir")
			default:
				err("この出力形式は指定できません")
		}

		switch o.optimize {
			case debug:
				if o.flags_level!=loose {
					o.add(
						"-C","opt-level=0",
						"-C","debuginfo=2",
						"-C","lto=off",
					)
				}
			case unoptimized:
				o.add(
					"-C","opt-level=0",
					"-C","debuginfo=0",
					"-C","lto=off",
				)
			case normal:
				if o.flags_level!=loose {
					o.add(
						"-C","opt-level=2",
						"-C","debuginfo=0",
						// "-C","lto=thin",
					)
				}
			case lto:
				o.add(
					"-C","opt-level=2",
					"-C","debuginfo=0",
					"-C","lto=fat",
				)
			case faster:
				o.add(
					"-C","opt-level=3",
					"-C","debuginfo=0",
					"-C","lto=fat",
				)
			case smaller:
				o.add(
					"-C","opt-level=z",
					"-C","debuginfo=0",
					"-C","lto=thin",
				)
		}

		if o.flags_level!=loose {
			o.add(
				"-C","overflow-checks=yes",
				"-C","relocation-model=pic",
				"-C","linker="+linker,
			)
		}

		if o.native {
			o.add("-C","target-cpu=native")
		}

		switch o.link_type {
			case statically:
				o.add("-C","prefer-dynamic=no")
			case dynamically:
				o.add("-C","prefer-dynamic=yes")
		}

	}

}

func (o *options) apple_platforms() {
	var dev_dir string = ""
	var plat string = ""
	var sdk string  = ""
	var arch string = ""
	switch o.clang {
		case clang_xcode:
			if !is_exist(xcode_dev_dir) { warn("Xcode が見つかりません") }
			dev_dir = xcode_dev_dir
		case clang_clt:
			if !is_exist(clt_dev_dir) { warn("Command Line Tools が見つかりません") }
			dev_dir = clt_dev_dir
			switch o.plat {
				case darwin,macos:
				default:
					err("このプラットフォームはこのコンパイラでサポートされていません")
			}
		default:
			err("Xcode 又は Command Line Tools を使う必要があります")
	}
	switch o.plat {
		case macos,maccatalyst,ios_simulator,watchos_simulator,tvos_simulator,driverkit:
			switch o.arch {
				case x86_64,x86_64h,arm64,arm64e,universal:
				default:
					err("サポートしていないアーキテクチャが指定されました")
			}
		case ios,watchos,tvos:
			switch o.arch {
				case arm64,arm64e,armv7,armv7s,armv7k,armv7m,armv7em,armv6,armv6m,arm64_32:
				default:
					err("サポートしていないアーキテクチャが指定されました")
			}
	}
	switch o.plat {
		case darwin:
			plat="darwin"
		case macos:
			plat="macos"
			sdk="macosx"
		case maccatalyst:
			plat="ios-macabi"
			sdk="macosx"
		case driverkit:
			plat="macos"
			sdk="driverkit"
		case ios:
			plat="ios"
			sdk="iphoneos"
		case ios_simulator:
			plat="ios"
			sdk="iphonesimulator"
		case watchos:
			plat="watchos"
			sdk="watchos"
		case watchos_simulator:
			plat="watchos"
			sdk="watchsimulator"
		case tvos:
			plat="tvos"
			sdk="appletvos"
		case tvos_simulator:
			plat="tvos"
			sdk="appletvsimulator"
	}
	o.add("env","DEVELOPER_DIR="+dev_dir,"xcrun")
	if sdk!="" { o.add("--sdk",sdk) }
	switch o.lang {
		case c:     o.add("clang")
		case cpp:   o.add("clang++")
		case swift: o.add("swiftc")
	}
	if o.arch==universal {
		o.add("-arch","x86_64h","-arch","arm64e","--target=apple-"+plat)
	} else {
		switch o.arch {
			case x86_64:   arch="x86_64"
			case x86_64h:  arch="x86_64h"
			case arm64:    arch="arm64"
			case arm64e:   arch="arm64e"
			case arm64_32: arch="arm64_32"
			case armv7:    arch="armv7"
			case armv7s:   arch="armv7s"
			case armv7k:   arch="armv7k"
			case armv7m:   arch="armv7m"
			case armv7em:  arch="armv7em"
			case armv6:    arch="armv6"
			case armv6m:   arch="armv6m"
		}
		o.add("--target="+arch+"-apple-"+plat)
	}
}

func (o *options) linux_gcc_target(musl bool,target_only bool) (string,string) {

	var target = ""
	var cc = ""
	var bin = "/usr/bin/"
	var options = ""

	if !musl {
		if !o.native {
			switch o.arch {
				case x86_64:
					target = "x86_64-linux-gnu"
				case i686:
					target = "i686-linux-gnu"
					options = "-march=i686"
				case i586:
					target = "i686-linux-gnu"
					options = "-march=i586"
				case x32:
					target = "x86_64-linux-gnux32"
				case arm64:
					target = "aarch64-linux-gnu"
					options = "-mlittle-endian"
				case arm64_be:
					target = "aarch64-linux-gnu"
					options = "-mbig-endian"
				case armv7hf,armv7,armv6hf,armv6,armv7el,armv6el,armv5:
					switch o.arch {
						case armv7hf,armv7,armv6hf,armv6:
							target = "arm-linux-gnueabihf"
						case armv7el,armv6el,armv5:
							target = "arm-linux-gnueabi"
					}
					switch o.arch {
						case armv7,armv7hf,armv7el:
							options = "-march=armv7"
						case armv6,armv6hf,armv6el:
							options = "-march=armv6"
						case armv5:
							options = "-march=armv5"
					}
				case mips64:
					target = "mips64-linux-gnuabi64"
				case mips64el:
					target = "mips64el-linux-gnuabi64"
				case mips:
					target = "mips-linux-gnu"
				case mipsel:
					target = "mipsel-linux-gnu"
				case s390x:
					target = "s390x-linux-gnu"
				case ppc:
					target = "powerpc-linux-gnu"
				case ppc64:
					target = "powerpc64-linux-gnu"
				case ppc64el:
					target = "powerpc64le-linux-gnu"
				case riscv64:
					target = "riscv64-linux-gnu"
				case sparc64:
					target = "sparc64-linux-gnu"
				default:
					err("サポートしていないアーキテクチャが指定されました")
			}
		}
	} else {
		switch o.arch {
			case x86_64:
				target = "x86_64-linux-musl"
			case i686,i586:
				target = "i686-linux-musl"
			case x32:
				target = "x86_64-linux-muslx32"
			case arm64:
				target = "aarch64-linux-musl"
			case arm64_be:
				target = "aarch64_be-linux-musl"
			case armv7hf,armv7:
				target = "armv7l-linux-musleabihf"
			case armv6hf,armv6:
				target = "armv6-linux-musleabihf"
			case armv7el:
				target = "armv7m-linux-musleabi"
			case armv6el:
				target = "armv6-linux-musleabi"
			case armv5:
				target = "armv5-linux-musleabi"
			case mips64:
				target = "mips64-linux-musl"
			case mips64el:
				target = "mips64el-linux-musl"
			case mips:
				target = "mips-linux-musl"
			case mipsel:
				target = "mipsel-linux-musl"
			case s390x:
				target = "s390x-linux-musl"
			case ppc:
				target = "powerpc-linux-musl"
			case ppc64:
				target = "powerpc64-linux-musl"
			case ppc64el:
				target = "powerpc64le-linux-musl"
			case riscv64:
				target = "riscv64-linux-musl"
			case sparc64:
				target = "sparc64-linux-musl"
			default:
				err("サポートしていないアーキテクチャが指定されました")
		}
	}

	if target_only {
		 return target,options
	}

	switch o.lang {
		case c,cargo,rustc: cc = "gcc"
		case cpp:           cc = "g++"
		case golang:        cc = "gccgo"
	}

	if target!="" {
		bin += target + "-" + cc
	} else {
		bin += cc
	}

	if is_exist(bin+"-12") {
		bin += "-12"
	} else if is_exist(bin+"-11") {
		bin += "-11"
	} else if is_exist(bin+"-10") {
		bin += "-10"
	} else if is_exist(bin+"-9") {
		bin += "-9"
	} else if !is_exist(bin) {
		warn("この環境で該当するコンパイラが見つかりません")
	}

	return bin,options
}

func (o *options) linux_rust_target() string {

	if !o.musl {
		switch o.arch {
			case x86_64:
				return "x86_64-unknown-linux-gnu"
			case i686:
				return "i686-unknown-linux-gnu"
			case i586:
				return "i586-unknown-linux-gnu"
			case x32:
				return "x86_64-unknown-linux-gnux32"
			case arm64:
				return "aarch64-unknown-linux-gnu"
			case armv7hf,armv7:
				return "armv7-unknown-linux-gnueabihf"
			case armv7el:
				return "armv7-unknown-linux-gnueabi"
			case armv6hf,armv6:
				return "arm-unknown-linux-gnueabihf"
			case armv6el:
				return "arm-unknown-linux-gnueabi"
			case armv5:
				return "armv5te-unknown-linux-gnueabi"
			case mips64:
				return "mips64-unknown-linux-gnuabi64"
			case mips64el:
				return "mips64el-unknown-linux-gnuabi64"
			case mips:
				return "mips-unknown-linux-gnu"
			case mipsel:
				return "mipsel-unknown-linux-gnu"
			case s390x:
				return "s390x-unknown-linux-gnu"
			case ppc:
				return "powerpc-unknown-linux-gnu"
			case ppc64:
				return "powerpc64-unknown-linux-gnu"
			case ppc64el:
				return "powerpc64le-unknown-linux-gnu"
			case riscv64:
				return "riscv64gc-unknown-linux-gnu"
			case sparc64:
				return "sparc64-unknown-linux-gnu"
			default:
				err("サポートしていないアーキテクチャが指定されました")
		}
	} else {
		switch o.arch {
			case x86_64:
				return "x86_64-unknown-linux-musl"
			case i686:
				return "i686-unknown-linux-musl"
			case i586:
				return "i586-unknown-linux-musl"
			case arm64:
				return "aarch64-unknown-linux-musl"
			case armv7hf,armv7:
				return "armv7-unknown-linux-musleabihf"
			case armv7el:
				return "armv7-unknown-linux-musleabi"
			case armv6hf,armv6:
				return "arm-unknown-linux-musleabihf"
			case armv6el:
				return "arm-unknown-linux-musleabi"
			case armv5:
				return "armv5te-unknown-linux-musleabi"
			case mips64:
				return "mips64-unknown-linux-muslabi64"
			case mips64el:
				return "mips64el-unknown-linux-muslabi64"
			case mips:
				return "mips-unknown-linux-musl"
			case mipsel:
				return "mipsel-unknown-linux-musl"
			case x32,s390x,ppc,ppc64,ppc64el,riscv64,sparc64:
				err("musl はこのアーキテクチャではサポートしていません")
			default:
				err("サポートしていないアーキテクチャが指定されました")
		}
	}

	return ""
}

func (o *options) optimize_flags() {
	if o.cc!=msvc {
		var debugger bool = false
		var with_lto bool = false
		switch o.optimize {
			case normal:
				o.add("-O2")
			case debug:
				o.add("-O0")
				debugger = true
			case unoptimized:
				o.add("-O0")
			case debug_optimized:
				o.add("-Og")
				debugger = true
			case lto:
				o.add("-O2")
				with_lto = true
			case faster:
				o.add("-Ofast")
				with_lto = true
			case smaller:
				if o.cc==gcc {
					o.add("-Os")
				} else {
					o.add("-Oz")
				}
		}
		if debugger {
			if o.flags_level!=loose {
				o.add(
					"-gfull","-fstandalone-debug","-gdwarf-5","-gz",
					"-fno-omit-frame-pointer",
				)
			} else { o.add("-g") }
		}
		if with_lto {
			o.add(
				"-flto=full", // LTO (リンク時最適化)
				"-fwhole-program-vtables",
				"-fforce-emit-vtables",
				"-fvirtual-function-elimination",
			)
		}
	} else {
		switch o.optimize {
			case normal:
				o.add("/O2")
			case debug:
				o.add("/Od","/Zi")
			case unoptimized:
				o.add("/Od")
			case debug_optimized:
				o.add("/O2","/Zi")
			case lto:
				o.add("/O2","/GL","/link","/LTCG")
			case faster:
				o.add("/Ox","/GL","/link","/LTCG")
			case smaller:
				o.add("/O1")
		}
	}
}

func (o *options) add(args ...string) {
	o.flags = append(o.flags,args...)
}