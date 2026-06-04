// =============================================
// NullCloud Console — App
// =============================================

// ----- State -----
let token = '';
let vpcs = [], subnets = [], instances = [], loadbalancers = [], buckets = [], databases = [], clusters = [];
let regions = [];
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
    const [a, b, c, d, e, f, g, h] = await Promise.all([
      api('/v1/vpcs'), api('/v1/subnets'), api('/v1/instances'),
      api('/v1/loadbalancers'), api('/v1/buckets'), api('/v1/databases'), api('/v1/clusters'),
      api('/v1/regions'),
    ]);
    regions       = (h.regions        || []);
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
    renderGraph();

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

// ----- Graph -----
function renderGraph() {
  const el = document.getElementById('graphCanvas');
  if (!el) return;

  const hasAny = vpcs.length || subnets.length || instances.length ||
                 loadbalancers.length || buckets.length || databases.length || clusters.length;
  if (!hasAny) {
    el.innerHTML = emptyState('No resources yet', 'Create some resources to see the infrastructure topology.');
    return;
  }

  // Layout constants
  const M    = 28;            // outer margin
  const NW   = 130, NH = 58; // node width/height
  const GN   = 12;            // gap between nodes
  const GS   = 16;            // gap between subnets
  const GZ   = 12;            // gap between zones
  const GV   = 24;            // gap between vpcs
  const GRG  = 36;            // gap between region boxes
  const GROW = 48;            // gap between rows (lb / regions / buckets)
  const SPX = 14, SPY = 34, SPB = 16; // subnet padding x / top / bottom
  const ZPX = 12, ZPY = 28, ZPB = 12; // zone padding x / top / bottom
  const VPX = 16, VPY = 44, VPB = 16; // vpc padding x / top / bottom
  const RPX = 20, RPY = 42, RPB = 16; // region padding x / top / bottom
  const MIN_SW = 160, MIN_ZW = 186, MIN_VW = 210, MIN_RW = 250;

  const pos = {}; // id -> {x, y, w, h, t, name, status}

  // --- Group resources by subnet ---
  const bySub = {};
  for (const s of subnets) bySub[s.id] = [];
  for (const vsi of instances) {
    if (bySub[vsi.subnet_id] !== undefined) bySub[vsi.subnet_id].push({...vsi, _t: 'instance'});
  }
  for (const db of databases) {
    for (const sid of (db.subnet_ids || [])) {
      if (bySub[sid] !== undefined) bySub[sid].push({...db, _t: 'database'});
    }
  }
  for (const cl of clusters) {
    for (const sid of (cl.subnet_ids || [])) {
      if (bySub[sid] !== undefined) bySub[sid].push({...cl, _t: 'cluster'});
    }
  }

  // --- Group subnets by zone within each VPC ---
  const byVpcZone = {}; // vpcId -> zoneName -> [subnets]
  for (const v of vpcs) byVpcZone[v.id] = {};
  for (const s of subnets) {
    if (!byVpcZone[s.vpc_id]) continue;
    const z = s.zone || 'unknown';
    if (!byVpcZone[s.vpc_id][z]) byVpcZone[s.vpc_id][z] = [];
    byVpcZone[s.vpc_id][z].push(s);
  }

  // --- Group VPCs by region ---
  const byRegion = {}; // regionName -> [vpcs]
  for (const v of vpcs) {
    const r = v.region || 'unknown';
    if (!byRegion[r]) byRegion[r] = [];
    byRegion[r].push(v);
  }
  const regionNames = Object.keys(byRegion).sort();

  // --- Compute sizes bottom-up ---
  const subSz = {};
  for (const s of subnets) {
    const n = (bySub[s.id] || []).length;
    const innerW = n > 0 ? n * NW + (n - 1) * GN : NW;
    subSz[s.id] = {
      w: Math.max(MIN_SW, innerW + 2 * SPX),
      h: SPY + (n > 0 ? NH : 0) + SPB,
    };
  }

  const zoneSz = {}; // vpcId+'|'+zone -> {w, h}
  for (const v of vpcs) {
    for (const [zname, zsubs] of Object.entries(byVpcZone[v.id] || {})) {
      const key = v.id + '|' + zname;
      const innerW = zsubs.length > 0
        ? zsubs.reduce((a, s) => a + subSz[s.id].w, 0) + (zsubs.length - 1) * GS
        : MIN_SW;
      const maxH = zsubs.length > 0 ? Math.max(...zsubs.map(s => subSz[s.id].h)) : SPY + SPB;
      zoneSz[key] = {
        w: Math.max(MIN_ZW, innerW + 2 * ZPX),
        h: ZPY + maxH + ZPB,
      };
    }
  }

  const vpcSz = {};
  for (const v of vpcs) {
    const zones = Object.entries(byVpcZone[v.id] || {}).sort(([a], [b]) => a.localeCompare(b));
    if (zones.length === 0) {
      vpcSz[v.id] = { w: MIN_VW, h: VPY + VPB };
    } else {
      const innerW = zones.reduce((a, [z]) => a + zoneSz[v.id + '|' + z].w, 0) + (zones.length - 1) * GZ;
      const maxH   = Math.max(...zones.map(([z]) => zoneSz[v.id + '|' + z].h));
      vpcSz[v.id] = {
        w: Math.max(MIN_VW, innerW + 2 * VPX),
        h: VPY + maxH + VPB,
      };
    }
  }

  const regionSz = {};
  for (const rname of regionNames) {
    const rvpcs = byRegion[rname] || [];
    if (rvpcs.length === 0) {
      regionSz[rname] = { w: MIN_RW, h: RPY + RPB };
    } else {
      const innerW = rvpcs.reduce((a, v) => a + vpcSz[v.id].w, 0) + (rvpcs.length - 1) * GV;
      const maxH   = Math.max(...rvpcs.map(v => vpcSz[v.id].h));
      regionSz[rname] = {
        w: Math.max(MIN_RW, innerW + 2 * RPX),
        h: RPY + maxH + RPB,
      };
    }
  }

  // --- Position LBs (top row) ---
  const hasLBs = loadbalancers.length > 0;
  const lbY = M + (hasLBs ? 16 : 0);
  let lbX = M;
  for (const lb of loadbalancers) {
    pos[lb.id] = {x: lbX, y: lbY, w: NW, h: NH, t: 'loadbalancer', name: lb.name, status: lb.status};
    lbX += NW + GN;
  }

  // --- Position regions (middle row) ---
  const regRowY = hasLBs ? lbY + NH + GROW : M;
  let regX = M;
  const regionR = {};
  const vpcR    = {};
  const zoneR   = {};
  const subR    = {};
  let maxRegH   = 0;

  for (const rname of regionNames) {
    const rsz = regionSz[rname];
    regionR[rname] = { x: regX, y: regRowY, ...rsz };
    maxRegH = Math.max(maxRegH, rsz.h);

    let vx = regX + RPX;
    const vy = regRowY + RPY;
    for (const v of (byRegion[rname] || [])) {
      vpcR[v.id] = { x: vx, y: vy, ...vpcSz[v.id] };

      const zones = Object.entries(byVpcZone[v.id] || {}).sort(([a], [b]) => a.localeCompare(b));
      let zx = vx + VPX;
      const zy = vy + VPY;
      for (const [zname, zsubs] of zones) {
        const key = v.id + '|' + zname;
        zoneR[key] = { x: zx, y: zy, ...zoneSz[key] };

        let sx = zx + ZPX;
        const sy = zy + ZPY;
        for (const s of zsubs) {
          subR[s.id] = { x: sx, y: sy, ...subSz[s.id] };
          sx += subSz[s.id].w + GS;
        }
        zx += zoneSz[key].w + GZ;
      }
      vx += vpcSz[v.id].w + GV;
    }
    regX += rsz.w + GRG;
  }

  // --- Position nodes within subnets ---
  for (const s of subnets) {
    const sr = subR[s.id];
    if (!sr) continue;
    const res = bySub[s.id] || [];
    const totalW = res.length * NW + (res.length - 1) * GN;
    let nx = sr.x + (sr.w - totalW) / 2;
    const ny = sr.y + SPY;
    for (const r of res) {
      pos[r.id] = {x: nx, y: ny, w: NW, h: NH, t: r._t, name: r.name, status: r.status};
      nx += NW + GN;
    }
  }

  // --- Position Buckets (bottom row) ---
  const hasBkts = buckets.length > 0;
  const bkY = regRowY + (maxRegH || NH + RPY + RPB) + GROW;
  let bkX = M;
  for (const b of buckets) {
    pos[b.id] = {x: bkX, y: bkY, w: NW, h: NH, t: 'bucket', name: b.name, status: b.status};
    bkX += NW + GN;
  }

  // --- SVG dimensions ---
  const allBoxes = [
    ...Object.values(regionR),
    ...Object.values(pos).map(p => ({x: p.x, y: p.y, w: p.w, h: p.h})),
  ];
  const svgW = allBoxes.length ? Math.max(...allBoxes.map(b => b.x + b.w)) + M : M * 6;
  const svgH = allBoxes.length ? Math.max(...allBoxes.map(b => b.y + b.h)) + M : M * 6;

  // --- Node type config ---
  const TC = {
    instance:     {bg: 'rgba(245,158,11,.13)',  stroke: '#fbbf24', label: 'VSI'},
    database:     {bg: 'rgba(234,88,12,.13)',   stroke: '#fb923c', label: 'DB'},
    cluster:      {bg: 'rgba(236,72,153,.13)',  stroke: '#f472b6', label: 'K8s'},
    loadbalancer: {bg: 'rgba(59,130,246,.13)',  stroke: '#60a5fa', label: 'LB'},
    bucket:       {bg: 'rgba(34,197,94,.13)',   stroke: '#4ade80', label: 'Bucket'},
  };

  const statusFill = st => {
    if (st === 'running' || st === 'available' || st === 'active') return '#34d399';
    if (st === 'stopped') return '#f87171';
    return 'var(--text-3)';
  };

  let out = '<defs>' +
    '<marker id="g-arr" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="5" markerHeight="5" orient="auto">' +
    '<path d="M0,0 L10,5 L0,10 Z" fill="#60a5fa" fill-opacity="0.55"/></marker>' +
    '</defs>';

  // Region containers (back-most layer)
  for (const rname of regionNames) {
    const r = regionR[rname];
    out += `<rect x="${r.x}" y="${r.y}" width="${r.w}" height="${r.h}" rx="14" ` +
      `style="fill:none;stroke:var(--text-3);stroke-opacity:0.25;stroke-width:1.5;stroke-dasharray:7,4"/>`;
    out += `<text x="${r.x+14}" y="${r.y+22}" ` +
      `style="font-size:10px;font-weight:700;letter-spacing:0.8px;fill:var(--text-3)">` +
      `REGION: ${esc(rname.toUpperCase())}</text>`;
  }

  // VPC containers
  for (const v of vpcs) {
    const r = vpcR[v.id];
    if (!r) continue;
    out += `<rect x="${r.x}" y="${r.y}" width="${r.w}" height="${r.h}" rx="10" ` +
      `style="fill:var(--brand-subtle);stroke:var(--brand);stroke-opacity:0.45;stroke-width:1.5"/>`;
    out += `<text x="${r.x+12}" y="${r.y+22}" ` +
      `style="font-size:13px;font-weight:700;fill:var(--brand-light)">${esc(v.name)}</text>`;
    out += `<text x="${r.x+12}" y="${r.y+34}" ` +
      `style="font-size:9px;font-family:monospace;fill:var(--text-3)">vpc · ${esc(v.id.slice(0,10))}…</text>`;
  }

  // Zone containers
  for (const [key, r] of Object.entries(zoneR)) {
    const zname = key.split('|')[1];
    out += `<rect x="${r.x}" y="${r.y}" width="${r.w}" height="${r.h}" rx="6" ` +
      `style="fill:rgba(100,116,139,.05);stroke:var(--border);stroke-width:1;stroke-dasharray:4,3"/>`;
    out += `<text x="${r.x+8}" y="${r.y+17}" ` +
      `style="font-size:9px;font-weight:600;letter-spacing:0.5px;fill:var(--text-3)">${esc(zname)}</text>`;
  }

  // Subnet containers
  for (const s of subnets) {
    const r = subR[s.id];
    if (!r) continue;
    out += `<rect x="${r.x}" y="${r.y}" width="${r.w}" height="${r.h}" rx="8" ` +
      `style="fill:var(--surface);stroke:var(--border);stroke-width:1"/>`;
    out += `<text x="${r.x+10}" y="${r.y+18}" ` +
      `style="font-size:11px;font-weight:600;fill:var(--text-2)">${esc(s.name)}</text>`;
    if (s.cidr_block) {
      out += `<text x="${r.x+10}" y="${r.y+29}" ` +
        `style="font-size:9px;font-family:monospace;fill:var(--text-3)">${esc(s.cidr_block)}</text>`;
    }
  }

  // Section labels
  if (hasLBs) {
    out += `<text x="${M}" y="${lbY-10}" style="font-size:10px;font-weight:600;fill:var(--text-3)">LOAD BALANCERS</text>`;
  }
  if (hasBkts) {
    out += `<text x="${M}" y="${bkY-10}" style="font-size:10px;font-weight:600;fill:var(--text-3)">OBJECT STORAGE</text>`;
  }

  // LB → target edges (drawn under nodes)
  for (const lb of loadbalancers) {
    const src = pos[lb.id];
    if (!src) continue;
    for (const tgt of (lb.targets || [])) {
      const dst = pos[tgt.id];
      if (!dst) continue;
      const x1 = src.x + NW / 2, y1 = src.y + NH;
      const x2 = dst.x + NW / 2, y2 = dst.y;
      const cy = (y1 + y2) / 2;
      out += `<path d="M${x1},${y1} C${x1},${cy} ${x2},${cy} ${x2},${y2}" ` +
        `style="fill:none;stroke:#60a5fa;stroke-width:1.5;stroke-dasharray:5,3;stroke-opacity:0.55" ` +
        `marker-end="url(#g-arr)"/>`;
    }
  }

  // Resource nodes (front-most layer)
  for (const [, n] of Object.entries(pos)) {
    const c = TC[n.t] || {bg: 'rgba(148,163,184,.1)', stroke: 'var(--border)', label: n.t};
    const label = n.name.length > 15 ? n.name.slice(0, 13) + '…' : n.name;
    out += `<g>` +
      `<rect x="${n.x}" y="${n.y}" width="${n.w}" height="${n.h}" rx="8" ` +
      `style="fill:${c.bg};stroke:${c.stroke};stroke-width:1.5"/>` +
      `<text x="${n.x + NW/2}" y="${n.y+17}" text-anchor="middle" ` +
      `style="font-size:10px;font-weight:700;fill:${c.stroke};letter-spacing:0.5px">${esc(c.label)}</text>` +
      `<text x="${n.x + NW/2}" y="${n.y+33}" text-anchor="middle" ` +
      `style="font-size:12px;font-weight:500;fill:var(--text)">${esc(label)}</text>` +
      `<circle cx="${n.x + NW/2 - 22}" cy="${n.y+47}" r="3" style="fill:${statusFill(n.status)}"/>` +
      `<text x="${n.x + NW/2 - 16}" y="${n.y+51}" ` +
      `style="font-size:10px;fill:var(--text-3)">${esc(n.status || '—')}</text>` +
      `</g>`;
  }

  el.innerHTML = `<svg viewBox="0 0 ${svgW} ${svgH}" width="${svgW}" height="${svgH}" ` +
    `xmlns="http://www.w3.org/2000/svg" style="display:block">${out}</svg>`;
}

// ----- Zone helpers -----
function zonesForVpc(vpcId) {
  const vpc = vpcs.find(v => v.id === vpcId);
  if (!vpc?.region) return [];
  const reg = regions.find(r => r.name === vpc.region);
  return (reg?.zones || []).map(z => z.name);
}

function refreshZoneSelect() {
  const vpcId = document.getElementById('f-vpc')?.value;
  const sel = document.getElementById('f-zone');
  if (!sel) return;
  const zones = zonesForVpc(vpcId);
  sel.innerHTML = zones.map(z => `<option value="${esc(z)}">${esc(z)}</option>`).join('');
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

  if (type === 'vpc' && !resource) {
    const opts = regions.map(r => `<option value="${esc(r.name)}">${esc(r.name)}</option>`).join('');
    html += `<div class="field">
      <label for="f-region">Region</label>
      ${opts
        ? `<select id="f-region">${opts}</select>`
        : `<p class="field-hint" style="color:var(--err)">No regions available.</p>`}
    </div>`;
  }

  if (type === 'subnet' && !resource) {
    const firstVpc = vpcs[0];
    const vpcOpts = vpcs.map(v => `<option value="${esc(v.id)}">${esc(v.name)}</option>`).join('');
    const initialZones = firstVpc ? zonesForVpc(firstVpc.id) : [];
    const zoneOpts = initialZones.map(z => `<option value="${esc(z)}">${esc(z)}</option>`).join('');
    html += `<div class="field">
      <label for="f-vpc">VPC</label>
      ${vpcOpts
        ? `<select id="f-vpc" onchange="refreshZoneSelect()">${vpcOpts}</select>`
        : `<p class="field-hint" style="color:var(--err)">No VPCs available — create a VPC first.</p>`}
    </div>
    <div class="field">
      <label for="f-zone">Zone</label>
      ${zoneOpts
        ? `<select id="f-zone">${zoneOpts}</select>`
        : `<p class="field-hint" style="color:var(--err)">No zones available.</p>`}
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
      if (type === 'vpc') {
        const region = document.getElementById('f-region')?.value;
        if (region) body.region = region;
      }
      if (type === 'subnet') {
        const vpcId = document.getElementById('f-vpc')?.value;
        if (!vpcId) { toast('Select a VPC', 'error'); return; }
        body.vpc = { id: vpcId };
        const zone = document.getElementById('f-zone')?.value;
        if (!zone) { toast('Select a zone', 'error'); return; }
        body.zone = zone;
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
