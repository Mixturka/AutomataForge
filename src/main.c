#include <stdio.h>

#include "../include/lexer/lexer.h"
#include "../include/parser.h"

int main() {
    const char* regex = "a(a|b)";
    LexerError error;
    Token* tokens = tokenize_regex(regex, strlen(regex), &error);

    if (tokens) {
        printf("First token: %s\n", tokens[0].value);
    } else {
        char error_message[256];
        // gen_lexer_error_message(&error, error_message, sizeof(error_message));
        printf("%s", error.message);
    }

    Token* rpn_tokens = parse_infix_to_rpn(tokens, strlen(regex) + 1);

    if (rpn_tokens) {
        for (size_t i = 0; i < strlen(regex) + 1; ++i) {
            printf("%s ", rpn_tokens[i].value);
        }
    }

    if (tokens) {
        for (int i = 0; i < strlen(regex) + 1; ++i) {
            TOKEN_FREE(&tokens[i]);
        }
        free(tokens);
    }

    // if (rpn_tokens) {
    //     for (int i = 0; i < strlen(regex) + 1; ++i) {
    //         TOKEN_FREE(&rpn_tokens[i]);
    //     }
    //     free(rpn_tokens);
    // }
}   