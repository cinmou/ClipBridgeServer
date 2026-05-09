const storageKey = "clipbridge.webui.token";

const state = {
  token: localStorage.getItem(storageKey) || "",
  latest: null,
  favorites: [],
  history: [],
  categories: [],
  historyFilter: "",
  cleanupSettings: null,
  cleanupStatus: null,
  storageStatus: null,
};

const nodes = {
  tokenInput: document.getElementById("token-input"),
  tokenStatus: document.getElementById("token-status"),
  noticeList: document.getElementById("notice-list"),
  healthState: document.getElementById("health-state"),
  healthVersion: document.getElementById("health-version"),
  serverOrigin: document.getElementById("server-origin"),
  quickInput: document.getElementById("quick-clipboard-input"),
  latestEmpty: document.getElementById("latest-empty"),
  latestCard: document.getElementById("latest-card"),
  latestId: document.getElementById("latest-id"),
  latestCreated: document.getElementById("latest-created"),
  latestType: document.getElementById("latest-type"),
  latestCategory: document.getElementById("latest-category"),
  latestFavorite: document.getElementById("latest-favorite"),
  latestText: document.getElementById("latest-text"),
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
};

nodes.tokenInput.value = state.token;
nodes.serverOrigin.textContent = window.location.origin;
updateTokenStatus();

document.getElementById("save-token-button").addEventListener("click", saveToken);
document.getElementById("clear-token-button").addEventListener("click", clearToken);
document.getElementById("refresh-health-button").addEventListener("click", loadHealth);
document.getElementById("upload-quick-button").addEventListener("click", uploadQuickClipboard);
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
document.getElementById("refresh-cleanup-button").addEventListener("click", loadCleanupData);
document.getElementById("run-cleanup-button").addEventListener("click", runCleanupNow);
document
  .getElementById("save-cleanup-settings-button")
  .addEventListener("click", saveCleanupSettings);
document.getElementById("clear-notices-button").addEventListener("click", clearNotices);
nodes.historyCategoryFilter.addEventListener("change", async (event) => {
  state.historyFilter = event.target.value;
  await loadHistory();
});

loadHealth();
loadAllProtectedData();

function saveToken() {
  state.token = nodes.tokenInput.value.trim();
  if (state.token) {
    localStorage.setItem(storageKey, state.token);
    addNotice("Saved token locally for this browser.", "success");
  } else {
    localStorage.removeItem(storageKey);
    addNotice("Cleared the saved token.", "info");
  }

  updateTokenStatus();
  loadAllProtectedData();
}

function clearToken() {
  nodes.tokenInput.value = "";
  saveToken();
}

function updateTokenStatus() {
  nodes.tokenStatus.textContent = state.token
    ? "Saved token will be used for API requests."
    : "No token saved yet.";
}

async function loadAllProtectedData() {
  await loadCategories();
  await Promise.all([loadLatest(), loadHistory(), loadFavorites(), loadDevices(), loadCleanupData()]);
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

async function uploadQuickClipboard() {
  const content = nodes.quickInput.value;
  if (!content.trim()) {
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
    addNotice("Uploaded browser text to the server clipboard stack.", "success");
    nodes.quickInput.value = "";
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
    if (!String(error.message).includes("no text clipboard item found")) {
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

  const suffix = state.historyFilter
    ? `?category=${encodeURIComponent(state.historyFilter)}`
    : "";

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
    state.cleanupSettings = null;
    state.cleanupStatus = null;
    state.storageStatus = null;
    renderCleanupStatus();
    renderCleanupSettings();
    return;
  }

  try {
    const [settingsResponse, statusResponse, storageResponse] = await Promise.all([
      apiFetch("/api/settings/cleanup"),
      apiFetch("/api/admin/cleanup/status"),
      apiFetch("/api/admin/storage/status"),
    ]);

    state.cleanupSettings = settingsResponse.data;
    state.cleanupStatus = statusResponse.data;
    state.storageStatus = storageResponse.data;
    renderCleanupSettings();
    renderCleanupStatus();
  } catch (error) {
    state.cleanupSettings = null;
    state.cleanupStatus = null;
    state.storageStatus = null;
    renderCleanupSettings();
    renderCleanupStatus();
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
    state.cleanupSettings = response.data;
    renderCleanupSettings();
    addNotice("Saved cleanup policy.", "success");
    await loadCleanupData();
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function runCleanupNow() {
  try {
    const response = await apiFetch("/api/admin/cleanup/run", { method: "POST" });
    state.cleanupStatus = response.data;
    renderCleanupStatus();
    addNotice("Triggered one manual cleanup run.", "success");
    await Promise.all([loadCleanupData(), loadLatest(), loadHistory(), loadFavorites()]);
  } catch (error) {
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
    addNotice("No clipboard text available to copy.", "info");
    return;
  }

  await copyItemText(state.latest);
}

async function copyItemText(item) {
  try {
    await navigator.clipboard.writeText(item.text || "");
    addNotice(`Copied clipboard item #${item.id}.`, "success");
  } catch (error) {
    addNotice("Copy failed. Your browser may block clipboard access.", "error");
  }
}

async function deleteLatest() {
  if (!state.latest) {
    return;
  }

  await deleteHistoryItem(state.latest.id);
}

async function deleteHistoryItem(id) {
  try {
    await apiFetch(`/api/clipboard/items/${id}`, { method: "DELETE" });
    addNotice(`Deleted clipboard item #${id}.`, "success");
    if (state.latest && state.latest.id === id) {
      state.latest = null;
    }
    await refreshClipboardViews();
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function setFavorite(id, favorite) {
  try {
    const response = await apiFetch(`/api/clipboard/items/${id}/favorite`, {
      method: favorite ? "POST" : "DELETE",
    });
    syncItemAcrossViews(response.data);
    renderLatest();
    renderHistory();
    await loadFavorites();
    addNotice(
      favorite ? `Favorited clipboard item #${id}.` : `Removed clipboard item #${id} from favorites.`,
      "success",
    );
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

async function updateLatestCategory() {
  if (!state.latest) {
    return;
  }
  await updateItemCategory(state.latest.id, nodes.latestCategorySelect.value);
}

async function updateItemCategory(id, category) {
  try {
    const response = await apiFetch(`/api/clipboard/items/${id}/category`, {
      method: "PATCH",
      body: JSON.stringify({ category }),
    });
    syncItemAcrossViews(response.data);
    renderLatest();
    renderHistory();
    await loadFavorites();
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
    addNotice(`Created category "${name}".`, "success");
    await loadCategories();
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function revokeDevice(id) {
  try {
    await apiFetch(`/api/auth/devices/${id}`, { method: "DELETE" });
    addNotice(`Revoked device #${id}.`, "success");
    await loadDevices();
  } catch (error) {
    addNotice(error.message, "error");
  }
}

async function refreshClipboardViews() {
  await Promise.all([loadLatest(), loadHistory(), loadFavorites(), loadCleanupData()]);
}

function syncItemAcrossViews(item) {
  if (state.latest && state.latest.id === item.id) {
    state.latest = item;
  }
  state.history = state.history.map((entry) => (entry.id === item.id ? item : entry));
  state.favorites = state.favorites
    .filter((entry) => entry.id !== item.id)
    .concat(item.is_favorite ? [item] : [])
    .sort((left, right) => right.id - left.id);
}

function renderLatest() {
  if (!state.latest) {
    nodes.latestEmpty.classList.remove("hidden");
    nodes.latestEmpty.textContent = state.token
      ? "No clipboard item loaded yet."
      : "Save an admin token or device token to load clipboard data.";
    nodes.latestCard.classList.add("hidden");
    return;
  }

  nodes.latestEmpty.classList.add("hidden");
  nodes.latestCard.classList.remove("hidden");
  nodes.latestId.textContent = `#${state.latest.id}`;
  nodes.latestCreated.textContent = state.latest.created_at;
  nodes.latestType.textContent = state.latest.type || "text";
  nodes.latestCategory.textContent = state.latest.category || "uncategorized";
  nodes.latestFavorite.textContent = state.latest.is_favorite ? "Favorite" : "Normal";
  nodes.latestText.textContent = state.latest.text || "";
  nodes.latestSourceName.textContent = state.latest.source_device_name || "Unknown Source";
  document.getElementById("toggle-latest-favorite-button").textContent = state.latest.is_favorite
    ? "Unfavorite"
    : "Favorite";
  populateCategorySelect(nodes.latestCategorySelect, state.latest.category || "");
}

function renderHistory() {
  renderClipboardCollection(
    nodes.historyList,
    state.history,
    state.token
      ? "No history available with the current token."
      : "Save an admin token or device token to view clipboard history.",
  );
}

function renderFavorites() {
  renderClipboardCollection(
    nodes.favoritesList,
    state.favorites,
    state.token
      ? "No favorites saved yet."
      : "Save an admin token or device token to view favorites.",
  );
}

function renderCleanupStatus() {
  const status = state.cleanupStatus;
  const storage = state.storageStatus;

  nodes.storageHistoryCount.textContent = storage ? String(storage.history_count) : "-";
  nodes.storageFavoriteCount.textContent = storage ? String(storage.favorite_count) : "-";
  nodes.storageTotalBytes.textContent = storage ? formatBytes(storage.total_bytes) : "-";
  nodes.storageFileBytes.textContent = storage ? formatBytes(storage.file_bytes) : "-";
  nodes.cleanupLastRun.textContent = status?.last_run_at || "Unavailable";

  if (!status) {
    nodes.cleanupLastResult.textContent = "Admin token required";
    return;
  }

  const bits = [
    `expired ${status.deleted_expired || 0}`,
    `count ${status.deleted_max_items || 0}`,
    `storage ${status.deleted_storage || 0}`,
  ];
  nodes.cleanupLastResult.textContent = status.last_error
    ? `Error: ${status.last_error}`
    : bits.join(" · ");
}

function renderCleanupSettings() {
  const settings = state.cleanupSettings;
  nodes.cleanupTTLHours.value = settings?.ttl_hours ?? "";
  nodes.cleanupMaxItems.value = settings?.max_items ?? "";
  nodes.cleanupMaxSizeMB.value = settings?.max_total_size_mb ?? "";
  nodes.cleanupIntervalMinutes.value = settings?.interval_minutes ?? "";
  nodes.cleanupEnabled.checked = Boolean(settings?.enabled);
}

function renderClipboardCollection(container, items, emptyMessage) {
  container.innerHTML = "";
  if (!items.length) {
    container.appendChild(renderEmpty(emptyMessage));
    return;
  }

  items.forEach((item) => {
    container.appendChild(renderClipboardItem(item));
  });
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

  const body = document.createElement("pre");
  body.textContent = item.text || "";

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

  const copyButton = document.createElement("button");
  copyButton.type = "button";
  copyButton.textContent = "Copy";
  copyButton.addEventListener("click", () => copyItemText(item));

  const favoriteButton = document.createElement("button");
  favoriteButton.type = "button";
  favoriteButton.className = "secondary";
  favoriteButton.textContent = item.is_favorite ? "Unfavorite" : "Favorite";
  favoriteButton.addEventListener("click", () => setFavorite(item.id, !item.is_favorite));

  const deleteButton = document.createElement("button");
  deleteButton.type = "button";
  deleteButton.className = "danger";
  deleteButton.textContent = "Delete";
  deleteButton.addEventListener("click", () => deleteHistoryItem(item.id));

  actions.append(copyButton, favoriteButton, deleteButton);
  li.append(title, tags, body, source, categoryRow, actions);
  return li;
}

function renderDevices(items) {
  nodes.devicesList.innerHTML = "";
  if (!items.length) {
    nodes.devicesList.appendChild(
      renderEmpty(
        state.token
          ? "No devices visible with the current token."
          : "Save an admin token to manage paired devices.",
      ),
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

  if (options.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(path, {
    method: options.method || "GET",
    headers,
    body: options.body,
  });

  let payload = {};
  try {
    payload = await response.json();
  } catch (error) {
    if (!response.ok) {
      throw new Error(`Request failed with status ${response.status}.`);
    }
    return payload;
  }

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
  if (bytes < 1024) {
    return `${bytes} B`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }
  if (bytes < 1024 * 1024 * 1024) {
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  }
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
}
