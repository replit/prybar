{
  description = "A universal interpreter front-end";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
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

      in {
        packages = {
          prybar-nodejs = buildPrybar {
            language = "nodejs";
            binaries = [ pkgs.nodejs-16_x ];
          };

          prybar-python3 = buildPrybar {
            language = "python3";
            buildInputs = [ pkgs.python38Full ];
          };

          prybar-python310 = buildPrybar {
            language = "python310";
            buildInputs = [ pkgs.python310Full ];
          };

          prybar-lua = buildPrybar {
            language = "lua";
            buildInputs = [ pkgs.lua5_1 pkgs.readline ];
          };

          prybar-sqlite = buildPrybar {
            language = "sqlite";
            binaries = [ pkgs.sqlite ];
          };

          prybar-tcl = buildPrybar {
            language = "tcl";
            buildInputs = [ pkgs.tcl ];
          };
        };

        devShells.default = pkgs.mkShell {
          nativeBuildInputs = [
            pkgs.pkg-config
          ];
          buildInputs = [
            pkgs.libxcrypt
            pkgs.nodejs-16_x
            pkgs.python2
            pkgs.python38Full
            pkgs.python310Full
            pkgs.zlib
            pkgs.lua5_1
            pkgs.readline
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

