#ifndef LEXER_H
#define LEXER_H

#include "token.h"

#define SET_LEXER_ERROR(error, error_code, msg, _pos) \
    do { \
        if (error) { \
            error->code = error_code; \
            error->message = msg; \
            error->pos = _pos; \
        } \
    } while(0)

typedef enum {
    LEXER_ERR_NONE = 0,
    LEXER_ERR_MEMORY,
    LEXER_ERR_INVALID_CHAR,
    LEXER_ERR_INVALID_ESCAPE,
    LEXER_ERR_UNEXPECTED_END,
    LEXER_ERR_UNKNOWN
} LexerErrorCode;

typedef struct {
    LexerErrorCode code;
    const char* message;
    const char* input;
    size_t pos;
} LexerError;

Token* tokenize_regex(const char* regex_src, size_t regex_len, LexerError* error);

void gen_lexer_error_message(const LexerError* error, char* output_buffer, size_t buffer_size);

#endif // LEXER_H