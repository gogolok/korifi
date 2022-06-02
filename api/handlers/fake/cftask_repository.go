// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/handlers"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type CFTaskRepository struct {
	CreateTaskStub        func(context.Context, authorization.Info, repositories.CreateTaskMessage) (repositories.TaskRecord, error)
	createTaskMutex       sync.RWMutex
	createTaskArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateTaskMessage
	}
	createTaskReturns struct {
		result1 repositories.TaskRecord
		result2 error
	}
	createTaskReturnsOnCall map[int]struct {
		result1 repositories.TaskRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFTaskRepository) CreateTask(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateTaskMessage) (repositories.TaskRecord, error) {
	fake.createTaskMutex.Lock()
	ret, specificReturn := fake.createTaskReturnsOnCall[len(fake.createTaskArgsForCall)]
	fake.createTaskArgsForCall = append(fake.createTaskArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateTaskMessage
	}{arg1, arg2, arg3})
	stub := fake.CreateTaskStub
	fakeReturns := fake.createTaskReturns
	fake.recordInvocation("CreateTask", []interface{}{arg1, arg2, arg3})
	fake.createTaskMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFTaskRepository) CreateTaskCallCount() int {
	fake.createTaskMutex.RLock()
	defer fake.createTaskMutex.RUnlock()
	return len(fake.createTaskArgsForCall)
}

func (fake *CFTaskRepository) CreateTaskCalls(stub func(context.Context, authorization.Info, repositories.CreateTaskMessage) (repositories.TaskRecord, error)) {
	fake.createTaskMutex.Lock()
	defer fake.createTaskMutex.Unlock()
	fake.CreateTaskStub = stub
}

func (fake *CFTaskRepository) CreateTaskArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateTaskMessage) {
	fake.createTaskMutex.RLock()
	defer fake.createTaskMutex.RUnlock()
	argsForCall := fake.createTaskArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFTaskRepository) CreateTaskReturns(result1 repositories.TaskRecord, result2 error) {
	fake.createTaskMutex.Lock()
	defer fake.createTaskMutex.Unlock()
	fake.CreateTaskStub = nil
	fake.createTaskReturns = struct {
		result1 repositories.TaskRecord
		result2 error
	}{result1, result2}
}

func (fake *CFTaskRepository) CreateTaskReturnsOnCall(i int, result1 repositories.TaskRecord, result2 error) {
	fake.createTaskMutex.Lock()
	defer fake.createTaskMutex.Unlock()
	fake.CreateTaskStub = nil
	if fake.createTaskReturnsOnCall == nil {
		fake.createTaskReturnsOnCall = make(map[int]struct {
			result1 repositories.TaskRecord
			result2 error
		})
	}
	fake.createTaskReturnsOnCall[i] = struct {
		result1 repositories.TaskRecord
		result2 error
	}{result1, result2}
}

func (fake *CFTaskRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createTaskMutex.RLock()
	defer fake.createTaskMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFTaskRepository) recordInvocation(key string, args []interface{}) {
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

var _ handlers.CFTaskRepository = new(CFTaskRepository)