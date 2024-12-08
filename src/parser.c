#include <stdio.h>

#include "../include/parser.h"

#define MAX_STACK 256

int precedence(TokenType operator) {
    switch (operator)
    {
        case T_STAR:
        case T_PLUS:
        case T_QUESTION:
            return 2;
        case T_UNION:
            return 1;
        default:
            return 0;
    }
}

Token* parse_infix_to_rpn(Token* tokens, size_t num_tokens) {
    Token* result = malloc(sizeof(Token) * num_tokens);

    if (!result) {
        fprintf(stderr, "Memory allocation failed for result\n");
        return NULL;
    }

    Token* op_stack = malloc(sizeof(Token) * MAX_STACK); // operator stack
    if (!op_stack) {
        fprintf(stderr, "Memory allocation failed for op_stack\n");
        free(result);
        return NULL;
    }
    int op_pos = 0;
    int output_pos = 0;

    for (size_t i = 0; i < num_tokens; ++i) {
        Token cur_token = tokens[i];

        switch (cur_token.type) {
            case T_CHAR:
                result[output_pos++] = cur_token;
                break;
            case T_LPAREN:
                op_stack[op_pos++] = cur_token;
                break;
            case T_RPAREN:
                while (op_pos > 0 && op_stack[op_pos - 1].type != T_LPAREN) {
                    result[output_pos++] = op_stack[--op_pos];
                }
                --op_pos;
                break;
            case T_STAR:
            case T_PLUS:
            case T_QUESTION:
            case T_UNION:
                while (op_pos > 0 && precedence(op_stack[op_pos - 1].type) >= precedence(cur_token.type)) {
                    if (output_pos < num_tokens) {
                        result[output_pos++] = op_stack[--op_pos];
                    }
                }
                if (op_pos < MAX_STACK) {
                    op_stack[op_pos++] = cur_token;
                }
                break;
            default:
                break;
        }
    }

    while (op_pos > 0) {
        if (output_pos < num_tokens) {
            result[output_pos++] = op_stack[--op_pos];
        }
    }

    free(op_stack);
    return result;
}