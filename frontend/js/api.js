/**
 * SGP - API Client
 * Wrapper sobre fetch para comunicacion con el backend
 */

const API = (() => {
  const BASE_URL = 'http://localhost:8080/api';
  const GRANJA_KEY = 'sgp_granja_activa';

  let session = {
    granjas: [],
    permisos: [],
    registerEnabled: true,
  };

  // --- Token management ---

  function getToken() {
    return localStorage.getItem('sgp_token');
  }

  function setToken(token) {
    localStorage.setItem('sgp_token', token);
  }

  function removeToken() {
    localStorage.removeItem('sgp_token');
    localStorage.removeItem('sgp_user');
    localStorage.removeItem(GRANJA_KEY);
    session = { granjas: [], permisos: [], registerEnabled: true };
    Permissions.load([]);
  }

  function getUser() {
    const data = localStorage.getItem('sgp_user');
    return data ? JSON.parse(data) : null;
  }

  function setUser(user) {
    localStorage.setItem('sgp_user', JSON.stringify(user));
  }

  function isAuthenticated() {
    return !!getToken();
  }

  // --- Sesion / auth/me ---

  async function refreshSession() {
    if (!isAuthenticated()) return null;
    const data = await get('/auth/me');
    const me = data.data || {};
    if (me.usuario) setUser(me.usuario);
    session.granjas = me.granjas || [];
    session.permisos = me.permisos || [];
    session.registerEnabled = me.register_enabled !== false;
    Permissions.load(session.permisos);
    ensureGranjaActiva();
    return me;
  }

  async function fetchPublicConfig() {
    try {
      const response = await fetch(`${BASE_URL}/auth/config`);
      const data = await response.json();
      if (data.data && typeof data.data.register_enabled === 'boolean') {
        session.registerEnabled = data.data.register_enabled;
      }
    } catch (_) {
      /* mantener default */
    }
    return session.registerEnabled;
  }

  function isRegisterEnabled() {
    return session.registerEnabled;
  }

  function getGranjas() {
    return session.granjas.slice();
  }

  function getGranjaActivaId() {
    const stored = localStorage.getItem(GRANJA_KEY);
    if (stored) {
      const id = parseInt(stored, 10);
      if (session.granjas.some(g => g.id === id)) return id;
    }
    return session.granjas[0] ? session.granjas[0].id : null;
  }

  function setGranjaActiva(id) {
    localStorage.setItem(GRANJA_KEY, String(id));
    window.dispatchEvent(new CustomEvent('sgp:granja-changed', { detail: { id: parseInt(id, 10) } }));
  }

  function ensureGranjaActiva() {
    if (session.granjas.length === 0) return;
    const current = getGranjaActivaId();
    if (current) return;
    setGranjaActiva(session.granjas[0].id);
  }

  function getGranjaActiva() {
    const id = getGranjaActivaId();
    return session.granjas.find(g => g.id === id) || null;
  }

  // --- HTTP helpers ---

  async function request(method, path, body = null) {
    const headers = {
      'Content-Type': 'application/json',
    };

    const token = getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const options = { method, headers };
    if (body && method !== 'GET') {
      options.body = JSON.stringify(body);
    }

    const response = await fetch(`${BASE_URL}${path}`, options);

    if (response.status === 401) {
      removeToken();
      if (!window.location.pathname.endsWith('index.html') && window.location.pathname !== '/') {
        window.location.href = 'index.html';
      }
      throw new Error('Sesion expirada');
    }

    const data = await response.json();

    if (!response.ok) {
      const errorMsg = data.error || data.message || `Error ${response.status}`;
      throw new Error(errorMsg);
    }

    return data;
  }

  function get(path) {
    return request('GET', path);
  }

  function post(path, body) {
    return request('POST', path, body);
  }

  function put(path, body) {
    return request('PUT', path, body);
  }

  function del(path) {
    return request('DELETE', path);
  }

  // --- Auth ---

  async function login(username, password) {
    const data = await post('/auth/login', { username, password });
    if (data.data && data.data.token) {
      setToken(data.data.token);
      if (data.data.usuario) setUser(data.data.usuario);
      await refreshSession();
    }
    return data;
  }

  async function register(input) {
    const data = await post('/auth/register', input);
    return data;
  }

  function logout() {
    removeToken();
    window.location.replace('index.html');
  }

  return {
    login,
    register,
    logout,
    isAuthenticated,
    getUser,
    getToken,
    refreshSession,
    fetchPublicConfig,
    isRegisterEnabled,
    getGranjas,
    getGranjaActivaId,
    setGranjaActiva,
    getGranjaActiva,
    get,
    post,
    put,
    del,
  };
})();
