/**
 * SGP - Modulo Padrillos
 * CRUD de padrillos + baja
 */

const Padrillos = (() => {
  let padrillos = [];
  let granjas = [];
  let granjaSeleccionada = null;
  let editingId = null;

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
          <p>Primero crea una granja para poder registrar padrillos.</p>
          <button class="btn btn-sgp" onclick="App.navigateTo('granjas')"><i class="bi bi-plus-lg me-2"></i>Crear Granja</button>
        </div>`;
      return;
    }

    granjaSeleccionada = granjas[0].id;

    content.innerHTML = `
      <div class="d-flex align-items-center gap-3 mb-4">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Granja:</label>
        <select class="form-select form-select-sm" id="selectGranjaPadrillo" style="max-width:300px;"></select>
      </div>

      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-gender-male me-2"></i>Padrillos</h5>
          <button class="btn btn-sgp" id="btnNuevoPadrillo"><i class="bi bi-plus-lg me-2"></i>Nuevo Padrillo</button>
        </div>
        <div id="padrillosTableBody"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Modal Crear/Editar -->
      <div class="modal fade" id="modalPadrillo" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="modalPadrilloTitle">Nuevo Padrillo</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalPadrilloAlert" class="alert d-none"></div>
              <form id="formPadrillo" novalidate>
                <div class="mb-3">
                  <label class="form-label">Numero de caravana <span class="text-danger">*</span></label>
                  <input type="text" class="form-control" id="padrilloCaravana" placeholder="Ej: P-001" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Nombre <span class="text-danger">*</span></label>
                  <input type="text" class="form-control" id="padrilloNombre" placeholder="Ej: Titan" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Genetica</label>
                  <input type="text" class="form-control" id="padrilloGenetica" placeholder="Ej: Duroc">
                </div>
                <div class="mb-3">
                  <label class="form-label">Fecha ultima vacunacion</label>
                  <input type="date" class="form-control" id="padrilloVacunacion">
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarPadrillo"><i class="bi bi-check-lg me-1"></i>Guardar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Baja -->
      <div class="modal fade" id="modalBajaPadrillo" tabindex="-1">
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Dar de baja</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <p>Dar de baja a <strong id="bajaPadrilloNombre"></strong>?</p>
              <div id="modalBajaPadrilloAlert" class="alert d-none"></div>
              <div class="mb-3">
                <label class="form-label">Motivo <span class="text-danger">*</span></label>
                <select class="form-select" id="bajaPadrilloMotivo">
                  <option value="muerte">Muerte</option>
                  <option value="venta">Venta</option>
                </select>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-danger btn-sm" id="btnConfirmarBajaPadrillo"><i class="bi bi-x-circle me-1"></i>Dar de baja</button>
            </div>
          </div>
        </div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaPadrillo');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => { granjaSeleccionada = parseInt(selectGranja.value); fetchPadrillos(); });

    document.getElementById('btnNuevoPadrillo').addEventListener('click', () => openModal());
    document.getElementById('btnGuardarPadrillo').addEventListener('click', handleSave);
    document.getElementById('formPadrillo').addEventListener('submit', (e) => { e.preventDefault(); handleSave(); });

    await fetchPadrillos();
  }

  async function fetchPadrillos() {
    try {
      const data = await API.get(`/granjas/${granjaSeleccionada}/padrillos`);
      padrillos = data.data || [];
      renderTable();
    } catch (err) {
      document.getElementById('padrillosTableBody').innerHTML = `<div class="p-4 text-center text-danger">Error: ${err.message}</div>`;
    }
  }

  function renderTable() {
    const container = document.getElementById('padrillosTableBody');
    if (padrillos.length === 0) {
      container.innerHTML = `<div class="empty-state"><i class="bi bi-gender-male d-block"></i><h6>No hay padrillos</h6><p>Registra padrillos en esta granja.</p></div>`;
      return;
    }

    const rows = padrillos.map(p => `
      <tr>
        <td><span class="fw-semibold">${esc(p.numero_caravana)}</span></td>
        <td>${esc(p.nombre)}</td>
        <td class="d-none d-md-table-cell"><small class="text-muted">${p.genetica ? esc(p.genetica) : '-'}</small></td>
        <td>
          ${p.activo
            ? '<span class="badge-estado badge-activo">Activo</span>'
            : `<span class="badge-estado badge-cerrado">${p.motivo_baja || 'Baja'}</span>`
          }
        </td>
        <td>
          <div class="d-flex gap-1">
            <button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="Padrillos.edit(${p.id})"><i class="bi bi-pencil"></i></button>
            ${p.activo ? `<button class="btn btn-sm btn-outline-danger" title="Dar de baja" onclick="Padrillos.confirmBaja(${p.id})"><i class="bi bi-x-circle"></i></button>` : ''}
          </div>
        </td>
      </tr>
    `).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr><th>Caravana</th><th>Nombre</th><th class="d-none d-md-table-cell">Genetica</th><th>Estado</th><th style="width:90px;">Acciones</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>
    `;
  }

  function openModal(padrillo = null) {
    editingId = padrillo ? padrillo.id : null;
    document.getElementById('modalPadrilloTitle').textContent = padrillo ? 'Editar Padrillo' : 'Nuevo Padrillo';
    document.getElementById('padrilloCaravana').value = padrillo ? padrillo.numero_caravana : '';
    document.getElementById('padrilloNombre').value = padrillo ? padrillo.nombre : '';
    document.getElementById('padrilloGenetica').value = padrillo ? (padrillo.genetica || '') : '';
    document.getElementById('padrilloVacunacion').value = '';
    document.getElementById('modalPadrilloAlert').classList.add('d-none');
    new bootstrap.Modal(document.getElementById('modalPadrillo')).show();
    setTimeout(() => document.getElementById('padrilloCaravana').focus(), 300);
  }

  async function handleSave() {
    const caravana = document.getElementById('padrilloCaravana').value.trim();
    const nombre = document.getElementById('padrilloNombre').value.trim();
    const genetica = document.getElementById('padrilloGenetica').value.trim();
    const vacunacion = document.getElementById('padrilloVacunacion').value;
    const alert = document.getElementById('modalPadrilloAlert');
    const btn = document.getElementById('btnGuardarPadrillo');

    if (!caravana || !nombre) {
      alert.className = 'alert alert-warning'; alert.textContent = 'Caravana y nombre son obligatorios'; alert.classList.remove('d-none'); return;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';

    try {
      const body = { numero_caravana: caravana, nombre };
      if (genetica) body.genetica = genetica;
      if (vacunacion) body.fecha_ultima_vacunacion = vacunacion;

      if (editingId) {
        await API.put(`/padrillos/${editingId}`, body);
        App.showToast('Padrillo actualizado');
      } else {
        await API.post(`/granjas/${granjaSeleccionada}/padrillos`, body);
        App.showToast('Padrillo registrado');
      }
      bootstrap.Modal.getInstance(document.getElementById('modalPadrillo')).hide();
      await fetchPadrillos();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';
    }
  }

  function edit(id) {
    const p = padrillos.find(x => x.id === id);
    if (p) openModal(p);
  }

  function confirmBaja(id) {
    const p = padrillos.find(x => x.id === id);
    if (!p) return;
    document.getElementById('bajaPadrilloNombre').textContent = `${p.nombre} (${p.numero_caravana})`;
    document.getElementById('modalBajaPadrilloAlert').classList.add('d-none');
    document.getElementById('bajaPadrilloMotivo').value = 'muerte';
    new bootstrap.Modal(document.getElementById('modalBajaPadrillo')).show();

    const btn = document.getElementById('btnConfirmarBajaPadrillo');
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener('click', async () => {
      const motivo = document.getElementById('bajaPadrilloMotivo').value;
      newBtn.disabled = true;
      try {
        await API.post(`/padrillos/${id}/baja`, { motivo_baja: motivo });
        App.showToast('Padrillo dado de baja');
        bootstrap.Modal.getInstance(document.getElementById('modalBajaPadrillo')).hide();
        await fetchPadrillos();
      } catch (err) {
        const al = document.getElementById('modalBajaPadrilloAlert');
        al.className = 'alert alert-danger'; al.textContent = err.message; al.classList.remove('d-none');
      } finally { newBtn.disabled = false; }
    });
  }

  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  return { load, edit, confirmBaja };
})();
