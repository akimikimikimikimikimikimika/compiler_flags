package main

import (
	"fmt"
	"strings"
	"os"
)

func help() {
	help := strings.ReplaceAll(`
		使い方
		   flags [keyword1] [keyword2] ...

		C / C++ / Go / Rust / Swift のコンパイラフラグを生成します。
		ターゲット、コンパイラごとのフラグや、各種機能を有効にするためのフラグを追加します。

		キーワードによりオプションを追加していきます。
		例えば、 "flags c macos amd64 thread" では
		  * 言語: C
		  * ターゲット: macOS
		  * アーキテクチャ: x86_64 (amd64)
		  * pthread を有効にする
		という指定になり、これを満たすコンパイラコマンドを返します。
		指定を省略すれば、現在の環境に合わせて自動的に設定します。

		Unix 系 OS では次のようにして利用します。
		 `+"`flags c`"+` -o bin src.c
		   → gcc -o bin src.c に相当し、さまざまなオプションが標準で付加されている

		Windows の PowerShell では次のようにサブシェルで利用します。
		 pwsh -Command "$(flags c) /Febin.exe src.c"

		キーワード一覧

		 * 言語/環境を指定
		   c      ... C のコンパイラ
		   cpp    ... C++ のコンパイラ
		   objc   ... Objective-C のコンパイラ
		   objcpp ... Objective-C++ のコンパイラ
		   swift  ... Swift のコンパイラ
		   go     ... Go のコンパイラ
		   gccgo  ... GCC Go
		   rustc  ... Rust のコンパイラ (Rustup を使用)
		   cargo  ... Cargo を使用した Rust のコンパイラ (cargo build に相当, Rustup を使用)

		 * プラットフォームを指定
		   darwin
		   linux
		   mingw
		   windows
		   wasm
		   macos
		   maccatalyst
		   ios
		   ios-simulator
		   watchos
		   watchos-simulator
		   tvos
		   tvos-simulator
		   driverkit

		   環境変数 FLAGS_PLATFORM でも指定できます

		 * アーキテクチャを指定
		   universal          ... macOS ユニバーサルバイナリ (x86_64h,arm64e, C/C++/Swift のみ)
		   native             ... ホスト環境に最適化 (-march=native、デフォルト)
		   musl               ... musl を使用 (Linux の Rust のみ)
		   x86_64,amd64
		   x86_64h
		   x32                ... x86_64 x32 ABI
		   x86,i686,i586,i386
		   arm64,aarch64
		   arm64e
		   arm64_32           ... watchOS
		   armv7hf
		   armv7el
		   armv6hf
		   armv6el
		   armv5
		   mips64
		   mips64el
		   mips
		   mipsel
		   ppc64
		   ppc64el
		   ppc,powerpc
		   s390x
		   riscv64
		   sparc64

		   環境変数 FLAGS_ARCH でも指定できます

		 * musl の使用 (Linux)
		   musl ... glibc++ の代わりに musl を使います

		 * リンクの仕方を指定
		   static  ... 静的リンクを明示する
		   dynamic ... 動的リンクを明示する

		   指定しない場合はコンパイラにより自動で設定されます

		 * 最適化オプションの指定
		   normal / 指定なし ... 一般的に十分な程度で最適化を実行します (-O2)
		   debug             ... デバッグ用です。最適化せず、デバッグ情報を付加します (-O0 -g)
		   unoptimized       ... 最適化を行いません (-O0)
		   debug-optimized   ... デバッグ情報を付加した上で最適化を行います (-Og)
		   lto               ... 最適化した上で LTO を有効化します (-O2 -flto)
		   faster            ... -O2 よりも速度重視で最適化します (-Ofast)
		   smaller           ... -O2 よりもバイナリサイズを小さくするよう最適化します (-Oz)

		   Rust では通常 debug を、リリース版では normal を使用します

		 * C / C++ コンパイラを指定
		   clang       ... Clang コンパイラ (C)
		   clang++     ... Clang コンパイラ (C++)
		   gcc         ... GCC コンパイラ (C)
		   g++         ... GCC コンパイラ (C++)
		   msvc        ... Visual C++ コンパイラ (Windows のみ)

		   macOS ターゲットでは次の値で使用する Clang コンパイラを指定できます
		   xcode    ... Clang として Xcode に付属のコンパイラを使用する
		   clt      ... Clang として Apple CLT に付属の Clang コンパイラを使用する
		   brew     ... Clang として Homebrew でインストールできる LLVM Clang コンパイラを使用する
		   macports ... Clang として MacPorts でインストールできる LLVM Clang コンパイラを使用する
		   graalvm  ... Clang として GraalVM の LLVM Clang コンパイラを使用する

		   環境変数 FLAGS_CC でも指定できます

		 * C / C++ 言語標準を指定
		   c90, gnu90, c99, gnu99, c11, gnu11, c17, gnu17, c2x, gnu2x
		   c++98, gnu++98, c++11, gnu++11, c++14, gnu++14, c++17, gnu++17, c++20, gnu++20, c++23, gnu++23

		   C で C++ 標準を指定した場合、 C++ で C 標準を指定した場合は無効になります

		 * Go コンパイラを指定
		   gc    ... Go 標準のコンパイラ (デフォルト)
		   gccgo ... GCC の Go コンパイラ

		   環境変数 FLAGS_GO でも指定できます

		 * 出力タイプを指定 (主に C / C++ )
		   指定しない   ... 実行ファイル
		   so         ... 共有ライブラリ (.so)
		   dylib      ... macOS ダイナミックライブラリ (.dylib)
		   dll        ... Windows ダイナミックリンクライブラリ (.dll)
		   bundle     ... macOS バンドル (.bundle)
		   obj        ... オブジェクトファイル (.o/.obj)
		   asm        ... アセンブリファイル (.s)
		   asm_intel  ... Intel 式アセンブリファイル (.s)
		   asm        ... アセンブリファイル (.s)
		   preprocess ... プリプロセスの結果を出力
		   precompile ... プリコンパイルの結果を出力
		   llvm-ir    ... LLVM 中間コード (.ll)
		   llvm-byte  ... LLVM バイトコード (.bc)
		   syntax     ... シンタックスチェックのみ行う
		   macros     ... 定義済みマクロを出力 (使用例: `+"`flags c macros` > output.h"+` )

		 * C / C++ オプションの程度
		   strict     ... 警告機能などを可能な限り追加する (過剰な場合があります)
		   loose      ... オプションをほとんど付加しない
		   指定しない ... 標準的なオプションを付加する

		 * C++ 標準ライブラリ
		   libc++
		   libstdc++

		 * リンカを指定
		   ld   ... 標準の ld リンカ
		   bfd  ... GNU BFD リンカ
		   gold ... GNU Gold リンカ
		   lld  ... LLVM lld リンカ

		   環境変数 FLAGS_LINKER でも指定できます

		 * その他のオプション (主に C / C++)
		   debug            ... デバッグ用のビルト (最適化を行わない)
		   release          ... 書き出し用のビルト (最適化を行う)
		   dry-run          ... 内部で実行されるコマンドの表示のみを行う
		   stack-usage      ... スタックトレース
		   sanitize         ... C / C++ でサニタイザを有効にする
		   math             ... Linux で math.h や cmath を使用する
		   openmp,omp       ... OpenMP を使用する
		   boost-random     ... Boost の random ライブラリを使用する
		   pthread          ... POSIX Thread を使用する
		   c++fs            ... C++ の filesystem を使用する
		   opencl           ... macOS で OpenCL を使用する
		   userlib          ... macOS でユーザーライブラリを使用する
		   no-unused-result ... 未使用の変数があっても警告しない
	`,"\t","")
	fmt.Println(help)
	os.Exit(0)
}