// Copyright 2016 the V8 project authors. All rights reserved.
// This file adapts the public parser entry points declared by V8's
// regexp-parser.h. See LICENSE.v8 for attribution and license terms.
package regexpengine

// V8 accepts both one-byte and UTF-16 strings. The byte entry points keep the
// source byte-backed; the cursor reads each byte as one Latin-1 code unit and
// widens only small metadata slices when a parser helper needs UTF-16 units.

func v8ASTParseRegExpOneByte(input []byte, flags Flags, stackOverflow func() bool) (v8ASTParserCompileData, error) {
	return v8ASTParseRegExpOneByteWithOptions(input, flags, stackOverflow, v8ASTParserOptions{})
}

func v8ASTParseRegExpOneByteWithOptions(input []byte, flags Flags, stackOverflow func() bool, options v8ASTParserOptions) (v8ASTParserCompileData, error) {
	return v8ASTNewParserOneByte(input, flags, v8ASTBuildAST, stackOverflow, options).parseCompileData()
}

func v8ASTParseRegExpUTF16(input []uint16, flags Flags, stackOverflow func() bool) (v8ASTParserCompileData, error) {
	return v8ASTParseRegExpUTF16WithOptions(input, flags, stackOverflow, v8ASTParserOptions{})
}

func v8ASTParseRegExpUTF16WithOptions(input []uint16, flags Flags, stackOverflow func() bool, options v8ASTParserOptions) (v8ASTParserCompileData, error) {
	return v8ASTNewParserWithOptions(input, flags, false, v8ASTBuildAST, stackOverflow, options).parseCompileData()
}

// v8ASTVerifyRegExpSyntaxOneByte corresponds to the uint8_t specialization of
// Parser::VerifyRegExpSyntax. Its CompileData equivalent contains diagnostics;
// AST-only metadata stays empty in syntax-only mode.
func v8ASTVerifyRegExpSyntaxOneByte(input []byte, flags Flags, stackOverflow func() bool) (v8ASTParserCompileData, error) {
	return v8ASTVerifyRegExpSyntaxOneByteWithOptions(input, flags, stackOverflow, v8ASTParserOptions{})
}

func v8ASTVerifyRegExpSyntaxOneByteWithOptions(input []byte, flags Flags, stackOverflow func() bool, options v8ASTParserOptions) (v8ASTParserCompileData, error) {
	return v8ASTNewParserOneByte(input, flags, v8ASTVerifySyntax, stackOverflow, options).parseCompileData()
}

// v8ASTVerifyRegExpSyntaxUTF16 corresponds to the base::uc16 specialization.
func v8ASTVerifyRegExpSyntaxUTF16(input []uint16, flags Flags, stackOverflow func() bool) (v8ASTParserCompileData, error) {
	return v8ASTVerifyRegExpSyntaxUTF16WithOptions(input, flags, stackOverflow, v8ASTParserOptions{})
}

func v8ASTVerifyRegExpSyntaxUTF16WithOptions(input []uint16, flags Flags, stackOverflow func() bool, options v8ASTParserOptions) (v8ASTParserCompileData, error) {
	return v8ASTNewParserWithOptions(input, flags, false, v8ASTVerifySyntax, stackOverflow, options).parseCompileData()
}

// v8ASTVerifySyntaxUTF16 retains the parser-only convenience API used by
// syntax checks that need only success or the first parser error.
func v8ASTVerifySyntaxUTF16(input []uint16, flags Flags, stackOverflow func() bool) error {
	_, err := v8ASTVerifyRegExpSyntaxUTF16(input, flags, stackOverflow)
	return err
}
