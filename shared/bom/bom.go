package bom

import (
	"bufio"
	"io"
)

// utf8BOM is the UTF-8 byte-order mark. Windows editors and PowerShell
// Set-Content -Encoding utf8 often write it; GNU tools ignore it.
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// Reader wraps r and strips a leading UTF-8 BOM if present.
func Reader(r io.Reader) io.Reader {
	br := bufio.NewReader(r)
	preview, err := br.Peek(3)
	if err == nil && len(preview) >= 3 &&
		preview[0] == utf8BOM[0] && preview[1] == utf8BOM[1] && preview[2] == utf8BOM[2] {
		_, _ = br.Discard(3)
	}
	return br
}
