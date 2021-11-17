{ language, buildInputs ? [ ], binaries ? [ ], setFlags ? false
, pkgName ? language }:

{ lib, buildGoModule, fetchFromGitHub, bash, expect, pkg-config, runCommand, git
, python3, copyPathToStore, rev, makeWrapper }:

buildGoModule {
  pname = "prybar-${language}";
  version = rev;

  src = ./.;

  inherit buildInputs;
  nativeBuildInputs = [ pkg-config makeWrapper ];

  subPackages = [ "languages/${language}" ];

  vendorSha256 = "1nrycpyjmmhs8cl6bmrajfwczxsv8z20y7a9k7js8p7z60fz7pya";

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

    ${bash}/bin/bash ./scripts/inject.sh ${language}
    go generate ./languages/${language}/main.go
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

    patchShebangs run_no_pty ./tests/${language}/*.exp

    # Currently the go tests won't work in nix because of something w/ nix's sandboxing.
    DISABLE_GO_TESTS=1 "${bash}/bin/bash" "./run_tests_language" "${language}" "${expect}/bin"

    runHook postCheck
  '';

  # Delete the test binary we copied in the preCheck
  postCheck = ''
    rm -f ./prybar-${language}
  '';

  postInstall = ''
    mv $out/bin/${language} $out/bin/prybar-${language}

    if [ -d "./prybar_assets/${language}" ] 
    then
      mkdir -p "$out/prybar_assets/${language}"
      cp -R "./prybar_assets/${language}" "$out/prybar_assets/"

      wrapProgram "$out/bin/prybar-${language}" \
        --set PRYBAR_ASSETS_DIR "$out/prybar_assets"
    fi
  '';
}
