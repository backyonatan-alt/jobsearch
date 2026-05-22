// Pursuit — vanilla JS frontend. Two pages share this file:
//   /            login page  (form id="login")
//   /app         app shell   (form id="new-app", views, etc.)

const STATUSES = ["wishlist","applied","screen","interview","offer","rejected","withdrawn"];

async function api(path, opts={}) {
  const r = await fetch(path, {
    headers: { "Content-Type": "application/json", ...(opts.headers||{}) },
    ...opts,
  });
  if (r.status === 401) { window.location = "/"; throw new Error("unauthorized"); }
  if (r.status === 204) return null;
  const body = await r.json().catch(() => ({}));
  if (!r.ok) throw new Error(body.error || ("http " + r.status));
  return body;
}

// ---------- Login page ----------
const login = document.getElementById("login");
if (login) {
  const status = document.getElementById("status");
  const params = new URLSearchParams(window.location.search);
  if (params.get("err") === "invalid_link") {
    status.textContent = "That link is invalid or expired. Request a new one.";
    status.className = "status error";
  }
  login.addEventListener("submit", async (e) => {
    e.preventDefault();
    status.textContent = "Sending…";
    status.className = "status";
    const email = document.getElementById("email").value;
    try {
      await api("/api/auth/request", { method: "POST", body: JSON.stringify({ email }) });
      status.textContent = "Check your email for a link. (In dev, check the server logs.)";
      status.className = "status ok";
    } catch (err) {
      status.textContent = err.message;
      status.className = "status error";
    }
  });
}

// ---------- App shell ----------
const newApp = document.getElementById("new-app");
if (newApp) {
  const who = document.getElementById("who");
  const tbody = document.querySelector("#apps-table tbody");
  const kanban = document.querySelector(".kanban");

  async function refresh() {
    const me = await api("/api/me");
    who.textContent = me.email;
    const apps = await api("/api/applications");
    renderList(apps);
    renderKanban(apps);
  }

  function renderList(apps) {
    tbody.innerHTML = "";
    if (apps.length === 0) {
      const tr = document.createElement("tr");
      tr.innerHTML = `<td colspan="6" class="empty">No applications yet. Add your first one above.</td>`;
      tbody.appendChild(tr);
      return;
    }
    for (const a of apps) {
      const tr = document.createElement("tr");
      tr.innerHTML = `
        <td>${escape(a.company)}</td>
        <td>${escape(a.role)}</td>
        <td>
          <select data-id="${a.id}" class="status-select">
            ${STATUSES.map(s => `<option value="${s}"${s===a.status?" selected":""}>${s}</option>`).join("")}
          </select>
        </td>
        <td>${a.applied_at ? new Date(a.applied_at).toLocaleDateString() : "—"}</td>
        <td>${escape(a.cv_variant || "")}</td>
        <td><button data-id="${a.id}" class="link delete">delete</button></td>
      `;
      tbody.appendChild(tr);
    }
    tbody.querySelectorAll(".status-select").forEach(sel => {
      sel.addEventListener("change", async () => {
        await api(`/api/applications/${sel.dataset.id}`, {
          method: "PATCH",
          body: JSON.stringify({ status: sel.value }),
        });
        refresh();
      });
    });
    tbody.querySelectorAll(".delete").forEach(btn => {
      btn.addEventListener("click", async () => {
        if (!confirm("Delete this application?")) return;
        await api(`/api/applications/${btn.dataset.id}`, { method: "DELETE" });
        refresh();
      });
    });
  }

  function renderKanban(apps) {
    kanban.innerHTML = "";
    const byStatus = Object.fromEntries(STATUSES.map(s => [s, []]));
    for (const a of apps) (byStatus[a.status] ||= []).push(a);
    for (const s of STATUSES) {
      const col = document.createElement("div");
      col.className = "kcol";
      col.innerHTML = `<h3>${s} <span class="count">${byStatus[s].length}</span></h3>`;
      for (const a of byStatus[s]) {
        const card = document.createElement("div");
        card.className = "kcard";
        card.innerHTML = `<strong>${escape(a.company)}</strong><br/><span>${escape(a.role)}</span>`;
        col.appendChild(card);
      }
      kanban.appendChild(col);
    }
  }

  function escape(s) {
    return String(s ?? "").replace(/[&<>"']/g, c =>
      ({ "&":"&amp;", "<":"&lt;", ">":"&gt;", "\"":"&quot;", "'":"&#39;" }[c]));
  }

  newApp.addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = Object.fromEntries(new FormData(newApp).entries());
    // Strip empty optionals so the server doesn't store empty strings.
    for (const k of ["jd_url","cv_variant","source","location","salary_note","notes"]) {
      if (data[k] === "") delete data[k];
    }
    await api("/api/applications", { method: "POST", body: JSON.stringify(data) });
    newApp.reset();
    refresh();
  });

  document.querySelectorAll("nav [data-view]").forEach(btn => {
    btn.addEventListener("click", () => {
      document.querySelectorAll("nav [data-view]").forEach(b => b.classList.remove("active"));
      btn.classList.add("active");
      document.getElementById("list-view").classList.toggle("hidden", btn.dataset.view !== "list");
      document.getElementById("kanban-view").classList.toggle("hidden", btn.dataset.view !== "kanban");
    });
  });

  document.getElementById("logout").addEventListener("click", async () => {
    await api("/api/auth/logout", { method: "POST" });
    window.location = "/";
  });

  refresh();
}
