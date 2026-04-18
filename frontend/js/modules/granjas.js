/**
 * SGP - Modulo Granjas
 * CRUD completo de granjas
 */

const Granjas = (() => {
  let granjas = [];
  let editingId = null;

  // --- Carga principal ---

  async function load() {
    const content = document.getElementById('contentArea');

    content.innerHTML = `
      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-building me-2"></i>Mis Granjas</h5>
          <button class="btn btn-sgp" id="btnNuevaGranja">
            <i class="bi bi-plus-lg me-2"></i>Nueva Granja
          </button>
        </div>
        <div id="granjasTableBody">
          <div class="loading-spinner">
            <div class="spinner-border text-success" role="status"></div>
          </div>
        </div>
      </div>

      <!-- Modal Crear/Editar -->
      <div class="modal fade" id="modalGranja" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="modalGranjaTitle">Nueva Granja</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalGranjaAlert" class="alert d-none"></div>
              <form id="formGranja" novalidate>
                <div class="mb-3">
                  <label for="granjaNombre" class="form-label">Nombre <span class="text-danger">*</span></label>
                  <input type="text" class="form-control" id="granjaNombre" placeholder="Ej: Granja Don Pedro" required>
                </div>
                <div class="mb-3">
                  <label for="granjaDescripcion" class="form-label">Descripcion</label>
                  <textarea class="form-control" id="granjaDescripcion" rows="2" placeholder="Descripcion de la granja"></textarea>
                </div>
                <div class="mb-3">
                  <label for="granjaUbicacion" class="form-label">Ubicacion</label>
                  <input type="text" class="form-control" id="granjaUbicacion" placeholder="Ej: Ruta 5 km 120, Pergamino">
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarGranja">
                <i class="bi bi-check-lg me-1"></i>Guardar
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Confirmar Eliminar -->
      <div class="modal fade" id="modalEliminarGranja" tabindex="-1">
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Confirmar eliminacion</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <p class="mb-0">Estas seguro de que queres eliminar la granja <strong id="eliminarGranjaNombre"></strong>?</p>
              <small class="text-muted">Solo se puede eliminar si no tiene cerdas, padrillos o corrales activos.</small>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-danger btn-sm" id="btnConfirmarEliminar">
                <i class="bi bi-trash me-1"></i>Eliminar
              </button>
            </div>
          </div>
        </div>
      </div>
    `;

    bindEvents();
    await fetchGranjas();
  }

  // --- Eventos ---

  function bindEvents() {
    document.getElementById('btnNuevaGranja').addEventListener('click', () => openModal());
    document.getElementById('btnGuardarGranja').addEventListener('click', handleSave);
    document.getElementById('formGranja').addEventListener('submit', (e) => {
      e.preventDefault();
      handleSave();
    });
  }

  // --- Fetch y render ---

  async function fetchGranjas() {
    try {
      const data = await API.get('/granjas');
      granjas = data.data || [];
      renderTable();
    } catch (err) {
      document.getElementById('granjasTableBody').innerHTML = `
        <div class="p-4 text-center text-danger">
          <i class="bi bi-exclamation-triangle me-2"></i>Error cargando granjas: ${err.message}
        </div>
      `;
    }
  }

  function renderTable() {
    const container = document.getElementById('granjasTableBody');

    if (granjas.length === 0) {
      container.innerHTML = `
        <div class="empty-state">
          <i class="bi bi-building d-block"></i>
          <h6>No hay granjas registradas</h6>
          <p>Crea tu primera granja con el boton de arriba.</p>
        </div>
      `;
      return;
    }

    const rows = granjas.map(g => `
      <tr>
        <td>
          <div class="fw-semibold">${escapeHtml(g.nombre)}</div>
          ${g.ubicacion ? `<small class="text-muted"><i class="bi bi-geo-alt me-1"></i>${escapeHtml(g.ubicacion)}</small>` : ''}
        </td>
        <td class="d-none d-md-table-cell">
          <small class="text-muted">${g.descripcion ? escapeHtml(g.descripcion) : '-'}</small>
        </td>
        <td>
          ${g.activo
            ? '<span class="badge-estado badge-activo">Activa</span>'
            : '<span class="badge-estado badge-cerrado">Inactiva</span>'
          }
        </td>
        <td class="d-none d-lg-table-cell">
          <small class="text-muted">${App.formatDate(g.created_at)}</small>
        </td>
        <td>
          <div class="d-flex gap-1">
            <button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="Granjas.edit(${g.id})">
              <i class="bi bi-pencil"></i>
            </button>
            <button class="btn btn-sm btn-outline-danger" title="Eliminar" onclick="Granjas.confirmDelete(${g.id})">
              <i class="bi bi-trash"></i>
            </button>
          </div>
        </td>
      </tr>
    `).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead>
          <tr>
            <th>Nombre</th>
            <th class="d-none d-md-table-cell">Descripcion</th>
            <th>Estado</th>
            <th class="d-none d-lg-table-cell">Creada</th>
            <th style="width: 90px;">Acciones</th>
          </tr>
        </thead>
        <tbody>${rows}</tbody>
      </table>
    `;
  }

  // --- Modal Crear/Editar ---

  function openModal(granja = null) {
    editingId = granja ? granja.id : null;

    document.getElementById('modalGranjaTitle').textContent = granja ? 'Editar Granja' : 'Nueva Granja';
    document.getElementById('granjaNombre').value = granja ? granja.nombre : '';
    document.getElementById('granjaDescripcion').value = granja ? (granja.descripcion || '') : '';
    document.getElementById('granjaUbicacion').value = granja ? (granja.ubicacion || '') : '';

    // Limpiar alertas
    const alert = document.getElementById('modalGranjaAlert');
    alert.classList.add('d-none');
    alert.textContent = '';

    const modal = new bootstrap.Modal(document.getElementById('modalGranja'));
    modal.show();

    // Focus en nombre
    setTimeout(() => document.getElementById('granjaNombre').focus(), 300);
  }

  async function handleSave() {
    const nombre = document.getElementById('granjaNombre').value.trim();
    const descripcion = document.getElementById('granjaDescripcion').value.trim();
    const ubicacion = document.getElementById('granjaUbicacion').value.trim();
    const alert = document.getElementById('modalGranjaAlert');
    const btn = document.getElementById('btnGuardarGranja');

    if (!nombre) {
      alert.className = 'alert alert-warning';
      alert.textContent = 'El nombre es obligatorio';
      alert.classList.remove('d-none');
      return;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';

    try {
      const body = { nombre };
      if (descripcion) body.descripcion = descripcion;
      if (ubicacion) body.ubicacion = ubicacion;

      if (editingId) {
        await API.put(`/granjas/${editingId}`, body);
        App.showToast('Granja actualizada correctamente');
      } else {
        await API.post('/granjas', body);
        App.showToast('Granja creada correctamente');
      }

      bootstrap.Modal.getInstance(document.getElementById('modalGranja')).hide();
      await fetchGranjas();
    } catch (err) {
      alert.className = 'alert alert-danger';
      alert.textContent = err.message;
      alert.classList.remove('d-none');
    } finally {
      btn.disabled = false;
      btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';
    }
  }

  // --- Editar ---

  function edit(id) {
    const granja = granjas.find(g => g.id === id);
    if (granja) openModal(granja);
  }

  // --- Eliminar ---

  function confirmDelete(id) {
    const granja = granjas.find(g => g.id === id);
    if (!granja) return;

    editingId = id;
    document.getElementById('eliminarGranjaNombre').textContent = granja.nombre;

    const modal = new bootstrap.Modal(document.getElementById('modalEliminarGranja'));
    modal.show();

    // Bind del boton confirmar (reemplazar para evitar duplicados)
    const btn = document.getElementById('btnConfirmarEliminar');
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener('click', async () => {
      newBtn.disabled = true;
      newBtn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Eliminando...';

      try {
        await API.del(`/granjas/${id}`);
        App.showToast('Granja eliminada correctamente');
        bootstrap.Modal.getInstance(document.getElementById('modalEliminarGranja')).hide();
        await fetchGranjas();
      } catch (err) {
        App.showToast(err.message, 'danger');
      } finally {
        newBtn.disabled = false;
        newBtn.innerHTML = '<i class="bi bi-trash me-1"></i>Eliminar';
      }
    });
  }

  // --- Utils ---

  function escapeHtml(str) {
    if (!str) return '';
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
  }

  return {
    load,
    edit,
    confirmDelete,
  };
})();
