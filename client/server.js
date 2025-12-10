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
  console.log(req.headers, req.body);

  res.json({ message: "Hello from Express!" });
});

app.post("/api/hello", (req, res) => {
  console.log(req.headers, req.body);

  res.json({
    message: "Hello from Express!", data: {
      firstName: "Jane",
      lastName: "Doe",
      age: 30,
      email: "jane.doe@example.com",
      isStudent: false,
      hobbies: ["reading", "hiking", "cooking"],
      address: {
        street: "123 Main St",
        city: "Anytown",
        zipCode: "12345"
      },
      greet: function () {
        return "Hello, my name is " + this.firstName + " " + this.lastName + ".";
      }
    }
  });
});

let clients = {};

app.get("/", (req, res) => {
  res.render("dashboard");
});


const wsSender = (client, data) => {
  client.send(
    JSON.stringify({
      event: "forward",
      message: ``,
      data
    })
  );

  return new Promise((resolve) => {
    client.on("message", (msg) => {
      const response = JSON.parse(msg.toString());
      if (response.event === "response") {
        resolve(response);
      }
    });
  });
}

app.post("/send", async (req, res) => {
  const { clientID, port, method, headers, body, path } = req.body;

  if (!clientID) {
    return res.status(400).json({ error: "clientID is required" });
  }

  const client = clients[clientID]

  if (!client || client.readyState !== 1) {
    return res.status(404).json({ error: "Client not connected" });
  }
  const response = await wsSender(client, {
    local_port: port,
    method,
    headers,
    body,
    path
  });
  // Send a message to this specific client
  res.json(response);
});


// Create WebSocket server
const wss = new WebSocketServer({ server, path: "/ws" });

wss.on("connection", (ws) => {
  const clientID = randomUUID();
  clients[clientID] = ws;
  // console.log("New client connected:", clientID);

  // Send welcome message
  ws.send(JSON.stringify({ event: "welcome", clientID }));

  ws.on("close", () => {
    // console.log("Client disconnected:", clientID);
    delete clients[clientID];
  });
});

const PORT = 8080;
server.listen(PORT, () => {
  // console.log(`Server running on port :${PORT}`);
});
