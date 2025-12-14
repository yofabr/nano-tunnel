if (window.location.pathname.includes("nano-tunnel")) {
  window.location.href = "/"
}

const responseBox = document.getElementById("jsonResponse")
const statusIndicator = document.getElementById("statusIndicator")
const statusText = statusIndicator.querySelector(".status-text")
const responseMeta = document.getElementById("responseMeta")
const responseHeadersSection = document.getElementById("responseHeadersSection")

const resizer = document.getElementById("resizer")
const leftPanel = document.getElementById("leftPanel")
const rightPanel = document.getElementById("rightPanel")

let isResizing = false

resizer.addEventListener("mousedown", (e) => {
  isResizing = true
  resizer.classList.add("dragging")
  document.body.style.cursor = "col-resize"
  document.body.style.userSelect = "none"
})

document.addEventListener("mousemove", (e) => {
  if (!isResizing) return

  const containerWidth = document.body.clientWidth
  const newLeftWidth = e.clientX

  // Constrain between min and max widths
  if (newLeftWidth >= 280 && newLeftWidth <= 600) {
    leftPanel.style.width = newLeftWidth + "px"
  }
})

document.addEventListener("mouseup", () => {
  if (isResizing) {
    isResizing = false
    resizer.classList.remove("dragging")
    document.body.style.cursor = ""
    document.body.style.userSelect = ""
  }
})

function syntaxHighlight(json) {
  json = json.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;")

  return json.replace(
    /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+-]?\d+)?)/g,
    (match) => {
      let cls = "number"
      if (/^"/.test(match)) {
        cls = /:$/.test(match) ? "key" : "string"
      } else if (/true|false/.test(match)) {
        cls = "boolean"
      } else if (/null/.test(match)) {
        cls = "null"
      }
      return `<span class="${cls}">${match}</span>`
    },
  )
}

function setStatus(status, text) {
  statusIndicator.className = `status-indicator status-${status}`
  statusText.textContent = text
}

function getStatusClass(code) {
  if (code >= 200 && code < 300) return "success"
  if (code >= 300 && code < 400) return "redirect"
  if (code >= 400 && code < 500) return "client-error"
  return "server-error"
}

function formatBytes(bytes) {
  if (bytes === 0) return "0 B"
  const k = 1024
  const sizes = ["B", "KB", "MB"]
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Number.parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i]
}

function displayResponseMeta(data) {
  responseMeta.classList.remove("hidden")
  responseHeadersSection.classList.remove("hidden")

  // Status code
  const statusCodeEl = document.getElementById("statusCode")
  statusCodeEl.textContent = data.status_code
  statusCodeEl.className = `meta-value status-code ${getStatusClass(data.status_code)}`

  // Time
  document.getElementById("responseTime").textContent = data.time_string || data.time_ms + "ms"

  // Size
  const contentLength = data.headers?.["Content-Length"]?.[0] || 0
  document.getElementById("responseSize").textContent = formatBytes(Number.parseInt(contentLength))

  // Content-Type
  const contentType = data.headers?.["Content-Type"]?.[0] || "-"
  document.getElementById("contentType").textContent = contentType.split(";")[0]

  // Headers
  const headersList = document.getElementById("responseHeadersList")
  const headersCount = document.getElementById("headersCount")
  const headers = data.headers || {}
  const headerKeys = Object.keys(headers)

  headersCount.textContent = headerKeys.length
  headersList.innerHTML = headerKeys
    .map(
      (key) => `
    <div class="header-item">
      <span class="header-key">${key}</span>
      <span class="header-value">${Array.isArray(headers[key]) ? headers[key].join(", ") : headers[key]}</span>
    </div>
  `,
    )
    .join("")
}

document.getElementById("toggleHeaders").onclick = () => {
  const list = document.getElementById("responseHeadersList")
  const chevron = document.getElementById("headersChevron")
  list.classList.toggle("collapsed")
  chevron.classList.toggle("rotated")
}

document.getElementById("copyBtn").onclick = async () => {
  const btn = document.getElementById("copyBtn")
  try {
    await navigator.clipboard.writeText(responseBox.textContent)
    btn.innerHTML = `
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="20 6 9 17 4 12"/>
      </svg>
      Copied!
    `
    btn.classList.add("copied")
    setTimeout(() => {
      btn.innerHTML = `
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
        </svg>
        Copy
      `
      btn.classList.remove("copied")
    }, 2000)
  } catch (err) {
    console.error("Failed to copy:", err)
  }
}

document.getElementById("addHeaderBtn").onclick = () => {
  const row = document.createElement("div")
  row.className = "header-row"

  row.innerHTML = `
    <input class="header-key" placeholder="Key" />
    <input class="header-value" placeholder="Value" />
    <button class="remove-header">✕</button>
  `

  row.querySelector(".remove-header").onclick = () => {
    row.style.opacity = "0"
    row.style.transform = "translateY(-8px)"
    setTimeout(() => row.remove(), 200)
  }
  document.getElementById("headersContainer").appendChild(row)
}

document.getElementById("clearBtn").onclick = () => {
  document.getElementById("clientId").value = ""
  document.getElementById("port").value = ""
  document.getElementById("path").value = ""
  document.getElementById("payload").value = ""
  document.getElementById("headersContainer").innerHTML = ""

  responseMeta.classList.add("hidden")
  responseHeadersSection.classList.add("hidden")
  setStatus("idle", "Ready")

  responseBox.innerHTML = syntaxHighlight(JSON.stringify({ message: "Cleared" }, null, 2))
}

function validateForm() {
  let valid = true

  document.getElementById("clientIdError").textContent = ""
  document.getElementById("portError").textContent = ""
  document.getElementById("payloadError").textContent = ""

  const clientId = document.getElementById("clientId").value.trim()
  const port = document.getElementById("port").value.trim()
  const payload = document.getElementById("payload").value.trim()

  if (!clientId) {
    document.getElementById("clientIdError").textContent = "Client ID is required."
    valid = false
  }
  if (!port) {
    document.getElementById("portError").textContent = "Port is required."
    valid = false
  }
  if (payload) {
    try {
      JSON.parse(payload)
    } catch {
      document.getElementById("payloadError").textContent = "Payload must be valid JSON."
      valid = false
    }
  }

  return valid
}

document.getElementById("startBtn").onclick = async () => {
  if (!validateForm()) return

  const clientID = document.getElementById("clientId").value.trim()
  const port = document.getElementById("port").value.trim()
  const method = document.getElementById("method").value
  const path = document.getElementById("path").value.trim()

  const headers = {}
  document.querySelectorAll(".header-row").forEach((row) => {
    const key = row.querySelector(".header-key").value.trim()
    const value = row.querySelector(".header-value").value.trim()
    if (key) headers[key] = value
  })

  let bodyJson = {}
  const rawPayload = document.getElementById("payload").value.trim()
  if (rawPayload) bodyJson = JSON.parse(rawPayload)

  const requestBody = {
    clientID,
    port,
    path,
    method,
    headers,
    body: bodyJson,
  }

  try {
    setStatus("loading", "Sending...")
    responseMeta.classList.add("hidden")
    responseHeadersSection.classList.add("hidden")
    responseBox.innerHTML = `<span class="loading-dots">Sending request</span>`

    const res = await fetch("/send", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(requestBody),
    })

    const responseData = await res.json()
    const data = responseData.data || responseData

    displayResponseMeta(responseData)

    let responseBody
    try {
      responseBody = JSON.parse(data.response)
    } catch {
      responseBody = data.response
    }

    const json = JSON.stringify(responseBody, null, 2)
    responseBox.innerHTML = syntaxHighlight(json)

    if (responseData.status_code >= 200 && responseData.status_code < 400) {
      setStatus("success", `${responseData.status_code} OK`)
    } else {
      setStatus("error", `${responseData.status_code} Error`)
    }
  } catch (err) {
    setStatus("error", "Failed")
    responseMeta.classList.add("hidden")
    responseHeadersSection.classList.add("hidden")

    const errorJson = JSON.stringify(
      {
        error: true,
        message: "Request Failed",
        details: err.message,
      },
      null,
      2,
    )

    responseBox.innerHTML = syntaxHighlight(errorJson)
  }
}
