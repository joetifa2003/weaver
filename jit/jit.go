package jit

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joetifa2003/weaver/parser"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

func Run(source string) error {
	// 1. Parse the weaver source
	p, err := parser.Parse(source)
	if err != nil {
		return fmt.Errorf("failed to parse weaver source: %w", err)
	}

	// 2. Generate Go source code
	goSrc, err := Generate(&p)
	if err != nil {
		return fmt.Errorf("failed to generate go source: %w", err)
	}

	// 3. Save the generated Go source to a temporary file
	tmpDir, err := os.MkdirTemp("", "weaver-jit")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	goFilePath := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(goFilePath, []byte(goSrc), 0644)
	if err != nil {
		return fmt.Errorf("failed to write go source file: %w", err)
	}

	// 4. Compile the Go code to Wasm
	wasmFilePath := filepath.Join(tmpDir, "main.wasm")
	cmd := exec.Command("tinygo", "build", "-target=wasi", "-o", wasmFilePath, goFilePath)
	cmd.Env = append(os.Environ(), "GOROOT=/usr/local/go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compile go source to wasm: %w\n%s", err, output)
	}

	// 5. Read the compiled Wasm file
	wasmBytes, err := os.ReadFile(wasmFilePath)
	if err != nil {
		return fmt.Errorf("failed to read wasm file: %w", err)
	}

	ctx := context.Background()

	// 6. Create a new wazero runtime
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// 7. Instantiate the Wasm module
	// We also need to configure the module to use the same stdin/stdout as the host.
	config := wazero.NewModuleConfig().WithStdout(os.Stdout).WithStderr(os.Stderr)
	mod, err := r.InstantiateWithConfig(ctx, wasmBytes, config)
	if err != nil {
		return fmt.Errorf("failed to instantiate wasm module: %w", err)
	}

	// 8. Run the main function of the Wasm module.
	// The Go compiler for wasm exports a `_start` function.
	start := mod.ExportedFunction("_start")
	if start == nil {
		return fmt.Errorf("_start function not found in wasm module")
	}
	_, err = start.Call(ctx)
	if err != nil {
		var exitErr *sys.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 0 {
			return nil
		}
		return fmt.Errorf("failed to run wasm module: %w", err)
	}

	return nil
}
