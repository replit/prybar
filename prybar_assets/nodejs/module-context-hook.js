const path = require("path");
const BuiltinModule = require("module");
const vm = require("vm");
const { Session } = require("inspector");
const Module =
  module.constructor.length > 1 ? module.constructor : BuiltinModule;
const nodeRepl = require("repl");

const kReplitPrybarInit = Symbol.for("replit.prybar.init");
const kReplitPrybarEvalCode = Symbol.for("replit.prybar.eval.code");
const kReplitPrybarEvalFile = Symbol.for("replit.prybar.eval.file");

/**
 * The footer we append at the end of files to run the code we want in the scope of the module its appended to.
 *
 * The code calls global[kReplitPrybarInit] with a callback to evaluate the code at `global`
 */
const interactiveFooterSnippet = `
global[Symbol.for('${kReplitPrybarInit.description}')](() =>
  eval(global[Symbol.for('${kReplitPrybarEvalCode.description}')])
);
`;

const moduleVariables = [
  // these are all added as arguments of the function modules are wrapped in.
  // function (exports, require, module, __filename, __dirname)
  "exports",
  "require",
  "module",
  "__filename",
  "__dirname",
];

/**
 * Regex to get anything that looks like it could be a top-level variable.
 * we use this on each evaluated line strings to test words / strings that could be variables
 * and see if they're in the local scope.  This is needed for the preview which `repl.js` shows & tab auto-complete.
 *
 * This will not match properties of values.
 */
const variableLikeRegex = /(?<=^|[^a-z_0-9$.])[a-z_$][a-z_$0-9]*/gi;

/**
 * Gets a wrapped version of the eval function which takes a string argument (rather than using a symbol on global)
 * and
 *
 * This is intended to removve any confusion from callsites which aren't referenced in the code of the file we're
 * running; calls to eval from lines which don't exist.
 *
 * @param {(code: string) => any} evalFn
 * @returns A a wrapper of the eval function which omits callsites from our proxy logic from stack traces.
 */
function wrapEval(evalFn) {
  return (str) => {
    global[kReplitPrybarEvalCode] = str;

    return evalFn();
  };
}

/**
 * Wraps a function so that none of the callsites under the function's caller are rendered in the error's stack.
 * @template T
 * @param {T} fn
 * @returns {T} A wrapped
 */
function withoutStack(fn) {
  const wrapper = function (...args) {
    const oldPrepareStackTrace = Error.prepareStackTrace;

    Error.prepareStackTrace = (error) => {
      return oldPrepareStackTrace(error, []);
    };

    let results;
    try {
      results = fn.apply(this, args);
    } catch (err) {
      // there's some getter magic on Error.prototype.stack which prevents it from being rendered until
      // its value is requested.
      err.stack;

      throw err;
    } finally {
      Error.prepareStackTrace = oldPrepareStackTrace;
    }

    return results;
  };

  return wrapper;
}

function hasOwnProperty(obj, prop) {
  return Object.hasOwnProperty.call(obj, prop);
}

/**
 * The source of the file which require appended the init hook too, if any.
 * @type {string | null}
 */
let hookedFileSource = null;

/**
 * The context which can be used to access the hooked module's in a repl.
 *
 * @type {nodeRepl.REPLServer | null}
 */
let repl = null;

/**
 * Registers a function to gbe called once as `Module._compile`.
 * @param {(
 *    this: typeof module,
 *    code: string,
 *    fileName: string,
 *    builtin: ((code: string, fileName: string) => any),
 *    mod: any
 *  ) => any} fn The compile function to be used instead of _compile on its next call.
 */
function mockCompileOnce(fn) {
  const defaultCompiler = Module._extensions[".js"];

  Module._extensions[".js"] = (mod, fileName) => {
    // Only the first import should have interp appended to it.
    Module._extensions[".js"] = defaultCompiler;

    const oldCompile = mod._compile;
    mod.id = ".";
    mod.parent = null;

    mod._compile = (code) => {
      mod._compile = oldCompile;

      return fn(code, fileName, oldCompile.bind(mod), mod);
    };

    return defaultCompiler(mod, fileName);
  };
}

function initRepl(evalFn) {
  // wrap the eval fn so that if it errors, it won't include anything from interp in the shell.
  const doEval = wrapEval(evalFn);

  /**
   * A list of variables referenced from the hooked module.  This should include all locals, and
   * will also end up including any referenced globals.
   */
  const locals = {};

  const variableMatches = hookedFileSource.match(variableLikeRegex);
  if (variableMatches) {
    for (const maybeVariable of variableMatches) {
      let value;

      try {
        value = doEval(maybeVariable);

        try {
          // asserts that the variable isn't a constant
          doEval(`${maybeVariable}=${maybeVariable}`);

          Reflect.defineProperty(locals, maybeVariable, {
            value,
            writeable: true,
            enumerable: true,
            configurable: false,
          });
        } catch (err) {
          // can't write to constants
          if (err instanceof TypeError) {
            Reflect.defineProperty(locals, maybeVariable, {
              value,
              writeable: false,
              enumerable: true,
              configurable: false,
            });
            continue;
          }

          throw err;
        }
      } catch (err) {
        if (err instanceof ReferenceError || err instanceof SyntaxError) {
          continue;
        }

        throw err;
      }
    }
  }

  for (const moduleVariable of moduleVariables) {
    if (!hasOwnProperty(locals, moduleVariable)) {
      Object.defineProperty(locals, moduleVariable, {
        writable: true,
        configurable: true,
        value: doEval(moduleVariable),
      });
    }
  }

  repl = nodeRepl.start({
    prompt: process.env.PRYBAR_PS1,
  });

  const replContext = new Proxy(
    repl.context,
    {
      // Returns a list of all properties of the proxy (enumerable or not).
      // https://tc39.es/ecma262/multipage/ordinary-and-exotic-objects-behaviours.html#sec-proxy-object-internal-methods-and-internal-slots-ownpropertykeys
      ownKeys(ctx) {
        return Array.from(
          new Set(
            Reflect.ownKeys(global).concat(
              Reflect.ownKeys(locals),
              Reflect.ownKeys(ctx)
            )
          )
        );
      },
      // Trap for the [[GetOwnProperty]] (used by the Object.getOwnPropertyDescriptor) function
      // Returns the descriptor for any variable.
      // https://tc39.es/ecma262/multipage/ordinary-and-exotic-objects-behaviours.html#sec-proxy-object-internal-methods-and-internal-slots-getownproperty-p
      getOwnPropertyDescriptor(ctx, variable) {
        if (hasOwnProperty(ctx, variable)) {
          return Reflect.getOwnPropertyDescriptor(ctx, variable);
        }

        // If the variable exists locally, prioritize it over global.
        if (hasOwnProperty(locals, variable)) {
          return {
            enumerable: true,
            configurable: true,
            get() {
              return locals[variable];
            },
            set(value) {
              // if the variable is a constant, this will throw
              replContext[variable] = value;
            },
          };
        }

        if (hasOwnProperty(global, variable)) {
          const descriptor = Reflect.getOwnPropertyDescriptor(global);

          return {
            ...descriptor,
            configurable: true,
          };
        }

        return undefined;
      },

      // Trap for determining if the proxy has a key or not.
      // https://tc39.es/ecma262/multipage/ordinary-and-exotic-objects-behaviours.html#sec-proxy-object-internal-methods-and-internal-slots-hasproperty-p
      has(ctx, variable) {
        return (
          hasOwnProperty(locals, variable) ||
          hasOwnProperty(global, variable) ||
          hasOwnProperty(ctx, variable)
        );
      },

      // Trap for getting a property from the proxy.
      // https://tc39.es/ecma262/multipage/ordinary-and-exotic-objects-behaviours.html#sec-proxy-object-internal-methods-and-internal-slots-get-p-receiver
      get(ctx, variable) {
        if (hasOwnProperty(ctx, variable)) {
          return ctx[variable];
        }

        if (hasOwnProperty(locals, variable)) {
          return locals[variable];
        }

        if (hasOwnProperty(global, variable)) {
          return global[variable];
        }

        throw new ReferenceError(`${variable} is not defined`);
      },

      // trap for setting a value in the proxy.
      // https://tc39.es/ecma262/multipage/ordinary-and-exotic-objects-behaviours.html#sec-proxy-object-internal-methods-and-internal-slots-set-p-v-receiver
      set: withoutStack((ctx, variable, value) => {
        if (hasOwnProperty(locals, variable)) {
          // To avoid running into potential name clashes of module variables and
          // variables owned by this REPL logic, we'll evaluate a callback which sets the variable
          // to its first argument, with an argument that has a different name.

          doEval(`_${variable}=>${variable}=_${variable}`)(value);

          locals[variable] = value;
        } else if (hasOwnProperty(global, variable)) {
          global[variable] = value;
        } else {
          ctx[variable] = value;
        }

        return true;
      }),
    }
  );

  repl.context = vm.createContext(replContext);

  for (const moduleVariable of moduleVariables) {
    delete repl.context[moduleVariable];
  }

  const kContextId = Object.getOwnPropertySymbols(repl).find(
    (v) => v.description === "contextId"
  );

  repl.context.repl = repl;

  // this should be 100% seafe since JS is synchronous.
  // we aren't using the repl's original context, but we create the new context directly after
  // repl.start creates one, so the ID of our context is the ID of the repl's context id + 1.
  const replContextId = ++repl[kContextId];

  // The nodejs inspector doesn't properly handle the proxy & scope magic,
  // so we need to mock its eval function if we want commodities such as preview and
  // tab-completion to work as intended.

  const post = Session.prototype.post;
  Session.prototype.post = function (method, params, callback) {
    if (method !== "Runtime.evaluate" || params.contextId !== replContextId) {
      return post.apply(this, arguments);
    }

    return post.call(this, method, params, (err, preview) => {
      const { result } = preview;
      if (preview.exceptionDetails && result.className === "EvalError") {
        return callback(err, preview);
      }

      if (
        (preview.exceptionDetails && result.className === "ReferenceError") ||
        !hasOwnProperty(result, "value")
      ) {
        try {
          const value = vm.runInContext(params.expression, repl.context, {
            displayErrors: false,
          });

          const newResult = { type: typeof value, value };

          return callback(null, { result: newResult });
        } catch {
          callback(err, preview);
        }
      }
      
      return callback(err, preview);
    });
  };
}

global[kReplitPrybarInit] = initRepl;

/**
 * Runs a module; adds a hook for an interactive repl if isInteractive is true.
 *
 * @param {string} moduleName The name / path of the target module
 * @param {boolean} isInteractive Whether the target file should init a repl or not.
 */
function runModule(moduleName, isInteractive) {
  // this will throw if the module doesn't exist.
  const absPath = require.resolve(path.resolve(process.cwd(), moduleName));
  // the the target module has already been `require`d this will drop it from the require cache;
  // otherwise this will do nothing.
  delete require.cache[absPath];

  if (isInteractive) {
    mockCompileOnce((code, fileName, compile) => {
      let compiled;

      try {
        hookedFileSource = code;
        compiled = compile(code + interactiveFooterSnippet, fileName);
      } catch (err) {
        // this error is probably unrelated to our added code
        // and should fail to compile again.
        // It's important that if this is the case, we didn't include our added code
        // otherwise users may see an error on a line of code that doesn't exist in the source.
        // (for example, an extra { which wasn't closed will show up on the last line of the file)
        compiled = compile(code, fileName);

        // If we've gotten here something's up w/ interp, not the module itself.
        console.warn(
          "running without interp",
          err,
          "please submit a bug report."
        );
      }

      return compiled;
    });
  }

  require(absPath);
}

/**
 * Executes the specified code in a fake module and return the results.
 *
 * @param {string} code The code to be executed
 * @param {boolean} isInterractive Whether a repl should be started in the same scope as the code or not.
 * @returns {any} The results of the evalutated code.
 */
function runCode(code, isInterractive) {
  mockCompileOnce((fileSource, __, compile) => {
    const fakeFileName = path.join(process.cwd(), "__replit_exec.js");
    global[kReplitPrybarEvalFile] = fakeFileName;

    hookedFileSource = code;
    delete require.cache[fakeFileName];
    // we need to `eval` here rather than directly input the code for
    global[kReplitPrybarEvalCode] = isInterractive
      ? code + interactiveFooterSnippet
      : code;

    return compile(
      fileSource,
      // lie to the `compile` function about the file we're compiling
      // so calls to `require` will be resolved from the current directory.
      fakeFileName
    );
  });

  delete require.cache[require.resolve("./runCode")];
  return require("./runCode");
}

/**
 * If available, returns a repl in the context of executed code or module.
 *
 * @returns {nodeRepl.REPLServer | null} A vm context if the repl init function was called; otherwise null.
 */
function getRepl() {
  return repl;
}

module.exports = {
  runCode,
  runModule,
  getRepl,
};
