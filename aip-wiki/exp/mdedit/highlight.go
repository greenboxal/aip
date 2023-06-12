package mdedit

import "bytes"

const (
	// blocks
	mdHeading uint16 = 1 << iota
	mdBlockquote
	mdCodeBlock
	mdThematicBreak
	// inlines
	mdItalic
	mdStrong
	mdCodeSpan
	mdListMarker
	mdLinkURL
)

type mdStyleMark struct {
	col   int
	value uint16
}

type styleBuilder struct {
	markers [][]mdStyleMark
	row     int
}

func (s *styleBuilder) startNewRow() {
	s.row++
	if s.row == len(s.markers) {
		s.markers = append(s.markers, []mdStyleMark{})
	}
}

func (s *styleBuilder) add(v uint16, col int) {
	s.markers[s.row] = append(s.markers[s.row], mdStyleMark{
		col:   col,
		value: v,
	})
}

type mdHighlighter struct{}

func (mdHighlighter) highlight(buf *buffer) [][]mdStyleMark {
	const (
		bqStarted uint8 = iota + 1
		bqHitChar
		codeSpan1 uint8 = iota + 1
		codeSpan2
	)
	var (
		marks        uint16
		maybeHeading bool
		bqState      uint8
		inCodeSpan   uint8
		inEmphasis1  byte
		inEmphasis2  byte
	)
	sb := styleBuilder{markers: make([][]mdStyleMark, len(buf.lines))}

lineloop:
	for row := 0; row < len(buf.lines); row++ {
		line := buf.lines[row].text
		if len(line) == 0 {
			if marks&mdBlockquote == mdBlockquote && bqState == bqHitChar {
				marks = marks &^ mdBlockquote
				bqState = 0
			}
		}
		marks = marks &^ mdHeading
		sb.add(marks, 0)

		var start int
		for start = 0; start < len(line); start++ {
			if line[start] != ' ' && line[start] != '\t' {
				break
			}
		}

		for col := start; col < len(line); col++ {
			char := line[col]

			if col == start {
				switch char {
				case '#':
					maybeHeading = true
				case '>':
					marks |= mdBlockquote
					sb.add(marks, col)
					bqState = bqStarted
				case '*', '+', '-':
					if col+1 < len(line)-1 && line[col+1] == ' ' {
						sb.add(marks|mdListMarker, col)
						sb.add(marks, col+1)
						col++
					}
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					if col+2 < len(line)-1 && line[col+1] == '.' && line[col+2] == ' ' {
						sb.add(marks|mdListMarker, col)
						sb.add(marks, col+2)
						col += 2
					}
				}
			}
			if bqState == bqStarted && char != ' ' && char != '\t' && char != '>' {
				bqState = bqHitChar
			}
			if maybeHeading && char != '#' {
				if char == ' ' {
					marks |= mdHeading
					sb.add(marks, start)
				}
				maybeHeading = false
			}
			switch char {
			case '*', '_':
				var prev, next byte
				if col > 0 {
					prev = line[col-1]
				}
				if col+1 < len(line) {
					next = line[col+1]
				}

				isPrevBlank := (prev == 0 || prev == ' ' || prev == '\t')
				isNextBlank := (next == 0 || next == ' ' || next == '\t')

				if char == next {
					switch {
					case inEmphasis2 == 0 && isPrevBlank:
						var next2 byte
						if col+2 < len(line) {
							next2 = line[col+2]
						}
						if next2 != ' ' && next2 != '\t' {
							inEmphasis2 = char
							marks |= mdStrong
							sb.add(marks, col)
							col++
						}
					case inEmphasis2 == char && !isPrevBlank:
						inEmphasis2 = 0
						marks &^= mdStrong
						sb.add(marks, col+2)
						col++
					}
				} else {
					switch {
					case inEmphasis1 == 0 && isPrevBlank && !isNextBlank:
						inEmphasis1 = char
						marks |= mdItalic
						sb.add(marks, col)
					case inEmphasis1 != 0 && !isPrevBlank && (isNextBlank || isPunct(next) || next == inEmphasis2):
						inEmphasis1 = 0
						marks &^= mdItalic
						if col < len(line)-1 {
							sb.add(marks, col+1)
						}
					}
				}
			case '`':
				var next byte
				if col+1 < len(line) {
					next = line[col+1]
				}
				if col == start && next == '`' && col+2 < len(line) && line[col+2] == '`' {
					sb.add(marks|mdCodeBlock, start)
					delimCodeBlock := []byte("```")
					for {
						row++
						if row >= len(buf.lines) {
							break lineloop
						}
						sb.startNewRow()
						sb.add(marks|mdCodeBlock, start)
						ln := &buf.lines[row]
						if bytes.Equal(ln.text[ln.startingIndex():], delimCodeBlock) {
							break
						}
					}
					sb.startNewRow()
					sb.add(marks, 0)
					continue lineloop
				}
				switch inCodeSpan {
				case 0:
					inCodeSpan = codeSpan1
					marks |= mdCodeSpan
					sb.add(marks, col)
					if next == '`' {
						inCodeSpan = codeSpan2
						col += 2
					}
				case codeSpan1:
					marks &^= mdCodeSpan
					inCodeSpan = 0
					if col < len(line)-1 {
						sb.add(marks, col+1)
					}
				case codeSpan2:
					if next == '`' {
						marks &^= mdCodeSpan
						inCodeSpan = 0
						col++
						if col+1 < len(line)-1 {
							sb.add(marks, col+2)
						}
					}
				}
			}
		}
		if maybeHeading {
			sb.add(marks|mdHeading, start)
			maybeHeading = false
		}
		sb.startNewRow()
	}

	return sb.markers
}

func isPunct(char byte) bool {
	switch char {
	case '.', ',', '!', '?':
		return true
	default:
		return false
	}
}
