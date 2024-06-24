{ language, target ? language, buildInputs ? [ ], binaries ? [ ], setFlags ? false
, pkgName ? language, cgoPkgs ? null, cgoExtraCflags ? "" }:

{ lib, buildGoModule, fetchFromGitHub, bash, expect, pkg-config, runCommand, git
, python3, copyPathToStore, rev, makeWrapper }:

buildGoModule {
  pname = "prybar-${target}";
  version = rev;

  src = ./.;

  inherit buildInputs;
  nativeBuildInputs = [ pkg-config makeWrapper ];

  subPackages = [ "languages/${language}" ];

  vendorHash = "sha256-yt/zHTD/XKTlmUkdD8RHW/fPuJMq12UoQxrWKv1lPts=";

  # This prebiuild hook will setup the compiler flags on demand based on the package
  # If a language requires this, it MUST expost a ${pkgName}.pc file in its `lib/pkg-config`
  # for pkg-config to detect it. Otherwise, we can't find the required flags.
  preBuild = ''
    export CGO_LDFLAGS_ALLOW="-Wl,--compress-debug-sections=zlib"
    ${if setFlags then ''
      export NIX_CFLAGS_COMPILE="$(pkg-config --cflags ${pkgName}) $NIX_CFLAGS_COMPILE"
      export NIX_LDFLAGS="$(pkg-config --libs-only-L  ${pkgName}) $(pkg-config --libs-only-l  ${pkgName}) $NIX_LDFLAGS"
    '' else
      ""}
    
    ${if cgoPkgs != null then ''
    export CGO_CFLAGS="$(pkg-config --cflags ${cgoPkgs}) ${cgoExtraCflags}"
    export CGO_LDFLAGS="$(pkg-config --libs ${cgoPkgs})"
    '' else ""}
    
    ${bash}/bin/bash ./scripts/inject.sh ${language}
    go generate ./languages/${language}/main.go

    if [[ -f "languages/${language}/compile" ]]; then
      patchShebangs languages/${language}/compile
      languages/${language}/compile;
    fi
  '';

  # The test file expect the binary to be in the current directory rather than bin
  preCheck = ''
    cp $GOPATH/bin/${language} ./prybar-${language}
  '';

  # Add all the dependencies to the check. Binaries should include the language interpreter/binaries
  # for testing the language.
  checkInputs = [ expect ] ++ buildInputs ++ binaries;

  # Run the end-to-end tests for this specific language
  checkPhase = ''
    runHook preCheck

    # Nix unsets HOME but some tests rely on it, let's set it to whatever tmp dir the check
    # is running in.
    export HOME=$(echo pwd)
    export PATH=$PATH:${lib.makeBinPath ([expect] ++ binaries)}

    # Currently the go tests won't work in nix because of something w/ nix's sandboxing.

    DISABLE_GO_TESTS=1 "${bash}/bin/bash" "./run_tests_language" "${language}"

    runHook postCheck
  '';

  # Delete the test binary we copied in the preCheck
  postCheck = ''
    rm -f ./prybar-${language}
  '';

  postInstall = ''
    mv $out/bin/${language} $out/bin/prybar-${target}

    if [ -d "./prybar_assets/${language}" ] 
    then
      mkdir -p "$out/prybar_assets/${language}"
      cp -R "./prybar_assets/${language}" "$out/prybar_assets/"

      wrapProgram "$out/bin/prybar-${target}" \
        --set PRYBAR_ASSETS_DIR "$out/prybar_assets"
    fi
  '';
}
