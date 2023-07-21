{
  description = "A universal interpreter front-end";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-22.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    {
      overlays.default = final: prev: {
        prybar = prev.lib.recurseIntoAttrs {
          inherit (self.packages.${prev.system})
            prybar-R prybar-clojure prybar-elisp prybar-julia prybar-lua prybar-nodejs
            prybar-ocaml prybar-python2 prybar-python3 prybar-python310 prybar-ruby
            prybar-scala prybar-sqlite prybar-tcl;
        };
      };
    } //
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
        
        python310Full = pkgs.python310.override {
          self = python310Full;
          pythonAttr = "python310Full";
          bluezSupport = true;
          x11Support = true;
        };

        ocaml = pkgs.ocaml-ng.ocamlPackages_4_12.ocaml;
        opam = pkgs.opam;
        findlib = pkgs.ocaml-ng.ocamlPackages_4_12.findlib;
        topkg = pkgs.ocaml-ng.ocamlPackages_4_12.topkg;
        reason = pkgs.ocaml-ng.ocamlPackages_4_12.reason;

      in {
        packages = {
          prybar-R = buildPrybar {
            language = "R";
            buildInputs = [ pkgs.R ];
            pkgName = "libR";
            setFlags = true;
          };

          prybar-nodejs = buildPrybar {
            language = "nodejs";
            binaries = [ pkgs.nodejs-16_x ];
          };

          prybar-python2 = buildPrybar {
            language = "python2";
            buildInputs = [ pkgs.python2 ];
          };

          prybar-python3 = buildPrybar {
            language = "python3";
            buildInputs = [ pkgs.libxcrypt pkgs.python38Full ];
          };

          prybar-python310 = buildPrybar {
            language = "python310";
            buildInputs = [ pkgs.libxcrypt python310Full ];
          };

          prybar-lua = buildPrybar {
            language = "lua";
            buildInputs = [ pkgs.lua5_1 pkgs.readline ];
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
            buildInputs = [
              opam
              ocaml
              findlib
              topkg
              reason
            ];
          };

          prybar-scala = buildPrybar { language = "scala"; };

          prybar-hoon = buildPrybar { language = "hoon"; };

          prybar-sqlite = buildPrybar {
            language = "sqlite";
            binaries = [ pkgs.sqlite ];
          };

          prybar-tcl = buildPrybar {
            language = "tcl";
            buildInputs = [ pkgs.tcl ];
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
              buildInputs = [ pkgs.ruby ];
              pkgName = "ruby-2.7";
              setFlags = true;
            };

          };


        devShells.default = pkgs.mkShell {
          nativeBuildInputs = [
            pkgs.pkg-config
          ];
          buildInputs = [
            pkgs.libxcrypt
            pkgs.R
            pkgs.nodejs-16_x
            pkgs.python2
            pkgs.python38Full
            pkgs.python310Full
            pkgs.lua5_1
            pkgs.readline
            clojureWithCP
            pkgs.jdk11_headless
            opam
            ocaml
            findlib
            topkg
            reason
            julia
            pkgs.zlib
            pkgs.ruby
            pkgs.sqlite
            pkgs.tcl
            pkgs.expect
            pkgs.urbit
          ];
          shellHook = ''
            export CGO_LDFLAGS_ALLOW="-Wl,--compress-debug-sections=zlib"
            export DISABLE_GO_TESTS=1
          '';
        };
      

      });
}

