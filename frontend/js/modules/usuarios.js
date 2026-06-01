/**
 * SGP - Modulo Usuarios
 * Gestion de usuarios (admin / propietario)
 */

const Usuarios = (() => {
  let usuarios = [];
  let editingId = null;

  const ROL_LABELS = {
    admin: 'Administrador',
    propietario: 'Propietario',
    empleado: 'Empleado',
  };

  function rolesDisponiblesCrear() {
    if (Permissions.isAdmin()) return ['admin', 'propietario', 'empleado'];
    return ['propietario', 'empleado'];
  }

  async function load() {
    const content = document.getElementById('contentArea');

    content.innerHTML = `
      <div class="table-container">
        <div class="table-header">
          <h5><i class="bi bi-people me-2"></i>Usuarios</h5>
          <button class="btn btn-sgp" id="btnNuevoUsuario">
            <i class="bi bi-person-plus me-2"></i>Nuevo usuario
          </button>
        </div>
        <div id="usuariosTableBody">
          <div class="loading-spinner"><div class="spinner-border text-success" role="status"></div></div>
        </div>
      </div>

      <div class="modal fade" id="modalUsuario" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title" id="modalUsuarioTitle">Nuevo usuario</h5>
              <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
              <div id="modalUsuarioAlert" class="alert d-none"></div>
              <form id="formUsuario" novalidate>
                <div class="mb-3">
                  <label class="form-label">Usuario <span class="text-danger">*</span></label>
                  <input type="text" class="form-control" id="usrUsername" required minlength="3">
                </div>
                <div class="mb-3">
                  <label class="form-label">Email <span class="text-danger">*</span></label>
                  <input type="email" class="form-control" id="usrEmail" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Nombre completo</label>
                  <input type="text" class="form-control" id="usrNombre">
                </div>
                <div class="mb-3">
                  <label class="form-label">Establecimiento</label>
                  <input type="text" class="form-control" id="usrEstablecimiento">
                </div>
                <div class="mb-3">
                  <label class="form-label">Rol <span class="text-danger">*</span></label>
                  <select class="form-select" id="usrRol"></select>
                </div>
                <div class="mb-3" id="usrPasswordGroup">
                  <label class="form-label">Contrasena <span class="text-danger">*</span></label>
                  <input type="password" class="form-control" id="usrPassword" minlength="6">
                </div>
                <div class="mb-3" id="usrActivoGroup" style="display:none;">
                  <div class="form-check">
                    <input class="form-check-input" type="checkbox" id="usrActivo" checked>
                    <label class="form-check-label" for="usrActivo">Usuario activo</label>
                  </div>
                </div>
              </form>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
              <button type="button" class="btn btn-sgp" id="btnGuardarUsuario">
                <i class="bi bi-check-lg me-1"></i>Guardar
              </button>
            </div>
          </div>
        </div>
      </div>
    `;

    document.getElementById('btnNuevoUsuario').addEventListener('click', () => openModal());
    document.getElementById('btnGuardarUsuario').addEventListener('click', handleSave);
    document.getElementById('formUsuario').addEventListener('submit', (e) => {
      e.preventDefault();
      handleSave();
    });

    await fetchUsuarios();
  }

  async function fetchUsuarios() {
    try {
      const data = await API.get('/usuarios');
      usuarios = data.data || [];
      renderTable();
    } catch (err) {
      document.getElementById('usuariosTableBody').innerHTML = `
        <div class="p-4 text-center text-danger"><i class="bi bi-exclamation-triangle me-2"></i>${esc(err.message)}</div>`;
    }
  }

  function renderTable() {
    const container = document.getElementById('usuariosTableBody');
    if (usuarios.length === 0) {
      container.innerHTML = `<div class="empty-state"><i class="bi bi-people d-block"></i><h6>No hay usuarios</h6></div>`;
      return;
    }

    const me = API.getUser();
    const rows = usuarios.map(u => {
      const isSelf = me && u.id === me.id;
      return `
      <tr>
        <td><span class="fw-semibold">${esc(u.username)}</span>${isSelf ? ' <span class="badge bg-secondary">yo</span>' : ''}</td>
        <td class="d-none d-md-table-cell"><small>${esc(u.email)}</small></td>
        <td><span class="badge bg-light text-dark border">${esc(ROL_LABELS[u.rol] || u.rol)}</span></td>
        <td>${u.activo ? '<span class="badge-estado badge-activo">Activo</span>' : '<span class="badge-estado badge-cerrado">Inactivo</span>'}</td>
        <td>
          <div class="d-flex gap-1">
            <button class="btn btn-sm btn-outline-secondary" title="Editar" onclick="Usuarios.edit(${u.id})"><i class="bi bi-pencil"></i></button>
            ${!isSelf && u.activo ? `<button class="btn btn-sm btn-outline-danger" title="Desactivar" onclick="Usuarios.confirmDesactivar(${u.id})"><i class="bi bi-person-x"></i></button>` : ''}
          </div>
        </td>
      </tr>`;
    }).join('');

    container.innerHTML = `
      <table class="table table-hover mb-0">
        <thead><tr><th>Usuario</th><th class="d-none d-md-table-cell">Email</th><th>Rol</th><th>Estado</th><th style="width:90px;">Acciones</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>`;
  }

  function fillRolSelect(selected) {
    const sel = document.getElementById('usrRol');
    const roles = rolesDisponiblesCrear();
    sel.innerHTML = roles.map(r => `<option value="${r}">${ROL_LABELS[r]}</option>`).join('');
    sel.value = selected && roles.includes(selected) ? selected : roles[0];
  }

  function openModal(usuario = null) {
    editingId = usuario ? usuario.id : null;
    document.getElementById('modalUsuarioTitle').textContent = usuario ? 'Editar usuario' : 'Nuevo usuario';
    document.getElementById('usrUsername').value = usuario ? usuario.username : '';
    document.getElementById('usrUsername').disabled = !!usuario;
    document.getElementById('usrEmail').value = usuario ? usuario.email : '';
    document.getElementById('usrNombre').value = usuario ? (usuario.nombre_completo || '') : '';
    document.getElementById('usrEstablecimiento').value = usuario ? (usuario.establecimiento || '') : '';
    fillRolSelect(usuario ? usuario.rol : null);
    document.getElementById('usrPassword').value = '';
    document.getElementById('usrPasswordGroup').style.display = usuario ? 'none' : '';
    document.getElementById('usrActivoGroup').style.display = usuario ? '' : 'none';
    if (usuario) document.getElementById('usrActivo').checked = usuario.activo !== false;
    document.getElementById('modalUsuarioAlert').classList.add('d-none');
    new bootstrap.Modal(document.getElementById('modalUsuario')).show();
  }

  async function handleSave() {
    const alert = document.getElementById('modalUsuarioAlert');
    const btn = document.getElementById('btnGuardarUsuario');
    const username = document.getElementById('usrUsername').value.trim();
    const email = document.getElementById('usrEmail').value.trim();
    const nombre = document.getElementById('usrNombre').value.trim();
    const establecimiento = document.getElementById('usrEstablecimiento').value.trim();
    const rol = document.getElementById('usrRol').value;
    const password = document.getElementById('usrPassword').value;

    if (!email || (!editingId && (!username || !password))) {
      alert.className = 'alert alert-warning';
      alert.textContent = 'Completa los campos obligatorios';
      alert.classList.remove('d-none');
      return;
    }

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Guardando...';

    try {
      if (editingId) {
        const body = { email, rol };
        if (nombre) body.nombre_completo = nombre;
        if (establecimiento) body.establecimiento = establecimiento;
        body.activo = document.getElementById('usrActivo').checked;
        await API.put(`/usuarios/${editingId}`, body);
        App.showToast('Usuario actualizado');
      } else {
        const body = { username, email, password, rol };
        if (nombre) body.nombre_completo = nombre;
        if (establecimiento) body.establecimiento = establecimiento;
        await API.post('/usuarios', body);
        App.showToast('Usuario creado');
      }
      bootstrap.Modal.getInstance(document.getElementById('modalUsuario')).hide();
      await fetchUsuarios();
    } catch (err) {
      alert.className = 'alert alert-danger';
      alert.textContent = err.message;
      alert.classList.remove('d-none');
    } finally {
      btn.disabled = false;
      btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Guardar';
    }
  }

  function edit(id) {
    const u = usuarios.find(x => x.id === id);
    if (u) openModal(u);
  }

  async function confirmDesactivar(id) {
    const u = usuarios.find(x => x.id === id);
    if (!u || !confirm(`Desactivar al usuario "${u.username}"?`)) return;
    try {
      await API.del(`/usuarios/${id}`);
      App.showToast('Usuario desactivado');
      await fetchUsuarios();
    } catch (err) {
      App.showToast(err.message, 'danger');
    }
  }

  function esc(s) {
    if (!s) return '';
    const d = document.createElement('div');
    d.textContent = s;
    return d.innerHTML;
  }

  return { load, edit, confirmDesactivar };
})();
