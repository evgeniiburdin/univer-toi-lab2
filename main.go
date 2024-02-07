package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Considering this is a program to encode text with Huffman algorithm,
// let's make a custom type representing a node of a binary tree.
// Field names should speak for themselves
type TreeNode struct {
	Symbol      string
	Probability float64
	BinaryValue int
	LeftNode    *TreeNode
	RightNode   *TreeNode
	Parent      *TreeNode
	// This one is used to respresent wether a node is already in tree or not yet
	Processed bool
	IsRoot    bool
	IsEnd     bool
}

func main() {

	// Let's give a program some string to encode
	var strToEncode string
	fmt.Println(`String to encode: `)
	reader := bufio.NewReader(os.Stdin)
	strToEncode, err := reader.ReadString('\n')
	strToEncode = strToEncode[:len(strToEncode)-1]
	if err != nil {
		log.Fatal(err)
	}

	// A map, representing each symbol probability in the given string
	symbolCountMap := make(map[string]int)
	strLen := len(strToEncode)
	for i := 0; i < strLen; i++ {
		symbolCountMap[string(strToEncode[i])]++
	}

	// A map, that will contain all tree nodes
	dynTreeMap := make(map[string]*TreeNode)

	// A map, that will contain codes for each letter in the given string alphabet
	codeMap := make(map[string]string)

	// This loop's purpose is to initialize our two maps with original symbols from the alphabet
	for j := range symbolCountMap {
		tempNode := TreeNode{j, float64(symbolCountMap[j]) / float64(strLen), 1, nil, nil, nil, false, false, true}
		dynTreeMap[j] = &tempNode
		codeMap[j] = ""
	}

	// A boolean, representing the completion of the whole tree building process
	var treeIsDone bool

	// A loop to build a tree
	for {

		// This "if" condition will activate only when the last, root node is built
		// and will start the encoding process.
		if treeIsDone {

			// This variable stores the symbol value of the root node of the tree
			var rootString string = ""
			for comb := range dynTreeMap {
				if len(dynTreeMap[comb].Symbol) > len(rootString) {
					rootString = dynTreeMap[comb].Symbol
				}
			}

			// Changing root node's "IsRoot" field to "true"
			dynTreeMap[rootString].IsRoot = true

			//fmt.Print("Making a key:value map... \n")

			// This loop is used to get codes for all symbols of the initial string alphabet
			for dj := range dynTreeMap {

				// This is how the loop understands wether the node is the end node or not
				if dynTreeMap[dj].IsEnd {
					var code string = ""
					//fmt.Println("\n\nGoing up the... ", dynTreeMap[dj].Symbol)
					code = GoToRoot(dynTreeMap[dj], code)

					// A recursive function returns binary code, which we immediately include right
					// into our map for symbol codes
					codeMap[dj] = code
				}
			}

			// The end of the program. User is given a code map for each symbol
			fmt.Println("\nDone!\n\nCode table: ")
			for e := range codeMap {
				fmt.Printf("%v %v\n", e, codeMap[e])
			}
			break
		}

		// Variables to store two symbols of the alphabet with less probability
		var LessProb0 string
		var LessProb1 string

		// These two loops' purpose is to pick aforementioned symbols from the tree.
		// We will use them to create a new parent node, which parameters will be taken
		// from both child nodes
		var InitProb float64 = 1.0
		for i := range dynTreeMap {
			if dynTreeMap[i].Probability < InitProb && !dynTreeMap[i].Processed {
				InitProb = dynTreeMap[i].Probability
				LessProb0 = i
			}
		}

		dynTreeMap[LessProb0].BinaryValue = 0
		dynTreeMap[LessProb0].Processed = true

		InitProb = 1.0
		for j := range dynTreeMap {
			if dynTreeMap[j].Probability < InitProb && dynTreeMap[j].Symbol != dynTreeMap[LessProb0].Symbol && !dynTreeMap[j].Processed {
				InitProb = dynTreeMap[j].Probability
				LessProb1 = j
			}
		}

		dynTreeMap[LessProb1].BinaryValue = 1
		dynTreeMap[LessProb1].Processed = true

		// Creating a parent node
		parentNode := TreeNode{
			dynTreeMap[LessProb0].Symbol + dynTreeMap[LessProb1].Symbol,
			dynTreeMap[LessProb0].Probability + dynTreeMap[LessProb1].Probability,
			1,
			dynTreeMap[LessProb0],
			dynTreeMap[LessProb1],
			nil,
			false,
			false,
			false,
		}

		// Making connection between out new parent node and its children
		dynTreeMap[LessProb0].Parent = &parentNode
		dynTreeMap[LessProb1].Parent = &parentNode

		// Including our newly-created parent node in our dynTreeMap map, containing all the nodes of the tree
		dynTreeMap[parentNode.Symbol] = &parentNode

		// When parent's node probability is more(float64 type issues) or equals 1.0
		// it means that the parent node's string contains all the alphabet symbols and
		// therefore we should end our loop
		if parentNode.Probability >= 1 {
			treeIsDone = true
		}
	}
}

// A function that simply reverses the given string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// A recursive function that takes an end node of the tree and goes right
// back to the root node, remembering nodes' bin values
// It returns code for requested symbol
func GoToRoot(node *TreeNode, code string) string {
	//fmt.Printf("\n symbol: %v    is root: %v    ", node.Symbol, node.IsRoot)

	if !node.IsRoot {
		//fmt.Printf("binary value: %v", node.BinaryValue)
		code += fmt.Sprint(node.BinaryValue)
		return GoToRoot(node.Parent, code)
	}

	if node.IsRoot {
		code = Reverse(code)
	}

	return code
}
