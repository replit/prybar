require = require("module").createRequire(
  global[Symbol.for("replit.prybar.eval.file")]
);
module.exports = eval(global[Symbol.for("replit.prybar.eval.code")]);
