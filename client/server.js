import express from "express";
import http from "http";
import { WebSocketServer } from "ws";
import { randomUUID } from "crypto";

const app = express();
const server = http.createServer(app);

app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Serve static files from the "public" directory
app.use(express.static("public"));
app.set("view engine", "ejs");
app.set('views', './views');

app.get("/api/hello", (req, res) => {
  res.json({ message: "Hello from Express!" });
});

let clients = {};

app.get("/dashboard", (req, res) => {
  res.render("dashboard", { clients: Object.keys(clients) });
});

app.post("/send", (req, res) => {
  console.log("Sending...", req.body);
  
  const { clientID, port, method, headers, body } = req.body;

  if (!clientID) {
    return res.status(400).json({ error: "clientID is required" });
  }

  const client = clients[clientID]

  if (!client || client.readyState !== 1) {
    return res.status(404).json({ error: "Client not connected" });
  }

  // Send a message to this specific client
  client.send(
    JSON.stringify({
      event: "forward",
      message: ``,
      data: {
        local_port: port,
        method, 
        headers,
        body
      }
    })
  );

  res.json({ message: "Message sent to client", clientID });
});


// Create WebSocket server
const wss = new WebSocketServer({ server, path: "/ws" });


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
