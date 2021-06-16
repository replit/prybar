{ language, buildInputs ? [] }:

{ lib, buildGoModule, fetchFromGitHub, bash, pkg-config, runCommand, git }:

let
    gitSrc = builtins.filterSource
               (path: type: true)
               ./.;
in
buildGoModule rec {
    pname = "prybar-${language}";

    revision = runCommand "get-rev" {
        nativeBuildInputs = [ git ];
        dummy = builtins.currentTime;
    } "GIT_DIR=${gitSrc}/.git git rev-parse --short HEAD | tr -d '\n' > $out";
    version = builtins.readFile revision;

    src = ./.;

    inherit buildInputs;
    nativeBuildInputs = [ pkg-config ];

    subPackages = [ "languages/${language}" ];

    vendorSha256 = null;

    preBuild = ''
        ${bash}/bin/bash ./scripts/inject.sh ${language}
        go generate ./languages/${language}/main.go
    '';

    postInstall = ''
        mv $out/bin/${language} $out/bin/prybar-${language}
    '';
}
