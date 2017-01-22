package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

const cellCount = 30000

type instruction struct {
	opCode byte
	op     func()
}

type machine struct {
	memory       [cellCount]byte
	dp           uint
	instructions []instruction
	ip           uint
	input        io.Reader
	output       io.Writer
}

func (m *machine) incData()    { m.memory[m.dp] += 1 }
func (m *machine) decData()    { m.memory[m.dp] -= 1 }
func (m *machine) incDataPtr() { m.dp += 1 }
func (m *machine) decDataPtr() { m.dp -= 1 }
func (m *machine) readInput() {
	_, err := m.input.Read(m.memory[m.dp : m.dp+1])
	if err != nil {
		log.Fatalln("Error reading", err)
	}
}
func (m *machine) writeOutput() {
	_, err := m.output.Write(m.memory[m.dp : m.dp+1])
	if err != nil {
		log.Fatalln("Error writing", err)
	}
}
func (m *machine) loop() {
	if m.memory[m.dp] == 0 {
		depth := 1
		for depth != 0 {
			m.ip++
			switch m.instructions[m.ip].opCode {
			case '[':
				depth++
			case ']':
				depth--
			}
		}
	}
}

func (m *machine) endLoop() {
	if m.memory[m.dp] != 0 {
		depth := 1
		for depth != 0 {
			m.ip--
			switch m.instructions[m.ip].opCode {
			case ']':
				depth++
			case '[':
				depth--
			}
		}
	}
}

func (m *machine) loadProgram(prog io.Reader) (err error) {
	scanner := bufio.NewScanner(prog)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		i := instruction{opCode: scanner.Text()[0]}
		switch i.opCode {
		case '>':
			i.op = m.incDataPtr
		case '<':
			i.op = m.decDataPtr
		case '+':
			i.op = m.incData
		case '-':
			i.op = m.decData
		case '.':
			i.op = m.writeOutput
		case ',':
			i.op = m.readInput
		case '[':
			i.op = m.loop
		case ']':
			i.op = m.endLoop
		default:
			continue // All other characters ignored
		}
		m.instructions = append(m.instructions, i)
	}
	return scanner.Err()
}

func (m *machine) run() {
	for int(m.ip) < len(m.instructions) {
		m.instructions[m.ip].op()
		m.ip++
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage bf-go program.bf")
	}
	prog, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	m := &machine{
		input:  os.Stdin,
		output: os.Stdout,
	}
	err = m.loadProgram(prog)
	if err != nil {
		log.Fatal(err)
	}
	m.run()
}
