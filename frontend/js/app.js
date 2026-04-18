/**
 * SGP - App Principal
 * Inicializacion, navegacion y utilidades globales
 */

const App = (() => {
  let currentPage = 'dashboard';

  // --- Inicializacion ---

  function init() {
    // Verificar autenticacion
    if (!API.isAuthenticated()) {
      window.location.href = 'index.html';
      return;
    }

    setupUser();
    setupNavigation();
    setupLogout();
    setupSidebar();

    // Cargar pagina inicial
    navigateTo('dashboard');
  }

  function setupUser() {
    const user = API.getUser();
    const nameEl = document.getElementById('sidebarUserName');
    if (user && nameEl) {
      nameEl.textContent = user.nombre_completo || user.username;
    }

    const estEl = document.getElementById('sidebarEstablecimiento');
    if (user && user.establecimiento && estEl) {
      estEl.textContent = user.establecimiento;
    }
  }

  // --- Navegacion ---

  function setupNavigation() {
    document.querySelectorAll('.nav-link[data-page]').forEach(link => {
      link.addEventListener('click', (e) => {
        e.preventDefault();
        const page = link.dataset.page;
        navigateTo(page);

        // Cerrar sidebar en mobile
        closeSidebar();
      });
    });
  }

  function navigateTo(page) {
    currentPage = page;

    // Actualizar active en sidebar
    document.querySelectorAll('.nav-link[data-page]').forEach(link => {
      link.classList.toggle('active', link.dataset.page === page);
    });

    // Actualizar titulo
    const titles = {
      dashboard: 'Inicio',
      granjas: 'Granjas',
      cerdas: 'Cerdas',
      padrillos: 'Padrillos',
      servicios: 'Servicios',
      partos: 'Partos',
      destetes: 'Destetes',
      corrales: 'Corrales',
      lotes: 'Lotes',
      muertes: 'Muertes de Animales',
      ventas: 'Ventas',
      calendario: 'Calendario',
      estadisticas: 'Estadisticas',
    };

    document.getElementById('pageTitle').textContent = titles[page] || page;

    // Cargar contenido
    loadPage(page);
  }

  async function loadPage(page) {
    const content = document.getElementById('contentArea');

    // Mostrar loading
    content.innerHTML = `
      <div class="loading-spinner">
        <div class="spinner-border text-success" role="status">
          <span class="visually-hidden">Cargando...</span>
        </div>
      </div>
    `;

    // Verificar si existe el modulo
    const moduleLoaders = {
      dashboard: loadDashboard,
      granjas: () => Granjas.load(),
      corrales: () => Corrales.load(),
      cerdas: () => Cerdas.load(),
      padrillos: () => Padrillos.load(),
      servicios: () => Servicios.load(),
      partos: () => Partos.load(),
      destetes: () => Destetes.load(),
      lotes: () => Lotes.load(),
      muertes: () => MuertesLechones.load(),
      ventas: () => Ventas.load(),
      calendario: () => Calendario.load(),
      estadisticas: () => Estadisticas.load(),
    };

    const loader = moduleLoaders[page];
    if (loader) {
      try {
        await loader();
      } catch (err) {
        content.innerHTML = `
          <div class="alert alert-danger">
            <i class="bi bi-exclamation-triangle me-2"></i>
            Error cargando la pagina: ${err.message}
          </div>
        `;
      }
    } else {
      // Pagina aun no implementada
      content.innerHTML = `
        <div class="empty-state">
          <i class="bi bi-tools d-block"></i>
          <h6>Modulo en construccion</h6>
          <p>El modulo de <strong>${titles[page] || page}</strong> se implementara proximamente.</p>
        </div>
      `;
    }
  }

  // --- Dashboard ---

  const DASH_TIPO_CFG = {
    parto_estimado: { bg: '#fef3c7', border: '#f59e0b', text: '#92400e', icon: 'bi-plus-circle-fill', label: 'Parto estimado' },
    destete_estimado: { bg: '#ede9fe', border: '#8b5cf6', text: '#5b21b6', icon: 'bi-box-arrow-right', label: 'Destete estimado' },
    confirmacion_pendiente: { bg: '#dbeafe', border: '#3b82f6', text: '#1e40af', icon: 'bi-question-circle-fill', label: 'Conf. pendiente' },
  };

  async function loadDashboard() {
    const content = document.getElementById('contentArea');

    let granjas = [];
    try {
      const data = await API.get('/granjas');
      granjas = data.data || [];
    } catch (e) {
      granjas = [];
    }

    if (granjas.length === 0) {
      content.innerHTML = renderStatCards(granjas.length, {}) + renderWelcome() + renderCtaGranja();
      document.getElementById('btnCrearPrimeraGranja').addEventListener('click', () => {
        navigateTo('granjas');
        setTimeout(() => { const btn = document.getElementById('btnNuevaGranja'); if (btn) btn.click(); }, 300);
      });
      return;
    }

    content.innerHTML = renderStatCards(granjas.length, {})
      + `<div class="row g-3 mb-4">
           <div class="col-lg-8">
             <div class="table-container h-100">
               <div class="table-header">
                 <h5><i class="bi bi-lightning-charge me-2" style="color:#f59e0b;"></i>Proximos 7 dias</h5>
                 <a href="#" class="btn btn-sm btn-sgp-outline" id="btnVerCalendario"><i class="bi bi-calendar-event me-1"></i>Ver calendario</a>
               </div>
               <div id="dashEventos" class="p-3"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
             </div>
           </div>
           <div class="col-lg-4">
             <div class="table-container h-100">
               <div class="table-header">
                 <h5><i class="bi bi-pie-chart me-2" style="color:var(--primary);"></i>Estado del rodeo</h5>
               </div>
               <div id="dashDistribucion" class="p-3"><div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div></div>
             </div>
           </div>
         </div>`;

    document.getElementById('btnVerCalendario').addEventListener('click', (e) => { e.preventDefault(); navigateTo('calendario'); });

    loadDashboardData(granjas[0].id);
  }

  async function loadDashboardData(granjaId) {
    const [statsRes, eventosRes] = await Promise.allSettled([
      API.get(`/estadisticas/granja/${granjaId}`),
      API.get(`/calendario?granja_id=${granjaId}&dias=7`),
    ]);

    if (statsRes.status === 'fulfilled') {
      const s = statsRes.value.data || {};
      document.getElementById('statCerdas').textContent = s.total_cerdas || 0;
      document.getElementById('statPadrillos').textContent = s.total_padrillos || 0;
      document.getElementById('statLechones').textContent = s.total_lechones || 0;
      renderDashDistribucion(s.cerdas_por_estado || {}, s.total_cerdas || 0);
    }

    const eventos = eventosRes.status === 'fulfilled' ? (eventosRes.value.data || []) : [];
    renderDashEventos(eventos);
  }

  function renderStatCards(cantGranjas, stats) {
    const items = [
      { id: 'statGranjas', label: 'Granjas', value: cantGranjas, icon: 'bi-building', bg: '#d1fae5', color: '#065f46' },
      { id: 'statCerdas', label: 'Cerdas activas', value: stats.total_cerdas || 0, icon: 'bi-gender-female', bg: '#fce7f3', color: '#9d174d' },
      { id: 'statPadrillos', label: 'Padrillos activos', value: stats.total_padrillos || 0, icon: 'bi-gender-male', bg: '#dbeafe', color: '#1e40af' },
      { id: 'statLechones', label: 'Lechones en lotes', value: stats.total_lechones || 0, icon: 'bi-collection', bg: '#fef3c7', color: '#92400e' },
    ];
    return `<div class="row g-3 mb-4">${items.map(c => `
      <div class="col-sm-6 col-xl-3">
        <div class="stat-card">
          <div class="d-flex align-items-center justify-content-between">
            <div>
              <div class="stat-value" id="${c.id}">${c.value}</div>
              <div class="stat-label">${c.label}</div>
            </div>
            <div class="stat-icon" style="background-color:${c.bg}; color:${c.color};"><i class="bi ${c.icon}"></i></div>
          </div>
        </div>
      </div>`).join('')}</div>`;
  }

  function renderDashEventos(eventos) {
    const container = document.getElementById('dashEventos');
    if (!eventos || eventos.length === 0) {
      container.innerHTML = `<div class="text-center py-4">
        <i class="bi bi-calendar-check d-block" style="font-size:2rem; opacity:0.3;"></i>
        <div class="text-muted mt-2">No hay eventos en los proximos 7 dias</div>
        <div class="text-muted small">Todo tranquilo por ahora</div>
      </div>`;
      return;
    }

    // Agrupar por fecha
    const porFecha = {};
    eventos.forEach(ev => {
      const fecha = ev.fecha_estimada ? ev.fecha_estimada.split('T')[0] : 'sin-fecha';
      if (!porFecha[fecha]) porFecha[fecha] = [];
      porFecha[fecha].push(ev);
    });

    const fechasOrdenadas = Object.keys(porFecha).sort();
    const hoyStr = new Date().toISOString().split('T')[0];

    let html = '<div class="dash-timeline">';
    fechasOrdenadas.forEach(fecha => {
      const d = new Date(fecha + 'T12:00:00');
      const esHoy = fecha === hoyStr;
      const diaSemana = d.toLocaleDateString('es-AR', { weekday: 'long' });
      const diaNum = d.toLocaleDateString('es-AR', { day: 'numeric', month: 'short' });
      const label = esHoy ? 'Hoy' : diaSemana.charAt(0).toUpperCase() + diaSemana.slice(1);

      html += `<div class="dash-timeline-day${esHoy ? ' dash-timeline-hoy' : ''}">
        <div class="dash-timeline-fecha">
          <span class="dash-timeline-dia">${label}</span>
          <span class="dash-timeline-num">${diaNum}</span>
        </div>
        <div class="dash-timeline-eventos">`;

      porFecha[fecha].forEach(ev => {
        const cfg = DASH_TIPO_CFG[ev.tipo] || { bg: '#e5e7eb', border: '#6b7280', text: '#374151', icon: 'bi-circle', label: ev.tipo };
        html += `<div class="dash-evento-item" style="border-left:3px solid ${cfg.border};">
          <div class="dash-evento-icono" style="background:${cfg.bg}; color:${cfg.text};"><i class="bi ${cfg.icon}"></i></div>
          <div class="dash-evento-info">
            <span class="fw-semibold" style="color:${cfg.text};">${cfg.label}</span>
            <span class="text-muted">-</span>
            <span>Cerda <strong>${escHtml(ev.cerda_caravana)}</strong></span>
          </div>
        </div>`;
      });

      html += '</div></div>';
    });
    html += '</div>';

    container.innerHTML = html;
  }

  function renderDashDistribucion(porEstado, total) {
    const container = document.getElementById('dashDistribucion');
    const cfgs = {
      disponible: { color: '#059669', label: 'Disponible' },
      servicio: { color: '#1e40af', label: 'En servicio' },
      gestacion: { color: '#92400e', label: 'Gestacion' },
      cria: { color: '#5b21b6', label: 'En cria' },
    };

    if (total === 0) {
      container.innerHTML = '<div class="text-center text-muted py-4"><i class="bi bi-gender-female d-block" style="font-size:2rem; opacity:0.3;"></i><div class="mt-2">Sin cerdas activas</div></div>';
      return;
    }

    let html = '';
    for (const [estado, cfg] of Object.entries(cfgs)) {
      const count = porEstado[estado] || 0;
      const pct = (count / total * 100).toFixed(0);
      html += `<div class="dash-estado-row">
        <div class="d-flex justify-content-between align-items-center mb-1">
          <span class="small fw-semibold">${cfg.label}</span>
          <span class="small"><strong>${count}</strong> <span class="text-muted">(${pct}%)</span></span>
        </div>
        <div class="dash-estado-bar-track">
          <div class="dash-estado-bar-fill" style="width:${pct}%; background:${cfg.color};"></div>
        </div>
      </div>`;
    }

    container.innerHTML = html;
  }

  function renderWelcome() {
    return `<div class="table-container mb-4">
      <div class="table-header"><h5><i class="bi bi-info-circle me-2"></i>Bienvenido al SGP</h5></div>
      <div class="p-4">
        <p class="mb-2">Para comenzar a usar el sistema:</p>
        <ol class="mb-0" style="font-size: 0.9rem;">
          <li class="mb-1">Crea una <strong>Granja</strong> desde el menu lateral</li>
          <li class="mb-1">Agrega <strong>Corrales</strong> a tu granja</li>
          <li class="mb-1">Registra tus <strong>Cerdas</strong> y <strong>Padrillos</strong></li>
          <li class="mb-1">Comienza a registrar <strong>Servicios</strong> para iniciar el ciclo reproductivo</li>
        </ol>
      </div>
    </div>`;
  }

  function renderCtaGranja() {
    return `<div class="table-container">
      <div class="p-4 text-center">
        <div style="font-size:3rem; color:var(--accent); margin-bottom:0.5rem;"><i class="bi bi-building-add"></i></div>
        <h5 class="mb-2">Aun no tenes granjas registradas</h5>
        <p class="text-muted mb-3" style="font-size:0.9rem;">Crea tu primera granja para empezar a gestionar tus animales y el ciclo reproductivo.</p>
        <button class="btn btn-sgp btn-lg" id="btnCrearPrimeraGranja"><i class="bi bi-plus-lg me-2"></i>Crear mi primera granja</button>
      </div>
    </div>`;
  }

  function escHtml(s) { if (!s) return ''; const d = document.createElement('div'); d.textContent = s; return d.innerHTML; }

  // --- Sidebar responsive ---

  function setupSidebar() {
    const toggle = document.getElementById('sidebarToggle');
    const overlay = document.getElementById('sidebarOverlay');

    if (toggle) {
      toggle.addEventListener('click', () => {
        document.getElementById('sidebar').classList.toggle('show');
        overlay.classList.toggle('show');
      });
    }

    if (overlay) {
      overlay.addEventListener('click', closeSidebar);
    }
  }

  function closeSidebar() {
    document.getElementById('sidebar').classList.remove('show');
    document.getElementById('sidebarOverlay').classList.remove('show');
  }

  // --- Logout ---

  function setupLogout() {
    const btn = document.getElementById('btnLogout');
    if (btn) {
      btn.addEventListener('click', () => {
        API.logout();
      });
    }
  }

  // --- Utilidades globales ---

  function showToast(message, type = 'success') {
    const container = document.getElementById('toastContainer');
    const icons = {
      success: 'bi-check-circle-fill',
      danger: 'bi-exclamation-triangle-fill',
      warning: 'bi-exclamation-circle-fill',
      info: 'bi-info-circle-fill',
    };

    const id = 'toast-' + Date.now();
    const html = `
      <div id="${id}" class="toast align-items-center text-bg-${type} border-0" role="alert">
        <div class="d-flex">
          <div class="toast-body">
            <i class="bi ${icons[type] || icons.info} me-2"></i>${message}
          </div>
          <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
      </div>
    `;

    container.insertAdjacentHTML('beforeend', html);
    const toastEl = document.getElementById(id);
    const toast = new bootstrap.Toast(toastEl, { delay: 3500 });
    toast.show();

    toastEl.addEventListener('hidden.bs.toast', () => toastEl.remove());
  }

  function formatDate(dateStr) {
    if (!dateStr) return '-';
    const datePart = typeof dateStr === 'string' ? dateStr.split('T')[0] : dateStr;
    const d = new Date(datePart + 'T12:00:00');
    return d.toLocaleDateString('es-AR', { day: '2-digit', month: '2-digit', year: 'numeric' });
  }

  function badgeEstado(estado) {
    return `<span class="badge-estado badge-${estado}">${estado}</span>`;
  }

  // --- Titulos para paginas no implementadas ---
  const titles = {
    dashboard: 'Dashboard',
    granjas: 'Granjas',
    cerdas: 'Cerdas',
    padrillos: 'Padrillos',
    servicios: 'Servicios',
    partos: 'Partos',
    destetes: 'Destetes',
    corrales: 'Corrales',
    lotes: 'Lotes',
    muertes: 'Muertes de Animales',
    ventas: 'Ventas',
    calendario: 'Calendario',
    estadisticas: 'Estadisticas',
  };

  return {
    init,
    navigateTo,
    showToast,
    formatDate,
    badgeEstado,
    get currentPage() { return currentPage; },
  };
})();

// Auto-init
document.addEventListener('DOMContentLoaded', App.init);
