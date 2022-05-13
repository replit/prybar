let didFail = false;

exports.assertEqual = function (val, exp, message) {
  if (exp !== val) {
    if (didFail) {
      // print a blank line
      console.error();
    }
    console.error("\x1b[31m[FAIL]\x1b[0m", message);
    console.error(`  \x1b[32mExpected:\x1b[0m ${exp}`);
    console.error(`  \x1b[31mReceived:\x1b[0m ${val}`);
    didFail = true;
  }
};

/**
 * @type {() => { line: number, column: number, file: string }}
 */
exports.getCallSite = function getCallSite() {
  let stack;
  const prevLimit = Error.stackTraceLimit;

  try {
    const obj = {};
    Error.prepareStackTrace = function (_err, [callSite]) {
      return {
        file: callSite.getFileName(),
        line: callSite.getLineNumber(),
        column: callSite.getColumnNumber(),
      };
    };
    Error.stackTraceLimit = 2;
    Error.captureStackTrace(obj, getCallSite);

    stack = obj.stack;
  } finally {
    // This shouldn't throw any errors error, but if it did, having these overwritten
    // would make the resulting error's stack trace useless.
    Error.stackTraceLimit = prevLimit;
    Error.prepareStackTrace = undefined;
  }

  return stack;
};

process.once("beforeExit", () => {
  if (didFail) {
    process.exit(1);
  }
});
