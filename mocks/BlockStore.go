// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	types "github.com/lazyledger/optimint/types"
)

// BlockStore is an autogenerated mock type for the BlockStore type
type BlockStore struct {
	mock.Mock
}

// Height provides a mock function with given fields:
func (_m *BlockStore) Height() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// LoadBlock provides a mock function with given fields: height
func (_m *BlockStore) LoadBlock(height uint64) *types.Block {
	ret := _m.Called(height)

	var r0 *types.Block
	if rf, ok := ret.Get(0).(func(uint64) *types.Block); ok {
		r0 = rf(height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Block)
		}
	}

	return r0
}

// LoadBlockByHash provides a mock function with given fields: hash
func (_m *BlockStore) LoadBlockByHash(hash [32]byte) *types.Block {
	ret := _m.Called(hash)

	var r0 *types.Block
	if rf, ok := ret.Get(0).(func([32]byte) *types.Block); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Block)
		}
	}

	return r0
}

// SaveBlock provides a mock function with given fields: block
func (_m *BlockStore) SaveBlock(block *types.Block) {
	_m.Called(block)
}
