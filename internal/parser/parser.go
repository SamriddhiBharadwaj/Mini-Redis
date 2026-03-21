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

// Command implements the behavior of the commands.
type Command struct {
	args []string
	conn net.Conn
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
	// consume '\n' at eol
	if _, err := p.r.ReadByte(); err != nil {
		return nil, err
	}
	// removes '\r'
	return line[:len(line)-1], nil
}


// command parses and returns a Command.
func (p *Parser) command() (Command, error) {
	b, err := p.r.ReadByte()
	if err != nil {
		return Command{}, err
	}
	// if first char is '*', handle via RESP
	if b == '*' {
		log.Println("resp array")
		return p.respArray()
	} else {
		// else handle like inline command
		// since first char is already read, read rest of line from buffer
		line, err := p.readLine()
		if err != nil {
			return Command{}, err
		}
		p.pos = 0
		// append first char b to line
		p.line = append([]byte{}, b)
		// line... used unpack entire slice (variadic expansion)
		p.line = append(p.line, line...)
		// call inline function
		return p.inline()
	}
}


// inline parses an inline message and returns a Command. Returns an error when there's
// a problem reading from the connection or parsing the command.
func (p *Parser) inline() (Command, error) {
	// skip initial whitespace if any
	for p.current() == ' ' {
		p.advance()
	}
	// initialize cmd to Command object and set conn to parser's conn
	cmd := Command{conn: p.conn}
	// extract till eol
	for !p.atEnd() {
		arg, err := p.consumeArg()
		if err != nil {
			return cmd, err
		}
		if arg != "" {
			cmd.args = append(cmd.args, arg)
		}
	}
	return cmd, nil
}

// consumeArg reads an argument from the current line.
func (p *Parser) consumeArg() (s string, err error) {
	// skip initial whitespace if any
	for p.current() == ' ' {
		p.advance()
	}
	// check for quoted string
	if p.current() == '"' {
		// skip opening '"' as required by consumeString
		p.advance()
		buf, err := p.consumeString()
		return string(buf), err
	}
	// pass arguments normally by building string
	for !p.atEnd() && p.current() != ' ' && p.current() != '\r' {
		s += string(p.current())
		p.advance()
	}
	return
}

// consumeString reads a string argument from the current line (handles \ and " chars)
func (p *Parser) consumeString() (s []byte, err error) {
	// assumes initial '"' has been consumed before entering function
	for p.current() != '"' && !p.atEnd() {
		cur := p.current()
		p.advance()
		next := p.current()
		// eg: "hello \"world\""
		if cur == '\\' && next == '"' {
			s = append(s, '"')
			p.advance()
		} else {
			s = append(s, cur)
		}
	}
	// request doesnt end with '"'
	if p.current() != '"' {
		return nil, errors.New("unbalanced quotes in request")
	}
	p.advance()
	return
}

//respArray parses a RESP array and returns a Command. Returns an error when there's
// a problem reading from the connection.
func (p *Parser) respArray() (Command, error) {
	cmd := Command{}
	elementsStr, err := p.readLine()
	if err != nil {
		return cmd, err
	}
	// read number of elements from string
	elements, _ := strconv.Atoi(string(elementsStr))
	log.Println("Elements", elements)
	// loop through elements
	for i := 0; i < elements; i++ {
		// read type of each element
		tp, err := p.r.ReadByte()
		if err != nil {
			return cmd, err
		}
		switch tp {
		// integer case
		case ':':
			arg, err := p.readLine()
			if err != nil {
				return cmd, err
			}
			cmd.args = append(cmd.args, string(arg))
		// bulk string: $5\r\nhello\r\n
		case '$':
			arg, err := p.readLine()
			if err != nil {
				return cmd, err
			}
			// read length
			length, _ := strconv.Atoi(string(arg))
			// create new empty slice
			text := make([]byte, 0)
			for i := 0; len(text) <= length; i++ {
				line, err := p.readLine()
				if err != nil {
					return cmd, err
				}
				// append to buffer 
				text = append(text, line...)
			}
			// trim to exact length append buffer to Command struct object
			cmd.args = append(cmd.args, string(text[:length]))
		// nested array
		case '*':
			// recursive call
			next, err := p.respArray()
			if err != nil {
				return cmd, err
			}
			cmd.args = append(cmd.args, next.args...)
		}
	}
	return cmd, nil
} 