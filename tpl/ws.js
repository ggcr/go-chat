let socket = new WebSocket("ws://localhost:8080/ws")
socket.onopen = function () {
    console.log("SOCKET")
    socket.send("Hi from socket")
}

let chat = document.getElementById("ulLog")
chat.onchange = function () {
    console.log("CHAT CHANGED BRO")
    socket.send("CHAT CHANGED BRO")
}