package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

const cellCount = 30000

type machine struct {
	memory       [cellCount]byte
	dp           uint
	instructions []byte
	ip           uint
	input        io.Reader
	output       io.Writer
}

func (m *machine) loadProgram(prog io.Reader) (err error) {
	m.instructions, err = ioutil.ReadAll(prog)
	return err
}

func (m *machine) run() {
	for int(m.ip) < len(m.instructions) {
		switch m.instructions[m.ip] {
		case '>':
			m.dp++
		case '<':
			m.dp--
		case '+':
			m.memory[m.dp]++
		case '-':
			m.memory[m.dp]--
		case '.':
			_, err := m.output.Write(m.memory[m.dp : m.dp+1])
			if err != nil {
				log.Fatalln("Error writing", err)
			}
		case ',':
			_, err := m.input.Read(m.memory[m.dp : m.dp+1])
			if err != nil {
				log.Fatalln("Error reading", err)
			}
		case '[':
			if m.memory[m.dp] == 0 {
				depth := 1
				for depth != 0 {
					m.ip++
					switch m.instructions[m.ip] {
					case '[':
						depth++
					case ']':
						depth--
					}
				}
			}
		case ']':
			if m.memory[m.dp] != 0 {
				depth := 1
				for depth != 0 {
					m.ip--
					switch m.instructions[m.ip] {
					case ']':
						depth++
					case '[':
						depth--
					}
				}
			}
		}
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
