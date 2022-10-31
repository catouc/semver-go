{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils/v1.0.0";
  };

  description = "A small CLI to fish out the current or next semver version from a git repository";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem ( system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        build = pkgs.buildGoModule {
          pname = "semver-go";
          version = "v0.4.0";
          modSha256 = pkgs.lib.fakeSha256;
          vendorSha256 = null;
          src = ./.;
          nativeBuildInputs = [ pkgs.git ];
        };
      in
        rec {
        packages = {
          semver-go = build;
          default = build;
        };

        devShells = {
          default = pkgs.mkShell {
            buildInputs = [
              pkgs.go
            ];
          };
        };
      }
    );
}