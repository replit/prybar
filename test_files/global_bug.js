// In this bug, the main module's globals were not the same
// as the global in other modules. Also some things were
// missing from the global like timing functions 
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

// Current implementation has this bug
// const { Object: importedObject, global: importedGlobal} = require('./global_bug_submodule')
// console.assert(importedObject === Object);
// console.assert(global === importedGlobal)
