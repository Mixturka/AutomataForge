#include <stdio.h>
#include <stdlib.h>
#include "../include/parser.h"
#include "../include/utilities/vector.h"

int precedence(TokenType operator) {
    switch (operator) {
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

Vector* parse_infix_to_rpn(Vector* tokens) {
    Vector* op_stack = vector_create(NULL);
    Vector* result = vector_create((void(*)(void*))TOKEN_FREE);

    for (size_t i = 0; i < tokens->size; ++i) {
        Token* cur_token = (Token*)vector_get(tokens, i);

        if (cur_token == NULL) {
            vector_free(op_stack);
            vector_free(result);
            return NULL;
        }

        switch (cur_token->type) {
            case T_CHAR:
                vector_push(result, cur_token);
                break;

            case T_LPAREN:
                vector_push(op_stack, cur_token);
                break;

            case T_RPAREN:
                while (op_stack->size > 0 && ((Token*)vector_back(op_stack))->type != T_LPAREN) {
                    vector_push(result, vector_back(op_stack));
                    vector_pop(op_stack);
                }
                vector_pop(op_stack);
                break;
            case T_STAR:
            case T_PLUS:
            case T_QUESTION:
            case T_UNION:
                while (op_stack->size > 0 && precedence(((Token*)vector_back(op_stack))->type) >= precedence(cur_token->type)) {
                    vector_push(result, vector_back(op_stack));
                    vector_pop(op_stack);
                }
                vector_push(op_stack, cur_token);
                break;

            default:
                break;
        }
    }

    while (op_stack->size > 0) {
        vector_push(result, vector_back(op_stack));
        vector_pop(op_stack);
    }

   vector_free(op_stack);

   return result;
}
