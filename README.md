# asciihex

[![Go Reference](https://pkg.go.dev/badge/github.com/0x5a17ed/asciihex.svg)](https://pkg.go.dev/github.com/0x5a17ed/sqltx)
[![License: 0BSD](https://img.shields.io/badge/License-0BSD-blue.svg)](https://opensource.org/licenses/0BSD)

**asciihex** is a simple and expressive Go library for encoding and decoding strings or byte slices into a *human-friendly*, roundtrip-safe format.

It strikes a careful balance between readability and precision — making it ideal for logs, debugging, and systems introspection.


## ✨ Features

The syntax is designed to be both human-readable and safe to work with:

* Prints printable ASCII characters as-is
* Encodes control characters using caret (`^`) notation (`^@`, `^A`, ..., `^?`)
* Encodes non-ASCII bytes using hex escape (`~8F`)
* Escapes special characters:
    * `~~` → `~`
    * `~^` → `^`
* Fully reversible and roundtrip-safe

The package interface is limited to two functions:
* `Encode(data []byte) string`  
  Encodes a byte slice into a string
* `Decode(data string) ([]byte, error)`   
  Decodes a string back into a byte slice


## 📦 Installation

No third-party dependencies. Drop the package into your project:

```bash
go get github.com/0x5a17ed/asciihex
```


## 🚀 Quick Start

```go
data := []byte{'H', 'e', 'l', 'l', 'o', 0x0A, 0x1E, 0x8F, '~', '^'}
encoded := asciihex.Encode(data)
// encoded == "Hello^J^^~8F~~~^"

decoded, err := asciihex.Decode(encoded)
// decoded == original data
```


## 🌱 Motivation

This package was built with care to provide a log-safe, human-legible format for debugging binary data. Whether you're inspecting packet payloads, serial protocols, or log traces, `asciihex` gives you a clean, reversible format that won’t make your eyes bleed.


## 📜 License

This project is licensed under the 0BSD Licence — see the [LICENCE](LICENSE) file for details.


## 🥇 Acknowledgments

The design and the implementation are roughly based on the idea and syntax of the caret notation binary data representation.

---

<p align="center">Made with ❤️ for data you want to understand</p>
