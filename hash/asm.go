//go:build ignore

package main

import (
	"go/types"
	"os"
	"slices"
	"strings"
	_ "unsafe" // FIXME: using this for linkname, delete after we fix local *build.Context

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

// FIXME
//
// //go:linkname avoCtx github.com/mmcloughlin/avo/build.ctx
var avoCtx *build.Context

// ===================================================================================================================

/*
# Avo notes

because they are strictly for `go generate` functionality, the avo examples are not exactly idiomatic golang.

 - when you write the generation functions you define a main package, but you add a build constraint to ignore it
   -\ this one we have to keep in order to use `go generate` without moving the file to a different directory

 - they use a global *build.Context as a primary way to pass around the state, along with an `init()`

 - the input to build.TEXT is just a macro that instantiates the global context with the function name and signature

// ===================================================================================================================

/*

# various asm notes for the gopher brained

despite it's name, MOV copies

(adr): *adr // dereference

ADD(src, dest): +=
SUB(src, dest): -=
MOV(src, dest): =
XOR(src, dest): ^=
OR(src, dest): |=
AND(src, dest): &=
NOT(src, dest): ^=
SHR(src, dest): >>=

JGE: Jump if result of CMPx greater or equal

movzx cx,ch  ; zero-extends ch into cx. the upper byte of cx will be filled with zeroes
movsx cx,ch  ; sign-extends ch into cx. the upper byte of cx will be filled with the most significant bit of ch

---

| cmp a,b | test a,b | Description      |
|---------|----------|------------------|
| je      | b == a   | b&a == 0         | Equal              |
| jne     | b != a   | b&a != 0         | Not equal          |
| js      | b-a < 0  | b&a < 0          | Sign (negative)    |
| jns     | b-a >= 0 | b&a >= 0         | Non-negative       |
| jg      | b > a    | b&a > 0          | Greater            |
| jge     | b >= a   | b&a >= 0         | Greater or equal   |
| jl      | b < a    | b&a < 0          | Less               |
| jle     | b <= a   | b&a <= 0         | Less or equal      |
| ja      | b > a    | b&a > 0U         | Above (unsigned >) |
| jb      | b < a    | b&a < 0U         | Below (unsigned <) |
| jz	  | b == 0   | b == 0           | Zero               |

	if (x < 3 && x == y) {
		return 1;
	} else {
		return 2;
	}

## ~=

    cmpq $3, %rdi  # Compare x with 3
    jge T2         # If x >= 3, jump to T2
    cmpq %rsi, %rdi # Compare x with y
    jne T2         # If x != y, jump to T2
	T1: # x < 3 && x == y:
    	movq $1, %rax  # Return 1
    	ret
	T2: # else
    	movq $2, %rax  # Return 2
    	ret

// --------------------------------------------------------

ADD and SUB carries any overflow by setting the carry flag (CF) in the EFLAGS register.
The carry flag indicates an overflow condition for unsigned-integer arithmetic.
The next instruction can test the carry flag with the JC (jump if carry) instruction.

ADC and SBB are similar to ADD and SUB, but they also add or subtract the value of the carry flag.

	- If the carry flag was set by ADD, a subequent ADD will add the value of the carry flag to the result.
	- If the carry flag was set by SUB, a subsequent SUB will subtract the value of the carry flag from the result.

In both the above cases, after the second operation, the carry flag is cleared;
assuming the operation did not set it again.

check what operations modify flags and what oprations modify or undefine them here:
	- ref: http://ref.x86asm.net/coder32.html
*/

// ===================================================================================================================

type checksumASM struct {
	// pass through methods from underlying global context
	// they don't export it so we linkname (cringe i know)
	// TODO: probably just create a new context lol
	ctx *build.Context

	name       string
	inputName  string
	outputName string
	doc        string
	args       *types.Signature

	data operand.Mem

	// note that because we're in go generate,
	// we don't need to worry about synchronization with regard to map access.
	registers      map[string]reg.Register
	sizedRegisters map[int]map[string]reg.Register
}

func (f *checksumASM) prepRegisters() {
	f.ctx.Comment("initialize registers")
	// ========= 64 bit registers =========
	// rdx: 64 bit register
	f.registers["rdx"] = build.GP64()
	f.registers["sum"] = f.registers["rdx"]
	f.ctx.XORQ(f.registers["rdx"], f.registers["rdx"])
	// rdi: 64 bit register representing our input data via a pointer
	//  - using ADDQ(operand.U8(2), reg.RDI)
	//    to increment the pointer to next byte pair
	//    during loop iteration
	f.registers["rdi"] = build.GP64()
	f.ctx.XORQ(f.registers["rdi"], f.registers["rdi"])
	// index register
	f.registers["i"] = build.GP64()
	f.ctx.XORQ(f.registers["i"], f.registers["i"])
	// r8: 64 bit register
	f.registers["r8"] = build.GP64()
	f.ctx.XORQ(f.registers["r8"], f.registers["r8"])
	// r9: 64 bit register
	f.registers["r9"] = build.GP64()
	f.ctx.XORQ(f.registers["r9"], f.registers["r9"])
	// rsi: 64 bit register for size of input data
	f.registers["rsi"] = build.GP64()
	f.ctx.XORQ(f.registers["rsi"], f.registers["rsi"])
	// =====================================

	// ========= 32 bit registers =========
	// edx: 32 bit register is the lower 32 bits of rdx
	// we'll accumulate the sum in this register
	f.registers["edx"] = f.registers["rdx"].(reg.GPVirtual).As32()
	// eax: 32 bit register for overflow
	f.registers["eax"] = build.GP32()
	f.ctx.XORL(f.registers["eax"], f.registers["eax"])
	// =====================================

	// ========= 16 bit registers =========
	// r8w: 16 bit register representing the low 16 bits of r8
	// 	- storage for high bytes of data processed
	//    - mov a byte in, shift left by 8 bits to acquire the high byte
	//  - combine it with r9w to form a 16 bit word
	f.registers["r8w"] = f.registers["r8"].(reg.GPVirtual).As16()
	// r9w: storage for low bytes of data processed to be combined with r8w
	f.registers["r9w"] = f.registers["r9"].(reg.GPVirtual).As16()
	// dx: lower 16 bits of edx
	f.registers["dx"] = f.registers["edx"].(reg.GPVirtual).As16()
	// ax: lower 16 bits of eax
	f.registers["ax"] = f.registers["eax"].(reg.GPVirtual).As16()
	// =====================================

	// ========= 8 bit registers =========
	// r8b: 8 bit register representing the low 8 bits of r8b
	//     - storage for low bytes of data processed
	f.registers["r8b"] = f.registers["r8w"].(reg.GPVirtual).As8()
	// r9b: storage for high bytes of data processed to be combined with r8b
	f.registers["r9b"] = f.registers["r9w"].(reg.GPVirtual).As8()
	// =====================================

	for key, value := range f.registers {
		if f.sizedRegisters[int(value.Size())*8] == nil {
			f.sizedRegisters[int(value.Size())*8] = make(map[string]reg.Register)
		}
		f.sizedRegisters[int(value.Size())*8][key] = value
	}

	// f.ctx.XORQ(f.registers["rdx"], f.registers["rdx"])
	// f.ctx.XORQ(f.registers["rsi"], f.registers["rsi"])

}

func (f *checksumASM) handle16BitRDXOverflow(register reg.Register) {
	swap := f.ctx.GP64()
	f.ctx.XORQ(swap, swap)
	// save register contents before we shift, effectively storing the overflow in another register
	f.ctx.MOVQ(register.(reg.GPVirtual).As64(), swap)
	// mask lower 16 of target
	f.ctx.ANDQ(operand.U32(0xFFFF), register.(reg.GPVirtual).As64())
	// shift work register right 16 bits so that the overflow is in the lower 16 bits
	f.ctx.SHRQ(operand.U8(16), swap)
	// snag our overflow from EAX and add it to EDX
	f.ctx.ADDQ(swap, register.(reg.GPVirtual).As64())

	// f.ctx.MOVL(f.registers["edx"], f.registers["eax"])
}

func (f *checksumASM) loadInput() {
	inputData := f.ctx.Param(f.inputName)
	// pointer to base of uint8 array ([]byte) input into new 64 bit register
	f.data = operand.Mem{
		Base: f.ctx.Load(
			inputData.Base(),
			f.sizedRegisters[64]["rdi"].(reg.GPVirtual),
		),
	}
	// length of input data
	dataLen := inputData.Len()
	// move it to dataLen register (rsi, we use an alias here, aka "remaining")
	f.ctx.Load(dataLen, f.sizedRegisters[64]["rsi"])
}

func newChecksumASM(name, inputName, outputName, doc string) *checksumASM {
	asmf := &checksumASM{}

	// asmf.ctx = build.NewContext()
	asmf.ctx = avoCtx
	asmf.name = name
	asmf.inputName = inputName
	asmf.outputName = outputName
	asmf.doc = doc
	asmf.registers = make(map[string]reg.Register)
	asmf.sizedRegisters = make(map[int]map[string]reg.Register)

	// equivelant to build.TEXT(name, build.NOSPLIT, "func(data []byte) uint16")
	// but local context instead of global
	asmf.ctx.Function(name)
	asmf.ctx.Attributes(build.NOSPLIT)
	asmf.ctx.SignatureExpr("func(" + inputName + " []byte) (" + outputName + " uint16)")
	asmf.ctx.Doc(doc)

	// we're using 64 bit registers so we're constrained to 64 bit architectures
	// ...and i don't even know how other instructions sets work lmao so amd64 it is
	// TODO: learn how to sign up for email
	asmf.ctx.ConstraintExpr("amd64")

	asmf.prepRegisters()
	asmf.loadInput()

	return asmf
}

func (f *checksumASM) AddLabeledFunc(name string, fnc func()) operand.LabelRef {
	f.ctx.Label(name)
	fnc()
	return operand.LabelRef(name)
}

func (f *checksumASM) nextb() {
	f.ctx.ADDQ(operand.Imm(2), f.data.Base)
	f.ctx.SUBQ(operand.Imm(2), f.sizedRegisters[64]["rsi"])
	f.ctx.JNC(operand.LabelRef("loop"))
	f.ctx.JMP(operand.LabelRef("fin"))
}

func (f *checksumASM) loop() {
	r8w := f.sizedRegisters[16]["r8w"]
	r9w := f.sizedRegisters[16]["r9w"]

	sum := f.sizedRegisters[64]["sum"]

	f.ctx.TESTQ(f.sizedRegisters[64]["rsi"], f.sizedRegisters[64]["rsi"])
	f.ctx.JZ(operand.LabelRef("fin"))

	f.ctx.CMPQ(f.sizedRegisters[64]["rsi"], operand.Imm(2))
	f.ctx.JL(operand.LabelRef("handle_odd"))

	// load first byte into r8w and fill the rest of r8w with zeros
	f.ctx.MOVB(f.data.Offset(0), r8w.(reg.GPVirtual).As8())
	// shift left by 8 bits to make room for the next byte
	f.ctx.SHLW(operand.Imm(8), r8w)
	// increment data pointer
	// f.ctx.INCQ(f.data.Base)
	// load second byte into r9w
	f.ctx.MOVB(f.data.Offset(1), r9w.(reg.GPVirtual).As8())
	// combine r8w and r9w to form a 16 bit word
	f.ctx.ORW(r9w, r8w)
	// add 16 bit word to 32 bit sum
	f.ctx.ADDQ(r8w.(reg.GPVirtual).As64(), sum)

	f.ctx.CMPL(f.sizedRegisters[64]["sum"].(reg.GPVirtual).As32(), operand.U32(0xFFFF))
	f.ctx.JA(operand.LabelRef("adjust_sum"))
}

func main_t(f *checksumASM, mode string) {
	switch {
	case strings.EqualFold("early_fail", mode):
		main_tEarlyFail(f)
	case strings.EqualFold("handle_odd", mode):
		f.AddLabeledFunc("adjust_sum", func() {
			f.handle16BitRDXOverflow(f.sizedRegisters[32]["edx"])
		})
		main_tCheckOdd(f)
	default:
		panic("unknown test mode: '" + mode + "'!")
	}
}

func main_tEarlyFail(f *checksumASM) {
	f.earlyCheck()

	retReg := build.GP16()
	f.ctx.XORW(retReg, retReg)
	f.ctx.MOVW(operand.U16(5), retReg)
	f.ctx.Store(retReg, f.ctx.Return(f.outputName))
	f.ctx.RET()
}

func main_tCheckOdd(f *checksumASM) {
	f.earlyCheck()

	f.ctx.MOVB(f.data, f.sizedRegisters[8]["r8b"])

	f.handle16BitRDXOverflow(f.sizedRegisters[64]["r8"])

	f.ctx.Label("fin")
	f.ctx.Store(f.registers["ax"], f.ctx.Return(f.outputName))
	f.ctx.RET()
}

func (f *checksumASM) earlyCheck() {
	f.ctx.TESTQ(f.registers["rsi"], f.registers["rsi"])
	f.ctx.JZ(operand.LabelRef("early_fail"))
}

func (f *checksumASM) earlyFail() {
	retReg := build.GP16()
	f.ctx.XORW(retReg, retReg)
	f.ctx.MOVW(operand.U16(0), retReg)
	f.ctx.Store(retReg, f.ctx.Return(f.outputName))
	f.ctx.RET()
}

func (f *checksumASM) handleOdd() {
	f.ctx.CMPQ(f.sizedRegisters[64]["rsi"].(reg.GPVirtual).As64(), operand.Imm(1))
	f.ctx.JNE(operand.LabelRef("fin"))

	f.ctx.MOVB(f.data.Offset(0), f.sizedRegisters[8]["r8b"])
	f.ctx.SHLW(operand.Imm(8), f.sizedRegisters[8]["r8b"].(reg.GPVirtual).As16())
	f.ctx.ADDQ(f.sizedRegisters[8]["r8b"].(reg.GPVirtual).As64(), f.sizedRegisters[64]["sum"].(reg.GPVirtual).As64())
	f.ctx.CMPQ(f.sizedRegisters[64]["sum"].(reg.GPVirtual).As64(), operand.U32(0xFFFF))
	f.ctx.JA(operand.LabelRef("adjust_sum"))
}

func main() {
	f := newChecksumASM("checksum", "data", "sum", "calculate RFC 1071 internet checksum for a byte slice")

	if os.Getenv("ASM_TEST_MODE") != "" || slices.Contains(os.Args, "-asmtest") {
		os.Args = os.Args[:len(os.Args)-1]
		println("running test mode: " + os.Getenv("ASM_TEST_MODE"))
		main_t(f, os.Getenv("ASM_TEST_MODE"))
		f.AddLabeledFunc("early_fail", f.earlyFail)
		goto gen
	}

	f.AddLabeledFunc("early_check", f.earlyCheck)

	f.AddLabeledFunc("loop", f.loop)

	f.AddLabeledFunc("nextb", f.nextb) // jumps to loop if we're not done

	f.AddLabeledFunc("fin", func() {
		f.ctx.CMPQ(f.sizedRegisters[64]["sum"].(reg.GPVirtual).As64(), operand.U32(0xFFFF))
		f.ctx.JA(operand.LabelRef("adjust_sum"))
		f.ctx.NOTW(f.sizedRegisters[64]["sum"].(reg.GPVirtual).As16())
		f.ctx.Store(f.sizedRegisters[64]["sum"].(reg.GPVirtual).As16(), f.ctx.Return(f.outputName))
		f.ctx.RET()
	})

	f.AddLabeledFunc("early_fail", f.earlyFail)

	f.AddLabeledFunc("handle_odd", f.handleOdd)

	// handle overflow if we didn't jump
	f.AddLabeledFunc("adjust_sum", func() {
		f.handle16BitRDXOverflow(f.sizedRegisters[64]["sum"])
		f.ctx.JMP(operand.LabelRef("nextb"))
	})

gen:

	build.Generate()
}
