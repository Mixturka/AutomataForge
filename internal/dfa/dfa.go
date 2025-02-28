package dfa

import (
	"fmt"
)

const Epsilon = 'ε'

// stateSet represents a set of states.
type stateSet map[int]struct{}

// DFA represents a deterministic finite automaton.
type DFA struct {
	Transitions map[int]map[rune]int
	Start       int
	Accepts     []int
}

// Minimize applies Hopcroft's algorithm to minimize the DFA.
func (dfa *DFA) Minimize() {
	alphabet := dfa.GetAlphabet()

	// Build the full set of states Q.
	Q := make(stateSet)
	for s := range dfa.Transitions {
		Q[s] = struct{}{}
	}

	// Build the set of accepting states.
	acceptSet := make(stateSet)
	for _, a := range dfa.Accepts {
		acceptSet[a] = struct{}{}
	}

	// Partition Q into accepting and non-accepting states.
	acceptStates := make(stateSet)
	nonAcceptStates := make(stateSet)
	for s := range Q {
		if _, ok := acceptSet[s]; ok {
			acceptStates[s] = struct{}{}
		} else {
			nonAcceptStates[s] = struct{}{}
		}
	}

	var P []stateSet
	if len(acceptStates) > 0 {
		P = append(P, acceptStates)
	}
	if len(nonAcceptStates) > 0 {
		P = append(P, nonAcceptStates)
	}

	W := make([]stateSet, len(P))
	copy(W, P)

	for len(W) > 0 {
		A := W[len(W)-1]
		W = W[:len(W)-1]

		for _, c := range alphabet {
			X := make(stateSet)
			for s := range Q {
				if next, ok := dfa.Transitions[s][c]; ok {
					if _, inA := A[next]; inA {
						X[s] = struct{}{}
					}
				}
			}

			var newP []stateSet
			for _, Y := range P {
				intersectYX := intersect(Y, X)
				if isEmpty(intersectYX) || len(intersectYX) == len(Y) {
					newP = append(newP, Y)
				} else {
					diffYX := difference(Y, X)
					newP = append(newP, intersectYX, diffYX)

					found := false
					for i, wSet := range W {
						if sameSet(wSet, Y) {
							// Replace Y with both new sets.
							W[i] = intersectYX
							W = append(W, diffYX)
							found = true
							break
						}
					}
					if !found {
						if len(intersectYX) <= len(diffYX) {
							W = append(W, intersectYX)
						} else {
							W = append(W, diffYX)
						}
					}
				}
			}
			P = newP
		}
	}

	stateToPart := make(map[int]int)
	for i, part := range P {
		for s := range part {
			stateToPart[s] = i
		}
	}

	newTransitions := make(map[int]map[rune]int)
	for i, part := range P {
		var rep int
		for s := range part {
			rep = s
			break
		}
		newTransitions[i] = make(map[rune]int)
		for _, c := range alphabet {
			if next, ok := dfa.Transitions[rep][c]; ok {
				newTransitions[i][c] = stateToPart[next]
			}
		}
	}

	newStart := stateToPart[dfa.Start]
	var newAccepts []int
	for i, part := range P {
		for a := range acceptSet {
			if _, ok := part[a]; ok {
				newAccepts = append(newAccepts, i)
				break
			}
		}
	}

	dfa.Transitions = newTransitions
	dfa.Start = newStart
	dfa.Accepts = unique(newAccepts)
}

func (dfa *DFA) GetAlphabet() []rune {
	symbolSet := make(map[rune]struct{})
	for _, trans := range dfa.Transitions {
		for sym := range trans {
			if sym != Epsilon {
				symbolSet[sym] = struct{}{}
			}
		}
	}
	alphabet := make([]rune, 0, len(symbolSet))
	for sym := range symbolSet {
		alphabet = append(alphabet, sym)
	}
	return alphabet
}

func intersect(s1, s2 stateSet) stateSet {
	res := make(stateSet)
	for x := range s1 {
		if _, ok := s2[x]; ok {
			res[x] = struct{}{}
		}
	}
	return res
}

func difference(s1, s2 stateSet) stateSet {
	res := make(stateSet)
	for x := range s1 {
		if _, ok := s2[x]; !ok {
			res[x] = struct{}{}
		}
	}
	return res
}

func isEmpty(s stateSet) bool {
	return len(s) == 0
}

func sameSet(a, b stateSet) bool {
	if len(a) != len(b) {
		return false
	}
	for x := range a {
		if _, ok := b[x]; !ok {
			return false
		}
	}
	return true
}

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	var result []int
	for _, entry := range intSlice {
		if !keys[entry] {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}

func (dfa *DFA) PrettyPrint() {
	fmt.Println("DFA Start State:", dfa.Start)
	fmt.Println("DFA Accepting States:", dfa.Accepts)
	fmt.Println("DFA Transitions:")
	for state, trans := range dfa.Transitions {
		fmt.Printf("State %d:\n", state)
		for sym, next := range trans {
			fmt.Printf("  On symbol '%c' -> State %d\n", sym, next)
		}
	}
}
