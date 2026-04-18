/**
 * SGP - Modulo Destetes
 * Registro de destetes: cerda en cria -> disponible, lechones a lote
 */

const Destetes = (() => {
  let destetes = [];
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
      content.innerHTML = `<div class="empty-state"><i class="bi bi-building d-block"></i><h6>No hay granjas</h6><p>Crea una granja primero.</p>
        <button class="btn btn-sgp" onclick="App.navigateTo('granjas')"><i class="bi bi-plus-lg me-2"></i>Crear Granja</button></div>`;
      return;
    }

    granjaSeleccionada = granjas[0].id;
    const hoy = new Date();
    const mesActual = hoy.getMonth() + 1;
    const anioActual = hoy.getFullYear();

    content.innerHTML = `
      <div class="d-flex align-items-center gap-3 mb-4 flex-wrap">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Granja:</label>
        <select class="form-select form-select-sm" id="selectGranjaDestete" style="max-width:250px;"></select>
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Mes:</label>
        <input type="number" class="form-control form-control-sm" id="filtroMesDestete" min="1" max="12" value="${mesActual}" style="max-width:70px;">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Ano:</label>
        <input type="number" class="form-control form-control-sm" id="filtroAnioDestete" min="2020" value="${anioActual}" style="max-width:90px;">
        <button class="btn btn-outline-secondary btn-sm" id="btnFiltrarDestete"><i class="bi bi-funnel me-1"></i>Filtrar</button>
      </div>

      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-box-arrow-right me-2"></i>Destetes</h5>
          <button class="btn btn-sgp" id="btnNuevoDestete"><i class="bi bi-plus-lg me-2"></i>Registrar Destete</button>
        </div>
        <div id="destetesTableBody"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Modal Crear Destete -->
      <div class="modal fade" id="modalDestete" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Registrar Destete</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalDesteteAlert" class="alert d-none"></div>
              <form id="formDestete" novalidate>
                <div class="mb-3">
                  <label class="form-label">Cerda <span class="text-danger">*</span></label>
                  <select class="form-select" id="desteteCerdaId"></select>
                  <small class="text-muted">Solo cerdas en estado "cria"</small>
                </div>
                <div class="mb-3">
                  <label class="form-label">Fecha de destete <span class="text-danger">*</span></label>
                  <input type="date" class="form-control" id="desteteFecha" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Lechones destetados <span class="text-danger">*</span></label>
                  <input type="number" class="form-control" id="desteteCantidad" min="0" value="0" required>
                  <small class="text-muted" id="desteteMaxInfo"></small>
                </div>

                <hr>
                <h6 class="fw-semibold mb-3">Asignacion a lote</h6>

                <div class="mb-3">
                  <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="desteteAsignacion" id="asignarLoteExistente" value="existente" checked>
                    <label class="form-check-label" for="asignarLoteExistente">Lote existente</label>
                  </div>
                  <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="desteteAsignacion" id="asignarLoteNuevo" value="nuevo">
                    <label class="form-check-label" for="asignarLoteNuevo">Crear lote nuevo</label>
                  </div>
                </div>

                <div id="seccionLoteExistente">
                  <div class="mb-3">
                    <label class="form-label">Lote <span class="text-danger">*</span></label>
                    <select class="form-select" id="desteteLoteId"></select>
                  </div>
                </div>

                <div id="seccionLoteNuevo" style="display:none;">
                  <div class="mb-3">
                    <label class="form-label">Corral <span class="text-danger">*</span></label>
                    <select class="form-select" id="desteteCorralId"></select>
                  </div>
                  <div class="mb-3">
                    <label class="form-label">Nombre del lote <span class="text-danger">*</span></label>
                    <input type="text" class="form-control" id="desteteNombreLote" placeholder="Ej: Lote Feb-2026">
                  </div>
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarDestete"><i class="bi bi-check-lg me-1"></i>Registrar</button>
            </div>
          </div>
        </div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaDestete');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => { granjaSeleccionada = parseInt(selectGranja.value); fetchDestetes(); });

    document.getElementById('btnFiltrarDestete').addEventListener('click', fetchDestetes);
    document.getElementById('btnNuevoDestete').addEventListener('click', openNuevoDestete);
    document.getElementById('btnGuardarDestete').addEventListener('click', handleCrear);
    document.getElementById('formDestete').addEventListener('submit', (e) => { e.preventDefault(); handleCrear(); });

    document.querySelectorAll('input[name="desteteAsignacion"]').forEach(r => {
      r.addEventListener('change', toggleAsignacion);
    });

    document.getElementById('desteteCerdaId').addEventListener('change', onCerdaChange);

    await fetchDestetes();
  }

  function toggleAsignacion() {
    const tipo = document.querySelector('input[name="desteteAsignacion"]:checked').value;
    document.getElementById('seccionLoteExistente').style.display = tipo === 'existente' ? '' : 'none';
    document.getElementById('seccionLoteNuevo').style.display = tipo === 'nuevo' ? '' : 'none';
  }

  async function onCerdaChange() {
    const cerdaId = document.getElementById('desteteCerdaId').value;
    const info = document.getElementById('desteteMaxInfo');
    if (!cerdaId) { info.textContent = ''; return; }

    try {
      const res = await API.get(`/cerdas/${cerdaId}/partos`);
      const partos = res.data || [];
      if (partos.length > 0) {
        const ultimo = partos[0];
        info.textContent = `Ultimo parto: ${ultimo.lechones_nacidos_vivos} nacidos vivos`;
        document.getElementById('desteteCantidad').max = ultimo.lechones_nacidos_vivos;
        document.getElementById('desteteCantidad').value = ultimo.lechones_nacidos_vivos;
      } else {
        info.textContent = '';
      }
    } catch (e) {
      info.textContent = '';
    }
  }

  async function fetchDestetes() {
    const mes = document.getElementById('filtroMesDestete').value;
    const anio = document.getElementById('filtroAnioDestete').value;
    try {
      const data = await API.get(`/destetes?granja_id=${granjaSeleccionada}&mes=${mes}&anio=${anio}`);
      destetes = data.data || [];
      renderTable();
    } catch (err) {
      document.getElementById('destetesTableBody').innerHTML = `<div class="p-4 text-center text-danger">Error: ${err.message}</div>`;
    }
  }

  function renderTable() {
    const container = document.getElementById('destetesTableBody');
    if (destetes.length === 0) {
      container.innerHTML = `<div class="empty-state"><i class="bi bi-box-arrow-right d-block"></i><h6>No hay destetes en este periodo</h6></div>`;
      return;
    }

    const rows = destetes.map(d => `
      <tr>
        <td>${fDate(d.fecha_destete)}</td>
        <td><span class="fw-semibold">${d.cerda ? esc(d.cerda.numero_caravana) : d.cerda_id}</span></td>
        <td class="fw-semibold">${d.cantidad_lechones_destetados}</td>
        <td>${d.lote ? esc(d.lote.nombre) : '-'}</td>
        <td class="d-none d-md-table-cell">${d.lote && d.lote.corral ? esc(d.lote.corral.nombre) : '-'}</td>
        <td><button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="Destetes.editar(${d.id})"><i class="bi bi-pencil"></i></button></td>
      </tr>
    `).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr><th>Fecha</th><th>Cerda</th><th>Destetados</th><th>Lote</th><th class="d-none d-md-table-cell">Corral</th><th style="width:50px;"></th></tr></thead>
        <tbody>${rows}</tbody>
      </table>`;
  }

  async function openNuevoDestete() {
    const alert = document.getElementById('modalDesteteAlert');
    alert.classList.add('d-none');

    try {
      const [cerdasRes, corralesRes, lotesRes] = await Promise.all([
        API.get(`/granjas/${granjaSeleccionada}/cerdas?estado=cria`),
        API.get(`/granjas/${granjaSeleccionada}/corrales`),
        API.get(`/lotes?estado=activo`),
      ]);

      const cerdasCria = (cerdasRes.data || []).filter(c => c.activo);
      const corrales = corralesRes.data || [];
      const lotes = (lotesRes.data || []).filter(l => l.estado === 'activo');

      document.getElementById('desteteCerdaId').innerHTML = cerdasCria.length
        ? cerdasCria.map(c => `<option value="${c.id}">${esc(c.numero_caravana)}</option>`).join('')
        : '<option value="">No hay cerdas en cria</option>';

      document.getElementById('desteteCorralId').innerHTML = corrales.length
        ? corrales.map(c => `<option value="${c.id}">${esc(c.nombre)}</option>`).join('')
        : '<option value="">No hay corrales</option>';

      document.getElementById('desteteLoteId').innerHTML = lotes.length
        ? lotes.map(l => `<option value="${l.id}">${esc(l.nombre)} (${l.cantidad_lechones} lechones)</option>`).join('')
        : '<option value="">No hay lotes activos</option>';

      // Pre-seleccionar "nuevo lote" si no hay lotes activos
      if (lotes.length === 0) {
        document.getElementById('asignarLoteNuevo').checked = true;
        toggleAsignacion();
      } else {
        document.getElementById('asignarLoteExistente').checked = true;
        toggleAsignacion();
      }
    } catch (e) {
      alert.className = 'alert alert-danger'; alert.textContent = 'Error cargando datos: ' + e.message; alert.classList.remove('d-none');
    }

    document.getElementById('desteteFecha').value = new Date().toISOString().split('T')[0];
    document.getElementById('desteteCantidad').value = '0';
    document.getElementById('desteteMaxInfo').textContent = '';
    document.getElementById('desteteNombreLote').value = '';

    new bootstrap.Modal(document.getElementById('modalDestete')).show();

    onCerdaChange();
  }

  async function handleCrear() {
    const alert = document.getElementById('modalDesteteAlert');
    const btn = document.getElementById('btnGuardarDestete');
    const cerdaId = parseInt(document.getElementById('desteteCerdaId').value);
    const fecha = document.getElementById('desteteFecha').value;
    const cantidad = parseInt(document.getElementById('desteteCantidad').value) || 0;
    const tipo = document.querySelector('input[name="desteteAsignacion"]:checked').value;

    if (!cerdaId || !fecha) {
      alert.className = 'alert alert-warning'; alert.textContent = 'Cerda y fecha son obligatorios'; alert.classList.remove('d-none'); return;
    }
    if (cantidad <= 0) {
      alert.className = 'alert alert-warning'; alert.textContent = 'La cantidad de lechones debe ser mayor a 0'; alert.classList.remove('d-none'); return;
    }

    const body = {
      cerda_id: cerdaId,
      fecha_destete: fecha,
      cantidad_lechones_destetados: cantidad,
    };

    if (tipo === 'existente') {
      const loteId = parseInt(document.getElementById('desteteLoteId').value);
      if (!loteId) { alert.className = 'alert alert-warning'; alert.textContent = 'Selecciona un lote'; alert.classList.remove('d-none'); return; }
      body.lote_id = loteId;
    } else {
      const corralId = parseInt(document.getElementById('desteteCorralId').value);
      const nombreLote = document.getElementById('desteteNombreLote').value.trim();
      if (!corralId) { alert.className = 'alert alert-warning'; alert.textContent = 'Selecciona un corral'; alert.classList.remove('d-none'); return; }
      if (!nombreLote) { alert.className = 'alert alert-warning'; alert.textContent = 'El nombre del lote es obligatorio'; alert.classList.remove('d-none'); return; }
      body.nuevo_lote = { corral_id: corralId, nombre: nombreLote };
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Registrando...';

    try {
      await API.post('/destetes', body);
      App.showToast('Destete registrado');
      bootstrap.Modal.getInstance(document.getElementById('modalDestete')).hide();
      await fetchDestetes();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Registrar';
    }
  }

  async function editar(id) {
    editingId = id;

    let d;
    try {
      const res = await API.get(`/destetes/${id}`);
      d = res.data;
    } catch (e) { App.showToast('Error cargando destete', 'danger'); return; }

    // Crear modal de edición dinámicamente si no existe
    let modalEl = document.getElementById('modalEditarDestete');
    if (!modalEl) {
      const div = document.createElement('div');
      div.innerHTML = `
        <div class="modal fade" id="modalEditarDestete" tabindex="-1">
          <div class="modal-dialog modal-sm">
            <div class="modal-content">
              <div class="modal-header">
                <h5 class="modal-title">Editar Destete</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
              </div>
              <div class="modal-body">
                <div id="modalEditarDesteteAlert" class="alert d-none"></div>
                <div class="mb-3">
                  <label class="form-label">Fecha de destete</label>
                  <input type="date" class="form-control" id="editDesteteFecha">
                </div>
                <div class="mb-3">
                  <label class="form-label">Lechones destetados</label>
                  <input type="number" class="form-control" id="editDesteteCantidad" min="0">
                </div>
              </div>
              <div class="modal-footer">
                <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancelar</button>
                <button type="button" class="btn btn-sgp btn-sm" id="btnGuardarEditarDestete"><i class="bi bi-check-lg me-1"></i>Guardar</button>
              </div>
            </div>
          </div>
        </div>`;
      document.getElementById('contentArea').appendChild(div.firstElementChild);
      modalEl = document.getElementById('modalEditarDestete');
      document.getElementById('btnGuardarEditarDestete').addEventListener('click', handleEditar);
    }

    document.getElementById('editDesteteFecha').value = d.fecha_destete ? d.fecha_destete.split('T')[0] : '';
    document.getElementById('editDesteteCantidad').value = d.cantidad_lechones_destetados;
    document.getElementById('modalEditarDesteteAlert').classList.add('d-none');

    new bootstrap.Modal(modalEl).show();
  }

  async function handleEditar() {
    const alert = document.getElementById('modalEditarDesteteAlert');
    const btn = document.getElementById('btnGuardarEditarDestete');
    alert.classList.add('d-none');

    const fecha = document.getElementById('editDesteteFecha').value;
    const cantidad = parseInt(document.getElementById('editDesteteCantidad').value);

    if (isNaN(cantidad) || cantidad < 0) {
      alert.className = 'alert alert-warning'; alert.textContent = 'Cantidad invalida'; alert.classList.remove('d-none'); return;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';

    try {
      const body = { cantidad_lechones_destetados: cantidad };
      if (fecha) body.fecha_destete = fecha;
      await API.put(`/destetes/${editingId}`, body);
      App.showToast('Destete actualizado');
      bootstrap.Modal.getInstance(document.getElementById('modalEditarDestete')).hide();
      await fetchDestetes();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';
    }
  }

  function fDate(d) { if (!d) return '-'; try { const p = d.split('T')[0]; return new Date(p + 'T12:00:00').toLocaleDateString('es-AR'); } catch { return d; } }
  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  return { load, editar };
})();
