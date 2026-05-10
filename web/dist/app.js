const storageKey = "clipbridge.webui.token"

const state = {
  token: localStorage.getItem(storageKey) || "",
  currentPage: "home",
  quickTab: "text",
  latest: null,
  history: [],
  favorites: [],
  devices: [],
  categories: [],
  appSettings: null,
  cleanupStatus: null,
  storageStatus: null,
  webdavStatus: null,
  health: null,
  historySearch: "",
  historyTypeFilter: "",
  historyCategoryFilter: "",
  historyFavoritesOnly: false,
  historyFocusID: null,
  linkSuggestion: { category: "link", tags: ["link"] },
  fileSuggestion: { category: "file", tags: ["file"] },
}

const nodes = {
  tokenInput: document.getElementById("token-input"),
  tokenStatus: document.getElementById("token-status"),
  serverOrigin: document.getElementById("server-origin"),
  healthState: document.getElementById("health-state"),
  connectedDevicesCount: document.getElementById("connected-devices-count"),
  storageTotalBytes: document.getElementById("storage-total-bytes"),
  homeGuide: document.getElementById("home-guide"),
  noticeList: document.getElementById("notice-list"),
  quickInput: document.getElementById("quick-clipboard-input"),
  quickLinkInput: document.getElementById("quick-link-input"),
  quickLinkCategory: document.getElementById("quick-link-category"),
  quickLinkTags: document.getElementById("quick-link-tags"),
  quickFileInput: document.getElementById("quick-file-input"),
  quickFileCategory: document.getElementById("quick-file-category"),
  quickFileTags: document.getElementById("quick-file-tags"),
  quickFilePreview: document.getElementById("quick-file-preview"),
  latestEmpty: document.getElementById("latest-empty"),
  latestCard: document.getElementById("latest-card"),
  latestType: document.getElementById("latest-type"),
  latestCategory: document.getElementById("latest-category"),
  latestCreated: document.getElementById("latest-created"),
  latestSourceName: document.getElementById("latest-source-name"),
  latestContent: document.getElementById("latest-content"),
  latestActions: document.getElementById("latest-actions"),
  historySearchInput: document.getElementById("history-search-input"),
  historyTypeFilter: document.getElementById("history-type-filter"),
  historyCategoryFilter: document.getElementById("history-category-filter"),
  historyFavoritesOnly: document.getElementById("history-favorites-only"),
  historyList: document.getElementById("history-list"),
  favoritesList: document.getElementById("favorites-list"),
  devicesList: document.getElementById("devices-list"),
  pairingOutput: document.getElementById("pairing-output"),
  pairingCodeValue: document.getElementById("pairing-code-value"),
  pairingCodeExpiry: document.getElementById("pairing-code-expiry"),
  pairingCodeURI: document.getElementById("pairing-code-uri"),
  newCategoryInput: document.getElementById("new-category-input"),
  adminTokenInput: document.getElementById("admin-token-input"),
  startupSettings: document.getElementById("startup-settings"),
  limitMinText: document.getElementById("limit-min-text"),
  limitMaxText: document.getElementById("limit-max-text"),
  limitMinImage: document.getElementById("limit-min-image"),
  limitMaxImage: document.getElementById("limit-max-image"),
  limitMinFile: document.getElementById("limit-min-file"),
  limitMaxFile: document.getElementById("limit-max-file"),
  limitMinLink: document.getElementById("limit-min-link"),
  limitMaxLink: document.getElementById("limit-max-link"),
  limitMaxRequest: document.getElementById("limit-max-request"),
  settingsHistoryCount: document.getElementById("settings-history-count"),
  settingsFavoriteCount: document.getElementById("settings-favorite-count"),
  settingsTotalBytes: document.getElementById("settings-total-bytes"),
  settingsFileBytes: document.getElementById("settings-file-bytes"),
  cleanupTTLHours: document.getElementById("cleanup-ttl-hours"),
  cleanupMaxItems: document.getElementById("cleanup-max-items"),
  cleanupMaxSizeMB: document.getElementById("cleanup-max-size-mb"),
  cleanupIntervalMinutes: document.getElementById("cleanup-interval-minutes"),
  cleanupEnabled: document.getElementById("cleanup-enabled"),
  cleanupLastRun: document.getElementById("cleanup-last-run"),
  cleanupLastResult: document.getElementById("cleanup-last-result"),
  webdavURL: document.getElementById("webdav-url"),
  webdavUsername: document.getElementById("webdav-username"),
  webdavPassword: document.getElementById("webdav-password"),
  webdavBasePath: document.getElementById("webdav-base-path"),
  webdavEnabled: document.getElementById("webdav-enabled"),
  webdavTestedAt: document.getElementById("webdav-tested-at"),
  webdavTestResult: document.getElementById("webdav-test-result"),
  webdavLastSync: document.getElementById("webdav-last-sync"),
  webdavLastSuccess: document.getElementById("webdav-last-success"),
  webdavSyncResult: document.getElementById("webdav-sync-result"),
  webdavTransferSummary: document.getElementById("webdav-transfer-summary"),
  pageButtons: Array.from(document.querySelectorAll("[data-page-target]")),
  pages: Array.from(document.querySelectorAll("[data-page]")),
  quickTabButtons: Array.from(document.querySelectorAll("[data-quick-tab]")),
  quickPanels: Array.from(document.querySelectorAll("[data-quick-panel]")),
}

nodes.tokenInput.value = state.token
nodes.serverOrigin.textContent = window.location.origin

document.getElementById("save-token-button").addEventListener("click", saveToken)
document.getElementById("clear-token-button").addEventListener("click", clearToken)
document.getElementById("refresh-home-button").addEventListener("click", refreshHome)
document.getElementById("refresh-latest-button").addEventListener("click", loadLatest)
document.getElementById("upload-quick-button").addEventListener("click", uploadQuickText)
document.getElementById("upload-link-button").addEventListener("click", uploadQuickLink)
document.getElementById("upload-file-button").addEventListener("click", uploadQuickFile)
document.getElementById("refresh-history-button").addEventListener("click", loadHistory)
document.getElementById("refresh-favorites-button").addEventListener("click", loadFavorites)
document.getElementById("refresh-devices-button").addEventListener("click", loadDevices)
document.getElementById("generate-pairing-button").addEventListener("click", generatePairingCode)
document.getElementById("create-category-button").addEventListener("click", createCategory)
document.getElementById("save-admin-token-button").addEventListener("click", saveAdminToken)
document.getElementById("save-limits-button").addEventListener("click", saveLimits)
document.getElementById("refresh-cleanup-button").addEventListener("click", loadCleanupData)
document.getElementById("run-cleanup-button").addEventListener("click", runCleanupNow)
document.getElementById("save-cleanup-settings-button").addEventListener("click", saveCleanupSettings)
document.getElementById("save-webdav-settings-button").addEventListener("click", saveWebDAVSettings)
document.getElementById("refresh-webdav-button").addEventListener("click", loadWebDAVData)
document.getElementById("test-webdav-button").addEventListener("click", testWebDAVConnection)
document.getElementById("sync-webdav-button").addEventListener("click", runWebDAVSync)
document.getElementById("clear-notices-button").addEventListener("click", clearNotices)

nodes.pageButtons.forEach((button) => {
  button.addEventListener("click", () => setPage(button.dataset.pageTarget || "home"))
})

nodes.quickTabButtons.forEach((button) => {
  button.addEventListener("click", () => setQuickTab(button.dataset.quickTab || "text"))
})

nodes.quickLinkInput.addEventListener("input", updateLinkSuggestion)
nodes.quickFileInput.addEventListener("change", updateFileSuggestion)
nodes.historySearchInput.addEventListener("input", (event) => {
  state.historySearch = event.target.value.trim().toLowerCase()
  renderHistory()
})
nodes.historyTypeFilter.addEventListener("change", (event) => {
  state.historyTypeFilter = event.target.value
  renderHistory()
})
nodes.historyCategoryFilter.addEventListener("change", (event) => {
  state.historyCategoryFilter = event.target.value
  renderHistory()
})
nodes.historyFavoritesOnly.addEventListener("change", (event) => {
  state.historyFavoritesOnly = event.target.checked
  renderHistory()
})

updateTokenStatus()
setPage("home")
setQuickTab("text")
renderHomeGuide()
renderLatest()
renderHistory()
renderFavorites()
renderDevices([])
renderAppSettings()
renderCleanupStatus()
renderCleanupSettings()
renderWebDAVSettings()
renderWebDAVStatus()
updateLinkSuggestion()
updateFileSuggestion()
loadHealth()
loadAllProtectedData()

function saveToken() {
  state.token = nodes.tokenInput.value.trim()
  persistToken()
  updateTokenStatus()
  renderHomeGuide()
  loadAllProtectedData()
}

function clearToken() {
  state.token = ""
  nodes.tokenInput.value = ""
  persistToken()
  updateTokenStatus()
  renderHomeGuide()
  loadAllProtectedData()
}

function persistToken() {
  if (state.token) {
    localStorage.setItem(storageKey, state.token)
    addNotice("Saved token locally for this browser.", "success")
  } else {
    localStorage.removeItem(storageKey)
    addNotice("Cleared the saved token.", "info")
  }
}

function updateTokenStatus() {
  nodes.tokenStatus.textContent = state.token
    ? "Saved token will be used for API requests."
    : "No token saved yet."
}

function setPage(page) {
  state.currentPage = page
  nodes.pageButtons.forEach((button) => {
    button.classList.toggle("active", button.dataset.pageTarget === page)
  })
  nodes.pages.forEach((section) => {
    section.classList.toggle("active", section.dataset.page === page)
  })
}

function setQuickTab(tab) {
  state.quickTab = tab
  nodes.quickTabButtons.forEach((button) => {
    button.classList.toggle("active", button.dataset.quickTab === tab)
  })
  nodes.quickPanels.forEach((panel) => {
    panel.classList.toggle("active", panel.dataset.quickPanel === tab)
  })
}

function renderHomeGuide() {
  if (!state.token) {
    nodes.homeGuide.textContent = "Save an admin token or device token first. Then upload a new item from Text, Link, or File."
    return
  }
  nodes.homeGuide.textContent = "Use Quick Clipboard to upload something new, then open History for favorite, delete, and category management."
}

async function refreshHome() {
  await Promise.all([loadHealth(), loadLatest(), loadDevices(), loadCleanupData()])
}

async function loadAllProtectedData() {
  await Promise.all([
    loadAppSettings(),
    loadCategories(),
    loadLatest(),
    loadHistory(),
    loadFavorites(),
    loadDevices(),
    loadCleanupData(),
    loadWebDAVData(),
  ])
}

async function loadHealth() {
  try {
    const response = await apiFetch("/api/health", { auth: false })
    state.health = response.data
  } catch (error) {
    state.health = null
    addNotice(error.message, "error")
  }
  renderHealth()
}

function renderHealth() {
  const isOnline = Boolean(state.health?.ok)
  nodes.healthState.textContent = isOnline ? "Online" : "Offline"
  nodes.healthState.classList.toggle("is-offline", !isOnline)
}

async function loadAppSettings() {
  if (!state.token) {
    state.appSettings = null
    renderAppSettings()
    renderCleanupSettings()
    renderWebDAVSettings()
    renderLimits()
    return
  }

  try {
    const response = await apiFetch("/api/settings")
    state.appSettings = response.data
  } catch (error) {
    state.appSettings = null
    if (!shouldSilenceAdminError(error)) {
      addNotice(error.message, "error")
    }
  }

  renderAppSettings()
  renderCleanupSettings()
  renderWebDAVSettings()
  renderLimits()
}

async function loadCategories() {
  if (!state.token) {
    state.categories = []
    renderCategoryControls()
    return
  }

  try {
    const response = await apiFetch("/api/categories")
    state.categories = response.data.items || []
  } catch (error) {
    state.categories = []
    addNotice(error.message, "error")
  }

  renderCategoryControls()
  updateLinkSuggestion()
  updateFileSuggestion()
}

async function loadLatest() {
  if (!state.token) {
    state.latest = null
    renderLatest()
    return
  }

  try {
    const response = await apiFetch("/api/clipboard/latest")
    state.latest = response.data
  } catch (error) {
    state.latest = null
    if (!String(error.message).includes("no clipboard item found")) {
      addNotice(error.message, "error")
    }
  }

  renderLatest()
}

async function loadHistory() {
  if (!state.token) {
    state.history = []
    renderHistory()
    return
  }

  try {
    const response = await apiFetch("/api/clipboard/history")
    state.history = response.data.items || []
  } catch (error) {
    state.history = []
    addNotice(error.message, "error")
  }

  renderHistory()
}

async function loadFavorites() {
  if (!state.token) {
    state.favorites = []
    renderFavorites()
    return
  }

  try {
    const response = await apiFetch("/api/favorites")
    state.favorites = response.data.items || []
  } catch (error) {
    state.favorites = []
    addNotice(error.message, "error")
  }

  renderFavorites()
}

async function loadDevices() {
  if (!state.token) {
    state.devices = []
    renderDevices([])
    renderDeviceSummary()
    return
  }

  try {
    const response = await apiFetch("/api/auth/devices")
    state.devices = response.data.items || []
  } catch (error) {
    state.devices = []
    if (!shouldSilenceAdminError(error)) {
      addNotice(error.message, "error")
    }
  }

  renderDevices(state.devices)
  renderDeviceSummary()
}

async function loadCleanupData() {
  if (!state.token) {
    state.cleanupStatus = null
    state.storageStatus = null
    renderCleanupStatus()
    return
  }

  try {
    const [statusResponse, storageResponse] = await Promise.all([
      apiFetch("/api/admin/cleanup/status"),
      apiFetch("/api/admin/storage/status"),
    ])
    state.cleanupStatus = statusResponse.data
    state.storageStatus = storageResponse.data
  } catch (error) {
    state.cleanupStatus = null
    state.storageStatus = null
    if (!shouldSilenceAdminError(error)) {
      addNotice(error.message, "error")
    }
  }

  renderCleanupStatus()
}

async function loadWebDAVData() {
  if (!state.token) {
    state.webdavStatus = null
    renderWebDAVStatus()
    return
  }

  try {
    const response = await apiFetch("/api/admin/webdav/status")
    state.webdavStatus = response.data
  } catch (error) {
    state.webdavStatus = null
    if (!shouldSilenceAdminError(error)) {
      addNotice(error.message, "error")
    }
  }

  renderWebDAVStatus()
}

async function uploadQuickText() {
  const content = nodes.quickInput.value.trim()
  if (!content) {
    addNotice("Enter some text before uploading.", "info")
    return
  }

  try {
    await apiFetch("/api/clipboard/text", {
      method: "POST",
      body: JSON.stringify({
        content,
        source_device_id: "web-ui",
        source_device_name: "Web UI",
      }),
    })
    nodes.quickInput.value = ""
    addNotice("Uploaded text to the clipboard stack.", "success")
    await refreshClipboardViews()
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function uploadQuickLink() {
  const url = nodes.quickLinkInput.value.trim()
  if (!url) {
    addNotice("Enter a link before uploading.", "info")
    return
  }

  try {
    const response = await apiFetch("/api/clipboard/link", {
      method: "POST",
      body: JSON.stringify({
        url,
        source_device_id: "web-ui",
        source_device_name: "Web UI",
      }),
    })
    await applySuggestedCategory(response.data, nodes.quickLinkCategory.value)
    nodes.quickLinkInput.value = ""
    updateLinkSuggestion()
    addNotice("Uploaded link clipboard item.", "success")
    await refreshClipboardViews()
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function uploadQuickFile() {
  const file = nodes.quickFileInput.files[0]
  if (!file) {
    addNotice("Choose one file first.", "info")
    return
  }

  const formData = new FormData()
  formData.append("file", file)
  formData.append("source_device_id", "web-ui")
  formData.append("source_device_name", "Web UI")

  try {
    const response = await apiFetch("/api/clipboard/file", { method: "POST", body: formData })
    await applySuggestedCategory(response.data, nodes.quickFileCategory.value)
    nodes.quickFileInput.value = ""
    updateFileSuggestion()
    addNotice("Uploaded image or file clipboard item.", "success")
    await refreshClipboardViews()
  } catch (error) {
    addNotice(error.message, "error")
  }
}

function updateLinkSuggestion() {
  const url = nodes.quickLinkInput.value.trim()
  state.linkSuggestion = classifyUrl(url)
  renderQuickLinkSuggestion()
}

function updateFileSuggestion() {
  const file = nodes.quickFileInput.files[0]
  state.fileSuggestion = classifyFile(file)
  renderQuickFileSuggestion(file)
}

function renderQuickLinkSuggestion() {
  renderTagList(nodes.quickLinkTags, state.linkSuggestion.tags)
  populateCategorySelect(nodes.quickLinkCategory, state.linkSuggestion.category, [state.linkSuggestion.category])
}

function renderQuickFileSuggestion(file) {
  renderTagList(nodes.quickFileTags, state.fileSuggestion.tags)
  populateCategorySelect(nodes.quickFileCategory, state.fileSuggestion.category, [state.fileSuggestion.category])
  nodes.quickFilePreview.innerHTML = ""

  if (!file) {
    nodes.quickFilePreview.classList.add("hidden")
    return
  }

  nodes.quickFilePreview.classList.remove("hidden")
  const summary = document.createElement("p")
  summary.className = "muted"
  summary.textContent = `${file.name} · ${formatBytes(file.size)} · ${file.type || "unknown type"}`

  if (file.type.startsWith("image/")) {
    const preview = document.createElement("img")
    preview.alt = file.name
    preview.src = URL.createObjectURL(file)
    nodes.quickFilePreview.append(preview, summary)
    return
  }

  nodes.quickFilePreview.appendChild(summary)
}

async function applySuggestedCategory(item, selectedCategory) {
  const category = normalizeCategoryName(selectedCategory)
  if (!item || !item.id || !category || category === normalizeCategoryName(item.category || "")) {
    return
  }

  await ensureCategoryExists(category)
  await apiFetch(`/api/clipboard/items/${item.id}/category`, {
    method: "PATCH",
    body: JSON.stringify({ category }),
  })
}

async function ensureCategoryExists(category) {
  const normalized = normalizeCategoryName(category)
  if (!normalized) {
    return
  }

  const exists = state.categories.some((entry) => normalizeCategoryName(entry.name) === normalized)
  if (exists) {
    return
  }

  try {
    await apiFetch("/api/categories", {
      method: "POST",
      body: JSON.stringify({ name: normalized }),
    })
  } catch (error) {
    if (!String(error.message).includes("already exists")) {
      throw error
    }
  }

  await loadCategories()
}

async function saveAdminToken() {
  const token = nodes.adminTokenInput.value.trim()
  if (!token) {
    addNotice("Enter a new admin token first.", "info")
    return
  }

  try {
    await apiFetch("/api/settings", {
      method: "PATCH",
      body: JSON.stringify({ admin_token: token }),
    })
    state.token = token
    nodes.tokenInput.value = token
    persistToken()
    updateTokenStatus()
    await loadAllProtectedData()
    addNotice("Updated admin token and switched this browser to the new token.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function saveLimits() {
  const payload = {
    min_text_bytes: Number(nodes.limitMinText.value),
    max_text_bytes: Number(nodes.limitMaxText.value),
    min_image_bytes: Number(nodes.limitMinImage.value),
    max_image_bytes: Number(nodes.limitMaxImage.value),
    min_file_bytes: Number(nodes.limitMinFile.value),
    max_file_bytes: Number(nodes.limitMaxFile.value),
    min_link_bytes: Number(nodes.limitMinLink.value),
    max_link_bytes: Number(nodes.limitMaxLink.value),
    max_request_bytes: Number(nodes.limitMaxRequest.value),
  }

  try {
    await apiFetch("/api/settings/limits", {
      method: "PATCH",
      body: JSON.stringify(payload),
    })
    await loadAppSettings()
    addNotice("Saved runtime limits.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function saveCleanupSettings() {
  const payload = {
    ttl_hours: Number(nodes.cleanupTTLHours.value),
    max_items: Number(nodes.cleanupMaxItems.value),
    max_total_size_mb: Number(nodes.cleanupMaxSizeMB.value),
    interval_minutes: Number(nodes.cleanupIntervalMinutes.value),
    enabled: nodes.cleanupEnabled.checked,
  }

  try {
    const response = await apiFetch("/api/settings/cleanup", {
      method: "PATCH",
      body: JSON.stringify(payload),
    })
    if (state.appSettings) {
      state.appSettings.cleanup = response.data
    }
    renderCleanupSettings()
    addNotice("Saved cleanup policy.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function runCleanupNow() {
  try {
    await apiFetch("/api/admin/cleanup/run", { method: "POST" })
    await Promise.all([loadCleanupData(), loadLatest(), loadHistory(), loadFavorites()])
    addNotice("Triggered one manual cleanup run.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function saveWebDAVSettings() {
  const payload = {
    enabled: nodes.webdavEnabled.checked,
    url: nodes.webdavURL.value.trim(),
    username: nodes.webdavUsername.value.trim(),
    password: nodes.webdavPassword.value,
    base_path: nodes.webdavBasePath.value.trim(),
  }

  try {
    const response = await apiFetch("/api/settings/webdav", {
      method: "PATCH",
      body: JSON.stringify(payload),
    })
    if (state.appSettings) {
      state.appSettings.webdav = response.data
    }
    renderWebDAVSettings()
    addNotice("Saved WebDAV settings.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function testWebDAVConnection() {
  try {
    const response = await apiFetch("/api/admin/webdav/test", { method: "POST" })
    state.webdavStatus = response.data
    renderWebDAVStatus()
    addNotice("WebDAV connection test succeeded.", "success")
  } catch (error) {
    await loadWebDAVData()
    addNotice(error.message, "error")
  }
}

async function runWebDAVSync() {
  try {
    const response = await apiFetch("/api/admin/webdav/sync", { method: "POST" })
    state.webdavStatus = response.data
    renderWebDAVStatus()
    await refreshClipboardViews()
    addNotice("WebDAV sync completed.", "success")
  } catch (error) {
    await loadWebDAVData()
    addNotice(error.message, "error")
  }
}

async function generatePairingCode() {
  try {
    const response = await apiFetch("/api/auth/pairing-codes", { method: "POST" })
    nodes.pairingOutput.classList.remove("hidden")
    nodes.pairingCodeValue.textContent = response.data.pairing_code
    nodes.pairingCodeExpiry.textContent = formatTime(response.data.expires_at)
    nodes.pairingCodeURI.textContent = response.data.pairing_uri
    addNotice("Generated a new pairing code.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function createCategory() {
  const name = normalizeCategoryName(nodes.newCategoryInput.value)
  if (!name) {
    addNotice("Enter a category name first.", "info")
    return
  }

  try {
    await apiFetch("/api/categories", {
      method: "POST",
      body: JSON.stringify({ name }),
    })
    nodes.newCategoryInput.value = ""
    await loadCategories()
    addNotice(`Created category "${name}".`, "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function copyText(value) {
  try {
    await navigator.clipboard.writeText(value || "")
    addNotice("Copied content into the browser clipboard.", "success")
  } catch (error) {
    addNotice("Copy failed. Your browser may block clipboard access.", "error")
  }
}

async function setFavorite(id, favorite) {
  try {
    await apiFetch(`/api/clipboard/items/${id}/favorite`, {
      method: favorite ? "POST" : "DELETE",
    })
    await refreshClipboardViews()
    addNotice(favorite ? "Added item to favorites." : "Removed item from favorites.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function deleteItem(id) {
  try {
    await apiFetch(`/api/clipboard/items/${id}`, { method: "DELETE" })
    if (state.historyFocusID === id) {
      state.historyFocusID = null
    }
    await refreshClipboardViews()
    addNotice("Deleted clipboard item.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function updateItemCategory(id, category) {
  const normalized = normalizeCategoryName(category)
  if (!normalized) {
    addNotice("Choose a category first.", "info")
    return
  }

  try {
    await ensureCategoryExists(normalized)
    await apiFetch(`/api/clipboard/items/${id}/category`, {
      method: "PATCH",
      body: JSON.stringify({ category: normalized }),
    })
    await refreshClipboardViews()
    addNotice("Updated item category.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function revokeDevice(id) {
  try {
    await apiFetch(`/api/auth/devices/${id}`, { method: "DELETE" })
    await loadDevices()
    addNotice("Revoked device access.", "success")
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function refreshClipboardViews() {
  await Promise.all([loadLatest(), loadHistory(), loadFavorites(), loadCleanupData()])
}

function renderLatest() {
  if (!state.latest) {
    nodes.latestCard.classList.add("hidden")
    nodes.latestEmpty.classList.remove("hidden")
    nodes.latestEmpty.textContent = state.token
      ? "No clipboard item loaded yet."
      : "Save an admin token or device token to load clipboard data."
    return
  }

  nodes.latestCard.classList.remove("hidden")
  nodes.latestEmpty.classList.add("hidden")
  nodes.latestType.textContent = itemTypeLabel(state.latest.type)
  nodes.latestCategory.textContent = state.latest.category || "uncategorized"
  nodes.latestCreated.textContent = formatTime(state.latest.created_at)
  nodes.latestSourceName.textContent = state.latest.source_device_name || "Unknown source"
  nodes.latestContent.innerHTML = ""
  nodes.latestContent.appendChild(renderItemContent(state.latest))

  nodes.latestActions.innerHTML = ""
  latestActionsForItem(state.latest).forEach((action) => {
    nodes.latestActions.appendChild(renderActionButton(action))
  })
}

function latestActionsForItem(item) {
  const actions = []
  if (item.type === "text") {
    actions.push({ label: "Copy", onClick: () => copyText(item.text || "") })
  }
  if (item.type === "link") {
    actions.push({ label: "Copy", onClick: () => copyText(item.url || item.text || "") })
    actions.push({ label: "Open", secondary: true, onClick: () => openLinkItem(item) })
  }
  if (item.type === "image") {
    actions.push({ label: "Open", secondary: true, onClick: () => openProtectedItem(item) })
  }
  if (item.type === "file") {
    actions.push({ label: "Download", secondary: true, onClick: () => downloadProtectedFile(item) })
  }
  if (item.type === "image" || item.type === "file" || item.type === "text" || item.type === "link") {
    actions.push({ label: "View in History", secondary: true, onClick: () => goToHistoryItem(item.id) })
  }
  return actions
}

function renderHistory() {
  renderClipboardCollection(
    nodes.historyList,
    getFilteredHistoryItems(),
    state.token ? "No history items match the current filters." : "Save a token to view clipboard history.",
    "history",
  )
}

function renderFavorites() {
  renderClipboardCollection(
    nodes.favoritesList,
    state.favorites,
    state.token ? "No favorites saved yet." : "Save a token to view favorites.",
    "favorites",
  )
}

function renderClipboardCollection(container, items, emptyMessage, mode) {
  container.innerHTML = ""
  if (!items.length) {
    container.appendChild(renderEmpty(emptyMessage))
    return
  }

  items.forEach((item) => {
    container.appendChild(renderClipboardItemCard(item, mode))
  })

  if (mode === "history" && state.historyFocusID) {
    const focusNode = container.querySelector(`[data-item-id="${state.historyFocusID}"]`)
    if (focusNode) {
      focusNode.scrollIntoView({ behavior: "smooth", block: "center" })
    }
  }
}

function renderClipboardItemCard(item, mode) {
  const li = document.createElement("li")
  li.className = "history-item"
  li.dataset.itemId = String(item.id)
  li.classList.toggle("is-focused", state.historyFocusID === item.id)

  const header = document.createElement("div")
  header.className = "clipboard-card-header"

  const titleBlock = document.createElement("div")
  titleBlock.className = "clipboard-card-title"
  const title = document.createElement("h3")
  title.textContent = itemPrimaryTitle(item)
  const subtitle = document.createElement("div")
  subtitle.className = "clipboard-card-subtitle"
  subtitle.textContent = `${formatTime(item.created_at)} · ${item.source_device_name || "Unknown source"}`
  titleBlock.append(title, subtitle)

  const tags = document.createElement("div")
  tags.className = "clipboard-tags"
  tags.append(
    renderTag(itemTypeLabel(item.type)),
    renderTag(item.category || "uncategorized", "tag-accent"),
  )
  if (item.is_favorite) {
    tags.appendChild(renderTag("Favorite", "tag-favorite"))
  }

  header.append(titleBlock, tags)

  const content = renderItemContent(item)
  const actions = document.createElement("div")
  actions.className = "clipboard-actions"
  itemActionsForMode(item, mode).forEach((action) => {
    actions.appendChild(renderActionButton(action))
  })

  li.append(header, content, actions)

  if (mode === "history") {
    const categoryRow = document.createElement("div")
    categoryRow.className = "inline-row"
    const categorySelect = document.createElement("select")
    populateCategorySelect(categorySelect, item.category || "", [item.category || ""])
    const categoryButton = renderActionButton({
      label: "Save Category",
      secondary: true,
      onClick: () => updateItemCategory(item.id, categorySelect.value),
    })
    categoryRow.append(categorySelect, categoryButton)
    li.appendChild(categoryRow)
  }

  li.appendChild(renderDetailsBlock(item))
  return li
}

function itemActionsForMode(item, mode) {
  const actions = []
  if (item.type === "text" || item.type === "link") {
    actions.push({
      label: "Copy",
      onClick: () => copyText(item.type === "link" ? item.url || item.text || "" : item.text || ""),
    })
  }
  if (item.type === "link") {
    actions.push({ label: "Open", secondary: true, onClick: () => openLinkItem(item) })
  }
  if (item.type === "image") {
    actions.push({ label: "Open", secondary: true, onClick: () => openProtectedItem(item) })
  }
  if (item.type === "file") {
    actions.push({ label: "Download", secondary: true, onClick: () => downloadProtectedFile(item) })
  }
  if (mode === "favorites") {
    actions.push({ label: "View in History", secondary: true, onClick: () => goToHistoryItem(item.id) })
    actions.push({ label: "Unfavorite", secondary: true, onClick: () => setFavorite(item.id, false) })
    return actions
  }
  if (mode === "history") {
    actions.push({
      label: item.is_favorite ? "Unfavorite" : "Favorite",
      secondary: true,
      onClick: () => setFavorite(item.id, !item.is_favorite),
    })
    actions.push({ label: "Delete", danger: true, onClick: () => deleteItem(item.id) })
  }
  return actions
}

function renderItemContent(item) {
  if (item.type === "image" && item.preview_url) {
    const wrapper = document.createElement("div")
    wrapper.className = "item-content"
    const preview = document.createElement("img")
    preview.className = "image-preview"
    preview.alt = item.filename || "Clipboard image"
    loadProtectedAssetURL(item.preview_url)
      .then((objectURL) => {
        preview.src = objectURL
      })
      .catch(() => {
        preview.alt = "Image preview requires a valid token."
      })
    const meta = document.createElement("div")
    meta.className = "muted"
    meta.textContent = item.filename || "Image"
    wrapper.append(preview, meta)
    return wrapper
  }

  if (item.type === "file") {
    const wrapper = document.createElement("div")
    wrapper.className = "item-content"
    const filename = document.createElement("strong")
    filename.textContent = item.filename || "File"
    const meta = document.createElement("div")
    meta.className = "muted"
    meta.textContent = [item.mime_type || "application/octet-stream", formatBytes(item.size_bytes || 0)]
      .filter(Boolean)
      .join(" · ")
    wrapper.append(filename, meta)
    return wrapper
  }

  if (item.type === "link") {
    const wrapper = document.createElement("div")
    wrapper.className = "item-content"
    const anchor = document.createElement("a")
    anchor.href = item.url || item.text || "#"
    anchor.target = "_blank"
    anchor.rel = "noreferrer"
    anchor.textContent = item.url || item.text || ""
    wrapper.appendChild(anchor)
    return wrapper
  }

  const pre = document.createElement("pre")
  pre.textContent = item.text || ""
  return pre
}

function renderDetailsBlock(item) {
  const details = document.createElement("details")
  const summary = document.createElement("summary")
  summary.textContent = "More details"
  const content = document.createElement("div")
  content.className = "details-content"

  const rows = [
    ["Type", item.type || "-"],
    ["Category", item.category || "-"],
    ["Source", item.source_device_name || "Unknown source"],
    ["Created", formatTime(item.created_at)],
    ["Updated", formatTime(item.updated_at)],
  ]

  if (item.filename) {
    rows.push(["Filename", item.filename])
  }
  if (item.mime_type) {
    rows.push(["MIME", item.mime_type])
  }
  if (item.size_bytes) {
    rows.push(["Size", formatBytes(item.size_bytes)])
  }
  if (item.url) {
    rows.push(["URL", item.url])
  }

  rows.forEach(([label, value]) => {
    const line = document.createElement("div")
    line.textContent = `${label}: ${value}`
    content.appendChild(line)
  })

  details.append(summary, content)
  return details
}

function renderDevices(items) {
  nodes.devicesList.innerHTML = ""
  if (!items.length) {
    nodes.devicesList.appendChild(
      renderEmpty(state.token ? "No devices visible with the current token." : "Save the admin token to manage devices."),
    )
    return
  }

  items.forEach((device) => {
    const li = document.createElement("li")
    li.className = "device-item"
    const title = document.createElement("strong")
    title.textContent = device.name || "Unnamed device"
    const created = document.createElement("div")
    created.className = "muted"
    created.textContent = `Created: ${formatTime(device.created_at)}`
    const lastSeen = document.createElement("div")
    lastSeen.className = "muted"
    lastSeen.textContent = `Last seen: ${device.last_seen_at ? formatTime(device.last_seen_at) : "Never"}`
    li.append(title, created, lastSeen)

    if (!device.revoked_at) {
      li.appendChild(
        renderActionButton({
          label: "Revoke",
          danger: true,
          onClick: () => revokeDevice(device.id),
        }),
      )
    } else {
      const revoked = document.createElement("div")
      revoked.className = "muted"
      revoked.textContent = `Revoked: ${formatTime(device.revoked_at)}`
      li.appendChild(revoked)
    }

    nodes.devicesList.appendChild(li)
  })
}

function renderDeviceSummary() {
  const activeCount = state.devices.filter((device) => !device.revoked_at).length
  nodes.connectedDevicesCount.textContent = state.devices.length ? String(activeCount) : "-"
}

function renderAppSettings() {
  const appSettings = state.appSettings
  nodes.adminTokenInput.value = appSettings?.admin_token || ""

  if (!appSettings) {
    nodes.startupSettings.textContent = "Admin token required."
    return
  }

  const startup = appSettings.startup || {}
  const restartFields = (appSettings.restart_required_fields || []).join(", ")
  nodes.startupSettings.textContent =
    `Host: ${startup.host} · Port: ${startup.port} · Data Dir: ${startup.data_dir} · DB: ${startup.database_path} · Restart required for: ${restartFields}`
}

function renderLimits() {
  const limits = state.appSettings?.limits
  nodes.limitMinText.value = limits?.min_text_bytes ?? ""
  nodes.limitMaxText.value = limits?.max_text_bytes ?? ""
  nodes.limitMinImage.value = limits?.min_image_bytes ?? ""
  nodes.limitMaxImage.value = limits?.max_image_bytes ?? ""
  nodes.limitMinFile.value = limits?.min_file_bytes ?? ""
  nodes.limitMaxFile.value = limits?.max_file_bytes ?? ""
  nodes.limitMinLink.value = limits?.min_link_bytes ?? ""
  nodes.limitMaxLink.value = limits?.max_link_bytes ?? ""
  nodes.limitMaxRequest.value = limits?.max_request_bytes ?? ""
}

function renderCleanupStatus() {
  const status = state.cleanupStatus
  const storage = state.storageStatus

  nodes.storageTotalBytes.textContent = storage ? formatBytes(storage.total_bytes) : "-"
  nodes.settingsHistoryCount.textContent = storage ? String(storage.history_count) : "-"
  nodes.settingsFavoriteCount.textContent = storage ? String(storage.favorite_count) : "-"
  nodes.settingsTotalBytes.textContent = storage ? formatBytes(storage.total_bytes) : "-"
  nodes.settingsFileBytes.textContent = storage ? formatBytes(storage.file_bytes) : "-"
  nodes.cleanupLastRun.textContent = status?.last_run_at ? formatTime(status.last_run_at) : "-"
  nodes.cleanupLastResult.textContent = status
    ? status.last_error || status.last_message || `expired ${status.deleted_expired || 0} · count ${status.deleted_max_items || 0} · storage ${status.deleted_storage || 0}`
    : "Admin token required"
}

function renderCleanupSettings() {
  const cleanup = state.appSettings?.cleanup
  nodes.cleanupTTLHours.value = cleanup?.ttl_hours ?? ""
  nodes.cleanupMaxItems.value = cleanup?.max_items ?? ""
  nodes.cleanupMaxSizeMB.value = cleanup?.max_total_size_mb ?? ""
  nodes.cleanupIntervalMinutes.value = cleanup?.interval_minutes ?? ""
  nodes.cleanupEnabled.checked = Boolean(cleanup?.enabled)
}

function renderWebDAVSettings() {
  const webdav = state.appSettings?.webdav
  nodes.webdavURL.value = webdav?.url ?? ""
  nodes.webdavUsername.value = webdav?.username ?? ""
  nodes.webdavPassword.value = webdav?.password ?? ""
  nodes.webdavBasePath.value = webdav?.base_path ?? "ClipBridgeServer"
  nodes.webdavEnabled.checked = Boolean(webdav?.enabled)
}

function renderWebDAVStatus() {
  const status = state.webdavStatus
  nodes.webdavTestedAt.textContent = status?.tested_at ? formatTime(status.tested_at) : "-"
  nodes.webdavTestResult.textContent = status
    ? status.last_test_success
      ? "OK"
      : status.last_test_error || "Failed"
    : "Admin token required"
  nodes.webdavLastSync.textContent = status?.last_sync_at ? formatTime(status.last_sync_at) : "-"
  nodes.webdavLastSuccess.textContent = status?.last_success_at ? formatTime(status.last_success_at) : "-"
  nodes.webdavSyncResult.textContent = status
    ? status.last_error || status.last_message || "No sync run yet"
    : "Admin token required"
  nodes.webdavTransferSummary.textContent = status
    ? `push ${status.pushed_items || 0}/${status.pushed_files || 0} files · pull ${status.pulled_items || 0}/${status.pulled_files || 0} files`
    : "-"
}

function renderCategoryControls() {
  populateCategorySelect(nodes.historyCategoryFilter, state.historyCategoryFilter, [])
  const allOption = document.createElement("option")
  allOption.value = ""
  allOption.textContent = "All categories"
  nodes.historyCategoryFilter.prepend(allOption)
  nodes.historyCategoryFilter.value = state.historyCategoryFilter
  renderQuickLinkSuggestion()
  renderQuickFileSuggestion(nodes.quickFileInput.files[0])
}

function populateCategorySelect(selectNode, selectedValue, extraOptions) {
  const values = new Set()
  const categories = []

  ;(state.categories || []).forEach((category) => {
    const name = normalizeCategoryName(category.name)
    if (!name || values.has(name)) {
      return
    }
    values.add(name)
    categories.push(name)
  })

  ;(extraOptions || []).forEach((value) => {
    const name = normalizeCategoryName(value)
    if (!name || values.has(name)) {
      return
    }
    values.add(name)
    categories.push(name)
  })

  categories.sort()
  selectNode.innerHTML = ""

  if (!categories.length) {
    const option = document.createElement("option")
    option.value = ""
    option.textContent = "No categories"
    selectNode.appendChild(option)
    return
  }

  categories.forEach((category) => {
    const option = document.createElement("option")
    option.value = category
    option.textContent = category
    selectNode.appendChild(option)
  })

  if (selectedValue) {
    selectNode.value = normalizeCategoryName(selectedValue)
  }
}

function renderTagList(container, tags) {
  container.innerHTML = ""
  ;(tags || []).forEach((tag) => {
    container.appendChild(renderTag(tag, "tag-predicted"))
  })
}

function renderTag(text, extraClass) {
  const tag = document.createElement("span")
  tag.className = `tag${extraClass ? ` ${extraClass}` : ""}`
  tag.textContent = text
  return tag
}

function renderActionButton(action) {
  const button = document.createElement("button")
  button.type = "button"
  button.textContent = action.label
  if (action.secondary) {
    button.classList.add("secondary")
  }
  if (action.danger) {
    button.classList.add("danger")
  }
  button.addEventListener("click", action.onClick)
  return button
}

function renderEmpty(message) {
  const li = document.createElement("li")
  li.className = "notice-item muted"
  li.textContent = message
  return li
}

function getFilteredHistoryItems() {
  return state.history.filter((item) => {
    if (state.historyTypeFilter && item.type !== state.historyTypeFilter) {
      return false
    }
    if (state.historyCategoryFilter && normalizeCategoryName(item.category) !== state.historyCategoryFilter) {
      return false
    }
    if (state.historyFavoritesOnly && !item.is_favorite) {
      return false
    }
    if (!state.historySearch) {
      return true
    }

    const haystack = [
      item.text,
      item.url,
      item.filename,
      item.source_device_name,
      item.category,
      item.type,
    ]
      .filter(Boolean)
      .join(" ")
      .toLowerCase()

    return haystack.includes(state.historySearch)
  })
}

function goToHistoryItem(id) {
  state.historyFocusID = id
  setPage("history")
  renderHistory()
}

function itemPrimaryTitle(item) {
  if (item.type === "text") {
    return summarizeText(item.text || "Text item")
  }
  if (item.type === "link") {
    return item.url || item.text || "Link item"
  }
  if (item.type === "image") {
    return item.filename || "Image item"
  }
  if (item.type === "file") {
    return item.filename || "File item"
  }
  return "Clipboard item"
}

function summarizeText(text) {
  const normalized = String(text || "").replace(/\s+/g, " ").trim()
  if (!normalized) {
    return "Text item"
  }
  if (normalized.length <= 84) {
    return normalized
  }
  return `${normalized.slice(0, 81)}...`
}

function itemTypeLabel(type) {
  switch (type) {
    case "text":
      return "Text"
    case "link":
      return "Link"
    case "image":
      return "Image"
    case "file":
      return "File"
    default:
      return type || "Item"
  }
}

function normalizeCategoryName(value) {
  return String(value || "").trim().toLowerCase()
}

function shouldSilenceAdminError(error) {
  const message = String(error?.message || "").toLowerCase()
  return message === "unauthorized" || message.includes("save an admin token or device token first")
}

function addNotice(message, tone) {
  const item = document.createElement("li")
  item.className = "notice-item"
  item.dataset.tone = tone
  item.textContent = `[${new Date().toLocaleTimeString()}] ${message}`
  nodes.noticeList.prepend(item)
}

function clearNotices() {
  nodes.noticeList.innerHTML = ""
}

async function apiFetch(path, options = {}) {
  const headers = new Headers(options.headers || {})
  const shouldAuth = options.auth !== false

  if (shouldAuth) {
    if (!state.token) {
      throw new Error("Save an admin token or device token first.")
    }
    headers.set("Authorization", `Bearer ${state.token}`)
  }

  const isFormData = typeof FormData !== "undefined" && options.body instanceof FormData
  if (options.body && !headers.has("Content-Type") && !isFormData) {
    headers.set("Content-Type", "application/json")
  }

  const response = await fetch(path, {
    method: options.method || "GET",
    headers,
    body: options.body,
  })

  const contentType = response.headers.get("Content-Type") || ""
  if (!contentType.includes("application/json")) {
    if (!response.ok) {
      throw new Error(`Request failed with status ${response.status}.`)
    }
    return {}
  }

  const payload = await response.json()
  if (!response.ok) {
    throw new Error(payload.error?.message || `Request failed with status ${response.status}.`)
  }
  return payload
}

async function loadProtectedAssetBlob(path) {
  if (!state.token) {
    throw new Error("Save an admin token or device token first.")
  }

  const response = await fetch(path, {
    headers: {
      Authorization: `Bearer ${state.token}`,
    },
  })

  if (!response.ok) {
    throw new Error(`Asset request failed with status ${response.status}.`)
  }

  return response.blob()
}

async function loadProtectedAssetURL(path) {
  const blob = await loadProtectedAssetBlob(path)
  return URL.createObjectURL(blob)
}

async function downloadProtectedFile(item) {
  if (!item?.download_url) {
    return
  }
  try {
    const objectURL = await loadProtectedAssetURL(item.download_url)
    const anchor = document.createElement("a")
    anchor.href = objectURL
    anchor.download = item.filename || "download"
    anchor.click()
    setTimeout(() => URL.revokeObjectURL(objectURL), 3000)
  } catch (error) {
    addNotice(error.message, "error")
  }
}

async function openProtectedItem(item) {
  if (!item?.preview_url && !item?.download_url) {
    return
  }
  try {
    const objectURL = await loadProtectedAssetURL(item.preview_url || item.download_url)
    window.open(objectURL, "_blank", "noopener")
    setTimeout(() => URL.revokeObjectURL(objectURL), 3000)
  } catch (error) {
    addNotice(error.message, "error")
  }
}

function openLinkItem(item) {
  const href = item?.url || item?.text
  if (!href) {
    return
  }
  window.open(href, "_blank", "noopener")
}

function formatBytes(value) {
  const bytes = Number(value || 0)
  if (!bytes) {
    return "0 B"
  }
  const units = ["B", "KB", "MB", "GB", "TB"]
  let size = bytes
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex += 1
  }
  return `${size >= 10 || unitIndex === 0 ? size.toFixed(0) : size.toFixed(1)} ${units[unitIndex]}`
}

function formatTime(value) {
  if (!value) {
    return "-"
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString()
}

function classifyUrl(rawURL) {
  if (!rawURL) {
    return { category: "link", tags: ["link"] }
  }

  try {
    const parsedURL = new URL(rawURL)
    const host = parsedURL.hostname.toLowerCase()

    if (matchesHost(host, ["github.com", "gitlab.com", "bitbucket.org"])) {
      return { category: "code", tags: ["code", "repo"] }
    }
    if (matchesHost(host, ["youtube.com", "youtu.be", "bilibili.com", "vimeo.com"])) {
      return { category: "video", tags: ["video"] }
    }
    if (matchesHost(host, ["arxiv.org", "doi.org", "scholar.google.com"])) {
      return { category: "paper", tags: ["research", "paper"] }
    }
    if (matchesHost(host, ["docs.google.com", "notion.so", "www.notion.so", "figma.com"])) {
      return { category: "document", tags: ["document"] }
    }
    if (matchesHost(host, ["x.com", "twitter.com", "reddit.com", "weibo.com"])) {
      return { category: "social", tags: ["social"] }
    }
    if (
      host.startsWith("amazon.") ||
      matchesHost(host, ["taobao.com", "jd.com", "ebay.com", "www.ebay.com"])
    ) {
      return { category: "shopping", tags: ["shopping"] }
    }
    if (
      matchesHost(host, [
        "nytimes.com",
        "bbc.com",
        "bbc.co.uk",
        "cnn.com",
        "theguardian.com",
        "reuters.com",
        "apnews.com",
      ])
    ) {
      return { category: "news", tags: ["news"] }
    }
  } catch (_error) {
    return { category: "link", tags: ["link"] }
  }

  return { category: "link", tags: ["link"] }
}

function classifyFile(file) {
  if (!file) {
    return { category: "file", tags: ["file"] }
  }

  const mimeType = (file.type || "").toLowerCase()
  const filename = (file.name || "").toLowerCase()
  let category = "file"
  let tags = ["file"]

  if (["image/png", "image/jpeg", "image/webp", "image/gif"].includes(mimeType)) {
    category = "image"
    tags = ["image"]
  } else if (mimeType === "image/svg+xml") {
    category = "image"
    tags = ["image", "svg"]
  } else if (mimeType === "application/pdf") {
    category = "document"
    tags = ["pdf", "document"]
  } else if (["text/plain", "text/markdown"].includes(mimeType)) {
    category = "document"
    tags = ["text", "document"]
  } else if (["application/zip", "application/x-7z-compressed", "application/x-rar-compressed"].includes(mimeType)) {
    category = "archive"
    tags = ["archive"]
  } else if (mimeType.startsWith("video/")) {
    category = "video"
    tags = ["video"]
  } else if (mimeType.startsWith("audio/")) {
    category = "audio"
    tags = ["audio"]
  } else if (["application/json", "application/xml", "text/csv"].includes(mimeType)) {
    category = "data"
    tags = ["data"]
  }

  if (filename.includes("screenshot") || filename.includes("截屏") || filename.includes("screen") || filename.includes("capture")) {
    tags.push("screenshot")
  }
  if (filename.includes("avatar") || filename.includes("profile")) {
    tags.push("avatar")
  }
  if (filename.includes("logo")) {
    tags.push("logo")
  }
  if (filename.includes("diagram") || filename.includes("chart") || filename.includes("plot") || filename.includes("graph")) {
    tags.push("diagram")
  }

  return { category, tags: Array.from(new Set(tags)) }
}

function matchesHost(host, domains) {
  return domains.some((domain) => host === domain || host.endsWith(`.${domain}`))
}
