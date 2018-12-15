package vm

import (
	"errors"

	"github.com/marmotini/monkey-lang/code"
	"github.com/marmotini/monkey-lang/compiler"
	"github.com/marmotini/monkey-lang/object"
)

const StackSize = 2048

var ErrStackOverflow = errors.New("stack overflow")

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int
}

func NewVM(bytecode *compiler.Bytecode) *VM {
	return &VM{
		stack:        make([]object.Object, StackSize),
		sp:           0,
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

// Run contains the fetch-decode-execute cycle
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.OpCode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()

			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			vm.push(&object.Integer{Value: leftValue + rightValue})
		}
	}
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return ErrStackOverflow
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}