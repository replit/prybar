(require 'ielm)

(defun prybar-config (key)
  "Read a Prybar configuration variable from the environment.
KEY is a symbol, like `quiet'. Return the value of the
environment variable \"PRYBAR_QUIET\" (or analogously for other
KEYs). If the environment variable is unset, or its value is the
empty string, then return nil instead."
  (let* ((name (format "PRYBAR_%s" (upcase (symbol-name key))))
         (val (getenv name)))
    (and val (/= (length val) 0) val)))

(defmacro with-prybar-config (vars &rest body)
  "Read Prybar configuration VARS and bind lexically while running BODY.
VARS is a list of symbols as could be passed to `prybar-config'.
Each symbol is prepended with `prybar-' and bound lexically to
the return value of `prybar-config'.

VARS may also contain lists of the form (SYMBOL DEFAULT), not
just symbols. In this case, the behavior is the same except that
DEFAULT is used as the binding for the symbol if `prybar-config'
returns nil.

For example, the following usage:

  (with-prybar-config (eval (ps1 \"-->\"))
    ...)

expands to:

  (let ((prybar-eval (prybar-config 'eval))
       (prybar-ps1 (or (prybar-config 'ps1) \"-->\")))
    ...)"
  (declare (indent 1))
  `(let (,@(mapcar
            (lambda (var)
              (if (listp var)
                  `(,(intern (format "prybar-%S" (nth 0 var)))
                    (or (prybar-config ',(nth 0 var))
                        ,(nth 1 var)))
                `(,(intern (format "prybar-%S" var))
                  (prybar-config ',var))))
            vars))
     ,@body))

(defun prybar-repl ()
  "Start a Prybar REPL at the top level. This mutates global state.
The REPL uses `ielm', which see. After version information is
printed and forms for immediate evaluation or execution are
handled, display the IELM buffer and return."
  ;; For some reason in the repl.it environment the text all gets
  ;; underlined sometimes. Redisplaying fixes this.
  (redisplay)
  (menu-bar-mode -1)
  ;; Make it so you can `load' or `require' files from the project
  ;; directory, and any lisp/ subdirectory (often seen in larger
  ;; projects).
  (add-to-list 'load-path default-directory)
  (add-to-list 'load-path (expand-file-name "lisp" default-directory))
  ;; IELM only supports PS1, not PS2.
  (with-prybar-config
      (eval exec (ps1 "--> ") quiet file)
    (setq ielm-prompt prybar-ps1)
    (with-current-buffer (get-buffer-create "*ielm*")
      (if prybar-quiet
          (setq ielm-header "")
        (insert (format "GNU Emacs %s\n" emacs-version)))
      (when prybar-file
        (condition-case e
            (load (expand-file-name prybar-file) nil 'nomessage)
          (error
           (insert (format "%s\n" (error-message-string e))))))
      (inferior-emacs-lisp-mode)
      (when prybar-exec
        (insert prybar-exec)
        (ielm-send-input 'for-effect))
      (when prybar-eval
        (insert prybar-eval)
        (ielm-send-input)))
    (pop-to-buffer-same-window "*ielm*")))
