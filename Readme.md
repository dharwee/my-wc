# my-wc 📝

A simplified version of the Unix `wc` (word count) command written in Go.  
It counts **lines**, **words**, and **bytes** from files or standard input.

---

## 🚀 Features
- Count **lines** (`-l`), **words** (`-w`), and **bytes** (`-c`)
- Works with one or more files
- Supports reading from **stdin** (pipe input or `-`)
- Prints totals when multiple files are given
- Behaves like the standard Unix `wc` command

---

## 📦 Installation

1. Clone the repo:
   ```bash
   git clone https://github.com/your-username/my-wc.git
   cd my-wc
   ```
2. Build the binary:
    ```bash
    go build -o my-wc
    ```
3. Run it:
```bash
./my-wc -l -w -c file.txt
```

## Flags

-l : print the newline counts

-w : print the word counts

-c : print the byte counts

## 📚 Concepts Learned

Command-line argument parsing with flag

File I/O with os.Open

Streaming input with io.Reader

Counting words using unicode.IsSpace

Handling stdin and multiple files

