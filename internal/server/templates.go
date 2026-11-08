package server

// browseTmpl is the main file browser HTML page
const browseTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ShareDrop 📦</title>
<style>
  :root {
    --bg: #0f0f13;
    --surface: #1a1a24;
    --surface2: #22223a;
    --border: #2e2e4a;
    --accent: #7c6af7;
    --accent2: #a78bfa;
    --text: #e2e0ff;
    --text2: #9994cc;
    --green: #4ade80;
    --red: #f87171;
    --orange: #fb923c;
    --radius: 10px;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    background: var(--bg);
    color: var(--text);
    font-family: 'Segoe UI', system-ui, -apple-system, sans-serif;
    min-height: 100vh;
  }
  header {
    background: var(--surface);
    border-bottom: 1px solid var(--border);
    padding: 14px 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    position: sticky; top: 0; z-index: 100;
  }
  .logo {
    font-size: 1.4rem;
    font-weight: 700;
    color: var(--accent2);
    display: flex;
    align-items: center;
    gap: 10px;
    letter-spacing: -0.5px;
  }
  .logo span { font-size: 1.6rem; }
  .header-actions { display: flex; gap: 10px; align-items: center; }
  .btn {
    padding: 7px 16px;
    border-radius: var(--radius);
    border: none;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 600;
    transition: all 0.15s;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }
  .btn-primary { background: var(--accent); color: white; }
  .btn-primary:hover { background: var(--accent2); }
  .btn-ghost { background: var(--surface2); color: var(--text); border: 1px solid var(--border); }
  .btn-ghost:hover { background: var(--border); }
  .btn-danger { background: #7f1d1d; color: var(--red); border: 1px solid var(--red); }
  .btn-danger:hover { background: var(--red); color: white; }
  .container { max-width: 1100px; margin: 0 auto; padding: 24px; }
  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 20px;
    font-size: 0.9rem;
    flex-wrap: wrap;
  }
  .breadcrumb a { color: var(--accent2); text-decoration: none; }
  .breadcrumb a:hover { text-decoration: underline; }
  .breadcrumb .sep { color: var(--text2); }
  .breadcrumb .current { color: var(--text2); }
  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
    flex-wrap: wrap;
    gap: 10px;
  }
  .search-box {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 8px 14px;
    color: var(--text);
    font-size: 0.9rem;
    width: 260px;
    outline: none;
    transition: border 0.15s;
  }
  .search-box:focus { border-color: var(--accent); }
  .file-table {
    background: var(--surface);
    border-radius: var(--radius);
    border: 1px solid var(--border);
    overflow: hidden;
  }
  table { width: 100%; border-collapse: collapse; }
  thead tr { background: var(--surface2); }
  th {
    padding: 12px 16px;
    text-align: left;
    font-size: 0.8rem;
    color: var(--text2);
    text-transform: uppercase;
    letter-spacing: 0.5px;
    border-bottom: 1px solid var(--border);
  }
  td { padding: 10px 16px; border-bottom: 1px solid var(--border); }
  tr:last-child td { border-bottom: none; }
  tr:hover td { background: var(--surface2); }
  .file-name { display: flex; align-items: center; gap: 10px; }
  .file-name a { color: var(--text); text-decoration: none; font-weight: 500; }
  .file-name a:hover { color: var(--accent2); }
  .file-icon { font-size: 1.2rem; min-width: 24px; text-align: center; }
  .file-size { color: var(--text2); font-size: 0.85rem; white-space: nowrap; }
  .file-date { color: var(--text2); font-size: 0.85rem; white-space: nowrap; }
  .file-actions { display: flex; gap: 6px; justify-content: flex-end; }
  .action-link {
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.78rem;
    text-decoration: none;
    font-weight: 600;
    transition: all 0.1s;
  }
  .dl { background: #1e3a5f; color: #60a5fa; }
  .dl:hover { background: #2563eb; color: white; }
  .zip-link { background: #1e3a5f; color: #818cf8; }
  .zip-link:hover { background: #4338ca; color: white; }
  .del { background: #3b1111; color: var(--red); }
  .del:hover { background: var(--red); color: white; }
  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: var(--text2);
  }
  .empty-state .icon { font-size: 4rem; margin-bottom: 16px; }
  .upload-zone {
    background: var(--surface);
    border: 2px dashed var(--border);
    border-radius: var(--radius);
    padding: 28px;
    margin-bottom: 20px;
    transition: all 0.2s;
  }
  .upload-zone.dragover { border-color: var(--accent); background: rgba(124,106,247,0.08); }
  .upload-zone-inner {
    display: flex;
    align-items: center;
    gap: 16px;
    flex-wrap: wrap;
  }
  .upload-zone label {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    color: var(--text2);
    font-size: 0.9rem;
  }
  .upload-zone input[type=file] { display: none; }
  .upload-zone .file-chosen { color: var(--accent2); font-size: 0.85rem; flex: 1; min-width: 100px; }
  .badge {
    display: inline-block;
    padding: 2px 8px;
    border-radius: 99px;
    font-size: 0.72rem;
    font-weight: 700;
    text-transform: uppercase;
  }
  .badge-dir { background: rgba(124,106,247,0.15); color: var(--accent2); }
  .footer {
    text-align: center;
    padding: 20px;
    color: var(--text2);
    font-size: 0.82rem;
    border-top: 1px solid var(--border);
    margin-top: 40px;
  }
  .stats { color: var(--text2); font-size: 0.85rem; }
  @media (max-width: 600px) {
    .file-date, .file-size, th:nth-child(2), th:nth-child(3), td:nth-child(2), td:nth-child(3) { display: none; }
    .container { padding: 12px; }
    header { padding: 10px 14px; }
  }
  /* notification toast */
  .toast {
    position: fixed; bottom: 24px; right: 24px;
    background: var(--green); color: #052e16;
    padding: 10px 20px; border-radius: var(--radius);
    font-weight: 600; font-size: 0.9rem;
    opacity: 0; transform: translateY(20px);
    transition: all 0.3s; z-index: 999;
    pointer-events: none;
  }
  .toast.show { opacity: 1; transform: translateY(0); }
</style>
</head>
<body>
<header>
  <div class="logo"><span>📦</span> ShareDrop</div>
  <div class="header-actions">
    <span class="stats" id="fileCount"></span>
    {{if .PasswordSet}}
    <a href="/logout" class="btn btn-ghost">🔓 Logout</a>
    {{end}}
  </div>
</header>

<div class="container">
  <!-- Breadcrumb -->
  <div class="breadcrumb">
    <a href="/">🏠 {{.RootName}}</a>
    {{range .Breadcrumbs}}
    <span class="sep">/</span>
    <a href="/?path={{.Path}}">{{.Name}}</a>
    {{end}}
  </div>

  <!-- Toolbar -->
  <div class="toolbar">
    <input class="search-box" type="text" id="searchBox" placeholder="🔍  Search files..." oninput="filterFiles()">
    <div style="display:flex;gap:8px;flex-wrap:wrap">
      {{if .AllowUpload}}
      <button class="btn btn-primary" onclick="document.getElementById('uploadZone').style.display=document.getElementById('uploadZone').style.display==='none'?'block':'none'">
        ⬆️ Upload
      </button>
      {{end}}
    </div>
  </div>

  <!-- Upload zone -->
  {{if .AllowUpload}}
  <div id="uploadZone" class="upload-zone" style="display:none" 
       ondragover="event.preventDefault();this.classList.add('dragover')"
       ondragleave="this.classList.remove('dragover')"
       ondrop="handleDrop(event)">
    <form id="uploadForm" method="POST" action="/upload" enctype="multipart/form-data">
      <input type="hidden" name="path" value="{{.CurrentPath}}">
      <div class="upload-zone-inner">
        <label>
          📎 Choose files (or drag & drop)
          <input type="file" name="files" id="fileInput" multiple onchange="updateFileChosen()">
        </label>
        <span class="file-chosen" id="fileChosen">No files selected</span>
        <button type="submit" class="btn btn-primary">Upload Now</button>
      </div>
    </form>
  </div>
  {{end}}

  <!-- File table -->
  <div class="file-table">
    {{if .Files}}
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Size</th>
          <th>Modified</th>
          <th style="text-align:right">Actions</th>
        </tr>
      </thead>
      <tbody id="fileList">
        {{range .Files}}
        <tr class="file-row" data-name="{{.Name}}">
          <td>
            <div class="file-name">
              <span class="file-icon">{{.Icon}}</span>
              {{if .IsDir}}
              <a href="/?path={{.Path}}">{{.Name}}</a>
              <span class="badge badge-dir">folder</span>
              {{else}}
              <a href="/download?path={{.Path}}">{{.Name}}</a>
              {{end}}
            </div>
          </td>
          <td class="file-size">{{if .IsDir}}—{{else}}{{.Size}}{{end}}</td>
          <td class="file-date">{{.ModTime}}</td>
          <td>
            <div class="file-actions">
              {{if .IsDir}}
              <a href="/zip?path={{.Path}}" class="action-link zip-link" title="Download as ZIP">🗜️ ZIP</a>
              {{else}}
              <a href="/download?path={{.Path}}" class="action-link dl" title="Download">⬇️ Download</a>
              {{end}}
              {{if $.AllowDelete}}
              <form method="POST" action="/delete" style="display:inline" onsubmit="return confirm('Delete {{.Name}}?')">
                <input type="hidden" name="path" value="{{.Path}}">
                <button type="submit" class="action-link del">🗑️</button>
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
      <div class="icon">📂</div>
      <div>This folder is empty</div>
    </div>
    {{end}}
  </div>
</div>

<div class="footer">
  ShareDrop — local network file sharing · <a href="/api/files?path={{.CurrentPath}}" style="color:var(--accent2)">JSON API</a>
</div>

<div class="toast" id="toast">✓ Done!</div>

<script>
function filterFiles() {
  const q = document.getElementById('searchBox').value.toLowerCase();
  const rows = document.querySelectorAll('.file-row');
  let shown = 0;
  rows.forEach(r => {
    const name = r.dataset.name.toLowerCase();
    const match = name.includes(q);
    r.style.display = match ? '' : 'none';
    if (match) shown++;
  });
  document.getElementById('fileCount').textContent = shown + ' item' + (shown !== 1 ? 's' : '');
}

function updateFileChosen() {
  const inp = document.getElementById('fileInput');
  const span = document.getElementById('fileChosen');
  if (inp.files.length === 0) {
    span.textContent = 'No files selected';
  } else if (inp.files.length === 1) {
    span.textContent = inp.files[0].name;
  } else {
    span.textContent = inp.files.length + ' files selected';
  }
}

function handleDrop(e) {
  e.preventDefault();
  document.querySelector('#uploadZone').classList.remove('dragover');
  const dt = e.dataTransfer;
  const input = document.getElementById('fileInput');
  input.files = dt.files;
  updateFileChosen();
  // auto submit on drop
  document.getElementById('uploadForm').submit();
}

function showToast(msg) {
  const t = document.getElementById('toast');
  t.textContent = msg;
  t.classList.add('show');
  setTimeout(() => t.classList.remove('show'), 2500);
}

// initial count
window.addEventListener('DOMContentLoaded', () => {
  const rows = document.querySelectorAll('.file-row');
  document.getElementById('fileCount').textContent = rows.length + ' item' + (rows.length !== 1 ? 's' : '');
});
</script>
</body>
</html>`

// loginTmpl is the password gate page
const loginTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ShareDrop — Login</title>
<style>
  :root {
    --bg: #0f0f13;
    --surface: #1a1a24;
    --border: #2e2e4a;
    --accent: #7c6af7;
    --accent2: #a78bfa;
    --text: #e2e0ff;
    --text2: #9994cc;
    --red: #f87171;
    --radius: 12px;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    background: var(--bg);
    color: var(--text);
    font-family: 'Segoe UI', system-ui, sans-serif;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 40px;
    width: 100%;
    max-width: 380px;
    text-align: center;
  }
  .logo { font-size: 3rem; margin-bottom: 12px; }
  h1 { font-size: 1.5rem; color: var(--accent2); margin-bottom: 6px; }
  .subtitle { color: var(--text2); font-size: 0.9rem; margin-bottom: 28px; }
  .form-group { margin-bottom: 16px; text-align: left; }
  label { display: block; font-size: 0.85rem; color: var(--text2); margin-bottom: 6px; }
  input[type=password] {
    width: 100%;
    background: #0f0f13;
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 11px 14px;
    color: var(--text);
    font-size: 1rem;
    outline: none;
    transition: border 0.15s;
  }
  input[type=password]:focus { border-color: var(--accent); }
  .btn {
    width: 100%;
    padding: 11px;
    background: var(--accent);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.15s;
    margin-top: 4px;
  }
  .btn:hover { background: var(--accent2); }
  .error {
    background: rgba(248,113,113,0.1);
    border: 1px solid var(--red);
    color: var(--red);
    border-radius: 8px;
    padding: 10px 14px;
    font-size: 0.85rem;
    margin-bottom: 16px;
  }
</style>
</head>
<body>
<div class="card">
  <div class="logo">🔒</div>
  <h1>ShareDrop</h1>
  <p class="subtitle">This share is password protected</p>
  {{if .WrongPassword}}
  <div class="error">⚠️ Wrong password, try again.</div>
  {{end}}
  <form method="POST" action="/login">
    <div class="form-group">
      <label>Password</label>
      <input type="password" name="password" autofocus placeholder="Enter password...">
    </div>
    <button type="submit" class="btn">🚀 Unlock</button>
  </form>
</div>
</body>
</html>`
