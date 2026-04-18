/**
 * SGP - Modulo Calendario
 * Vista mensual de eventos futuros: partos estimados, destetes estimados
 */

const Calendario = (() => {
  let granjas = [];
  let granjaSeleccionada = null;
  let eventos = [];
  let currentYear = new Date().getFullYear();
  let currentMonth = new Date().getMonth(); // 0-based

  const TIPO_COLORES = {
    parto_estimado: { bg: '#fef3c7', border: '#f59e0b', text: '#92400e', icon: 'bi-plus-circle-fill', label: 'Parto estimado' },
    destete_estimado: { bg: '#ede9fe', border: '#8b5cf6', text: '#5b21b6', icon: 'bi-box-arrow-right', label: 'Destete estimado' },
    confirmacion_pendiente: { bg: '#dbeafe', border: '#3b82f6', text: '#1e40af', icon: 'bi-question-circle-fill', label: 'Conf. pendiente' },
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
    currentYear = new Date().getFullYear();
    currentMonth = new Date().getMonth();

    content.innerHTML = `
      <div class="calendario-wrapper">
        <div class="calendario-toolbar">
          <div class="d-flex align-items-center gap-3">
            <label class="form-label mb-0 fw-semibold" style="white-space:nowrap;">Granja:</label>
            <select class="form-select form-select-sm" id="selectGranjaCalendario" style="max-width:220px;"></select>
          </div>
          <div class="calendario-leyenda d-none d-md-flex">
            <span class="leyenda-item"><span class="leyenda-dot" style="background:#f59e0b;"></span>Parto est.</span>
            <span class="leyenda-item"><span class="leyenda-dot" style="background:#8b5cf6;"></span>Destete est.</span>
            <span class="leyenda-item"><span class="leyenda-dot" style="background:#3b82f6;"></span>Conf. pendiente</span>
          </div>
        </div>

        <div class="calendario-nav">
          <button class="btn btn-outline-secondary btn-sm" id="btnMesAnterior"><i class="bi bi-chevron-left"></i></button>
          <h4 class="calendario-mes-titulo" id="tituloMes"></h4>
          <button class="btn btn-outline-secondary btn-sm" id="btnMesSiguiente"><i class="bi bi-chevron-right"></i></button>
          <button class="btn btn-outline-secondary btn-sm ms-2" id="btnHoy">Hoy</button>
        </div>

        <div class="calendario-grid-container">
          <div class="calendario-header-row">
            <div class="calendario-header-cell">Dom</div>
            <div class="calendario-header-cell">Lun</div>
            <div class="calendario-header-cell">Mar</div>
            <div class="calendario-header-cell">Mie</div>
            <div class="calendario-header-cell">Jue</div>
            <div class="calendario-header-cell">Vie</div>
            <div class="calendario-header-cell">Sab</div>
          </div>
          <div class="calendario-body" id="calendarioBody"></div>
        </div>
      </div>

      <!-- Modal detalle dia -->
      <div class="modal fade" id="modalDiaDetalle" tabindex="-1">
        <div class="modal-dialog modal-dialog-centered">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="modalDiaTitulo">Eventos del dia</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body" id="modalDiaBody"></div>
          </div>
        </div>
      </div>
    `;

    const selectGranja = document.getElementById('selectGranjaCalendario');
    selectGranja.innerHTML = granjas.map(g => `<option value="${g.id}">${esc(g.nombre)}</option>`).join('');
    selectGranja.addEventListener('change', () => {
      granjaSeleccionada = parseInt(selectGranja.value);
      fetchEventos();
    });

    document.getElementById('btnMesAnterior').addEventListener('click', () => {
      currentMonth--;
      if (currentMonth < 0) { currentMonth = 11; currentYear--; }
      renderMes();
      fetchEventos();
    });

    document.getElementById('btnMesSiguiente').addEventListener('click', () => {
      currentMonth++;
      if (currentMonth > 11) { currentMonth = 0; currentYear++; }
      renderMes();
      fetchEventos();
    });

    document.getElementById('btnHoy').addEventListener('click', () => {
      currentYear = new Date().getFullYear();
      currentMonth = new Date().getMonth();
      renderMes();
      fetchEventos();
    });

    renderMes();
    await fetchEventos();
  }

  async function fetchEventos() {
    const diasNeeded = calcDiasParam();

    try {
      const data = await API.get(`/calendario?granja_id=${granjaSeleccionada}&dias=${diasNeeded}`);
      eventos = data.data || [];
    } catch (e) {
      eventos = [];
    }

    renderEventosEnCalendario();
  }

  function calcDiasParam() {
    const hoy = new Date();
    const inicioMes = new Date(currentYear, currentMonth, 1);
    const finMes = new Date(currentYear, currentMonth + 1, 0);

    const diffInicio = Math.abs(Math.ceil((inicioMes - hoy) / (1000 * 60 * 60 * 24)));
    const diffFin = Math.abs(Math.ceil((finMes - hoy) / (1000 * 60 * 60 * 24)));

    return Math.max(diffInicio, diffFin) + 15;
  }

  function renderMes() {
    const meses = ['Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
      'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'];

    document.getElementById('tituloMes').textContent = `${meses[currentMonth]} ${currentYear}`;

    const primerDia = new Date(currentYear, currentMonth, 1);
    const ultimoDia = new Date(currentYear, currentMonth + 1, 0);
    const diasEnMes = ultimoDia.getDate();
    const primerDiaSemana = primerDia.getDay(); // 0=Dom
    const hoy = new Date();
    const esHoyMes = hoy.getFullYear() === currentYear && hoy.getMonth() === currentMonth;

    let html = '';
    let diaActual = 1;
    const totalCeldas = Math.ceil((diasEnMes + primerDiaSemana) / 7) * 7;

    for (let i = 0; i < totalCeldas; i++) {
      if (i % 7 === 0) html += '<div class="calendario-week-row">';

      if (i < primerDiaSemana || diaActual > diasEnMes) {
        html += '<div class="calendario-cell calendario-cell-empty"></div>';
      } else {
        const esHoy = esHoyMes && diaActual === hoy.getDate();
        const dateStr = `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}-${String(diaActual).padStart(2, '0')}`;

        html += `<div class="calendario-cell${esHoy ? ' calendario-cell-hoy' : ''}" data-date="${dateStr}">
          <div class="calendario-dia-numero${esHoy ? ' hoy' : ''}">${diaActual}</div>
          <div class="calendario-eventos" id="eventos-${dateStr}"></div>
        </div>`;
        diaActual++;
      }

      if (i % 7 === 6) html += '</div>';
    }

    document.getElementById('calendarioBody').innerHTML = html;
  }

  function renderEventosEnCalendario() {
    // Limpiar eventos anteriores
    document.querySelectorAll('.calendario-eventos').forEach(el => { el.innerHTML = ''; });
    document.querySelectorAll('.calendario-cell').forEach(el => {
      el.classList.remove('tiene-eventos');
      el.onclick = null;
    });

    if (!eventos || eventos.length === 0) return;

    const eventosPorDia = {};
    eventos.forEach(ev => {
      const fecha = ev.fecha_estimada ? ev.fecha_estimada.split('T')[0] : null;
      if (!fecha) return;
      if (!eventosPorDia[fecha]) eventosPorDia[fecha] = [];
      eventosPorDia[fecha].push(ev);
    });

    Object.keys(eventosPorDia).forEach(fecha => {
      const container = document.getElementById(`eventos-${fecha}`);
      if (!container) return;

      const cell = container.closest('.calendario-cell');
      cell.classList.add('tiene-eventos');

      const eventsForDay = eventosPorDia[fecha];
      const tipos = {};
      eventsForDay.forEach(ev => {
        if (!tipos[ev.tipo]) tipos[ev.tipo] = 0;
        tipos[ev.tipo]++;
      });

      let dotsHtml = '';
      Object.keys(tipos).forEach(tipo => {
        const cfg = TIPO_COLORES[tipo] || { bg: '#e5e7eb', border: '#6b7280', text: '#374151' };
        const count = tipos[tipo];
        dotsHtml += `<span class="calendario-evento-badge" style="background:${cfg.bg}; color:${cfg.text}; border-color:${cfg.border};">${count}</span>`;
      });

      container.innerHTML = dotsHtml;

      cell.addEventListener('click', () => abrirDetalleDia(fecha, eventsForDay));
      cell.style.cursor = 'pointer';
    });
  }

  function abrirDetalleDia(fecha, eventosDelDia) {
    const d = new Date(fecha + 'T12:00:00');
    const fechaFormateada = d.toLocaleDateString('es-AR', {
      weekday: 'long', day: 'numeric', month: 'long', year: 'numeric'
    });

    document.getElementById('modalDiaTitulo').textContent = fechaFormateada.charAt(0).toUpperCase() + fechaFormateada.slice(1);

    let html = '';
    eventosDelDia.forEach(ev => {
      const cfg = TIPO_COLORES[ev.tipo] || { bg: '#e5e7eb', border: '#6b7280', text: '#374151', icon: 'bi-circle', label: ev.tipo };
      const diasTxt = ev.dias_restantes >= 0
        ? `<span class="text-muted">Faltan <strong>${ev.dias_restantes}</strong> dias</span>`
        : `<span class="text-warning">Hace <strong>${Math.abs(ev.dias_restantes)}</strong> dias</span>`;

      html += `
        <div class="calendario-evento-card" style="border-left: 4px solid ${cfg.border}; background: ${cfg.bg}20;">
          <div class="d-flex align-items-start gap-3">
            <div class="calendario-evento-icono" style="background:${cfg.bg}; color:${cfg.text};">
              <i class="bi ${cfg.icon}"></i>
            </div>
            <div class="flex-grow-1">
              <div class="fw-semibold" style="color:${cfg.text};">${cfg.label}</div>
              <div class="mt-1">
                <span class="fw-semibold"><i class="bi bi-tag me-1"></i>Caravana: ${esc(ev.cerda_caravana)}</span>
              </div>
              <div class="mt-1 small">${diasTxt}</div>
              ${ev.descripcion ? `<div class="mt-1 small text-muted">${esc(ev.descripcion)}</div>` : ''}
            </div>
          </div>
        </div>
      `;
    });

    if (eventosDelDia.length === 0) {
      html = '<div class="text-center text-muted py-3">No hay eventos para este dia.</div>';
    }

    document.getElementById('modalDiaBody').innerHTML = html;
    new bootstrap.Modal(document.getElementById('modalDiaDetalle')).show();
  }

  function esc(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  return { load };
})();
