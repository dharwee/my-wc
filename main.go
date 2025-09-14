package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

//count reads from r and returns (lines,words,bytes,error).
//-lines : number of '\n' bytes encountered
//-words: seq of non-whitespaces (unicode.IsSpace) characters
//-bytes: total bytes read

//Implementation uses bufio.Reader.ReadRune so we correctly handle UTF-8
//decoding while still counting raw bytes (size returned by ReadRune)

func count(r io.Reader) (lines, words, bytesCount int, err error) {
	br := bufio.NewReader(r)

	//Am I currently inside a word? no! then set inWord=false
	inWord := false

	for {
		ch, size, e := br.ReadRune()
		if e != nil {
			if e == io.EOF {
				break
			}
			return lines, words, bytesCount, e
		}

		//add number of bytes for this rune
		bytesCount += size

		//lines:count newlines characters '\n'
		if ch == '\n' {
			lines++
		}

		//word detection: enter a word when we see a non-space and were not in a word
		if unicode.IsSpace(ch) {
			//If yes, it means we are definitely not in a word. So, the code sets inWord = false
			inWord = false
			// This 'else' block runs only when we find a non-space character.
		} else {
			//inWord = true means "We are inside a word."
			if !inWord {
				words++
				inWord = true
			}
		}
	}
	return lines, words, bytesCount, nil
}

// printCounts prints selected counts (inorder lines, words, bytes) formatted in columns
// similar to Unix wc. If name is empty it's omitted (useful for stdin)
func printCounts(lines, words, bytesCount int, name string, showL, showW, showC bool) {
	if showL {
		fmt.Printf("%8d", lines)
	}
	if showW {
		fmt.Printf("%8d", words)
	}
	if showC {
		fmt.Printf("%8d", bytesCount)
	}
	if name != "" {
		fmt.Printf(" %s", name)
	}
	fmt.Println()
}

func main() {
	//flags
	lFlag := flag.Bool("l", false, "print the newline counts")
	wFlag := flag.Bool("w", false, "print the word counts")
	cFlag := flag.Bool("c", false, "print the byte counts")
	//flag.Parse(): This is the crucial command that
	// actually reads the command line (like my-wc -l -w file.txt)
	// and sets the values of the flag variables (lFlag, wFlag, cFlag) accordingly.
	flag.Parse()

	//default: if the user didn't specify any of -l -w -c, show all three
	if !*lFlag && !*wFlag && !*cFlag {
		//	The * is used to get the boolean value from the pointers
		*lFlag, *wFlag, *cFlag = true, true, true
	}

	//args:=flag.Args(): This gets a list of
	// all command-line arguments that were not flags.
	// These are typically filenames
	args := flag.Args()

	//if len(args)==0: This checks if the list
	// of filenames is empty. An empty list means the user didn't
	// provide a file, so the program should read from standard input.
	// This is what allows you to "pipe" data into it (e.g., cat somefile.txt | ./my-wc).
	if len(args) == 0 {
		//read from stdin
		lines, words, bytesCount, err := count(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "my-wc", err)
			os.Exit(1)
		}
		printCounts(lines, words, bytesCount, "", *lFlag, *wFlag, *cFlag)
		return
	}

	//handle one or more files; show totals when multiple
	//files processed successfully
	var totalLines, totalWords, totalBytes int
	var filesProcessed int

	for _, name := range args {
		var f *os.File
		var err error

		if name == "-" {
			//convention: "-" means read from stdin
			f = os.Stdin
		} else {
			f, err = os.Open(name) //eg. file not found error
			if err != nil {
				// report error but continue with next file (like many Unix tools)
				fmt.Fprintf(os.Stderr, "my-wc: cannot opem %s:%v\n", name, err)
				continue
			}
		}

		lines, words, bytesCount, err := count(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "my-wc:error reading %s:%v\n", name, err)
			if f != os.Stdin {
				f.Close()
			}
			continue
		}

		printCounts(lines, words, bytesCount, name, *lFlag, *wFlag, *cFlag)
		totalLines += lines
		totalWords += words
		totalBytes += bytesCount
		filesProcessed++

		//Is this file something I opened myself,
		// or is it the shared standard input?
		if f != os.Stdin {
			//If f is a regular file you opened, the condition is true, and 
			// f.Close() is called to clean up properly.
			//If f is os.Stdin, the condition is false, 
			//and the f.Close() call is skipped, safely leaving the main input channel alone.

			//Rule is simple: Don't close what you didn't open
			f.Close()
		}

	}
	if filesProcessed > 1 {
		printCounts(totalLines, totalWords, totalBytes, "total", *lFlag, *wFlag, *cFlag)
	}
}
