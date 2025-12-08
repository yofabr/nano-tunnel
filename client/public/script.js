if (window.location.pathname.includes("nano-tunnel")) {
    window.location.href = "/";
}

const responseBox = document.getElementById("jsonResponse");

document.getElementById("addHeaderBtn").onclick = () => {
    const row = document.createElement("div");
    row.className = "header-row";

    row.innerHTML = `
          <input class="header-key" placeholder="Key" />
          <input class="header-value" placeholder="Value" />
          <button class="remove-header">✕</button>
        `;

    row.querySelector(".remove-header").onclick = () => row.remove();
    document.getElementById("headersContainer").appendChild(row);
};

document.getElementById("clearBtn").onclick = () => {
    document.getElementById("clientId").value = "";
    document.getElementById("port").value = "";
    document.getElementById("payload").value = "";
    document.getElementById("headersContainer").innerHTML = "";
    responseBox.textContent = JSON.stringify(
        { message: "Cleared" },
        null,
        2
    );
};

function validateForm() {
    let valid = true;

    document.getElementById("clientIdError").textContent = "";
    document.getElementById("portError").textContent = "";
    document.getElementById("payloadError").textContent = "";

    const clientId = document.getElementById("clientId").value.trim();
    const port = document.getElementById("port").value.trim();
    const payload = document.getElementById("payload").value.trim();

    if (!clientId) {
        document.getElementById("clientIdError").textContent =
            "Client ID is required.";
        valid = false;
    }
    if (!port) {
        document.getElementById("portError").textContent =
            "Port is required.";
        valid = false;
    }
    if (payload) {
        try {
            JSON.parse(payload);
        } catch {
            document.getElementById("payloadError").textContent =
                "Payload must be valid JSON.";
            valid = false;
        }
    }

    return valid;
}

document.getElementById("startBtn").onclick = async () => {
    if (!validateForm()) return;

    const clientID = document.getElementById("clientId").value.trim();
    const port = document.getElementById("port").value.trim();
    const method = document.getElementById("method").value;
    const path = document.getElementById("path").value.trim();

    const headers = {};
    document.querySelectorAll(".header-row").forEach((row) => {
        const key = row.querySelector(".header-key").value.trim();
        const value = row.querySelector(".header-value").value.trim();
        if (key) headers[key] = value;
    });

    let bodyJson = {};
    const rawPayload = document.getElementById("payload").value.trim();
    if (rawPayload) bodyJson = JSON.parse(rawPayload);

    const requestBody = {
        clientID,
        port,
        path, // <- added path here
        method,
        headers,
        body: bodyJson,
    };

    try {
        responseBox.textContent = "Sending request...";

        const res = await fetch("/send", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(requestBody),
        });

        const data = await res.json();
        responseBox.textContent = JSON.stringify(data, null, 2);
    } catch (err) {
        responseBox.textContent = JSON.stringify(
            {
                error: true,
                message: "Request Failed",
                details: err.message,
            },
            null,
            2
        );
    }
};
