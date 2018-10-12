package mock

// copy from tendermint/rpc/core/dev.go

import (
	"errors"
	"os"
	"runtime/pprof"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UnsafeFlushMempool unsafe
func UnsafeFlushMempool() (*ctypes.ResultUnsafeFlushMempool, error) {
	// mempool.Flush()
	// return &ctypes.ResultUnsafeFlushMempool{}, nil
	return nil, errors.New("not implemented yet")
}

var profFile *os.File

// UnsafeStartCPUProfiler unsafe
func UnsafeStartCPUProfiler(filename string) (*ctypes.ResultUnsafeProfile, error) {
	var err error
	profFile, err = os.Create(filename)
	if err != nil {
		return nil, err
	}
	err = pprof.StartCPUProfile(profFile)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultUnsafeProfile{}, nil
}

// UnsafeStopCPUProfiler unsafe
func UnsafeStopCPUProfiler() (*ctypes.ResultUnsafeProfile, error) {
	pprof.StopCPUProfile()
	if err := profFile.Close(); err != nil {
		return nil, err
	}
	return &ctypes.ResultUnsafeProfile{}, nil
}

// UnsafeWriteHeapProfile unsafe
func UnsafeWriteHeapProfile(filename string) (*ctypes.ResultUnsafeProfile, error) {
	memProfFile, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	if err := pprof.WriteHeapProfile(memProfFile); err != nil {
		return nil, err
	}
	if err := memProfFile.Close(); err != nil {
		return nil, err
	}

	return &ctypes.ResultUnsafeProfile{}, nil
}
