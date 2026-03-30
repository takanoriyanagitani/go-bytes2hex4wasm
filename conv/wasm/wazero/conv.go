package conv

import (
	"context"
	"errors"

	"github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"
)

var (
	ErrNilMem  error = errors.New("nil memory")
	ErrNilFunc error = errors.New("nil function")

	ErrUnableToRead  error = errors.New("unable to read the converted")
	ErrUnableToWrite error = errors.New("unable to write the original")
)

type WasmFn struct{ wa.Function }

func (f WasmFn) Call(ctx context.Context) error {
	_, err := f.Function.Call(ctx)
	return err
}

type WasmMem struct{ wa.Memory }

func (m WasmMem) ReadPage(offset uint32) ([]byte, bool) {
	return m.Memory.Read(offset, 65536)
}

func (m WasmMem) ReadConverted() ([]byte, error) {
	conv, ok := m.ReadPage(65536)
	if !ok {
		return nil, ErrUnableToRead
	}
	return conv, nil
}

func (m WasmMem) WriteBytes(offset uint32, data []byte) bool {
	return m.Memory.Write(offset, data)
}

type HalfPage [32768]byte

func (m WasmMem) WriteOriginal(original *HalfPage) error {
	ok := m.WriteBytes(0, original[:])
	if !ok {
		return ErrUnableToWrite
	}
	return nil
}

type WasmMod struct{ wa.Module }

func (m WasmMod) Close(ctx context.Context) error {
	return m.Module.Close(ctx)
}

func (m WasmMod) Memory() (WasmMem, error) {
	var mem wa.Memory = m.Module.Memory()
	if nil == mem {
		return WasmMem{}, ErrNilMem
	}
	return WasmMem{Memory: mem}, nil
}

func (m WasmMod) GetFunction(name string) (WasmFn, error) {
	var fnc wa.Function = m.Module.ExportedFunction(name)
	if nil == fnc {
		return WasmFn{}, ErrNilFunc
	}
	return WasmFn{Function: fnc}, nil
}

func (m WasmMod) GetConverter() (WasmFn, error) {
	return m.GetFunction("bytes2hex_hpage2page")
}

type Compiled struct{ wazero.CompiledModule }

func (c Compiled) Close(ctx context.Context) error {
	return c.CompiledModule.Close(ctx)
}

type WasmRuntime struct{ wazero.Runtime }

func (r WasmRuntime) Close(ctx context.Context) error {
	return r.Runtime.Close(ctx)
}

func (r WasmRuntime) Compile(
	ctx context.Context,
	wasm []byte,
) (Compiled, error) {
	cmod, err := r.Runtime.CompileModule(ctx, wasm)
	return Compiled{CompiledModule: cmod}, err
}

func (r WasmRuntime) Instantiate(
	ctx context.Context,
	compiled Compiled,
	cfg wazero.ModuleConfig,
) (WasmMod, error) {
	amod, err := r.Runtime.InstantiateModule(
		ctx,
		compiled.CompiledModule,
		cfg,
	)

	return WasmMod{Module: amod}, err
}

type WasmConfig struct{ wazero.RuntimeConfig }
