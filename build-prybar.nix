{ language, buildInputs ? [] }:

{ lib, buildGoModule, fetchFromGitHub, bash, pkg-config, runCommand, git }:

let
    src = ./.;

    gitSrc = builtins.filterSource
                   (path: type: true)
                   ./.;
    revision = runCommand "get-rev" {
        nativeBuildInputs = [ git ];
        dummy = builtins.currentTime;
    } "GIT_DIR=${gitSrc}/.git git rev-parse --short HEAD | tr -d '\n' > $out";
in buildGoModule {
    pname = "prybar-${language}";
    version = builtins.readFile revision;

    inherit src;

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
