;; This is a slightly modified version of
;; https://github.com/clojure/clojure/blob/master/src/clj/clojure/main.clj#repl-opt.

;; This script is loaded and run by `clojure.main/script-opt` and since it calls
;; `clojure.main/initialize` in prior to this script, we don't need to call it here.

;; `ps1` and `q` ("quiet") options are passed in as system properties.

(def prybar-ps1 (System/getProperty "PRYBAR_PS1" "--> "))
(def prybar-quiet? (System/getProperty "PRYBAR_QUIET" "false"))

(when-not (Boolean/valueOf prybar-quiet?)
  (println "Clojure" (clojure-version)))

(clojure.main/repl :init #(apply require clojure.main/repl-requires)
                   :prompt #(print prybar-ps1))

(prn)
(System/exit 0)

;; How to bootstrap a REPL with this script:
;;     clj -J-DPRYBAR_PS1="==> " -J-DPRYBAR_QUIET=true prybar_repl.clj

