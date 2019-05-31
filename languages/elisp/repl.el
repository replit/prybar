(defun prybar-config (key)
  (let* ((name (format "PRYBAR_%s" (upcase (symbol-name key))))
         (val (getenv name)))
    (and val (/= (length val) 0) val)))

(defmacro with-prybar-config (vars &rest body)
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

(defun prybar-read-string (ps)
  (condition-case e
      (read-string ps)
    (error
     (if (equal (cadr e) "Error reading from stdin")
         (progn
           (terpri)
           (throw :eof nil))
       (signal (car e) (cdr e))))))

(defun prybar-parse-forms (str)
  (with-temp-buffer
    (insert str)
    (goto-char (point-min))
    (let ((forms nil))
      (catch :return
        (while t
          (skip-chars-forward "[:space:]")
          (when (= (point) (point-max))
            (throw :return (nreverse forms)))
          (condition-case e
              (push (read (current-buffer)) forms)
            (error
             (if (eq (car e) 'end-of-file)
                 (throw :return :unfinished)
               (message "Read error: %s" (error-message-string e))
               (throw :error nil)))))))))

(defun prybar-parse-single-form (str)
  (let ((forms (prybar-parse-forms)))
    (prog1 forms
      (when (eq forms :unfinished)
        (message "End of string during parsing")
        (throw :error nil)))))

(defun prybar-read-forms (ps1 ps2)
  (let ((str (prybar-read-string ps1))
        (forms nil))
    (while (eq :unfinished (setq forms (prybar-parse-forms str)))
      (setq str (concat str "\n" (prybar-read-string ps2))))
    forms))

(defun prybar-eval-or-backtrace (form)
  (let* ((debug-on-error t)
         (debugger
          (lambda (&rest args)
            (message "Lisp error: %S" args)
            ;; For some reason, binding `standard-output' to t does
            ;; not do the trick.
            (with-temp-buffer
              (let ((standard-output (current-buffer)))
                (backtrace)
                (message "%s" (buffer-string))))
            (exit-recursive-edit))))
    (condition-case-unless-debug _
        (eval form)
      (error (throw :error nil)))))

(defun prybar-repl ()
  (let ((debug-on-error t))
    (with-prybar-config
        (eval exec (ps1 "--> ") (ps2 "... ") quiet interactive)
      (unless prybar-quiet
        (message "GNU Emacs %s" emacs-version))
      (catch :error
        (when prybar-exec
          (prybar-eval-or-backtrace (prybar-parse-single-form prybar-exec)))
        (when prybar-eval
          (prin1
           (prybar-eval-or-backtrace (prybar-parse-single-form prybar-eval)))
          (terpri))
        (when prybar-interactive
          (catch :eof
            (while t
              (catch :error
                (let ((forms (prybar-read-forms prybar-ps1 prybar-ps2)))
                  (dolist (form forms)
                    (let ((result (prybar-eval-or-backtrace form)))
                      (prin1 result)
                      (terpri))))))))))))

;; TODO: deal with extra newlines
;; TODO: fix debugger getting disabled
