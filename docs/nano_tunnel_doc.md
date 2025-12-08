# Nano-tunnel Documentation

## Overview

**Nano-tunnel** is a simple tool that allows you to forward requests to your local ports from the internet.  
It enables you to fetch resources running locally on your device from anywhere, using a secure WebSocket-based tunnel.

Think of it as **Postman for your local machine**, but accessible remotely.

---

## Features

- Expose local ports to the internet securely  
- Forward HTTP requests from hosted instance to your device  
- Supports HTTP methods: GET, POST, PUT, DELETE  
- Configure request body, headers, and API paths  
- Lightweight CLI built with Go

---

## Installation

### 1. Build from source

If you want to build the CLI locally:

```bash
# Clone the repository
git clone https://github.com/yourusername/nano-tunnel.git
cd nano-tunnel

# Build the CLI
go build -o nano-tunnel .

# Make it globally available
sudo mv nano-tunnel /usr/local/bin
chmod +x /usr/local/bin/nano-tunnel

# Verify installation
nano-tunnel --help
```

### 2. Install from hosted binary (no Go required)

If you prefer to install directly:

```bash
curl -fsSL https://nano-tunnel.onrender.com/install.sh | bash
```

**Note:** This downloads the prebuilt binary and installs it globally.

---

## Uninstallation

To uninstall Nano-tunnel:

```bash
# Using uninstall script
curl -fsSL https://nano-tunnel.onrender.com/uninstall.sh | bash

# Or manually
sudo rm /usr/local/bin/nano-tunnel
```

---

## Configuration

Before starting the CLI, you need a configuration file.

### Example config file (`your_config_file.json`):

```json
{
  "remote_url": "nano-tunnel.onrender.com"
}
```

- `remote_url` is the URL of your hosted Nano-tunnel server.

---

## Running the CLI

Start Nano-tunnel and connect your device to the remote server:

```bash
nano-tunnel start ./your_config_file.json
```

- After connecting, the CLI prints a **Client ID** in the terminal.
- This **Client ID** identifies your device for the hosted instance.

---

## Using the Hosted Nano-tunnel

1. Open the hosted Nano-tunnel client website.  
2. Enter your **Client ID** in the input field.  
3. Specify the **local port** you want to fetch.  
4. Enter the **API URL path**, **request body**, and **headers**.  
5. Submit the request — it will be forwarded to your local device via WebSockets.  

**Example use cases:**

- Access a local development server (`localhost:8080`) from outside your network  
- Test APIs on your local machine remotely  
- Forward multiple ports securely  

---

## Notes

- Nano-tunnel uses **WebSockets** to maintain a persistent connection between your device and the hosted server.  
- Keep your CLI running while using the hosted interface.  
- This project is ongoing; features may be updated frequently.

---

## Troubleshooting

- **Cannot connect to remote server:**  
  Make sure `remote_url` is correct and reachable.

- **Client ID not showing:**  
  Ensure your local CLI is running and connected.

- **Permission denied installing CLI:**  
  Use `sudo` when copying the binary to `/usr/local/bin`.

---

## Contributing

Nano-tunnel is an open side project. Contributions, bug reports, and feature requests are welcome.

- Clone the repo, make your changes, and submit a pull request.  
- Ensure your code passes `go build` and any tests.

---

## License

Include your license here, for example:

```
MIT License
```

