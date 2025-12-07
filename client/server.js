import express from "express";
import http from "http";
import { WebSocketServer } from "ws";
import { randomUUID } from "crypto";

const app = express();
const server = http.createServer(app);
app.use(express.json());

app.get("/api/hello", (req, res) => {
  res.json({ message: "Hello from Express!" });
});

app.get("/send", (req, res) => {
  // Broadcast message to all WebSocket clients
  wss.clients.forEach((client) => {
    if (client.readyState === 1) { // WebSocket.OPEN
      client.send(JSON.stringify({ event: "forward", message: "This is broadcaster!" }));
    }
  });
  res.json({ message: "From send" });
});

// Create WebSocket server
const wss = new WebSocketServer({ server, path: "/ws" });

let clients = {};

wss.on("connection", (ws) => {
  const clientID = randomUUID();
  clients[clientID] = ws;
  console.log("New client connected:", clientID);

  // Send welcome message
  ws.send(JSON.stringify({ event: "welcome", clientID }));

  // Receive messages from client
  ws.on("message", (msg) => {
    console.log("Received from client:", msg.toString());
  });

  ws.on("close", () => {
    console.log("Client disconnected:", clientID);
    delete clients[clientID];
  });
});

const PORT = 8080;
server.listen(PORT, () => {
  console.log(`Server running on port :${PORT}`);
});
