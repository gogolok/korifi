// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"sync"

	"code.cloudfoundry.org/korifi/controllers/controllers/workloads"
)

type SeqIdGenerator struct {
	GenerateStub        func() (int64, error)
	generateMutex       sync.RWMutex
	generateArgsForCall []struct {
	}
	generateReturns struct {
		result1 int64
		result2 error
	}
	generateReturnsOnCall map[int]struct {
		result1 int64
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SeqIdGenerator) Generate() (int64, error) {
	fake.generateMutex.Lock()
	ret, specificReturn := fake.generateReturnsOnCall[len(fake.generateArgsForCall)]
	fake.generateArgsForCall = append(fake.generateArgsForCall, struct {
	}{})
	stub := fake.GenerateStub
	fakeReturns := fake.generateReturns
	fake.recordInvocation("Generate", []interface{}{})
	fake.generateMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *SeqIdGenerator) GenerateCallCount() int {
	fake.generateMutex.RLock()
	defer fake.generateMutex.RUnlock()
	return len(fake.generateArgsForCall)
}

func (fake *SeqIdGenerator) GenerateCalls(stub func() (int64, error)) {
	fake.generateMutex.Lock()
	defer fake.generateMutex.Unlock()
	fake.GenerateStub = stub
}

func (fake *SeqIdGenerator) GenerateReturns(result1 int64, result2 error) {
	fake.generateMutex.Lock()
	defer fake.generateMutex.Unlock()
	fake.GenerateStub = nil
	fake.generateReturns = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *SeqIdGenerator) GenerateReturnsOnCall(i int, result1 int64, result2 error) {
	fake.generateMutex.Lock()
	defer fake.generateMutex.Unlock()
	fake.GenerateStub = nil
	if fake.generateReturnsOnCall == nil {
		fake.generateReturnsOnCall = make(map[int]struct {
			result1 int64
			result2 error
		})
	}
	fake.generateReturnsOnCall[i] = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *SeqIdGenerator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.generateMutex.RLock()
	defer fake.generateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SeqIdGenerator) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ workloads.SeqIdGenerator = new(SeqIdGenerator)