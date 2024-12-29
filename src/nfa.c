#include "../include/nfa.h"

void add_transition(NFAState* state, NFAState* to_state, Token* token) {
    if (!state || !to_state) return;

    NFATransition* transition = NFA_TRANSITION(token, to_state);
    if (!transition) return;

    vector_push(state->transitions, transition);
    to_state->ref_count++;
}

NFA* nfa_union(NFA* nfa1, NFA* nfa2) {
    NFA* result = NFA_CREATE();
    
    add_transition(result->start, nfa1->start, TOKEN_SINGLE(T_EPSILON, '\0'));
    add_transition(nfa1->accept, result->accept, TOKEN_SINGLE(T_EPSILON, '\0'));

    add_transition(result->start, nfa2->start, TOKEN_SINGLE(T_EPSILON, '\0'));
    add_transition(nfa2->accept, result->accept, TOKEN_SINGLE(T_EPSILON, '\0'));

    return result;
}

NFA* nfa_concatenate(NFA* nfa1, NFA* nfa2) {
    add_transition(nfa1->accept, nfa2->start, TOKEN_SINGLE(T_EPSILON, '\0'));
    NFA* result = malloc(sizeof(NFA));
    result->ref_count++;

    result->start = nfa1->start;
    result->accept = nfa2->accept;

    return result;
}

NFA* nfa_closure(NFA* nfa) {
    NFA* result = malloc(sizeof(NFA));
    result->ref_count++;
    result->start = NFA_STATE();
    result->accept = NFA_STATE();

    add_transition(result->start, nfa->start, TOKEN_SINGLE(T_EPSILON, '\0'));
    add_transition(result->start, result->accept, TOKEN_SINGLE(T_EPSILON, '\0'));
    add_transition(nfa->accept, nfa->start, TOKEN_SINGLE(T_EPSILON, '\0'));
    add_transition(nfa->accept, result->accept, TOKEN_SINGLE(T_EPSILON, '\0'));

    return result;
}

NFA* nfa_build(Vector* rpn_tokens) {
    Vector* stack = vector_create((void(*)(void*))NFA_FREE);

    for (size_t i = 0; i < rpn_tokens->size; ++i) {
        Token* token = (Token*)vector_get(rpn_tokens, i);
        printf("DEBUG: %s\n", token->value);
        switch (token->type) {
            case T_CHAR: {
                NFA* nfa = NFA_BASIC(token);
                vector_push(stack, nfa);
                break;
            }
            case T_CONCAT: {
                NFA* nfa2 = (NFA*)vector_back(stack);
                nfa2->ref_count++;
                vector_pop(stack);
                NFA* nfa1 = (NFA*)vector_back(stack);
                nfa1->ref_count++;
                vector_pop(stack);

                NFA* nfa = nfa_concatenate(nfa1, nfa2);
                vector_push(stack, nfa);
                break;
            }
            case T_UNION: {
                NFA* nfa2 = (NFA*)vector_back(stack);

                nfa2->ref_count++;
                vector_pop(stack);
                NFA* nfa1 = (NFA*)vector_back(stack);
                nfa1->ref_count++;
                vector_pop(stack);
                
                NFA* result = nfa_union(nfa1, nfa2);
                vector_push(stack, result);

                break;
            }
            case T_STAR: {
                NFA* nfa = (NFA*)vector_back(stack);
                nfa->ref_count++;
                vector_pop(stack);
                NFA* result = nfa_closure(nfa);
                vector_push(stack, result);
                break;
            }
            case T_QUESTION:
            case T_PLUS:
            case T_LPAREN:
            case T_RPAREN:
                break;
            case T_EPSILON:
                break;
            default:
                break; // TODO: add logging and exiting   
        }
    }
    NFA* result = (NFA*)vector_back(stack);
    result->ref_count++;
    vector_pop(stack);
    vector_free(stack);

    return result;
}

void nfa_print(NFA* nfa) {
    if (!nfa) {
        printf("NFA is NULL.\n");
        return;
    }

    printf("NFA Start State: %d\n", nfa->start->id);
    printf("NFA Accept State: %d\n", nfa->accept->id);

    size_t max_states = nfa->accept->id + 1;
    int visited_states[max_states];
    memset(visited_states, 0, sizeof(visited_states));

    Vector* queue = vector_create((void(*)(void*))NFA_STATE_FREE);

    vector_push(queue, nfa->start);
    nfa->start->ref_count++;
    visited_states[nfa->start->id] = 1;

    size_t front_index = 0;

    while (front_index < queue->size) {
        NFAState* current_state = (NFAState*)vector_get(queue, front_index);
        front_index++;

        printf("State %d:\n", current_state->id);

        for (size_t i = 0; i < current_state->transitions->size; ++i) {
            NFATransition* transition = (NFATransition*)vector_get(current_state->transitions, i);
            if (!transition || !transition->token || !transition->next_state) {
                printf("  Invalid transition encountered.\n");
                continue;
            }
            printf("  Transition on token '%s' to state %d\n", transition->token->value, transition->next_state->id);

            if (!visited_states[transition->next_state->id]) {
                vector_push(queue, transition->next_state);
                transition->next_state->ref_count++;
                visited_states[transition->next_state->id] = 1;
            }
        }
    }

    vector_free(queue);
}

