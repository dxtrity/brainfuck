package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	defaultMemorySize = 30000
	snapshotSize      = 2
	defaultImageWidth = 16
	helpText          = `
    Simple Brainfuck Interpreter 

    BF [options] <input>
         ^ can be a file or brainfuck string:
      e.g. BF program.bf
      e.g. BF +++.

    Options:
    -memory <size>  set the memory size (default: 30000)
    -snapx <size>  set the snapshot start size (default: 2)
    -snapy <size>  set the snapshot end size (default: 2)

    Commands:
    >  increment the data pointer
    <  decrement the data pointer
    +  increment the byte at the data pointer
    -  decrement the byte at the data pointer
    .  output the byte at the data pointer
    ,  accept one byte of input
    [  jump forward past the matching ] if the byte at the pointer is 0
    ]  jump back to the matching [ if the byte at the pointer is nonzero
    #  print debug information
`
)

type BrainfuckInterpreter struct {
	memory            []byte
	imageWidth        int
	snapshotStartX    int
	snapshotEndX      int
	pointer           int
	input             io.Reader
	output            io.Writer
	maxPointerReached int
	nonZeroCells      map[int]bool
}

func NewInterpreter(memSize int, imgWidth int, snapStartX int, snapEndX int, input io.Reader, output io.Writer) *BrainfuckInterpreter {
	return &BrainfuckInterpreter{
		memory:         make([]byte, memSize),
		imageWidth:     imgWidth,
		snapshotStartX: snapStartX,
		snapshotEndX:   snapEndX,
		input:          input,
		output:         output,
		nonZeroCells:   make(map[int]bool),
	}
}

func (bf *BrainfuckInterpreter) CalculateMemoryUsage() (int, int) {
	nonZeroCount := 0
	for i := 0; i <= bf.maxPointerReached; i++ {
		if bf.memory[i] != 0 {
			nonZeroCount++
		}
	}

	highestAddressUsed := bf.maxPointerReached + 1

	return nonZeroCount, highestAddressUsed
}

func (bf *BrainfuckInterpreter) PrintDebugInfo() {
	nonZeroCount, highestAddress := bf.CalculateMemoryUsage()

	fmt.Fprintf(bf.output, "\n\n")
	color.Set(color.BgHiWhite, color.Bold)
	fmt.Fprintf(bf.output, "    debug information    ")
	color.Unset()

	fmt.Fprintf(bf.output, "\n")

	color.Set(color.BgBlue)
	fmt.Fprintf(bf.output, " pointer location ")
	color.Unset()

	color.Set(color.Bold)
	fmt.Fprintf(bf.output, "  %d\n", bf.pointer)
	color.Unset()

	color.Set(color.BgRed)
	fmt.Fprintf(bf.output, " pointer value    ")
	color.Unset()

	color.Set(color.Bold)
	fmt.Fprintf(bf.output, "  %d", bf.memory[bf.pointer])
	if bf.memory[bf.pointer] >= 32 && bf.memory[bf.pointer] <= 126 {
		fmt.Fprintf(bf.output, " (%s)", string(bf.memory[bf.pointer]))
	}
	fmt.Fprintf(bf.output, "\n")
	color.Unset()

	color.Set(color.BgGreen)
	fmt.Fprintf(bf.output, " memory usage     ")
	color.Unset()

	color.Set(color.Bold)
	fmt.Fprintf(bf.output, "  %.2f%%\n", float64(highestAddress)/float64(len(bf.memory))*100)
	color.Unset()

	color.Set(color.BgCyan)
	fmt.Fprintf(bf.output, " non-zero cells   ")
	color.Unset()

	color.Set(color.Bold)
	fmt.Fprintf(bf.output, "  %d\n", nonZeroCount)
	color.Unset()
}

func (bf *BrainfuckInterpreter) PrintMemorySnapshot() {
	fmt.Fprintf(bf.output, "\n\n")
	color.Set(color.BgHiMagenta, color.Bold)
	fmt.Fprintf(bf.output, "      memory snapshot      ")
	color.Unset()
	fmt.Fprintf(bf.output, "\n")

	start := max(bf.pointer-bf.snapshotStartX, 0)
	end := bf.pointer + bf.snapshotEndX
	if end >= len(bf.memory) {
		end = len(bf.memory) - 1
	}

	for i := start; i <= end; i++ {
		if i == bf.pointer {
			color.Set(color.BgRed)
		} else {
			color.Set(color.BgBlue)
		}
		if i < 10 {
			fmt.Fprintf(bf.output, " [0%d] ", i)
		} else {
			fmt.Fprintf(bf.output, " [%d] ", i)
		}
		color.Unset()

		color.Set(color.Bold)
		fmt.Fprintf(bf.output, " %d", bf.memory[i])
		if bf.memory[i] >= 32 && bf.memory[i] <= 126 {
			fmt.Fprintf(bf.output, " (%s)", string(bf.memory[i]))
		}
		fmt.Fprintf(bf.output, "\n")
		color.Unset()
	}
	fmt.Fprintf(bf.output, "\n")
}

func (bf *BrainfuckInterpreter) Run(code string) error {
	for i := 0; i < len(code); i++ {
		switch code[i] {
		case '>':
			if bf.pointer >= len(bf.memory)-1 {
				return errors.New("pointer out of bounds (out of memory)")
			}
			bf.pointer++
			if bf.pointer > bf.maxPointerReached {
				bf.maxPointerReached = bf.pointer
			}
		case '<':
			if bf.pointer <= 0 {
				return errors.New("pointer out of bounds")
			}
			bf.pointer--
		case '+':
			bf.memory[bf.pointer]++
		case '-':
			bf.memory[bf.pointer]--
		case '.':
			fmt.Fprintf(bf.output, "%c", bf.memory[bf.pointer])
		case ',':
			var inputChar [1]byte
			_, err := bf.input.Read(inputChar[:])
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			bf.memory[bf.pointer] = inputChar[0]
		case '#':
			bf.PrintDebugInfo()
		case '@':
			bf.PrintMemorySnapshot()
		case '[':
			if bf.memory[bf.pointer] == 0 {
				loop := 1
				for loop > 0 && i < len(code)-1 {
					i++
					if code[i] == '[' {
						loop++
					} else if code[i] == ']' {
						loop--
					}
				}
				if loop != 0 {
					return errors.New("unmatched brackets (missing ']')")
				}
			}
		case ']':
			if bf.memory[bf.pointer] != 0 {
				loop := 1
				for loop > 0 && i > 0 {
					i--
					if code[i] == '[' {
						loop--
					} else if code[i] == ']' {
						loop++
					}
				}
				if loop != 0 {
					return errors.New("unmatched brackets: missing '['")
				}
			}
		}
	}
	return nil
}

func (bf *BrainfuckInterpreter) RunImage(code string) error {
	squareCount := 0 // Counter for rendered squares
	for i := 0; i < len(code); i++ {
		switch code[i] {
		case '>':
			if bf.pointer >= len(bf.memory)-1 {
				return errors.New("pointer out of bounds (out of memory)")
			}
			bf.pointer++
			if bf.pointer > bf.maxPointerReached {
				bf.maxPointerReached = bf.pointer
			}
		case '<':
			if bf.pointer <= 0 {
				return errors.New("pointer out of bounds")
			}
			bf.pointer--
		case '+':
			bf.memory[bf.pointer]++
		case '-':
			bf.memory[bf.pointer]--
		case '.':
			colorValue := bf.memory[bf.pointer]
			bf.drawColoredSquare(colorValue)
			squareCount++
			if squareCount%bf.imageWidth == 0 {
				fmt.Fprintf(bf.output, "\n") // Add newline after 16 squares
			}
		case ',':
			var inputChar [1]byte
			_, err := bf.input.Read(inputChar[:])
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			bf.memory[bf.pointer] = inputChar[0]
		case '#':
			bf.PrintDebugInfo()
		case '@':
			bf.PrintMemorySnapshot()
		case '[':
			if bf.memory[bf.pointer] == 0 {
				loop := 1
				for loop > 0 && i < len(code)-1 {
					i++
					if code[i] == '[' {
						loop++
					} else if code[i] == ']' {
						loop--
					}
				}
				if loop != 0 {
					return errors.New("unmatched brackets (missing ']')")
				}
			}
		case ']':
			if bf.memory[bf.pointer] != 0 {
				loop := 1
				for loop > 0 && i > 0 {
					i--
					if code[i] == '[' {
						loop--
					} else if code[i] == ']' {
						loop++
					}
				}
				if loop != 0 {
					return errors.New("unmatched brackets: missing '['")
				}
			}
		}
	}
	return nil
}

func (bf *BrainfuckInterpreter) drawColoredSquare(value byte) {
	// Map the byte value to a color. You can customize this mapping.
	var colorAttribute color.Attribute
	switch {
	case value < 32:
		colorAttribute = color.BgHiBlack
	case value < 64:
		colorAttribute = color.BgRed
	case value < 96:
		colorAttribute = color.BgGreen
	case value < 128:
		colorAttribute = color.BgYellow
	case value < 160:
		colorAttribute = color.BgBlue
	case value < 192:
		colorAttribute = color.BgMagenta
	case value < 224:
		colorAttribute = color.BgCyan
	default:
		colorAttribute = color.BgHiWhite
	}

	color.Set(colorAttribute)
	fmt.Fprintf(bf.output, "  ")
	color.Unset()
}

func IsBrainfuckCommand(input string) bool {
	if len(input) == 0 {
		return false
	}

	switch input[0] {
	case '+', '-', '[', ']', '.', ',', '>', '<', '#':
		return true
	default:
		return false
	}
}

func IsBrainfuckFile(filename string) bool {
	return len(filename) > 3 && strings.HasSuffix(filename, ".bf")
}

func main() {
	_title := color.New(color.BgBlue)
	_info := color.New(color.BgWhite)
	_eg := color.New(color.FgWhite)

	memSize := flag.Int("memory", defaultMemorySize, "-memory <size>")
	helpFlag := flag.Bool("help", false, "-help")
	snapX := flag.Int("snapx", snapshotSize, "-snapx <size>")
	snapY := flag.Int("snapy", snapshotSize, "-snapy <size>")
	devFlag := flag.Bool("proc", false, "-tr")
	imageFlag := flag.Bool("image", false, "-image")
	imageWidthFlag := flag.Int("w", defaultImageWidth, "-w <size>")

	flag.Parse()

	args := flag.Args()

	if *devFlag {
		//TODO:  DEV STUFF

		var a byte

		fmt.Println("Enter your input: ")
		_, err := fmt.Scan(&a)
		check(err)
		fmt.Printf("\n%c\n", a)
		os.Exit(0)
	}

	if *helpFlag {
		_title.Println(" Simple Brainfuck Interpreter ")
		fmt.Println()

		_info.Println("BF [options] <input>")
		fmt.Println("               ^ can be a file or brainfuck string:")
		_eg.Println("  e.g. BF program.bf")
		_eg.Println("  e.g. BF +++#")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -image           change to image rendering mode")
		fmt.Println("  -w <size>        set the image width (default:", defaultImageWidth, ")")
		fmt.Println("  -memory <size>   set the memory size (default:", defaultMemorySize, ")")
		fmt.Println("  -help            display help message")
		fmt.Println("  -snapx <size>    set the snapshot start size (default:", snapshotSize, ")")
		fmt.Println("  -snapy <size>    set the snapshot end size (default:", snapshotSize, ")")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  >   increment the data pointer")
		fmt.Println("  <   decrement the data pointer")
		fmt.Println("  +   increment the byte at the data pointer")
		fmt.Println("  -   decrement the byte at the data pointer")
		fmt.Println("  .   output the byte at the data pointer")
		fmt.Println("  ,   accept one byte of input")
		fmt.Println("  [   jump forward past the matching ] if the byte at the pointer is 0")
		fmt.Println("  ]   jump back to the matching [ if the byte at the pointer is nonzero")
		fmt.Println("  #   print debug information")
		fmt.Println("  @   print a memory snapshot")

		return
	}

	if len(args) == 0 {
		fmt.Print("[USAGE]: BF [options] <input>\n\n")
		fmt.Println("for more help: BF -help")
		return
	}

	var code string
	var err error

	if IsBrainfuckFile(args[0]) {
		data, readErr := os.ReadFile(args[0])
		if readErr != nil {
			fmt.Printf("Error reading file: %v\n", readErr)
			os.Exit(1)
		}
		code = string(data)
	} else if IsBrainfuckCommand(args[0]) {
		code = args[0]
	} else {
		fmt.Println("Invalid input. Please provide a valid Brainfuck command or file with .bf extension.")
		os.Exit(1)
	}

	if *imageFlag {
		interpreter := NewInterpreter(*memSize, *imageWidthFlag, *snapX, *snapY, os.Stdin, os.Stdout)
		err = interpreter.RunImage(code)
		if err != nil {
			fmt.Printf("\n")
			color.Set(color.BgRed, color.FgHiWhite, color.Bold)
			fmt.Printf("  Execution error  ")
			color.Unset()
			fmt.Printf(": %v\n", err)
			os.Exit(1)
		}
	} else {
		interpreter := NewInterpreter(*memSize, *imageWidthFlag, *snapX, *snapY, os.Stdin, os.Stdout)
		err = interpreter.Run(code)
		if err != nil {
			fmt.Printf("\n")
			color.Set(color.BgRed, color.FgHiWhite, color.Bold)
			fmt.Printf("  Execution error  ")
			color.Unset()
			fmt.Printf(": %v\n", err)
			os.Exit(1)
		}
	}
}

func check(e error) {
	if e != nil {
		fmt.Printf("%v", e)
	}
}
