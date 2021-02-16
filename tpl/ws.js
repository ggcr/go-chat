let socket = new.WebSocket("ws://localhost:8080/ws")
socket.onopen = () => {
    socked.send("Hi from the client")
}

socked.send("Lol")