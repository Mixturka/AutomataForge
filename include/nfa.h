#ifndef NFA_H
#define NFA_H

#include "lexer/token.h"
#include "utilities/vector.h"

typedef struct {
    int id;
    Vector* transitions;
    size_t ref_count;
} NFAState;

typedef struct {
    Token* token;
    NFAState* next_state;
    size_t ref_count;
} NFATransition;

typedef struct {
    NFAState* start;
    NFAState* accept;
    size_t ref_count;
} NFA;

static int state_id_counter = 0;

void add_transition(NFAState* state, NFAState* to_state, Token* token);
NFA* nfa_union(NFA* nfa1, NFA* nfa2);
NFA* nfa_concatenate(NFA* nfa1, NFA* nfa2);
NFA* nfa_closure(NFA* nfa);
void nfa_print(NFA* nfa);
NFA* nfa_build(Vector* rpn_tokens);
static inline void NFA_STATE_FREE(NFAState* state);
static inline void NFA_TRANSITION_FREE(NFATransition* transition);

static inline NFAState* NFA_STATE() {
    NFAState* state = (NFAState*)malloc(sizeof(NFAState));
    if (!state) {
        return NULL;
    }

    state->id = state_id_counter++;
    state->transitions = vector_create((void(*)(void*))NFA_TRANSITION_FREE);
    state->ref_count = 1;

    return state;
}

static inline NFATransition* NFA_TRANSITION(Token* token, NFAState* next_state) {
    NFATransition* transition = (NFATransition*)malloc(sizeof(NFATransition));
    if (!transition) {
        return NULL;
    }
    
    transition->token = TOKEN_COPY(token); // Creating a copy of token to make sure token will be freed once
    transition->next_state = next_state;
    transition->ref_count = 1;
    
    return transition;
}

static inline void NFA_TRANSITION_FREE(NFATransition* transition) {
    if (transition) {
        if (--transition->ref_count == 0) {
            TOKEN_FREE(transition->token);
            NFA_STATE_FREE(transition->next_state);
            free(transition);
        }
    }
}

static inline void NFA_STATE_FREE(NFAState* state) {
    if (state) {
        if (--state->ref_count == 0) {
            if (state->transitions) {
                for (size_t i = 0; i < state->transitions->size; ++i) {
                    NFATransition* transition = (NFATransition*)vector_get(state->transitions, i);
                    NFA_TRANSITION_FREE(transition);
                }
                vector_free(state->transitions);
            }

            free(state);
        }
    }
}

static inline void NFA_FREE(NFA* nfa) {
    if (nfa) {
        if (--nfa->ref_count == 0) {
            NFA_STATE_FREE(nfa->start);
            NFA_STATE_FREE(nfa->accept);
            free(nfa);
            nfa = NULL;
        }
    }
}

static inline NFA* NFA_CREATE() {
    NFA* nfa = (NFA*)malloc(sizeof(NFA));
    if (!nfa) {
        return NULL;
    }
    
    nfa->start = NFA_STATE();
    nfa->accept = NFA_STATE();
    nfa->ref_count = 1;
    
    return nfa;
}

static inline NFA* NFA_BASIC(Token* token) {
    NFA* nfa = (NFA*)malloc(sizeof(NFA));
    nfa->start = NFA_STATE();
    nfa->accept = NFA_STATE();
    nfa->ref_count = 1;
    add_transition(nfa->start, nfa->accept, token);
    return nfa;
}

#endif // NFA_H
