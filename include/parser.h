#ifndef PARSER_H
#define PARSER_H

#include "./lexer/token.h"

typedef enum {
    PARSER_ERR_NONE = 0,
} ParserErrorCode;

typedef struct {
    ParserErrorCode code;
    const char* message;
    const char* input;
    size_t pos;
} ParserError;

Token* parse_infix_to_rpn(Token* tokens, size_t num_tokens); // rpn - reverse polish notation

#endif // PARSER_H