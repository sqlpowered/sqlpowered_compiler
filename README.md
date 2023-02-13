# sqlpowered_compiler
The sqlpowered compiler engine

## Aims
To make a SQL sanitisation compiler which accepts SQL and compiles to verify the user is allowed to run each part of the statement. If they're allowed we can then build SQL from the compiled representation of the query. 

Safety is the main goal. Achieved by applying user roles and general restrictions, this can make SQL databases accessible without needing to have a bespoke API endpoint each time, a single generalised endpoint should be able to cover most cases. 

**DB access like this has broad usecases**
- Edge computing
- Frontend
- Other traditional backends
- IOT/embedded
- Mobile

## Authentication JWTs

- Authentication gives user a signed, verifiable JWT with roles inside.
- Roles correspond to permissions sets in the API the user can define.
- Permissions like what the user can SELECT, INSERT, DELETE, UPDATE.
- Permissions can define what columns they are allowed to access/operate on too.
- Permissions can require that actions use a particular column as well to comply with security policy or a particular cost tier the user is in.

## How can we trust the SQL?

- RAW SQL is a rich and powerful language, with *many* security problems as it was not designed to work on the internet!
- Our aim is to support a secure subset of SQL which is widely supported across database vendors.

## Trusted subset of SQL

| SQL Features in select    | Supported                   |
|---------------------------|-----------------------------|
| select                    | ✓* (columns must be permitted)|
| from                      | ✓ |
| join                      | ✓ (inner, left, right, full) |
| where                     | ✓* (columns must be permitted) |
| group by                  | ✓ |
| having                    | ✓ |
| limit                     | ✓* (with API defined limits on max) |
| offset                    | ✓ |
| common table expressions  | ✓ (for security, requires a specific name prefix) |
| column naming             | ✓ |
| table naming              | ✗ (makes verification more difficult) |
| table unions              | ✗ (poor security profile) |
| subqueries                | ✗ (CTEs are faster and give better functionality, simplify and standardise the queries) |

| SQL Features in update    | Supported                   |
|---------------------------|-----------------------------|
| update                    | ✓* (columns must be permitted) |
| set                       | ✓* (mandatory, cannot be omitted) |
| where                     | ✓* (mandatory, cannot be omitted, columns must be permitted) |
| common table expressions  | ✗ (easier to secure, separation of concerns by query type) |
| subqueries                | ✗ (easier to secure, separation of concerns by query type) |

| SQL Features in delete    | Supported                   |
|---------------------------|-----------------------------|
| delete                    | ✓* (where user permissions allow) |
| where                     | ✓* (mandatory, cannot be omitted, cannot be empty, columns must be permitted) |
| common table expressions  | ✗ (easier to secure, separation of concerns by query type) |
| subqueries                | ✗ (easier to secure, separation of concerns by query type) |

| SQL Features in insert    | Supported                   |
|---------------------------|-----------------------------|
| insert                    | ✓* (where user permissions allow) |
| values                    | ✓* (mandatory, cannot be omitted, cannot be empty) |
| common table expressions  | ✗ (easier to secure, separation of concerns by query type) |
| subqueries                | ✗ (easier to secure, separation of concerns by query type) |


| SQL Functions             | Supported                   |
|---------------------------|-----------------------------|
| min                       | ✓ |
| max                       | ✓ |
| sum                       | ✓ |
| avg                       | ✓ |
| count                     | ✓ |
| distinct                  | ✓ |
More to follow later...

| SQL Cast Types            | Supported                   |
|---------------------------|-----------------------------|
| varchar                   | ✓ |
| text                      | ✓ |
| integer                   | ✓ |
| numeric                   | ✓ |
| float                     | ✓ |
| money                     | ✓ |
| boolean                   | ✓ |
More to follow later...


| SQL Operators             | Supported                   |
|---------------------------|-----------------------------|
| >                         | ✓ |
| >=                        | ✓ |
| <                         | ✓ |
| <=                        | ✓ |
| != or <>                  | ✓ |
| exists                    | ✓ |
| not exists                | ✓ |
| in                        | ✓ |
| not in                    | ✓ |
| is null                   | ✓ |
| is not null               | ✓ |

| SQL Conditionals          | Supported                   |
|---------------------------|-----------------------------|
| case                      | ✓ |
| coalesce                  | ✓ |
| nullif                    | ✓ |


## Tech stack considering:
| Component | Reason | URLs |
|-----------|--------|------|
| Go lang   | Type safety, performance and relative simplicity | https://go.dev |
| Postgres | Database with relatively good security profile | https://www.postgresql.org/ |
| pq | Well established postgres driver, support for extracting arbitrary query results | https://github.com/lib/pq |



## "Lexers" / "Tokenisers"
| Component | Reason | URLs |
|-----------|--------|------|
| bzick/tokenizer | In Golang and appears maintained | https://github.com/bzick/tokenizer |
| lexmachine | In Golang and appears maintained | https://github.com/timtadh/lexmachine https://blog.gopheracademy.com/advent-2017/lexmachine-advent/ |
| go-lexer | Not very active but could be simple? Compatibility with Goyacc.| https://github.com/bbuck/go-lexer |
| blynn/nex | Compatibility with Goyacc. | https://github.com/blynn/nex http://www-cs-students.stanford.edu/~blynn//nex/ |


 

## "Parser" / "Syntactic Analysis" generators 
Using a "generating parser" is likely easier than doing a parser manually!
| Component | Reason | URLs |
|-----------|--------|------|
| Goyacc    | Implements parsers in Go, uses industry standard yacc tool     |  Compatible license. Industry standard parser | https://pkg.go.dev/golang.org/x/tools/cmd/goyacc

## What compiler phases must we support?
### "Lexers" / "Tokenisers"
These recognise the language keywords, operators, expressions etc in the input.
This is essential cleaning and preparation of the input into the "Parser" / "Syntactic Analysis" next.

> For instance with the input `2 + 22` statement, it would recognise:
> - the `+` operator in the expression and it's position in the source code
> - the first `2` expression and it's position in the source code
> - the second `2` expression and it's position in the source code

- The output: for the first "`2`" in `2 + 22`

| Loc                       | Value |
|---------------------------|-------|
| start: Line 1, Column 1   | "2"   |
| end:   Line 1, Column 2   |       |

- The output: for the "`+`" in `2 + 22`

| Loc                       | Value |
|---------------------------|-------|
| start: Line 1, Column 3   | "+"   |
| end:   Line 1, Column 4   |       |

- The output: for the "`22`" in `2 + 22`

| Loc                       | Value |
|---------------------------|-------|
| start: Line 1, Column 5   | "22"   |
| end:   Line 1, Column 7   |       |



### "Parser" / "Syntactic Analysis"
This understands the language's grammar can tell us about: 
- If the lexer output ("lexemes") are valid according to our grammar, eg `2`, `+`, `22`
- The meaning in our language of each type of valid lexer output ("lexeme").

> Running `2 + 22` in https://astexplorer.net/ with the Go language:

- In our example `2` is a `"BasicLit"` or basic literal, of Kind `"INT"`, and the source code location information.
- In our example `22` is a `"BasicLit"` or basic literal, of Kind `"INT"`, and the source code location information.
- In our example `+` is a `"BinaryExpr"` or binary expression, with the `Op: "+"` or plus operator, it then lists the two `BasicLit` for `2` and `22` nested underneath, this is how the tree forms, by this nesting, in an abstract syntax tree.
- `BinaryExpr(Op:"+", X: BasicLit, Y: BasicLit)`


This represents the relationship between the different parts of the input, and is called an "abstract syntax tree" (AST). There are some really good examples from: https://astexplorer.net/, check the Go lang example, click on the code on the left and it will show you the entry in the AST on the right, pretty amazing!

### Further stages may be not required (yet)
Once we have the AST let's evaluate if we can perform the checks we need to:
- Allowed functions
- Allowed operators
- Allowed columns
- Allowed tables

To simplify things initially we could require that data is referred to like `table_name.column_name` so we can easily verify if the column is allowed from that table for that user. 
Future work could be to do some more complex work to deduce `table_name`, and throw errors when it is ambiguous. We know the schema of the tables in the API, which enables this sort of check.