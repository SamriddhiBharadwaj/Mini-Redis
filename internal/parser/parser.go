import "bufio"

// Parser contains the logic to read from a raw tcp connection and parse commands.
// Raw TCP bytes → Buffered Reader → Line parsing → Command parsing
type Parser struct {
	conn net.Conn
	r    *bufio.Reader
	// Used for inline parsing
	line []byte
	pos  int
}

// NewParser returns a new Parser that reads from the given connection.
// wraps connection with a buffer
func NewParser(conn net.Conn) *Parser {
	return &Parser{
		conn: conn,
		r:    bufio.NewReader(conn),
		line: make([]byte, 0),
		pos:  0,
	}
}

// helper functions for inline command parsing


// returns current char or eof ('\r')
func (p *Parser) current() byte {
	if p.atEnd() {
		return '\r'
	}
	return p.line[p.pos]
}

// moves pointer position by one
func (p *Parser) advance() {
	p.pos++
}

// checks if pointer at end of line i.e if parsing is finished
func (p *Parser) atEnd() bool {
	return p.pos >= len(p.line)
}


func (p *Parser) readLine() ([]byte, error) {
	// reads till '\r'
	line, err := p.r.ReadBytes('\r')
	if err != nil {
		return nil, err
	}
	// consume '\n' at eof
	if _, err := p.r.ReadByte(); err != nil {
		return nil, err
	}
	// removes '\r'
	return line[:len(line)-1], nil
}