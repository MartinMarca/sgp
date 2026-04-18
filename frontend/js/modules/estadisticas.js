/**
 * SGP - Modulo Estadisticas
 * Resumen general de la granja y metricas del periodo
 */

const Estadisticas = (() => {
  let granjas = [];
  let granjaSeleccionada = null;

  const ESTADOS_CONFIG = {
    disponible: { color: '#059669', bg: '#d1fae5', label: 'Disponible', icon: 'bi-check-circle' },
    servicio:   { color: '#1e40af', bg: '#dbeafe', label: 'En servicio', icon: 'bi-heart-pulse' },
    gestacion:  { color: '#92400e', bg: '#fef3c7', label: 'En gestacion', icon: 'bi-hourglass-split' },
    cria:       { color: '#5b21b6', bg: '#ede9fe', label: 'En cria', icon: 'bi-egg' },
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

    content.innerHTML = `
      <div class="d-flex align-items-center gap-3 mb-4 flex-wrap">
        <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Granja:</label>
        <select class="form-select form-select-sm" id="selectGranjaStats" style="max-width:250px;"></select>
      </div>

      <!-- Resumen General -->
      <div class="stats-section-title"><i class="bi bi-bar-chart me-2"></i>Resumen General</div>
      <div class="row g-3 mb-4" id="resumenCards">
        <div class="col-12"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>

      <!-- Distribucion de Cerdas -->
      <div class="stats-section-title"><i class="bi bi-pie-chart me-2"></i>Distribucion de Cerdas por Estado</div>
      <div class="table-container mb-4">
        <div class="p-3" id="distribucionCerdas">
          <div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div>
        </div>
      </div>

      <!-- Estadisticas del Periodo -->
      <div class="d-flex align-items-center gap-3 mb-3 flex-wrap">
        <div class="stats-section-title mb-0"><i class="bi bi-calendar3 me-2"></i>Estadisticas del Periodo</div>
        <div class="d-flex align-items-center gap-2 ms-auto">
          <label class="form-label mb-0 fw-semibold small" style="white-space:nowrap;">Mes:</label>
          <input type="number" class="form-control form-control-sm" id="statsMes" min="1" max="12" value="${hoy.getMonth() + 1}" style="max-width:70px;">
          <label class="form-label mb-0 fw-semibold small" style="white-space:nowrap;">Ano:</label>
          <input type="number" class="form-control form-control-sm" id="statsAnio" min="2020" value="${hoy.getFullYear()}" style="max-width:90px;">
          <button class="btn btn-outline-secondary btn-sm" id="btnFiltrarStats"><i class="bi bi-funnel me-1"></i>Filtrar</button>
        </div>
      </div>
      <div class="row g-3 mb-4" id="periodoStats">
        <div class="col-12"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaStats');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => {
      granjaSeleccionada = parseInt(selectGranja.value);
      fetchAll();
    });

    document.getElementById('btnFiltrarStats').addEventListener('click', fetchPeriodo);

    await fetchAll();
  }

  async function fetchAll() {
    await Promise.all([fetchResumen(), fetchPeriodo()]);
  }

  async function fetchResumen() {
    try {
      const data = await API.get(`/estadisticas/granja/${granjaSeleccionada}`);
      const r = data.data || {};
      renderResumen(r);
      renderDistribucion(r.cerdas_por_estado || {}, r.total_cerdas || 0);
    } catch (e) {
      document.getElementById('resumenCards').innerHTML = `<div class="col-12"><div class="alert alert-danger">Error cargando resumen: ${e.message}</div></div>`;
    }
  }

  function renderResumen(r) {
    const totalAnimales = (r.total_cerdas || 0) + (r.total_padrillos || 0) + (r.total_lechones || 0);
    const cards = [
      { label: 'Total animales', value: totalAnimales, icon: 'bi-bar-chart-fill', iconBg: '#d1fae5', iconColor: '#065f46' },
      { label: 'Cerdas activas', value: r.total_cerdas || 0, icon: 'bi-gender-female', iconBg: '#fce7f3', iconColor: '#9d174d' },
      { label: 'Padrillos activos', value: r.total_padrillos || 0, icon: 'bi-gender-male', iconBg: '#dbeafe', iconColor: '#1e40af' },
      { label: 'Lechones en lotes', value: r.total_lechones || 0, icon: 'bi-piggy-bank', iconBg: '#ede9fe', iconColor: '#5b21b6' },
      { label: 'Corrales', value: r.total_corrales || 0, icon: 'bi-grid-3x3-gap', iconBg: '#f3f4f6', iconColor: '#374151' },
      { label: 'Lotes activos', value: r.total_lotes_activos || 0, icon: 'bi-collection', iconBg: '#fef3c7', iconColor: '#92400e' },
    ];

    document.getElementById('resumenCards').innerHTML = cards.map(c => `
      <div class="col-6 col-md-4 col-xl">
        <div class="stat-card">
          <div class="d-flex align-items-center justify-content-between">
            <div>
              <div class="stat-value">${c.value}</div>
              <div class="stat-label">${c.label}</div>
            </div>
            <div class="stat-icon" style="background-color:${c.iconBg}; color:${c.iconColor};">
              <i class="bi ${c.icon}"></i>
            </div>
          </div>
        </div>
      </div>
    `).join('');
  }

  function renderDistribucion(porEstado, total) {
    const container = document.getElementById('distribucionCerdas');

    if (total === 0) {
      container.innerHTML = '<div class="text-center text-muted py-3">No hay cerdas activas en esta granja.</div>';
      return;
    }

    let html = '<div class="stats-bar-container">';

    // Barra visual stacked
    html += '<div class="stats-stacked-bar">';
    for (const [estado, cfg] of Object.entries(ESTADOS_CONFIG)) {
      const count = porEstado[estado] || 0;
      const pct = total > 0 ? (count / total * 100) : 0;
      if (pct > 0) {
        html += `<div class="stats-bar-segment" style="width:${pct}%; background:${cfg.color};" title="${cfg.label}: ${count} (${pct.toFixed(1)}%)"></div>`;
      }
    }
    html += '</div>';

    // Desglose detallado
    html += '<div class="stats-estado-grid">';
    for (const [estado, cfg] of Object.entries(ESTADOS_CONFIG)) {
      const count = porEstado[estado] || 0;
      const pct = total > 0 ? (count / total * 100) : 0;
      html += `
        <div class="stats-estado-item">
          <div class="stats-estado-dot" style="background:${cfg.color};"></div>
          <div class="stats-estado-info">
            <div class="stats-estado-label">${cfg.label}</div>
            <div class="stats-estado-values">
              <span class="fw-bold" style="font-size:1.1rem;">${count}</span>
              <span class="text-muted small">(${pct.toFixed(1)}%)</span>
            </div>
          </div>
        </div>
      `;
    }
    html += '</div></div>';

    container.innerHTML = html;
  }

  async function fetchPeriodo() {
    const mes = document.getElementById('statsMes').value;
    const anio = document.getElementById('statsAnio').value;
    const qs = `granja_id=${granjaSeleccionada}&mes=${mes}&anio=${anio}`;

    try {
      const [periodoRes, muertesRes, ventasRes] = await Promise.all([
        API.get(`/estadisticas/periodo?${qs}`),
        API.get(`/muertes-lechones/estadisticas?${qs}`),
        API.get(`/ventas/estadisticas?${qs}`),
      ]);
      renderPeriodo(
        periodoRes.data || {},
        muertesRes.data || {},
        ventasRes.data  || {},
        mes, anio
      );
    } catch (e) {
      document.getElementById('periodoStats').innerHTML = `<div class="col-12"><div class="alert alert-danger">Error: ${e.message}</div></div>`;
    }
  }

  const REF_PRENEZ = [
    { label: 'Excelente',            min: 90,  max: 95,  color: '#059669', bg: '#d1fae5' },
    { label: 'Buena',                min: 85,  max: 90,  color: '#0284c7', bg: '#e0f2fe' },
    { label: 'Regular',              min: 80,  max: 85,  color: '#d97706', bg: '#fef3c7' },
    { label: 'Problema reproductivo',min: 0,   max: 80,  color: '#dc2626', bg: '#fee2e2' },
  ];

  const REF_PARTO = [
    { label: 'Excelente',  min: 88, max: 92, color: '#059669', bg: '#d1fae5' },
    { label: 'Buena',      min: 85, max: 88, color: '#0284c7', bg: '#e0f2fe' },
    { label: 'Aceptable',  min: 80, max: 85, color: '#d97706', bg: '#fef3c7' },
    { label: 'Problema',   min: 0,  max: 80, color: '#dc2626', bg: '#fee2e2' },
  ];

  function clasificarTasa(tasa, referencias) {
    for (const ref of referencias) {
      if (tasa >= ref.min) return ref;
    }
    return referencias[referencias.length - 1];
  }

  function rangoLabel(r) {
    if (r.min === 0) return '&lt;' + r.max + '%';
    return r.min + '\u2013' + r.max + '%';
  }

  function renderTasaCard(titulo, icono, iconoColor, tasa, totalServicios, numerador, referencias) {
    const sinDatos = totalServicios === 0;
    const clasificacion = sinDatos ? null : clasificarTasa(tasa, referencias);

    let valorHtml;
    if (sinDatos) {
      valorHtml = '<div class="tasa-valor" style="color:#9ca3af;">\u2014</div>'
        + '<div class="tasa-sin-datos">Sin servicios en el periodo</div>';
    } else {
      valorHtml = '<div class="tasa-valor" style="color:' + clasificacion.color + ';">' + tasa.toFixed(1) + '%</div>'
        + '<div class="tasa-badge" style="background:' + clasificacion.bg + '; color:' + clasificacion.color + ';">'
        + esc(clasificacion.label) + '</div>'
        + '<div class="tasa-detalle">' + numerador + ' / ' + totalServicios + ' servicios</div>';
    }

    let refRows = '';
    for (let i = 0; i < referencias.length; i++) {
      const r = referencias[i];
      const esActual = !sinDatos && clasificacion === r;
      refRows += '<tr' + (esActual ? ' class="tasa-ref-activa"' : '') + '>'
        + '<td><span class="tasa-ref-dot" style="background:' + r.color + ';"></span>' + esc(r.label) + '</td>'
        + '<td>' + rangoLabel(r) + '</td>'
        + '</tr>';
    }

    return '<div class="col-md-6">'
      + '<div class="table-container h-100">'
      + '<div class="table-header"><h5><i class="bi ' + icono + ' me-2" style="color:' + iconoColor + ';"></i>' + esc(titulo) + '</h5></div>'
      + '<div class="p-3">'
      + '<div class="tasa-display">' + valorHtml + '</div>'
      + '<div class="tasa-ref-titulo">Valores de referencia</div>'
      + '<table class="tasa-ref-table"><tbody>' + refRows + '</tbody></table>'
      + '</div></div></div>';
  }

  function renderPeriodo(s, muertes, ventas, mes, anio) {
    const partos   = s.partos    || {};
    const destetes = s.destetes  || {};
    const srvStats = s.servicios || {};

    const mesesNombres = ['', 'Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
      'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'];
    const periodoLabel = mesesNombres[parseInt(mes)] + ' ' + anio;

    const totalServicios      = parseInt(srvStats.total_servicios)      || 0;
    const totalConfirmaciones = parseInt(srvStats.total_confirmaciones)  || 0;
    const totalPartos         = parseInt(partos.total_partos)            || 0;

    const tasaPrenez = totalServicios > 0 ? (totalConfirmaciones / totalServicios * 100) : 0;
    const tasaParto  = totalServicios > 0 ? (totalPartos         / totalServicios * 100) : 0;

    const cardPrenez = renderTasaCard('Tasa de Pre\u00f1ez', 'bi-heart-pulse', '#ec4899',
      tasaPrenez, totalServicios, totalConfirmaciones + ' confirmaciones', REF_PRENEZ);
    const cardParto = renderTasaCard('Tasa de Parto', 'bi-plus-circle', '#f59e0b',
      tasaParto, totalServicios, totalPartos + ' partos', REF_PARTO);

    // --- Muertes ---
    const totalMuertes     = parseInt(muertes.total_muertes)     || 0;
    const muertesLactancia = parseInt(muertes.muertes_lactancia) || 0;
    const muertesEngorde   = parseInt(muertes.muertes_engorde)   || 0;
    const porCausa         = muertes.muertes_por_causa           || [];

    const causaLabels = { aplastamiento: 'Aplastamiento', enfermedad: 'Enfermedad', inanicion: 'Inanición', otro: 'Otro' };
    let causaRows = porCausa.map(c => {
      const pct = totalMuertes > 0 ? (c.cantidad / totalMuertes * 100).toFixed(0) : 0;
      return renderMetricRow((causaLabels[c.causa] || c.causa) + ' (' + pct + '%)', c.cantidad, '#dc2626');
    }).join('');
    if (!causaRows) causaRows = '<div class="text-center text-muted py-2 small">Sin muertes en este periodo</div>';

    const cardMuertes = '<div class="col-md-6">'
      + '<div class="table-container h-100">'
      + '<div class="table-header"><h5><i class="bi bi-clipboard2-pulse me-2" style="color:#dc2626;"></i>Muertes de Animales</h5></div>'
      + '<div class="p-3">'
      + renderMetricRow('Total muertes', totalMuertes, '#dc2626')
      + renderMetricRow('En lactancia', muertesLactancia, '#f87171')
      + renderMetricRow('En engorde', muertesEngorde, '#fca5a5')
      + '<div style="border-top:1px solid #f3f4f6; margin:0.5rem 0;"></div>'
      + causaRows
      + '</div></div></div>';

    // --- Ventas ---
    const totalAnimalesVendidos = parseInt(ventas.total_animales) || 0;
    const totalKgVendidos       = parseFloat(ventas.total_kg)     || 0;
    const totalMontoVentas      = parseFloat(ventas.total_monto)  || 0;
    const ventasPorTipo         = ventas.por_tipo                 || [];
    const kgProm = totalAnimalesVendidos > 0 ? (totalKgVendidos / totalAnimalesVendidos).toFixed(1) : '—';

    const tipoLabels = { lechon: 'Lechones', cerda: 'Cerdas', padrillo: 'Padrillos' };
    let tipoRows = ventasPorTipo.map(t =>
      renderMetricRow((tipoLabels[t.tipo_animal] || t.tipo_animal) + ' — $' + formatMonto(t.monto), t.cantidad, '#059669')
    ).join('');
    if (!tipoRows) tipoRows = '<div class="text-center text-muted py-2 small">Sin ventas en este periodo</div>';

    const cardVentas = '<div class="col-md-6">'
      + '<div class="table-container h-100">'
      + '<div class="table-header"><h5><i class="bi bi-cart-check me-2" style="color:#059669;"></i>Ventas</h5></div>'
      + '<div class="p-3">'
      + renderMetricRow('Animales vendidos', totalAnimalesVendidos, '#059669')
      + renderMetricRow('KG totales', totalKgVendidos.toFixed(1) + ' kg', '#0284c7')
      + renderMetricRow('KG promedio/animal', kgProm + ' kg', '#0284c7')
      + renderMetricRow('Monto total', '$' + formatMonto(totalMontoVentas), '#7c3aed')
      + '<div style="border-top:1px solid #f3f4f6; margin:0.5rem 0;"></div>'
      + tipoRows
      + '</div></div></div>';

    const html = cardPrenez + cardParto
      + '<div class="col-md-6"><div class="table-container h-100"><div class="table-header">'
      + '<h5><i class="bi bi-plus-circle me-2" style="color:#f59e0b;"></i>Partos</h5></div><div class="p-3">'
      + renderMetricRow('Total de partos', partos.total_partos || 0, '#f59e0b')
      + renderMetricRow('Lechones nacidos vivos', partos.total_lechones_nacidos || 0, '#059669')
      + renderMetricRow('Promedio vivos/parto', formatDec(partos.promedio_lechones_vivos), '#3b82f6')
      + renderMetricRow('Promedio totales/parto', formatDec(partos.promedio_lechones_totales), '#6366f1')
      + (partos.total_partos === 0 ? '<div class="text-center text-muted py-2 small">Sin partos en este periodo</div>' : '')
      + '</div></div></div>'
      + '<div class="col-md-6"><div class="table-container h-100"><div class="table-header">'
      + '<h5><i class="bi bi-box-arrow-right me-2" style="color:#8b5cf6;"></i>Destetes</h5></div><div class="p-3">'
      + renderMetricRow('Total de destetes', destetes.total_destetes || 0, '#8b5cf6')
      + renderMetricRow('Lechones destetados', destetes.total_lechones_destetados || 0, '#059669')
      + renderMetricRow('Promedio lech./destete', formatDec(destetes.promedio_lechones_destetados), '#3b82f6')
      + (destetes.total_destetes === 0 ? '<div class="text-center text-muted py-2 small">Sin destetes en este periodo</div>' : '')
      + '</div></div></div>'
      + cardMuertes
      + cardVentas;

    document.getElementById('periodoStats').innerHTML = html;
  }

  function formatMonto(v) {
    return parseFloat(v || 0).toLocaleString('es-AR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  }

  function renderMetricRow(label, value, color) {
    return `
      <div class="stats-metric-row">
        <div class="stats-metric-indicator" style="background:${color};"></div>
        <span class="stats-metric-label">${label}</span>
        <span class="stats-metric-value">${value}</span>
      </div>
    `;
  }

  function formatDec(val) {
    if (val === undefined || val === null) return '0';
    return parseFloat(val).toFixed(1);
  }

  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  return { load };
})();
