package tokenizer

var (
	wsLUT = [256]bool{' ': true, '\t': true, '\n': true, '\f': true, '\r': true}

	dqStopLUT = [256]bool{'"': true, '\\': true, '\n': true, '\r': true, '\f': true, 0: true}
	sqStopLUT = [256]bool{'\'': true, '\\': true, '\n': true, '\r': true, '\f': true, 0: true}

	urlSpecialLUT = func() (t [256]bool) {
		for _, c := range []byte{'"', '\'', '(', ')', '\\', ' ', '\t', '\n', '\f', '\r'} {
			t[c] = true
		}
		for c := 0; c <= 0x08; c++ {
			t[c] = true
		}
		for c := 0x0E; c <= 0x1F; c++ {
			t[c] = true
		}
		t[0x0B], t[0x7F] = true, true
		return
	}()

	identStartLUT = func() (t [256]bool) {
		for c := 'a'; c <= 'z'; c++ {
			t[c] = true
			t[c-0x20] = true
		}
		for c := 0x80; c < 0x100; c++ {
			t[c] = true
		}
		t['_'] = true
		return
	}()

	identContinueLUT = func() (t [256]bool) {
		for c := '0'; c <= '9'; c++ {
			t[c] = true
		}
		for c := 'a'; c <= 'z'; c++ {
			t[c] = true
			t[c-0x20] = true
		}
		for c := 0x80; c < 0x100; c++ {
			t[c] = true
		}
		t['-'], t['_'] = true, true
		return
	}()
)

func isIdentStartLUT(ch byte) bool    { return identStartLUT[ch] }
func isIdentContinueLUT(ch byte) bool { return identContinueLUT[ch] }

func stopLUT(quote byte) *[256]bool {
	if quote == '\'' {
		return &sqStopLUT
	}
	return &dqStopLUT
}
