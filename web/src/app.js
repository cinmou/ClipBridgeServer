const storageKey = "clipbridge.webui.token";
const languageKey = "clipbridge.webui.language";
const webUISourceID = "web-ui";
const webUISourceName = "Web UI";
const routes = ["/", "/history", "/favorites", "/devices", "/settings", "/webdav"];

const translations = {
  en: {
    "hero.eyebrow": "v0.2.0-beta.1 · Cherwell",
    "hero.title": "ClipBridge",
    "hero.summary": "A browser clipboard client and management console for your embedded single-binary server.",
    "hero.badge.webdock": "WebDock",
    "hero.badge.embedded": "Embedded UI",
    "hero.badge.mobile": "Mobile Ready",
    "nav.dashboard": "Dashboard",
    "nav.history": "History",
    "nav.favorites": "Favorites",
    "nav.devices": "Devices",
    "nav.settings": "Settings",
    "nav.webdav": "WebDAV",
    "dashboard.quick.eyebrow": "Quick Clipboard",
    "dashboard.quick.title": "Push and Pull",
    "dashboard.quick.inputLabel": "Upload text to server clipboard",
    "dashboard.quick.placeholder": "Paste or type text here. It will be sent with source name Web UI.",
    "dashboard.quick.fetchLatest": "Fetch Latest",
    "dashboard.latest.eyebrow": "Latest Item",
    "dashboard.latest.title": "Clipboard Snapshot",
    "history.eyebrow": "History",
    "history.title": "Clipboard Timeline",
    "history.searchPlaceholder": "Search text, link, filename, source, or MIME type",
    "favorites.eyebrow": "Favorites",
    "favorites.title": "Protected Items",
    "favorites.note": "Favorites are protected from automatic cleanup until you remove the favorite flag.",
    "devices.pairing.eyebrow": "Pairing",
    "devices.pairing.title": "New Device",
    "devices.generatePairingCode": "Generate Pairing Code",
    "devices.qrPlaceholder": "Reserved area for QR pairing support in a later beta.",
    "devices.list.eyebrow": "Paired Devices",
    "devices.list.title": "Trust List",
    "settings.language.eyebrow": "Language",
    "settings.language.title": "Language",
    "settings.language.label": "Interface language",
    "settings.session.eyebrow": "Session",
    "settings.session.title": "Access Token",
    "settings.session.signedOut": "No token saved",
    "settings.session.tokenLoaded": "Token loaded",
    "settings.session.inputLabel": "Admin or device token",
    "settings.session.placeholder": "Paste token for API access",
    "settings.session.save": "Save Token",
    "settings.session.clear": "Clear",
    "settings.session.note": "The saved token is stored in this browser and used for API requests. Avoid using this on shared browsers.",
    "settings.session.statusLoaded": "The saved token is stored in this browser and used for API requests. Avoid using this on shared browsers.",
    "settings.session.statusEmpty": "No token saved in this browser yet.",
    "settings.security.eyebrow": "Security",
    "settings.security.title": "Security & Access",
    "settings.security.insecureWarning": "This server appears to be reachable over the network without HTTPS. Use only on trusted LANs or place it behind an HTTPS reverse proxy.",
    "settings.security.helper": "For public or remote access, use an HTTPS reverse proxy such as Caddy, Nginx, Traefik, Tailscale, or your NAS/router HTTPS gateway. Built-in TLS is optional for advanced users.",
    "settings.server.eyebrow": "Server",
    "settings.server.title": "Runtime Details",
    "settings.cleaner.eyebrow": "Cleanup",
    "settings.cleaner.title": "Retention Policy",
    "settings.cleaner.ttl": "TTL (hours)",
    "settings.cleaner.maxItems": "Max items",
    "settings.cleaner.maxSize": "Max total size (MB)",
    "settings.cleaner.interval": "Cleaner interval (minutes)",
    "settings.cleaner.enabled": "Cleaner enabled",
    "settings.cleaner.save": "Save Cleanup Settings",
    "settings.cleaner.runNow": "Run Cleanup Now",
    "settings.storage.eyebrow": "Storage",
    "settings.storage.title": "Current Usage",
    "settings.cleanupStatus.eyebrow": "Cleanup Status",
    "settings.cleanupStatus.title": "Latest Run",
    "webdav.config.eyebrow": "Configuration",
    "webdav.config.title": "WebDAV Preview",
    "webdav.enabled": "WebDAV sync enabled",
    "webdav.url": "Server URL",
    "webdav.username": "Username",
    "webdav.password": "Password",
    "webdav.basePath": "Base path",
    "webdav.save": "Save Config",
    "webdav.test": "Test Connection",
    "webdav.sync": "Manual Sync",
    "webdav.note": "WebDAV credentials are used only to connect to your storage provider. Passwords should not be shown again after saving and must not appear in logs.",
    "webdav.status.eyebrow": "Sync Status",
    "webdav.status.title": "Latest Result",
    "dialog.detailEyebrow": "Item Detail",
    "common.refresh": "Refresh",
    "common.search": "Search",
    "common.type": "Type",
    "common.source": "Source",
    "common.size": "Size",
    "common.created": "Created",
    "common.updated": "Updated",
    "common.close": "Close",
    "common.upload": "Upload",
    "common.download": "Download",
    "common.open": "Open",
    "common.copy": "Copy",
    "common.copyLatest": "Copy Latest",
    "common.revoke": "Revoke Device",
    "common.favorite": "Favorite",
    "common.unfavorite": "Unfavorite",
    "common.details": "Details",
    "common.delete": "Delete",
    "common.future": "Future",
    "common.empty": "Empty",
    "common.error": "Error",
    "common.unknown": "Unknown",
    "common.none": "None",
    "common.notAvailable": "Not available",
    "common.yes": "Yes",
    "common.no": "No",
    "filter.all": "All",
    "type.text": "Text",
    "type.link": "Link",
    "type.image": "Image",
    "type.file": "File",
    "status.healthy": "Healthy",
    "status.waiting": "Waiting",
    "status.localOnly": "Local only",
    "status.lanExposed": "LAN exposed",
    "status.customBind": "Custom bind",
    "status.tlsOn": "On",
    "status.tlsOff": "Off",
    "metric.service": "Service",
    "metric.serverUrl": "Server URL",
    "metric.connectedDevices": "Connected Devices",
    "metric.storedData": "Stored Data",
    "metric.version": "Version",
    "metric.dataDir": "Data Directory",
    "metric.databasePath": "Database Path",
    "metric.restartFields": "Restart Required Fields",
    "metric.historyItems": "History Items",
    "metric.favorites": "Favorites",
    "metric.totalBytes": "Total Bytes",
    "metric.fileBytes": "File Bytes",
    "metric.serverUrlHint": "Embedded Web UI origin",
    "metric.connectedDevicesHint": "Active non-revoked devices",
    "metric.storedDataHint": "{count} history items",
    "metric.versionHint": "Current server release",
    "metric.dataDirHint": "Embedded storage root",
    "metric.databasePathHint": "SQLite database location",
    "metric.restartFieldsHint": "Startup-only values",
    "metric.historyItemsHint": "Reverse chronological stack",
    "metric.favoritesHint": "Protected from cleanup",
    "metric.totalBytesHint": "All clipboard payloads",
    "metric.fileBytesHint": "Image and file payloads",
    "empty.saveTokenLatest": "Save a token to fetch the latest clipboard item.",
    "empty.noLatest": "No clipboard item is available yet.",
    "empty.saveTokenHistory": "Save a token to browse clipboard history.",
    "empty.noHistory": "No items match the current history filters.",
    "empty.saveTokenFavorites": "Save a token to load favorites.",
    "empty.noFavorites": "No favorite items yet.",
    "empty.saveTokenPairing": "Save an admin token to generate a pairing code.",
    "empty.noPairing": "Generate a short-lived pairing code for a new device.",
    "empty.saveTokenDevices": "Save an admin token to load paired devices.",
    "empty.noDevices": "No paired devices are recorded yet.",
    "empty.saveTokenSettings": "Save an admin token to load server settings.",
    "empty.saveTokenStorage": "Save an admin token to load storage status.",
    "empty.saveTokenCleanup": "Save an admin token to load cleanup status.",
    "empty.saveTokenWebDAV": "Save an admin token to load WebDAV sync status.",
    "toast.tokenSaved": "Token saved in this browser.",
    "toast.tokenCleared": "Token cleared.",
    "toast.enterText": "Enter text before uploading.",
    "toast.textUploaded": "Text uploaded to the clipboard stack.",
    "toast.pairingGenerated": "Pairing code generated.",
    "toast.confirmRevoke": "Revoke this device token now?",
    "toast.deviceRevoked": "Device revoked.",
    "toast.cleanupSaved": "Cleanup settings saved.",
    "toast.cleanupCompleted": "Cleanup run completed.",
    "toast.webdavSaved": "WebDAV settings saved.",
    "toast.webdavTested": "WebDAV connection test completed.",
    "toast.webdavSynced": "Manual WebDAV sync completed.",
    "toast.favoriteAdded": "Added to favorites.",
    "toast.favoriteRemoved": "Favorite removed.",
    "toast.confirmDelete": "Delete this clipboard item?",
    "toast.itemDeleted": "Item deleted.",
    "toast.copied": "Copied to browser clipboard.",
    "toast.copyFailed": "This item cannot be copied into the browser clipboard.",
    "label.bindAddress": "Bind address",
    "label.publicUrl": "Server URL / Public URL",
    "label.tlsStatus": "TLS status",
    "label.accessMode": "Access mode",
    "label.pairingPolicy": "Pairing code policy",
    "label.devicePolicy": "Device token policy",
    "label.pairingCode": "Pairing code",
    "label.expiresAt": "Expires at",
    "label.pairingUri": "Pairing URI",
    "label.id": "ID",
    "label.lastSeen": "Last seen",
    "label.revokedAt": "Revoked at",
    "label.reason": "Reason",
    "label.deletedExpired": "Deleted expired",
    "label.deletedMaxItems": "Deleted max items",
    "label.deletedStorage": "Deleted storage",
    "label.deletedFiles": "Deleted files",
    "label.lastError": "Latest error",
    "label.lastRun": "Last run",
    "label.lastSync": "Last sync",
    "label.lastSuccess": "Last success",
    "label.lastTested": "Last tested",
    "label.lastMessage": "Last message",
    "label.pushedItems": "Pushed items",
    "label.pulledItems": "Pulled items",
    "label.pushedFiles": "Pushed files",
    "label.pulledFiles": "Pulled files",
    "label.remoteItems": "Remote items",
    "label.localItems": "Local items",
    "label.conflictSkips": "Conflict skips",
    "label.type": "Type",
    "label.category": "Category",
    "label.sourceDevice": "Source device",
    "label.filename": "Filename",
    "label.mimeType": "MIME type",
    "label.fileSize": "File size",
    "label.sha256": "SHA-256",
    "label.created": "Created",
    "label.updated": "Updated",
    "policy.pairing": "5-minute expiry · single-use",
    "policy.device": "Generated during pairing",
    "policy.oneTimeCode": "One-time code",
    "policy.deepLinkReserved": "Reserved for app deep link flows",
    "device.active": "Active",
    "device.revoked": "Revoked",
    "detail.titlePrefix": "Detail",
    "detail.clipboardItem": "Clipboard Item",
    "pairing.expired": "Expired",
    "pairing.noExpiry": "No expiry reported",
    "pairing.expiresIn": "Expires in {minutes}m {seconds}s",
    "item.storedFile": "Stored file",
    "item.unnamedFile": "Unnamed file"
  },
  "zh-CN": {
    "hero.eyebrow": "v0.2.0-beta.1 · Cherwell",
    "hero.title": "ClipBridge",
    "hero.summary": "一个适用于内置单文件服务端的浏览器剪贴板客户端与管理界面。",
    "hero.badge.webdock": "WebDock",
    "hero.badge.embedded": "内置 Web UI",
    "hero.badge.mobile": "支持移动端",
    "nav.dashboard": "概览",
    "nav.history": "历史记录",
    "nav.favorites": "收藏",
    "nav.devices": "设备",
    "nav.settings": "设置",
    "nav.webdav": "WebDAV",
    "dashboard.quick.eyebrow": "快速剪贴板",
    "dashboard.quick.title": "上传与获取",
    "dashboard.quick.inputLabel": "上传文本到服务端剪贴板",
    "dashboard.quick.placeholder": "在这里粘贴或输入文本。来源会标记为 Web UI。",
    "dashboard.quick.fetchLatest": "获取最新内容",
    "dashboard.latest.eyebrow": "最新剪贴板",
    "dashboard.latest.title": "最新快照",
    "history.eyebrow": "历史记录",
    "history.title": "剪贴板时间线",
    "history.searchPlaceholder": "搜索文本、链接、文件名、来源或 MIME 类型",
    "favorites.eyebrow": "收藏",
    "favorites.title": "受保护的条目",
    "favorites.note": "收藏内容不会被自动清理，直到你取消收藏。",
    "devices.pairing.eyebrow": "配对",
    "devices.pairing.title": "新设备",
    "devices.generatePairingCode": "生成配对码",
    "devices.qrPlaceholder": "这里预留给后续版本的二维码配对支持。",
    "devices.list.eyebrow": "已配对设备",
    "devices.list.title": "信任列表",
    "settings.language.eyebrow": "语言",
    "settings.language.title": "语言",
    "settings.language.label": "界面语言",
    "settings.session.eyebrow": "会话",
    "settings.session.title": "访问令牌",
    "settings.session.signedOut": "未保存令牌",
    "settings.session.tokenLoaded": "令牌已加载",
    "settings.session.inputLabel": "管理员或设备令牌",
    "settings.session.placeholder": "粘贴用于 API 访问的令牌",
    "settings.session.save": "保存令牌",
    "settings.session.clear": "清除",
    "settings.session.note": "保存的令牌会存储在当前浏览器并用于 API 请求。请避免在共享浏览器中使用。",
    "settings.session.statusLoaded": "保存的令牌会存储在当前浏览器并用于 API 请求。请避免在共享浏览器中使用。",
    "settings.session.statusEmpty": "当前浏览器尚未保存令牌。",
    "settings.security.eyebrow": "安全",
    "settings.security.title": "安全与访问",
    "settings.security.insecureWarning": "此服务器看起来可以通过未启用 HTTPS 的网络访问。请仅在可信局域网中使用，或放在 HTTPS 反向代理之后。",
    "settings.security.helper": "如果需要公网或远程访问，请使用 Caddy、Nginx、Traefik、Tailscale，或 NAS/路由器自带的 HTTPS 网关作为反向代理。内置 TLS 适合高级用户按需启用。",
    "settings.server.eyebrow": "服务端",
    "settings.server.title": "运行时信息",
    "settings.cleaner.eyebrow": "清理",
    "settings.cleaner.title": "保留策略",
    "settings.cleaner.ttl": "TTL（小时）",
    "settings.cleaner.maxItems": "最大条目数",
    "settings.cleaner.maxSize": "最大总大小（MB）",
    "settings.cleaner.interval": "清理间隔（分钟）",
    "settings.cleaner.enabled": "启用清理器",
    "settings.cleaner.save": "保存清理设置",
    "settings.cleaner.runNow": "立即执行清理",
    "settings.storage.eyebrow": "存储",
    "settings.storage.title": "当前使用情况",
    "settings.cleanupStatus.eyebrow": "清理状态",
    "settings.cleanupStatus.title": "最近一次运行",
    "webdav.config.eyebrow": "配置",
    "webdav.config.title": "WebDAV 预览",
    "webdav.enabled": "启用 WebDAV 同步",
    "webdav.url": "服务地址",
    "webdav.username": "用户名",
    "webdav.password": "密码",
    "webdav.basePath": "基础路径",
    "webdav.save": "保存配置",
    "webdav.test": "测试连接",
    "webdav.sync": "手动同步",
    "webdav.note": "WebDAV 凭据仅用于连接你的存储服务。密码在保存后不应再次显示，也不应出现在日志中。",
    "webdav.status.eyebrow": "同步状态",
    "webdav.status.title": "最新结果",
    "dialog.detailEyebrow": "详情",
    "common.refresh": "刷新",
    "common.search": "搜索",
    "common.type": "类型",
    "common.source": "来源",
    "common.size": "大小",
    "common.created": "创建时间",
    "common.updated": "更新时间",
    "common.close": "关闭",
    "common.upload": "上传",
    "common.download": "下载",
    "common.open": "打开",
    "common.copy": "复制",
    "common.copyLatest": "复制最新内容",
    "common.revoke": "撤销设备",
    "common.favorite": "收藏",
    "common.unfavorite": "取消收藏",
    "common.details": "详情",
    "common.delete": "删除",
    "common.future": "后续功能",
    "common.empty": "空状态",
    "common.error": "错误",
    "common.unknown": "未知",
    "common.none": "无",
    "common.notAvailable": "不可用",
    "common.yes": "是",
    "common.no": "否",
    "filter.all": "全部",
    "type.text": "文本",
    "type.link": "链接",
    "type.image": "图片",
    "type.file": "文件",
    "status.healthy": "正常",
    "status.waiting": "等待中",
    "status.localOnly": "仅本机",
    "status.lanExposed": "局域网可访问",
    "status.customBind": "自定义绑定",
    "status.tlsOn": "开启",
    "status.tlsOff": "关闭",
    "metric.service": "服务状态",
    "metric.serverUrl": "服务地址",
    "metric.connectedDevices": "在线设备",
    "metric.storedData": "已存数据",
    "metric.version": "版本",
    "metric.dataDir": "数据目录",
    "metric.databasePath": "数据库路径",
    "metric.restartFields": "需重启字段",
    "metric.historyItems": "历史条目",
    "metric.favorites": "收藏数量",
    "metric.totalBytes": "总字节数",
    "metric.fileBytes": "文件字节数",
    "metric.serverUrlHint": "当前内置 Web UI 地址",
    "metric.connectedDevicesHint": "未撤销的活跃设备",
    "metric.storedDataHint": "{count} 条历史记录",
    "metric.versionHint": "当前服务端版本",
    "metric.dataDirHint": "内置存储根目录",
    "metric.databasePathHint": "SQLite 数据库位置",
    "metric.restartFieldsHint": "仅启动时生效",
    "metric.historyItemsHint": "按时间倒序排列",
    "metric.favoritesHint": "不会被自动清理",
    "metric.totalBytesHint": "所有剪贴板负载",
    "metric.fileBytesHint": "图片与文件负载",
    "empty.saveTokenLatest": "请先保存令牌，然后再读取最新剪贴板。",
    "empty.noLatest": "目前还没有剪贴板内容。",
    "empty.saveTokenHistory": "请先保存令牌，然后再查看历史记录。",
    "empty.noHistory": "当前筛选条件下没有匹配内容。",
    "empty.saveTokenFavorites": "请先保存令牌，然后再查看收藏。",
    "empty.noFavorites": "当前还没有收藏内容。",
    "empty.saveTokenPairing": "请先保存管理员令牌，然后再生成配对码。",
    "empty.noPairing": "为新设备生成一个短时有效的配对码。",
    "empty.saveTokenDevices": "请先保存管理员令牌，然后再查看设备列表。",
    "empty.noDevices": "当前还没有已配对设备。",
    "empty.saveTokenSettings": "请先保存管理员令牌，然后再读取服务端设置。",
    "empty.saveTokenStorage": "请先保存管理员令牌，然后再读取存储状态。",
    "empty.saveTokenCleanup": "请先保存管理员令牌，然后再读取清理状态。",
    "empty.saveTokenWebDAV": "请先保存管理员令牌，然后再读取 WebDAV 同步状态。",
    "toast.tokenSaved": "令牌已保存到当前浏览器。",
    "toast.tokenCleared": "令牌已清除。",
    "toast.enterText": "请先输入要上传的文本。",
    "toast.textUploaded": "文本已上传到剪贴板堆栈。",
    "toast.pairingGenerated": "配对码已生成。",
    "toast.confirmRevoke": "确定要撤销这个设备令牌吗？",
    "toast.deviceRevoked": "设备已撤销。",
    "toast.cleanupSaved": "清理设置已保存。",
    "toast.cleanupCompleted": "清理已执行完成。",
    "toast.webdavSaved": "WebDAV 设置已保存。",
    "toast.webdavTested": "WebDAV 连接测试已完成。",
    "toast.webdavSynced": "手动 WebDAV 同步已完成。",
    "toast.favoriteAdded": "已加入收藏。",
    "toast.favoriteRemoved": "已取消收藏。",
    "toast.confirmDelete": "确定要删除这个剪贴板条目吗？",
    "toast.itemDeleted": "条目已删除。",
    "toast.copied": "已复制到浏览器剪贴板。",
    "toast.copyFailed": "这个条目无法写入浏览器剪贴板。",
    "label.bindAddress": "绑定地址",
    "label.publicUrl": "服务地址 / 公网地址",
    "label.tlsStatus": "TLS 状态",
    "label.accessMode": "访问模式",
    "label.pairingPolicy": "配对码策略",
    "label.devicePolicy": "设备令牌策略",
    "label.pairingCode": "配对码",
    "label.expiresAt": "过期时间",
    "label.pairingUri": "配对 URI",
    "label.id": "ID",
    "label.lastSeen": "最近在线",
    "label.revokedAt": "撤销时间",
    "label.reason": "原因",
    "label.deletedExpired": "过期删除数",
    "label.deletedMaxItems": "超条目删除数",
    "label.deletedStorage": "超存储删除数",
    "label.deletedFiles": "删除文件数",
    "label.lastError": "最近错误",
    "label.lastRun": "最近运行",
    "label.lastSync": "最近同步",
    "label.lastSuccess": "最近成功",
    "label.lastTested": "最近测试",
    "label.lastMessage": "最近消息",
    "label.pushedItems": "推送条目",
    "label.pulledItems": "拉取条目",
    "label.pushedFiles": "推送文件",
    "label.pulledFiles": "拉取文件",
    "label.remoteItems": "远端条目",
    "label.localItems": "本地条目",
    "label.conflictSkips": "冲突跳过数",
    "label.type": "类型",
    "label.category": "分类",
    "label.sourceDevice": "来源设备",
    "label.filename": "文件名",
    "label.mimeType": "MIME 类型",
    "label.fileSize": "文件大小",
    "label.sha256": "SHA-256",
    "label.created": "创建时间",
    "label.updated": "更新时间",
    "policy.pairing": "5 分钟有效 · 单次使用",
    "policy.device": "配对时生成",
    "policy.oneTimeCode": "一次性配对码",
    "policy.deepLinkReserved": "预留给客户端深链接使用",
    "device.active": "有效",
    "device.revoked": "已撤销",
    "detail.titlePrefix": "详情",
    "detail.clipboardItem": "剪贴板条目",
    "pairing.expired": "已过期",
    "pairing.noExpiry": "未提供过期时间",
    "pairing.expiresIn": "{minutes} 分 {seconds} 秒后过期",
    "item.storedFile": "已存储文件",
    "item.unnamedFile": "未命名文件"
  }
};

const state = {
  token: localStorage.getItem(storageKey) || "",
  language: detectInitialLanguage(),
  route: normalizeRoute(window.location.pathname),
  health: null,
  latest: null,
  history: [],
  favorites: [],
  devices: [],
  settings: null,
  cleanup: null,
  cleanupStatus: null,
  storageStatus: null,
  webdavSettings: null,
  webdavStatus: null,
  pairing: null,
  historySearch: "",
  historyType: "all",
  activeDetailItem: null
};

const nodes = {
  tokenInput: document.getElementById("token-input"),
  tokenStatus: document.getElementById("token-status"),
  sessionBadge: document.getElementById("session-badge"),
  languageSelect: document.getElementById("language-select"),
  navLinks: Array.from(document.querySelectorAll(".nav-link")),
  pages: Array.from(document.querySelectorAll(".page")),
  dashboardOverview: document.getElementById("dashboard-overview"),
  quickTextInput: document.getElementById("quick-text-input"),
  latestPanel: document.getElementById("latest-panel"),
  historySearchInput: document.getElementById("history-search-input"),
  historyTypeFilter: document.getElementById("history-type-filter"),
  historyPanel: document.getElementById("history-panel"),
  favoritesPanel: document.getElementById("favorites-panel"),
  pairingPanel: document.getElementById("pairing-panel"),
  devicesPanel: document.getElementById("devices-panel"),
  securityPanel: document.getElementById("security-panel"),
  serverSettingsPanel: document.getElementById("server-settings-panel"),
  cleanupForm: document.getElementById("cleanup-form"),
  cleanupTTLHours: document.getElementById("cleanup-ttl-hours"),
  cleanupMaxItems: document.getElementById("cleanup-max-items"),
  cleanupMaxSizeMB: document.getElementById("cleanup-max-size-mb"),
  cleanupIntervalMinutes: document.getElementById("cleanup-interval-minutes"),
  cleanupEnabled: document.getElementById("cleanup-enabled"),
  storagePanel: document.getElementById("storage-panel"),
  cleanupStatusPanel: document.getElementById("cleanup-status-panel"),
  webdavForm: document.getElementById("webdav-form"),
  webdavEnabled: document.getElementById("webdav-enabled"),
  webdavURL: document.getElementById("webdav-url"),
  webdavUsername: document.getElementById("webdav-username"),
  webdavPassword: document.getElementById("webdav-password"),
  webdavBasePath: document.getElementById("webdav-base-path"),
  webdavStatusPanel: document.getElementById("webdav-status-panel"),
  detailDialog: document.getElementById("detail-dialog"),
  detailTitle: document.getElementById("detail-title"),
  detailBody: document.getElementById("detail-body"),
  toastRegion: document.getElementById("toast-region")
};

nodes.tokenInput.value = state.token;
nodes.languageSelect.value = state.language;
applyStaticTranslations();
updateSessionUI();
renderRoute();
renderAll();

document.getElementById("save-token-button").addEventListener("click", saveToken);
document.getElementById("clear-token-button").addEventListener("click", clearToken);
document.getElementById("refresh-dashboard-button").addEventListener("click", refreshDashboard);
document.getElementById("upload-text-button").addEventListener("click", uploadQuickText);
document.getElementById("fetch-latest-button").addEventListener("click", loadLatest);
document.getElementById("refresh-history-button").addEventListener("click", loadHistory);
document.getElementById("refresh-favorites-button").addEventListener("click", loadFavorites);
document.getElementById("generate-pairing-button").addEventListener("click", generatePairingCode);
document.getElementById("refresh-devices-button").addEventListener("click", loadDevices);
document.getElementById("refresh-settings-button").addEventListener("click", refreshSettingsPage);
document.getElementById("run-cleanup-button").addEventListener("click", runCleanupNow);
document.getElementById("refresh-webdav-button").addEventListener("click", refreshWebDAVPage);
document.getElementById("test-webdav-button").addEventListener("click", testWebDAV);
document.getElementById("sync-webdav-button").addEventListener("click", syncWebDAV);
document.getElementById("close-detail-button").addEventListener("click", () => nodes.detailDialog.close());
nodes.languageSelect.addEventListener("change", onLanguageChange);

window.addEventListener("popstate", () => {
  state.route = normalizeRoute(window.location.pathname);
  renderRoute();
});

nodes.navLinks.forEach((button) => {
  button.addEventListener("click", () => navigate(button.dataset.route || "/"));
});

nodes.historySearchInput.addEventListener("input", (event) => {
  state.historySearch = event.target.value.trim().toLowerCase();
  renderHistoryPanel();
});

nodes.historyTypeFilter.addEventListener("change", (event) => {
  state.historyType = event.target.value;
  renderHistoryPanel();
});

nodes.cleanupForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  await saveCleanupSettings();
});

nodes.webdavForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  await saveWebDAVSettings();
});

loadInitialData();
window.setInterval(updateCountdowns, 1000);

function detectInitialLanguage() {
  const saved = localStorage.getItem(languageKey);
  if (saved === "en" || saved === "zh-CN") {
    return saved;
  }
  return navigator.language && navigator.language.toLowerCase().startsWith("zh") ? "zh-CN" : "en";
}

function t(key, vars = {}) {
  const dict = translations[state.language] || translations.en;
  let value = dict[key] || translations.en[key] || key;
  for (const [name, replacement] of Object.entries(vars)) {
    value = value.replaceAll(`{${name}}`, String(replacement));
  }
  return value;
}

function applyStaticTranslations() {
  document.documentElement.lang = state.language;
  document.querySelectorAll("[data-i18n]").forEach((node) => {
    node.textContent = t(node.dataset.i18n);
  });
  document.querySelectorAll("[data-i18n-placeholder]").forEach((node) => {
    node.setAttribute("placeholder", t(node.dataset.i18nPlaceholder));
  });
}

function onLanguageChange(event) {
  const next = event.target.value === "zh-CN" ? "zh-CN" : "en";
  state.language = next;
  localStorage.setItem(languageKey, next);
  applyStaticTranslations();
  updateSessionUI();
  renderRoute();
  renderAll();
}

async function loadInitialData() {
  await loadHealth();
  await loadProtectedData();
}

async function loadProtectedData() {
  if (!state.token) {
    resetProtectedState();
    return;
  }

  await Promise.allSettled([
    loadLatest(),
    loadHistory(),
    loadFavorites(),
    loadDevices(),
    loadSettings(),
    loadCleanupSettings(),
    loadStorageStatus(),
    loadCleanupStatus(),
    loadWebDAVSettings(),
    loadWebDAVStatus()
  ]);
}

function resetProtectedState() {
  state.latest = null;
  state.history = [];
  state.favorites = [];
  state.devices = [];
  state.settings = null;
  state.cleanup = null;
  state.cleanupStatus = null;
  state.storageStatus = null;
  state.webdavSettings = null;
  state.webdavStatus = null;
  state.pairing = null;
  renderAll();
}

function renderAll() {
  renderDashboardOverview();
  renderLatestPanel();
  renderHistoryPanel();
  renderFavoritesPanel();
  renderPairingPanel();
  renderDevicesPanel();
  renderSecurityPanel();
  renderServerSettingsPanel();
  renderStoragePanel();
  renderCleanupStatusPanel();
  renderWebDAVStatusPanel();
}

function navigate(route) {
  const nextRoute = normalizeRoute(route);
  if (nextRoute === state.route) {
    return;
  }
  state.route = nextRoute;
  window.history.pushState({}, "", nextRoute);
  renderRoute();
}

function normalizeRoute(pathname) {
  if (!pathname || pathname === "/") {
    return "/";
  }
  const trimmed = pathname.replace(/\/+$/, "");
  return routes.includes(trimmed) ? trimmed : "/";
}

function renderRoute() {
  nodes.navLinks.forEach((button) => {
    button.classList.toggle("is-active", (button.dataset.route || "/") === state.route);
  });
  nodes.pages.forEach((page) => {
    page.classList.toggle("is-active", page.dataset.page === state.route);
  });
}

function saveToken() {
  state.token = nodes.tokenInput.value.trim();
  if (state.token) {
    localStorage.setItem(storageKey, state.token);
    toast(t("toast.tokenSaved"), "success");
  } else {
    localStorage.removeItem(storageKey);
    toast(t("toast.tokenCleared"), "success");
  }
  updateSessionUI();
  loadProtectedData();
}

function clearToken() {
  nodes.tokenInput.value = "";
  state.token = "";
  localStorage.removeItem(storageKey);
  updateSessionUI();
  resetProtectedState();
  toast(t("toast.tokenCleared"), "success");
}

function updateSessionUI() {
  const signedIn = Boolean(state.token);
  nodes.sessionBadge.textContent = signedIn ? t("settings.session.tokenLoaded") : t("settings.session.signedOut");
  nodes.sessionBadge.className = `badge ${signedIn ? "badge-accent" : "badge-muted"}`.trim();
  nodes.tokenStatus.textContent = signedIn ? t("settings.session.statusLoaded") : t("settings.session.statusEmpty");
}

async function refreshDashboard() {
  await Promise.allSettled([loadHealth(), loadLatest(), loadDevices(), loadStorageStatus()]);
}

async function refreshSettingsPage() {
  await Promise.allSettled([loadHealth(), loadSettings(), loadCleanupSettings(), loadStorageStatus(), loadCleanupStatus()]);
}

async function refreshWebDAVPage() {
  await Promise.allSettled([loadWebDAVSettings(), loadWebDAVStatus()]);
}

async function loadHealth() {
  try {
    state.health = await apiFetch("/api/health", { auth: false });
  } catch (error) {
    state.health = null;
    toast(error.message, "error");
  }
  renderDashboardOverview();
  renderSecurityPanel();
  renderServerSettingsPanel();
}

async function loadLatest() {
  if (!state.token) {
    state.latest = null;
    renderLatestPanel();
    return;
  }

  try {
    state.latest = await apiFetch("/api/clipboard/latest");
  } catch (error) {
    state.latest = null;
    if (!isNotFound(error)) {
      toast(error.message, "error");
    }
  }
  renderLatestPanel();
}

async function loadHistory() {
  if (!state.token) {
    state.history = [];
    renderHistoryPanel();
    return;
  }

  try {
    const response = await apiFetch("/api/clipboard/history");
    state.history = response.items || [];
  } catch (error) {
    state.history = [];
    toast(error.message, "error");
  }
  renderHistoryPanel();
}

async function loadFavorites() {
  if (!state.token) {
    state.favorites = [];
    renderFavoritesPanel();
    return;
  }

  try {
    const response = await apiFetch("/api/favorites");
    state.favorites = response.items || [];
  } catch (error) {
    state.favorites = [];
    toast(error.message, "error");
  }
  renderFavoritesPanel();
}

async function loadDevices() {
  if (!state.token) {
    state.devices = [];
    renderDevicesPanel();
    return;
  }

  try {
    const response = await apiFetch("/api/auth/devices");
    state.devices = response.items || [];
  } catch (error) {
    state.devices = [];
    renderDevicesPanel(error.message);
    return;
  }
  renderDevicesPanel();
}

async function loadSettings() {
  if (!state.token) {
    state.settings = null;
    renderSecurityPanel();
    renderServerSettingsPanel();
    return;
  }

  try {
    state.settings = await apiFetch("/api/settings");
  } catch (error) {
    state.settings = null;
    renderSecurityPanel();
    renderServerSettingsPanel(error.message);
    return;
  }
  renderSecurityPanel();
  renderServerSettingsPanel();
}

async function loadCleanupSettings() {
  if (!state.token) {
    state.cleanup = null;
    populateCleanupForm();
    return;
  }

  try {
    state.cleanup = await apiFetch("/api/settings/cleanup");
  } catch (error) {
    state.cleanup = null;
    toast(error.message, "error");
  }
  populateCleanupForm();
}

async function loadStorageStatus() {
  if (!state.token) {
    state.storageStatus = null;
    renderStoragePanel();
    return;
  }

  try {
    state.storageStatus = await apiFetch("/api/admin/storage/status");
  } catch (error) {
    state.storageStatus = null;
    renderStoragePanel(error.message);
    return;
  }
  renderStoragePanel();
}

async function loadCleanupStatus() {
  if (!state.token) {
    state.cleanupStatus = null;
    renderCleanupStatusPanel();
    return;
  }

  try {
    state.cleanupStatus = await apiFetch("/api/admin/cleanup/status");
  } catch (error) {
    state.cleanupStatus = null;
    renderCleanupStatusPanel(error.message);
    return;
  }
  renderCleanupStatusPanel();
}

async function loadWebDAVSettings() {
  if (!state.token) {
    state.webdavSettings = null;
    populateWebDAVForm();
    return;
  }

  try {
    state.webdavSettings = await apiFetch("/api/settings/webdav");
  } catch (error) {
    state.webdavSettings = null;
    toast(error.message, "error");
  }
  populateWebDAVForm();
}

async function loadWebDAVStatus() {
  if (!state.token) {
    state.webdavStatus = null;
    renderWebDAVStatusPanel();
    return;
  }

  try {
    state.webdavStatus = await apiFetch("/api/admin/webdav/status");
  } catch (error) {
    state.webdavStatus = null;
    renderWebDAVStatusPanel(error.message);
    return;
  }
  renderWebDAVStatusPanel();
}

async function uploadQuickText() {
  const content = nodes.quickTextInput.value.trim();
  if (!content) {
    toast(t("toast.enterText"), "error");
    return;
  }

  try {
    await apiFetch("/api/clipboard/text", {
      method: "POST",
      body: JSON.stringify({
        content,
        source_device_id: webUISourceID,
        source_device_name: webUISourceName
      })
    });
    nodes.quickTextInput.value = "";
    toast(t("toast.textUploaded"), "success");
    await Promise.allSettled([loadLatest(), loadHistory(), loadFavorites(), loadStorageStatus()]);
  } catch (error) {
    toast(error.message, "error");
  }
}

async function generatePairingCode() {
  try {
    state.pairing = await apiFetch("/api/auth/pairing-codes", { method: "POST" });
    renderPairingPanel();
    toast(t("toast.pairingGenerated"), "success");
  } catch (error) {
    renderPairingPanel(error.message);
  }
}

async function revokeDevice(id) {
  if (!window.confirm(t("toast.confirmRevoke"))) {
    return;
  }

  try {
    await apiFetch(`/api/auth/devices/${id}`, { method: "DELETE" });
    toast(t("toast.deviceRevoked"), "success");
    await loadDevices();
  } catch (error) {
    toast(error.message, "error");
  }
}

async function saveCleanupSettings() {
  try {
    const payload = {
      ttl_hours: Number(nodes.cleanupTTLHours.value),
      max_items: Number(nodes.cleanupMaxItems.value),
      max_total_size_mb: Number(nodes.cleanupMaxSizeMB.value),
      interval_minutes: Number(nodes.cleanupIntervalMinutes.value),
      enabled: nodes.cleanupEnabled.checked
    };
    state.cleanup = await apiFetch("/api/settings/cleanup", {
      method: "PATCH",
      body: JSON.stringify(payload)
    });
    populateCleanupForm();
    toast(t("toast.cleanupSaved"), "success");
  } catch (error) {
    toast(error.message, "error");
  }
}

async function runCleanupNow() {
  try {
    state.cleanupStatus = await apiFetch("/api/admin/cleanup/run", { method: "POST" });
    renderCleanupStatusPanel();
    await loadStorageStatus();
    toast(t("toast.cleanupCompleted"), "success");
  } catch (error) {
    toast(error.message, "error");
  }
}

async function saveWebDAVSettings() {
  try {
    state.webdavSettings = await apiFetch("/api/settings/webdav", {
      method: "PATCH",
      body: JSON.stringify({
        enabled: nodes.webdavEnabled.checked,
        url: nodes.webdavURL.value.trim(),
        username: nodes.webdavUsername.value.trim(),
        password: nodes.webdavPassword.value,
        base_path: nodes.webdavBasePath.value.trim()
      })
    });
    populateWebDAVForm();
    toast(t("toast.webdavSaved"), "success");
  } catch (error) {
    toast(error.message, "error");
  }
}

async function testWebDAV() {
  try {
    state.webdavStatus = await apiFetch("/api/admin/webdav/test", { method: "POST" });
    renderWebDAVStatusPanel();
    toast(t("toast.webdavTested"), "success");
  } catch (error) {
    toast(error.message, "error");
  }
}

async function syncWebDAV() {
  try {
    state.webdavStatus = await apiFetch("/api/admin/webdav/sync", { method: "POST" });
    renderWebDAVStatusPanel();
    toast(t("toast.webdavSynced"), "success");
  } catch (error) {
    toast(error.message, "error");
  }
}

async function toggleFavorite(item) {
  const method = item.is_favorite ? "DELETE" : "POST";
  try {
    const updated = await apiFetch(`/api/clipboard/items/${item.id}/favorite`, { method });
    patchItemAcrossCollections(updated);
    toast(item.is_favorite ? t("toast.favoriteRemoved") : t("toast.favoriteAdded"), "success");
  } catch (error) {
    toast(error.message, "error");
  }
}

async function deleteItem(item) {
  if (!window.confirm(t("toast.confirmDelete"))) {
    return;
  }

  try {
    await apiFetch(`/api/clipboard/items/${item.id}`, { method: "DELETE" });
    state.history = state.history.filter((entry) => entry.id !== item.id);
    state.favorites = state.favorites.filter((entry) => entry.id !== item.id);
    if (state.latest && state.latest.id === item.id) {
      state.latest = null;
    }
    renderLatestPanel();
    renderHistoryPanel();
    renderFavoritesPanel();
    toast(t("toast.itemDeleted"), "success");
  } catch (error) {
    toast(error.message, "error");
  }
}

function patchItemAcrossCollections(updated) {
  state.history = state.history.map((entry) => (entry.id === updated.id ? updated : entry));
  if (updated.is_favorite) {
    const favoriteExists = state.favorites.some((entry) => entry.id === updated.id);
    state.favorites = favoriteExists ? state.favorites.map((entry) => (entry.id === updated.id ? updated : entry)) : [updated, ...state.favorites];
  } else {
    state.favorites = state.favorites.filter((entry) => entry.id !== updated.id);
  }
  if (state.latest && state.latest.id === updated.id) {
    state.latest = updated;
  }
  renderLatestPanel();
  renderHistoryPanel();
  renderFavoritesPanel();
}

async function handleCopy(item) {
  const text = item.type === "link" ? item.url : item.text;
  if (!text) {
    toast(t("toast.copyFailed"), "error");
    return;
  }

  try {
    await navigator.clipboard.writeText(text);
    toast(t("toast.copied"), "success");
  } catch (error) {
    toast(error.message || t("toast.copyFailed"), "error");
  }
}

function openItemDetails(item) {
  state.activeDetailItem = item;
  nodes.detailTitle.textContent = detailTitle(item);
  nodes.detailBody.innerHTML = renderDetailBody(item);
  nodes.detailDialog.showModal();

  const detailCopyButton = document.getElementById("detail-copy-button");
  const detailFavoriteButton = document.getElementById("detail-favorite-button");
  const detailDeleteButton = document.getElementById("detail-delete-button");

  if (detailCopyButton) {
    detailCopyButton.addEventListener("click", () => handleCopy(item));
  }
  if (detailFavoriteButton) {
    detailFavoriteButton.addEventListener("click", async () => {
      await toggleFavorite(item);
      nodes.detailDialog.close();
    });
  }
  if (detailDeleteButton) {
    detailDeleteButton.addEventListener("click", async () => {
      await deleteItem(item);
      nodes.detailDialog.close();
    });
  }
}

function renderDashboardOverview() {
  const cards = [
    metricCard(t("metric.service"), state.health?.ok ? t("status.healthy") : t("status.waiting"), state.health?.version || t("common.notAvailable")),
    metricCard(t("metric.serverUrl"), window.location.origin, t("metric.serverUrlHint")),
    metricCard(t("metric.connectedDevices"), String(activeDeviceCount()), t("metric.connectedDevicesHint")),
    metricCard(t("metric.storedData"), formatBytes(state.storageStatus?.total_bytes || 0), t("metric.storedDataHint", { count: state.storageStatus?.history_count || 0 }))
  ];
  nodes.dashboardOverview.innerHTML = cards.join("");
}

function renderLatestPanel() {
  if (!state.token) {
    nodes.latestPanel.innerHTML = renderEmptyState(t("empty.saveTokenLatest"));
    return;
  }
  if (!state.latest) {
    nodes.latestPanel.innerHTML = renderEmptyState(t("empty.noLatest"));
    return;
  }
  nodes.latestPanel.innerHTML = renderItemCard(state.latest, { latest: true });
  bindItemCardActions(nodes.latestPanel, [state.latest]);
}

function renderHistoryPanel() {
  if (!state.token) {
    nodes.historyPanel.innerHTML = renderEmptyState(t("empty.saveTokenHistory"));
    return;
  }
  const filtered = filteredHistory();
  if (!filtered.length) {
    nodes.historyPanel.innerHTML = renderEmptyState(t("empty.noHistory"));
    return;
  }
  nodes.historyPanel.innerHTML = `<div class="item-list">${filtered.map((item) => renderItemCard(item)).join("")}</div>`;
  bindItemCardActions(nodes.historyPanel, filtered);
}

function renderFavoritesPanel() {
  if (!state.token) {
    nodes.favoritesPanel.innerHTML = renderEmptyState(t("empty.saveTokenFavorites"));
    return;
  }
  if (!state.favorites.length) {
    nodes.favoritesPanel.innerHTML = renderEmptyState(t("empty.noFavorites"));
    return;
  }
  nodes.favoritesPanel.innerHTML = `<div class="item-list">${state.favorites.map((item) => renderItemCard(item)).join("")}</div>`;
  bindItemCardActions(nodes.favoritesPanel, state.favorites);
}

function renderPairingPanel(errorMessage = "") {
  if (errorMessage) {
    nodes.pairingPanel.innerHTML = renderErrorState(errorMessage);
    return;
  }
  if (!state.token) {
    nodes.pairingPanel.innerHTML = renderEmptyState(t("empty.saveTokenPairing"));
    return;
  }
  if (!state.pairing) {
    nodes.pairingPanel.innerHTML = renderEmptyState(t("empty.noPairing"));
    return;
  }

  nodes.pairingPanel.innerHTML = `
    <div class="pairing-grid">
      ${metricCard(t("metric.serverUrl"), window.location.origin, t("metric.serverUrlHint"))}
      ${metricCard(t("label.pairingCode"), escapeHTML(state.pairing.pairing_code), countdownLabel(state.pairing.expires_at))}
      ${metricCard(t("label.expiresAt"), formatDateTime(state.pairing.expires_at), t("policy.oneTimeCode"))}
      ${metricCard(t("label.pairingUri"), `<span class="mono">${escapeHTML(state.pairing.pairing_uri)}</span>`, t("policy.deepLinkReserved"))}
    </div>
  `;
}

function renderDevicesPanel(errorMessage = "") {
  if (errorMessage) {
    nodes.devicesPanel.innerHTML = renderErrorState(errorMessage);
    return;
  }
  if (!state.token) {
    nodes.devicesPanel.innerHTML = renderEmptyState(t("empty.saveTokenDevices"));
    return;
  }
  if (!state.devices.length) {
    nodes.devicesPanel.innerHTML = renderEmptyState(t("empty.noDevices"));
    return;
  }

  nodes.devicesPanel.innerHTML = `<div class="item-list">${state.devices.map((device) => `
    <article class="item-card device-card" data-device-id="${device.id}">
      <div class="item-header">
        <div class="detail-stack">
          <h3>${escapeHTML(device.name || `${t("nav.devices")} ${device.id}`)}</h3>
          <div class="item-badges">
            <span class="badge">${device.revoked_at ? t("device.revoked") : t("device.active")}</span>
            <span class="badge badge-muted">${escapeHTML(`${t("label.id")} ${device.id}`)}</span>
          </div>
        </div>
        <button class="button ${device.revoked_at ? "button-ghost" : "button-danger"}" type="button" data-action="revoke-device" ${device.revoked_at ? "disabled" : ""}>${t("common.revoke")}</button>
      </div>
      <div class="meta-list">
        ${metaLine(t("label.created"), formatDateTime(device.created_at))}
        ${metaLine(t("label.lastSeen"), formatDateTime(device.last_seen_at))}
        ${metaLine(t("label.revokedAt"), formatDateTime(device.revoked_at))}
      </div>
    </article>
  `).join("")}</div>`;

  nodes.devicesPanel.querySelectorAll("[data-action='revoke-device']").forEach((button) => {
    button.addEventListener("click", () => revokeDevice(Number(button.closest("[data-device-id]").dataset.deviceId)));
  });
}

function renderSecurityPanel() {
  const bindAddress = state.settings?.startup?.host || t("common.notAvailable");
  const publicURL = window.location.origin;
  const tlsOn = window.location.protocol === "https:";
  const accessMode = bindAddress === "127.0.0.1" ? t("status.localOnly") : bindAddress === "0.0.0.0" ? t("status.lanExposed") : t("status.customBind");
  const showInsecureLANWarning = !tlsOn && !isLocalBrowserOrigin();

  nodes.securityPanel.innerHTML = `
    ${showInsecureLANWarning ? `<div class="warning-banner">${escapeHTML(t("settings.security.insecureWarning"))}</div>` : ""}
    <div class="meta-list">
      ${metaLine(t("label.bindAddress"), escapeHTML(bindAddress))}
      ${metaLine(t("label.publicUrl"), `<span class="mono">${escapeHTML(publicURL)}</span>`)}
      ${metaLine(t("label.tlsStatus"), tlsOn ? t("status.tlsOn") : t("status.tlsOff"))}
      ${metaLine(t("label.accessMode"), accessMode)}
      ${metaLine(t("label.pairingPolicy"), t("policy.pairing"))}
      ${metaLine(t("label.devicePolicy"), t("policy.device"))}
    </div>
    <p class="support-text">${escapeHTML(t("settings.security.helper"))}</p>
  `;
}

function renderServerSettingsPanel(errorMessage = "") {
  if (errorMessage) {
    nodes.serverSettingsPanel.innerHTML = renderErrorState(errorMessage);
    return;
  }
  if (!state.settings) {
    nodes.serverSettingsPanel.innerHTML = renderEmptyState(t("empty.saveTokenSettings"));
    return;
  }

  nodes.serverSettingsPanel.innerHTML = `
    <div class="stats-grid">
      ${metricCard(t("metric.version"), state.health?.version || "v0.2.0-beta.1", t("metric.versionHint"))}
      ${metricCard(t("metric.dataDir"), `<span class="mono">${escapeHTML(state.settings.startup?.data_dir || "-")}</span>`, t("metric.dataDirHint"))}
      ${metricCard(t("metric.databasePath"), `<span class="mono">${escapeHTML(state.settings.startup?.database_path || "-")}</span>`, t("metric.databasePathHint"))}
      ${metricCard(t("metric.restartFields"), escapeHTML((state.settings.restart_required_fields || []).join(", ") || t("common.none")), t("metric.restartFieldsHint"))}
    </div>
  `;
}

function renderStoragePanel(errorMessage = "") {
  if (errorMessage) {
    nodes.storagePanel.innerHTML = renderErrorState(errorMessage);
    return;
  }
  if (!state.storageStatus) {
    nodes.storagePanel.innerHTML = renderEmptyState(t("empty.saveTokenStorage"));
    return;
  }

  nodes.storagePanel.innerHTML = `
    <div class="stats-grid">
      ${metricCard(t("metric.historyItems"), String(state.storageStatus.history_count || 0), t("metric.historyItemsHint"))}
      ${metricCard(t("metric.favorites"), String(state.storageStatus.favorite_count || 0), t("metric.favoritesHint"))}
      ${metricCard(t("metric.totalBytes"), formatBytes(state.storageStatus.total_bytes || 0), t("metric.totalBytesHint"))}
      ${metricCard(t("metric.fileBytes"), formatBytes(state.storageStatus.file_bytes || 0), t("metric.fileBytesHint"))}
    </div>
  `;
}

function renderCleanupStatusPanel(errorMessage = "") {
  if (errorMessage) {
    nodes.cleanupStatusPanel.innerHTML = renderErrorState(errorMessage);
    return;
  }
  if (!state.cleanupStatus) {
    nodes.cleanupStatusPanel.innerHTML = renderEmptyState(t("empty.saveTokenCleanup"));
    return;
  }

  nodes.cleanupStatusPanel.innerHTML = `
    <div class="meta-list">
      ${metaLine(t("label.lastRun"), formatDateTime(state.cleanupStatus.last_run_at))}
      ${metaLine(t("label.reason"), escapeHTML(state.cleanupStatus.last_run_reason || "-"))}
      ${metaLine(t("label.deletedExpired"), String(state.cleanupStatus.deleted_expired || 0))}
      ${metaLine(t("label.deletedMaxItems"), String(state.cleanupStatus.deleted_max_items || 0))}
      ${metaLine(t("label.deletedStorage"), String(state.cleanupStatus.deleted_storage || 0))}
      ${metaLine(t("label.deletedFiles"), String(state.cleanupStatus.deleted_files || 0))}
      ${metaLine(t("label.lastError"), escapeHTML(state.cleanupStatus.last_error || t("common.none")))}
    </div>
  `;
}

function renderWebDAVStatusPanel(errorMessage = "") {
  if (errorMessage) {
    nodes.webdavStatusPanel.innerHTML = renderErrorState(errorMessage);
    return;
  }
  if (!state.webdavStatus) {
    nodes.webdavStatusPanel.innerHTML = renderEmptyState(t("empty.saveTokenWebDAV"));
    return;
  }

  nodes.webdavStatusPanel.innerHTML = `
    <div class="meta-list">
      ${metaLine(t("label.lastSync"), formatDateTime(state.webdavStatus.last_sync_at))}
      ${metaLine(t("label.lastSuccess"), formatDateTime(state.webdavStatus.last_success_at))}
      ${metaLine(t("label.lastTested"), formatDateTime(state.webdavStatus.tested_at))}
      ${metaLine(t("label.lastMessage"), escapeHTML(state.webdavStatus.last_message || "-"))}
      ${metaLine(t("label.lastError"), escapeHTML(state.webdavStatus.last_error || state.webdavStatus.last_test_error || t("common.none")))}
      ${metaLine(t("label.pushedItems"), String(state.webdavStatus.pushed_items || 0))}
      ${metaLine(t("label.pulledItems"), String(state.webdavStatus.pulled_items || 0))}
      ${metaLine(t("label.pushedFiles"), String(state.webdavStatus.pushed_files || 0))}
      ${metaLine(t("label.pulledFiles"), String(state.webdavStatus.pulled_files || 0))}
      ${metaLine(t("label.remoteItems"), String(state.webdavStatus.remote_item_count || 0))}
      ${metaLine(t("label.localItems"), String(state.webdavStatus.local_item_count || 0))}
      ${metaLine(t("label.conflictSkips"), String(state.webdavStatus.conflict_skips || 0))}
    </div>
  `;
}

function populateCleanupForm() {
  const cleanup = state.cleanup;
  nodes.cleanupTTLHours.value = cleanup?.ttl_hours ?? "";
  nodes.cleanupMaxItems.value = cleanup?.max_items ?? "";
  nodes.cleanupMaxSizeMB.value = cleanup?.max_total_size_mb ?? "";
  nodes.cleanupIntervalMinutes.value = cleanup?.interval_minutes ?? "";
  nodes.cleanupEnabled.checked = Boolean(cleanup?.enabled);
}

function populateWebDAVForm() {
  const settings = state.webdavSettings;
  nodes.webdavEnabled.checked = Boolean(settings?.enabled);
  nodes.webdavURL.value = settings?.url || "";
  nodes.webdavUsername.value = settings?.username || "";
  nodes.webdavPassword.value = "";
  nodes.webdavBasePath.value = settings?.base_path || "";
}

function bindItemCardActions(container, items) {
  container.querySelectorAll("[data-action='copy']").forEach((button) => {
    button.addEventListener("click", () => {
      const item = findItem(items, button);
      if (item) {
        handleCopy(item);
      }
    });
  });

  container.querySelectorAll("[data-action='favorite']").forEach((button) => {
    button.addEventListener("click", async () => {
      const item = findItem(items, button);
      if (item) {
        await toggleFavorite(item);
      }
    });
  });

  container.querySelectorAll("[data-action='delete']").forEach((button) => {
    button.addEventListener("click", async () => {
      const item = findItem(items, button);
      if (item) {
        await deleteItem(item);
      }
    });
  });

  container.querySelectorAll("[data-action='detail']").forEach((button) => {
    button.addEventListener("click", () => {
      const item = findItem(items, button);
      if (item) {
        openItemDetails(item);
      }
    });
  });
}

function findItem(items, button) {
  const id = Number(button.closest("[data-item-id]")?.dataset.itemId);
  return items.find((item) => item.id === id);
}

function filteredHistory() {
  return state.history.filter((item) => {
    if (state.historyType !== "all" && item.type !== state.historyType) {
      return false;
    }
    if (!state.historySearch) {
      return true;
    }
    const haystack = [item.text, item.url, item.filename, item.source_device_name, item.mime_type, item.category, item.type]
      .filter(Boolean)
      .join(" ")
      .toLowerCase();
    return haystack.includes(state.historySearch);
  });
}

function renderItemCard(item, options = {}) {
  const content = describeItemContent(item);
  const preview = renderPreview(item, content);
  const actions = renderActions(item, options.latest);
  const categoryBadge = item.category && item.category !== item.type ? `<span class="badge badge-muted">${escapeHTML(item.category)}</span>` : "";
  return `
    <article class="item-card" data-item-id="${item.id}">
      <div class="item-header">
        <div class="detail-stack">
          <div class="item-badges">
            ${typeBadge(item.type)}
            ${item.is_favorite ? `<span class="badge badge-warning">${escapeHTML(t("common.favorite"))}</span>` : ""}
            ${categoryBadge}
          </div>
          ${content.title ? `<h3 class="item-title">${escapeHTML(content.title)}</h3>` : ""}
        </div>
        <div class="meta-row">
          <span>${formatDateTime(item.created_at)}</span>
        </div>
      </div>
      ${preview}
      <div class="meta-row">
        <span>${t("common.source")}: ${escapeHTML(item.source_device_name || t("common.unknown"))}</span>
        ${item.mime_type ? `<span>${t("label.mimeType")}: ${escapeHTML(item.mime_type)}</span>` : ""}
        ${item.size_bytes ? `<span>${t("common.size")}: ${formatBytes(item.size_bytes)}</span>` : ""}
      </div>
      <div class="item-actions">${actions}</div>
    </article>
  `;
}

function describeItemContent(item) {
  if (item.type === "link") {
    const url = item.url || "";
    return { title: url, preview: "" };
  }
  if (item.type === "file" || item.type === "image") {
    return { title: item.filename || t("item.storedFile"), preview: "" };
  }
  const text = truncate(item.text || t("type.text"), 320);
  return { title: text, preview: "" };
}

function renderPreview(item, content) {
  if (item.type === "image" && item.preview_url) {
    return `
      <div class="media-preview">
        <img src="${item.preview_url}" alt="${escapeHTML(item.filename || t("item.unnamedFile"))}" loading="lazy" />
      </div>
      ${content.preview ? `<div class="item-preview">${content.preview}</div>` : ""}
    `;
  }
  if (!content.preview) {
    return "";
  }
  return `<div class="item-preview">${content.preview}</div>`;
}

function renderActions(item, latest = false) {
  const buttons = [];
  const copyable = item.type === "text" || item.type === "link";

  if (copyable) {
    buttons.push(`<button class="button button-primary" type="button" data-action="copy">${latest ? t("common.copyLatest") : t("common.copy")}</button>`);
  }
  if (item.type === "link" && item.url) {
    buttons.push(`<a class="button button-secondary" href="${item.url}" target="_blank" rel="noreferrer">${t("common.open")}</a>`);
  }
  if ((item.type === "image" || item.type === "file") && item.download_url) {
    buttons.push(`<a class="button button-secondary" href="${item.download_url}" target="_blank" rel="noreferrer">${t("common.open")}</a>`);
    buttons.push(`<a class="button button-secondary" href="${item.download_url}" download>${t("common.download")}</a>`);
  }
  buttons.push(`<button class="button button-secondary" type="button" data-action="favorite">${item.is_favorite ? t("common.unfavorite") : t("common.favorite")}</button>`);
  buttons.push(`<button class="button button-secondary" type="button" data-action="detail">${t("common.details")}</button>`);
  buttons.push(`<button class="button button-danger" type="button" data-action="delete">${t("common.delete")}</button>`);
  return buttons.join("");
}

function renderDetailBody(item) {
  return `
    <div class="detail-stack">
      ${item.type === "image" && item.preview_url ? `<div class="media-preview"><img src="${item.preview_url}" alt="${escapeHTML(item.filename || t("item.unnamedFile"))}" /></div>` : ""}
      <div class="detail-panel">
        <div class="meta-list">
          ${metaLine(t("label.type"), escapeHTML(typeLabel(item.type)))}
          ${metaLine(t("label.category"), escapeHTML(item.category || "-"))}
          ${metaLine(t("label.sourceDevice"), escapeHTML(item.source_device_name || "-"))}
          ${metaLine(t("label.created"), formatDateTime(item.created_at))}
          ${metaLine(t("label.updated"), formatDateTime(item.updated_at))}
          ${metaLine(t("label.filename"), escapeHTML(item.filename || "-"))}
          ${metaLine(t("label.mimeType"), escapeHTML(item.mime_type || "-"))}
          ${metaLine(t("label.fileSize"), item.size_bytes ? formatBytes(item.size_bytes) : "-")}
          ${metaLine(t("label.sha256"), escapeHTML(item.sha256 || "-"))}
        </div>
      </div>
      <pre class="detail-pre">${escapeHTML(item.type === "link" ? item.url : item.text || item.filename || "")}</pre>
      <div class="item-actions">
        ${(item.type === "text" || item.type === "link") ? `<button id="detail-copy-button" class="button button-primary" type="button">${t("common.copy")}</button>` : ""}
        <button id="detail-favorite-button" class="button button-secondary" type="button">${item.is_favorite ? t("common.unfavorite") : t("common.favorite")}</button>
        <button id="detail-delete-button" class="button button-danger" type="button">${t("common.delete")}</button>
      </div>
    </div>
  `;
}

function metricCard(label, value, hint) {
  return `
    <article class="metric-card">
      <span class="field-label">${label}</span>
      <div class="status-value">${value}</div>
      <p class="support-text">${hint}</p>
    </article>
  `;
}

function renderEmptyState(message) {
  return `<div class="empty-state"><span class="badge badge-muted">${t("common.empty")}</span><p>${escapeHTML(message)}</p></div>`;
}

function renderErrorState(message) {
  return `<div class="error-state"><span class="badge badge-danger">${t("common.error")}</span><p>${escapeHTML(message)}</p></div>`;
}

function metaLine(label, value) {
  return `<div class="meta-line"><span class="meta-label">${label}</span><span>${value || "-"}</span></div>`;
}

function typeBadge(type) {
  const className = type === "text" ? "badge badge-accent" : type === "link" ? "badge badge-accent" : type === "image" ? "badge badge-success" : "badge badge-muted";
  return `<span class="${className}">${escapeHTML(typeLabel(type))}</span>`;
}

function typeLabel(type) {
  return t(`type.${type}`) || type;
}

function detailTitle(item) {
  return `${t("detail.clipboardItem")} #${item.id}`;
}

function activeDeviceCount() {
  return state.devices.filter((device) => !device.revoked_at).length;
}

function countdownLabel(expiresAt) {
  if (!expiresAt) {
    return t("pairing.noExpiry");
  }
  const diff = new Date(expiresAt).getTime() - Date.now();
  if (diff <= 0) {
    return t("pairing.expired");
  }
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const remainder = seconds % 60;
  return t("pairing.expiresIn", { minutes, seconds: remainder });
}

function updateCountdowns() {
  if (!state.pairing) {
    return;
  }
  renderPairingPanel();
}

function isLocalBrowserOrigin() {
  const hostname = window.location.hostname;
  return hostname === "127.0.0.1" || hostname === "localhost" || hostname === "::1";
}

function formatDateTime(value) {
  if (!value) {
    return "-";
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return new Intl.DateTimeFormat(state.language, {
    dateStyle: "medium",
    timeStyle: "short"
  }).format(date);
}

function formatBytes(bytes) {
  const numeric = Number(bytes);
  if (!numeric) {
    return "0 B";
  }
  const units = ["B", "KB", "MB", "GB", "TB"];
  const exponent = Math.min(Math.floor(Math.log(numeric) / Math.log(1024)), units.length - 1);
  const value = numeric / 1024 ** exponent;
  return `${value >= 10 || exponent === 0 ? value.toFixed(0) : value.toFixed(1)} ${units[exponent]}`;
}

function truncate(value, maxLength) {
  return value.length > maxLength ? `${value.slice(0, maxLength - 1)}…` : value;
}

function isNotFound(error) {
  return error?.status === 404;
}

async function apiFetch(url, options = {}) {
  const { auth = true, method = "GET", body } = options;
  const headers = new Headers();
  if (body) {
    headers.set("Content-Type", "application/json");
  }
  if (auth) {
    if (!state.token) {
      throw new Error(t("settings.session.signedOut"));
    }
    headers.set("Authorization", `Bearer ${state.token}`);
  }

  const response = await fetch(url, { method, headers, body });
  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    const error = new Error(payload?.error?.message || `Request failed with status ${response.status}`);
    error.status = response.status;
    throw error;
  }
  return payload.data;
}

function toast(message, tone = "success") {
  const element = document.createElement("div");
  element.className = `toast ${tone}`;
  element.textContent = message;
  nodes.toastRegion.appendChild(element);
  window.setTimeout(() => {
    element.remove();
  }, 3200);
}

function escapeHTML(value) {
  return String(value || "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}
