const Module = require('module');
const { readFileSync } = require('fs');
const { parseSync, transformSync } = require('@swc/core');
const sourceMapSupport = require('source-map-support');

// Source maps won't always be enabled.  This is to stay consistent with normal nodejs.
// However, we may want to support source maps natively in the future (e.g. with typescript)
// so prybar will also allow it, behind a flag.
const sourceMapsEnabled = Boolean(
  (process.env.ENABLE_SOURCE_MAPS && process.env.ENABLE_SOURCE_MAPS !== '0') ||
    (process.env.NODE_OPTIONS && getCallSite.includes('--enable-source-maps'))
);

/**
 * The maps for files with esm syntax which have been transpiled to work with cjs.
 * @type {Map<string, {url: string, map: string}>}
 */
const maps = new Map();

/**
 * @type {string | undefined}
 */
let fileContent;

sourceMapSupport.install({
  environment: 'node',
  handleUncaughtExceptions: false,
  hookRequire: false,
  retrieveFile() {
    if (fileContent) {
      const tmp = fileContent;
      fileContent = undefined;
      return tmp;
    }
  },
  // If source maps aren't enabled, we'll override the default behavior of resolving source maps.
  overrideRetrieveSourceMap: !sourceMapsEnabled,
  retrieveSourceMap(filename) {
    return maps.get(filename);
  },
});

const sourceMapSupportPrepareStackTrace = Error.prepareStackTrace;
delete Error.prepareStackTrace;

// prybar assets shouldn't be added to require.cache.
const prybarFilenames = Object.keys(require.cache);

for (const filename of prybarFilenames) {
  delete require.cache[filename];
}

const { getPrototypeOf, setPrototypeOf, defineProperty } = Object;

const defaultCompile = Module.prototype._compile;

Module.prototype._compile = function (code, filename) {
  const parsed = parseSync(code, {
    target: 'es2021',
    syntax: 'ecmascript',
  });

  for (const stmt of parsed.body) {
    if (
      stmt.type === 'ExportDefaultDeclaration' ||
      stmt.type === 'ExportDefaultExpression' ||
      stmt.type === 'ExportNamedDeclaration' ||
      stmt.type === 'ExportDeclaration' ||
      stmt.type === 'ImportDeclaration' ||
      stmt.type === 'ExportAllDeclaration'
    ) {
      // Despite what the documentation clearly states, if inputSourceMap
      // is explicitly set to true swc will throw an error when unable to load source maps for a file.
      // (additionally, this seems to default to false; again contrary to the documentation)
      let resolvedMap = sourceMapsEnabled;

      if (sourceMapsEnabled) {
        fileContent = code;
        // source-map-support has working source map logic so we'll use that instead of swc.
        resolvedMap = sourceMapSupport.retrieveSourceMap(filename);
      }

      const tr = transformSync(code, {
        filename,
        sourceFileName: resolvedMap ? resolvedMap.url : filename,

        sourceMaps: true,
        // We ONLY want to map back what we changed, otherwise this may interfere with handlers
        // created by users. (e.g. the user may try to map to the original source after we already did,
        // causing confusing results)
        inputSourceMap: resolvedMap ? resolvedMap.map : false,
        env: {
          targets: ['node 16'],
        },
        module: { type: 'commonjs' },
      });

      maps.set(filename, {
        map: tr.map,
        url: resolvedMap.url,
      });

      return defaultCompile.call(this, tr.code, filename);
    }
  }

  return defaultCompile.call(this, code, filename);
};

Module._extensions['.js'] = function (module, filename) {
  const source = readFileSync(filename, 'utf8');

  // Require hooks could be defined outside of prybar.  If so, they almost always end up calling
  // Module.prototype._compile.  To make this as transparent (and non-intrusive) as possible, we need to
  // make sure that we're the last wrapper of Module._compile.
  //
  // One example of where this could otherwise cause issues is `source-map-support` which can use require hooks
  // to prime its cache.
  // https://github.com/evanw/node-source-map-support/blob/ac2c3e4c633c66931981ac94b44e6963addbe3f4/source-map-support.js#L565-L571
  return module._compile(source, filename);
};

Module._extensions['.mjs'] = Module._extensions['.js'];

let userSpacePrepareStackTrace;
let isNestedCall = false;

function prepareStackTrace(error, callSites) {
  if (isNestedCall) {
    if (!(error instanceof Error)) {
      const tmp = new Error();

      tmp.message = error.message;
      tmp.name = error.name;
      error = tmp;
    }

    return callSites.length > 0
      ? `${error.toString()}\n    at ${callSites.join('\n    at ')}`
      : error.toString();
  }

  // this logic is sourced from the internal repl error handler

  // Search from the bottom of the call stack to
  // find the first frame with a null function name
  callSites.reverse();

  const idx = callSites.findIndex((frame) => frame.getFunctionName() === null);
  const domainIndex = callSites.findIndex(
    (site) =>
      site.getFileName() === 'domain.js' || site.getFileName() === 'node:domain'
  );

  if (domainIndex !== -1 && domainIndex < idx) {
    // If found, get rid of it and everything below it
    callSites = callSites.slice(idx);
  }

  callSites.reverse();

  // <begin hand-written code>
  const lowestPrybarFileIndex = callSites.findIndex((site) =>
    prybarFilenames.includes(site.getFileName())
  );

  callSites = callSites.slice(0, lowestPrybarFileIndex);

  if (userSpacePrepareStackTrace) {
    try {
      isNestedCall = true;
      return userSpacePrepareStackTrace(
        error,
        callSites.map((site) => {
          const filename = site.getFileName();

          if (sourceMapsEnabled || maps.has(filename)) {
            const result = setPrototypeOf(
              sourceMapSupport.wrapCallSite(site),
              getPrototypeOf(site)
            );

            return result;
          }

          return site;
        })
      );
    } finally {
      isNestedCall = false;
    }
  }

  return sourceMapSupportPrepareStackTrace(error, callSites);
}

defineProperty(Error, 'prepareStackTrace', {
  enumerable: false,
  configurable: false,
  set(fn) {
    userSpacePrepareStackTrace = fn;
  },
  get() {
    return prepareStackTrace;
  },
});
