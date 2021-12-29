const repl = require('repl');
const path = require('path');
const { isatty } = require('tty');
const assets_dir =
  process.env.PRYBAR_ASSETS_DIR || path.join(process.cwd(), 'prybar_assets');
/**
 * @type {import("../../prybar_assets/nodejs/input-sync")}
 */
const rl = require(path.join(assets_dir, 'nodejs', 'input-sync.js'));
/**
 * @type {import("../../prybar_assets/nodejs/module-context-hook")}
 */
const { runCode, runModule, getRepl } = require(path.join(
  assets_dir,
  'nodejs',
  'module-context-hook.js',
));

// imports to builtin modules don't get added to require.cache.
const prybarFilenames = Object.keys(require.cache);

for (const filename of prybarFilenames) {
  delete require.cache[filename];
}

Error.prepareStackTrace = function prepareStackTrace(error, callSites) {
  // this logic is sourced from the internal repl error handler

  // Search from the bottom of the call stack to
  // find the first frame with a null function name
  callSites.reverse();

  const idx = callSites.findIndex((frame) => frame.getFunctionName() === null);
  const domainIndex = callSites.findIndex(
    (site) => site.getFileName() === 'domain.js',
  );

  if (domainIndex !== -1 && domainIndex < idx) {
    // If found, get rid of it and everything below it
    callSites = callSites.slice(idx);
  }

  callSites.reverse();

  const lowestPrybarFileIndex = callSites.findIndex((site) =>
    prybarFilenames.includes(site.getFileName()),
  );

  callSites = callSites.slice(0, lowestPrybarFileIndex);

  if (!(error instanceof Error)) {
    const tmp = new Error();

    tmp.message = error.message;
    tmp.name = error.name;
    error = tmp;
  }

  return callSites.length > 0
    ? `${error.toString()}\n    at ${callSites.join('\n    at ')}`
    : error.toString();
};

const isInterractive = Boolean(process.env.PRYBAR_INTERACTIVE);

let r;
if (!process.env.PRYBAR_QUIET) {
  console.log('Node ' + process.version + ' on ' + process.platform);
}

const isTTY = isatty(process.stdin.fd);

// Red errors (if stdout is a TTY)
function logError(msg) {
  if (isTTY) {
    process.stderr.write(`\u001b[0m\u001b[31m${msg}\u001b[0m`);
  } else {
    process.stderr.write(msg);
  }

  if (!msg.endsWith('\n')) {
    process.stderr.write('\n');
  }

  if (!msg.endsWith('\n')) {
    process.stdout.write('\n');
  }
}

// The nodejs repl operates in raw mode and does some funky stuff to
// the terminal. This ns the repl and forces non-raw mode.
function pauseRepl() {
  if (!r) return;

  r.pause();
}

// Forces raw mode and resumes the repl.
function resumeRepl() {
  if (!r) return;
  r.resume();
}

// Clear the line if it has anything on it.
function clearLine() {
  if (isTTY && r && r.line) r.clearLine();
}

// Adapted from the internal node repl code just a lot simpler and adds
// red errors (see https://bit.ly/2FRM86S)
function handleError(e) {
  if (r) {
    r.lastError = e;
  }

  if (e && typeof e === 'object' && e.stack && e.name) {
    if (e.name === 'SyntaxError') {
      e.stack = e.stack
        .replace(/^repl:\d+\r?\n/, '')
        .replace(/^\s+at\s.*\n?/gm, '');
    }

    logError(e.stack);
  } else {
    // For some reason needs a newline to flush.
    logError('Thrown: ' + r.writer(e));
  }

  if (r) {
    r.clearBufferedCommand();
    r.lines.level = [];
    r.displayPrompt();
  }
}

function start() {
  /** @type { repl.REPLServer} */
  r =
    getRepl() ||
    repl.start({
      useGlobal: true,
      prompt: process.env.PRYBAR_PS1,
    });

  // remove the internal error and ours for red etc.
  r._domain.removeListener('error', r._domain.listeners('error')[0]);

  r._domain.on('error', handleError);
  process.on('uncaughtException', handleError);
  return false;
}

global.alert = console.log;
global.prompt = (p) => {
  pauseRepl();
  clearLine();

  let ret = rl.question(`${p}> `);

  resumeRepl();

  // Display prompt on the next turn.
  if (r) setImmediate(() => r.displayPrompt());

  return ret;
};

global.confirm = (q) => {
  pauseRepl();
  clearLine();

  const ret = rl.keyInYNStrict(q);

  resumeRepl();

  // Display prompt on the next turn.
  if (r) setImmediate(() => r.displayPrompt());
  return ret;
};

if (process.env.PRYBAR_CODE) {
  try {
    runCode(process.env.PRYBAR_CODE, isInterractive);
  } catch (err) {
    handleError(err);

    if (!isInterractive) {
      process.exit(1)
    }
  }

  if (isInterractive) {
    process.once('beforeExit', start);
  }
} else if (process.env.PRYBAR_EXP) {
  try {
    console.log(runCode(process.env.PRYBAR_EXP, false));
  } catch (err) {
    handleError(err);

    if (!isInterractive) {
      process.exit(1)
    }
  }
} else if (process.env.PRYBAR_FILE) {
  try {
    runModule(process.env.PRYBAR_FILE, isInterractive);
  } catch (err) {
    handleError(err);

    if (!isInterractive) {
      process.exit(1)
    }
  }

  if (isInterractive) {
    if (isTTY) {
      console.log(
        '\u001b[0m\u001b[90mHint: hit control+c anytime to enter REPL.\u001b[0m',
      );
    }

    process.once('beforeExit', start);
    process.once('SIGINT', start);
  }
} else if (isInterractive) {
  start();
}
