// In this bug the main module was adding module scope variables
// to the global scope!
const shouldNotBeAvailableInOtherModules = 'a';

require('./main_module_scoped_to_global_submodule');