package mdedit

import "bytes"

type action struct {
	cmd    command
	change change
}

type command struct {
	modChar     byte
	opCount     int
	opChar      byte
	cmdCount    int
	cmdChar     byte
	motionCount int
	motionChar1 byte
	motionChar2 byte
}

func (c *command) process(char byte) {
	if char == '0' && c.motionCount == 0 {
		c.motionChar1 = '0'
	} else if char >= '0' && char <= '9' {
		c.motionCount = (c.motionCount * 10) + int(char-'0')
	}
	switch char {
	case 'g', 'z':
		c.modChar = char
	case '.', 'u', 'I', 'S', 'o', 'O', 'C', 'A', 'x':
		c.cmdChar = char
	case 'c', 'd', 'y':
		if c.opChar == char {
			c.motionChar1 = char
		} else {
			c.opChar = char
		}
	case 'i', 'a':
		if c.opChar == 0 {
			c.cmdChar = char
		} else {
			c.motionChar1 = char
		}
	case 'w', 'e', 'b':
		if c.motionChar1 == 0 {
			c.motionChar1 = char
		} else {
			c.motionChar2 = char
		}
	case ' ':
		if c.modChar != 0 {
			c.cmdChar = char
			break
		}
		fallthrough
	case 'h', 'l', '$':
		c.motionChar1 = char
	case 'j', 'k', 'H', 'L':
		c.motionChar1 = char
	}
}

var motionChars = []byte("jkhl LHweWEb0$")

func (c *command) hasMotion() bool {
	return bytes.IndexByte(motionChars, c.motionChar1) != -1 ||
		(c.motionChar2 != 0 && (c.motionChar1 == 'i' || c.motionChar1 == 'a'))
}

type change struct {
	typ     changeType
	from    position
	content content
}

type changeType byte

const (
	changeAddition changeType = iota
	changeDeletion
)

type content struct {
	simple []byte
	lines  []line
}
