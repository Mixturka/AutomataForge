package nfa

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Mixturka/AutomataForge/internal/dfa"
	"github.com/Mixturka/AutomataForge/internal/lexer/token"
	"github.com/Mixturka/AutomataForge/pkg/queue"
	"github.com/Mixturka/AutomataForge/pkg/sliceutils"
	"github.com/Mixturka/AutomataForge/pkg/stack"
	"slices"
)

const Epsilon = 'ε'

var curId = 0

type NFA struct {
	transitions map[int]map[rune][]int
	Start       int
	Accept      int
}

func NewNFA(t token.Token) *NFA {
	start := curId
	curId++
	accept := curId
	curId++
	nfa := &NFA{
		transitions: make(map[int]map[rune][]int),
		Start:       start,
		Accept:      accept,
	}
	nfa.AddTransition(start, accept, t)
	return nfa
}

func (nfa *NFA) Concatenate(other *NFA) {
	nfa.AddTransition(nfa.Accept, other.Start, token.NewToken(token.Eps, Epsilon, 0))
	nfa.mergeTransitions(other)
	nfa.Accept = other.Accept
}

func (nfa *NFA) Alterate(other *NFA) {
	newStart := curId
	curId++
	newAccept := curId
	curId++
	nfa.AddTransition(newStart, nfa.Start, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(nfa.Accept, newAccept, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(newStart, other.Start, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(other.Accept, newAccept, token.NewToken(token.Eps, Epsilon, 0))
	nfa.mergeTransitions(other)
	nfa.Start = newStart
	nfa.Accept = newAccept
}

func (nfa *NFA) Closure() {
	newStart := curId
	curId++
	newAccept := curId
	curId++
	nfa.AddTransition(newStart, nfa.Start, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(nfa.Accept, nfa.Start, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(newStart, newAccept, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(nfa.Accept, newAccept, token.NewToken(token.Eps, Epsilon, 0))
	nfa.Start = newStart
	nfa.Accept = newAccept
}

func (nfa *NFA) Optional() {
	newStart := curId
	curId++
	newAccept := curId
	curId++
	nfa.AddTransition(newStart, newAccept, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(newStart, nfa.Start, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(nfa.Accept, newAccept, token.NewToken(token.Eps, Epsilon, 0))
	nfa.Start = newStart
	nfa.Accept = newAccept
}

func (nfa *NFA) Plus() {
	newStart := curId
	curId++
	newAccept := curId
	curId++
	nfa.AddTransition(newStart, nfa.Start, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(nfa.Accept, nfa.Start, token.NewToken(token.Eps, Epsilon, 0))
	nfa.AddTransition(nfa.Accept, newAccept, token.NewToken(token.Eps, Epsilon, 0))
	nfa.Start = newStart
	nfa.Accept = newAccept
}

func (nfa *NFA) AddTransition(from, to int, t token.Token) {
	if _, exists := nfa.transitions[from]; !exists {
		nfa.transitions[from] = make(map[rune][]int)
	}
	nfa.transitions[from][t.Value] = append(nfa.transitions[from][t.Value], to)
}

func (nfa *NFA) mergeTransitions(other *NFA) {
	for state, transMap := range other.transitions {
		if nfa.transitions[state] == nil {
			nfa.transitions[state] = make(map[rune][]int)
		}
		for symbol, nextStates := range transMap {
			nfa.transitions[state][symbol] = append(nfa.transitions[state][symbol], nextStates...)
		}
	}
}

func (nfa *NFA) GetAlphabet() []rune {
	symbolSet := make(map[rune]struct{})
	alphabet := make([]rune, 0)

	for _, possibleTransitions := range nfa.transitions {
		for key := range possibleTransitions {
			if _, ok := symbolSet[key]; !ok && key != Epsilon {
				alphabet = append(alphabet, key)
				symbolSet[key] = struct{}{}
			}
		}
	}

	return alphabet
}

// Returns NFA from given token slice in Reverse Polish Notation (RPN)
func BuildNFA(rpn []token.Token) *NFA {
	stack := stack.NewStack[*NFA]()
	for _, tok := range rpn {
		if tok.Type == token.Char {
			stack.Push(NewNFA(tok))
		} else if token.IsOperator(tok) {
			switch tok.Type {
			case token.Concat:
				nfa1 := stack.Pop()
				nfa2 := stack.Pop()
				nfa2.Concatenate(nfa1)
				stack.Push(nfa2)
			case token.Alter:
				nfa1 := stack.Pop()
				nfa2 := stack.Pop()
				nfa2.Alterate(nfa1)
				stack.Push(nfa2)
			case token.Star:
				nfa := stack.Pop()
				nfa.Closure()
				stack.Push(nfa)
			case token.Optional:
				nfa := stack.Pop()
				nfa.Optional()
				stack.Push(nfa)
			case token.Plus:
				nfa := stack.Pop()
				nfa.Plus()
				stack.Push(nfa)
			}
		}
	}

	return stack.Top()
}

// Performs subset construction algorithm to build up a DFA from given NFA
func (nfa *NFA) BuildDFA() *dfa.DFA {
	curId = 0
	alphabet := nfa.GetAlphabet()
	dfaTransitions := make(map[int]map[rune]int)
	var dfaAccepts []int

	q0 := nfa.computeEpsilonClosure([]int{nfa.Start})
	Q := [][]int{q0}
	worklist := queue.NewQueue[[]int]()
	worklist.Push(q0)
	stateMap := map[string]int{sliceutils.HashSlice(q0): curId}
	curId++

	for worklist.Length() > 0 {
		q := worklist.Pop()
		sort.Ints(q)
		dfaState := stateMap[sliceutils.HashSlice(q)]
		dfaTransitions[dfaState] = make(map[rune]int)

		if nfa.isAcceptState(q) {
			dfaAccepts = append(dfaAccepts, dfaState)
		}

		for _, symbol := range alphabet {
			reachableStates := nfa.computeEpsilonClosure(nfa.Delta(q, symbol))

			if len(reachableStates) > 0 {
				sort.Ints(reachableStates)
				hashedKey := sliceutils.HashSlice(reachableStates)
				if _, exists := stateMap[hashedKey]; !exists {
					stateMap[hashedKey] = curId
					curId++
					Q = append(Q, reachableStates)
					worklist.Push(reachableStates)
				}

				dfaTransitions[dfaState][symbol] = stateMap[hashedKey]
			}
		}
	}

	return &dfa.DFA{
		Transitions: dfaTransitions,
		Start:       0,
		Accepts:     dfaAccepts,
	}
}

func (nfa *NFA) computeEpsilonClosure(states []int) []int {
	visited := make(map[int]struct{})
	result := make([]int, 0, len(states))

	for _, state := range states {
		if _, exists := visited[state]; !exists {
			visited[state] = struct{}{}
			result = append(result, state)
		}
	}
	for i := range result {
		state := result[i]
		if transitionMap, ok := nfa.transitions[state]; ok {
			if nextStates, ok := transitionMap[Epsilon]; ok {
				for _, nextState := range nextStates {
					if _, exists := visited[nextState]; !exists {
						visited[nextState] = struct{}{}
						result = append(result, nextState)
					}
				}
			}
		}
	}

	return result
}

// NFA transition function. Returns all reachable states if transition by symbol from given states
func (nfa *NFA) Delta(states []int, symbol rune) []int {
	var result []int

	for _, state := range states {
		if possibleTransitions, ok := nfa.transitions[state]; ok {
			if reachableStates, ok := possibleTransitions[symbol]; ok {
				result = append(result, reachableStates...)
			}
		}
	}

	return result
}

// Checks if given state slice contains accept state or not
func (nfa *NFA) isAcceptState(states []int) bool {
	return slices.Contains(states, nfa.Accept)
}

func (nfa *NFA) PrettyPrint() string {
	visited := make(map[int]bool)
	var builder strings.Builder

	builder.WriteString("NFA:\n")
	builder.WriteString(fmt.Sprintf("Start State: %d\n", nfa.Start))
	builder.WriteString(fmt.Sprintf("Accept State: %d\n", nfa.Accept))
	builder.WriteString("Transitions:\n")

	nfa.prettyPrintState(&builder, nfa.Start, visited)
	return builder.String()
}

func (nfa *NFA) prettyPrintState(builder *strings.Builder, state int, visited map[int]bool) {
	if visited[state] {
		return
	}
	visited[state] = true

	for symbol, nextStates := range nfa.transitions[state] {
		displaySymbol := symbol

		for _, nextState := range nextStates {
			builder.WriteString(fmt.Sprintf("  %d --[%c]--> %d\n", state, displaySymbol, nextState))
			nfa.prettyPrintState(builder, nextState, visited)
		}
	}
}
