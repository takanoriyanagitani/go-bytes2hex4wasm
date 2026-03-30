(module

  ;; bytes 0,1,2, ..., 32767: the original bytes
  ;; bytes 65536 - 131071:    the converted hex
  (memory (export "memory") 2)

  (func (export "bytes2hex_hpage2page")

    (local $optr i32)
    (local $cptr i32)

    (local $chunk v128)
    (local $dup v128)
    (local $shifted v128)
    (local $merged v128)

    (local $ltable v128)

    (local $mask8 v128)
    (local $mask4 v128)

    (local.set $ltable (
      v128.const i32x4 0x33323130 0x37363534 0x62613938 0x66656463
    ))

    (local.set $mask8 (
      v128.const i32x4 0x00ff00ff 0x00ff00ff 0x00ff00ff 0x00ff00ff
    ))

    (local.set $mask4 (
      v128.const i32x4 0x0f0f0f0f 0x0f0f0f0f 0x0f0f0f0f 0x0f0f0f0f
    ))

    (local.set $optr (i32.const 0)) ;; pointer to the original
    (local.set $cptr (i32.const 65536)) ;; pointer to the converted

    (block $exit
      (loop $process

        (br_if $exit (i32.ge_u (local.get $optr) (i32.const 32768)))

        (local.set $chunk (v128.load64_zero (local.get $optr)))

        (local.set
          $dup
          (i8x16.shuffle
            0 0 1 1 2 2 3 3 4 4 5 5 6 6 7 7
            (local.get $chunk)
            (local.get $chunk)
          )
        )

        (local.set
          $shifted
          (i16x8.shr_u (local.get $dup) (i32.const 4))
        )

        (local.set
          $merged
          (v128.and
            (v128.bitselect
              (local.get $shifted)
              (local.get $dup)
              (local.get $mask8)
            )
            (local.get $mask4)
          )
        )

        (v128.store
          (local.get $cptr)
          (i8x16.swizzle (local.get $ltable) (local.get $merged))
        )

        (local.set $optr (i32.add (local.get $optr) (i32.const 8)))
        (local.set $cptr (i32.add (local.get $cptr) (i32.const 16)))

        (br $process)

      )
    )

  )

)
