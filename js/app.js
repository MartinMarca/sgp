console.log("App Granja Porcina iniciada");

// Utilidad simple para debug
function log(mensaje) {
  console.log(`[Granja] ${mensaje}`);
}

// Simulación de datos (temporal)
const data = {
  cerdas: [],
  padrillos: [],
  servicios: [],
  partos: [],
  muertes: []
};

// Ejemplo de función
function agregarCerda(cerda) {
  data.cerdas.push(cerda);
  log(`Cerda agregada: ${cerda.identificacion}`);
}
