const storageKey = "clipbridge.webui.token";

const state = {
  token: localStorage.getItem(storageKey) || "",
  latest: null,
  favorites: [],
  history: [],
  categories: [],
  historyFilter: "",
  appSettings: null,
  cleanupStatus: null,
  storageStatus: null,
  webdavStatus: null,
};

const nodes = {
  tokenInput: document.getElementById("token-input"),
  tokenStatus: document.getElementById("token-status"),
  noticeList: document.getElementById("notice-list"),
  healthState: document.getElementById("health-state"),
  healthVersion: document.getElementById("health-version"),
  serverOrigin: document.getElementById("server-origin"),
  quickInput: document.getElementById("quick-clipboard-input"),
  quickLinkInput: document.getElementById("quick-link-input"),
  quickFileInput: document.getElementById("quick-file-input"),
  latestEmpty: document.getElementById("latest-empty"),
  latestCard: document.getElementById("latest-card"),
  latestId: document.getElementById("latest-id"),
  latestCreated: document.getElementById("latest-created"),
  latestType: document.getElementById("latest-type"),
  latestCategory: document.getElementById("latest-category"),
  latestFavorite: document.getElementById("latest-favorite"),
  latestContent: document.getElementById("latest-content"),
  latestSourceName: document.getElementById("latest-source-name"),
  latestCategorySelect: document.getElementById("latest-category-select"),
  historyCategoryFilter: document.getElementById("history-category-filter"),
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
  cleanupTTLHours: document.getElementById("cleanup-ttl-hours"),
  cleanupMaxItems: document.getElementById("cleanup-max-items"),
  cleanupMaxSizeMB: document.getElementById("cleanup-max-size-mb"),
  cleanupIntervalMinutes: document.getElementById("cleanup-interval-minutes"),
  cleanupEnabled: document.getElementById("cleanup-enabled"),
  storageHistoryCount: document.getElementById("storage-history-count"),
  storageFavoriteCount: document.getElementById("storage-favorite-count"),
  storageTotalBytes: document.getElementById("storage-total-bytes"),
  storageFileBytes: document.getElementById("storage-file-bytes"),
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
};

nodes.tokenInput.value = state.token;
nodes.serverOrigin.textContent = window.location.origin;
updateTokenStatus();

document.getElementById("save-token-button").addEventListener("click", saveToken);
document.getElementById("clear-token-button").addEventListener("click", clearToken);
document.getElementById("refresh-health-button").addEventListener("click", loadHealth);
document.getElementById("upload-quick-button").addEventListener("click", uploadQuickClipboard);
document.getElementById("upload-link-button").addEventListener("click", uploadLinkClipboard);
document.getElementById("upload-file-button").addEventListener("click", uploadFileClipboard);
document.getElementById("copy-latest-button").addEventListener("click", copyLatestText);
document.getElementById("refresh-latest-button").addEventListener("click", loadLatest);
document.getElementById("toggle-latest-favorite-button").addEventListener("click", toggleLatestFavorite);
document.getElementById("delete-latest-button").addEventListener("click", deleteLatest);
document.getElementById("save-latest-category-button").addEventListener("click", updateLatestCategory);
document.getElementById("refresh-history-button").addEventListener("click", loadHistory);
document.getElementById("refresh-favorites-button").addEventListener("click", loadFavorites);
document.getElementById("refresh-devices-button").addEventListener("click", loadDevices);
document.getElementById("generate-pairing-button").addEventListener("click", generatePairingCode);
document.getElementById("create-category-button").addEventListener("click", createCategory);
document.getElementById("save-admin-token-button").addEventListener("click", saveAdminToken);
document.getElementById("save-limits-button").addEventListener("click", saveLimits);
document.getElementById("refresh-cleanup-button").addEventListener("click", loadCleanupData);
document.getElementById("run-cleanup-button").addEventListener("click", runCleanupNow);
document.getElementById("save-cleanup-settings-button").addEventListener("click", saveCleanupSettings);
document.getElementById("save-webdav-settings-button").addEventListener("click", saveWebDAVSettings);
document.getElementById("refresh-webdav-button").addEventListener("click", loadWebDAVData);
document.getElementById("test-webdav-button").addEventListener("click", testWebDAVConnection);
document.getElementById("sync-webdav-button").addEventListener("click", runWebDAVSync);
document.getElementById("clear-notices-button").addEventListener("click", clearNotices);
nodes.historyCategoryFilter.addEventListener("change", async (event) => {
  state.historyFilter = event.target.value;
  await loadHistory();
});

loadHealth();
loadAllProtectedData();

function saveToken() {
  state.token = nodes.tokenInput.value.trim();
  persistToken();
  updateTokenStatus();
  loadAllProtectedData();
}

function clearToken() {
  state.token = "";
  nodes.tokenInput.value = "";
  persistToken();
  updateTokenStatus();
  loadAllProtectedData();
}

function persistToken() {
  if (state.token) {
    localStorage.setItem(storageKey, state.token);
    addNotice("Saved token locally for this browser.", "success");
  } else {
    localStorage.removeItem(storageKey);
    addNotice("Cleared the saved token.", "info");
  }
}

function updateTokenStatus() {
  nodes.tokenStatus.textContent = state.token
    ? "Saved token will be used for API requests."
    : "No token saved yet.";
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
  ]);
}

async function loadHealth() {
  try {
    const response = await apiFetch("/api/health", { auth: false });
    nodes.healthState.textContent = response.data.ok ? "Healthy" : "Unavailable";
    nodes.healthVersion.textContent = response.data.version;
  } catch (error) {
    nodes.healthState.textContent = "Unavailable";
    nodes.healthVersion.textContent = "-";
    addNotice(error.message, "error");
  }
}

async function loadAppSettings() {
  if (!state.token) {
    state.appSettings = null;
    renderAppSettings();
    renderLimits();
    return;
  }

  try {
    const response = await apiFetch("/api/settings");
    state.appSettings = response.data;
    renderAppSettings();
    renderLimits();
    renderWebDAVSettings();
  } catch (error) {
    state.appSettings = null;
    renderAppSettings();
    renderLimits();
    renderWebDAVSettings();
    addNotice(error.message, "error");
  }
}

async function uploadQuickClipboard() {
  const content = nodes.quickInput.value.trim();
  if (!content) {
    addNotice("Enter some text before uploading.", "info");
    return;
  }

  try {
    await apiFetch("/api/clipboard/text", {
      method: "POST",
      body: JSON.stringify({
        content,
        source_device_id: "web-ui",
        source_device_name: "Web UI",
      }),
    });
    nodes.quickInput.value = "";
    addNotice("Uploaded browser text to the clipboard stack.", "success");
    await refreshClipboardViews();
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function uploadLinkClipboard() {
  const url = nodes.quickLinkInput.value.trim();
  if (!url) {
    addNotice("Enter a link before uploading.", "info");
    return;
  }

  try {
    await apiFetch("/api/clipboard/link", {
      method: "POST",
      body: JSON.stringify({
        url,
        source_device_id: "web-ui",
        source_device_name: "Web UI",
      }),
    });
    nodes.quickLinkInput.value = "";
    addNotice("Uploaded link clipboard item.", "success");
    await refreshClipboardViews();
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function uploadFileClipboard() {
  const file = nodes.quickFileInput.files[0];
  if (!file) {
    addNotice("Choose one file first.", "info");
    return;
  }

  const formData = new FormData();
  formData.append("file", file);
  formData.append("source_device_id", "web-ui");
  formData.append("source_device_name", "Web UI");

  try {
    await apiFetch("/api/clipboard/file", { method: "POST", body: formData });
    nodes.quickFileInput.value = "";
    addNotice("Uploaded image or file clipboard item.", "success");
    await refreshClipboardViews();
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function loadCategories() {
  if (!state.token) {
    state.categories = [];
    renderCategoryControls();
    return;
  }
  try {
    const response = await apiFetch("/api/categories");
    state.categories = response.data.items || [];
    renderCategoryControls();
  } catch (error) {
    state.categories = [];
    renderCategoryControls();
    addNotice(error.message, "error");
  }
}

async function loadLatest() {
  if (!state.token) {
    state.latest = null;
    renderLatest();
    return;
  }
  try {
    const response = await apiFetch("/api/clipboard/latest");
    state.latest = response.data;
    renderLatest();
  } catch (error) {
    state.latest = null;
    renderLatest();
    if (!String(error.message).includes("no clipboard item found")) {
      addNotice(error.message, "error");
    }
  }
}

async function loadHistory() {
  if (!state.token) {
    state.history = [];
    renderHistory();
    return;
  }
  const suffix = state.historyFilter ? `?category=${encodeURIComponent(state.historyFilter)}` : "";
  try {
    const response = await apiFetch(`/api/clipboard/history${suffix}`);
    state.history = response.data.items || [];
    renderHistory();
  } catch (error) {
    state.history = [];
    renderHistory();
    addNotice(error.message, "error");
  }
}

async function loadFavorites() {
  if (!state.token) {
    state.favorites = [];
    renderFavorites();
    return;
  }
  try {
    const response = await apiFetch("/api/favorites");
    state.favorites = response.data.items || [];
    renderFavorites();
  } catch (error) {
    state.favorites = [];
    renderFavorites();
    addNotice(error.message, "error");
  }
}

async function loadDevices() {
  if (!state.token) {
    renderDevices([]);
    return;
  }
  try {
    const response = await apiFetch("/api/auth/devices");
    renderDevices(response.data.items || []);
  } catch (error) {
    renderDevices([]);
    addNotice(error.message, "error");
  }
}

async function loadCleanupData() {
  if (!state.token) {
    state.cleanupStatus = null;
    state.storageStatus = null;
    renderCleanupStatus();
    renderCleanupSettings();
    return;
  }

  try {
    const [statusResponse, storageResponse] = await Promise.all([
      apiFetch("/api/admin/cleanup/status"),
      apiFetch("/api/admin/storage/status"),
    ]);
    state.cleanupStatus = statusResponse.data;
    state.storageStatus = storageResponse.data;
    renderCleanupStatus();
    renderCleanupSettings();
  } catch (error) {
    state.cleanupStatus = null;
    state.storageStatus = null;
    renderCleanupStatus();
    renderCleanupSettings();
    addNotice(error.message, "error");
  }
}

async function loadWebDAVData() {
  if (!state.token) {
    state.webdavStatus = null;
    renderWebDAVSettings();
    renderWebDAVStatus();
    return;
  }

  try {
    const response = await apiFetch("/api/admin/webdav/status");
    state.webdavStatus = response.data;
    renderWebDAVSettings();
    renderWebDAVStatus();
  } catch (error) {
    state.webdavStatus = null;
    renderWebDAVSettings();
    renderWebDAVStatus();
    addNotice(error.message, "error");
  }
}

async function saveAdminToken() {
  const token = nodes.adminTokenInput.value.trim();
  if (!token) {
    addNotice("Enter a new admin token first.", "info");
    return;
  }

  try {
    await apiFetch("/api/settings", {
      method: "PATCH",
      body: JSON.stringify({ admin_token: token }),
    });
    state.token = token;
    nodes.tokenInput.value = token;
    persistToken();
    updateTokenStatus();
    await loadAppSettings();
    addNotice("Updated admin token and switched this browser to the new token.", "success");
  } catch (error) {
    addNotice(error.message, "error");
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
  };

  try {
    await apiFetch("/api/settings/limits", {
      method: "PATCH",
      body: JSON.stringify(payload),
    });
    await loadAppSettings();
    addNotice("Saved runtime limits.", "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function saveCleanupSettings() {
  const payload = {
    ttl_hours: Number(nodes.cleanupTTLHours.value),
    max_items: Number(nodes.cleanupMaxItems.value),
    max_total_size_mb: Number(nodes.cleanupMaxSizeMB.value),
    interval_minutes: Number(nodes.cleanupIntervalMinutes.value),
    enabled: nodes.cleanupEnabled.checked,
  };

  try {
    const response = await apiFetch("/api/settings/cleanup", {
      method: "PATCH",
      body: JSON.stringify(payload),
    });
    if (state.appSettings) {
      state.appSettings.cleanup = response.data;
    }
    renderCleanupSettings();
    addNotice("Saved cleanup policy.", "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function runCleanupNow() {
  try {
    await apiFetch("/api/admin/cleanup/run", { method: "POST" });
    await Promise.all([loadCleanupData(), loadLatest(), loadHistory(), loadFavorites()]);
    addNotice("Triggered one manual cleanup run.", "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function saveWebDAVSettings() {
  const payload = {
    enabled: nodes.webdavEnabled.checked,
    url: nodes.webdavURL.value.trim(),
    username: nodes.webdavUsername.value.trim(),
    password: nodes.webdavPassword.value,
    base_path: nodes.webdavBasePath.value.trim(),
  };

  try {
    const response = await apiFetch("/api/settings/webdav", {
      method: "PATCH",
      body: JSON.stringify(payload),
    });
    if (state.appSettings) {
      state.appSettings.webdav = response.data;
    }
    renderWebDAVSettings();
    addNotice("Saved WebDAV settings.", "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function testWebDAVConnection() {
  try {
    const response = await apiFetch("/api/admin/webdav/test", { method: "POST" });
    state.webdavStatus = response.data;
    renderWebDAVStatus();
    addNotice("WebDAV connection test succeeded.", "success");
  } catch (error) {
    await loadWebDAVData();
    addNotice(error.message, "error");
  }
}

async function runWebDAVSync() {
  try {
    const response = await apiFetch("/api/admin/webdav/sync", { method: "POST" });
    state.webdavStatus = response.data;
    renderWebDAVStatus();
    await refreshClipboardViews();
    addNotice("WebDAV sync completed.", "success");
  } catch (error) {
    await loadWebDAVData();
    addNotice(error.message, "error");
  }
}

async function generatePairingCode() {
  try {
    const response = await apiFetch("/api/auth/pairing-codes", { method: "POST" });
    nodes.pairingOutput.classList.remove("hidden");
    nodes.pairingCodeValue.textContent = response.data.pairing_code;
    nodes.pairingCodeExpiry.textContent = response.data.expires_at;
    nodes.pairingCodeURI.textContent = response.data.pairing_uri;
    addNotice("Generated a new pairing code.", "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function copyLatestText() {
  if (!state.latest) {
    addNotice("No clipboard item available to copy.", "info");
    return;
  }
  if (state.latest.type !== "text" && state.latest.type !== "link") {
    addNotice("Only text or link items can be copied into the browser clipboard.", "info");
    return;
  }
  await copyText(state.latest.text || state.latest.url || "");
}

async function copyText(value) {
  try {
    await navigator.clipboard.writeText(value || "");
    addNotice("Copied content into the browser clipboard.", "success");
  } catch (error) {
    addNotice("Copy failed. Your browser may block clipboard access.", "error");
  }
}

async function deleteLatest() {
  if (!state.latest) {
    return;
  }
  await deleteItem(state.latest.id);
}

async function deleteItem(id) {
  try {
    await apiFetch(`/api/clipboard/items/${id}`, { method: "DELETE" });
    await refreshClipboardViews();
    addNotice(`Deleted clipboard item #${id}.`, "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function toggleLatestFavorite() {
  if (!state.latest) {
    return;
  }
  await setFavorite(state.latest.id, !state.latest.is_favorite);
}

async function setFavorite(id, favorite) {
  try {
    await apiFetch(`/api/clipboard/items/${id}/favorite`, {
      method: favorite ? "POST" : "DELETE",
    });
    await refreshClipboardViews();
    addNotice(favorite ? `Favorited clipboard item #${id}.` : `Removed favorite from item #${id}.`, "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function updateLatestCategory() {
  if (!state.latest) {
    return;
  }
  await updateItemCategory(state.latest.id, nodes.latestCategorySelect.value);
}

async function updateItemCategory(id, category) {
  try {
    await apiFetch(`/api/clipboard/items/${id}/category`, {
      method: "PATCH",
      body: JSON.stringify({ category }),
    });
    await refreshClipboardViews();
    addNotice(`Updated category for clipboard item #${id}.`, "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function createCategory() {
  const name = nodes.newCategoryInput.value.trim();
  if (!name) {
    addNotice("Enter a category name first.", "info");
    return;
  }

  try {
    await apiFetch("/api/categories", {
      method: "POST",
      body: JSON.stringify({ name }),
    });
    nodes.newCategoryInput.value = "";
    await loadCategories();
    addNotice(`Created category "${name}".`, "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function revokeDevice(id) {
  try {
    await apiFetch(`/api/auth/devices/${id}`, { method: "DELETE" });
    await loadDevices();
    addNotice(`Revoked device #${id}.`, "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function refreshClipboardViews() {
  await Promise.all([loadLatest(), loadHistory(), loadFavorites(), loadCleanupData()]);
}

function renderLatest() {
  if (!state.latest) {
    nodes.latestEmpty.classList.remove("hidden");
    nodes.latestCard.classList.add("hidden");
    nodes.latestEmpty.textContent = state.token
      ? "No clipboard item loaded yet."
      : "Save an admin token or device token to load clipboard data.";
    return;
  }

  nodes.latestEmpty.classList.add("hidden");
  nodes.latestCard.classList.remove("hidden");
  nodes.latestId.textContent = `#${state.latest.id}`;
  nodes.latestCreated.textContent = state.latest.created_at;
  nodes.latestType.textContent = state.latest.type || "text";
  nodes.latestCategory.textContent = state.latest.category || "uncategorized";
  nodes.latestFavorite.textContent = state.latest.is_favorite ? "Favorite" : "Normal";
  nodes.latestSourceName.textContent = state.latest.source_device_name || "Unknown Source";
  document.getElementById("toggle-latest-favorite-button").textContent = state.latest.is_favorite
    ? "Unfavorite"
    : "Favorite";
  populateCategorySelect(nodes.latestCategorySelect, state.latest.category || "");

  nodes.latestContent.innerHTML = "";
  nodes.latestContent.appendChild(renderItemContent(state.latest));
}

function renderHistory() {
  renderClipboardCollection(
    nodes.historyList,
    state.history,
    state.token ? "No history available with the current token." : "Save a token to view clipboard history.",
  );
}

function renderFavorites() {
  renderClipboardCollection(
    nodes.favoritesList,
    state.favorites,
    state.token ? "No favorites saved yet." : "Save a token to view favorites.",
  );
}

function renderClipboardCollection(container, items, emptyMessage) {
  container.innerHTML = "";
  if (!items.length) {
    container.appendChild(renderEmpty(emptyMessage));
    return;
  }
  items.forEach((item) => container.appendChild(renderClipboardItem(item)));
}

function renderClipboardItem(item) {
  const li = document.createElement("li");
  li.className = "history-item";

  const title = document.createElement("strong");
  title.textContent = `#${item.id} · ${item.created_at}`;

  const tags = document.createElement("div");
  tags.className = "clipboard-tags";
  tags.innerHTML = `
    <span class="tag">${escapeHTML(item.type || "text")}</span>
    <span class="tag tag-accent">${escapeHTML(item.category || "uncategorized")}</span>
    <span class="tag tag-favorite">${item.is_favorite ? "Favorite" : "Normal"}</span>
  `;

  const content = renderItemContent(item);
  const source = document.createElement("div");
  source.className = "muted";
  source.textContent = `From: ${item.source_device_name || "Unknown Source"}`;

  const categoryRow = document.createElement("div");
  categoryRow.className = "category-row";
  const categoryLabel = document.createElement("label");
  categoryLabel.textContent = "Category";
  const categorySelect = document.createElement("select");
  populateCategorySelect(categorySelect, item.category || "");
  const categoryButton = document.createElement("button");
  categoryButton.type = "button";
  categoryButton.className = "secondary";
  categoryButton.textContent = "Update Category";
  categoryButton.addEventListener("click", () => updateItemCategory(item.id, categorySelect.value));
  categoryRow.append(categoryLabel, categorySelect, categoryButton);

  const actions = document.createElement("div");
  actions.className = "panel-actions";

  if (item.type === "text" || item.type === "link") {
    const copyButton = document.createElement("button");
    copyButton.type = "button";
    copyButton.textContent = "Copy";
    copyButton.addEventListener("click", () => copyText(item.text || item.url || ""));
    actions.appendChild(copyButton);
  }

  const favoriteButton = document.createElement("button");
  favoriteButton.type = "button";
  favoriteButton.className = "secondary";
  favoriteButton.textContent = item.is_favorite ? "Unfavorite" : "Favorite";
  favoriteButton.addEventListener("click", () => setFavorite(item.id, !item.is_favorite));

  const deleteButton = document.createElement("button");
  deleteButton.type = "button";
  deleteButton.className = "danger";
  deleteButton.textContent = "Delete";
  deleteButton.addEventListener("click", () => deleteItem(item.id));

  actions.append(favoriteButton, deleteButton);
  li.append(title, tags, content, source, categoryRow, actions);
  return li;
}

function renderItemContent(item) {
  if (item.type === "image" && item.preview_url) {
    const wrapper = document.createElement("div");
    wrapper.className = "item-content";
    const preview = document.createElement("img");
    preview.className = "image-preview";
    preview.alt = item.filename || `Clipboard image ${item.id}`;
    preview.dataset.loading = "true";
    loadProtectedAssetURL(item.preview_url)
      .then((objectURL) => {
        preview.src = objectURL;
      })
      .catch(() => {
        preview.alt = "Image preview requires a valid token.";
      });
    const meta = document.createElement("div");
    meta.className = "muted";
    meta.textContent = item.filename || "image";
    const downloadButton = document.createElement("button");
    downloadButton.type = "button";
    downloadButton.className = "secondary";
    downloadButton.textContent = "Download Image";
    downloadButton.addEventListener("click", () => downloadProtectedFile(item));
    wrapper.append(preview, meta, downloadButton);
    return wrapper;
  }

  if (item.type === "file") {
    const wrapper = document.createElement("div");
    wrapper.className = "item-content";
    const title = document.createElement("div");
    title.innerHTML = `<strong>${escapeHTML(item.filename || "file")}</strong>`;
    const type = document.createElement("div");
    type.className = "muted";
    type.textContent = item.mime_type || "application/octet-stream";
    const downloadButton = document.createElement("button");
    downloadButton.type = "button";
    downloadButton.className = "secondary";
    downloadButton.textContent = "Download File";
    downloadButton.addEventListener("click", () => downloadProtectedFile(item));
    wrapper.append(title, type, downloadButton);
    return wrapper;
  }

  if (item.type === "link") {
    const wrapper = document.createElement("div");
    wrapper.className = "item-content";
    const anchor = document.createElement("a");
    anchor.href = item.url || item.text || "#";
    anchor.target = "_blank";
    anchor.rel = "noreferrer";
    anchor.textContent = item.url || item.text || "";
    wrapper.appendChild(anchor);
    return wrapper;
  }

  const pre = document.createElement("pre");
  pre.textContent = item.text || "";
  return pre;
}

function renderDevices(items) {
  nodes.devicesList.innerHTML = "";
  if (!items.length) {
    nodes.devicesList.appendChild(
      renderEmpty(state.token ? "No devices visible with the current token." : "Save the admin token to manage devices."),
    );
    return;
  }
  items.forEach((device) => {
    const li = document.createElement("li");
    li.className = "device-item";
    li.innerHTML = `
      <strong>#${device.id} · ${escapeHTML(device.name)}</strong>
      <div class="muted">Created: ${escapeHTML(device.created_at)}</div>
      <div class="muted">Last Seen: ${escapeHTML(device.last_seen_at || "never")}</div>
      <div class="muted">Revoked: ${escapeHTML(device.revoked_at || "active")}</div>
    `;
    if (!device.revoked_at) {
      const revokeButton = document.createElement("button");
      revokeButton.type = "button";
      revokeButton.className = "danger";
      revokeButton.textContent = "Revoke";
      revokeButton.addEventListener("click", () => revokeDevice(device.id));
      li.appendChild(revokeButton);
    }
    nodes.devicesList.appendChild(li);
  });
}

function renderAppSettings() {
  const appSettings = state.appSettings;
  nodes.adminTokenInput.value = appSettings?.admin_token || "";
  if (!appSettings) {
    nodes.startupSettings.textContent = "Admin token required.";
    return;
  }
  const startup = appSettings.startup || {};
  const restartFields = (appSettings.restart_required_fields || []).join(", ");
  nodes.startupSettings.textContent =
    `Host: ${startup.host} · Port: ${startup.port} · Data Dir: ${startup.data_dir} · DB: ${startup.database_path} · Restart required for: ${restartFields}`;
}

function renderLimits() {
  const limits = state.appSettings?.limits;
  nodes.limitMinText.value = limits?.min_text_bytes ?? "";
  nodes.limitMaxText.value = limits?.max_text_bytes ?? "";
  nodes.limitMinImage.value = limits?.min_image_bytes ?? "";
  nodes.limitMaxImage.value = limits?.max_image_bytes ?? "";
  nodes.limitMinFile.value = limits?.min_file_bytes ?? "";
  nodes.limitMaxFile.value = limits?.max_file_bytes ?? "";
  nodes.limitMinLink.value = limits?.min_link_bytes ?? "";
  nodes.limitMaxLink.value = limits?.max_link_bytes ?? "";
  nodes.limitMaxRequest.value = limits?.max_request_bytes ?? "";
}

function renderCleanupStatus() {
  const status = state.cleanupStatus;
  const storage = state.storageStatus;
  nodes.storageHistoryCount.textContent = storage ? String(storage.history_count) : "-";
  nodes.storageFavoriteCount.textContent = storage ? String(storage.favorite_count) : "-";
  nodes.storageTotalBytes.textContent = storage ? formatBytes(storage.total_bytes) : "-";
  nodes.storageFileBytes.textContent = storage ? formatBytes(storage.file_bytes) : "-";
  nodes.cleanupLastRun.textContent = status?.last_run_at || "-";
  nodes.cleanupLastResult.textContent = status
    ? status.last_error || `expired ${status.deleted_expired} · count ${status.deleted_max_items} · storage ${status.deleted_storage}`
    : "Admin token required";
}

function renderCleanupSettings() {
  const cleanup = state.appSettings?.cleanup;
  nodes.cleanupTTLHours.value = cleanup?.ttl_hours ?? "";
  nodes.cleanupMaxItems.value = cleanup?.max_items ?? "";
  nodes.cleanupMaxSizeMB.value = cleanup?.max_total_size_mb ?? "";
  nodes.cleanupIntervalMinutes.value = cleanup?.interval_minutes ?? "";
  nodes.cleanupEnabled.checked = Boolean(cleanup?.enabled);
}

function renderWebDAVSettings() {
  const webdav = state.appSettings?.webdav;
  nodes.webdavURL.value = webdav?.url ?? "";
  nodes.webdavUsername.value = webdav?.username ?? "";
  nodes.webdavPassword.value = webdav?.password ?? "";
  nodes.webdavBasePath.value = webdav?.base_path ?? "ClipBridgeServer";
  nodes.webdavEnabled.checked = Boolean(webdav?.enabled);
}

function renderWebDAVStatus() {
  const status = state.webdavStatus;
  nodes.webdavTestedAt.textContent = status?.tested_at || "-";
  nodes.webdavTestResult.textContent = status
    ? status.last_test_success
      ? "OK"
      : status.last_test_error || "Failed"
    : "Admin token required";
  nodes.webdavLastSync.textContent = status?.last_sync_at || "-";
  nodes.webdavLastSuccess.textContent = status?.last_success_at || "-";
  nodes.webdavSyncResult.textContent = status
    ? status.last_error || status.last_message || "No sync run yet"
    : "Admin token required";
  nodes.webdavTransferSummary.textContent = status
    ? `push ${status.pushed_items || 0}/${status.pushed_files || 0} files · pull ${status.pulled_items || 0}/${status.pulled_files || 0} files`
    : "-";
}

function renderCategoryControls() {
  populateCategorySelect(nodes.latestCategorySelect, state.latest ? state.latest.category : "");
  populateFilterSelect();
}

function populateFilterSelect() {
  const currentValue = state.historyFilter || "";
  nodes.historyCategoryFilter.innerHTML = "";
  const allOption = document.createElement("option");
  allOption.value = "";
  allOption.textContent = "All Categories";
  nodes.historyCategoryFilter.appendChild(allOption);
  state.categories.forEach((category) => {
    const option = document.createElement("option");
    option.value = category.name;
    option.textContent = category.name;
    nodes.historyCategoryFilter.appendChild(option);
  });
  nodes.historyCategoryFilter.value = currentValue;
}

function populateCategorySelect(selectNode, selectedValue) {
  selectNode.innerHTML = "";
  if (!state.categories.length) {
    const option = document.createElement("option");
    option.value = "";
    option.textContent = "No categories";
    selectNode.appendChild(option);
    return;
  }
  state.categories.forEach((category) => {
    const option = document.createElement("option");
    option.value = category.name;
    option.textContent = category.name;
    selectNode.appendChild(option);
  });
  if (selectedValue) {
    selectNode.value = selectedValue;
  }
}

function renderEmpty(message) {
  const li = document.createElement("li");
  li.className = "notice-item muted";
  li.textContent = message;
  return li;
}

function addNotice(message, tone) {
  const item = document.createElement("li");
  item.className = "notice-item";
  item.dataset.tone = tone;
  item.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
  nodes.noticeList.prepend(item);
}

function clearNotices() {
  nodes.noticeList.innerHTML = "";
}

async function apiFetch(path, options = {}) {
  const headers = new Headers(options.headers || {});
  const shouldAuth = options.auth !== false;

  if (shouldAuth) {
    if (!state.token) {
      throw new Error("Save an admin token or device token first.");
    }
    headers.set("Authorization", `Bearer ${state.token}`);
  }

  const isFormData = typeof FormData !== "undefined" && options.body instanceof FormData;
  if (options.body && !headers.has("Content-Type") && !isFormData) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(path, {
    method: options.method || "GET",
    headers,
    body: options.body,
  });

  const contentType = response.headers.get("Content-Type") || "";
  if (!contentType.includes("application/json")) {
    if (!response.ok) {
      throw new Error(`Request failed with status ${response.status}.`);
    }
    return {};
  }

  const payload = await response.json();
  if (!response.ok) {
    throw new Error(payload.error?.message || `Request failed with status ${response.status}.`);
  }
  return payload;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

function formatBytes(value) {
  const bytes = Number(value || 0);
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
}

async function loadProtectedAssetURL(path) {
  const response = await fetch(path, {
    headers: {
      Authorization: `Bearer ${state.token}`,
    },
  });
  if (!response.ok) {
    throw new Error(`Asset request failed with status ${response.status}.`);
  }
  const blob = await response.blob();
  return URL.createObjectURL(blob);
}

async function downloadProtectedFile(item) {
  try {
    const objectURL = await loadProtectedAssetURL(item.download_url);
    const anchor = document.createElement("a");
    anchor.href = objectURL;
    anchor.download = item.filename || `clipbridge-${item.id}`;
    document.body.appendChild(anchor);
    anchor.click();
    anchor.remove();
    addNotice(`Downloaded clipboard item #${item.id}.`, "success");
  } catch (error) {
    addNotice(error.message, "error");
  }
}
