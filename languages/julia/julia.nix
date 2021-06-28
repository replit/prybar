{ pkgs }:

pkgs.stdenv.mkDerivation rec {
    pname = "julia-sources";
    version = "1.5.4";

    src = {
        x86_64-linux = pkgs.fetchurl {
            url = "https://julialang-s3.julialang.org/bin/linux/x64/${pkgs.lib.versions.majorMinor version}/julia-${version}-linux-x86_64.tar.gz";
            sha256 = "1icb3rpn2qs6c3rqfb5rzby1pj6h90d9fdi62nnyi4x5s58w7pl0";
        };
    }.${pkgs.stdenv.hostPlatform.system} or (throw "Unsupported system: ${pkgs.stdenv.hostPlatform.system}");

    installPhase = ''
        cp -r . $out

        mkdir $out/lib/pkgconfig
        cat > $out/lib/pkgconfig/julia.pc <<- EOM
            Name: julia
            Cflags: -I$out/include/julia
            Version: 1.5.4
            Description: it's julia
            Libs: -L$out\/lib -ljulia
        EOM
    '';
}