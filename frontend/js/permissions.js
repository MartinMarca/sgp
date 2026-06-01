/**
 * SGP - Permisos (RBAC)
 * Espejo de backend/internal/authz/permissions.go
 */

const Permissions = (() => {
  const P = {
    USERS_MANAGE: 'users:manage',
    USERS_MANAGE_ALL: 'users:manage_all',
    GRANJA_READ: 'granja:read',
    GRANJA_WRITE: 'granja:write',
    GRANJA_DELETE: 'granja:delete',
    CERDA_WRITE: 'cerda:write',
    CERDA_DELETE: 'cerda:delete',
    PADRILLO_WRITE: 'padrillo:write',
    PADRILLO_DELETE: 'padrillo:delete',
    VENTA_READ: 'venta:read',
    VENTA_WRITE: 'venta:write',
    STATS_READ: 'stats:read',
    STATS_VENTAS_READ: 'stats:ventas:read',
    RECURSO_WRITE: 'recurso:write',
  };

  let permSet = new Set();

  function load(perms) {
    permSet = new Set(perms || []);
  }

  function can(perm) {
    if (isAdmin()) return true;
    return permSet.has(perm);
  }

  function rol() {
    const user = API.getUser();
    return user ? user.rol : null;
  }

  function isAdmin() {
    return rol() === 'admin';
  }

  function isPropietario() {
    return rol() === 'propietario';
  }

  function isEmpleado() {
    return rol() === 'empleado';
  }

  function rolLabel(r) {
    const labels = { admin: 'Administrador', propietario: 'Propietario', empleado: 'Empleado' };
    return labels[r] || r || '';
  }

  /** Permiso requerido por página del menú (data-page). */
  const PAGE_PERM = {
    ventas: P.VENTA_READ,
    usuarios: P.USERS_MANAGE,
  };

  function canAccessPage(page) {
    if (isAdmin()) return true;
    const perm = PAGE_PERM[page];
    return !perm || can(perm);
  }

  return {
    P,
    load,
    can,
    rol,
    isAdmin,
    isPropietario,
    isEmpleado,
    rolLabel,
    canAccessPage,
  };
})();
