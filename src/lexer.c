#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>

#include "../include/lexer/lexer.h"

Token* tokenize_regex(const char* regex_src, size_t regex_len, LexerError* error) {
    if (!regex_src || regex_len == 0) {
        SET_LEXER_ERROR(error, LEXER_ERR_UNKNOWN, "Invalid input: NULL or empty regex", 0);
        return NULL;
    }

    Token* tokens = malloc(sizeof(Token) * (regex_len + 1)); // + 1 for T_END token

    if (!tokens) {
        SET_LEXER_ERROR(error, LEXER_ERR_MEMORY, "Memory allocation failed", 0);
        return NULL;
    }

    size_t tokens_ptr = 0;
    int open_parens = 0;

    for (size_t i = 0; i < regex_len; ++i) {
        char c = regex_src[i];

        if (isalnum(c)) {
            tokens[tokens_ptr++] = TOKEN_SINGLE(T_CHAR, c);
        } else if (c == '*') {
            tokens[tokens_ptr++] = TOKEN_SINGLE(T_STAR, c);
        } else if (c == '+') {
            tokens[tokens_ptr++] = TOKEN_SINGLE(T_PLUS, c);
        } else if (c == '?') {
            tokens[tokens_ptr++] = TOKEN_SINGLE(T_QUESTION, c);
        } else if (c == '(') {
            tokens[tokens_ptr++] = TOKEN_SINGLE(T_LPAREN, c);
            ++open_parens;
        } else if (c == ')') {
            if (open_parens <= 0) {
                SET_LEXER_ERROR(error, LEXER_ERR_INVALID_CHAR, "Unmatched closing parenthesis", i);
                goto cleanup;
            }
            tokens[tokens_ptr++] = TOKEN_SINGLE(T_RPAREN, c);
            --open_parens;
        } else if (c == '|') {
            if (tokens_ptr == 0 || (tokens[tokens_ptr - 1].type == T_UNION || tokens[tokens_ptr - 1].type == T_LPAREN) ||
                (i + 1 < regex_len && regex_src[i + 1] == ')')) {
                SET_LEXER_ERROR(error, LEXER_ERR_INVALID_CHAR, "Invalid alternation operator '|'", i);
                goto cleanup;
            }
            tokens[tokens_ptr++] = TOKEN_SINGLE(T_UNION, c);
        } else if (c == '\\') {
            if (i + 1 >= regex_len) {
                SET_LEXER_ERROR(error, LEXER_ERR_INVALID_ESCAPE, "Invalid escape sequence", i);
                goto cleanup;
            }
            size_t escape_len = 1;

            while (i + escape_len < regex_len && isalnum(regex_src[i + escape_len])) {
                ++escape_len;
            }

            tokens[tokens_ptr++] = TOKEN_MULTI(T_ESCAPE, &regex_src[i + 1], escape_len - 1);
            i += escape_len - 1;
        } else {
            SET_LEXER_ERROR(error, LEXER_ERR_INVALID_CHAR, "Invalid character in regex", i);
            goto cleanup;
        }
    }

    if (open_parens > 0) {
        SET_LEXER_ERROR(error, LEXER_ERR_UNEXPECTED_END, "Unmatched opening parenthesis", regex_len);
        goto cleanup;
    }

    tokens[tokens_ptr++] = TOKEN_SINGLE(T_END, '\0');
    return tokens;

cleanup:
    for (size_t j = 0; j < tokens_ptr; ++j) {
        TOKEN_FREE(&tokens[j]);
    }
    free(tokens);
    return NULL;
}

void gen_lexer_error_message(const LexerError* error, char* output_buffer, size_t buffer_size) {
    memset(output_buffer, 0, buffer_size);

    snprintf(output_buffer, buffer_size, "\033[31mLexing error: %s\033[0m\n", error->message);
    strncat(output_buffer, error->input, buffer_size - strlen(output_buffer) - 1);

    if (error->pos < strlen(error->input)) {
        for (size_t i = 0; i < error->pos; ++i) {
            strncat(output_buffer, " ", buffer_size - strlen(output_buffer) - 1);
        }
        strncat(output_buffer, "^\n", buffer_size - strlen(output_buffer) - 1);
    }
}