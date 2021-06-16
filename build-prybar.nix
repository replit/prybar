{ language, buildInputs ? [] }:

{ lib, buildGoModule, fetchFromGitHub, bash, pkg-config }:

buildGoModule rec {
    name = "prybar-${language}";
    version = "98148ff";

    src = fetchFromGitHub{
        owner = "replit";
        repo = "prybar";
        rev = "${version}";
        sha256 = "14l4bhdlssp22wdxx1ycz9wl64z7bf3qwpqqp3npcva969191c48";
    };

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
