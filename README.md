# 🔢 Big Number Calculator — High-Precision Calculator in Go

A simple yet powerful calculator built in Go that supports **arbitrarily large numbers** using `math/big`. Perform calculations like `100!`, `2^100`, `√2`, and high-precision division with ease. Comes with a clean web interface and a REST API.

---

## 🛠️ Features

- ✅ Addition, subtraction, multiplication, division, power, factorial, square root
- ✅ Supports **very large numbers** via `math/big`
- ✅ High precision arithmetic (512-bit precision)
- ✅ Modern, responsive UI with light/dark mode toggle
- ✅ REST API + HTML frontend
- ✅ Unit and integration tests
- ✅ CI with GitHub Actions
- ✅ Makefile for easy automation

---

## 🚀 How to Run

### 1. Clone the repository

```bash
git clone https://github.com/DanRulev/calculator.git
cd calculator
```

### 2. Download dependencies
```bash
go mod download
```

### 3. Start the server
```bash
go run cmd/main.go
```

Open in your browser: 🔗 http://localhost:8080

---

## 📦 Tech Stack
- Go – Backend logic
- Gin – Web framework
- math/big – Arbitrary-precision arithmetic
- HTML/CSS/JavaScript – Frontend UI
- GitHub Actions – CI
- Makefile – Build automation

## 📄 License
This project is licensed under the MIT License – see the LICENSE file for details.