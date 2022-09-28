var socket = new WebSocket("ws://localhost:9090/echo");

socket.onopen = function () {
    console.log("Status: Connected\n")
};

socket.onmessage = function (e) {
    console.log("Server: " + e.data + "\n")
};

function send(message) {
    socket.send(message);
}