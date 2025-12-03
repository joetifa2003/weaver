// Weaver Language ANTLR4 Grammar
grammar Weaver;

// ==================== Parser Rules ====================

program
    : statement* EOF
    ;

// ==================== Statements ====================

statement
    : varDeclStmt
    | blockStmt
    | whileStmt
    | forStmt
    | ifStmt
    | matchStmt
    | labelStmt
    | gotoStmt
    | continueStmt
    | breakStmt
    | exprStmt
    ;

varDeclStmt
    : IDENT VARDECL expr SEMI?
    ;

blockStmt
    : LBRACE statement* RBRACE
    ;

whileStmt
    : 'while' LPAREN expr RPAREN blockStmt
    ;

forStmt
    : 'for' LPAREN statement expr SEMI expr RPAREN blockStmt    # ForClassic
    | 'for' LPAREN IDENT 'in' expr DOTDOT expr RPAREN statement # ForRange
    ;

ifStmt
    : 'if' LPAREN expr RPAREN blockStmt elseClause?
    ;

elseClause
    : 'else' (blockStmt | ifStmt)
    ;

matchStmt
    : 'match' expr LBRACE matchCaseList? RBRACE
    ;

matchCaseList
    : matchCase (COMMA matchCase)* COMMA?
    ;

matchCase
    : matchCondition ('if' expr)? ARROW statement
    ;

matchCondition
    : matchConditionBase (PIPE matchConditionBase)*
    ;

matchConditionBase
    : matchCaseTypeError
    | matchCaseTypeNumber
    | matchCaseTypeString
    | matchRangeCondition
    | matchCaseInt
    | matchCaseFloat
    | matchCaseString
    | matchCaseArray
    | matchCaseObject
    | matchCaseIdent
    ;

matchCaseTypeError
    : 'error' LPAREN RPAREN                                          # ErrorEmpty
    | 'error' LPAREN matchCondition RPAREN                           # ErrorMessage
    | 'error' LPAREN matchCondition COMMA matchCondition RPAREN      # ErrorMessageData
    ;

matchCaseTypeNumber
    : 'number' LPAREN matchCondition? RPAREN
    ;

matchCaseTypeString
    : 'string' LPAREN matchCondition? RPAREN
    ;

matchRangeCondition
    : (INT | FLOAT) DOTDOT (INT | FLOAT)
    ;

matchCaseInt
    : INT
    ;

matchCaseFloat
    : FLOAT
    ;

matchCaseString
    : STRING
    ;

matchCaseArray
    : LBRACKET (matchCondition (COMMA matchCondition)*)? RBRACKET
    ;

matchCaseObject
    : LBRACE (matchObjectKV (COMMA matchObjectKV)*)? RBRACE
    ;

matchObjectKV
    : IDENT COLON matchCondition   # MatchObjectKVFull
    | IDENT                        # MatchObjectKVShort
    ;

matchCaseIdent
    : IDENT
    ;

labelStmt
    : 'label' IDENT COLON
    ;

gotoStmt
    : 'goto' IDENT
    ;

continueStmt
    : 'continue'
    ;

breakStmt
    : 'break'
    ;

exprStmt
    : expr SEMI?
    ;

// ==================== Expressions ====================

expr
    : returnExpr
    ;

returnExpr
    : 'return' raiseExpr?
    | raiseExpr
    ;

raiseExpr
    : 'raise' tryExpr
    | tryExpr
    ;

tryExpr
    : 'try' ternaryExpr
    | ternaryExpr
    ;

ternaryExpr
    : orExpr (QUESTION expr PIPE expr)?
    ;

orExpr
    : andExpr (OR andExpr)*
    ;

andExpr
    : pipeExpr (AND pipeExpr)*
    ;

pipeExpr
    : equalityExpr (PIPEOP equalityExpr)*
    ;

equalityExpr
    : nequalExpr (EQUAL nequalExpr)*
    ;

nequalExpr
    : lessThanExpr (NOTEQUAL lessThanExpr)*
    ;

lessThanExpr
    : lessThanEqualExpr (LT lessThanEqualExpr)*
    ;

lessThanEqualExpr
    : greaterThanEqualExpr (LTE greaterThanEqualExpr)*
    ;

greaterThanEqualExpr
    : greaterThanExpr (GTE greaterThanExpr)*
    ;

greaterThanExpr
    : addExpr (GT addExpr)*
    ;

addExpr
    : subExpr (PLUS subExpr)*
    ;

subExpr
    : modExpr (MINUS modExpr)*
    ;

modExpr
    : mulExpr (MOD mulExpr)*
    ;

mulExpr
    : divExpr (MUL divExpr)*
    ;

divExpr
    : unaryExpr (DIV unaryExpr)*
    ;

unaryExpr
    : NOT unaryExpr
    | MINUS unaryExpr
    | assignExpr
    ;

assignExpr
    : postFixExpr ASSIGN expr
    | postFixExpr
    ;

postFixExpr
    : incrementExpr postFixOp*
    ;

postFixOp
    : LBRACKET expr RBRACKET     # IndexOp
    | LPAREN exprList? RPAREN trailingBlock?   # CallOp
    | DOT IDENT                  # DotOp
    ;

trailingBlock
    : LBRACE statement* RBRACE
    ;

incrementExpr
    : IDENT INCREMENT    # PostIncrement
    | IDENT DECREMENT    # PostDecrement
    | atom               # AtomExpr
    ;

atom
    : INT                   # IntLiteral
    | FLOAT                 # FloatLiteral
    | STRING                # StringLiteral
    | 'true'                # TrueLiteral
    | 'false'               # FalseLiteral
    | 'nil'                 # NilLiteral
    | IDENT                 # IdentExpr
    | paramList blockStmt   # FunctionExpr
    | paramList expr        # LambdaExpr
    | arrayExpr             # ArrayExprAtom
    | objectExpr            # ObjectExprAtom
    | LPAREN expr RPAREN    # ParenExpr
    ;

paramList
    : PIPE (IDENT (COMMA IDENT)*)? PIPE
    ;

arrayExpr
    : LBRACKET exprList? RBRACKET
    ;

objectExpr
    : LBRACE objectKVList? RBRACE
    ;

objectKVList
    : objectKV (COMMA objectKV)*
    ;

objectKV
    : (IDENT | STRING) COLON expr
    ;

exprList
    : expr (COMMA expr)*
    ;

// ==================== Lexer Rules ====================

// Keywords (must be defined before IDENT to have higher priority)
// Note: ANTLR handles keyword vs identifier automatically by matching longest and earlier rule

// Operators
VARDECL     : ':=' ;
ARROW       : '=>' ;
DOTDOT      : '..' ;
PIPEOP      : '|>' ;
INCREMENT   : '++' ;
DECREMENT   : '--' ;
AND         : '&&' ;
OR          : '||' ;
EQUAL       : '==' ;
NOTEQUAL    : '!=' ;
LTE         : '<=' ;
GTE         : '>=' ;
LT          : '<' ;
GT          : '>' ;
ASSIGN      : '=' ;
PLUS        : '+' ;
MINUS       : '-' ;
MUL         : '*' ;
DIV         : '/' ;
MOD         : '%' ;
NOT         : '!' ;
QUESTION    : '?' ;
PIPE        : '|' ;
DOT         : '.' ;
COMMA       : ',' ;
SEMI        : ';' ;
COLON       : ':' ;
LPAREN      : '(' ;
RPAREN      : ')' ;
LBRACE      : '{' ;
RBRACE      : '}' ;
LBRACKET    : '[' ;
RBRACKET    : ']' ;

// Literals
FLOAT       : [0-9]+ '.' [0-9]+ ;
INT         : [0-9]+ ;
STRING      : '"' (~["\\\r\n] | '\\' .)* '"' ;

// Identifiers - matches the original lexer pattern: letters/underscores followed by optional digits
IDENT       : [a-zA-Z_]+ [0-9]* ;

// Whitespace and Comments
WS          : [ \t\r\n]+ -> skip ;
COMMENT     : '#' ~[\r\n]* -> skip ;
