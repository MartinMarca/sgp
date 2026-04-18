/**
 * SGP - Modulo Mortalidad de Lechones
 * Registro de muertes en lactancia (por parto) y en engorde (por lote)
 */

const MuertesLechones = (() => {
  let muertes = [];
  let granjas = [];
  let granjaSeleccionada = null;
  let editingId = null;

  const CAUSAS = [
    { value: 'aplastamiento', label: 'Aplastamiento' },
    { value: 'enfermedad', label: 'Enfermedad' },
    { value: 'inanicion', label: 'Inanicion' },
    { value: 'otro', label: 'Otro' },
  ];

  const CAUSA_BADGES = {
    aplastamiento: { bg: '#fef3c7', color: '#92400e' },
    enfermedad: { bg: '#fee2e2', color: '#991b1b' },
    inanicion: { bg: '#ede9fe', color: '#5b21b6' },
    otro: { bg: '#e5e7eb', color: '#374151' },
  };

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
        <select class="form-select form-select-sm" id="selectGranjaMuertes" style="max-width:250px;"></select>
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Mes:</label>
        <input type="number" class="form-control form-control-sm" id="filtroMesMuertes" min="1" max="12" value="${mesActual}" style="max-width:70px;">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Ano:</label>
        <input type="number" class="form-control form-control-sm" id="filtroAnioMuertes" min="2020" value="${anioActual}" style="max-width:90px;">
        <button class="btn btn-outline-secondary btn-sm" id="btnFiltrarMuertes"><i class="bi bi-funnel me-1"></i>Filtrar</button>
        <button class="btn btn-outline-info btn-sm ms-auto" id="btnEstadisticasMuertes"><i class="bi bi-graph-up me-1"></i>Estadisticas</button>
      </div>

      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-clipboard2-x me-2"></i>Mortalidad de animales</h5>
          <button class="btn btn-sgp" id="btnNuevaMuerte"><i class="bi bi-plus-lg me-2"></i>Registrar muerte</button>
        </div>
        <div id="muertesTableBody"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Modal Registrar Muerte -->
      <div class="modal fade" id="modalMuerte" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="modalMuerteTitle">Registrar muerte de animales</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalMuerteAlert" class="alert d-none"></div>
              <form id="formMuerte" novalidate>

                <div class="mb-3">
                  <label class="form-label fw-semibold">Etapa <span class="text-danger">*</span></label>
                  <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="muerteEtapa" id="etapaLactancia" value="lactancia" checked>
                    <label class="form-check-label" for="etapaLactancia">Lactancia (parto)</label>
                  </div>
                  <div class="form-check form-check-inline">
                    <input class="form-check-input" type="radio" name="muerteEtapa" id="etapaEngorde" value="engorde">
                    <label class="form-check-label" for="etapaEngorde">Engorde (lote)</label>
                  </div>
                </div>

                <div id="seccionLactancia">
                  <div class="mb-3">
                    <label class="form-label">Cerda en cria <span class="text-danger">*</span></label>
                    <select class="form-select" id="muerteCerdaId"></select>
                  </div>
                  <div class="mb-3" id="muertePartoInfo">
                    <small class="text-muted" id="muertePartoDetalle"></small>
                  </div>
                </div>

                <div id="seccionEngorde" style="display:none;">
                  <div class="mb-3">
                    <label class="form-label">Corral activo <span class="text-danger">*</span></label>
                    <select class="form-select" id="muerteCorralId"></select>
                    <small class="text-muted" id="muerteCorralDetalle"></small>
                  </div>
                </div>

                <div class="mb-3">
                  <label class="form-label">Fecha <span class="text-danger">*</span></label>
                  <input type="date" class="form-control" id="muerteFecha" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Cantidad <span class="text-danger">*</span></label>
                  <input type="number" class="form-control" id="muerteCantidad" min="1" value="1" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Causa <span class="text-danger">*</span></label>
                  <select class="form-select" id="muerteCausa">
                    ${CAUSAS.map(c => `<option value="${c.value}">${c.label}</option>`).join('')}
                  </select>
                </div>
                <div class="mb-3">
                  <label class="form-label">Notas</label>
                  <textarea class="form-control" id="muerteNotas" rows="2" placeholder="Observaciones (opcional)"></textarea>
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarMuerte"><i class="bi bi-check-lg me-1"></i>Registrar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Estadisticas -->
      <div class="modal fade" id="modalEstadisticasMuertes" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title"><i class="bi bi-graph-up me-2"></i>Estadisticas de mortalidad</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body" id="statsBody">
              <div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cerrar</button>
            </div>
          </div>
        </div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaMuertes');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => { granjaSeleccionada = parseInt(selectGranja.value); fetchMuertes(); });

    document.getElementById('btnFiltrarMuertes').addEventListener('click', fetchMuertes);
    document.getElementById('btnNuevaMuerte').addEventListener('click', openNuevaMuerte);
    document.getElementById('btnGuardarMuerte').addEventListener('click', handleCrear);
    document.getElementById('formMuerte').addEventListener('submit', (e) => { e.preventDefault(); handleCrear(); });
    document.getElementById('btnEstadisticasMuertes').addEventListener('click', showEstadisticas);

    document.querySelectorAll('input[name="muerteEtapa"]').forEach(r => {
      r.addEventListener('change', toggleEtapa);
    });

    document.getElementById('muerteCerdaId').addEventListener('change', onCerdaChange);
    document.getElementById('muerteCorralId').addEventListener('change', onCorralChange);

    await fetchMuertes();
  }

  function toggleEtapa() {
    const tipo = document.querySelector('input[name="muerteEtapa"]:checked').value;
    document.getElementById('seccionLactancia').style.display = tipo === 'lactancia' ? '' : 'none';
    document.getElementById('seccionEngorde').style.display = tipo === 'engorde' ? '' : 'none';
  }

  // --- Datos auxiliares para lactancia ---

  let partoActual = null;

  async function onCerdaChange() {
    const cerdaId = document.getElementById('muerteCerdaId').value;
    const info = document.getElementById('muertePartoDetalle');
    partoActual = null;
    if (!cerdaId) { info.textContent = ''; return; }

    try {
      const res = await API.get(`/cerdas/${cerdaId}/partos`);
      const partos = res.data || [];
      if (partos.length > 0) {
        partoActual = partos[0];
        // Ver muertes ya registradas para este parto
        let muertesExistentes = 0;
        try {
          const mr = await API.get(`/partos/${partoActual.id}/muertes-lechones`);
          const ml = mr.data || [];
          muertesExistentes = ml.reduce((sum, m) => sum + m.cantidad, 0);
        } catch (e) { /* ignore */ }
        const disponibles = partoActual.lechones_nacidos_vivos - muertesExistentes;
        info.textContent = `Parto: ${partoActual.lechones_nacidos_vivos} nacidos vivos, ${muertesExistentes} muertes registradas, ${disponibles} restantes`;
        document.getElementById('muerteCantidad').max = disponibles;
      } else {
        info.textContent = 'No se encontraron partos para esta cerda';
      }
    } catch (e) {
      info.textContent = '';
    }
  }

  function onCorralChange() {
    const select = document.getElementById('muerteCorralId');
    const info = document.getElementById('muerteCorralDetalle');
    const opt = select.options[select.selectedIndex];
    if (opt && opt.dataset.animales) {
      info.textContent = `Animales en el corral: ${opt.dataset.animales}`;
      document.getElementById('muerteCantidad').max = parseInt(opt.dataset.animales);
    } else {
      info.textContent = '';
      document.getElementById('muerteCantidad').removeAttribute('max');
    }
  }

  // --- Fetch y render ---

  async function fetchMuertes() {
    const mes = document.getElementById('filtroMesMuertes').value;
    const anio = document.getElementById('filtroAnioMuertes').value;
    try {
      const data = await API.get(`/muertes-lechones?granja_id=${granjaSeleccionada}&mes=${mes}&anio=${anio}`);
      muertes = data.data || [];
      renderTable();
    } catch (err) {
      document.getElementById('muertesTableBody').innerHTML = `<div class="p-4 text-center text-danger">Error: ${err.message}</div>`;
    }
  }

  function renderTable() {
    const container = document.getElementById('muertesTableBody');
    if (muertes.length === 0) {
      container.innerHTML = `<div class="empty-state"><i class="bi bi-clipboard2-x d-block"></i><h6>No hay muertes registradas en este periodo</h6><p>Registra muertes para hacer seguimiento de la mortalidad.</p></div>`;
      return;
    }

    const rows = muertes.map(m => {
      const etapa = m.parto_id ? 'Lactancia' : 'Engorde';
      const etapaBadge = m.parto_id
        ? '<span class="badge" style="background:#dbeafe;color:#1e40af;">Lactancia</span>'
        : '<span class="badge" style="background:#fef3c7;color:#92400e;">Engorde</span>';

      let referencia = '-';
      if (m.parto && m.parto.cerda) {
        referencia = `Cerda <strong>${esc(m.parto.cerda.numero_caravana)}</strong>`;
      } else if (m.corral) {
        referencia = `Corral <strong>${esc(m.corral.nombre)}</strong>`;
      }

      const causaCfg = CAUSA_BADGES[m.causa] || CAUSA_BADGES.otro;
      const causaBadge = `<span class="badge" style="background:${causaCfg.bg};color:${causaCfg.color};">${capitalize(m.causa)}</span>`;

      return `<tr>
        <td>${fDate(m.fecha)}</td>
        <td>${etapaBadge}</td>
        <td>${referencia}</td>
        <td class="fw-semibold">${m.cantidad}</td>
        <td>${causaBadge}</td>
        <td class="d-none d-md-table-cell"><small class="text-muted">${m.notas ? esc(m.notas) : '-'}</small></td>
        <td>
          <div class="d-flex gap-1">
            <button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="MuertesLechones.editar(${m.id})"><i class="bi bi-pencil"></i></button>
            <button class="btn btn-sm btn-outline-danger" title="Eliminar" onclick="MuertesLechones.eliminar(${m.id})"><i class="bi bi-trash"></i></button>
          </div>
        </td>
      </tr>`;
    }).join('');

    const totalMuertes = muertes.reduce((sum, m) => sum + m.cantidad, 0);

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr>
          <th>Fecha</th><th>Etapa</th><th>Referencia</th><th>Cantidad</th><th>Causa</th>
          <th class="d-none d-md-table-cell">Notas</th><th style="width:90px;"></th>
        </tr></thead>
        <tbody>${rows}</tbody>
        <tfoot><tr class="table-light">
          <td colspan="3" class="fw-semibold text-end">Total del periodo:</td>
          <td class="fw-bold">${totalMuertes}</td>
          <td colspan="3"></td>
        </tr></tfoot>
      </table>`;
  }

  // --- Crear ---

  async function openNuevaMuerte() {
    editingId = null;
    document.getElementById('modalMuerteTitle').textContent = 'Registrar muerte de animales';
    const alert = document.getElementById('modalMuerteAlert');
    alert.classList.add('d-none');

    // Mostrar radios de etapa
    document.querySelectorAll('input[name="muerteEtapa"]').forEach(r => r.disabled = false);
    document.getElementById('etapaLactancia').checked = true;
    toggleEtapa();

    try {
      const [cerdasRes, corralesRes] = await Promise.all([
        API.get(`/granjas/${granjaSeleccionada}/cerdas?estado=cria`),
        API.get(`/granjas/${granjaSeleccionada}/corrales?activo=true`),
      ]);

      const cerdasCria = (cerdasRes.data || []).filter(c => c.activo);
      const corralesActivos = corralesRes.data || [];

      document.getElementById('muerteCerdaId').innerHTML = cerdasCria.length
        ? '<option value="">Seleccionar cerda...</option>' + cerdasCria.map(c => `<option value="${c.id}">${esc(c.numero_caravana)}</option>`).join('')
        : '<option value="">No hay cerdas en cria</option>';

      document.getElementById('muerteCorralId').innerHTML = corralesActivos.length
        ? '<option value="">Seleccionar corral...</option>' + corralesActivos.map(c => `<option value="${c.id}" data-animales="${c.total_animales ?? 0}">${esc(c.nombre)} (${c.total_animales ?? 0} animales)</option>`).join('')
        : '<option value="">No hay corrales activos</option>';

      if (cerdasCria.length === 0 && corralesActivos.length > 0) {
        document.getElementById('etapaEngorde').checked = true;
        toggleEtapa();
      }
    } catch (e) {
      alert.className = 'alert alert-danger'; alert.textContent = 'Error cargando datos: ' + e.message; alert.classList.remove('d-none');
    }

    document.getElementById('muerteFecha').value = new Date().toISOString().split('T')[0];
    document.getElementById('muerteCantidad').value = '1';
    document.getElementById('muerteCantidad').removeAttribute('max');
    document.getElementById('muerteCausa').value = 'aplastamiento';
    document.getElementById('muerteNotas').value = '';
    document.getElementById('muertePartoDetalle').textContent = '';
    document.getElementById('muerteCorralDetalle').textContent = '';
    partoActual = null;

    const btn = document.getElementById('btnGuardarMuerte');
    btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Registrar';

    new bootstrap.Modal(document.getElementById('modalMuerte')).show();
  }

  async function handleCrear() {
    const alert = document.getElementById('modalMuerteAlert');
    const btn = document.getElementById('btnGuardarMuerte');
    alert.classList.add('d-none');

    const etapa = document.querySelector('input[name="muerteEtapa"]:checked').value;
    const fecha = document.getElementById('muerteFecha').value;
    const cantidad = parseInt(document.getElementById('muerteCantidad').value) || 0;
    const causa = document.getElementById('muerteCausa').value;
    const notas = document.getElementById('muerteNotas').value.trim();

    if (!fecha) {
      alert.className = 'alert alert-warning'; alert.textContent = 'La fecha es obligatoria'; alert.classList.remove('d-none'); return;
    }
    if (cantidad < 1) {
      alert.className = 'alert alert-warning'; alert.textContent = 'La cantidad debe ser al menos 1'; alert.classList.remove('d-none'); return;
    }

    const body = {
      granja_id: granjaSeleccionada,
      fecha: fecha,
      cantidad: cantidad,
      causa: causa,
      notas: notas,
    };

    if (editingId) {
      // Actualizar
      const updateBody = { fecha: fecha, cantidad: cantidad, causa: causa, notas: notas };
      btn.disabled = true;
      btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';
      try {
        await API.put(`/muertes-lechones/${editingId}`, updateBody);
        App.showToast('Registro actualizado');
        bootstrap.Modal.getInstance(document.getElementById('modalMuerte')).hide();
        await fetchMuertes();
      } catch (err) {
        alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
      } finally {
        btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';
      }
      return;
    }

    if (etapa === 'lactancia') {
      if (!partoActual) {
        alert.className = 'alert alert-warning'; alert.textContent = 'Selecciona una cerda con parto'; alert.classList.remove('d-none'); return;
      }
      body.parto_id = partoActual.id;
    } else {
      const corralId = parseInt(document.getElementById('muerteCorralId').value);
      if (!corralId) {
        alert.className = 'alert alert-warning'; alert.textContent = 'Selecciona un corral'; alert.classList.remove('d-none'); return;
      }
      body.corral_id = corralId;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Registrando...';

    try {
      await API.post('/muertes-lechones', body);
      App.showToast('Muerte registrada');
      bootstrap.Modal.getInstance(document.getElementById('modalMuerte')).hide();
      await fetchMuertes();
    } catch (err) {
      alert.className = 'alert alert-danger'; alert.textContent = err.message; alert.classList.remove('d-none');
    } finally {
      btn.disabled = false; btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Registrar';
    }
  }

  // --- Editar ---

  async function editar(id) {
    let m;
    try {
      const res = await API.get(`/muertes-lechones/${id}`);
      m = res.data;
    } catch (e) { App.showToast('Error cargando registro', 'danger'); return; }

    editingId = id;
    document.getElementById('modalMuerteTitle').textContent = 'Editar registro de muerte';
    const alert = document.getElementById('modalMuerteAlert');
    alert.classList.add('d-none');

    // En edicion no se puede cambiar la etapa ni el parto/corral
    const esLactancia = !!m.parto_id;
    document.getElementById(esLactancia ? 'etapaLactancia' : 'etapaEngorde').checked = true;
    document.querySelectorAll('input[name="muerteEtapa"]').forEach(r => r.disabled = true);
    toggleEtapa();

    if (esLactancia) {
      const cerda = m.parto && m.parto.cerda ? m.parto.cerda : null;
      document.getElementById('muerteCerdaId').innerHTML = cerda
        ? `<option value="${cerda.id}" selected>${esc(cerda.numero_caravana)}</option>`
        : '<option value="">-</option>';
      document.getElementById('muerteCerdaId').disabled = true;
      document.getElementById('muertePartoDetalle').textContent = cerda ? `Parto #${m.parto_id}` : '';
      partoActual = m.parto || null;
    } else {
      document.getElementById('muerteCorralId').innerHTML = m.corral
        ? `<option value="${m.corral.id}" selected>${esc(m.corral.nombre)}</option>`
        : '<option value="">-</option>';
      document.getElementById('muerteCorralId').disabled = true;
      document.getElementById('muerteCorralDetalle').textContent = '';
    }

    document.getElementById('muerteFecha').value = m.fecha ? m.fecha.split('T')[0] : '';
    document.getElementById('muerteCantidad').value = m.cantidad;
    document.getElementById('muerteCantidad').removeAttribute('max');
    document.getElementById('muerteCausa').value = m.causa;
    document.getElementById('muerteNotas').value = m.notas || '';

    const btn = document.getElementById('btnGuardarMuerte');
    btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';

    new bootstrap.Modal(document.getElementById('modalMuerte')).show();
  }

  // Resetear campos disabled al cerrar el modal
  function onModalHidden() {
    document.getElementById('muerteCerdaId').disabled = false;
    document.getElementById('muerteCorralId').disabled = false;
    document.querySelectorAll('input[name="muerteEtapa"]').forEach(r => r.disabled = false);
    editingId = null;
  }

  // --- Eliminar ---

  async function eliminar(id) {
    const m = muertes.find(x => x.id === id);
    if (!m) return;

    const desc = m.parto_id ? 'muerte en lactancia' : 'muerte en engorde';
    if (!confirm(`Eliminar este registro de ${desc} (${m.cantidad} animales)?`)) return;

    try {
      await API.del(`/muertes-lechones/${id}`);
      App.showToast('Registro eliminado');
      await fetchMuertes();
    } catch (err) {
      App.showToast('Error: ' + err.message, 'danger');
    }
  }

  // --- Estadisticas ---

  async function showEstadisticas() {
    const mes = document.getElementById('filtroMesMuertes').value;
    const anio = document.getElementById('filtroAnioMuertes').value;
    const body = document.getElementById('statsBody');
    body.innerHTML = '<div class="text-center py-4"><div class="spinner-border text-success" role="status"></div></div>';

    new bootstrap.Modal(document.getElementById('modalEstadisticasMuertes')).show();

    try {
      const res = await API.get(`/muertes-lechones/estadisticas?granja_id=${granjaSeleccionada}&mes=${mes}&anio=${anio}`);
      const s = res.data || {};
      renderEstadisticas(s, mes, anio);
    } catch (err) {
      body.innerHTML = `<div class="alert alert-danger">Error: ${err.message}</div>`;
    }
  }

  function renderEstadisticas(s, mes, anio) {
    const body = document.getElementById('statsBody');
    const periodoLabel = mes > 0 ? `${mes}/${anio}` : anio;

    const porCausa = s.muertes_por_causa || [];

    body.innerHTML = `
      <div class="mb-3 text-center">
        <small class="text-muted">Periodo: ${periodoLabel}</small>
      </div>

      <div class="row g-3 mb-4">
        <div class="col-6">
          <div class="p-3 rounded-3 text-center" style="background:#fee2e2;">
            <div class="fs-3 fw-bold" style="color:#991b1b;">${s.total_muertes || 0}</div>
            <small style="color:#991b1b;">Total muertes</small>
          </div>
        </div>
        <div class="col-6">
          <div class="p-3 rounded-3 text-center" style="background:#e5e7eb;">
            <div class="fs-3 fw-bold" style="color:#374151;">${s.total_registros || 0}</div>
            <small style="color:#374151;">Registros</small>
          </div>
        </div>
      </div>

      <div class="row g-3 mb-4">
        <div class="col-6">
          <div class="p-3 rounded-3 text-center" style="background:#dbeafe;">
            <div class="fs-4 fw-bold" style="color:#1e40af;">${s.muertes_lactancia || 0}</div>
            <small style="color:#1e40af;">En lactancia</small>
          </div>
        </div>
        <div class="col-6">
          <div class="p-3 rounded-3 text-center" style="background:#fef3c7;">
            <div class="fs-4 fw-bold" style="color:#92400e;">${s.muertes_engorde || 0}</div>
            <small style="color:#92400e;">En engorde</small>
          </div>
        </div>
      </div>

      ${porCausa.length > 0 ? `
        <h6 class="fw-semibold mb-3">Por causa</h6>
        ${porCausa.map(c => {
          const cfg = CAUSA_BADGES[c.causa] || CAUSA_BADGES.otro;
          const pct = s.total_muertes > 0 ? (c.cantidad / s.total_muertes * 100).toFixed(0) : 0;
          return `<div class="d-flex justify-content-between align-items-center mb-2">
            <div class="d-flex align-items-center gap-2">
              <span class="badge" style="background:${cfg.bg};color:${cfg.color};">${capitalize(c.causa)}</span>
            </div>
            <div>
              <strong>${c.cantidad}</strong> <small class="text-muted">(${pct}%)</small>
            </div>
          </div>`;
        }).join('')}
      ` : '<div class="text-center text-muted py-2">Sin datos de causas</div>'}
    `;
  }

  // --- Helpers ---

  function fDate(d) { if (!d) return '-'; try { const p = d.split('T')[0]; return new Date(p + 'T12:00:00').toLocaleDateString('es-AR'); } catch { return d; } }
  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }
  function capitalize(s) { return s ? s.charAt(0).toUpperCase() + s.slice(1) : ''; }

  return { load, editar, eliminar };
})();
