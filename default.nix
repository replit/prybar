{ pkgs ? import <nixpkgs>{} } :
let
    callPackage = pkgs.callPackage;
    python3 = pkgs.python3;
    python2 = pkgs.python2;
    lua5_1 = pkgs.lua5_1;
    readline = pkgs.readline;
    clojure = pkgs.clojure;
    tcl = pkgs.tcl;
in
let
  buildPrybar = attrs: callPackage(import ./build-prybar.nix attrs) {};
in
{
    # TODO: These don't quite work yet

    # Won't build
    # prybar-R = buildPrybar { language = "R"; buildInputs = [ pkgs.R ]; };

    # Julia 1.1 doesn't seem to be on nix
    # prybar-julia = buildPrybar { language = "julia"; buildInputs = [ julia_13 ]; };

    # Has problem reading from prybar_assets. Seems like a bug with the impl.
    # prybar-nodejs = buildPrybar { language = "nodejs"; buildInputs = [ ]; };

    # Ruby 2.5 is deprecated!
    # prybar-ruby = buildPrybar { language = "ruby"; buildInputs = [ ruby_2_5 ]; };

    # Untested, spidermonkey not supported on darwin. Feel free to test on linux.
    # prybar-spidermonkey = buildPrybar { language = "spidermonkey"; buildInputs = [ spidermonkey ]; };

    prybar-python2 = buildPrybar { language = "python2"; buildInputs = [ python2 ]; };

    prybar-python3 = buildPrybar { language = "python3"; buildInputs = [ python3 ]; };

    prybar-lua = buildPrybar { language = "lua"; buildInputs = [ lua5_1 readline ]; };

    prybar-clojure = buildPrybar { language = "clojure"; buildInputs = [ clojure ]; };

    prybar-elisp = buildPrybar { language = "elisp"; };

    prybar-ocaml = buildPrybar { language = "ocaml"; };

    prybar-scala = buildPrybar { language = "scala"; };

    prybar-sqlite = buildPrybar { language = "sqlite"; };

    prybar-tcl = buildPrybar { language = "tcl"; buildInputs = [ tcl ]; };
}