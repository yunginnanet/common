//go:build amd64

package hash

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

const (
	goGen         = "go run asm.go -out checksum_amd64.s -stubs checksum_amd64.go"
	testEarlyFail = "early_fail"
	green         = "\033[32m"
	red           = "\033[31m"
	reset         = "\033[0m"
	testPassed    = green + "\n\ntest passed: " + reset + "%v\n"
	testFailed    = red + "\n\ntest failed: " + reset + "%v\n"
)

func readPipe(ctx context.Context, pipePath string, u16bChan chan []byte, t *testing.T) {
	tryRead := func() (goOn bool) {
		f, err := os.OpenFile(pipePath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			t.Errorf("%sfailed to open named pipe: %s%s", red, err, reset)
			return true
		}
		defer func() {
			if err = f.Close(); err != nil {
				panic(err)
			}
		}()
		b := make([]byte, 2)
		if _, err = f.Read(b); err != nil {
			t.Errorf("%sfailed to read from named pipe: %s%s", red, err, reset)
			return true
		}
		u16bChan <- b
		return false
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !tryRead() {
				return
			}
		}
	}
}

func TestASMChecksumComponents(t *testing.T) {
	t.Cleanup(func() {
		t.Setenv("ASM_TEST_MODE", "")
		cmdFields := strings.Fields(goGen)
		cmd := exec.Command(cmdFields[0], cmdFields[1:]...)
		cmd.Env = append(cmd.Environ(), `ASM_TEST_MODE=`)
		t.Logf("cleaning up with: %s", cmd.String())
		if cmdOut, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("%scleanup failed: %s%s", red, cmdOut, reset)
		}
	})

	type asmTest struct {
		name   string
		mode   string
		input  string
		expect uint16
	}

	testModes := []asmTest{
		{
			name:   "zero length input",
			mode:   "early_fail",
			input:  "[]byte{}",
			expect: 0,
		},
		{
			name:   "single byte input",
			mode:   "early_fail",
			input:  `[]byte{0x5}`,
			expect: 5,
		},
		/*		{
				name:   "three byte input",
				mode:   "handle_odd",
				input:  `[]byte{0x5, 0x7, 0x9}`,
				expect: 21,
			},*/
	}

	for _, mode := range testModes {
		t.Run(mode.mode+"/"+mode.name, func(t *testing.T) {
			t.Setenv("ASM_TEST_MODE", mode.mode)
			cmdFields := append(strings.Fields(goGen), "-asmtest")
			cmd := exec.Command(cmdFields[0], cmdFields[1:]...)
			t.Logf("generating test ASM with: %s", cmd.String())
			cmdOut, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("%sgeneration failed:\n %s%s", red, cmdOut, reset)
			}
			asmOut, err := os.ReadFile("checksum_amd64.s")
			if err != nil {
				t.Fatalf("failed to read generated ASM: %s", err)
			}
			t.Logf("generated test ASM: \n%s", asmOut)

			tmpDir := t.TempDir()
			pipePath := filepath.Join(tmpDir, "checksum")
			if err = syscall.Mkfifo(pipePath, 0777); err != nil {
				t.Fatalf("failed to create named pipe: %s", err)
			}

			test := `
package hash
import "testing"
import "os"
import "fmt"
import "strings"
func TestChecksumTmp(t *testing.T) {
	data := ` + mode.input + `
	expected := uint16(` + strconv.Itoa(int(mode.expect)) + `)
	actual := checksum(data)
	if actual != expected {
		t.Errorf("` + red + `Expected %v, but got %v` + reset + `", expected, actual)
	}
	f, err := os.OpenFile("` + pipePath + `", os.O_WRONLY, os.ModePerm)
	if err != nil {
		t.Fatalf("` + red + `failed to open named pipe: %s` + reset + `", err)
	}
	actualS := strings.TrimSpace(fmt.Sprintf("%d", actual))
	if _, err = f.WriteString(actualS+"\n"); err != nil {
		t.Fatalf("` + red + `failed to write to named pipe: %s` + reset + `", err)
	}
	_ = f.Close()
}
`

			if err = os.WriteFile("checksum_tmp_test.go", []byte(test), 0644); err != nil {
				t.Fatalf("%sfailed to write test file: %s%s", red, err, reset)
			}
			defer func() {
				if err = os.Remove("checksum_tmp_test.go"); err != nil {
					panic(err)
				}
			}()

			cmd = exec.Command("go", "test", "-run", "TestChecksumTmp")
			var actualB []byte
			var actual64 uint64
			var actual uint16

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			u16bChan := make(chan []byte, 1)
			go readPipe(ctx, pipePath, u16bChan, t)
			var tOut []byte
			if tOut, err = cmd.CombinedOutput(); err != nil {
				cancel()
				t.Errorf(testFailed+"\n%s", err, string(tOut))
				if strings.Contains(string(tOut), "checksum_amd64.s:") {
					xerox := bufio.NewScanner(strings.NewReader(string(tOut)))
					for xerox.Scan() {
						if strings.Contains(xerox.Text(), "checksum_amd64.s:") {
							split := strings.Split(xerox.Text(), ":")
							lineStr := strings.TrimSpace(strings.Fields(split[1])[0])
							line, _ := strconv.Atoi(lineStr)
							t.Logf("test failed at line %d", line)
							testXerox := bufio.NewScanner(strings.NewReader(string(asmOut)))
							lc := 0
							for testXerox.Scan() {
								l := testXerox.Text()
								lc++
								if lc == line {
									l += "; <------ FAILED HERE"
								}
								t.Logf("%s", l)
							}
							break
						}
					}
				}
			}

			select {
			case <-ctx.Done():
				t.Fatalf("timed out waiting for test to complete")
			case actualB = <-u16bChan:
			}

			actualS := strings.TrimSpace(string(actualB))
			if actual64, err = strconv.ParseUint(actualS, 10, 16); err != nil {
				t.Errorf(testFailed+" failed to parse output: ", err)
			}
			actual = uint16(actual64)
			if actual != mode.expect {
				t.Errorf(testFailed, fmt.Sprintf("Expected %v, but got %v", mode.expect, actual))
			} else {
				t.Logf(testPassed, actual)
			}
		})
	}
}

func TestASMChecksum(t *testing.T) {
	cmd := exec.Command(strings.Fields(goGen)[0], strings.Fields(goGen)[1:]...)
	t.Logf("generating test ASM with: %s", cmd.String())
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generation failed: %s", cmdOut)
	}
	asmOut, err := os.ReadFile("checksum_amd64.s")
	if err != nil {
		t.Fatalf("failed to read generated ASM: %s", err)
	}
	t.Logf("generated ASM: \n%s", asmOut)
	type test struct {
		name   string
		input  []byte
		expect uint16
	}
	tests := []test{
		{
			name:   "zero length input",
			input:  []byte{},
			expect: 0,
		},
		{
			name:   "zero length input twice",
			input:  []byte{},
			expect: 0,
		},
		{
			name:   "hello",
			input:  []byte("hello"),
			expect: 48173,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := checksum(testCase.input)
			if actual != testCase.expect {
				t.Errorf(testFailed, fmt.Sprintf("Expected %v, but got %s%v%s", testCase.expect, red, actual, reset))
			} else {
				t.Logf(testPassed, actual)
			}
		})
	}
}
