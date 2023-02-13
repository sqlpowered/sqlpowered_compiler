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

| SQL Features in insert    | Supported                   |
|---------------------------|-----------------------------|
| insert                    | ✓* (where user permissions allow) |
| values                    | ✓* (mandatory, cannot be omitted, cannot be empty) |



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

