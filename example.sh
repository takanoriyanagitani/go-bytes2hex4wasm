#!/bin/sh

rtm=wazero
wsm="./bytes2hex_cli.wasm"

ENV_WASM_LOC=/guest.d/read.d/bytes2hex.wasm
ENV_WASM_LOC_NATIVE=./opt.wasm

ex1(){
	echo helo
	printf helo | wrun
	echo
	echo
}

ex2(){
	echo 4095 bytes

	pass='generate-insecure-random-number'

	dd if=/dev/zero bs=4095 count=1 status=none |
		openssl \
			enc \
			-nosalt \
			-aes-256-ctr \
			-pass pass:"${pass}" |
		tail --bytes=16 |
		xxd -ps

	dd if=/dev/zero bs=4095 count=1 status=none |
		openssl \
			enc \
			-nosalt \
			-aes-256-ctr \
			-pass pass:"${pass}" |
		wrun |
		tail --bytes=32

	echo
	echo
}

ex3(){
	echo 1048575 bytes

	pass='generate-insecure-random-number'

	dd if=/dev/zero bs=1048575 count=1 status=none |
		openssl \
			enc \
			-nosalt \
			-aes-256-ctr \
			-pass pass:"${pass}" |
		tail --bytes=16 |
		xxd -ps

	dd if=/dev/zero bs=1048575 count=1 status=none |
		openssl \
			enc \
			-nosalt \
			-aes-256-ctr \
			-pass pass:"${pass}" |
		wrun |
		tail --bytes=32

	echo
}

wzrun(){
  wazero \
    run \
    -env ENV_WASM_LOC="${ENV_WASM_LOC}" \
    -mount "${PWD}:/guest.d/read.d:ro" \
    "${wsm}"
}

wasmtimerun(){
  wasmtime \
    run \
    --env ENV_WASM_LOC="${ENV_WASM_LOC}" \
    --dir "${PWD}::/guest.d/read.d" \
    "${wsm}"
}

wrun(){
  wasmtimerun
}

wrun_native_wasm(){
  ENV_WASM_LOC="${ENV_WASM_LOC_NATIVE}" \
  ./cmd/bytes2hex/bytes2hex
}

bench_wasi_wasm(){
  time dd \
    if=/dev/zero \
    bs=1048576 \
    count=16 \
    status=none |
    openssl \
      enc \
      -nosalt \
      -aes-256-ctr \
      -pass pass:'generate-insecure-random-number' |
    wrun |
    dd \
      of=/dev/null \
      bs=1048576 \
      status=progress
}

bench_native_wasm(){
  time dd \
    if=/dev/zero \
    bs=1048576 \
    count=1024 \
    status=none |
    openssl \
      enc \
      -nosalt \
      -aes-256-ctr \
      -pass pass:'generate-insecure-random-number' |
    wrun_native_wasm |
    dd \
      of=/dev/null \
      bs=1048576 \
      status=progress
}

ex1
ex2
ex3

#bench_wasi_wasm
bench_native_wasm
