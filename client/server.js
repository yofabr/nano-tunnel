import { randomUUID } from "crypto";
import express from "express";
import http from "http";
import { Server } from "socket.io";

const app = express();
const server = http.createServer(app);
app.use(express.json());

app.get("/api/hello", (req, res) => {
  res.json({ message: "Hello from Express!" });
});

app.get("/send", (req, res) => {
    io.emit("broad", "This is broadcaster!")
    res.json({
        message: "From send"
    })
})

const io = new Server(server, {
  cors: { origin: "*" }
});

let clients = {};

io.path("/socket.io").on("connection", (socket) => {
  console.log("New client connected:", socket.id);

  let clientID = randomUUID();
  clients[socket.id] = clientID;

  socket.emit("welcome", {
    clientID: clientID
  });

  socket.on("disconnect", () => {
    console.log("Client disconnected:", socket.id);
    delete clients[socket.id];
  });
});

const PORT = 8080;
server.listen(PORT, () => {
  console.log(`Server running on port :${PORT}`);
});
