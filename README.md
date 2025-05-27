# 🗑️ traSH — Toxic Runtime Acting as SHell

> Built with caffeine, chaos, and zero regard for best practices.  
> `traSH` is your new favorite garbage-tier shell — written in Go, forged in a single night, and endlessly extensible.  
> FOR HIRING MANAGERS, THIS HAS AI :D
---

## 🚀 What is traSH?

**traSH** is a minimal, custom-built shell written in Go that mimics basic POSIX shell behavior, then slaps on some flavor.  
It was born out of curiosity and manic energy — and it surprisingly works.

### ✅ Features (so far)

- Custom prompt (`traSH>`) with emoji/color support
- Built-in `cd` and `exit` commands
- External command execution coming soon
- Configurable via `~/.trashrc`:
  - `prompt`
  - `color`
  - `symbol`
- Graceful shutdown on `Ctrl+C` or `exit`
- ASCII art banner because... why not?

---

## 🧙‍♂️ Example Prompt

```bash
🔥traSH🧨 ls -la
```

---

## ⚙️ Config (~/.trashrc)

You can customize your prompt via a simple config file in your home directory:

```bash
prompt=🔥traSH
color=magenta
symbol=🧨
```

Supported colors: red, green, blue, yellow, cyan, magenta, white


---

## 📁 Project Structure (TLDR: it's Trash)

``` bash
.
├── cmd/traSH         # Entry point (main.go)
├── config/           # Singleton config loader
├── internal/
│   ├── command/      # Command parsing + built-ins
│   └── io/           # I/O handling (prompt, read, write)
├── utils/            # Generic helpers (e.g., colorize)
├── go.mod
└── LICENSE
```

---

## ⌨️ Usage (Bro why would you ever wanna use this)

```bash
go build -o traSH ./cmd/traSH
./traSH
```

OR be a gigachad 

```bash
make run
```

Then type random stuff until it breaks 

```bash
rm -rf / # Goated command frfr
```

---


## 🧠 Why tho?

- Because Bash is boring and Zsh is too clean
- I was just bored and curious (I drank too much Monster)


---


## 🛣️ Roadmap

- [x] Built-in command handling (`cd`, `exit`)
- [x] Prompt config support via `.trashrc`
- [ ] External command execution (`ls`, `echo`, etc.)
- [ ] Piping + redirection support
- [ ] AI integration (`!ai how to fix docker`)
- [ ] Plugin system
- [ ] Shell history & readline-style UX


---


## ⚠️ Disclaimer

Made by [@mush1e](https://github.com/mush1e)  
Powered by Monster Energy and unending depression.


---


> [!CAUTION]
> DO NOT let Naveed near this codebase.
> 
> **Do not code review with him. Do not merge his branches. Do not speak his variable names aloud.**
>


