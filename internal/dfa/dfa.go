package dfa

import (
	"fmt"
	"maps"
	"slices"
	"sort"
)

const Epsilon = 'ε'

type stateSet map[int]struct{}

type DFA struct {
	Transitions map[int]map[rune]int
	Start       int
	Accepts     map[int]string
}

func (dfa *DFA) Minimize() {
	alphabet := dfa.GetAlphabet()
	stateToSet := make(map[int]*stateSet)
	T := make(map[*stateSet]struct{})

	partitions := make(map[string]*stateSet)
	for state, token := range dfa.Accepts {
		key := token
		if set, exists := partitions[key]; exists {
			(*set)[state] = struct{}{}
			stateToSet[state] = set
		} else {
			newSet := make(stateSet)
			newSet[state] = struct{}{}
			partitions[key] = &newSet
			stateToSet[state] = &newSet
			T[&newSet] = struct{}{}
		}
	}

	nonAcceptStates := make(stateSet)
	for state := range dfa.Transitions {
		if _, isAccept := dfa.Accepts[state]; !isAccept {
			nonAcceptStates[state] = struct{}{}
		}
	}
	if len(nonAcceptStates) > 0 {
		T[&nonAcceptStates] = struct{}{}
		stateToSetPtr := &nonAcceptStates
		for state := range nonAcceptStates {
			stateToSet[state] = stateToSetPtr
		}
	}

	P := make(map[*stateSet]struct{})
	for !maps.Equal(T, P) {
		P = make(map[*stateSet]struct{})
		maps.Copy(P, T)
		T = make(map[*stateSet]struct{})

		for p := range P {
			split := dfa.Split(p, stateToSet, alphabet)
			if len(split) > 0 {
				for _, newSet := range split {
					for state := range *newSet {
						stateToSet[state] = newSet
					}
					T[newSet] = struct{}{}
				}
			} else {
				T[p] = struct{}{}
			}
		}
	}

	newDfaTransitions := make(map[int]map[rune]int)
	newStateIdCounter := 0
	stateSetToNewStateIds := make(map[*stateSet]int)

	for stateSet := range T {
		stateSetToNewStateIds[stateSet] = newStateIdCounter
		newStateIdCounter++
	}

	for _, newStateId := range stateSetToNewStateIds {
		newDfaTransitions[newStateId] = make(map[rune]int)
	}

	for _, c := range alphabet {
		for stateSet := range T {
			for state := range *stateSet {
				if toState, exists := dfa.Transitions[state][c]; exists {
					newDfaTransitions[stateSetToNewStateIds[stateSet]][c] = stateSetToNewStateIds[stateToSet[toState]]
				}
			}
		}
	}

	dfa.Transitions = newDfaTransitions

	newAccepts := make(map[int]string)
	for stateSet, newStateId := range stateSetToNewStateIds {
		for state := range *stateSet {
			if token, exists := dfa.Accepts[state]; exists {
				newAccepts[newStateId] = token
				break
			}
		}
	}

	dfa.Accepts = newAccepts
	dfa.Start = stateSetToNewStateIds[stateToSet[dfa.Start]]
}
func (dfa *DFA) Split(states *stateSet, stateToSet map[int]*stateSet, alphabet []rune) []*stateSet {
	for _, c := range alphabet {
		newStateSet1 := make(stateSet)
		newStateSet2 := make(stateSet)

		var refPartition *stateSet
		isFirst := true

		for state := range *states {
			var curPartition *stateSet

			if transitions, exists := dfa.Transitions[state]; exists {
				if next, exists := transitions[c]; exists {
					curPartition = stateToSet[next]
				}
			}

			if isFirst {
				isFirst = false
				refPartition = curPartition
			}

			if refPartition == curPartition {
				newStateSet1[state] = struct{}{}
			} else {
				newStateSet2[state] = struct{}{}
			}
		}

		if len(newStateSet1) > 0 && len(newStateSet2) > 0 {
			return []*stateSet{&newStateSet1, &newStateSet2}
		}
	}

	return nil
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

func (dfa *DFA) BuildClassifierTable() map[rune]int {
	states := dfa.getStates()

	alphabet := dfa.GetAlphabet()
	slices.Sort(alphabet)

	transitions := make(map[rune][]int)
	for _, r := range alphabet {
		statesTo := make([]int, len(states))
		for i, state := range states {
			if stateTo, exists := dfa.Transitions[state][r]; exists {
				statesTo[i] = stateTo
			} else {
				statesTo[i] = -1
			}
		}
		transitions[r] = statesTo
	}

	runeGroups := make(map[string][]rune)
	for r, transitionOn := range transitions {
		statesTo := fmt.Sprint(transitionOn)
		runeGroups[statesTo] = append(runeGroups[statesTo], r)
	}

	runeGroupKeys := make([]string, 0, len(runeGroups))
	for key := range runeGroups {
		runeGroupKeys = append(runeGroupKeys, key)
	}

	classifierTable := make(map[rune]int)
	for id, key := range runeGroupKeys {
		for _, r := range runeGroups[key] {
			classifierTable[r] = id
		}
	}

	return classifierTable
}

func (dfa *DFA) BuildTransitionTable(classifierTable map[rune]int) [][]int {
	states := dfa.getStates()

	classIds := make(map[int]struct{})
	for _, id := range classifierTable {
		classIds[id] = struct{}{}
	}
	classes := make([]int, 0, len(classIds))
	for class := range classIds {
		classes = append(classes, class)
	}
	sort.Ints(classes)

	reverseClassifier := make(map[int]rune) // stores random rune from class since they all behave the same
	for r, class := range classifierTable {
		reverseClassifier[class] = r
	}

	transitionTable := make([][]int, len(states))

	for i, state := range states {
		transitionTable[i] = make([]int, len(classes))
		for j, class := range classes {
			transitionTable[i][j] = -1
			if r, exists := reverseClassifier[class]; exists {
				if stateTo, exists := dfa.Transitions[state][r]; exists {
					transitionTable[i][j] = stateTo
				}
			}
		}
	}

	return transitionTable
}

func (dfa *DFA) BuildTypeTable() map[int]string {
	typeTable := make(map[int]string)

	for state, tokenName := range dfa.Accepts {
		typeTable[state] = tokenName
	}

	return typeTable
}

func (dfa *DFA) getStates() []int {
	states := make([]int, 0, len(dfa.Transitions))
	for state := range dfa.Transitions {
		states = append(states, state)
	}
	sort.Ints(states)

	return states
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
