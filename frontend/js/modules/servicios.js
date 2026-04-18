/**
 * SGP - Modulo Servicios
 * Registro de servicios, listado, confirmar/cancelar prenez
 */

const Servicios = (() => {
  let servicios = [];
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
        <select class="form-select form-select-sm" id="selectGranjaServ" style="max-width:250px;"></select>
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Mes:</label>
        <input type="number" class="form-control form-control-sm" id="filtroMesServ" min="1" max="12" value="${mesActual}" style="max-width:70px;">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Ano:</label>
        <input type="number" class="form-control form-control-sm" id="filtroAnioServ" min="2020" value="${anioActual}" style="max-width:90px;">
        <button class="btn btn-outline-secondary btn-sm" id="btnFiltrarServ"><i class="bi bi-funnel me-1"></i>Filtrar</button>
      </div>

      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-heart-pulse me-2"></i>Servicios</h5>
          <button class="btn btn-sgp" id="btnNuevoServicio"><i class="bi bi-plus-lg me-2"></i>Nuevo Servicio</button>
        </div>
        <div id="serviciosTableBody"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Pendientes de confirmacion -->
      <div class="table-container mt-4">
        <div class="table-header">
          <h5><i class="bi bi-hourglass-split me-2"></i>Pendientes de confirmacion</h5>
        </div>
        <div id="pendientesTableBody"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Modal Crear Servicio -->
      <div class="modal fade" id="modalServicio" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Registrar Servicio</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalServicioAlert" class="alert d-none"></div>
              <form id="formServicio" novalidate>
                <div class="mb-3">
                  <label class="form-label">Cerda <span class="text-danger">*</span></label>
                  <select class="form-select" id="servCerdaId"></select>
                  <small class="text-muted">Solo cerdas en estado "disponible"</small>
                </div>
                <div class="mb-3">
                  <label class="form-label">Fecha de servicio <span class="text-danger">*</span></label>
                  <input type="date" class="form-control" id="servFecha" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Tipo de monta <span class="text-danger">*</span></label>
                  <select class="form-select" id="servTipoMonta">
                    <option value="natural">Monta natural</option>
                    <option value="inseminacion">Inseminacion artificial</option>
                  </select>
                </div>
                <div id="camposNatural">
                  <div class="mb-3">
                    <label class="form-label">Padrillo <span class="text-danger">*</span></label>
                    <select class="form-select" id="servPadrilloId"></select>
                  </div>
                  <div class="mb-3">
                    <label class="form-label">Cantidad de saltos</label>
                    <input type="number" class="form-control" id="servSaltos" min="0" placeholder="Opcional">
                  </div>
                </div>
                <div id="camposInseminacion" style="display:none;">
                  <div class="mb-3">
                    <label class="form-label">Numero de pajuela <span class="text-danger">*</span></label>
                    <input type="text" class="form-control" id="servPajuela" placeholder="Ej: PAJ-2024-001">
                  </div>
                </div>
                <div class="form-check mb-3">
                  <input type="checkbox" class="form-check-input" id="servRepeticiones">
                  <label class="form-check-label" for="servRepeticiones">Tiene repeticiones</label>
                </div>
                <div class="mb-3" id="servRepeticionesGroup" style="display:none;">
                  <label class="form-label">Cantidad de repeticiones</label>
                  <input type="number" class="form-control" id="servCantRepeticiones" min="0" value="0">
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarServicio"><i class="bi bi-check-lg me-1"></i>Registrar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Confirmar Prenez -->
      <div class="modal fade" id="modalConfirmarPrenez" tabindex="-1">
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Confirmar prenez</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalConfirmarAlert" class="alert d-none"></div>
              <p>Confirmar prenez del servicio de <strong id="confirmarCerdaNombre"></strong>?</p>
              <div class="mb-3 d-none" id="confirmarFechaEstimadaRow">
                <div class="p-2 rounded" style="background:#fef9c3;border:1px solid #fde047;">
                  <small><i class="bi bi-calendar-check me-1" style="color:#92400e;"></i>
                  <strong style="color:#92400e;">Parto estimado:</strong>
                  <span id="confirmarFechaEstimada" style="color:#92400e;"></span>
                  <span class="text-muted ms-1">(fecha servicio + 114 días)</span></small>
                </div>
              </div>
              <div class="mb-3">
                <label class="form-label">Fecha de confirmacion</label>
                <input type="date" class="form-control" id="confirmarFecha">
                <small class="text-muted">Vacio = fecha de hoy</small>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-success btn-sm" id="btnConfirmarPrenez"><i class="bi bi-check-circle me-1"></i>Confirmar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Cancelar Prenez -->
      <div class="modal fade" id="modalCancelarPrenez" tabindex="-1">
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Cancelar prenez</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalCancelarAlert" class="alert d-none"></div>
              <p>Cancelar prenez del servicio de <strong id="cancelarCerdaNombre"></strong>?</p>
              <div class="mb-3">
                <label class="form-label">Motivo <span class="text-danger">*</span></label>
                <textarea class="form-control" id="cancelarMotivo" rows="2" placeholder="Motivo de cancelacion"></textarea>
              </div>
              <div class="mb-3">
                <label class="form-label">Fecha de cancelacion</label>
                <input type="date" class="form-control" id="cancelarFecha">
                <small class="text-muted">Vacio = fecha de hoy</small>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-danger btn-sm" id="btnCancelarPrenez"><i class="bi bi-x-circle me-1"></i>Cancelar prenez</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Editar Servicio -->
      <div class="modal fade" id="modalEditarServicio" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Editar Servicio</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalEditarServAlert" class="alert d-none"></div>
              <div id="editServPrenezInfo" class="alert alert-info py-2 d-none" style="font-size:0.85rem;">
                <i class="bi bi-info-circle me-1"></i>Prenez confirmada: solo se pueden editar repeticiones.
              </div>
              <form id="formEditarServicio" novalidate>
                <div class="mb-3">
                  <label class="form-label">Fecha de servicio</label>
                  <input type="date" class="form-control" id="editServFecha">
                </div>
                <div class="mb-3">
                  <label class="form-label">Tipo de monta</label>
                  <select class="form-select" id="editServTipoMonta">
                    <option value="natural">Monta natural</option>
                    <option value="inseminacion">Inseminacion artificial</option>
                  </select>
                </div>
                <div id="editCamposNatural">
                  <div class="mb-3">
                    <label class="form-label">Padrillo</label>
                    <select class="form-select" id="editServPadrilloId"></select>
                  </div>
                  <div class="mb-3">
                    <label class="form-label">Cantidad de saltos</label>
                    <input type="number" class="form-control" id="editServSaltos" min="0">
                  </div>
                </div>
                <div id="editCamposInseminacion" style="display:none;">
                  <div class="mb-3">
                    <label class="form-label">Numero de pajuela</label>
                    <input type="text" class="form-control" id="editServPajuela">
                  </div>
                </div>
                <div class="form-check mb-3">
                  <input type="checkbox" class="form-check-input" id="editServRepeticiones">
                  <label class="form-check-label" for="editServRepeticiones">Tiene repeticiones</label>
                </div>
                <div class="mb-3" id="editServRepeticionesGroup" style="display:none;">
                  <label class="form-label">Cantidad de repeticiones</label>
                  <input type="number" class="form-control" id="editServCantRepeticiones" min="0" value="0">
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarEditarServ"><i class="bi bi-check-lg me-1"></i>Guardar</button>
            </div>
          </div>
        </div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaServ');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => { granjaSeleccionada = parseInt(selectGranja.value); fetchAll(); });

    document.getElementById('btnFiltrarServ').addEventListener('click', fetchServicios);
    document.getElementById('btnNuevoServicio').addEventListener('click', openNuevoServicio);
    document.getElementById('btnGuardarServicio').addEventListener('click', handleCrear);
    document.getElementById('formServicio').addEventListener('submit', (e) => { e.preventDefault(); handleCrear(); });

    document.getElementById('servTipoMonta').addEventListener('change', toggleCamposMonta);
    document.getElementById('servRepeticiones').addEventListener('change', () => {
      document.getElementById('servRepeticionesGroup').style.display = document.getElementById('servRepeticiones').checked ? '' : 'none';
    });

    document.getElementById('editServTipoMonta').addEventListener('change', toggleEditCamposMonta);
    document.getElementById('editServRepeticiones').addEventListener('change', () => {
      document.getElementById('editServRepeticionesGroup').style.display = document.getElementById('editServRepeticiones').checked ? '' : 'none';
    });
    document.getElementById('btnGuardarEditarServ').addEventListener('click', handleEditar);
    document.getElementById('formEditarServicio').addEventListener('submit', (e) => { e.preventDefault(); handleEditar(); });

    await fetchAll();
  }

  async function fetchAll() {
    await Promise.all([fetchServicios(), fetchPendientes()]);
  }

  async function fetchServicios() {
    const mes = document.getElementById('filtroMesServ').value;
    const anio = document.getElementById('filtroAnioServ').value;
    try {
      const data = await API.get(`/servicios?granja_id=${granjaSeleccionada}&mes=${mes}&anio=${anio}`);
      servicios = data.data || [];
      renderServicios();
    } catch (err) {
      document.getElementById('serviciosTableBody').innerHTML = `<div class="p-4 text-center text-danger">Error: ${err.message}</div>`;
    }
  }

  async function fetchPendientes() {
    try {
      const data = await API.get(`/servicios/pendientes?granja_id=${granjaSeleccionada}`);
      renderPendientes(data.data || []);
    } catch (err) {
      document.getElementById('pendientesTableBody').innerHTML = `<div class="p-4 text-center text-danger">Error: ${err.message}</div>`;
    }
  }

  function renderServicios() {
    const container = document.getElementById('serviciosTableBody');
    if (servicios.length === 0) {
      container.innerHTML = `<div class="empty-state"><i class="bi bi-heart-pulse d-block"></i><h6>No hay servicios en este periodo</h6></div>`;
      return;
    }

    const rows = servicios.map(s => {
      let detalle = s.tipo_monta === 'natural' ? 'Natural' : 'IA';
      if (s.padrillo) detalle += ' — ' + esc(s.padrillo.nombre);
      else if (s.numero_pajuela) detalle += ' — ' + esc(s.numero_pajuela);

      let prenez;
      if (s.prenez_confirmada) {
        const fechaEst = s.fecha_estimada_parto ? ` <small class="text-muted">(Parto est. ${fDate(s.fecha_estimada_parto)})</small>` : '';
        prenez = `<span class="badge bg-success">Confirmada</span>${fechaEst}`;
      } else if (s.prenez_cancelada) prenez = '<span class="badge bg-danger">Cancelada</span>';
      else prenez = '<span class="badge bg-secondary">Pendiente</span>';

      return `<tr>
        <td>${fDate(s.fecha_servicio)}</td>
        <td><span class="fw-semibold">${s.cerda ? esc(s.cerda.numero_caravana) : s.cerda_id}</span></td>
        <td>${detalle}</td>
        <td>${prenez}</td>
        <td>
          <div class="d-flex gap-1">
            <button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="Servicios.editar(${s.id})"><i class="bi bi-pencil"></i></button>
            ${!s.prenez_confirmada && !s.prenez_cancelada ? `
              <button class="btn btn-sm btn-outline-success" title="Confirmar prenez" onclick="Servicios.confirmarPrenez(${s.id})"><i class="bi bi-check-circle"></i></button>
              <button class="btn btn-sm btn-outline-danger" title="Cancelar prenez" onclick="Servicios.cancelarPrenez(${s.id})"><i class="bi bi-x-circle"></i></button>
            ` : ''}
          </div>
        </td>
      </tr>`;
    }).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr><th>Fecha</th><th>Cerda</th><th>Tipo</th><th>Prenez</th><th style="width:120px;">Acciones</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>`;
  }

  function renderPendientes(pendientes) {
    const container = document.getElementById('pendientesTableBody');
    if (pendientes.length === 0) {
      container.innerHTML = `<div class="empty-state" style="padding:1.5rem;"><i class="bi bi-check-circle d-block" style="font-size:1.5rem;"></i><p class="mb-0 mt-1">No hay servicios pendientes de confirmacion</p></div>`;
      return;
    }

    const rows = pendientes.map(s => {
      const dias = Math.floor((Date.now() - new Date(s.fecha_servicio.split('T')[0] + 'T12:00:00').getTime()) / 86400000);
      return `<tr>
        <td>${fDate(s.fecha_servicio)}</td>
        <td><span class="fw-semibold">${s.cerda ? esc(s.cerda.numero_caravana) : s.cerda_id}</span></td>
        <td>${s.tipo_monta === 'natural' ? 'Natural' : 'IA'}</td>
        <td><small class="text-muted">${dias} dias</small></td>
        <td>
          <div class="d-flex gap-1">
            <button class="btn btn-sm btn-outline-success" title="Confirmar" onclick="Servicios.confirmarPrenez(${s.id})"><i class="bi bi-check-circle"></i></button>
            <button class="btn btn-sm btn-outline-danger" title="Cancelar" onclick="Servicios.cancelarPrenez(${s.id})"><i class="bi bi-x-circle"></i></button>
          </div>
        </td>
      </tr>`;
    }).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr><th>Fecha</th><th>Cerda</th><th>Tipo</th><th>Transcurridos</th><th style="width:90px;">Acciones</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>`;
  }

  function toggleCamposMonta() {
    const tipo = document.getElementById('servTipoMonta').value;
    document.getElementById('camposNatural').style.display = tipo === 'natural' ? '' : 'none';
    document.getElementById('camposInseminacion').style.display = tipo === 'inseminacion' ? '' : 'none';
  }

  async function openNuevoServicio() {
    const alert = document.getElementById('modalServicioAlert');
    alert.classList.add('d-none');

    // Cargar cerdas disponibles y padrillos activos de la granja
    try {
      const [cerdasRes, padrillosRes] = await Promise.all([
        API.get(`/granjas/${granjaSeleccionada}/cerdas?estado=disponible`),
        API.get(`/granjas/${granjaSeleccionada}/padrillos`),
      ]);
      const cerdasDisp = (cerdasRes.data || []).filter(c => c.activo);
      const padrillosAct = (padrillosRes.data || []).filter(p => p.activo);

      document.getElementById('servCerdaId').innerHTML = cerdasDisp.length
        ? cerdasDisp.map(c => `<option value="${c.id}">${esc(c.numero_caravana)}</option>`).join('')
        : '<option value="">No hay cerdas disponibles</option>';

      document.getElementById('servPadrilloId').innerHTML = padrillosAct.length
        ? padrillosAct.map(p => `<option value="${p.id}">${esc(p.nombre)} (${esc(p.numero_caravana)})</option>`).join('')
        : '<option value="">No hay padrillos activos</option>';
    } catch (e) {
      alert.className = 'alert alert-danger'; alert.textContent = 'Error cargando datos: ' + e.message; alert.classList.remove('d-none');
    }

    document.getElementById('servFecha').value = new Date().toISOString().split('T')[0];
    document.getElementById('servTipoMonta').value = 'natural';
    document.getElementById('servSaltos').value = '';
    document.getElementById('servPajuela').value = '';
    document.getElementById('servRepeticiones').checked = false;
    document.getElementById('servCantRepeticiones').value = '0';
    document.getElementById('servRepeticionesGroup').style.display = 'none';
    toggleCamposMonta();
    new bootstrap.Modal(document.getElementById('modalServicio')).show();
  }

  async function handleCrear() {
    const alert = document.getElementById('modalServicioAlert');
    const btn = document.getElementById('btnGuardarServicio');
    const cerdaId = parseInt(document.getElementById('servCerdaId').value);
    const fecha = document.getElementById('servFecha').value;
    const tipoMonta = document.getElementById('servTipoMonta').value;

    if (!cerdaId || !fecha) {
      alert.className = 'alert alert-warning'; alert.textContent = 'Cerda y fecha son obligatorios'; alert.classList.remove('d-none'); return;
    }

    const body = { cerda_id: cerdaId, fecha_servicio: fecha, tipo_monta: tipoMonta };

    if (tipoMonta === 'natural') {
      const padrilloId = parseInt(document.getElementById('servPadrilloId').value);
      if (!padrilloId) { alert.className = 'alert alert-warning'; alert.textContent = 'Selecciona un padrillo'; alert.classList.remove('d-none'); return; }
      body.padrillo_id = padrilloId;
      const saltos = document.getElementById('servSaltos').value;
      if (saltos !== '') body.cantidad_saltos = parseInt(saltos);
    } else {
      const pajuela = document.getElementById('servPajuela').value.trim();
      if (!pajuela) { alert.className = 'alert alert-warning'; alert.textContent = 'Numero de pajuela es obligatorio'; alert.classList.remove('d-none'); return; }
      body.numero_pajuela = pajuela;
    }

    if (document.getElementById('servRepeticiones').checked) {
      body.tiene_repeticiones = true;
      body.cantidad_repeticiones = parseInt(document.getElementById('servCantRepeticiones').value) || 0;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Registrando...';

    try {
      await API.post('/servicios', body);
      App.showToast('Servicio registrado');
      bootstrap.Modal.getInstance(document.getElementById('modalServicio')).hide();
      await fetchAll();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Registrar';
    }
  }

  function confirmarPrenez(servicioId) {
    const s = [...servicios].find(x => x.id === servicioId);
    document.getElementById('confirmarCerdaNombre').textContent = s && s.cerda ? s.cerda.numero_caravana : '#' + servicioId;
    document.getElementById('confirmarFecha').value = '';
    document.getElementById('modalConfirmarAlert').classList.add('d-none');

    // Calcular y mostrar fecha estimada de parto (fecha_servicio + 114 días)
    if (s && s.fecha_servicio) {
      const fechaServicio = new Date(s.fecha_servicio.split('T')[0] + 'T12:00:00');
      fechaServicio.setDate(fechaServicio.getDate() + 114);
      const fechaEst = fechaServicio.toLocaleDateString('es-AR');
      document.getElementById('confirmarFechaEstimada').textContent = fechaEst;
      document.getElementById('confirmarFechaEstimadaRow').classList.remove('d-none');
    } else {
      document.getElementById('confirmarFechaEstimadaRow').classList.add('d-none');
    }

    new bootstrap.Modal(document.getElementById('modalConfirmarPrenez')).show();

    const btn = document.getElementById('btnConfirmarPrenez');
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener('click', async () => {
      newBtn.disabled = true;
      const al = document.getElementById('modalConfirmarAlert');
      try {
        const body = {};
        const f = document.getElementById('confirmarFecha').value;
        if (f) body.fecha_confirmacion = f;
        await API.post(`/servicios/${servicioId}/confirmar-prenez`, body);
        App.showToast('Prenez confirmada');
        bootstrap.Modal.getInstance(document.getElementById('modalConfirmarPrenez')).hide();
        await fetchAll();
      } catch (err) {
        al.className = 'alert alert-danger'; al.textContent = err.message; al.classList.remove('d-none');
      } finally { newBtn.disabled = false; }
    });
  }

  function cancelarPrenez(servicioId) {
    const s = [...servicios].find(x => x.id === servicioId);
    document.getElementById('cancelarCerdaNombre').textContent = s && s.cerda ? s.cerda.numero_caravana : '#' + servicioId;
    document.getElementById('cancelarMotivo').value = '';
    document.getElementById('cancelarFecha').value = '';
    document.getElementById('modalCancelarAlert').classList.add('d-none');
    new bootstrap.Modal(document.getElementById('modalCancelarPrenez')).show();

    const btn = document.getElementById('btnCancelarPrenez');
    const newBtn = btn.cloneNode(true);
    btn.parentNode.replaceChild(newBtn, btn);
    newBtn.addEventListener('click', async () => {
      const motivo = document.getElementById('cancelarMotivo').value.trim();
      if (!motivo) {
        const al = document.getElementById('modalCancelarAlert');
        al.className = 'alert alert-warning'; al.textContent = 'El motivo es obligatorio'; al.classList.remove('d-none'); return;
      }
      newBtn.disabled = true;
      const al = document.getElementById('modalCancelarAlert');
      try {
        const body = { motivo };
        const f = document.getElementById('cancelarFecha').value;
        if (f) body.fecha_cancelacion = f;
        await API.post(`/servicios/${servicioId}/cancelar-prenez`, body);
        App.showToast('Prenez cancelada');
        bootstrap.Modal.getInstance(document.getElementById('modalCancelarPrenez')).hide();
        await fetchAll();
      } catch (err) {
        al.className = 'alert alert-danger'; al.textContent = err.message; al.classList.remove('d-none');
      } finally { newBtn.disabled = false; }
    });
  }

  function toggleEditCamposMonta() {
    const tipo = document.getElementById('editServTipoMonta').value;
    document.getElementById('editCamposNatural').style.display = tipo === 'natural' ? '' : 'none';
    document.getElementById('editCamposInseminacion').style.display = tipo === 'inseminacion' ? '' : 'none';
  }

  async function editar(id) {
    editingId = id;
    const alert = document.getElementById('modalEditarServAlert');
    alert.classList.add('d-none');

    let s;
    try {
      const res = await API.get(`/servicios/${id}`);
      s = res.data;
    } catch (e) { App.showToast('Error cargando servicio', 'danger'); return; }

    const prenezConfirmada = s.prenez_confirmada;
    const infoEl = document.getElementById('editServPrenezInfo');
    infoEl.classList.toggle('d-none', !prenezConfirmada);

    document.getElementById('editServFecha').value = s.fecha_servicio ? s.fecha_servicio.split('T')[0] : '';
    document.getElementById('editServFecha').disabled = prenezConfirmada;
    document.getElementById('editServTipoMonta').value = s.tipo_monta;
    document.getElementById('editServTipoMonta').disabled = prenezConfirmada;

    // Cargar padrillos
    try {
      const padrillosRes = await API.get(`/granjas/${granjaSeleccionada}/padrillos`);
      const padrillosAct = (padrillosRes.data || []).filter(p => p.activo);
      document.getElementById('editServPadrilloId').innerHTML = padrillosAct.length
        ? padrillosAct.map(p => `<option value="${p.id}">${esc(p.nombre)} (${esc(p.numero_caravana)})</option>`).join('')
        : '<option value="">No hay padrillos</option>';
    } catch (e) {}

    if (s.padrillo_id) document.getElementById('editServPadrilloId').value = s.padrillo_id;
    document.getElementById('editServPadrilloId').disabled = prenezConfirmada;
    document.getElementById('editServSaltos').value = s.cantidad_saltos != null ? s.cantidad_saltos : '';
    document.getElementById('editServSaltos').disabled = prenezConfirmada;
    document.getElementById('editServPajuela').value = s.numero_pajuela || '';
    document.getElementById('editServPajuela').disabled = prenezConfirmada;

    document.getElementById('editServRepeticiones').checked = s.tiene_repeticiones;
    document.getElementById('editServCantRepeticiones').value = s.cantidad_repeticiones || 0;
    document.getElementById('editServRepeticionesGroup').style.display = s.tiene_repeticiones ? '' : 'none';

    toggleEditCamposMonta();
    new bootstrap.Modal(document.getElementById('modalEditarServicio')).show();
  }

  async function handleEditar() {
    const alert = document.getElementById('modalEditarServAlert');
    const btn = document.getElementById('btnGuardarEditarServ');
    alert.classList.add('d-none');

    const body = {};
    const fecha = document.getElementById('editServFecha').value;
    if (fecha && !document.getElementById('editServFecha').disabled) body.fecha_servicio = fecha;

    const tipoMonta = document.getElementById('editServTipoMonta').value;
    if (!document.getElementById('editServTipoMonta').disabled) {
      body.tipo_monta = tipoMonta;
      if (tipoMonta === 'natural') {
        const pid = parseInt(document.getElementById('editServPadrilloId').value);
        if (pid) body.padrillo_id = pid;
        const saltos = document.getElementById('editServSaltos').value;
        if (saltos !== '') body.cantidad_saltos = parseInt(saltos);
      } else {
        const paj = document.getElementById('editServPajuela').value.trim();
        if (paj) body.numero_pajuela = paj;
      }
    }

    body.tiene_repeticiones = document.getElementById('editServRepeticiones').checked;
    body.cantidad_repeticiones = parseInt(document.getElementById('editServCantRepeticiones').value) || 0;

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';

    try {
      await API.put(`/servicios/${editingId}`, body);
      App.showToast('Servicio actualizado');
      bootstrap.Modal.getInstance(document.getElementById('modalEditarServicio')).hide();
      await fetchAll();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';
    }
  }

  function fDate(d) { if (!d) return '-'; try { const p = d.split('T')[0]; return new Date(p + 'T12:00:00').toLocaleDateString('es-AR'); } catch { return d; } }
  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  return { load, confirmarPrenez, cancelarPrenez, editar };
})();
