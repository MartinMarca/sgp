/**
 * SGP - API Client
 * Wrapper sobre fetch para comunicacion con el backend
 */

const API = (() => {
  const BASE_URL = 'http://localhost:8080/api';

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

    // Si es 401, redirigir al login
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
      setUser(data.data.usuario);
    }
    return data;
  }

  async function register(input) {
    const data = await post('/auth/register', input);
    return data;
  }

  function logout() {
    removeToken();
    window.location.href = 'index.html';
  }

  // --- Public API ---

  return {
    // Auth
    login,
    register,
    logout,
    isAuthenticated,
    getUser,
    getToken,

    // HTTP
    get,
    post,
    put,
    del,
  };
})();
