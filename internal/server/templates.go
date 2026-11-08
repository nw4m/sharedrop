package server

const browseTmpl = `<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ShareDrop</title>
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

/* ---- dark theme (default) ---- */
:root {
  --bg:       #09090b;
  --bg1:      #111113;
  --bg2:      #18181b;
  --bg3:      #27272a;
  --border:   #3f3f46;
  --border2:  #52525b;
  --accent:   #6d28d9;
  --accent-h: #7c3aed;
  --accent-l: #ede9fe;
  --text:     #fafafa;
  --text2:    #a1a1aa;
  --text3:    #71717a;
  --green:    #22c55e;
  --red:      #ef4444;
  --blue:     #3b82f6;
  --r:        8px;
  --r-sm:     5px;
  --shadow:   0 4px 24px rgba(0,0,0,.5);
}

/* ---- light theme ---- */
[data-theme="light"] {
  --bg:       #f8f8f8;
  --bg1:      #ffffff;
  --bg2:      #f1f1f3;
  --bg3:      #e5e5ea;
  --border:   #d1d1d6;
  --border2:  #b0b0b8;
  --accent:   #6d28d9;
  --accent-h: #5b21b6;
  --accent-l: #4c1d95;
  --text:     #09090b;
  --text2:    #52525b;
  --text3:    #a1a1aa;
  --green:    #16a34a;
  --red:      #dc2626;
  --blue:     #2563eb;
  --shadow:   0 4px 24px rgba(0,0,0,.12);
}

html { height: 100%; }
body {
  background: var(--bg);
  color: var(--text);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  font-size: 14px;
  line-height: 1.5;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  transition: background .2s, color .2s;
}

::-webkit-scrollbar { width: 6px; }
::-webkit-scrollbar-track { background: var(--bg); }
::-webkit-scrollbar-thumb { background: var(--bg3); border-radius: 99px; }

/* ---- topbar ---- */
.topbar {
  height: 52px;
  background: var(--bg1);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  padding: 0 20px;
  gap: 12px;
  position: sticky;
  top: 0;
  z-index: 200;
}
.topbar-logo {
  font-size: 15px;
  font-weight: 700;
  color: var(--text);
  letter-spacing: -.3px;
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
  text-decoration: none;
}
.topbar-sep { width: 1px; height: 20px; background: var(--border); flex-shrink: 0; }
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  overflow: hidden;
  flex: 1;
  min-width: 0;
}
.breadcrumb-item a {
  color: var(--text2);
  text-decoration: none;
  padding: 2px 6px;
  border-radius: var(--r-sm);
  transition: color .12s, background .12s;
}
.breadcrumb-item a:hover { color: var(--text); background: var(--bg3); }
.breadcrumb-sep { color: var(--text3); font-size: 11px; }
.topbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
  flex-shrink: 0;
}

/* ---- theme toggle ---- */
.theme-toggle {
  width: 32px;
  height: 32px;
  border-radius: var(--r-sm);
  border: 1px solid var(--border);
  background: var(--bg2);
  color: var(--text2);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background .12s, color .12s;
  flex-shrink: 0;
}
.theme-toggle:hover { background: var(--bg3); color: var(--text); }
.icon-moon, .icon-sun { display: none; }
[data-theme="dark"]  .icon-moon { display: block; }
[data-theme="light"] .icon-sun  { display: block; }

/* ---- buttons ---- */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 0 14px;
  height: 32px;
  border-radius: var(--r-sm);
  border: none;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: background .12s, color .12s, border-color .12s;
  text-decoration: none;
  white-space: nowrap;
  flex-shrink: 0;
}
.btn-primary { background: var(--accent); color: #fff; }
.btn-primary:hover { background: var(--accent-h); }
.btn-ghost { background: transparent; color: var(--text2); border: 1px solid var(--border); }
.btn-ghost:hover { background: var(--bg3); color: var(--text); border-color: var(--border2); }
.btn-danger { background: transparent; color: var(--red); border: 1px solid transparent; }
.btn-danger:hover { background: rgba(239,68,68,.1); border-color: var(--red); }
.btn-sm { height: 26px; padding: 0 10px; font-size: 12px; }

/* ---- search ---- */
.search-wrap { position: relative; flex: 1; max-width: 360px; }
.search-icon {
  position: absolute; left: 10px; top: 50%;
  transform: translateY(-50%); color: var(--text3); pointer-events: none;
}
.search-input {
  width: 100%; height: 32px;
  background: var(--bg2); border: 1px solid var(--border);
  border-radius: var(--r-sm); padding: 0 10px 0 32px;
  color: var(--text); font-size: 13px; outline: none;
  transition: border-color .12s;
}
.search-input::placeholder { color: var(--text3); }
.search-input:focus { border-color: var(--accent); }

/* ---- layout ---- */
.layout { display: flex; flex: 1; min-height: 0; }
.sidebar {
  width: 220px; flex-shrink: 0;
  background: var(--bg1); border-right: 1px solid var(--border);
  padding: 16px 10px;
  display: flex; flex-direction: column; gap: 2px;
}
.sidebar-label {
  font-size: 11px; font-weight: 600; color: var(--text3);
  text-transform: uppercase; letter-spacing: .8px;
  padding: 0 8px; margin: 8px 0 4px;
}
.sidebar-item {
  display: flex; align-items: center; gap: 8px;
  padding: 7px 10px; border-radius: var(--r-sm);
  color: var(--text2); text-decoration: none; font-size: 13px;
  cursor: pointer; border: none; background: none;
  width: 100%; text-align: left;
  transition: background .1s, color .1s;
}
.sidebar-item:hover { background: var(--bg3); color: var(--text); }
.sidebar-item.active { background: rgba(109,40,217,.15); color: #c4b5fd; }
.sidebar-item svg { flex-shrink: 0; }

.main { flex: 1; display: flex; flex-direction: column; min-width: 0; overflow: auto; }

/* ---- toolbar ---- */
.toolbar {
  display: flex; align-items: center; gap: 10px;
  padding: 14px 20px; border-bottom: 1px solid var(--border);
  flex-wrap: wrap;
}
.toolbar-left  { display: flex; align-items: center; gap: 8px; flex: 1; min-width: 0; }
.toolbar-right { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }

/* ---- upload zone ---- */
.upload-wrap { padding: 12px 20px 0; }
.upload-zone {
  border: 1.5px dashed var(--border2); border-radius: var(--r);
  padding: 20px 24px; background: var(--bg2);
  transition: border-color .15s, background .15s;
}
.upload-zone.drag-over { border-color: var(--accent); background: rgba(109,40,217,.06); }
.upload-form-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.upload-file-label {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 0 14px; height: 32px; border-radius: var(--r-sm);
  border: 1px solid var(--border2); background: var(--bg3);
  color: var(--text2); font-size: 13px; cursor: pointer;
  transition: background .1s, color .1s;
}
.upload-file-label:hover { background: var(--border); color: var(--text); }
input[type=file] { display: none; }
.upload-chosen { font-size: 13px; color: var(--text2); flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.upload-hint { font-size: 12px; color: var(--text3); margin-top: 8px; }

/* ---- progress bar ---- */
.progress-wrap {
  margin-top: 12px;
  display: none;
}
.progress-bar-bg {
  height: 4px;
  background: var(--bg3);
  border-radius: 99px;
  overflow: hidden;
}
.progress-bar-fill {
  height: 100%;
  background: var(--accent);
  border-radius: 99px;
  width: 0%;
  transition: width .1s linear;
}
.progress-info {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--text3);
  margin-top: 5px;
}

/* ---- file table ---- */
.file-container { padding: 16px 20px; flex: 1; }
.file-table { width: 100%; border-collapse: collapse; }
.file-table thead th {
  text-align: left; padding: 8px 12px;
  font-size: 11px; font-weight: 600; color: var(--text3);
  text-transform: uppercase; letter-spacing: .5px;
  border-bottom: 1px solid var(--border);
  white-space: nowrap; cursor: pointer; user-select: none;
}
.file-table thead th:hover { color: var(--text2); }
.file-table thead th.sort-asc::after  { content: " ^"; }
.file-table thead th.sort-desc::after { content: " v"; }
.file-table tbody tr { border-bottom: 1px solid var(--border); transition: background .08s; }
.file-table tbody tr:last-child { border-bottom: none; }
.file-table tbody tr:hover { background: var(--bg2); }
.file-table tbody tr:hover .row-actions { opacity: 1; }
.file-table td { padding: 9px 12px; vertical-align: middle; }
.cell-name { display: flex; align-items: center; gap: 10px; min-width: 0; }
.file-type-badge {
  display: inline-flex; align-items: center; justify-content: center;
  width: 32px; height: 24px; border-radius: 4px;
  font-size: 10px; font-weight: 700;
  font-family: 'SF Mono', 'Fira Code', monospace;
  letter-spacing: .2px; flex-shrink: 0; text-transform: uppercase;
}
.type-dir  { background: rgba(109,40,217,.18); color: #c4b5fd; }
.type-img  { background: rgba(34,197,94,.15);  color: #86efac; }
.type-vid  { background: rgba(239,68,68,.15);  color: #fca5a5; }
.type-aud  { background: rgba(249,115,22,.15); color: #fdba74; }
.type-doc  { background: rgba(59,130,246,.15); color: #93c5fd; }
.type-zip  { background: rgba(234,179,8,.15);  color: #fde047; }
.type-src  { background: rgba(20,184,166,.15); color: #5eead4; }
.type-cfg  { background: rgba(156,163,175,.15);color: #d1d5db; }
.type-file { background: var(--bg3);           color: var(--text3); }
.file-name-link {
  color: var(--text); text-decoration: none;
  font-size: 13px; font-weight: 450;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
  transition: color .1s;
}
.file-name-link:hover { color: #c4b5fd; }
.cell-size, .cell-date { color: var(--text3); font-size: 12px; white-space: nowrap; font-variant-numeric: tabular-nums; }
.row-actions { display: flex; align-items: center; gap: 4px; justify-content: flex-end; opacity: 0; transition: opacity .1s; }

/* ---- empty state ---- */
.empty-state {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; padding: 80px 20px; gap: 12px; color: var(--text3);
}
.empty-icon {
  width: 48px; height: 48px; border-radius: 12px;
  background: var(--bg2); border: 1px solid var(--border);
  display: flex; align-items: center; justify-content: center;
}
.empty-state p { font-size: 14px; color: var(--text2); }
.empty-state span { font-size: 12px; }

/* ---- status bar ---- */
.statusbar {
  height: 28px; background: var(--bg1); border-top: 1px solid var(--border);
  display: flex; align-items: center; padding: 0 20px; gap: 16px;
  font-size: 12px; color: var(--text3); flex-shrink: 0;
}
.statusbar-item { display: flex; align-items: center; gap: 5px; }
.statusbar-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--green); }
.statusbar-sep { color: var(--border2); }

/* ---- toast ---- */
.toast {
  position: fixed; bottom: 20px; right: 20px;
  background: var(--bg2); border: 1px solid var(--border2);
  color: var(--text); padding: 10px 16px; border-radius: var(--r);
  font-size: 13px; box-shadow: var(--shadow);
  opacity: 0; transform: translateY(8px);
  transition: opacity .2s, transform .2s;
  pointer-events: none; z-index: 999;
  display: flex; align-items: center; gap: 8px;
}
.toast.show { opacity: 1; transform: translateY(0); }
.toast-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--green); flex-shrink: 0; }

/* ---- lock badge ---- */
.lock-badge {
  display: inline-flex; align-items: center; gap: 5px;
  font-size: 11px; font-weight: 600; color: var(--text3);
  padding: 3px 8px; border-radius: 99px; border: 1px solid var(--border);
}

/* ---- responsive ---- */
@media (max-width: 768px) { .sidebar { display: none; } .cell-date { display: none; } }
@media (max-width: 560px) { .cell-size { display: none; } .topbar { padding: 0 12px; } .file-container { padding: 10px 12px; } .toolbar { padding: 10px 12px; } }
</style>
</head>
<body>

<header class="topbar">
  <a href="/" class="topbar-logo">
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
      <polyline points="17 8 12 3 7 8"/>
      <line x1="12" y1="3" x2="12" y2="15"/>
    </svg>
    ShareDrop
  </a>
  <div class="topbar-sep"></div>
  <nav class="breadcrumb">
    <div class="breadcrumb-item"><a href="/">{{.RootName}}</a></div>
    {{range .Breadcrumbs}}
    <span class="breadcrumb-sep">›</span>
    <div class="breadcrumb-item"><a href="/?path={{.Path}}">{{.Name}}</a></div>
    {{end}}
  </nav>
  <div class="topbar-right">
    {{if .PasswordSet}}
    <span class="lock-badge">
      <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
      Protected
    </span>
    <a href="/logout" class="btn btn-ghost btn-sm">Logout</a>
    {{end}}
    <button class="theme-toggle" onclick="toggleTheme()" title="Toggle light/dark mode">
      <!-- moon icon (shown in dark mode) -->
      <svg class="icon-moon" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
      <!-- sun icon (shown in light mode) -->
      <svg class="icon-sun" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
    </button>
  </div>
</header>

<div class="layout">
  <aside class="sidebar">
    <span class="sidebar-label">Navigate</span>
    <a href="/" class="sidebar-item {{if eq .CurrentPath "."}}active{{end}}">
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
      Root
    </a>
    {{range .Breadcrumbs}}
    <a href="/?path={{.Path}}" class="sidebar-item">
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
      {{.Name}}
    </a>
    {{end}}

    <span class="sidebar-label" style="margin-top:16px">Actions</span>
    {{if .AllowUpload}}
    <button class="sidebar-item" onclick="toggleUpload()">
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
      Upload files
    </button>
    {{end}}
    <a href="/api/files?path={{.CurrentPath}}" class="sidebar-item" target="_blank">
      <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>
      JSON API
    </a>
  </aside>

  <div class="main">
    <div class="toolbar">
      <div class="toolbar-left">
        <div class="search-wrap">
          <svg class="search-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          <input class="search-input" id="searchInput" type="text" placeholder="Search in this folder..." oninput="filterFiles()" autocomplete="off">
        </div>
      </div>
      <div class="toolbar-right">
        {{if .AllowUpload}}
        <button class="btn btn-primary btn-sm" onclick="toggleUpload()">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/><path d="M20 21H4"/></svg>
          Upload
        </button>
        {{end}}
      </div>
    </div>

    {{if .AllowUpload}}
    <div id="uploadZone" class="upload-wrap" style="display:none">
      <div class="upload-zone" id="dropZone">
        <form id="uploadForm" method="POST" action="/upload" enctype="multipart/form-data">
          <input type="hidden" name="path" value="{{.CurrentPath}}">
          <div class="upload-form-row">
            <label class="upload-file-label" for="fileInput">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
              Choose files
            </label>
            <input type="file" id="fileInput" name="files" multiple>
            <span class="upload-chosen" id="uploadChosen">No files selected</span>
            <button type="button" id="uploadBtn" class="btn btn-primary btn-sm" onclick="startUpload()" disabled>Upload</button>
          </div>
          <p class="upload-hint">Or drag and drop files here</p>
          <!-- progress bar -->
          <div class="progress-wrap" id="progressWrap">
            <div class="progress-bar-bg">
              <div class="progress-bar-fill" id="progressFill"></div>
            </div>
            <div class="progress-info">
              <span id="progressText">Uploading...</span>
              <span id="progressPct">0%</span>
            </div>
          </div>
        </form>
      </div>
    </div>
    {{end}}

    <div class="file-container">
      {{if .Files}}
      <table class="file-table" id="fileTable">
        <thead>
          <tr>
            <th onclick="sortBy('name')" id="th-name">Name</th>
            <th onclick="sortBy('size')" id="th-size" style="width:90px">Size</th>
            <th onclick="sortBy('date')" id="th-date" style="width:140px">Modified</th>
            <th style="width:120px;cursor:default;text-align:right">Actions</th>
          </tr>
        </thead>
        <tbody id="fileList">
          {{range .Files}}
          <tr class="file-row" data-name="{{.Name}}" data-size="{{.RawSize}}" data-date="{{.ModTime}}" data-isdir="{{.IsDir}}">
            <td>
              <div class="cell-name">
                <span class="file-type-badge {{typeClass .Icon}}">{{typeLabel .Icon}}</span>
                {{if .IsDir}}
                <a class="file-name-link" href="/?path={{.Path}}">{{.Name}}</a>
                {{else}}
                <a class="file-name-link" href="/download?path={{.Path}}">{{.Name}}</a>
                {{end}}
              </div>
            </td>
            <td class="cell-size">{{if .IsDir}}--{{else}}{{.Size}}{{end}}</td>
            <td class="cell-date">{{.ModTime}}</td>
            <td>
              <div class="row-actions">
                {{if .IsDir}}
                <a href="/zip?path={{.Path}}" class="btn btn-ghost btn-sm">ZIP</a>
                {{else}}
                <a href="/download?path={{.Path}}" class="btn btn-ghost btn-sm">Download</a>
                {{end}}
                {{if $.AllowDelete}}
                <form method="POST" action="/delete" style="display:contents" onsubmit="return confirm('Delete &quot;{{.Name}}&quot;?')">
                  <input type="hidden" name="path" value="{{.Path}}">
                  <button type="submit" class="btn btn-danger btn-sm">Delete</button>
                </form>
                {{end}}
              </div>
            </td>
          </tr>
          {{end}}
        </tbody>
      </table>
      {{else}}
      <div class="empty-state">
        <div class="empty-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" color="#52525b"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
        </div>
        <p>This folder is empty</p>
        {{if .AllowUpload}}<span>Drop files here or click Upload</span>{{end}}
      </div>
      {{end}}
    </div>

    <footer class="statusbar">
      <div class="statusbar-item"><span class="statusbar-dot"></span>Running</div>
      <span class="statusbar-sep">|</span>
      <span id="statusCount" class="statusbar-item"></span>
      <span class="statusbar-sep">|</span>
      <span class="statusbar-item">{{.RootName}}</span>
    </footer>
  </div>
</div>

<div class="toast" id="toast"><span class="toast-dot"></span><span id="toastMsg">Done!</span></div>

<script>
var totalFiles = {{len .Files}};

// ---- theme ----
(function() {
  var saved = localStorage.getItem('sd-theme') || 'dark';
  document.documentElement.setAttribute('data-theme', saved);
})();

function toggleTheme() {
  var current = document.documentElement.getAttribute('data-theme');
  var next = current === 'dark' ? 'light' : 'dark';
  document.documentElement.setAttribute('data-theme', next);
  localStorage.setItem('sd-theme', next);
}

// ---- sort ----
var sortCol = 'name', sortDir = 1;
function sortBy(col) {
  if (sortCol === col) sortDir *= -1; else { sortCol = col; sortDir = 1; }
  ['name','size','date'].forEach(function(c) {
    var th = document.getElementById('th-' + c);
    if (th) th.className = (c === col) ? (sortDir === 1 ? 'sort-asc' : 'sort-desc') : '';
  });
  var tbody = document.getElementById('fileList');
  var rows = Array.from(tbody.querySelectorAll('.file-row'));
  rows.sort(function(a, b) {
    var ad = a.dataset.isdir === 'true', bd = b.dataset.isdir === 'true';
    if (ad !== bd) return ad ? -1 : 1;
    var av, bv;
    if (col === 'size') { av = parseInt(a.dataset.size)||0; bv = parseInt(b.dataset.size)||0; return (av-bv)*sortDir; }
    av = col === 'date' ? a.dataset.date : a.dataset.name.toLowerCase();
    bv = col === 'date' ? b.dataset.date : b.dataset.name.toLowerCase();
    return av < bv ? -sortDir : av > bv ? sortDir : 0;
  });
  rows.forEach(function(r) { tbody.appendChild(r); });
}

// ---- search ----
function filterFiles() {
  var q = document.getElementById('searchInput').value.toLowerCase();
  var rows = document.querySelectorAll('.file-row');
  var shown = 0;
  rows.forEach(function(r) {
    var match = r.dataset.name.toLowerCase().includes(q);
    r.style.display = match ? '' : 'none';
    if (match) shown++;
  });
  document.getElementById('statusCount').textContent = shown + ' item' + (shown !== 1 ? 's' : '');
}

// ---- upload ----
function toggleUpload() {
  var z = document.getElementById('uploadZone');
  if (!z) return;
  z.style.display = (z.style.display === 'none' || z.style.display === '') ? 'block' : 'none';
}

function onFileChange() {
  var inp = document.getElementById('fileInput');
  var sp  = document.getElementById('uploadChosen');
  var btn = document.getElementById('uploadBtn');
  if (!inp || !inp.files.length) {
    if (sp) sp.textContent = 'No files selected';
    if (btn) btn.disabled = true;
    return;
  }
  sp.textContent = inp.files.length === 1 ? inp.files[0].name : inp.files.length + ' files selected';
  btn.disabled = false;
}

function startUpload() {
  var form = document.getElementById('uploadForm');
  var inp  = document.getElementById('fileInput');
  if (!inp.files.length) return;

  var fd = new FormData(form);

  // show progress UI
  var progressWrap = document.getElementById('progressWrap');
  var fill  = document.getElementById('progressFill');
  var pText = document.getElementById('progressText');
  var pPct  = document.getElementById('progressPct');
  var btn   = document.getElementById('uploadBtn');
  progressWrap.style.display = 'block';
  btn.disabled = true;

  // calculate total size for display
  var totalSize = 0;
  for (var i = 0; i < inp.files.length; i++) totalSize += inp.files[i].size;

  var xhr = new XMLHttpRequest();

  xhr.upload.addEventListener('progress', function(e) {
    if (!e.lengthComputable) return;
    var pct = Math.round((e.loaded / e.total) * 100);
    fill.style.width = pct + '%';
    pPct.textContent = pct + '%';
    pText.textContent = 'Uploading ' + formatBytes(e.loaded) + ' / ' + formatBytes(e.total);
  });

  xhr.addEventListener('load', function() {
    fill.style.width = '100%';
    pPct.textContent = '100%';
    pText.textContent = 'Upload complete!';
    setTimeout(function() { window.location.reload(); }, 600);
  });

  xhr.addEventListener('error', function() {
    pText.textContent = 'Upload failed. Please try again.';
    pPct.textContent = '';
    btn.disabled = false;
  });

  xhr.open('POST', form.action);
  xhr.send(fd);
}

function formatBytes(b) {
  if (b < 1024) return b + ' B';
  if (b < 1048576) return (b/1024).toFixed(1) + ' KB';
  if (b < 1073741824) return (b/1048576).toFixed(1) + ' MB';
  return (b/1073741824).toFixed(2) + ' GB';
}

// ---- toast ----
function showToast(msg) {
  var t = document.getElementById('toast');
  document.getElementById('toastMsg').textContent = msg;
  t.classList.add('show');
  setTimeout(function() { t.classList.remove('show'); }, 2500);
}

// ---- init ----
document.addEventListener('DOMContentLoaded', function() {
  var sc = document.getElementById('statusCount');
  if (sc) sc.textContent = totalFiles + ' item' + (totalFiles !== 1 ? 's' : '');

  var inp = document.getElementById('fileInput');
  if (inp) inp.addEventListener('change', onFileChange);

  var dropZone = document.getElementById('dropZone');
  if (dropZone) {
    dropZone.addEventListener('dragover', function(e) { e.preventDefault(); dropZone.classList.add('drag-over'); });
    dropZone.addEventListener('dragleave', function(e) { dropZone.classList.remove('drag-over'); });
    dropZone.addEventListener('drop', function(e) {
      e.preventDefault();
      dropZone.classList.remove('drag-over');
      var files = e.dataTransfer.files;
      if (!files || !files.length) return;
      var inp = document.getElementById('fileInput');
      // build a new FileList via DataTransfer
      var dt = new DataTransfer();
      for (var i = 0; i < files.length; i++) dt.items.add(files[i]);
      inp.files = dt.files;
      onFileChange();
      startUpload();
    });
  }
});
</script>
</body>
</html>`

const loginTmpl = `<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ShareDrop - Login</title>
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
:root {
  --bg: #09090b; --bg1: #111113; --bg2: #18181b; --bg3: #27272a;
  --border: #3f3f46; --accent: #6d28d9; --accent-h: #5b21b6;
  --text: #fafafa; --text2: #a1a1aa; --red: #ef4444; --r: 8px;
}
[data-theme="light"] {
  --bg: #f8f8f8; --bg1: #ffffff; --bg2: #f1f1f3; --bg3: #e5e5ea;
  --border: #d1d1d6; --accent: #6d28d9; --accent-h: #5b21b6;
  --text: #09090b; --text2: #52525b; --red: #dc2626;
}
body {
  background: var(--bg); color: var(--text);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 20px;
}
.card {
  background: var(--bg1); border: 1px solid var(--border); border-radius: 12px;
  padding: 36px 32px; width: 100%; max-width: 360px;
}
.card-logo {
  width: 40px; height: 40px; border-radius: var(--r);
  background: rgba(109,40,217,.2); border: 1px solid rgba(109,40,217,.3);
  display: flex; align-items: center; justify-content: center; margin-bottom: 16px;
}
h1 { font-size: 17px; font-weight: 600; margin-bottom: 4px; }
.subtitle { font-size: 13px; color: var(--text2); margin-bottom: 24px; }
.form-group { margin-bottom: 14px; }
label { display: block; font-size: 12px; font-weight: 500; color: var(--text2); margin-bottom: 6px; }
input[type=password] {
  width: 100%; background: var(--bg2); border: 1px solid var(--border);
  border-radius: var(--r); padding: 9px 12px; color: var(--text);
  font-size: 14px; outline: none; transition: border-color .12s;
}
input[type=password]:focus { border-color: var(--accent); }
.btn {
  width: 100%; padding: 10px; background: var(--accent); color: white;
  border: none; border-radius: var(--r); font-size: 14px; font-weight: 600;
  cursor: pointer; transition: background .12s; margin-top: 4px;
}
.btn:hover { background: var(--accent-h); }
.error {
  background: rgba(239,68,68,.08); border: 1px solid rgba(239,68,68,.3);
  color: var(--red); border-radius: var(--r); padding: 9px 12px;
  font-size: 13px; margin-bottom: 14px;
}
</style>
<script>
(function() {
  var t = localStorage.getItem('sd-theme') || 'dark';
  document.documentElement.setAttribute('data-theme', t);
})();
</script>
</head>
<body>
<div class="card">
  <div class="card-logo">
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#a78bfa" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
  </div>
  <h1>ShareDrop</h1>
  <p class="subtitle">Enter your password to access this share</p>
  {{if .WrongPassword}}
  <div class="error">Incorrect password. Please try again.</div>
  {{end}}
  <form method="POST" action="/login">
    <div class="form-group">
      <label for="password">Password</label>
      <input type="password" id="password" name="password" autofocus placeholder="Enter password...">
    </div>
    <button type="submit" class="btn">Continue</button>
  </form>
</div>
</body>
</html>`
