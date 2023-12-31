package main

import (
	"fmt"
	"os"
	"strconv"
)

// FullAdderCNF generates CNF clauses for a full adder.
func FullAdderCNF(a, b, cin, s, cout,t,maxBits int) [][]int {
	t2 := t + maxBits
	t3 := t + maxBits*2
	t4 := t + maxBits*3
	clauses := [][]int{
		// s = a XOR b XOR cin // where s = true
		// {-a, -b, s}, {a, b, s}, {-a, -cin, s}, {a, cin, s}, {-b, -cin, s}, {b, cin, s},
		// s = a XOR b XOR cin // where s = false
		// {-s, a, b}, {-s, a, cin}, {-s, b, cin},

		//for s
		//a XOR b = t
		//t XOR cin = s
		{-a, -b, t}, {a, b, t}, {a, -b, -t}, {-a, b, -t},
		{-t, -cin, s}, {t, cin, s}, {t, -cin, -s}, {-t, cin, -s},


		// for c_out
		// t2 = a AND b
		// (-a OR t2) AND (-b OR t2) AND (a OR b OR -t2)
		{-a,t2},{-b,t2},{a,b,-t2},
		// t3 = a AND c_in
		// (-a OR t3) AND (-c OR t3) AND (a OR c OR -t3)
		{-a,t3},{-cin,t3},{a,cin,-t3},
		// t4 = b AND c_in
		//(-b OR t4) AND (-c OR t4) AND (b OR c OR -t4)
		{-b,t4},{-cin,t4},{b,cin,-t4},
		//c_out = t2 OR t3 OR t4
		// (-t2 OR c_out) AND (-t3 OR c_out) AND (-t4 OR c_out) AND (t2 OR t3 OR t4 OR -c_out)
		{-t2,cout},{-t3,cout},{-t4,cout},{t2,t3,t4,-cout},

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
	carryVars := make([]int, maxBits+1)

	//ex:3bit 3bit 1~4 for s, 5~9 for c 
	for i := 0; i < maxBits; i++ {
		sumVars[i] = varIndex
		varIndex++
	}
	for i := 0; i < len(carryVars); i++ {
		carryVars[i] = varIndex
		varIndex++
	}

	// Generate clauses for each full adder
	cnt := 0
	for i := 0; i < maxBits; i++ {
		//ex:3bit 3bit 9~13 for a, 13~17 for b
		// a := getVarIndex(i, n, varIndex)
		a := varIndex+cnt
		cnt++
		b := varIndex+cnt
		cnt++
		cin := 0
		// if i > 0 {
			cin = carryVars[i]
		// }
		s := sumVars[i]
		cout := carryVars[i+1]

		//ex:3bit 3bit 18~
		t:=maxBits*4+2+i
		fmt.Println("s: ", s)
		fmt.Println("cout: ", cout)
		fmt.Println("a: ", a)
		fmt.Println("b: ", b)
		fmt.Println("t: ", t)
		fmt.Println("")
		clauses = append(clauses, FullAdderCNF(a, b, cin, s, cout,t,maxBits)...)
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
	clauses := GenerateAdderCNF(3, 3)

	// Output to a file
	file, err := os.Create("adder_main.txt")
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

