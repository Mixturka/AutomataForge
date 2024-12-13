#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>

#include "../include/lexer/lexer.h"
#include "../include/utilities/vector.h"

Vector* tokenize_regex(const char* regex_src, size_t regex_len, LexerError* error) {
    if (!regex_src || regex_len == 0) {
        SET_LEXER_ERROR(error, LEXER_ERR_UNKNOWN, regex_src, "Invalid input: NULL or empty regex", 0);
        return NULL;
    }

    Vector* tokens = vector_create((void(*)(void*))TOKEN_FREE);
    int open_parens = 0;

    for (size_t i = 0; i < regex_len; ++i) {
        char c = regex_src[i];

        if (isalnum(c)) {
            Token* token = TOKEN_SINGLE(T_CHAR, c);
            vector_push(tokens, token);
        } else if (c == '*') {
            Token* token = TOKEN_SINGLE(T_STAR, c);
            vector_push(tokens, token);
        } else if (c == '+') {
            Token* token = TOKEN_SINGLE(T_PLUS, c);
            vector_push(tokens, token);
        } else if (c == '?') {
            Token* token = TOKEN_SINGLE(T_QUESTION, c);
            vector_push(tokens, token);
        } else if (c == '(') {
            Token* token = TOKEN_SINGLE(T_LPAREN, c);
            vector_push(tokens, token);
            ++open_parens;
        } else if (c == ')') {
            if (open_parens <= 0) {
                SET_LEXER_ERROR(error, LEXER_ERR_INVALID_CHAR, regex_src, "Unmatched closing parenthesis", i);
                goto cleanup;
            }
            Token* token = TOKEN_SINGLE(T_RPAREN, c);
            vector_push(tokens, token);
            --open_parens;
        } else if (c == '|') {
            if (tokens->size == 0 || (((Token*)vector_back(tokens))->type == T_UNION || ((Token*)vector_back(tokens))->type == T_LPAREN) ||
                (i + 1 < regex_len && regex_src[i + 1] == ')')) {
                SET_LEXER_ERROR(error, LEXER_ERR_INVALID_CHAR, regex_src, "Invalid alternation operator '|'", i);
                goto cleanup;
            }
            Token* token = TOKEN_SINGLE(T_UNION, c);
            vector_push(tokens, token);
        } else if (c == '\\') {
            if (i + 1 >= regex_len) {
                SET_LEXER_ERROR(error, LEXER_ERR_INVALID_ESCAPE, regex_src, "Invalid escape sequence", i);
                goto cleanup;
            }
            size_t escape_len = 1;

            while (i + escape_len < regex_len && isalnum(regex_src[i + escape_len])) {
                ++escape_len;
            }

            Token* token = TOKEN_MULTI(T_ESCAPE, &regex_src[i + 1], escape_len - 1);
            vector_push(tokens, token);
            i += escape_len - 1;
        } else {
            SET_LEXER_ERROR(error, LEXER_ERR_INVALID_CHAR, regex_src, "Invalid character in regex", i);
            goto cleanup;
        }
    }

    if (open_parens > 0) {
        SET_LEXER_ERROR(error, LEXER_ERR_UNEXPECTED_END, regex_src, "Unmatched opening parenthesis", regex_len);
        goto cleanup;
    }

    Token* end_token = TOKEN_SINGLE(T_END, '\0');
    vector_push(tokens, end_token);

    return tokens;

cleanup:
    vector_free(tokens);
    return NULL;
}

void gen_lexer_error_message(const LexerError* error, char* output_buffer, size_t buffer_size) {
    if (error == NULL) {
        snprintf(output_buffer, buffer_size, "No lexer error struct provided.\n");
        return;
    }

    memset(output_buffer, 0, buffer_size);
    snprintf(output_buffer, buffer_size, "\033[31mLexing error: %s\033[0m\n", error->message);
    size_t remaining_space = buffer_size - strlen(output_buffer) - 1;

    if (remaining_space > 0 && error->input != NULL) {
        strncat(output_buffer, error->input, remaining_space);
        strncat(output_buffer, "\n", remaining_space);
    }

    if (error->input != NULL && error->pos <= strlen(error->input)) {
        remaining_space = buffer_size - strlen(output_buffer) - 1;

        if (remaining_space > 0) {
            for (size_t i = 0; i < error->pos && remaining_space > 0; ++i) {
                strncat(output_buffer, " ", remaining_space);
                remaining_space--;
            }
            if (remaining_space > 0) {
                strncat(output_buffer, "^\n", remaining_space);
            }
        }
    }
}

