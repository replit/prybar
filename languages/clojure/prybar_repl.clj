(def prybar-ps1 (System/getProperty "PRYBAR_PS1" "--> "))
(def prybar-quiet? (System/getProperty "PRYBAR_QUIET" "false"))

(when-not (Boolean/valueOf prybar-quiet?)
  (println "Clojure" (clojure-version)))

(clojure.main/repl :init #(apply require clojure.main/repl-requires)
                   :prompt #(print prybar-ps1))
(prn)
(System/exit 0)

