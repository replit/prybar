{ pkgs }:

let
    cljdeps = import ./deps.nix { inherit pkgs; };
    classp = cljdeps.makeClasspaths {};
in pkgs.stdenv.mkDerivation {
    name = "wrapped-clojure-with-deps";

    src = ./.;
    nativeBuildInputs = [ pkgs.makeWrapper ];
    buildInputs = [ pkgs.clojure ];

    installPhase = ''
        mkdir -p $out/bin

        makeWrapper ${pkgs.clojure}/bin/clojure $out/bin/clojure --add-flags "-Scp ${builtins.toString classp}"
    '';
}
