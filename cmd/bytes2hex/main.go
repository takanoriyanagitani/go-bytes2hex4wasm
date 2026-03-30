package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	bh "github.com/takanoriyanagitani/go-bytes2hex4wasm/conv/wasm/wazero"
	"github.com/tetratelabs/wazero"
)

var (
	ErrUnableToRead error = errors.New("unable to read the wasm memory")
)

func sub(ctx context.Context) error {

	var wasmLoc string = os.Getenv("ENV_WASM_LOC")
	var wasmMax int64 = 1024
	file, err := os.Open(wasmLoc)
	if nil != err {
		return err
	}
	defer file.Close()

	var limited io.Reader = &io.LimitedReader{
		R: file,
		N: wasmMax,
	}

	wasmBytes, err := io.ReadAll(limited)
	if nil != err {
		return err
	}

	var limitPages uint32 = 16
	var cfg wazero.RuntimeConfig = wazero.
		NewRuntimeConfig().
		WithMemoryLimitPages(limitPages)
	var rtm wazero.Runtime = wazero.NewRuntimeWithConfig(
		ctx,
		cfg,
	)

	wrtm := bh.WasmRuntime{Runtime: rtm}
	defer wrtm.Close(ctx)

	compiled, err := wrtm.Compile(ctx, wasmBytes)
	if nil != err {
		return err
	}
	defer compiled.Close(ctx)

	var mcfg wazero.ModuleConfig = wazero.NewModuleConfig()
	instance, err := wrtm.Instantiate(ctx, compiled, mcfg)
	if nil != err {
		return err
	}
	defer instance.Close(ctx)

	conv, err := instance.GetConverter()
	if nil != err {
		return err
	}

	wmem, err := instance.Memory()
	if nil != err {
		return err
	}

	var ibuf [32768]byte
	for {
		cnt, err := io.ReadFull(os.Stdin, ibuf[:])

		if 0 < cnt {
			woriginal, ok := wmem.Memory.Read(0, 32768)
			if !ok {
				return ErrUnableToRead
			}
			copy(woriginal, ibuf[:cnt])

			err = conv.Call(ctx)
			if nil != err {
				return err
			}

			var hlen int = (cnt << 1) & 0xffff_ffff
			wconverted, ok := wmem.Memory.Read(65536, uint32(hlen))
			if !ok {
				return ErrUnableToRead
			}

			_, err = os.Stdout.Write(wconverted)
			if nil != err {
				return err
			}
		}

		switch {
		case nil == err:
			continue
		case io.EOF == err:
			return nil
		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil
		default:
			return err
		}
	}
}

func main() {
	err := sub(context.Background())
	if nil != err {
		log.Printf("%v\n", err)
	}
}
