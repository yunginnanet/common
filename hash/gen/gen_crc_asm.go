// //go:build ignore

package main

import (
	"fmt"
	"go/types"
	_ "unsafe" // FIXME: using this for linkname, delete after we fix local *build.Context

	"github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/pass"
	"github.com/mmcloughlin/avo/reg"
)

// ===================================================================================================================

/*
# Avo notes

because they are strictly for `go generate` functionality, the avo examples are not exactly idiomatic golang.

 - when you write the generation functions you define a main package, but you add a build constraint to ignore it
   -\ this one we have to keep in order to use `go generate` without moving the file to a different directory

 - they use a global *build.Context as a primary way to pass around the state, along with an `init()`

 - the input to build.TEXT is just a macro that instantiates the global context with the function name and signature

we don't have to use that global context though, so we don't. much more idiomatic to declare a local context

### other notes

the following code is definitely longer than it needs to be, lots of syntactic sugar and abstraction.
it just works better in my head to define structs and methods on them to represent the state of various things.
this is inconsequential; as the generaeted goasm will be the same as if it was written as the examples in the Avo repo.
*/

// ===================================================================================================================

/*

# various asm notes for the gopher brained

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

*/

// ===================================================================================================================

// WordSize represents the bit size of x86 components
type WordSize uint

const (
	BWord WordSize = 8
	Word  WordSize = 16
	DWord WordSize = 32
	QWord WordSize = 64
)

func (ws WordSize) U() uint {
	return uint(ws / 8)
}

func (ws WordSize) String() string {
	switch ws {
	case BWord:
		return "byte"
	case Word:
		return "word"
	case DWord:
		return "dword"
	case QWord:
		return "qword"
	default:
		return "unknown"
	}
}

type asmFunc interface {
	Ctx() *build.Context
	GetRegister(name string) reg.Register
	NewRegister(name string, size WordSize) reg.Register
	AddLabeledFunc(name string, op func())
	GetLabel(string) operand.LabelRef
	ZeroByName(names ...string)
	Zero(regs ...reg.Register) error
}

var _ asmFunc = &checksumASM{}

type checksumASM struct {
	// pass through methods from underlying global context
	// they don't export it so we linkname (cringe i know)
	// TODO: probably just create a new context lol
	ctx *build.Context

	name string
	doc  string
	args *types.Signature

	data            reg.Register
	dataLen         reg.Register
	dataLenRegister reg.GPVirtual

	// note that because we're in go generate,
	// we don't need to worry about synchronization with regard to map access.
	registers map[string]reg.Register
}

func (f *checksumASM) Ctx() *build.Context {
	return f.ctx
}

func (f *checksumASM) GetRegister(name string) reg.Register {
	if r, ok := f.registers[name]; ok {
		return r
	}
	panic(fmt.Errorf("could not find register named '%s'", name))
	return nil // unreachable
}

func (f *checksumASM) NewRegister(name string, size WordSize) reg.Register {
	var r reg.GPVirtual
	switch size {
	case BWord:
		r = f.ctx.GP8()
	case Word:
		r = f.ctx.GP16()
	case DWord:
		r = f.ctx.GP32()
	case QWord:
		r = f.ctx.GP64()
	default:
		panic("unsupported word size")
	}
	f.registers[name] = r
	return r
}

func (f *checksumASM) Zero(reg ...reg.Register) error {
	for _, rg := range reg {
		switch rg.Size() {
		case BWord.U(): // 8 bits, 1 byte
			f.ctx.XORB(rg, rg)
		case Word.U(): // 16 bits, 2 bytes
			f.ctx.XORW(rg, rg)
		case DWord.U(): // 32 bits, 4 bytes
			f.ctx.XORL(rg, rg)
		case QWord.U(): // 64 bits, 8 bytes
			f.ctx.XORQ(rg, rg)
		default:
			return fmt.Errorf("unsupported word size: %d", rg.Size())
		}
	}
	return nil
}

func (f *checksumASM) ZeroByName(names ...string) {
	for _, name := range names {
		rg := f.GetRegister(name)
		if err := f.Zero(); err != nil {
			panic(fmt.Errorf("could not zero register '%s' (%d): %w", name, rg.Size(), err))
		}

	}
}

func (f *checksumASM) GetLabel(name string) operand.LabelRef {
	return operand.LabelRef(name) // syntactic sugar
}

func (f *checksumASM) JumpToLabel(name string) {
	f.ctx.JMP(f.GetLabel(name))
}

func (f *checksumASM) AddLabeledFunc(name string, fnc func()) {
	f.ctx.Label(name)
	fnc()
}

func (f *checksumASM) jumpIfGreaterOrEqual(r1, r2 reg.Register, label string) {
	f.ctx.CMPQ(r1, r2)
	f.ctx.JGE(f.GetLabel(label))
}

func (f *checksumASM) prepRegisters() {
	f.NewRegister("sum", QWord)                 // sum
	f.registers["high"] = reg.Register(reg.R8)  // high 16 bits // r8
	f.registers["index"] = reg.Register(reg.R9) // loop index // r9

	// auxiliary registers
	// note: aux64 initially holds the length of the input
	f.registers["aux64"] = reg.Register(reg.R10) // r10
	f.registers["aux32"] = reg.Register(reg.R11) // r11

	f.ZeroByName("sum", "high", "index", "aux64", "aux32")
}

func (f *checksumASM) loadInput() {
	// load pointer to uint8 array ([]byte) input into new 64 bit register
	f.data = f.ctx.Load(f.ctx.Param("data").Base(), f.ctx.GP64())
	// length of uint8 array ([]byte) input into new 64 bit register
	f.dataLen = f.ctx.Load(f.ctx.Param("data").Len(), f.ctx.GP64())
	f.dataLenRegister = f.ctx.GP64()

	f.ctx.DECQ(f.dataLen)
	f.ctx.MOVQ(f.dataLen, f.GetRegister("aux64"))
}

func (f *checksumASM) registerInteropVars() (sum, r8, r9, r10, r11 reg.Register) {
	sum = f.GetRegister("sum")
	r8 = f.GetRegister("high")
	r9 = f.GetRegister("index")
	r10 = f.GetRegister("aux64")
	r11 = f.GetRegister("aux32")
	return
}

// temporary
//
//go:linkname avoCtx github.com/mmcloughlin/avo/build.ctx
var avoCtx *build.Context

func newChecksumASM(name, doc string) *checksumASM {
	asmf := &checksumASM{}

	// asmf.ctx = build.NewContext()
	asmf.ctx = avoCtx
	asmf.name = name
	asmf.doc = doc
	asmf.registers = make(map[string]reg.Register, 5)

	// equivelant to build.TEXT(name, build.NOSPLIT, "func(data []byte) uint16")
	// but local context instead of global
	asmf.ctx.Function(name)
	asmf.ctx.Attributes(build.NOSPLIT)
	asmf.ctx.SignatureExpr("func(data []byte) uint16")
	asmf.ctx.Doc(doc)

	// we're using 64 bit registers so we're constrained to 64 bit architectures
	// ...and i don't even know how other instructions sets work lmao so amd64 it is
	// TODO: learn how to sign up for email
	asmf.ctx.ConstraintExpr("amd64")

	asmf.prepRegisters()
	asmf.loadInput()

	asmf.ctx.DECQ(asmf.GetRegister("aux64"))

	return asmf
}

func (f *checksumASM) loop() {
	sum, r8, r9, r10, r11 := f.registerInteropVars()
	f.jumpIfGreaterOrEqual(r9.(reg.GPPhysical), r10.(reg.GPPhysical), "check_odd")

	/*
		MOVBLZX notes:
			MOVBLZX(src, dest)
			MOVe Byte to Lower, Zero eXtend (8 bits to 64 bits)
	*/

	/*
		MOVBQZX notes:
			MOVBQZX(src, dest)
			MOVe Byte to Quadword, Zero eXtend (8 bits to 64 bits)

			zero extend fills the upper bits with 0s
			so move 8 bits into a 64 bit register, the upper 56 bits get zeroed

	*/

	/*
		var sum [8]byte // uint64
		var fdata []byte
		var r9 uint64 // [8]byte // quadword
		var r11 [4]byte // doubleword
	*/

	/*
		r11[0] = fdata[r9]
		r11[1:] = [7]byte{0, 0, 0, 0, 0, 0, 0}
	*/
	f.ctx.MOVBQZX(operand.Mem{Base: f.data, Index: r9, Scale: 1}, r11)

	//	r9++
	f.ctx.INCQ(r9)

	/*
		sum[0] = fdata[r9]
		sum[1:] = [7]byte{0, 0, 0, 0, 0, 0, 0}
	*/
	f.ctx.MOVBQZX(operand.Mem{Base: f.data, Index: r9, Scale: 1}, r8)

	/*
		SHLQ notes:
			SHLQ(count, dest)
			SHift Left a Quadword

			shift the bits in the destination left by the number of bits specified in the count operand
			and store the result in the destination operand. you probably wanna deref dest
	*/

	// sum <<= bitlength(byte)
	f.ctx.SHLQ(operand.U8(8), r8) // $8, %r8 // uint8(8), *r8

	// sum |= r11
	f.ctx.ORQ(r11, r8)

	// sum += r8
	f.ctx.ADDQ(r8, sum)

	// r9++
	f.ctx.INCQ(r9)

	f.ctx.JMP(f.GetLabel("loop"))
}

func (f *checksumASM) checkOdd() {
	_, _, _, r10, _ := f.registerInteropVars()

	/*
		TESTQ notes:
			TESTQ(src, dest)
			TEST Quadword

			Performs a bitwise AND of the two operands (first operand is the source operand)
			OF and CF are cleared; SF, ZF, and PF are set according to the result.
		    If the result is zero, the ZF flag is set; otherwise, it is cleared.
	*/

	// res := r10 & r10
	f.ctx.TESTQ(r10, r10)
	// 		if res == 0 { jump to fold_sum }
	f.ctx.JZ(f.GetLabel("fold_sum"))
	// }
}

func (f *checksumASM) addSumOp() {
	// only happens if we have an odd byte after the loop

	sum, r8, _, r10, _ := f.registerInteropVars()

	// sum += f.data[r10]
	f.ctx.MOVBQZX(operand.Mem{Base: f.data, Index: r10, Scale: 1}, r8)
	f.ctx.ADDQ(r8, sum)
}

func (f *checksumASM) foldSum() {
	sum, r8, _, _, _ := f.registerInteropVars()

	f.ctx.MOVQ(sum, r8)
	// shift right quadword by 16 bits
	f.ctx.SHRQ(operand.U8(16), r8)
	f.ctx.ADDQ(r8, sum)
	f.ctx.MOVQ(sum, r8)
	// r8 = r8 & 0xffff
	f.ctx.ANDQ(operand.U32(0xffff), r8)
	f.ctx.ADDQ(r8, sum)

	/*
		NOTQ notes:
			NOTQ(dest)
			NOT Quadword

			Performs a bitwise NOT bitwise operation (each 1 is changed to 0 and each 0 is changed to 1)
			on the destination operand and stores the result in the destination operand.
	*/

	// sum ^= sum
	f.ctx.NOTQ(sum)
	// return sum
	f.ctx.Store(sum.(reg.GPVirtual).As16(), f.ctx.ReturnIndex(0))
}

func (f *checksumASM) _return() {
	f.ctx.RET()
}

func (f *checksumASM) Generate() {
	res, err := f.ctx.Result()
	if err != nil {
		panic(err)
	}

	p := pass.Concat([]pass.Interface{pass.Compile}...)
	if err = p.Execute(res); err != nil {
		panic(err)
	}
}

func main() {
	// registers are initialized/zeroed in the constructor when it calls asm.prepRegisters()
	f := newChecksumASM("checksum", "calculate RFC 1071 internet checksum for a byte slice")

	f.AddLabeledFunc("loop", f.loop)

	/*
		res := r10 & r10
			if res == 0 { jump to fold_sum } // not odd
		}
	*/
	f.AddLabeledFunc("check_odd", f.checkOdd)

	f.addSumOp()

	// 	fold_sum:
	f.AddLabeledFunc("fold_sum", f.foldSum)

	f._return()

	// f.Generate()
	// FIXME
	build.Generate()
}
