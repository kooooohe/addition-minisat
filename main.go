package main

import (
	"fmt"
	"os"
	"strconv"
)

// FullAdderCNF generates CNF clauses for a full adder.
func FullAdderCNF(a, b, cin, s, cout int) [][]int {
	clauses := [][]int{
		{-a, -b, s}, {a, b, s}, {-a, -cin, s}, {a, cin, s}, {-b, -cin, s}, {b, cin, s},
		{-s, a, b}, {-s, a, cin}, {-s, b, cin},
		{-a, -b, -cin, cout}, {a, b, cout}, {a, cin, cout}, {b, cin, cout},
	}
	return clauses
}

// GenerateAdderCNF generates CNF for adding two binary numbers.
func GenerateAdderCNF(n, m int) [][]int {
	var clauses [][]int
	maxBits := max(n, m) + 1

	// Variables for the adder
	varIndex := 1
	sumVars := make([]int, maxBits)
	carryVars := make([]int, maxBits)
	for i := 0; i < maxBits; i++ {
		sumVars[i] = varIndex
		varIndex++
		carryVars[i] = varIndex
		varIndex++
	}

	// Generate clauses for each full adder
	for i := 0; i < maxBits; i++ {
		a := getVarIndex(i, n, varIndex)
		b := getVarIndex(i, m, varIndex)
		cin := 0
		if i > 0 {
			cin = carryVars[i-1]
		}
		s := sumVars[i]
		cout := carryVars[i]

		clauses = append(clauses, FullAdderCNF(a, b, cin, s, cout)...)
	}

	return clauses
}

// Utility functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getVarIndex(bit, size, varIndex int) int {
	if bit < size {
		return varIndex + bit*2
	}
	return 0 // Zero for bits beyond the size
}

func main() {
	// Example usage: Generate CNF for 3-bit and 4-bit addition
	clauses := GenerateAdderCNF(3, 4)

	// Output to a file
	file, err := os.Create("adder_3bit_4bit.cnf")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// DIMACS format header
	_, err = file.WriteString(fmt.Sprintf("p cnf %d %d\n", len(clauses)*4, len(clauses)))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// Write clauses
	for _, clause := range clauses {
		for _, literal := range clause {
			_, err = file.WriteString(strconv.Itoa(literal) + " ")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
		_, err = file.WriteString("0\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	fmt.Println("CNF file created successfully")
}

