// Make sure that the module has the right values
console.assert(module.id === '.')
console.assert(Object.keys(module.exports).length === 0)
console.assert(module.parent === null)
console.assert(module.loaded === false)
console.assert(module.children.length === 0)
console.assert(module.filename === __filename)
console.assert(module.filename.endsWith('/test_files/module_patching.js'));