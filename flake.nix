{
  description = "A universal interpreter front-end";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?rev=3e0ce8c5d478d06b37a4faa7a4cc8642c6bb97de";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        rev = if self ? rev then
          "0.0.0-${builtins.substring 0 7 self.rev}"
        else
          "0.0.0-dirty";

        buildPrybar = attrs:
          pkgs.callPackage (import ./build-prybar.nix attrs) { inherit rev; };

        julia = import ./languages/julia/julia.nix { inherit pkgs; };

        clojureWithCP =
          import ./languages/clojure/wrappedClojure.nix { inherit pkgs; };

      in {
        packages = {
          prybar-R = buildPrybar {
            language = "R";
            binaries = [ pkgs.R ];
            pkgName = "libR";
            setFlags = true;
          };

          prybar-nodejs = buildPrybar {
            language = "nodejs";
            binaries = [ pkgs.nodejs ];
          };

          prybar-python2 = buildPrybar {
            language = "python2";
            binaries = [ pkgs.python2 ];
          };

          prybar-python3 = buildPrybar {
            language = "python3";
            binaries = [ pkgs.python3 ];
          };

          prybar-lua = buildPrybar {
            language = "lua";
            binaries = [ pkgs.lua5_1 pkgs.readline ];
          };

          prybar-clojure = buildPrybar {
            language = "clojure";
            buildInputs = [ clojureWithCP ];
            binaries = [ pkgs.jdk11_headless ];
          };

          prybar-elisp = buildPrybar { language = "elisp"; };

          prybar-ocaml = buildPrybar {
            language = "ocaml";
            # This has to be in the build inputs so it set-ups the library paths before the check phase
            # It will not be copied over to the final package.
            binaries = [
              pkgs.opam
              pkgs.ocaml
              pkgs.ocamlPackages.findlib
              pkgs.ocamlPackages.topkg
              pkgs.ocamlPackages.reason
            ];
          };

          prybar-scala = buildPrybar { language = "scala"; };

          prybar-sqlite = buildPrybar {
            language = "sqlite";
            binaries = [ pkgs.sqlite ];
          };

          prybar-tcl = buildPrybar {
            language = "tcl";
            binaries = [ pkgs.tcl ];
          };
        }
        # These packages have issues on macOS
          // nixpkgs.lib.optionalAttrs (!pkgs.stdenv.isDarwin) {
            prybar-julia = buildPrybar {
              language = "julia";
              buildInputs = [ julia ];
              setFlags = true;
              binaries = [ pkgs.zlib ];
            };

            # The official Ruby packages is running ruby 2.7
            prybar-ruby = buildPrybar {
              language = "ruby";
              binaries = [ pkgs.ruby ];
              pkgName = "ruby-2.7";
              setFlags = true;
            };

          };

      });
}

