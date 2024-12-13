#ifndef PARSER_H
#define PARSER_H

#include "./lexer/token.h"
#include "./utilities/vector.h"

#define SET_PARSER_ERROR(error, error_code) \
    do { \
        if (error) { \
            error->code = error_code; \
        } \
    } while(0)

typedef enum {
    PARSER_ERR_NONE = 0,
    PARSER_ERR_MEMORY
} ParserErrorCode;

typedef struct {
    ParserErrorCode code;
} ParserError;

Vector* parse_infix_to_rpn(Vector* tokens); // rpn - reverse polish notation

#endif // PARSER_H
