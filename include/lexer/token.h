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
    T_EPSILON,
    T_CONCAT,
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

static inline Token* TOKEN_COPY(Token* token) {
    if (!token) return NULL;

    Token* copy = malloc(sizeof(Token));
    if (!copy) return NULL;

    copy->type = token->type;
    copy->len = token->len;

    if (token->value) {
        copy->value = (char*)malloc(token->len + 1);
        if (!copy->value) {
            free(copy);
            return NULL;
        }
        strncpy(copy->value, token->value, token->len);
        copy->value[token->len] = '\0';
    } else {
        copy->value = NULL;
    }

    return copy;
}

#endif // TOKEN_H
