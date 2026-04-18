/**
 * SGP - Modulo Partos
 * Registro de partos para cerdas en gestacion
 */

const Partos = (() => {
  let partos = [];
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
        <select class="form-select form-select-sm" id="selectGranjaParto" style="max-width:250px;"></select>
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Mes:</label>
        <input type="number" class="form-control form-control-sm" id="filtroMesParto" min="1" max="12" value="${mesActual}" style="max-width:70px;">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Ano:</label>
        <input type="number" class="form-control form-control-sm" id="filtroAnioParto" min="2020" value="${anioActual}" style="max-width:90px;">
        <button class="btn btn-outline-secondary btn-sm" id="btnFiltrarParto"><i class="bi bi-funnel me-1"></i>Filtrar</button>
      </div>

      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-plus-circle me-2"></i>Partos</h5>
          <button class="btn btn-sgp" id="btnNuevoParto"><i class="bi bi-plus-lg me-2"></i>Registrar Parto</button>
        </div>
        <div id="partosTableBody"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Modal Crear Parto -->
      <div class="modal fade" id="modalParto" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Registrar Parto</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalPartoAlert" class="alert d-none"></div>
              <form id="formParto" novalidate>
                <div class="mb-3">
                  <label class="form-label">Cerda <span class="text-danger">*</span></label>
                  <select class="form-select" id="partoCerdaId"></select>
                  <small class="text-muted">Solo cerdas en estado "gestacion"</small>
                </div>
                <div class="mb-3">
                  <label class="form-label">Fecha de parto <span class="text-danger">*</span></label>
                  <input type="date" class="form-control" id="partoFecha" required>
                </div>
                <div class="row mb-3">
                  <div class="col-6">
                    <label class="form-label">Hembras vivas <span class="text-danger">*</span></label>
                    <input type="number" class="form-control" id="partoHembras" min="0" value="0" required>
                  </div>
                  <div class="col-6">
                    <label class="form-label">Machos vivos <span class="text-danger">*</span></label>
                    <input type="number" class="form-control" id="partoMachos" min="0" value="0" required>
                  </div>
                </div>
                <div class="row mb-3">
                  <div class="col-6">
                    <label class="form-label">Nacidos vivos</label>
                    <input type="number" class="form-control" id="partoVivos" min="0" value="0" readonly style="background:#f0f0f0;">
                    <small class="text-muted">Hembras + Machos (automatico)</small>
                  </div>
                  <div class="col-6">
                    <label class="form-label">Nacidos totales <span class="text-danger">*</span></label>
                    <input type="number" class="form-control" id="partoTotales" min="0" value="0" required>
                    <small class="text-muted">Vivos + muertos</small>
                  </div>
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarParto"><i class="bi bi-check-lg me-1"></i>Registrar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Editar Parto -->
      <div class="modal fade" id="modalEditarParto" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Editar Parto</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalEditarPartoAlert" class="alert d-none"></div>
              <form id="formEditarParto" novalidate>
                <div class="mb-3">
                  <label class="form-label">Fecha de parto</label>
                  <input type="date" class="form-control" id="editPartoFecha">
                </div>
                <div class="row mb-3">
                  <div class="col-6">
                    <label class="form-label">Hembras vivas</label>
                    <input type="number" class="form-control" id="editPartoHembras" min="0" value="0">
                  </div>
                  <div class="col-6">
                    <label class="form-label">Machos vivos</label>
                    <input type="number" class="form-control" id="editPartoMachos" min="0" value="0">
                  </div>
                </div>
                <div class="row mb-3">
                  <div class="col-6">
                    <label class="form-label">Nacidos vivos</label>
                    <input type="number" class="form-control" id="editPartoVivos" min="0" value="0" readonly style="background:#f0f0f0;">
                  </div>
                  <div class="col-6">
                    <label class="form-label">Nacidos totales</label>
                    <input type="number" class="form-control" id="editPartoTotales" min="0" value="0">
                    <small class="text-muted">Vivos + muertos</small>
                  </div>
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarEditarParto"><i class="bi bi-check-lg me-1"></i>Guardar</button>
            </div>
          </div>
        </div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaParto');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => { granjaSeleccionada = parseInt(selectGranja.value); fetchPartos(); });

    document.getElementById('btnFiltrarParto').addEventListener('click', fetchPartos);
    document.getElementById('btnNuevoParto').addEventListener('click', openNuevoParto);
    document.getElementById('btnGuardarParto').addEventListener('click', handleCrear);
    document.getElementById('formParto').addEventListener('submit', (e) => { e.preventDefault(); handleCrear(); });

    document.getElementById('partoHembras').addEventListener('input', calcVivos);
    document.getElementById('partoMachos').addEventListener('input', calcVivos);

    document.getElementById('editPartoHembras').addEventListener('input', calcEditVivos);
    document.getElementById('editPartoMachos').addEventListener('input', calcEditVivos);
    document.getElementById('btnGuardarEditarParto').addEventListener('click', handleEditar);
    document.getElementById('formEditarParto').addEventListener('submit', (e) => { e.preventDefault(); handleEditar(); });

    await fetchPartos();
  }

  async function fetchPartos() {
    const mes = document.getElementById('filtroMesParto').value;
    const anio = document.getElementById('filtroAnioParto').value;
    try {
      const data = await API.get(`/partos?granja_id=${granjaSeleccionada}&mes=${mes}&anio=${anio}`);
      partos = data.data || [];
      renderTable();
    } catch (err) {
      document.getElementById('partosTableBody').innerHTML = `<div class="p-4 text-center text-danger">Error: ${err.message}</div>`;
    }
  }

  function renderTable() {
    const container = document.getElementById('partosTableBody');
    if (partos.length === 0) {
      container.innerHTML = `<div class="empty-state"><i class="bi bi-plus-circle d-block"></i><h6>No hay partos en este periodo</h6></div>`;
      return;
    }

    const rows = partos.map(p => {
      const muertos = p.lechones_nacidos_totales - p.lechones_nacidos_vivos;
      return `<tr>
        <td>${fDate(p.fecha_parto)}</td>
        <td><span class="fw-semibold">${p.cerda ? esc(p.cerda.numero_caravana) : p.cerda_id}</span></td>
        <td class="fw-semibold">${p.lechones_nacidos_vivos}</td>
        <td>${p.lechones_nacidos_totales}${muertos > 0 ? ` <small class="text-danger">(${muertos} muertos)</small>` : ''}</td>
        <td>${p.lechones_hembras}♀ / ${p.lechones_machos}♂</td>
        <td><small class="text-muted">${fDate(p.fecha_estimada)}</small></td>
        <td><button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="Partos.editar(${p.id})"><i class="bi bi-pencil"></i></button></td>
      </tr>`;
    }).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr><th>Fecha</th><th>Cerda</th><th>Vivos</th><th>Totales</th><th>H / M</th><th>Est. destete</th><th style="width:50px;"></th></tr></thead>
        <tbody>${rows}</tbody>
      </table>`;
  }

  async function openNuevoParto() {
    const alert = document.getElementById('modalPartoAlert');
    alert.classList.add('d-none');

    try {
      const cerdasRes = await API.get(`/granjas/${granjaSeleccionada}/cerdas?estado=gestacion`);
      const cerdasGest = (cerdasRes.data || []).filter(c => c.activo);

      document.getElementById('partoCerdaId').innerHTML = cerdasGest.length
        ? cerdasGest.map(c => `<option value="${c.id}">${esc(c.numero_caravana)}</option>`).join('')
        : '<option value="">No hay cerdas en gestacion</option>';
    } catch (e) {
      alert.className = 'alert alert-danger'; alert.textContent = 'Error cargando datos: ' + e.message; alert.classList.remove('d-none');
    }

    document.getElementById('partoFecha').value = new Date().toISOString().split('T')[0];
    document.getElementById('partoVivos').value = '0';
    document.getElementById('partoTotales').value = '0';
    document.getElementById('partoHembras').value = '0';
    document.getElementById('partoMachos').value = '0';

    new bootstrap.Modal(document.getElementById('modalParto')).show();
  }

  function calcVivos() {
    const h = parseInt(document.getElementById('partoHembras').value) || 0;
    const m = parseInt(document.getElementById('partoMachos').value) || 0;
    const vivos = h + m;
    document.getElementById('partoVivos').value = vivos;
    const totalesEl = document.getElementById('partoTotales');
    if (parseInt(totalesEl.value) < vivos) totalesEl.value = vivos;
  }

  async function handleCrear() {
    const alert = document.getElementById('modalPartoAlert');
    const btn = document.getElementById('btnGuardarParto');
    const cerdaId = parseInt(document.getElementById('partoCerdaId').value);
    const fecha = document.getElementById('partoFecha').value;
    const hembras = parseInt(document.getElementById('partoHembras').value) || 0;
    const machos = parseInt(document.getElementById('partoMachos').value) || 0;
    const vivos = hembras + machos;
    const totales = parseInt(document.getElementById('partoTotales').value) || 0;

    if (!cerdaId || !fecha) {
      alert.className = 'alert alert-warning'; alert.textContent = 'Cerda y fecha son obligatorios'; alert.classList.remove('d-none'); return;
    }
    if (totales < vivos) {
      alert.className = 'alert alert-warning'; alert.textContent = 'Nacidos totales no puede ser menor que nacidos vivos (' + vivos + ')'; alert.classList.remove('d-none'); return;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Registrando...';

    try {
      await API.post('/partos', {
        cerda_id: cerdaId,
        fecha_parto: fecha,
        lechones_nacidos_vivos: vivos,
        lechones_nacidos_totales: totales,
        lechones_hembras: hembras,
        lechones_machos: machos,
      });
      App.showToast('Parto registrado');
      bootstrap.Modal.getInstance(document.getElementById('modalParto')).hide();
      await fetchPartos();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Registrar';
    }
  }

  function calcEditVivos() {
    const h = parseInt(document.getElementById('editPartoHembras').value) || 0;
    const m = parseInt(document.getElementById('editPartoMachos').value) || 0;
    const vivos = h + m;
    document.getElementById('editPartoVivos').value = vivos;
    const totalesEl = document.getElementById('editPartoTotales');
    if (parseInt(totalesEl.value) < vivos) totalesEl.value = vivos;
  }

  async function editar(id) {
    editingId = id;
    const alert = document.getElementById('modalEditarPartoAlert');
    alert.classList.add('d-none');

    let p;
    try {
      const res = await API.get(`/partos/${id}`);
      p = res.data;
    } catch (e) { App.showToast('Error cargando parto', 'danger'); return; }

    document.getElementById('editPartoFecha').value = p.fecha_parto ? p.fecha_parto.split('T')[0] : '';
    document.getElementById('editPartoHembras').value = p.lechones_hembras;
    document.getElementById('editPartoMachos').value = p.lechones_machos;
    document.getElementById('editPartoVivos').value = p.lechones_nacidos_vivos;
    document.getElementById('editPartoTotales').value = p.lechones_nacidos_totales;

    new bootstrap.Modal(document.getElementById('modalEditarParto')).show();
  }

  async function handleEditar() {
    const alert = document.getElementById('modalEditarPartoAlert');
    const btn = document.getElementById('btnGuardarEditarParto');
    alert.classList.add('d-none');

    const fecha = document.getElementById('editPartoFecha').value;
    const hembras = parseInt(document.getElementById('editPartoHembras').value) || 0;
    const machos = parseInt(document.getElementById('editPartoMachos').value) || 0;
    const vivos = hembras + machos;
    const totales = parseInt(document.getElementById('editPartoTotales').value) || 0;

    if (totales < vivos) {
      alert.className = 'alert alert-warning'; alert.textContent = 'Nacidos totales no puede ser menor que nacidos vivos (' + vivos + ')'; alert.classList.remove('d-none'); return;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';

    try {
      const body = {
        lechones_nacidos_vivos: vivos,
        lechones_nacidos_totales: totales,
        lechones_hembras: hembras,
        lechones_machos: machos,
      };
      if (fecha) body.fecha_parto = fecha;
      await API.put(`/partos/${editingId}`, body);
      App.showToast('Parto actualizado');
      bootstrap.Modal.getInstance(document.getElementById('modalEditarParto')).hide();
      await fetchPartos();
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
