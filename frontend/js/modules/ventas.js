/**
 * SGP - Modulo Ventas
 * Registro de ventas de animales (cerdas, padrillos, lechones)
 */

const Ventas = (() => {
  let granjas = [];
  let granjaSeleccionada = null;
  let ventas = [];
  let lotes = [];
  let corrales = [];
  let editingId = null;

  const TIPOS = {
    lechon:   { label: 'Lechón',   color: '#92400e', bg: '#fef3c7' },
    capon:    { label: 'Capón',    color: '#065f46', bg: '#d1fae5' },
    cerda:    { label: 'Cerda',    color: '#9d174d', bg: '#fce7f3' },
    padrillo: { label: 'Padrillo', color: '#1e40af', bg: '#dbeafe' },
  };

  async function load() {
    const content = document.getElementById('contentArea');

    try {
      const data = await API.get('/granjas');
      granjas = data.data || [];
    } catch (e) { granjas = []; }

    if (granjas.length === 0) {
      content.innerHTML = `<div class="empty-state"><i class="bi bi-building d-block"></i><h6>No hay granjas</h6>
        <button class="btn btn-sgp" onclick="App.navigateTo('granjas')"><i class="bi bi-plus-lg me-2"></i>Crear Granja</button></div>`;
      return;
    }

    granjaSeleccionada = granjas[0].id;
    const hoy = new Date();
    const mesActual = hoy.getMonth() + 1;
    const anioActual = hoy.getFullYear();

    content.innerHTML = `
      <div class="d-flex align-items-center gap-3 mb-3 flex-wrap">
        <select class="form-select form-select-sm" id="selectGranjaVentas" style="max-width:250px;"></select>
        <input type="number" class="form-control form-control-sm" id="filtroMesVentas" min="1" max="12" value="${mesActual}" style="max-width:70px;">
        <input type="number" class="form-control form-control-sm" id="filtroAnioVentas" min="2020" value="${anioActual}" style="max-width:90px;">
        <button class="btn btn-outline-secondary btn-sm" id="btnFiltrarVentas"><i class="bi bi-funnel me-1"></i>Filtrar</button>
        <button class="btn btn-outline-info btn-sm ms-auto" id="btnEstadisticasVentas"><i class="bi bi-graph-up me-1"></i>Estadísticas</button>
      </div>

      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-cart-check me-2" style="color:#059669;"></i>Ventas</h5>
          <button class="btn btn-sgp btn-sm" id="btnNuevaVenta"><i class="bi bi-plus-lg me-2"></i>Registrar venta</button>
        </div>
        <div id="ventasTableBody"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Modal Registrar / Editar Venta -->
      <div class="modal fade" id="modalVenta" tabindex="-1">
        <div class="modal-dialog modal-dialog-centered">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="modalVentaTitle">Registrar venta</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalVentaAlert" class="alert d-none"></div>
              <form id="formVenta" novalidate>
                <div class="row g-3">
                  <div class="col-md-6">
                    <label class="form-label">Tipo de animal <span class="text-danger">*</span></label>
                    <select class="form-select" id="ventaTipo">
                      <option value="lechon">Lechón</option>
                      <option value="capon">Capón</option>
                      <option value="cerda">Cerda</option>
                      <option value="padrillo">Padrillo</option>
                    </select>
                  </div>
                  <div class="col-md-6">
                    <label class="form-label">Fecha <span class="text-danger">*</span></label>
                    <input type="date" class="form-control" id="ventaFecha" required>
                  </div>
                  <div class="col-md-6">
                    <label class="form-label">Cantidad <span class="text-danger">*</span></label>
                    <input type="number" class="form-control" id="ventaCantidad" min="1" value="1" required>
                  </div>
                  <div class="col-md-6">
                    <label class="form-label">Comprador <span class="text-danger">*</span></label>
                    <input type="text" class="form-control" id="ventaComprador" placeholder="Nombre o razón social" required>
                  </div>
                  <div class="col-md-6">
                    <label class="form-label">KG totales <span class="text-danger">*</span></label>
                    <input type="number" class="form-control" id="ventaKg" min="0" step="0.1" value="0" required>
                  </div>
                  <div class="col-md-6">
                    <label class="form-label">Monto ($) <span class="text-danger">*</span></label>
                    <input type="number" class="form-control" id="ventaMonto" min="0" step="0.01" value="0" required>
                  </div>

                  <!-- Referencia (lote/corral) solo para lechones -->
                  <div class="col-12" id="ventaReferenciaRow">
                    <label class="form-label">Referencia (opcional)</label>
                    <div class="d-flex gap-2">
                      <select class="form-select form-select-sm" id="ventaRefTipo" style="max-width:130px;">
                        <option value="">Ninguna</option>
                        <option value="lote">Lote</option>
                        <option value="corral">Corral</option>
                      </select>
                      <select class="form-select form-select-sm d-none" id="ventaLoteId"></select>
                      <select class="form-select form-select-sm d-none" id="ventaCorralId"></select>
                    </div>
                    <small class="text-muted" id="ventaLoteInfo"></small>
                  </div>

                  <div class="col-12">
                    <label class="form-label">Notas</label>
                    <textarea class="form-control" id="ventaNotas" rows="2" placeholder="Observaciones (opcional)"></textarea>
                  </div>
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-outline-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarVenta"><i class="bi bi-check-lg me-1"></i>Guardar</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Estadísticas -->
      <div class="modal fade" id="modalEstadisticasVentas" tabindex="-1">
        <div class="modal-dialog modal-dialog-centered modal-lg">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">Estadísticas de Ventas</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body" id="modalEstadisticasVentasBody">
              <div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div>
            </div>
          </div>
        </div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaVentas');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => { granjaSeleccionada = parseInt(selectGranja.value); fetchVentas(); });

    document.getElementById('btnFiltrarVentas').addEventListener('click', fetchVentas);
    document.getElementById('btnNuevaVenta').addEventListener('click', openNuevaVenta);
    document.getElementById('btnGuardarVenta').addEventListener('click', handleGuardar);
    document.getElementById('formVenta').addEventListener('submit', e => { e.preventDefault(); handleGuardar(); });
    document.getElementById('btnEstadisticasVentas').addEventListener('click', showEstadisticas);
    document.getElementById('ventaTipo').addEventListener('change', onTipoChange);
    document.getElementById('ventaRefTipo').addEventListener('change', onRefTipoChange);
    document.getElementById('ventaLoteId').addEventListener('change', onLoteChange);

    await fetchVentas();
  }

  async function fetchVentas() {
    const mes = document.getElementById('filtroMesVentas').value;
    const anio = document.getElementById('filtroAnioVentas').value;
    try {
      const data = await API.get(`/ventas?granja_id=${granjaSeleccionada}&mes=${mes}&anio=${anio}`);
      ventas = data.data || [];
      renderTabla();
    } catch (err) {
      document.getElementById('ventasTableBody').innerHTML =
        `<div class="p-4 text-center text-danger">Error: ${esc(err.message)}</div>`;
    }
  }

  function renderTabla() {
    const container = document.getElementById('ventasTableBody');
    if (ventas.length === 0) {
      container.innerHTML = `<div class="empty-state"><i class="bi bi-cart-x d-block"></i>
        <h6>No hay ventas registradas en este periodo</h6></div>`;
      return;
    }

    let totalAnimales = 0, totalKg = 0, totalMonto = 0;
    const rows = ventas.map(v => {
      totalAnimales += v.cantidad;
      totalKg += parseFloat(v.kg_totales) || 0;
      totalMonto += parseFloat(v.monto) || 0;
      const cfg = TIPOS[v.tipo_animal] || { label: v.tipo_animal, color: '#374151', bg: '#f3f4f6' };
      const kgProm = v.cantidad > 0 ? (v.kg_totales / v.cantidad).toFixed(1) : '—';
      const ref = v.lote ? `<span class="badge bg-warning text-dark">Lote: ${esc(v.lote.nombre)}</span>`
        : v.corral ? `<span class="badge bg-secondary">Corral: ${esc(v.corral.nombre)}</span>` : '—';
      return `<tr>
        <td>${fDate(v.fecha)}</td>
        <td><span class="badge" style="background:${cfg.bg}; color:${cfg.color};">${cfg.label}</span></td>
        <td class="text-center">${v.cantidad}</td>
        <td>${esc(v.comprador)}</td>
        <td class="text-end">${parseFloat(v.kg_totales).toFixed(1)} kg</td>
        <td class="text-center text-muted small">${kgProm} kg</td>
        <td class="text-end fw-semibold">$${formatMonto(v.monto)}</td>
        <td>${ref}</td>
        <td>
          <button class="btn btn-sm btn-outline-secondary" onclick="Ventas.editar(${v.id})"><i class="bi bi-pencil"></i></button>
          <button class="btn btn-sm btn-outline-danger ms-1" onclick="Ventas.eliminar(${v.id})"><i class="bi bi-trash"></i></button>
        </td>
      </tr>`;
    }).join('');

    container.innerHTML = `
      <div class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead class="table-light">
            <tr>
              <th>Fecha</th><th>Tipo</th><th class="text-center">Cant.</th>
              <th>Comprador</th><th class="text-end">KG Total</th><th class="text-center">KG Prom.</th>
              <th class="text-end">Monto</th><th>Referencia</th><th></th>
            </tr>
          </thead>
          <tbody>${rows}</tbody>
          <tfoot class="table-light fw-bold">
            <tr>
              <td colspan="2">Total</td>
              <td class="text-center">${totalAnimales}</td>
              <td></td>
              <td class="text-end">${totalKg.toFixed(1)} kg</td>
              <td></td>
              <td class="text-end">$${formatMonto(totalMonto)}</td>
              <td colspan="2"></td>
            </tr>
          </tfoot>
        </table>
      </div>`;
  }

  async function openNuevaVenta() {
    editingId = null;
    document.getElementById('modalVentaTitle').textContent = 'Registrar venta';
    resetAlert();
    document.getElementById('ventaTipo').value = 'lechon';
    document.getElementById('ventaFecha').value = new Date().toISOString().split('T')[0];
    document.getElementById('ventaCantidad').value = '1';
    document.getElementById('ventaComprador').value = '';
    document.getElementById('ventaKg').value = '0';
    document.getElementById('ventaMonto').value = '0';
    document.getElementById('ventaNotas').value = '';
    document.getElementById('ventaRefTipo').value = '';
    document.getElementById('ventaLoteInfo').textContent = '';
    await cargarReferencias();
    onTipoChange();
    onRefTipoChange();
    new bootstrap.Modal(document.getElementById('modalVenta')).show();
  }

  async function cargarReferencias() {
    try {
      const [lotesRes, corralesRes] = await Promise.all([
        API.get(`/granjas/${granjaSeleccionada}/lotes`),
        API.get(`/granjas/${granjaSeleccionada}/corrales`),
      ]);
      lotes = (lotesRes.data || []).filter(l => l.estado === 'activo');
      corrales = corralesRes.data || [];
    } catch { lotes = []; corrales = []; }

    document.getElementById('ventaLoteId').innerHTML = lotes.length
      ? lotes.map(l => `<option value="${l.id}" data-cantidad="${l.cantidad_lechones}">${esc(l.nombre)} (${l.cantidad_lechones} lech.)</option>`).join('')
      : '<option value="">Sin lotes activos</option>';

    document.getElementById('ventaCorralId').innerHTML = corrales.length
      ? corrales.map(c => `<option value="${c.id}">${esc(c.nombre)}</option>`).join('')
      : '<option value="">Sin corrales</option>';
  }

  function onTipoChange() {
    const tipo = document.getElementById('ventaTipo').value;
    const refRow = document.getElementById('ventaReferenciaRow');
    refRow.style.display = (tipo === 'lechon' || tipo === 'capon') ? '' : 'none';
    if (tipo !== 'lechon' && tipo !== 'capon') {
      document.getElementById('ventaRefTipo').value = '';
      onRefTipoChange();
    }
  }

  function onRefTipoChange() {
    const refTipo = document.getElementById('ventaRefTipo').value;
    document.getElementById('ventaLoteId').classList.toggle('d-none', refTipo !== 'lote');
    document.getElementById('ventaCorralId').classList.toggle('d-none', refTipo !== 'corral');
    document.getElementById('ventaLoteInfo').textContent = '';
    if (refTipo === 'lote') onLoteChange();
  }

  function onLoteChange() {
    const sel = document.getElementById('ventaLoteId');
    const opt = sel.options[sel.selectedIndex];
    const cantidad = opt ? parseInt(opt.dataset.cantidad) || 0 : 0;
    document.getElementById('ventaLoteInfo').textContent = cantidad > 0
      ? `Disponibles: ${cantidad} lechones` : '';
    document.getElementById('ventaCantidad').max = cantidad > 0 ? cantidad : '';
  }

  async function handleGuardar() {
    const alert = document.getElementById('modalVentaAlert');
    const btn = document.getElementById('btnGuardarVenta');

    const fecha = document.getElementById('ventaFecha').value;
    const tipo = document.getElementById('ventaTipo').value;
    const cantidad = parseInt(document.getElementById('ventaCantidad').value) || 0;
    const comprador = document.getElementById('ventaComprador').value.trim();
    const kg = parseFloat(document.getElementById('ventaKg').value) || 0;
    const monto = parseFloat(document.getElementById('ventaMonto').value) || 0;
    const notas = document.getElementById('ventaNotas').value.trim();
    const refTipo = document.getElementById('ventaRefTipo').value;

    if (!fecha || !comprador || cantidad < 1) {
      showAlert(alert, 'Completá los campos obligatorios (fecha, comprador, cantidad).', 'danger');
      return;
    }

    const body = {
      granja_id: granjaSeleccionada,
      fecha,
      tipo_animal: tipo,
      cantidad,
      kg_totales: kg,
      monto,
      comprador,
      notas,
    };

    if (tipo === 'lechon' || tipo === 'capon') {
      if (refTipo === 'lote') {
        const loteId = parseInt(document.getElementById('ventaLoteId').value);
        if (loteId) body.lote_id = loteId;
      } else if (refTipo === 'corral') {
        const corralId = parseInt(document.getElementById('ventaCorralId').value);
        if (corralId) body.corral_id = corralId;
      }
    }

    btn.disabled = true;
    try {
      if (editingId) {
        await API.put(`/ventas/${editingId}`, body);
        App.showToast('Venta actualizada');
      } else {
        await API.post('/ventas', body);
        App.showToast('Venta registrada');
      }
      bootstrap.Modal.getInstance(document.getElementById('modalVenta')).hide();
      await fetchVentas();
    } catch (err) {
      showAlert(alert, err.message || 'Error al guardar.', 'danger');
    } finally {
      btn.disabled = false;
    }
  }

  async function editar(id) {
    editingId = id;
    const v = ventas.find(x => x.id === id);
    if (!v) return;

    document.getElementById('modalVentaTitle').textContent = 'Editar venta';
    resetAlert();
    await cargarReferencias();

    document.getElementById('ventaTipo').value = v.tipo_animal;
    document.getElementById('ventaFecha').value = v.fecha ? v.fecha.split('T')[0] : '';
    document.getElementById('ventaCantidad').value = v.cantidad;
    document.getElementById('ventaCantidad').removeAttribute('max');
    document.getElementById('ventaComprador').value = v.comprador;
    document.getElementById('ventaKg').value = v.kg_totales;
    document.getElementById('ventaMonto').value = v.monto;
    document.getElementById('ventaNotas').value = v.notas || '';

    if (v.tipo_animal === 'lechon' || v.tipo_animal === 'capon') {
      if (v.lote_id) {
        document.getElementById('ventaRefTipo').value = 'lote';
        document.getElementById('ventaLoteId').value = v.lote_id;
      } else if (v.corral_id) {
        document.getElementById('ventaRefTipo').value = 'corral';
        document.getElementById('ventaCorralId').value = v.corral_id;
      } else {
        document.getElementById('ventaRefTipo').value = '';
      }
    } else {
      document.getElementById('ventaRefTipo').value = '';
    }

    onTipoChange();
    onRefTipoChange();
    new bootstrap.Modal(document.getElementById('modalVenta')).show();
  }

  async function eliminar(id) {
    const v = ventas.find(x => x.id === id);
    if (!v) return;
    const cfg = TIPOS[v.tipo_animal] || { label: v.tipo_animal };
    if (!confirm(`¿Eliminar la venta de ${v.cantidad} ${cfg.label.toLowerCase()}(s) a "${v.comprador}"?`)) return;
    try {
      await API.del(`/ventas/${id}`);
      App.showToast('Venta eliminada');
      await fetchVentas();
    } catch (err) {
      App.showToast(err.message || 'Error al eliminar', 'danger');
    }
  }

  async function showEstadisticas() {
    const mes = document.getElementById('filtroMesVentas').value;
    const anio = document.getElementById('filtroAnioVentas').value;
    const body = document.getElementById('modalEstadisticasVentasBody');
    body.innerHTML = '<div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div>';
    new bootstrap.Modal(document.getElementById('modalEstadisticasVentas')).show();

    try {
      const res = await API.get(`/ventas/estadisticas?granja_id=${granjaSeleccionada}&mes=${mes}&anio=${anio}`);
      const s = res.data || {};
      const meses = ['','Enero','Febrero','Marzo','Abril','Mayo','Junio','Julio','Agosto','Septiembre','Octubre','Noviembre','Diciembre'];
      const periodo = `${meses[parseInt(mes)] || ''} ${anio}`;
      const porTipo = s.por_tipo || [];

      let tipoRows = porTipo.map(t => {
        const cfg = TIPOS[t.tipo_animal] || { label: t.tipo_animal, color: '#374151', bg: '#f3f4f6' };
        const kgProm = t.cantidad > 0 ? (t.kg_totales / t.cantidad).toFixed(1) : '—';
        return `<tr>
          <td><span class="badge" style="background:${cfg.bg}; color:${cfg.color};">${cfg.label}</span></td>
          <td class="text-center">${t.cantidad}</td>
          <td class="text-end">${parseFloat(t.kg_totales).toFixed(1)} kg</td>
          <td class="text-center">${kgProm} kg</td>
          <td class="text-end fw-semibold">$${formatMonto(t.monto)}</td>
        </tr>`;
      }).join('');

      if (!tipoRows) tipoRows = '<tr><td colspan="5" class="text-center text-muted">Sin ventas en el periodo</td></tr>';

      body.innerHTML = `
        <p class="text-muted small mb-3">Periodo: <strong>${esc(periodo)}</strong></p>
        <div class="row g-3 mb-4">
          <div class="col-6 col-md-3">
            <div class="stat-card"><div class="stat-value">${s.total_ventas || 0}</div><div class="stat-label">Registros</div></div>
          </div>
          <div class="col-6 col-md-3">
            <div class="stat-card"><div class="stat-value">${s.total_animales || 0}</div><div class="stat-label">Animales</div></div>
          </div>
          <div class="col-6 col-md-3">
            <div class="stat-card"><div class="stat-value">${parseFloat(s.total_kg || 0).toFixed(1)}</div><div class="stat-label">KG totales</div></div>
          </div>
          <div class="col-6 col-md-3">
            <div class="stat-card"><div class="stat-value">$${formatMonto(s.total_monto || 0)}</div><div class="stat-label">Monto total</div></div>
          </div>
        </div>
        <h6 class="fw-semibold mb-2">Desglose por tipo</h6>
        <div class="table-responsive">
          <table class="table table-sm align-middle">
            <thead class="table-light">
              <tr><th>Tipo</th><th class="text-center">Cant.</th><th class="text-end">KG Total</th><th class="text-center">KG Prom.</th><th class="text-end">Monto</th></tr>
            </thead>
            <tbody>${tipoRows}</tbody>
          </table>
        </div>`;
    } catch (err) {
      body.innerHTML = `<div class="alert alert-danger">Error: ${esc(err.message)}</div>`;
    }
  }

  // --- Helpers ---

  function resetAlert() {
    const a = document.getElementById('modalVentaAlert');
    a.classList.add('d-none');
    a.textContent = '';
  }

  function showAlert(el, msg, type) {
    el.className = `alert alert-${type}`;
    el.textContent = msg;
  }

  function formatMonto(v) {
    return parseFloat(v).toLocaleString('es-AR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  }

  function fDate(d) {
    if (!d) return '-';
    try { const p = d.split('T')[0]; return new Date(p + 'T12:00:00').toLocaleDateString('es-AR'); } catch { return d; }
  }

  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  return { load, editar, eliminar };
})();
