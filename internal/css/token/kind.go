package token

type Kind uint8

//go:generate stringer -type Kind kind.go

const (
	None Kind = iota

	EOF                // <EOF-token>
	Ident              // <ident-token>
	Function           // <function-token>
	AtKeyword          // <at-keyword-token>
	Hash               // <hash-token>
	String             // <string-token>
	BadString          // <bad-string-token>
	Url                // <url-token>
	BadUrl             // <bad-url-token>
	Delim              // <delim-token>
	Number             // <number-token>
	Percentage         // <percentage-token>
	Dimension          // <dimension-token>
	Whitespace         // <whitespace-token>
	Comment            // <comment-token>
	CDO                // <CDO-token>
	CDC                // <CDC-token>
	Colon              // <colon-token>
	Semicolon          // <semicolon-token>
	Comma              // <comma-token>
	OpenParenthesis    // <(-token>
	OpenSquareBracket  // <[-token>
	OpenCurlyBracket   // <{-token>
	CloseParenthesis   // <)-token>
	CloseSquareBracket // <]-token>
	CloseCurlyBracket  // <}-token>
	// UnicodeRange       // <unicode-range-token>
)
