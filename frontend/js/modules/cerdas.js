/**
 * SGP - Modulo Cerdas
 * CRUD de cerdas + baja + filtro por estado
 */

const Cerdas = (() => {
  let cerdas = [];
  let granjas = [];
  let granjaSeleccionada = null;
  let editingId = null;

  const ESTADOS = ['disponible', 'servicio', 'gestacion', 'cria'];

  async function load() {
    const content = document.getElementById('contentArea');

    try {
      const data = await API.get('/granjas');
      granjas = data.data || [];
    } catch (e) { granjas = []; }

    if (granjas.length === 0) {
      content.innerHTML = `
        <div class="empty-state">
          <i class="bi bi-building d-block"></i><h6>No hay granjas registradas</h6>
          <p>Primero crea una granja para poder registrar cerdas.</p>
          <button class="btn btn-sgp" onclick="App.navigateTo('granjas')"><i class="bi bi-plus-lg me-2"></i>Crear Granja</button>
        </div>`;
      return;
    }

    granjaSeleccionada = granjas[0].id;

    content.innerHTML = `
      <div class="d-flex align-items-center gap-3 mb-4 flex-wrap">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Granja:</label>
        <select class="form-select form-select-sm" id="selectGranjaCerda" style="max-width:300px;"></select>
        <div class="input-group input-group-sm ms-auto" style="max-width:200px;">
          <span class="input-group-text"><i class="bi bi-search"></i></span>
          <input type="text" class="form-control" id="filtroCaravanaCerda" placeholder="Buscar caravana...">
        </div>
        <select class="form-select form-select-sm" id="filtroEstadoCerda" style="max-width:160px;">
          <option value="">Todos</option>
          ${ESTADOS.map(e => `<option value="${e}">${e}</option>`).join('')}
        </select>
      </div>

      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-gender-female me-2"></i>Cerdas</h5>
          <button class="btn btn-sgp" id="btnNuevaCerda"><i class="bi bi-plus-lg me-2"></i>Nueva Cerda</button>
        </div>
        <div id="cerdasTableBody"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Modal Crear/Editar -->
      <div class="modal fade" id="modalCerda" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="modalCerdaTitle">Nueva Cerda</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalCerdaAlert" class="alert d-none"></div>
              <form id="formCerda" novalidate>
                <div class="mb-3">
                  <label class="form-label">Numero de caravana <span class="text-danger">*</span></label>
                  <input type="text" class="form-control" id="cerdaCaravana" placeholder="Ej: C-001" required>
                </div>
                <div class="mb-3" id="cerdaEstadoGroup">
                  <label class="form-label">Estado inicial</label>
                  <select class="form-select" id="cerdaEstado">
                    ${ESTADOS.map(e => `<option value="${e}">${e}</option>`).join('')}
                  </select>
                </div>
                <div class="mb-3">
                  <label class="form-label">Genetica</label>
                  <input type="text" class="form-control" id="cerdaGenetica" placeholder="Ej: Landrace x Large White">
                </div>
                <div class="mb-3">
                  <label class="form-label">Detalle de pelaje</label>
                  <input type="text" class="form-control" id="cerdaPelaje" placeholder="Ej: Blanca con manchas negras">
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarCerda"><i class="bi bi-check-lg me-1"></i>Guardar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Baja -->
      <div class="modal fade" id="modalBajaCerda" tabindex="-1">
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Dar de baja</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <p>Dar de baja a <strong id="bajaCerdaNombre"></strong>?</p>
              <div id="modalBajaCerdaAlert" class="alert d-none"></div>
              <div class="mb-3">
                <label class="form-label">Motivo <span class="text-danger">*</span></label>
                <select class="form-select" id="bajaCerdaMotivo">
                  <option value="muerte">Muerte</option>
                  <option value="venta">Venta</option>
                </select>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-danger btn-sm" id="btnConfirmarBajaCerda"><i class="bi bi-x-circle me-1"></i>Dar de baja</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Historial -->
      <div class="modal fade" id="modalHistorialCerda" tabindex="-1">
        <div class="modal-dialog modal-lg modal-dialog-scrollable">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title"><i class="bi bi-clock-history me-2"></i>Historial de <span id="historialCerdaNombre"></span></h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body" id="historialCerdaBody">
              <div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cerrar</button>
            </div>
          </div>
        </div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaCerda');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => { granjaSeleccionada = parseInt(selectGranja.value); fetchCerdas(); });

    document.getElementById('filtroEstadoCerda').addEventListener('change', () => fetchCerdas());
    document.getElementById('filtroCaravanaCerda').addEventListener('input', () => renderTable());
    document.getElementById('btnNuevaCerda').addEventListener('click', () => openModal());
    document.getElementById('btnGuardarCerda').addEventListener('click', handleSave);
    document.getElementById('formCerda').addEventListener('submit', (e) => { e.preventDefault(); handleSave(); });

    await fetchCerdas();
  }

  async function fetchCerdas() {
    const estado = document.getElementById('filtroEstadoCerda').value;
    try {
      let url = `/granjas/${granjaSeleccionada}/cerdas`;
      if (estado) url += `?estado=${estado}`;
      const data = await API.get(url);
      cerdas = data.data || [];
      renderTable();
    } catch (err) {
      document.getElementById('cerdasTableBody').innerHTML = `<div class="p-4 text-center text-danger">Error: ${err.message}</div>`;
    }
  }

  function renderTable() {
    const container = document.getElementById('cerdasTableBody');
    const filtroCaravana = (document.getElementById('filtroCaravanaCerda')?.value || '').trim().toLowerCase();
    const filtradas = filtroCaravana
      ? cerdas.filter(c => c.numero_caravana.toLowerCase().includes(filtroCaravana))
      : cerdas;

    if (filtradas.length === 0) {
      container.innerHTML = cerdas.length === 0
        ? `<div class="empty-state"><i class="bi bi-gender-female d-block"></i><h6>No hay cerdas</h6><p>Registra cerdas en esta granja.</p></div>`
        : `<div class="empty-state"><i class="bi bi-search d-block"></i><h6>Sin resultados</h6><p>No se encontraron cerdas con esa caravana.</p></div>`;
      return;
    }

    const rows = filtradas.map(c => `
      <tr>
        <td><span class="fw-semibold">${esc(c.numero_caravana)}</span></td>
        <td>${App.badgeEstado(c.estado)}</td>
        <td class="d-none d-md-table-cell"><small class="text-muted">${c.genetica ? esc(c.genetica) : '-'}</small></td>
        <td class="d-none d-md-table-cell"><small class="text-muted">${c.detalle_pelaje ? esc(c.detalle_pelaje) : '-'}</small></td>
        <td>
          <div class="d-flex gap-1">
            <button class="btn btn-sm btn-outline-primary" title="Historial" onclick="Cerdas.showHistorial(${c.id})"><i class="bi bi-clock-history"></i></button>
            <button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="Cerdas.edit(${c.id})"><i class="bi bi-pencil"></i></button>
            ${c.activo ? `<button class="btn btn-sm btn-outline-danger" title="Dar de baja" onclick="Cerdas.confirmBaja(${c.id})"><i class="bi bi-x-circle"></i></button>` : ''}
          </div>
        </td>
      </tr>
    `).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr><th>Caravana</th><th>Estado</th><th class="d-none d-md-table-cell">Genetica</th><th class="d-none d-md-table-cell">Pelaje</th><th style="width:120px;">Acciones</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>
    `;
  }

  function openModal(cerda = null) {
    editingId = cerda ? cerda.id : null;
    document.getElementById('modalCerdaTitle').textContent = cerda ? 'Editar Cerda' : 'Nueva Cerda';
    document.getElementById('cerdaCaravana').value = cerda ? cerda.numero_caravana : '';
    document.getElementById('cerdaEstado').value = cerda ? cerda.estado : 'disponible';
    document.getElementById('cerdaGenetica').value = cerda ? (cerda.genetica || '') : '';
    document.getElementById('cerdaPelaje').value = cerda ? (cerda.detalle_pelaje || '') : '';
    // Ocultar estado en edicion (no se cambia manualmente, se cambia por el ciclo)
    document.getElementById('cerdaEstadoGroup').style.display = cerda ? 'none' : '';
    document.getElementById('modalCerdaAlert').classList.add('d-none');
    new bootstrap.Modal(document.getElementById('modalCerda')).show();
    setTimeout(() => document.getElementById('cerdaCaravana').focus(), 300);
  }

  async function handleSave() {
    const caravana = document.getElementById('cerdaCaravana').value.trim();
    const estado = document.getElementById('cerdaEstado').value;
    const genetica = document.getElementById('cerdaGenetica').value.trim();
    const pelaje = document.getElementById('cerdaPelaje').value.trim();
    const alert = document.getElementById('modalCerdaAlert');
    const btn = document.getElementById('btnGuardarCerda');

    if (!caravana) { alert.className = 'alert alert-warning'; alert.textContent = 'La caravana es obligatoria'; alert.classList.remove('d-none'); return; }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';

    try {
      if (editingId) {
        const body = { numero_caravana: caravana };
        if (genetica) body.genetica = genetica;
        if (pelaje) body.detalle_pelaje = pelaje;
        await API.put(`/cerdas/${editingId}`, body);
        App.showToast('Cerda actualizada');
      } else {
        const body = { numero_caravana: caravana, estado };
        if (genetica) body.genetica = genetica;
        if (pelaje) body.detalle_pelaje = pelaje;
        await API.post(`/granjas/${granjaSeleccionada}/cerdas`, body);
        App.showToast('Cerda registrada');
      }
      bootstrap.Modal.getInstance(document.getElementById('modalCerda')).hide();
      await fetchCerdas();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';
    }
  }

  function edit(id) {
    const c = cerdas.find(x => x.id === id);
    if (c) openModal(c);
  }

  function confirmBaja(id) {
    const c = cerdas.find(x => x.id === id);
    if (!c) return;
    editingId = id;
    document.getElementById('bajaCerdaNombre').textContent = c.numero_caravana;
    document.getElementById('modalBajaCerdaAlert').classList.add('d-none');
    document.getElementById('bajaCerdaMotivo').value = 'muerte';
    new bootstrap.Modal(document.getElementById('modalBajaCerda')).show();

    const btn = document.getElementById('btnConfirmarBajaCerda');
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener('click', async () => {
      const motivo = document.getElementById('bajaCerdaMotivo').value;
      newBtn.disabled = true;
      try {
        await API.post(`/cerdas/${id}/baja`, { motivo_baja: motivo });
        App.showToast('Cerda dada de baja');
        bootstrap.Modal.getInstance(document.getElementById('modalBajaCerda')).hide();
        await fetchCerdas();
      } catch (err) {
        const al = document.getElementById('modalBajaCerdaAlert');
        al.className = 'alert alert-danger'; al.textContent = err.message; al.classList.remove('d-none');
      } finally { newBtn.disabled = false; }
    });
  }

  async function showHistorial(id) {
    const c = cerdas.find(x => x.id === id);
    if (!c) return;

    document.getElementById('historialCerdaNombre').textContent = c.numero_caravana;
    const body = document.getElementById('historialCerdaBody');
    body.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-success" role="status"></div></div>';
    new bootstrap.Modal(document.getElementById('modalHistorialCerda')).show();

    try {
      const res = await API.get(`/cerdas/${id}/historial`);
      const data = res.data || {};

      const eventos = [];

      (data.servicios || []).forEach(s => {
        let detalle = s.tipo_monta === 'natural' ? 'Monta natural' : 'Inseminacion';
        if (s.tipo_monta === 'natural' && s.padrillo) {
          detalle += ' — ' + esc(s.padrillo.nombre) + ' (' + esc(s.padrillo.numero_caravana) + ')';
        } else if (s.tipo_monta === 'inseminacion' && s.numero_pajuela) {
          detalle += ' — Pajuela: ' + esc(s.numero_pajuela);
        }
        if (s.tiene_repeticiones) detalle += ' · ' + s.cantidad_repeticiones + ' rep.';
        let prenez = '';
        if (s.prenez_confirmada) prenez = '<span class="badge bg-success ms-2">Prenez confirmada</span>';
        else if (s.prenez_cancelada) prenez = '<span class="badge bg-danger ms-2">Prenez cancelada</span>';
        else prenez = '<span class="badge bg-secondary ms-2">Prenez pendiente</span>';

        eventos.push({ fecha: s.fecha_servicio, tipo: 'servicio', icon: 'bi-heart-pulse', color: '#dc3545', label: 'Servicio', detalle: detalle + prenez });
      });

      (data.partos || []).forEach(p => {
        const detalle = `${p.lechones_nacidos_vivos} vivos de ${p.lechones_nacidos_totales} totales (${p.lechones_hembras}♀ ${p.lechones_machos}♂)`;
        eventos.push({ fecha: p.fecha_parto, tipo: 'parto', icon: 'bi-plus-circle', color: '#0d6efd', label: 'Parto', detalle });
      });

      (data.destetes || []).forEach(d => {
        let detalle = `${d.cantidad_lechones_destetados} lechones destetados`;
        if (d.lote) detalle += ' → Lote: ' + esc(d.lote.nombre || '#' + d.lote.id);
        eventos.push({ fecha: d.fecha_destete, tipo: 'destete', icon: 'bi-box-arrow-right', color: '#198754', label: 'Destete', detalle });
      });

      if (eventos.length === 0) {
        body.innerHTML = `
          <div class="text-center text-muted py-4">
            <i class="bi bi-inbox" style="font-size:2.5rem;"></i>
            <p class="mt-2 mb-0">Esta cerda aun no tiene registros en su historial.</p>
          </div>`;
        return;
      }

      eventos.sort((a, b) => new Date(b.fecha.split('T')[0]) - new Date(a.fecha.split('T')[0]));

      body.innerHTML = `
        <div class="timeline-historial">
          ${eventos.map(e => `
            <div class="d-flex align-items-start mb-3">
              <div class="text-center me-3" style="min-width:38px;">
                <div style="width:38px;height:38px;border-radius:50%;background:${e.color}15;display:flex;align-items:center;justify-content:center;">
                  <i class="bi ${e.icon}" style="color:${e.color};font-size:1.1rem;"></i>
                </div>
              </div>
              <div style="flex:1;border-bottom:1px solid #eee;padding-bottom:0.75rem;">
                <div class="d-flex align-items-center gap-2 mb-1">
                  <span class="badge" style="background:${e.color};font-size:0.7rem;">${e.label}</span>
                  <small class="text-muted">${fDate(e.fecha)}</small>
                </div>
                <div style="font-size:0.9rem;">${e.detalle}</div>
              </div>
            </div>
          `).join('')}
        </div>`;
    } catch (err) {
      body.innerHTML = `<div class="alert alert-danger">Error al cargar historial: ${err.message}</div>`;
    }
  }

  function fDate(d) {
    if (!d) return '-';
    try { const p = d.split('T')[0]; return new Date(p + 'T12:00:00').toLocaleDateString('es-AR'); } catch { return d; }
  }

  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  return { load, edit, confirmBaja, showHistorial };
})();
