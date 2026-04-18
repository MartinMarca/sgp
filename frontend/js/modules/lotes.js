/**
 * SGP - Modulo Lotes
 * CRUD de lotes: listar, crear, editar, cerrar/vender, ver destetes
 */

const Lotes = (() => {
  let lotes = [];
  let granjas = [];
  let corrales = [];
  let granjaSeleccionada = null;
  let corralSeleccionado = null;
  let editingId = null;

  async function load() {
    const content = document.getElementById('contentArea');

    try {
      const data = await API.get('/granjas');
      granjas = data.data || [];
    } catch (e) { granjas = []; }

    if (granjas.length === 0) {
      content.innerHTML = `<div class="empty-state"><i class="bi bi-building d-block"></i><h6>No hay granjas</h6><p>Crea una granja primero.</p>
        <button class="btn btn-sgp" onclick="App.navigateTo('granjas')"><i class="bi bi-plus-lg me-2"></i>Crear Granja</button></div>`;
      return;
    }

    granjaSeleccionada = granjas[0].id;

    content.innerHTML = `
      <div class="d-flex align-items-center gap-3 mb-4 flex-wrap">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Granja:</label>
        <select class="form-select form-select-sm" id="selectGranjaLote" style="max-width:250px;"></select>
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Corral:</label>
        <select class="form-select form-select-sm" id="selectCorralLote" style="max-width:250px;"></select>
        <label class="form-label mb-0 fw-semibold ms-auto" style="white-space:nowrap;">Estado:</label>
        <select class="form-select form-select-sm" id="filtroEstadoLote" style="max-width:140px;">
          <option value="">Todos</option>
          <option value="activo" selected>Activo</option>
          <option value="cerrado">Cerrado</option>
          <option value="vendido">Vendido</option>
        </select>
      </div>

      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-collection me-2"></i>Lotes</h5>
          <button class="btn btn-sgp" id="btnNuevoLote"><i class="bi bi-plus-lg me-2"></i>Nuevo Lote</button>
        </div>
        <div id="lotesTableBody"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Modal Crear Lote -->
      <div class="modal fade" id="modalLote" tabindex="-1">
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="modalLoteTitle">Nuevo Lote</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalLoteAlert" class="alert d-none"></div>
              <form id="formLote" novalidate>
                <div class="mb-3" id="loteCorralGroup">
                  <label class="form-label">Corral <span class="text-danger">*</span></label>
                  <select class="form-select" id="loteCorralId"></select>
                </div>
                <div class="mb-3">
                  <label class="form-label">Nombre <span class="text-danger">*</span></label>
                  <input type="text" class="form-control" id="loteNombre" placeholder="Ej: Lote Feb-2026" required>
                </div>
                <div class="mb-3" id="loteFechaGroup">
                  <label class="form-label">Fecha de creacion</label>
                  <input type="date" class="form-control" id="loteFecha">
                  <small class="text-muted">Vacio = hoy</small>
                </div>
                <div class="mb-3" id="loteCantidadGroup" style="display:none;">
                  <label class="form-label">Cantidad de lechones</label>
                  <input type="number" class="form-control" id="loteCantidad" min="0" value="0">
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp btn-sm" id="btnGuardarLote"><i class="bi bi-check-lg me-1"></i>Guardar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Cerrar Lote -->
      <div class="modal fade" id="modalCerrarLote" tabindex="-1">
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Cerrar / Vender Lote</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalCerrarLoteAlert" class="alert d-none"></div>
              <p>Cerrar el lote <strong id="cerrarLoteNombre"></strong>?</p>
              <div class="mb-3">
                <label class="form-label">Accion <span class="text-danger">*</span></label>
                <select class="form-select" id="cerrarLoteEstado">
                  <option value="cerrado">Cerrar</option>
                  <option value="vendido">Marcar como vendido</option>
                </select>
              </div>
              <div class="mb-3">
                <label class="form-label">Motivo <span class="text-danger">*</span></label>
                <textarea class="form-control" id="cerrarLoteMotivo" rows="2" placeholder="Motivo de cierre o venta"></textarea>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-warning btn-sm" id="btnConfirmarCerrarLote"><i class="bi bi-lock me-1"></i>Confirmar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Detalle Destetes -->
      <div class="modal fade" id="modalDestetesLote" tabindex="-1">
        <div class="modal-dialog modal-dialog-scrollable">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title"><i class="bi bi-list-ul me-2"></i>Destetes del lote <span id="destetesLoteNombre"></span></h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body" id="destetesLoteBody">
              <div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cerrar</button>
            </div>
          </div>
        </div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaLote');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => {
      granjaSeleccionada = parseInt(selectGranja.value);
      loadCorrales();
    });

    document.getElementById('selectCorralLote').addEventListener('change', () => {
      corralSeleccionado = document.getElementById('selectCorralLote').value || null;
      fetchLotes();
    });
    document.getElementById('filtroEstadoLote').addEventListener('change', fetchLotes);

    document.getElementById('btnNuevoLote').addEventListener('click', openNuevoLote);
    document.getElementById('btnGuardarLote').addEventListener('click', handleSave);
    document.getElementById('formLote').addEventListener('submit', (e) => { e.preventDefault(); handleSave(); });

    await loadCorrales();
  }

  async function loadCorrales() {
    try {
      const data = await API.get(`/granjas/${granjaSeleccionada}/corrales`);
      corrales = data.data || [];
    } catch (e) { corrales = []; }

    const select = document.getElementById('selectCorralLote');
    select.innerHTML = '<option value="">Todos los corrales</option>' +
      corrales.map(c => `<option value="${c.id}">${esc(c.nombre)}</option>`).join('');
    corralSeleccionado = null;

    await fetchLotes();
  }

  async function fetchLotes() {
    const estado = document.getElementById('filtroEstadoLote').value;

    try {
      let allLotes = [];
      if (corralSeleccionado) {
        const url = `/corrales/${corralSeleccionado}/lotes` + (estado ? `?estado=${estado}` : '');
        const data = await API.get(url);
        allLotes = data.data || [];
      } else {
        // Traer lotes de todos los corrales de la granja
        const promises = corrales.map(c => API.get(`/corrales/${c.id}/lotes` + (estado ? `?estado=${estado}` : '')));
        const results = await Promise.all(promises);
        results.forEach(r => { allLotes = allLotes.concat(r.data || []); });
      }
      lotes = allLotes;
      renderTable();
    } catch (err) {
      document.getElementById('lotesTableBody').innerHTML = `<div class="p-4 text-center text-danger">Error: ${err.message}</div>`;
    }
  }

  function renderTable() {
    const container = document.getElementById('lotesTableBody');
    if (lotes.length === 0) {
      container.innerHTML = `<div class="empty-state"><i class="bi bi-collection d-block"></i><h6>No hay lotes</h6><p>Los lotes se crean al registrar destetes o manualmente.</p></div>`;
      return;
    }

    const rows = lotes.map(l => {
      let badge;
      if (l.estado === 'activo') badge = '<span class="badge-estado badge-activo">Activo</span>';
      else if (l.estado === 'vendido') badge = '<span class="badge-estado" style="background:#fff3cd;color:#856404;">Vendido</span>';
      else badge = '<span class="badge-estado badge-cerrado">Cerrado</span>';

      return `<tr>
        <td class="fw-semibold">${esc(l.nombre)}</td>
        <td>${l.corral ? esc(l.corral.nombre) : '-'}</td>
        <td class="fw-semibold">${l.cantidad_lechones}</td>
        <td>${badge}</td>
        <td><small class="text-muted">${fDate(l.fecha_creacion)}</small></td>
        <td>
          <div class="d-flex gap-1">
            <button class="btn btn-sm btn-outline-primary" title="Ver destetes" onclick="Lotes.verDestetes(${l.id})"><i class="bi bi-list-ul"></i></button>
            ${l.estado === 'activo' ? `
              <button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="Lotes.editar(${l.id})"><i class="bi bi-pencil"></i></button>
              <button class="btn btn-sm btn-outline-warning" title="Cerrar/Vender" onclick="Lotes.cerrar(${l.id})"><i class="bi bi-lock"></i></button>
            ` : ''}
          </div>
        </td>
      </tr>`;
    }).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr><th>Nombre</th><th>Corral</th><th>Lechones</th><th>Estado</th><th>Creacion</th><th style="width:120px;">Acciones</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>`;
  }

  function openNuevoLote() {
    editingId = null;
    document.getElementById('modalLoteTitle').textContent = 'Nuevo Lote';
    document.getElementById('loteNombre').value = '';
    document.getElementById('loteFecha').value = new Date().toISOString().split('T')[0];
    document.getElementById('loteCorralGroup').style.display = '';
    document.getElementById('loteFechaGroup').style.display = '';
    document.getElementById('loteCantidadGroup').style.display = 'none';
    document.getElementById('modalLoteAlert').classList.add('d-none');

    document.getElementById('loteCorralId').innerHTML = corrales.length
      ? corrales.map(c => `<option value="${c.id}">${esc(c.nombre)}</option>`).join('')
      : '<option value="">No hay corrales</option>';

    new bootstrap.Modal(document.getElementById('modalLote')).show();
    setTimeout(() => document.getElementById('loteNombre').focus(), 300);
  }

  function editar(id) {
    const l = lotes.find(x => x.id === id);
    if (!l) return;
    editingId = id;
    document.getElementById('modalLoteTitle').textContent = 'Editar Lote';
    document.getElementById('loteNombre').value = l.nombre;
    document.getElementById('loteCorralGroup').style.display = 'none';
    document.getElementById('loteFechaGroup').style.display = 'none';
    document.getElementById('loteCantidadGroup').style.display = '';
    document.getElementById('loteCantidad').value = l.cantidad_lechones;
    document.getElementById('modalLoteAlert').classList.add('d-none');

    new bootstrap.Modal(document.getElementById('modalLote')).show();
    setTimeout(() => document.getElementById('loteNombre').focus(), 300);
  }

  async function handleSave() {
    const nombre = document.getElementById('loteNombre').value.trim();
    const alert = document.getElementById('modalLoteAlert');
    const btn = document.getElementById('btnGuardarLote');

    if (!nombre) {
      alert.className = 'alert alert-warning'; alert.textContent = 'El nombre es obligatorio'; alert.classList.remove('d-none'); return;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';

    try {
      if (editingId) {
        const body = { nombre };
        const cant = parseInt(document.getElementById('loteCantidad').value);
        if (!isNaN(cant)) body.cantidad_lechones = cant;
        await API.put(`/lotes/${editingId}`, body);
        App.showToast('Lote actualizado');
      } else {
        const corralId = parseInt(document.getElementById('loteCorralId').value);
        if (!corralId) { alert.className = 'alert alert-warning'; alert.textContent = 'Selecciona un corral'; alert.classList.remove('d-none'); btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar'; return; }
        const body = { nombre };
        const fecha = document.getElementById('loteFecha').value;
        if (fecha) body.fecha = fecha;
        await API.post(`/corrales/${corralId}/lotes`, body);
        App.showToast('Lote creado');
      }
      bootstrap.Modal.getInstance(document.getElementById('modalLote')).hide();
      await fetchLotes();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';
    }
  }

  function cerrar(id) {
    const l = lotes.find(x => x.id === id);
    if (!l) return;
    document.getElementById('cerrarLoteNombre').textContent = l.nombre;
    document.getElementById('cerrarLoteEstado').value = 'cerrado';
    document.getElementById('cerrarLoteMotivo').value = '';
    document.getElementById('modalCerrarLoteAlert').classList.add('d-none');
    new bootstrap.Modal(document.getElementById('modalCerrarLote')).show();

    const btn = document.getElementById('btnConfirmarCerrarLote');
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener('click', async () => {
      const estado = document.getElementById('cerrarLoteEstado').value;
      const motivo = document.getElementById('cerrarLoteMotivo').value.trim();
      const al = document.getElementById('modalCerrarLoteAlert');

      if (!motivo) { al.className = 'alert alert-warning'; al.textContent = 'El motivo es obligatorio'; al.classList.remove('d-none'); return; }

      newBtn.disabled = true;
      try {
        await API.post(`/lotes/${id}/cerrar`, { estado, motivo_cierre: motivo });
        App.showToast(estado === 'vendido' ? 'Lote vendido' : 'Lote cerrado');
        bootstrap.Modal.getInstance(document.getElementById('modalCerrarLote')).hide();
        await fetchLotes();
      } catch (err) {
        al.className = 'alert alert-danger'; al.textContent = err.message; al.classList.remove('d-none');
      } finally { newBtn.disabled = false; }
    });
  }

  async function verDestetes(id) {
    const l = lotes.find(x => x.id === id);
    document.getElementById('destetesLoteNombre').textContent = l ? l.nombre : '#' + id;
    const body = document.getElementById('destetesLoteBody');
    body.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-success" role="status"></div></div>';
    new bootstrap.Modal(document.getElementById('modalDestetesLote')).show();

    try {
      const res = await API.get(`/lotes/${id}/destetes`);
      const destetes = res.data || [];

      if (destetes.length === 0) {
        body.innerHTML = '<div class="text-center text-muted py-3"><i class="bi bi-inbox" style="font-size:2rem;"></i><p class="mt-2 mb-0">No hay destetes asociados a este lote.</p></div>';
        return;
      }

      body.innerHTML = `
        <div class="table-responsive">
          <table class="table table-sm table-bordered mb-0">
            <thead class="table-light">
              <tr><th>Fecha</th><th>Cerda</th><th>Destetados</th></tr>
            </thead>
            <tbody>
              ${destetes.map(d => `<tr>
                <td>${fDate(d.fecha_destete)}</td>
                <td class="fw-semibold">${d.cerda ? esc(d.cerda.numero_caravana) : d.cerda_id}</td>
                <td>${d.cantidad_lechones_destetados}</td>
              </tr>`).join('')}
            </tbody>
          </table>
        </div>`;
    } catch (err) {
      body.innerHTML = `<div class="alert alert-danger">Error: ${err.message}</div>`;
    }
  }

  function fDate(d) { if (!d) return '-'; try { const p = d.split('T')[0]; return new Date(p + 'T12:00:00').toLocaleDateString('es-AR'); } catch { return d; } }
  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  return { load, editar, cerrar, verDestetes };
})();
