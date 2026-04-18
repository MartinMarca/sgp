// Cliente API para comunicación con el backend
const API_BASE_URL = 'http://localhost:3000/api';

// Función auxiliar para obtener el token de autenticación
function getAuthToken() {
  return localStorage.getItem('authToken');
}

// Función auxiliar para hacer peticiones HTTP
async function apiRequest(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  const token = getAuthToken();

  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` })
    }
  };

  const config = {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...(options.headers || {})
    }
  };

  try {
    const response = await fetch(url, config);
    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error?.message || 'Error en la petición');
    }

    return data;
  } catch (error) {
    console.error('Error en API:', error);
    throw error;
  }
}

// API de Autenticación
export const authAPI = {
  login: async (username, password) => {
    return apiRequest('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    });
  },
  logout: () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
  },
  getCurrentUser: async () => {
    return apiRequest('/auth/me');
  }
};

// API de Cerdas
export const cerdasAPI = {
  getAll: async (estado = null) => {
    const params = estado ? `?estado=${estado}` : '';
    return apiRequest(`/cerdas${params}`);
  },
  getById: async (id) => {
    return apiRequest(`/cerdas/${id}`);
  },
  create: async (cerdaData) => {
    return apiRequest('/cerdas', {
      method: 'POST',
      body: JSON.stringify(cerdaData)
    });
  },
  update: async (id, cerdaData) => {
    return apiRequest(`/cerdas/${id}`, {
      method: 'PUT',
      body: JSON.stringify(cerdaData)
    });
  },
  delete: async (id) => {
    return apiRequest(`/cerdas/${id}`, {
      method: 'DELETE'
    });
  },
  confirmarPrenez: async (servicioId) => {
    return apiRequest(`/cerdas/servicios/${servicioId}/confirmar-prenez`, {
      method: 'POST'
    });
  },
  cancelarPrenez: async (servicioId, motivo) => {
    return apiRequest(`/cerdas/servicios/${servicioId}/cancelar-prenez`, {
      method: 'POST',
      body: JSON.stringify({ motivo })
    });
  },
  getEstadisticas: async (cerdaId) => {
    return apiRequest(`/cerdas/${cerdaId}/estadisticas`);
  }
};

// API de Padrillos
export const padrillosAPI = {
  getAll: async () => {
    return apiRequest('/padrillos');
  },
  getById: async (id) => {
    return apiRequest(`/padrillos/${id}`);
  },
  create: async (padrilloData) => {
    return apiRequest('/padrillos', {
      method: 'POST',
      body: JSON.stringify(padrilloData)
    });
  },
  update: async (id, padrilloData) => {
    return apiRequest(`/padrillos/${id}`, {
      method: 'PUT',
      body: JSON.stringify(padrilloData)
    });
  },
  delete: async (id) => {
    return apiRequest(`/padrillos/${id}`, {
      method: 'DELETE'
    });
  }
};

// API de Servicios
export const serviciosAPI = {
  getAll: async (mes = null, año = null) => {
    const params = new URLSearchParams();
    if (mes) params.append('mes', mes);
    if (año) params.append('año', año);
    const query = params.toString() ? `?${params.toString()}` : '';
    return apiRequest(`/servicios${query}`);
  },
  getById: async (id) => {
    return apiRequest(`/servicios/${id}`);
  },
  create: async (servicioData) => {
    return apiRequest('/servicios', {
      method: 'POST',
      body: JSON.stringify(servicioData)
    });
  },
  update: async (id, servicioData) => {
    return apiRequest(`/servicios/${id}`, {
      method: 'PUT',
      body: JSON.stringify(servicioData)
    });
  }
};

// API de Partos
export const partosAPI = {
  getAll: async (mes = null, año = null) => {
    const params = new URLSearchParams();
    if (mes) params.append('mes', mes);
    if (año) params.append('año', año);
    const query = params.toString() ? `?${params.toString()}` : '';
    return apiRequest(`/partos${query}`);
  },
  getById: async (id) => {
    return apiRequest(`/partos/${id}`);
  },
  create: async (partoData) => {
    return apiRequest('/partos', {
      method: 'POST',
      body: JSON.stringify(partoData)
    });
  },
  update: async (id, partoData) => {
    return apiRequest(`/partos/${id}`, {
      method: 'PUT',
      body: JSON.stringify(partoData)
    });
  },
  getFuturos: async () => {
    return apiRequest('/partos/futuros');
  }
};

// API de Destetes
export const destetesAPI = {
  getAll: async (mes = null, año = null) => {
    const params = new URLSearchParams();
    if (mes) params.append('mes', mes);
    if (año) params.append('año', año);
    const query = params.toString() ? `?${params.toString()}` : '';
    return apiRequest(`/destetes${query}`);
  },
  getById: async (id) => {
    return apiRequest(`/destetes/${id}`);
  },
  create: async (desteteData) => {
    return apiRequest('/destetes', {
      method: 'POST',
      body: JSON.stringify(desteteData)
    });
  },
  update: async (id, desteteData) => {
    return apiRequest(`/destetes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(desteteData)
    });
  },
  getFuturos: async () => {
    return apiRequest('/destetes/futuros');
  },
  getEstadisticas: async () => {
    return apiRequest('/destetes/estadisticas');
  }
};

// API de Estadísticas
export const estadisticasAPI = {
  getCerdas: async () => {
    return apiRequest('/estadisticas/cerdas');
  },
  getServicios: async () => {
    return apiRequest('/estadisticas/servicios');
  },
  getPartos: async () => {
    return apiRequest('/estadisticas/partos');
  }
};

// API de Calendario
export const calendarioAPI = {
  getEventosFuturos: async () => {
    return apiRequest('/calendario/eventos-futuros');
  }
};

// API de Reportes
export const reportesAPI = {
  exportarExcel: async (tipo, filtros = {}) => {
    const params = new URLSearchParams(filtros).toString();
    const url = `${API_BASE_URL}/reportes/exportar/${tipo}${params ? `?${params}` : ''}`;
    const token = getAuthToken();
    
    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Error al exportar');
    }
    
    const blob = await response.blob();
    const downloadUrl = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = downloadUrl;
    a.download = `reporte_${tipo}_${new Date().toISOString().split('T')[0]}.xlsx`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(downloadUrl);
  }
};
