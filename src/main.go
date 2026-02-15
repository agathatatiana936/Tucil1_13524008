package main

import "fmt"

func main() {
	// parse
	b, err := parse_board_from_file("input.txt")
	if err != nil {
		fmt.Println("ERR:", err)
		return
	}

	// validate
	if err := validate_board(b); err != nil {
		fmt.Println("INVALID BOARD:", err)
		return
	}

	fmt.Println("N:", b.n)
	fmt.Println("ID Unique:", b.idUnique)
	for i := 0; i < b.n; i++ {
		fmt.Println(string(b.raw[i]), b.id[i])
	}
}