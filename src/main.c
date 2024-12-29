#include <stdio.h>

#include "../include/lexer/lexer.h"
#include "../include/parser.h"
#include "../include/utilities/vector.h"
#include "../include/nfa.h"

int main() {
    const char* regex = "a*";
    LexerError error;
    Vector* tokens = tokenize_regex(regex, strlen(regex), &error);

    if (tokens) {
        printf("First token: %s\n", ((Token*)tokens->data[0])->value);
    } else {
        char error_message[256];
        gen_lexer_error_message(&error, error_message, sizeof(error_message));
        printf("%s\n", error_message);
        exit(1);
    }

    for (size_t i = 0; i < tokens->size; ++i) {
     printf("Token %zu: %s\n", i, ((Token*)vector_get(tokens, i))->value);
    }

    Vector* a = parse_infix_to_rpn(tokens);

    printf("%zu\n", a->size);

    for (size_t i = 0; i < a->size; ++i) {
        printf("%s ", ((Token*)vector_get(a, i))->value);
    }

    NFA* nfa = nfa_build(a);

    nfa_print(nfa);

    // if (tokens) {
    //     for (int i = 0; i < strlen(regex) + 1; ++i) {
    //         TOKEN_FREE(&tokens[i]);
    //     }
    //     free(tokens);
    // }

    // if (rpn_tokens) {
    //     for (int i = 0; i < strlen(regex) + 1; ++i) {
    //         TOKEN_FREE(&rpn_tokens[i]);
    //     }
    //     free(rpn_tokens);
    // }
}   