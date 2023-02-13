package main

import (
	"fmt"
	"log"

	"github.com/bzick/tokenizer"
)

// TODO: explore Goyacc
// similar work references:
// 		https://pkg.go.dev/modernc.org/goyacc
//		https://github.com/sjjian/oracle-sql-parser
// 		https://github.com/slrtbtfs/goyacc-tutorial
//		https://github.com/sjjian/yacc-examples
// 		https://github.com/monsterxx03/sqlpar
// 		https://github.com/cdstelly/goyacc-sample
// 		https://github.com/pingcap/tidb https://github.com/pingcap/tidb/tree/master/parser -- apparently hard to read

// Exploring lexing/tokenising with: github.com/bzick/tokenizer
func main() {

	// define custom tokens keys
	const (
		tokenString = iota
		tokenDot
		tokenQuotedIdentifier
		tokenMathFunc
		tokenMathBinaryOp
		tokenSetUnaryOp
		tokenSetBinaryOp
		tokenArguments
		tokenBracketStart
		tokenBracketEnd
		tokenSqlSelect
		tokenSqlFrom
		tokenSqlJoin
		tokenSqlJoinLeft
		tokenSqlJoinRight
		tokenSqlJoinFull
		tokenSqlGroupby
		tokenSqlHaving
		tokenSqlLimit
		tokenSqlWhere
		tokenConditionalOp
		tokenSqlInsert
		tokenSqlInsertInto
		tokenSqlValues
		tokenSqlUpdate
		tokenSqlUpdateSet
		tokenCase
		tokenCaseWhen
		tokenCaseElse
		tokenCaseThen
		tokenCaseEnd
	)
	parser := tokenizer.New()
	parser.AllowKeywordUnderscore()
	parser.DefineStringToken(tokenString, `'`, `'`)

	parser.DefineTokens(tokenQuotedIdentifier, []string{`"`})
	parser.DefineTokens(tokenDot, []string{`.`})
	parser.DefineTokens(tokenMathFunc, []string{`sum`, `avg`, `min`, `max`, `count`})
	parser.DefineTokens(tokenMathBinaryOp, []string{`<`, `<=`, `=`, `>=`, `>`, `!=`, `<>`})
	parser.DefineTokens(tokenConditionalOp, []string{`and`, `not`, `or`})
	parser.DefineTokens(tokenSetUnaryOp, []string{`is null`, `is not null`, `is true`, `is false`})
	parser.DefineTokens(tokenSetBinaryOp, []string{`in`, `not in`})
	parser.DefineTokens(tokenBracketStart, []string{`(`})
	parser.DefineTokens(tokenBracketEnd, []string{`)`})
	// Parts of a SQL statement
	parser.DefineTokens(tokenSqlSelect, []string{`select`})
	parser.DefineTokens(tokenSqlFrom, []string{`from`})
	parser.DefineTokens(tokenSqlJoin, []string{`join`})
	parser.DefineTokens(tokenSqlJoinLeft, []string{`left`})
	parser.DefineTokens(tokenSqlJoinRight, []string{`right`})
	parser.DefineTokens(tokenSqlJoinFull, []string{`full`})
	parser.DefineTokens(tokenSqlWhere, []string{`where`})
	parser.DefineTokens(tokenSqlGroupby, []string{`group by`})
	parser.DefineTokens(tokenSqlHaving, []string{`having`})
	parser.DefineTokens(tokenSqlLimit, []string{`limit`})
	parser.DefineTokens(tokenSqlInsert, []string{`insert`})
	parser.DefineTokens(tokenSqlInsertInto, []string{`into`})
	parser.DefineTokens(tokenSqlValues, []string{`values`})
	parser.DefineTokens(tokenSqlUpdate, []string{`update`})
	parser.DefineTokens(tokenSqlUpdateSet, []string{`set`})
	parser.DefineTokens(tokenCase, []string{`case`})
	parser.DefineTokens(tokenCaseWhen, []string{`when`})
	parser.DefineTokens(tokenCaseElse, []string{`else`})
	parser.DefineTokens(tokenCaseThen, []string{`then`})
	parser.DefineTokens(tokenCaseEnd, []string{`end`})

	// create tokens stream
	// stream := parser.ParseString(`user_id = 119 and modified > "2020-01-01 00:00:00" or amount >= 122.34`)
	// stream := parser.ParseString(`avg(year)`)
	stream := parser.ParseString(`select "year" from clients where year > '2012-04-22'`)
	defer stream.Close()

	for stream.IsValid() {
		fmt.Print(stream.CurrentToken().String())
		log.Print(stream.CurrentToken().Key() == tokenQuotedIdentifier)
		stream.GoNext()
	}
}
