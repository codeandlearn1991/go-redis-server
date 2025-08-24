# A Redis Server from Scratch in Go 🚀

This project is a personal exploration into the internals of Redis, implemented from the ground up in Go. It's an educational endeavor to understand how a key-value store like Redis functions under the hood, including its data structures, network protocols, and persistence mechanisms.

## 🌟 Features

  * **REPL Protocol Implementation:** Understand and handle the Redis Serialization Protocol (RESP).
  * **In-Memory Key-Value Store:** A basic in-memory data store using Go's built-in `map`.
  * **Concurrency Safe:** The server is built to handle multiple client connections concurrently using goroutines and mutexes.

## 🛠️ How to Run

### Prerequisites

  * Go 1.24 or higher

### Steps

1.  Clone the repository:

    ```bash
    git clone https://github.com/your-username/go-redis-from-scratch.git
    cd go-redis-from-scratch
    ```

2.  Run the server:

    ```bash
    make run
    ```

    The server will start and listen on `localhost:6379`.

3.  Connect to the server using `redis-cli`:

    ```bash
    redis-cli
    ```

    Now you can interact with your server\! Try simple commands like `PING` or `SET mykey "hello"`.

## ✍️ Contributing

This project is primarily a learning tool and isn't intended to be a production-ready Redis replacement. However, if you find a bug or have an idea for an educational feature, feel free to open an issue or submit a pull request.

## 📄 License

This project is licensed under the MIT License - see the `LICENSE` file for details.