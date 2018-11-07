var repl = require('repl');
var exit = require('async-exit-hook');


exit((reallyExit) => {
	let r = repl.start(process.env.NODE_PROMPT || '> ');
	r.on('exit', reallyExit);	
});
