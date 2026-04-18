/**
 * SGP - Auth Module
 * Maneja login, registro y redirecciones
 */

const Auth = (() => {
  let currentForm = 'login'; // 'login' | 'register'

  function init() {
    // Si ya esta autenticado, ir al app
    if (API.isAuthenticated()) {
      window.location.href = 'app.html';
      return;
    }

    bindEvents();
  }

  function bindEvents() {
    // Login form
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
      loginForm.addEventListener('submit', handleLogin);
    }

    // Register form
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
      registerForm.addEventListener('submit', handleRegister);
    }

    // Toggle between login/register
    const toggleLinks = document.querySelectorAll('[data-toggle-form]');
    toggleLinks.forEach(link => {
      link.addEventListener('click', (e) => {
        e.preventDefault();
        toggleForm(link.dataset.toggleForm);
      });
    });
  }

  function toggleForm(form) {
    const loginSection = document.getElementById('loginSection');
    const registerSection = document.getElementById('registerSection');

    if (form === 'register') {
      loginSection.classList.add('d-none');
      registerSection.classList.remove('d-none');
      currentForm = 'register';
    } else {
      registerSection.classList.add('d-none');
      loginSection.classList.remove('d-none');
      currentForm = 'login';
    }

    clearAlerts();
  }

  async function handleLogin(e) {
    e.preventDefault();
    clearAlerts();

    const username = document.getElementById('loginUsername').value.trim();
    const password = document.getElementById('loginPassword').value;
    const btn = e.target.querySelector('button[type="submit"]');

    if (!username || !password) {
      showAlert('loginAlert', 'Completa todos los campos', 'warning');
      return;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Ingresando...';

    try {
      await API.login(username, password);
      window.location.href = 'app.html';
    } catch (err) {
      showAlert('loginAlert', err.message || 'Error al iniciar sesion', 'danger');
    } finally {
      btn.disabled = false;
      btn.innerHTML = 'Ingresar';
    }
  }

  async function handleRegister(e) {
    e.preventDefault();
    clearAlerts();

    const username = document.getElementById('regUsername').value.trim();
    const email = document.getElementById('regEmail').value.trim();
    const nombre = document.getElementById('regNombre').value.trim();
    const establecimiento = document.getElementById('regEstablecimiento').value.trim();
    const password = document.getElementById('regPassword').value;
    const password2 = document.getElementById('regPassword2').value;
    const btn = e.target.querySelector('button[type="submit"]');

    if (!username || !email || !password) {
      showAlert('registerAlert', 'Completa los campos obligatorios', 'warning');
      return;
    }

    if (password.length < 6) {
      showAlert('registerAlert', 'La contrasena debe tener al menos 6 caracteres', 'warning');
      return;
    }

    if (password !== password2) {
      showAlert('registerAlert', 'Las contrasenas no coinciden', 'warning');
      return;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Registrando...';

    try {
      await API.register({
        username,
        email,
        password,
        nombre_completo: nombre,
        establecimiento: establecimiento || undefined,
      });
      showAlert('registerAlert', 'Registro exitoso. Ahora podes iniciar sesion.', 'success');
      setTimeout(() => toggleForm('login'), 1500);
    } catch (err) {
      showAlert('registerAlert', err.message || 'Error al registrar', 'danger');
    } finally {
      btn.disabled = false;
      btn.innerHTML = 'Crear cuenta';
    }
  }

  function showAlert(id, message, type) {
    const el = document.getElementById(id);
    if (el) {
      el.className = `alert alert-${type}`;
      el.textContent = message;
      el.classList.remove('d-none');
    }
  }

  function clearAlerts() {
    document.querySelectorAll('.auth-alert').forEach(el => {
      el.classList.add('d-none');
      el.textContent = '';
    });
  }

  return { init };
})();

// Auto-init cuando el DOM carga
document.addEventListener('DOMContentLoaded', Auth.init);
