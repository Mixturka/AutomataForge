package nfa

import (
	"fmt"
	"regexp/syntax"
	"sort"
	"strings"

	"github.com/Mixturka/AutomataForge/internal/dfa"
	"github.com/Mixturka/AutomataForge/pkg/queue"
	"github.com/Mixturka/AutomataForge/pkg/sliceutils"
)

const Epsilon = 'ε'

type StateIdGenerator struct {
	nextId int
}

func NewStateIdGenerator() StateIdGenerator {
	return StateIdGenerator{
		nextId: 0,
	}
}

func (sid *StateIdGenerator) NextId() int {
	id := sid.nextId
	sid.nextId++
	return id
}

type BaseNfa struct {
	transitions map[int]map[rune][]int
	stateIdGen  *StateIdGenerator
}

func NewBaseNfa(stateIdGen *StateIdGenerator) *BaseNfa {
	return &BaseNfa{
		transitions: make(map[int]map[rune][]int),
		stateIdGen:  stateIdGen,
	}
}

func (bnfa *BaseNfa) AddTransition(from, to int, symbol rune) {
	if _, exists := bnfa.transitions[from]; !exists {
		bnfa.transitions[from] = make(map[rune][]int)
	}
	bnfa.transitions[from][symbol] = append(bnfa.transitions[from][symbol], to)
}

func (bnfa *BaseNfa) MergeTransitions(other *BaseNfa) {
	for state, transMap := range other.transitions {
		if bnfa.transitions[state] == nil {
			bnfa.transitions[state] = make(map[rune][]int)
		}
		for symbol, nextStates := range transMap {
			bnfa.transitions[state][symbol] = append(bnfa.transitions[state][symbol], nextStates...)
		}
	}
}

type RegexNfa struct {
	*BaseNfa
	Start     int
	Accept    int
	TokenInfo dfa.TokenInfo
}

func (rnfa *RegexNfa) Concatenate(other *RegexNfa) {
	rnfa.AddTransition(rnfa.Accept, other.Start, Epsilon)
	rnfa.MergeTransitions(other.BaseNfa)
	rnfa.Accept = other.Accept
}

func (rnfa *RegexNfa) Alterate(other *RegexNfa) {
	newStart := rnfa.stateIdGen.NextId()
	newAccept := rnfa.stateIdGen.NextId()

	rnfa.AddTransition(newStart, rnfa.Start, Epsilon)
	rnfa.AddTransition(rnfa.Accept, newAccept, Epsilon)
	rnfa.AddTransition(newStart, other.Start, Epsilon)
	rnfa.AddTransition(other.Accept, newAccept, Epsilon)
	rnfa.MergeTransitions(other.BaseNfa)
	rnfa.Start = newStart
	rnfa.Accept = newAccept
}

func (rnfa *RegexNfa) Closure() {
	newStart := rnfa.stateIdGen.NextId()
	newAccept := rnfa.stateIdGen.NextId()

	rnfa.AddTransition(newStart, rnfa.Start, Epsilon)
	rnfa.AddTransition(rnfa.Accept, rnfa.Start, Epsilon)
	rnfa.AddTransition(newStart, newAccept, Epsilon)
	rnfa.AddTransition(rnfa.Accept, newAccept, Epsilon)
	rnfa.Start = newStart
	rnfa.Accept = newAccept
}

func (rnfa *RegexNfa) Optional() {
	newStart := rnfa.stateIdGen.NextId()
	newAccept := rnfa.stateIdGen.NextId()

	rnfa.AddTransition(newStart, newAccept, Epsilon)
	rnfa.AddTransition(newStart, rnfa.Start, Epsilon)
	rnfa.AddTransition(rnfa.Accept, newAccept, Epsilon)
	rnfa.Start = newStart
	rnfa.Accept = newAccept
}

func (rnfa *RegexNfa) Plus() {
	newStart := rnfa.stateIdGen.NextId()
	newAccept := rnfa.stateIdGen.NextId()

	rnfa.AddTransition(newStart, rnfa.Start, Epsilon)
	rnfa.AddTransition(rnfa.Accept, rnfa.Start, Epsilon)
	rnfa.AddTransition(rnfa.Accept, newAccept, Epsilon)
	rnfa.Start = newStart
	rnfa.Accept = newAccept
}

func buildNfaFromRegexp(re *syntax.Regexp, stateIdGen *StateIdGenerator) *RegexNfa {
	switch re.Op {
	case syntax.OpEmptyMatch:
		nfa := NewBaseNfa(stateIdGen)
		start := stateIdGen.NextId()
		accept := stateIdGen.NextId()
		nfa.AddTransition(start, accept, Epsilon)
		return &RegexNfa{
			BaseNfa: nfa,
			Start:   start,
			Accept:  accept,
		}

	case syntax.OpLiteral:
		nfa := NewBaseNfa(stateIdGen)
		start := stateIdGen.NextId()
		current := start
		for _, r := range re.Rune {
			next := stateIdGen.NextId()
			nfa.AddTransition(current, next, r)
			current = next
		}
		return &RegexNfa{
			BaseNfa: nfa,
			Start:   start,
			Accept:  current,
		}

	case syntax.OpCharClass:
		nfa := NewBaseNfa(stateIdGen)
		start := stateIdGen.NextId()
		accept := stateIdGen.NextId()
		for _, r := range getCharClassRunes(re) {
			nfa.AddTransition(start, accept, r)
		}
		return &RegexNfa{
			BaseNfa: nfa,
			Start:   start,
			Accept:  accept,
		}

	case syntax.OpConcat:
		nfa := buildNfaFromRegexp(re.Sub[0], stateIdGen)
		for _, sub := range re.Sub[1:] {
			nextNfa := buildNfaFromRegexp(sub, stateIdGen)
			nfa.Concatenate(nextNfa)
		}
		return nfa

	case syntax.OpAlternate:
		nfa1 := buildNfaFromRegexp(re.Sub[0], stateIdGen)
		nfa2 := buildNfaFromRegexp(re.Sub[1], stateIdGen)
		nfa1.Alterate(nfa2)
		return nfa1

	case syntax.OpStar:
		nfa := buildNfaFromRegexp(re.Sub[0], stateIdGen)
		nfa.Closure()
		return nfa

	case syntax.OpPlus:
		nfa := buildNfaFromRegexp(re.Sub[0], stateIdGen)
		nfa.Plus()
		return nfa

	case syntax.OpQuest:
		nfa := buildNfaFromRegexp(re.Sub[0], stateIdGen)
		nfa.Optional()
		return nfa

	default:
		panic(fmt.Sprintf("Unsupported regex operation: %v", re.Op))
	}
}

func getCharClassRunes(re *syntax.Regexp) []rune {
	var runes []rune
	for i := 0; i < len(re.Rune); i += 2 {
		start := re.Rune[i]
		end := re.Rune[i+1]
		for r := start; r <= end; r++ {
			runes = append(runes, r)
		}
	}
	return runes
}

func BuildNFA(stateIdGen *StateIdGenerator, pattern string, tokenType string, tokenPriority int) *RegexNfa {
	re, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		panic(err)
	}
	re = re.Simplify()

	nfa := buildNfaFromRegexp(re, stateIdGen)
	nfa.TokenInfo = dfa.TokenInfo{Name: tokenType, Priority: tokenPriority}
	return nfa
}

type UnifiedNfa struct {
	*BaseNfa
	Start   int
	Accepts map[int]dfa.TokenInfo
}

func NewUnifiedNfa(stateIdGen *StateIdGenerator) *UnifiedNfa {
	base := NewBaseNfa(stateIdGen)
	start := stateIdGen.NextId()
	return &UnifiedNfa{
		BaseNfa: base,
		Start:   start,
		Accepts: make(map[int]dfa.TokenInfo),
	}
}

func (unfa *UnifiedNfa) AddRegex(rnfa *RegexNfa) {
	unfa.AddTransition(unfa.Start, rnfa.Start, Epsilon)
	unfa.MergeTransitions(rnfa.BaseNfa)
	unfa.Accepts[rnfa.Accept] = rnfa.TokenInfo
}

func (unfa *UnifiedNfa) GetAlphabet() []rune {
	symbolSet := make(map[rune]struct{})
	alphabet := make([]rune, 0)

	for _, possibleTransitions := range unfa.transitions {
		for key := range possibleTransitions {
			if _, ok := symbolSet[key]; !ok && key != Epsilon {
				alphabet = append(alphabet, key)
				symbolSet[key] = struct{}{}
			}
		}
	}

	return alphabet
}

func (unfa *UnifiedNfa) BuildDFA() *dfa.DFA {
	curId := 0
	alphabet := unfa.GetAlphabet()
	dfaTransitions := make(map[int]map[rune]int)
	dfaAccepts := make(map[int]dfa.TokenInfo)

	q0 := unfa.computeEpsilonClosure([]int{unfa.Start})
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

		tokenMap := make(map[dfa.TokenInfo]struct{})
		for _, state := range q {
			if tokenInfo, exists := unfa.Accepts[state]; exists {
				tokenMap[tokenInfo] = struct{}{}
			}
		}
		acceptInfos := make([]dfa.TokenInfo, 0, len(tokenMap))
		for token := range tokenMap {
			acceptInfos = append(acceptInfos, token)
		}
		sort.Slice(acceptInfos, func(i, j int) bool {
			return acceptInfos[i].Priority < acceptInfos[j].Priority
		})

		if len(acceptInfos) > 0 {
			dfaAccepts[dfaState] = acceptInfos[0]
		}

		for _, symbol := range alphabet {
			reachableStates := unfa.computeEpsilonClosure(unfa.Delta(q, symbol))

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

func (unfa *UnifiedNfa) computeEpsilonClosure(states []int) []int {
	visited := make(map[int]struct{})
	queue := queue.NewQueue[int]()
	result := make([]int, 0)

	for _, state := range states {
		if _, exists := visited[state]; !exists {
			visited[state] = struct{}{}
			queue.Push(state)
			result = append(result, state)
		}
	}

	for queue.Length() > 0 {
		state := queue.Pop()

		if transitions, ok := unfa.transitions[state]; ok {
			if nextStates, ok := transitions[Epsilon]; ok {
				for _, nextState := range nextStates {
					if _, exists := visited[nextState]; !exists {
						visited[nextState] = struct{}{}
						queue.Push(nextState)
						result = append(result, nextState)
					}
				}
			}
		}
	}

	sort.Ints(result)
	return result
}

func (unfa *UnifiedNfa) Delta(states []int, symbol rune) []int {
	var result []int

	for _, state := range states {
		if possibleTransitions, ok := unfa.transitions[state]; ok {
			if reachableStates, ok := possibleTransitions[symbol]; ok {
				result = append(result, reachableStates...)
			}
		}
	}

	return result
}

func (unfa *UnifiedNfa) PrettyPrint() string {
	visited := make(map[int]bool)
	var builder strings.Builder

	builder.WriteString("OneAcceptNfa:\n")
	builder.WriteString(fmt.Sprintf("Start State: %d\n", unfa.Start))
	builder.WriteString(fmt.Sprintf("Accept State: %v\n", unfa.Accepts))
	builder.WriteString("Transitions:\n")

	unfa.prettyPrintState(&builder, unfa.Start, visited)
	return builder.String()
}

func (unfa *UnifiedNfa) prettyPrintState(builder *strings.Builder, state int, visited map[int]bool) {
	if visited[state] {
		return
	}
	visited[state] = true

	for symbol, nextStates := range unfa.transitions[state] {
		displaySymbol := symbol

		for _, nextState := range nextStates {
			builder.WriteString(fmt.Sprintf("  %d --[%c]--> %d\n", state, displaySymbol, nextState))
			unfa.prettyPrintState(builder, nextState, visited)
		}
	}
}
