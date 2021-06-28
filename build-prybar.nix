{ language, buildInputs ? [], binaries ? [], setFlags ? false, pkgName ? language, versionArg ? "" }:

{ lib, buildGoModule, fetchFromGitHub, bash, expect, pkg-config, runCommand, git, python3, copyPathToStore }:

let
    src = copyPathToStore ./.;
    revision = runCommand "get-rev" {
        nativeBuildInputs = [ git python3 ];
        dummy = builtins.currentTime;
    } ''
        if [ -d ${src}/.git ]; then
            cd ${src}
            git rev-parse --short HEAD | tr -d '\n' > $out
        else
            echo ${versionArg} | tr -d '\n' > $out
        fi
    '';
in buildGoModule {
    pname = "prybar-${language}";
    version = builtins.readFile revision;

    inherit src;

    inherit buildInputs;
    nativeBuildInputs = [ pkg-config ];

    subPackages = [ "languages/${language}" ];

    vendorSha256 = null;

    # This prebiuild hook will setup the compiler flags on demand based on the package
    # If a language requires this, it MUST expost a ${pkgName}.pc file in its `lib/pkg-config`
    # for pkg-config to detect it. Otherwise, we can't find the required flags.
    preBuild = ''
        export CGO_LDFLAGS_ALLOW="-Wl,--compress-debug-sections=zlib"
        ${if setFlags then ''
            export NIX_CFLAGS_COMPILE="$(pkg-config --cflags ${pkgName}) $NIX_CFLAGS_COMPILE"
            export NIX_LDFLAGS="$(pkg-config --libs-only-L  ${pkgName}) $(pkg-config --libs-only-l  ${pkgName}) $NIX_LDFLAGS"
        '' else ""}

        ${bash}/bin/bash ./scripts/inject.sh ${language}
        go generate ./languages/${language}/main.go
    '';

    # The test file expect the binary to be in the current directory rather than bin
    preCheck = ''
        cp $GOPATH/bin/${language} ./prybar-${language}
        ln -s ${expect} /usr
    '';

    # Add all the dependencies to the check. Binaries should include the language interpreter/binaries
    # for testing the language.
    checkInputs = [ expect ] ++ buildInputs ++ binaries;

    # Run the end-to-end tests for this specific language
    checkPhase = ''
        runHook preCheck

        ${bash}/bin/bash ./run_tests_language ${language}

        runHook postCheck
    '';

    # Delete the test binary we copied in the preCheck
    postCheck = ''
        rm -f ./prybar-${language}
    '';

    postInstall = ''
        mv $out/bin/${language} $out/bin/prybar-${language}
    '';
}
