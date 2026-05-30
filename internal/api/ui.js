// =============================================
// NullCloud Console — App
// =============================================

// ----- State -----
let token = '';
let vpcs = [], subnets = [], instances = [], loadbalancers = [], buckets = [], databases = [], clusters = [];
let modalState = null;
let pendingDelete = null;

// ----- Init -----
(function init() {
  const savedToken = localStorage.getItem('nc_token');
  if (savedToken) document.getElementById('tok').value = savedToken;

  const savedTheme = localStorage.getItem('nc_theme');
  if (savedTheme) document.documentElement.setAttribute('data-theme', savedTheme);

  document.getElementById('tok').addEventListener('keydown', e => {
    if (e.key === 'Enter') loadAll();
  });

  document.addEventListener('keydown', e => {
    if (e.key === 'Escape') {
      pendingDelete = null;
      resetModalOk();
      closeModal();
    }
    if (e.key === 'Enter' && !document.getElementById('overlay').classList.contains('hidden')) {
      if (['INPUT', 'SELECT'].includes(document.activeElement?.tagName)) submitModal();
    }
  });

  document.getElementById('overlay').addEventListener('click', ev => {
    if (ev.target === ev.currentTarget) {
      pendingDelete = null;
      resetModalOk();
      closeModal();
    }
  });

  if (savedToken) loadAll();
})();

// ----- Theme -----
function toggleTheme() {
  const current = document.documentElement.getAttribute('data-theme');
  const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  const isDark = current === 'dark' || (!current && systemDark);
  const next = isDark ? 'light' : 'dark';
  document.documentElement.setAttribute('data-theme', next);
  localStorage.setItem('nc_theme', next);
}

// ----- Tab navigation -----
function switchTabByName(name) {
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  document.getElementById('tab-btn-' + name)?.classList.add('active');
  document.querySelectorAll('.tab-panel').forEach(p => p.classList.add('hidden'));
  document.getElementById('tab-' + name)?.classList.remove('hidden');
  document.querySelectorAll('.stat-card').forEach(c =>
    c.classList.toggle('active', c.dataset.tab === name));
}

// ----- Data loading -----
async function loadAll() {
  const v = document.getElementById('tok').value.trim();
  if (!v) { toast('Please enter an API token', 'error'); return; }
  token = v;
  localStorage.setItem('nc_token', token);

  document.getElementById('welcome').classList.add('hidden');
  document.getElementById('view').classList.remove('hidden');
  showSkeletons();

  try {
    const [a, b, c, d, e, f, g] = await Promise.all([
      api('/v1/vpcs'), api('/v1/subnets'), api('/v1/instances'),
      api('/v1/loadbalancers'), api('/v1/buckets'), api('/v1/databases'), api('/v1/clusters'),
    ]);
    vpcs          = (a.vpcs           || []).sort(byCreated);
    subnets       = (b.subnets        || []).sort(byCreated);
    instances     = (c.instances      || []).sort(byCreated);
    loadbalancers = (d.load_balancers || []).sort(byCreated);
    buckets       = (e.buckets        || []).sort(byCreated);
    databases     = (f.databases      || []).sort(byCreated);
    clusters      = (g.clusters       || []).sort(byCreated);

    updateCounts();
    renderVPCs(); renderSubnets(); renderInstances();
    renderLoadBalancers(); renderBuckets(); renderDatabases(); renderClusters();

    const pill = document.getElementById('connPill');
    pill.classList.remove('hidden');
    const t = token;
    document.getElementById('connToken').textContent = t.length > 22 ? t.slice(0, 22) + '…' : t;
  } catch (err) {
    toast('Failed to load: ' + err.message, 'error');
  }
}

const byCreated = (a, b) => new Date(a.created_at) - new Date(b.created_at);

async function refreshAll(btn) {
  if (!token) return;
  btn = btn?.closest ? btn.closest('button') : btn;
  btn?.classList.add('spinning');
  try { await loadAll(); } finally { btn?.classList.remove('spinning'); }
}

function showSkeletons() {
  const rows = n => Array.from({length: n}, (_, i) => `
    <div class="skeleton-row">
      <div style="display:flex;flex-direction:column;gap:6px;flex:1">
        <div class="skel" style="width:${110+i*28}px;height:13px"></div>
        <div class="skel" style="width:${80+i*18}px;height:11px"></div>
      </div>
      <div class="skel" style="width:68px;height:22px;border-radius:20px;margin-left:auto"></div>
      <div class="skel" style="width:80px;height:28px;border-radius:6px"></div>
    </div>`).join('');
  document.getElementById('vpcsBody').innerHTML          = `<div class="skel-table">${rows(4)}</div>`;
  document.getElementById('subnetsBody').innerHTML       = `<div class="skel-table">${rows(4)}</div>`;
  document.getElementById('instancesBody').innerHTML     = `<div class="skel-table">${rows(4)}</div>`;
  document.getElementById('loadbalancersBody').innerHTML = `<div class="skel-table">${rows(3)}</div>`;
  document.getElementById('bucketsBody').innerHTML       = `<div class="skel-table">${rows(3)}</div>`;
  document.getElementById('databasesBody').innerHTML     = `<div class="skel-table">${rows(3)}</div>`;
  document.getElementById('clustersBody').innerHTML      = `<div class="skel-table">${rows(3)}</div>`;
}

function updateCounts() {
  document.getElementById('cVpc').textContent = vpcs.length;
  document.getElementById('cSub').textContent = subnets.length;
  document.getElementById('cVsi').textContent = instances.length;
  document.getElementById('cLb').textContent  = loadbalancers.length;
  document.getElementById('cBkt').textContent = buckets.length;
  document.getElementById('cDb').textContent  = databases.length;
  document.getElementById('cK8s').textContent = clusters.length;
}

// ----- Rendering helpers -----
const esc = s => String(s ?? '').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
const fmt = s => new Date(s).toLocaleString(undefined, {dateStyle:'medium', timeStyle:'short'});
const badge = st => `<span class="badge badge-${esc(st)}">${esc(st)}</span>`;

const PENCIL = `<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>`;
const TRASH  = `<svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6"/><path d="M10 11v6M14 11v6M9 6V4a1 1 0 011-1h4a1 1 0 011 1v2"/></svg>`;

function emptyState(title, sub) {
  return `<div class="empty">
    <svg class="empty-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2" stroke-linecap="round">
      <path d="M3 15a4 4 0 004 4h9a5 5 0 10-4.9-6H7a4 4 0 00-4 2z"/>
    </svg>
    <strong>${esc(title)}</strong>
    <p>${esc(sub)}</p>
  </div>`;
}

// ----- Renderers -----
function renderVPCs() {
  const el = document.getElementById('vpcsBody');
  if (!vpcs.length) { el.innerHTML = emptyState('No VPCs yet', 'Create your first virtual private cloud.'); return; }
  el.innerHTML = `<table>
    <thead><tr><th>Name</th><th>Status</th><th>CRN</th><th>Created</th><th></th></tr></thead>
    <tbody>${vpcs.map(v => `<tr>
      <td>
        <div class="rname">${esc(v.name)}</div>
        <div class="rid">${esc(v.id)}</div>
      </td>
      <td>${badge(v.status)}</td>
      <td><span class="crn">${esc(v.crn)}</span></td>
      <td style="color:var(--text-3);white-space:nowrap">${fmt(v.created_at)}</td>
      <td>
        <div class="acts row-actions">
          <button class="btn-icon" title="Rename" onclick='openEdit("vpc",${JSON.stringify(v)})'>${PENCIL}</button>
          <button class="btn-icon danger" title="Delete" onclick="confirmDelete('vpcs','${esc(v.id)}','${esc(v.name)}')">${TRASH}</button>
        </div>
      </td>
    </tr>`).join('')}</tbody>
  </table>`;
}

function renderSubnets() {
  const el = document.getElementById('subnetsBody');
  if (!subnets.length) { el.innerHTML = emptyState('No subnets yet', 'Create a subnet inside one of your VPCs.'); return; }
  const vpcMap = Object.fromEntries(vpcs.map(v => [v.id, v.name]));
  el.innerHTML = `<table>
    <thead><tr><th>Name</th><th>Status</th><th>CIDR</th><th>VPC</th><th>Created</th><th></th></tr></thead>
    <tbody>${subnets.map(s => `<tr>
      <td>
        <div class="rname">${esc(s.name)}</div>
        <div class="rid">${esc(s.id)}</div>
      </td>
      <td>${badge(s.status)}</td>
      <td><code>${esc(s.cidr_block)}</code></td>
      <td>
        <div style="font-size:13px">${esc(vpcMap[s.vpc_id] || s.vpc_id)}</div>
        <div class="rid">${esc(s.vpc_id)}</div>
      </td>
      <td style="color:var(--text-3);white-space:nowrap">${fmt(s.created_at)}</td>
      <td>
        <div class="acts row-actions">
          <button class="btn-icon" title="Rename" onclick='openEdit("subnet",${JSON.stringify(s)})'>${PENCIL}</button>
          <button class="btn-icon danger" title="Delete" onclick="confirmDelete('subnets','${esc(s.id)}','${esc(s.name)}')">${TRASH}</button>
        </div>
      </td>
    </tr>`).join('')}</tbody>
  </table>`;
}

function renderInstances() {
  const el = document.getElementById('instancesBody');
  if (!instances.length) { el.innerHTML = emptyState('No instances yet', 'Launch a virtual server instance.'); return; }
  const vpcMap = Object.fromEntries(vpcs.map(v => [v.id, v.name]));
  el.innerHTML = `<table>
    <thead><tr><th>Name</th><th>Status</th><th>IP</th><th>Profile / Image</th><th>VPC</th><th>Created</th><th></th></tr></thead>
    <tbody>${instances.map(i => `<tr>
      <td>
        <div class="rname">${esc(i.name)}</div>
        <div class="rid">${esc(i.id)}</div>
      </td>
      <td>${badge(i.status)}</td>
      <td><code>${esc(i.primary_ip)}</code></td>
      <td>
        ${i.profile ? `<div style="font-size:13px">${esc(i.profile)}</div>` : ''}
        ${i.image   ? `<div class="rid">${esc(i.image)}</div>` : (!i.profile ? `<span style="color:var(--text-3)">—</span>` : '')}
      </td>
      <td>
        <div style="font-size:13px">${esc(vpcMap[i.vpc_id] || i.vpc_id)}</div>
        <div class="rid">${esc(i.vpc_id)}</div>
      </td>
      <td style="color:var(--text-3);white-space:nowrap">${fmt(i.created_at)}</td>
      <td>
        <div class="acts row-actions">
          ${i.status === 'running'
            ? `<button class="btn btn-ghost btn-sm" onclick="vsiAct('${esc(i.id)}','stop')">Stop</button>
               <button class="btn btn-ghost btn-sm" onclick="vsiAct('${esc(i.id)}','restart')">Restart</button>`
            : `<button class="btn btn-ghost btn-sm" onclick="vsiAct('${esc(i.id)}','start')">Start</button>`}
          <button class="btn-icon" title="Rename" onclick='openEdit("instance",${JSON.stringify(i)})'>${PENCIL}</button>
          <button class="btn-icon danger" title="Delete" onclick="confirmDelete('instances','${esc(i.id)}','${esc(i.name)}')">${TRASH}</button>
        </div>
      </td>
    </tr>`).join('')}</tbody>
  </table>`;
}

function renderLoadBalancers() {
  const el = document.getElementById('loadbalancersBody');
  if (!loadbalancers.length) { el.innerHTML = emptyState('No load balancers yet', 'Create a load balancer to distribute traffic.'); return; }
  el.innerHTML = `<table>
    <thead><tr><th>Name</th><th>Status</th><th>Protocol</th><th>Port</th><th>CRN</th><th>Created</th><th></th></tr></thead>
    <tbody>${loadbalancers.map(lb => `<tr>
      <td>
        <div class="rname">${esc(lb.name)}</div>
        <div class="rid">${esc(lb.id)}</div>
      </td>
      <td>${badge(lb.status)}</td>
      <td><code>${esc(lb.protocol)}</code></td>
      <td><code>${esc(lb.port)}</code></td>
      <td><span class="crn">${esc(lb.crn)}</span></td>
      <td style="color:var(--text-3);white-space:nowrap">${fmt(lb.created_at)}</td>
      <td>
        <div class="acts row-actions">
          <button class="btn-icon" title="Rename" onclick='openEdit("loadbalancer",${JSON.stringify(lb)})'>${PENCIL}</button>
          <button class="btn-icon danger" title="Delete" onclick="confirmDelete('loadbalancers','${esc(lb.id)}','${esc(lb.name)}')">${TRASH}</button>
        </div>
      </td>
    </tr>`).join('')}</tbody>
  </table>`;
}

function renderBuckets() {
  const el = document.getElementById('bucketsBody');
  if (!buckets.length) { el.innerHTML = emptyState('No buckets yet', 'Create an object storage bucket.'); return; }
  el.innerHTML = `<table>
    <thead><tr><th>Name</th><th>Status</th><th>Region</th><th>CRN</th><th>Created</th><th></th></tr></thead>
    <tbody>${buckets.map(b => `<tr>
      <td>
        <div class="rname">${esc(b.name)}</div>
        <div class="rid">${esc(b.id)}</div>
      </td>
      <td>${badge(b.status)}</td>
      <td><code>${esc(b.region)}</code></td>
      <td><span class="crn">${esc(b.crn)}</span></td>
      <td style="color:var(--text-3);white-space:nowrap">${fmt(b.created_at)}</td>
      <td>
        <div class="acts row-actions">
          <button class="btn-icon" title="Rename" onclick='openEdit("bucket",${JSON.stringify(b)})'>${PENCIL}</button>
          <button class="btn-icon danger" title="Delete" onclick="confirmDelete('buckets','${esc(b.id)}','${esc(b.name)}')">${TRASH}</button>
        </div>
      </td>
    </tr>`).join('')}</tbody>
  </table>`;
}

function renderDatabases() {
  const el = document.getElementById('databasesBody');
  if (!databases.length) { el.innerHTML = emptyState('No databases yet', 'Create a managed database instance.'); return; }
  el.innerHTML = `<table>
    <thead><tr><th>Name</th><th>Status</th><th>Engine</th><th>Version</th><th>Plan</th><th>CRN</th><th>Created</th><th></th></tr></thead>
    <tbody>${databases.map(db => `<tr>
      <td>
        <div class="rname">${esc(db.name)}</div>
        <div class="rid">${esc(db.id)}</div>
      </td>
      <td>${badge(db.status)}</td>
      <td><code>${esc(db.engine)}</code></td>
      <td><code>${esc(db.version)}</code></td>
      <td><code>${esc(db.plan)}</code></td>
      <td><span class="crn">${esc(db.crn)}</span></td>
      <td style="color:var(--text-3);white-space:nowrap">${fmt(db.created_at)}</td>
      <td>
        <div class="acts row-actions">
          <button class="btn-icon" title="Rename" onclick='openEdit("database",${JSON.stringify(db)})'>${PENCIL}</button>
          <button class="btn-icon danger" title="Delete" onclick="confirmDelete('databases','${esc(db.id)}','${esc(db.name)}')">${TRASH}</button>
        </div>
      </td>
    </tr>`).join('')}</tbody>
  </table>`;
}

function renderClusters() {
  const el = document.getElementById('clustersBody');
  if (!clusters.length) { el.innerHTML = emptyState('No clusters yet', 'Create a Kubernetes cluster.'); return; }
  el.innerHTML = `<table>
    <thead><tr><th>Name</th><th>Status</th><th>Version</th><th>Nodes</th><th>CRN</th><th>Created</th><th></th></tr></thead>
    <tbody>${clusters.map(cl => `<tr>
      <td>
        <div class="rname">${esc(cl.name)}</div>
        <div class="rid">${esc(cl.id)}</div>
      </td>
      <td>${badge(cl.status)}</td>
      <td><code>${esc(cl.version)}</code></td>
      <td>${esc(cl.node_count)}</td>
      <td><span class="crn">${esc(cl.crn)}</span></td>
      <td style="color:var(--text-3);white-space:nowrap">${fmt(cl.created_at)}</td>
      <td>
        <div class="acts row-actions">
          <button class="btn-icon" title="Rename" onclick='openEdit("cluster",${JSON.stringify(cl)})'>${PENCIL}</button>
          <button class="btn-icon danger" title="Delete" onclick="confirmDelete('clusters','${esc(cl.id)}','${esc(cl.name)}')">${TRASH}</button>
        </div>
      </td>
    </tr>`).join('')}</tbody>
  </table>`;
}

// ----- Modal -----
function openCreate(type) {
  pendingDelete = null;
  resetModalOk();
  modalState = { mode: 'create', type, resource: null };
  const labels = { vpc: 'VPC', subnet: 'Subnet', instance: 'Instance', loadbalancer: 'Load Balancer', bucket: 'Bucket', database: 'Database', cluster: 'Cluster' };
  document.getElementById('modalTitle').textContent = `Create ${labels[type]}`;
  document.getElementById('modalOk').textContent = 'Create';
  document.getElementById('modalBody').innerHTML = buildForm(type, null);
  showOverlay();
  document.getElementById('f-name')?.focus();
}

function openEdit(type, resource) {
  pendingDelete = null;
  resetModalOk();
  modalState = { mode: 'edit', type, resource };
  const labels = { vpc: 'VPC', subnet: 'Subnet', instance: 'Instance', loadbalancer: 'Load Balancer', bucket: 'Bucket', database: 'Database', cluster: 'Cluster' };
  document.getElementById('modalTitle').textContent = `Rename ${labels[type]}`;
  document.getElementById('modalOk').textContent = 'Save';
  document.getElementById('modalBody').innerHTML = buildForm(type, resource);
  showOverlay();
  const inp = document.getElementById('f-name');
  inp?.focus();
  inp?.select();
}

function buildForm(type, resource) {
  const val = v => v ? ` value="${esc(v)}"` : '';

  let html = `<div class="field">
    <label for="f-name">Name</label>
    <input type="text" id="f-name" placeholder="my-resource"${val(resource?.name)} autocomplete="off" spellcheck="false"/>
  </div>`;

  if (type === 'subnet' && !resource) {
    const opts = vpcs.map(v => `<option value="${esc(v.id)}">${esc(v.name)}</option>`).join('');
    html += `<div class="field">
      <label for="f-vpc">VPC</label>
      ${opts
        ? `<select id="f-vpc">${opts}</select>`
        : `<p class="field-hint" style="color:var(--err)">No VPCs available — create a VPC first.</p>`}
    </div>`;
  }

  if (type === 'instance' && !resource) {
    const subOpts = subnets.map(s => `<option value="${esc(s.id)}">${esc(s.name)}</option>`).join('');
    html += `<div class="field">
      <label for="f-subnet">Subnet</label>
      ${subOpts
        ? `<select id="f-subnet">${subOpts}</select>`
        : `<p class="field-hint" style="color:var(--err)">No subnets available — create a subnet first.</p>`}
    </div>
    <div class="field-row">
      <div class="field">
        <label for="f-profile">Profile <span class="field-hint">(optional)</span></label>
        <input type="text" id="f-profile" placeholder="cx2-2x4" autocomplete="off"/>
      </div>
      <div class="field">
        <label for="f-image">Image <span class="field-hint">(optional)</span></label>
        <input type="text" id="f-image" placeholder="ibm-ubuntu-22-04" autocomplete="off"/>
      </div>
    </div>`;
  }

  if (type === 'loadbalancer' && !resource) {
    html += `<div class="field-row">
      <div class="field">
        <label for="f-protocol">Protocol</label>
        <select id="f-protocol">
          <option value="tcp">tcp</option>
          <option value="http">http</option>
          <option value="https" selected>https</option>
        </select>
      </div>
      <div class="field">
        <label for="f-port">Port</label>
        <input type="number" id="f-port" placeholder="443" min="1" max="65535" value="443" autocomplete="off"/>
      </div>
    </div>`;
  }

  if (type === 'bucket' && !resource) {
    html += `<div class="field">
      <label for="f-region">Region <span class="field-hint">(optional, default us-east-1)</span></label>
      <select id="f-region">
        <option value="">us-east-1 (default)</option>
        <option value="us-east-1">us-east-1</option>
        <option value="us-west-2">us-west-2</option>
        <option value="eu-west-1">eu-west-1</option>
        <option value="ap-southeast-1">ap-southeast-1</option>
      </select>
    </div>`;
  }

  if (type === 'database' && !resource) {
    html += `<div class="field-row">
      <div class="field">
        <label for="f-engine">Engine</label>
        <select id="f-engine">
          <option value="postgres">postgres</option>
          <option value="mysql">mysql</option>
          <option value="mariadb">mariadb</option>
        </select>
      </div>
      <div class="field">
        <label for="f-version">Version</label>
        <input type="text" id="f-version" placeholder="15" autocomplete="off"/>
      </div>
    </div>
    <div class="field">
      <label for="f-plan">Plan</label>
      <select id="f-plan">
        <option value="small">small</option>
        <option value="medium" selected>medium</option>
        <option value="large">large</option>
      </select>
    </div>`;
  }

  if (type === 'cluster' && !resource) {
    html += `<div class="field-row">
      <div class="field">
        <label for="f-version">Kubernetes Version</label>
        <input type="text" id="f-version" placeholder="1.30" autocomplete="off"/>
      </div>
      <div class="field">
        <label for="f-nodes">Node Count</label>
        <input type="number" id="f-nodes" placeholder="3" min="1" value="3" autocomplete="off"/>
      </div>
    </div>`;
  }

  return html;
}

function showOverlay() {
  document.getElementById('overlay').classList.remove('hidden');
}

function closeModal() {
  modalState = null;
  document.getElementById('overlay').classList.add('hidden');
}

function resetModalOk() {
  const btn = document.getElementById('modalOk');
  btn.className = 'btn btn-primary';
  btn.textContent = 'Save';
  btn.disabled = false;
}

async function submitModal() {
  if (pendingDelete) {
    const { path, id } = pendingDelete;
    pendingDelete = null;
    resetModalOk();
    closeModal();
    try {
      await api(`/v1/${path}/${id}`, { method: 'DELETE' });
      const labels = { vpcs: 'VPC', subnets: 'Subnet', instances: 'Instance', loadbalancers: 'Load Balancer', buckets: 'Bucket', databases: 'Database', clusters: 'Cluster' };
      toast(`${labels[path]} deleted`, 'success');
      await loadAll();
    } catch (err) { toast('Delete failed: ' + err.message, 'error'); }
    return;
  }

  if (!modalState) return;
  const { mode, type, resource } = modalState;

  const nameEl = document.getElementById('f-name');
  const name = nameEl?.value.trim();
  if (!name) {
    nameEl?.classList.add('error');
    nameEl?.focus();
    return;
  }
  nameEl?.classList.remove('error');

  const pathMap = { vpc: 'vpcs', subnet: 'subnets', instance: 'instances', loadbalancer: 'loadbalancers', bucket: 'buckets', database: 'databases', cluster: 'clusters' };
  const path = pathMap[type];

  document.getElementById('modalOk').disabled = true;
  try {
    if (mode === 'edit') {
      const res = await api(`/v1/${path}/${resource.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name }),
      });
      toast('Renamed successfully', 'success');
      closeModal();
      if (type === 'vpc')          { const i = vpcs.findIndex(x => x.id === resource.id);          if (i !== -1) { vpcs[i] = res;          renderVPCs(); } }
      if (type === 'subnet')       { const i = subnets.findIndex(x => x.id === resource.id);       if (i !== -1) { subnets[i] = res;       renderSubnets(); } }
      if (type === 'instance')     { const i = instances.findIndex(x => x.id === resource.id);     if (i !== -1) { instances[i] = res;     renderInstances(); } }
      if (type === 'loadbalancer') { const i = loadbalancers.findIndex(x => x.id === resource.id); if (i !== -1) { loadbalancers[i] = res; renderLoadBalancers(); } }
      if (type === 'bucket')       { const i = buckets.findIndex(x => x.id === resource.id);       if (i !== -1) { buckets[i] = res;       renderBuckets(); } }
      if (type === 'database')     { const i = databases.findIndex(x => x.id === resource.id);     if (i !== -1) { databases[i] = res;     renderDatabases(); } }
      if (type === 'cluster')      { const i = clusters.findIndex(x => x.id === resource.id);      if (i !== -1) { clusters[i] = res;      renderClusters(); } }
    } else {
      let body = { name };
      if (type === 'subnet') {
        const vpcId = document.getElementById('f-vpc')?.value;
        if (!vpcId) { toast('Select a VPC', 'error'); return; }
        body.vpc = { id: vpcId };
      }
      if (type === 'instance') {
        const subnetId = document.getElementById('f-subnet')?.value;
        if (!subnetId) { toast('Select a subnet', 'error'); return; }
        body.subnet = { id: subnetId };
        const profile = document.getElementById('f-profile')?.value.trim();
        const image   = document.getElementById('f-image')?.value.trim();
        if (profile) body.profile = { name: profile };
        if (image)   body.image   = { id: image };
      }
      if (type === 'loadbalancer') {
        const protocol = document.getElementById('f-protocol')?.value;
        const port     = parseInt(document.getElementById('f-port')?.value, 10);
        if (!protocol) { toast('Select a protocol', 'error'); return; }
        if (!port || port < 1 || port > 65535) { toast('Enter a valid port (1-65535)', 'error'); return; }
        body.protocol = protocol;
        body.port = port;
      }
      if (type === 'bucket') {
        const region = document.getElementById('f-region')?.value;
        if (region) body.region = region;
      }
      if (type === 'database') {
        const engine  = document.getElementById('f-engine')?.value;
        const version = document.getElementById('f-version')?.value.trim();
        const plan    = document.getElementById('f-plan')?.value;
        if (!version) { toast('Enter a version', 'error'); document.getElementById('f-version')?.focus(); return; }
        body.engine  = engine;
        body.version = version;
        body.plan    = plan;
      }
      if (type === 'cluster') {
        const version   = document.getElementById('f-version')?.value.trim();
        const nodeCount = parseInt(document.getElementById('f-nodes')?.value, 10);
        if (!version) { toast('Enter a Kubernetes version', 'error'); document.getElementById('f-version')?.focus(); return; }
        if (!nodeCount || nodeCount < 1) { toast('Node count must be at least 1', 'error'); return; }
        body.version    = version;
        body.node_count = nodeCount;
      }
      const res = await api(`/v1/${path}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
      const labels = { vpc: 'VPC', subnet: 'Subnet', instance: 'Instance', loadbalancer: 'Load Balancer', bucket: 'Bucket', database: 'Database', cluster: 'Cluster' };
      toast(`${labels[type]} created`, 'success');
      closeModal();
      if (type === 'vpc')          { vpcs.push(res);          vpcs.sort(byCreated);          updateCounts(); renderVPCs(); }
      if (type === 'subnet')       { subnets.push(res);       subnets.sort(byCreated);       updateCounts(); renderSubnets(); }
      if (type === 'instance')     { instances.push(res);     instances.sort(byCreated);     updateCounts(); renderInstances(); }
      if (type === 'loadbalancer') { loadbalancers.push(res); loadbalancers.sort(byCreated); updateCounts(); renderLoadBalancers(); }
      if (type === 'bucket')       { buckets.push(res);       buckets.sort(byCreated);       updateCounts(); renderBuckets(); }
      if (type === 'database')     { databases.push(res);     databases.sort(byCreated);     updateCounts(); renderDatabases(); }
      if (type === 'cluster')      { clusters.push(res);      clusters.sort(byCreated);      updateCounts(); renderClusters(); }
    }
  } catch (err) {
    toast((mode === 'edit' ? 'Rename' : 'Create') + ' failed: ' + err.message, 'error');
  } finally {
    document.getElementById('modalOk').disabled = false;
  }
}

// ----- Delete -----
function confirmDelete(path, id, name) {
  pendingDelete = { path, id };
  modalState = null;
  const labels = { vpcs: 'VPC', subnets: 'Subnet', instances: 'Instance', loadbalancers: 'Load Balancer', buckets: 'Bucket', databases: 'Database', clusters: 'Cluster' };
  document.getElementById('modalTitle').textContent = `Delete ${labels[path]}`;
  document.getElementById('modalOk').textContent = 'Delete';
  document.getElementById('modalOk').className = 'btn btn-danger';
  document.getElementById('modalBody').innerHTML = `
    <p style="color:var(--text-2);line-height:1.6">
      Delete <strong style="color:var(--text)">${esc(name)}</strong>?
      <br><span class="rid">${esc(id)}</span>
    </p>
    <div style="font-size:13px;color:var(--err);background:var(--err-bg);border:1px solid rgba(220,38,38,.15);border-radius:var(--r-sm);padding:10px 12px;margin-top:4px">
      This action cannot be undone.
    </div>`;
  showOverlay();
}

// ----- VSI actions -----
async function vsiAct(id, type) {
  try {
    const res = await api(`/v1/instances/${id}/actions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type }),
    });
    toast(`Instance ${type}ed`, 'success');
    const idx = instances.findIndex(i => i.id === id);
    if (idx !== -1) { instances[idx] = res; renderInstances(); }
  } catch (err) { toast('Action failed: ' + err.message, 'error'); }
}

// ----- API helper -----
async function api(path, opts = {}) {
  opts.headers = { ...(opts.headers || {}), Authorization: `Bearer ${token}` };
  const r = await fetch(path, opts);
  if (r.status === 204) return {};
  const d = await r.json();
  if (!r.ok) throw new Error(d.errors?.[0]?.message || r.statusText);
  return d;
}

// ----- Toasts -----
function toast(msg, type = 'success') {
  const el = document.createElement('div');
  el.className = `toast toast-${type}`;
  el.textContent = msg;
  document.getElementById('toasts').appendChild(el);
  setTimeout(() => {
    el.style.transition = 'opacity .18s ease, transform .18s ease';
    el.style.opacity = '0';
    el.style.transform = 'translateX(110%)';
    setTimeout(() => el.remove(), 200);
  }, 3000);
}
