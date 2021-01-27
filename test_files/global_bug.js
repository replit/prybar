console.assert([] instanceof global.Array);
console.assert([] instanceof Array);
console.assert(process);
console.assert(URL);
console.assert(URLSearchParams);
console.assert(clearImmediate);
console.assert(clearInterval);
console.assert(clearTimeout);
console.assert(setImmediate);
console.assert(setInterval);
console.assert(setTimeout);
console.assert(Buffer);
console.assert(globalThis === global);

const { Object: importedObject, global: importedGlobal} = require('./global_require')
console.assert(importedObject === Object);
console.assert(global === importedGlobal)
