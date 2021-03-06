# generated by clj2nix-1.0.6
{ pkgs ? import <nixpkgs> {} }:

let repos = [
        "http://jcenter.bintray.com"
        "https://repo1.maven.org/maven2/"
        "https://repo.clojars.org/" ];

  in rec {
      makePaths = {extraClasspaths ? null}:
        (pkgs.lib.concatMap
          (dep:
            builtins.map
            (path:
              if builtins.isString path then
                path
              else if builtins.hasAttr "jar" path then
                path.jar
              else if builtins.hasAttr "outPath" path then
                path.outPath
              else
                path
                )
            dep.paths)
          packages)
        ++ (if extraClasspaths != null then [ extraClasspaths ] else []);
      makeClasspaths = {extraClasspaths ? null}: builtins.concatStringsSep ":" (makePaths {inherit extraClasspaths;});
      packageSources = builtins.map (dep: dep.src) packages;

      packages = [
      rec {
        name = "clojure/org.clojure";
        src = pkgs.fetchMavenArtifact {
          inherit repos;
          artifactId = "clojure";
          groupId = "org.clojure";
          sha512 = "4bb567b9262d998f554f44e677a8628b96e919bc8bcfb28ab2e80d9810f8adf8f13a8898142425d92f3515e58c57b16782cff12ba1b5ffb38b7d0ccd13d99bbc";
          version = "1.10.3";
        };
        paths = [ src ];
      }
      rec {
        name = "spec.alpha/org.clojure";
        src = pkgs.fetchMavenArtifact {
          inherit repos;
          artifactId = "spec.alpha";
          groupId = "org.clojure";
          sha512 = "0740dc3a755530f52e32d27139a9ebfd7cbdb8d4351c820de8d510fe2d52a98acd6e4dfc004566ede3d426e52ec98accdca1156965218f269e60dd1cd4242a73";
          version = "0.2.194";
        };
        paths = [ src ];
      }
      rec {
        name = "core.specs.alpha/org.clojure";
        src = pkgs.fetchMavenArtifact {
          inherit repos;
          artifactId = "core.specs.alpha";
          groupId = "org.clojure";
          sha512 = "c1d2a740963896d97cd6b9a8c3dcdcc84459ea66b44170c05b8923e5fbb731b4b292b217ed3447bbc9e744c9a496552f77a6c38aea232e5e69f8faa627dea4b5";
          version = "0.2.56";
        };
        paths = [ src ];
      }
  ];
}