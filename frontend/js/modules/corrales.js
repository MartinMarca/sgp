/**
 * SGP - Modulo Corrales
 * CRUD de corrales asociados a una granja
 */

const Corrales = (() => {
  let corrales = [];
  let granjas = [];
  let granjaSeleccionada = null;
  let editingId = null;

  async function load() {
    const content = document.getElementById('contentArea');

    // Cargar granjas para el selector
    try {
      const data = await API.get('/granjas');
      granjas = data.data || [];
    } catch (e) {
      granjas = [];
    }

    if (granjas.length === 0) {
      content.innerHTML = `
        <div class="empty-state">
          <i class="bi bi-building d-block"></i>
          <h6>No hay granjas registradas</h6>
          <p>Primero crea una granja para poder agregar corrales.</p>
          <button class="btn btn-sgp" onclick="App.navigateTo('granjas')">
            <i class="bi bi-plus-lg me-2"></i>Crear Granja
          </button>
        </div>
      `;
      return;
    }

    granjaSeleccionada = granjas[0].id;

    content.innerHTML = `
      <!-- Selector de granja -->
      <div class="d-flex align-items-center gap-3 mb-4">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Granja:</label>
        <select class="form-select form-select-sm" id="selectGranjaCorral" style="max-width: 300px;"></select>
      </div>

      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-grid-3x3-gap me-2"></i>Corrales</h5>
          <button class="btn btn-sgp" id="btnNuevoCorral">
            <i class="bi bi-plus-lg me-2"></i>Nuevo Corral
          </button>
        </div>
        <div id="corralesTableBody">
          <div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div>
        </div>
      </div>

      <!-- Modal Crear/Editar -->
      <div class="modal fade" id="modalCorral" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="modalCorralTitle">Nuevo Corral</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalCorralAlert" class="alert d-none"></div>
              <form id="formCorral" novalidate>
                <div class="mb-3">
                  <label class="form-label">Nombre <span class="text-danger">*</span></label>
                  <input type="text" class="form-control" id="corralNombre" placeholder="Ej: Corral A1" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Descripcion</label>
                  <textarea class="form-control" id="corralDescripcion" rows="2" placeholder="Descripcion opcional"></textarea>
                </div>
                <div class="mb-3">
                  <label class="form-label">Capacidad maxima</label>
                  <input type="number" class="form-control" id="corralCapacidad" placeholder="Opcional" min="0">
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarCorral"><i class="bi bi-check-lg me-1"></i>Guardar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Eliminar -->
      <div class="modal fade" id="modalEliminarCorral" tabindex="-1">
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Confirmar eliminacion</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <p class="mb-0">Eliminar el corral <strong id="eliminarCorralNombre"></strong>?</p>
              <small class="text-muted">Solo si no tiene lotes activos.</small>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-danger btn-sm" id="btnConfirmarEliminarCorral"><i class="bi bi-trash me-1"></i>Eliminar</button>
            </div>
          </div>
        </div>
      </div>
    `;

    // Llenar selector de granjas
    const select = document.getElementById('selectGranjaCorral');
    select.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    select.addEventListener('change', () => {
      granjaSeleccionada = parseInt(select.value);
      fetchCorrales();
    });

    document.getElementById('btnNuevoCorral').addEventListener('click', () => openModal());
    document.getElementById('btnGuardarCorral').addEventListener('click', handleSave);
    document.getElementById('formCorral').addEventListener('submit', (e) => { e.preventDefault(); handleSave(); });

    await fetchCorrales();
  }

  async function fetchCorrales() {
    try {
      const data = await API.get(`/granjas/${granjaSeleccionada}/corrales`);
      corrales = data.data || [];
      renderTable();
    } catch (err) {
      document.getElementById('corralesTableBody').innerHTML = `<div class="p-4 text-center text-danger">Error: ${err.message}</div>`;
    }
  }

  function renderTable() {
    const container = document.getElementById('corralesTableBody');
    if (corrales.length === 0) {
      container.innerHTML = `<div class="empty-state"><i class="bi bi-grid-3x3-gap d-block"></i><h6>No hay corrales</h6><p>Agrega corrales a esta granja.</p></div>`;
      return;
    }

    const rows = corrales.map(c => `
      <tr>
        <td class="fw-semibold">${esc(c.nombre)}</td>
        <td class="d-none d-md-table-cell"><small class="text-muted">${c.descripcion ? esc(c.descripcion) : '-'}</small></td>
        <td>${c.total_animales ?? 0}</td>
        <td>${c.activo ? '<span class="badge-estado badge-activo">Activo</span>' : '<span class="badge-estado badge-cerrado">Inactivo</span>'}</td>
        <td>
          <div class="d-flex gap-1">
            <button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="Corrales.edit(${c.id})"><i class="bi bi-pencil"></i></button>
            <button class="btn btn-sm btn-outline-danger" title="Eliminar" onclick="Corrales.confirmDelete(${c.id})"><i class="bi bi-trash"></i></button>
          </div>
        </td>
      </tr>
    `).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr><th>Nombre</th><th class="d-none d-md-table-cell">Descripcion</th><th>Animales</th><th>Estado</th><th style="width:90px;">Acciones</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>
    `;
  }

  function openModal(corral = null) {
    editingId = corral ? corral.id : null;
    document.getElementById('modalCorralTitle').textContent = corral ? 'Editar Corral' : 'Nuevo Corral';
    document.getElementById('corralNombre').value = corral ? corral.nombre : '';
    document.getElementById('corralDescripcion').value = corral ? (corral.descripcion || '') : '';
    document.getElementById('corralCapacidad').value = corral && corral.capacidad_maxima != null ? corral.capacidad_maxima : '';
    document.getElementById('modalCorralAlert').classList.add('d-none');
    new bootstrap.Modal(document.getElementById('modalCorral')).show();
    setTimeout(() => document.getElementById('corralNombre').focus(), 300);
  }

  async function handleSave() {
    const nombre = document.getElementById('corralNombre').value.trim();
    const descripcion = document.getElementById('corralDescripcion').value.trim();
    const capacidad = document.getElementById('corralCapacidad').value;
    const alert = document.getElementById('modalCorralAlert');
    const btn = document.getElementById('btnGuardarCorral');

    if (!nombre) { alert.className = 'alert alert-warning'; alert.textContent = 'El nombre es obligatorio'; alert.classList.remove('d-none'); return; }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';

    try {
      const body = { nombre };
      if (descripcion) body.descripcion = descripcion;
      if (capacidad !== '') body.capacidad_maxima = parseInt(capacidad);

      if (editingId) {
        await API.put(`/corrales/${editingId}`, body);
        App.showToast('Corral actualizado');
      } else {
        await API.post(`/granjas/${granjaSeleccionada}/corrales`, body);
        App.showToast('Corral creado');
      }

      bootstrap.Modal.getInstance(document.getElementById('modalCorral')).hide();
      await fetchCorrales();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';
    }
  }

  function edit(id) {
    const c = corrales.find(x => x.id === id);
    if (c) openModal(c);
  }

  function confirmDelete(id) {
    const c = corrales.find(x => x.id === id);
    if (!c) return;
    document.getElementById('eliminarCorralNombre').textContent = c.nombre;
    new bootstrap.Modal(document.getElementById('modalEliminarCorral')).show();

    const btn = document.getElementById('btnConfirmarEliminarCorral');
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener('click', async () => {
      newBtn.disabled = true;
      try {
        await API.del(`/corrales/${id}`);
        App.showToast('Corral eliminado');
        bootstrap.Modal.getInstance(document.getElementById('modalEliminarCorral')).hide();
        await fetchCorrales();
      } catch (err) { App.showToast(err.message, 'danger'); }
      finally { newBtn.disabled = false; }
    });
  }

  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  return { load, edit, confirmDelete };
})();
