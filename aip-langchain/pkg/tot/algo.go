package tot

import (
	"fmt"
)

// Thought represents a single thought step in the tree.
type Thought struct {
	Text string
}

// State represents a node in the tree, representing a partial solution.
type State struct {
	Input    string
	Thoughts []Thought
}

// ThoughtGenerator generates potential thoughts from a given state.
func ThoughtGenerator(pTheta func(State) []Thought, s State, k int) []Thought {
	return pTheta(s)
}

// StateEvaluator evaluates the progress made by a state towards solving the problem.
func StateEvaluator(pValue func(State) float64, s State) float64 {
	return pValue(s)
}

// ToT-BFS implements the Tree of Thoughts breadth-first search algorithm.
func ToTBFS(x string, pTheta func(State) []Thought, k int, pValue func(State) float64, T int, b int) Thought {
	S0 := State{
		Input:    x,
		Thoughts: nil,
	}

	St := []State{S0}

	for t := 1; t <= T; t++ {
		StPrime := make([]State, 0)

		for _, s := range St {
			zCandidates := ThoughtGenerator(pTheta, s, k)

			for _, zt := range zCandidates {
				sNew := State{
					Input:    s.Input,
					Thoughts: append(s.Thoughts, zt),
				}
				StPrime = append(StPrime, sNew)
			}
		}

		Vt := make([]float64, len(StPrime))
		for i, s := range StPrime {
			Vt[i] = StateEvaluator(pValue, s)
		}

		bestStates := make([]State, 0, b)
		bestValues := make([]float64, 0, b)

		for _, s := range StPrime {
			if len(bestStates) < b {
				bestStates = append(bestStates, s)
				bestValues = append(bestValues, StateEvaluator(pValue, s))
			} else {
				minValue := bestValues[0]
				minIndex := 0

				for j := 1; j < b; j++ {
					if bestValues[j] < minValue {
						minValue = bestValues[j]
						minIndex = j
					}
				}

				if StateEvaluator(pValue, s) > minValue {
					bestStates[minIndex] = s
					bestValues[minIndex] = StateEvaluator(pValue, s)
				}
			}
		}

		St = bestStates
	}

	return St[0].Thoughts[0]
}

// ToT-DFS implements the Tree of Thoughts depth-first search algorithm.
func ToTDFS(s State, t int, pTheta func(State) []Thought, k int, pValue func(State) float64, T int, vth float64) {
	if t > T {
		fmt.Println(s.Thoughts[0].Text)
		return
	}

	zCandidates := ThoughtGenerator(pTheta, s, k)

	for _, sPrime := range zCandidates {
		if StateEvaluator(pValue, sPrime) > vth {
			ToTDFS(sPrime, t+1, pTheta, k, pValue, T, vth)
		}
	}
}

func main() {
	// Define the functions pTheta and pValue according to your specific problem.

	// Example usage of ToT-BFS
	x := "input"
	k := 5
	T := 10
	b := 3

	pTheta := func(s State) []Thought {
		// Implement your thought generation logic based on the state.
		return []Thought{}
	}

	pValue := func(s State) float64 {
		// Implement your state evaluation logic.
		return 0.0
	}

	result := ToTBFS(x, pTheta, k, pValue, T, b)
	fmt.Println(result.Text)

	// Example usage of ToT-DFS
	s := State{
		Input:    "input",
		Thoughts: nil,
	}
	t := 0
	vth := 0.5

	ToTDFS(s, t, pTheta, k, pValue, T, vth)
}
