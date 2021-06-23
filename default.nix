{ pkgs ? import <nixpkgs>{}, versionArg ? "" } :
let
  buildPrybar = attrs: pkgs.callPackage(import ./build-prybar.nix attrs) {};

  julia = import ./languages/julia/julia.nix { inherit pkgs; };
  clojureWithCP = import ./languages/clojure/wrappedClojure.nix { inherit pkgs; };
in
{
    prybar-R = buildPrybar {
        language = "R";
        buildInputs = [ pkgs.R ];
        versionArg = versionArg;
        pkgName = "libR";
        setFlags = true;
    };

    prybar-julia = buildPrybar {
        language = "julia";
        buildInputs = [ julia ];
        versionArg = versionArg;
        setFlags = true;
    };

    prybar-nodejs = buildPrybar { language = "nodejs"; binaries = [ pkgs.nodejs ]; versionArg = versionArg;  };

    # The official Ruby packages is running ruby 2.7
    prybar-ruby = buildPrybar {
        language = "ruby";
        buildInputs = [ pkgs.ruby ];
        pkgName = "ruby-2.7";
        versionArg = versionArg;
        setFlags = true;
    };

    prybar-python2 = buildPrybar { language = "python2"; buildInputs = [ pkgs.python2 ]; versionArg = versionArg; };

    prybar-python3 = buildPrybar { language = "python3"; buildInputs = [ pkgs.python3 ]; versionArg = versionArg; };

    prybar-lua = buildPrybar { language = "lua"; buildInputs = [ pkgs.lua5_1 pkgs.readline ]; versionArg = versionArg; };

    prybar-clojure = buildPrybar {
        language = "clojure";
        buildInputs = [ clojureWithCP ];
        binaries = [ pkgs.jdk11_headless ];
        versionArg = versionArg;
    };

    prybar-elisp = buildPrybar { language = "elisp"; versionArg = versionArg; };

    prybar-ocaml = buildPrybar {
        language = "ocaml";
        # This has to be in the build inputs so it set-ups the library paths before the check phase
        # It will not be copied over to the final package.
        buildInputs = [ pkgs.opam pkgs.ocaml pkgs.ocamlPackages.findlib pkgs.ocamlPackages.topkg pkgs.ocamlPackages.reason ];
        versionArg = versionArg;
    };

    prybar-scala = buildPrybar { language = "scala"; versionArg = versionArg; };

    prybar-sqlite = buildPrybar { language = "sqlite"; binaries = [ pkgs.sqlite ]; versionArg = versionArg; };

    prybar-tcl = buildPrybar { language = "tcl"; buildInputs = [ pkgs.tcl ]; versionArg = versionArg; };
}