#ifndef TOKEN_H
#define TOKEN_H

#include <string.h>
#include <stdlib.h>

typedef enum {
    T_CHAR,
    T_UNION,
    T_STAR,
    T_PLUS,
    T_QUESTION,
    T_LPAREN,
    T_RPAREN,
    T_ESCAPE,
    T_END,
} TokenType;

typedef struct {
    TokenType type;
    size_t len;
    char* value;
} Token;

// Dynamically allocates Token instance
static inline Token* TOKEN_SINGLE(TokenType type, char c) {
    Token* token = malloc(sizeof(Token));
    token->type = type;
    token->len = 1;
    token->value = malloc(2);

    if (token->value) {
        token->value[0] = c;
        token->value[1] = '\0';
    }

    return token;
}

// Dynamically allocates Token instance
static inline Token* TOKEN_MULTI(TokenType type, const char* str, size_t len) {
    Token* token = malloc(sizeof(Token));
    token->type = type;
    token->len = len;
    token->value = malloc(len + 1);

    if (token->value) {
        strncpy(token->value, str, len);
        token->value[len] = '\0';
    }

    return token;
}
/* Frees ONLY dynamically allocated tokens
*  such as ones that were created via TOKEN_SINGLE / TOKEN_MULTI
*/
static inline void TOKEN_FREE(Token* token) {
    if (token) {
        if (token->value) {
            free(token->value);
            token->value = NULL;
        }
        free(token);
    }
}

#endif // TOKEN_H
