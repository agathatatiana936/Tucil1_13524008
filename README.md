# Queens Puzzle Solver — Brute Force Approach

> **IF2211 – Algorithm Strategies (Tucil 1)**  
> Institut Teknologi Bandung

## Program Overview

This program implements a **pure brute force (exhaustive search) algorithm** to solve the **Queens Puzzle**, inspired by the Queens game available on LinkedIn.

The goal of the puzzle is to place **exactly one queen** on each:

- Row  
- Column  
- Colored region  

with the additional constraint that **no two queens may be adjacent**, including diagonally.

The program reads a board configuration from a `.txt` file and systematically enumerates **all possible queen placements** without using:

- heuristics  
- pruning  
- backtracking optimizations  

If a valid solution is found, the program displays the solution, execution time, and number of configurations examined.  
Additional features include **live brute-force tracing**, **text output**, and **PNG image export** (bonus).

---

## Requirements & Installation

### System Requirements
- **Programming Language**: Go (Golang)
- **Go Version**: Go 1.20 or newer
- **Operating System**:
  - Windows (primary)
  - Linux (supported)

### Dependencies
- Uses **only Go standard libraries**
- No external libraries required

### Go Installation
Download and install Go from:
https://go.dev/dl/

Verify installation:
```bash
go version
``` 

### Repository Structure
Tucil1_NIM/
│
├── bin/
│   └── queen.exe              # Compiled executable
│
├── src/
│   ├── main.go                # Program entry point (CLI)
│   ├── algo.go                # Solver
│   ├── parser.go              # Input parsing & board validation
│   ├── render.go              # Text output & live report (ASCII)
│   ├── render_image.go        # PNG image rendering 
│   ├── cli_ui.go              # Colored ASCII UI utilities
│   ├── go.mod
│   └── go.sum
│
├── test/
│   ├── input.txt              # Sample input board
│   └── input2.txt
│
├── doc/
│   └── Tucil1_IF2211_Report.pdf
│
└── README.md

### Compilation Guide
From src directory, run : 
```bash
cd src
go build -o ../bin/queen.exe .
```

After successful compilation, the executable file will be created at:
```bash
bin/queen.exe
```


### How to Run the Program
Run the program from the root directory of the repository.
```bash
bin\queen.exe
```

### Program Flow
1) The program prompts the user to input the board file path
2) The user chooses whether to:
    a. Save live brute-force tracing (trace.txt)
    b. Save the solution as a .txt file
    c. Save the solution as a .png image
3) The program find the solution with brute force
4) If a solution is found: The solution board is displayed (# represents a queen)
    a. Execution time (milliseconds) is printed
    b. Number of examined configurations is reported
5) The user is asked whether they want to solve another Queens puzzle or exit

### Author
Name : Agatha Tatianingseto
Student ID (NIM) : 13524008
Course : IF2211 – Algorithm Strategies
Institution : Institut Teknologi Bandung