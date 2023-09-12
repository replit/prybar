{
  description = "A universal interpreter front-end";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    {
      overlays.default = final: prev: {
        prybar = prev.lib.recurseIntoAttrs {
          inherit (self.packages.${prev.system})
            prybar-elisp prybar-julia prybar-nodejs
            prybar-python2 prybar-python3 prybar-python38 prybar-python310
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

        python38Full = pkgs.python38Full.override {
          self = python38Full;
          pythonAttr = "python38Full";
          bluezSupport = true;
          x11Support = true;
        };
        
        python310Full = pkgs.python310.override {
          self = python310Full;
          pythonAttr = "python310Full";
          bluezSupport = true;
          x11Support = true;
        };

        python311Full = pkgs.python311.override {
          self = python311Full;
          pythonAttr = "python311Full";
          bluezSupport = true;
          x11Support = true;
        };

      in {
        packages = rec {

          prybar-nodejs = buildPrybar {
            language = "nodejs";
            binaries = [ pkgs.nodejs ];
          };

          prybar-python38 = buildPrybar {
            language = "python3";
            target = "python38";
            cgoPkgs = "python-3.8-embed";
            cgoExtraCflags = "-DPYTHON_3_8";
            buildInputs = [ pkgs.libxcrypt python38Full ];
          };

          prybar-python3 = prybar-python38;

          prybar-python310 = buildPrybar {
            language = "python3";
            target = "python310";
            cgoPkgs = "python-3.10-embed";
            cgoExtraCflags = "-DPYTHON_3_10";
            buildInputs = [ pkgs.libxcrypt python310Full ];
          };

          prybar-python311 = buildPrybar {
            language = "python3";
            target = "python311";
            cgoPkgs = "python-3.11-embed";
            cgoExtraCflags = "-DPYTHON_3_11";
            buildInputs = [ pkgs.libxcrypt python311Full ];
          };

          prybar-elisp = buildPrybar { language = "elisp"; };

          prybar-scala = buildPrybar { language = "scala"; };

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

          };


        devShells.default = pkgs.mkShell {
          nativeBuildInputs = [
            pkgs.pkg-config
          ];
          buildInputs = [
            pkgs.libxcrypt
            # pkgs.R
            pkgs.nodejs
            pkgs.python38Full
            pkgs.python310Full
            pkgs.python311Full
            # pkgs.lua5_1
            pkgs.readline
            # clojureWithCP
            # pkgs.jdk11_headless
            # julia
            pkgs.zlib
            # pkgs.ruby
            pkgs.sqlite
            pkgs.tcl
            pkgs.expect
          ];
          shellHook = ''
            export CGO_LDFLAGS_ALLOW="-Wl,--compress-debug-sections=zlib"
            export DISABLE_GO_TESTS=1
          '';
        };
      });
}

