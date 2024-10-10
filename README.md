# TCP Chat Application

A simple TCP chat server implemented in Go that allows multiple clients to connect, send messages, and receive chat history.

## How to Connect

### Using Telnet

1. **Start the Server**:
   - Open a terminal and navigate to the directory where the application is located.
   - Run the server:
     ```go
     go run .
     ```
   - The server will listen on port **8989** by default.

2. **Open a Telnet Client**:
   - In another terminal window, type the following command to connect to the server:
     ```bash
     telnet localhost 8989
     ```
   - If you specified a different port when starting the server, replace `8989` with that port number.

3. **Enter Your Name**:
   - When prompted, type your name and press `Enter`.

4. **Start Chatting**:
   - You can now send messages to other connected clients.

### Using Netcat (nc)

If you prefer using Netcat instead of Telnet, follow these steps:

1. **Start the Server**:
   - Ensure the server is running as described above.

2. **Open a Netcat Client**:
   - In a new terminal window, connect to the server using:
     ```bash
     nc localhost 8989
     ```
   - Again, replace `8989` with your specified port if different.

3. **Enter Your Name**:
   - Type your name when prompted and press `Enter`.

4. **Start Chatting**:
   - You can now communicate with other clients.

## How to Quit

1. **Exit Command**:
   - Type `exit` at the prompt and press `Enter` (if using Telnet).

2. **Keyboard Shortcut**:
   - For Telnet:
     - Press `Ctrl + ]` to access the Telnet prompt.
     - Type `quit` and press `Enter`.

3. **Force Close**:
   - You can close the terminal window or use `Ctrl + C` to terminate the session in Netcat.

## License

This project is licensed under the MIT License.
