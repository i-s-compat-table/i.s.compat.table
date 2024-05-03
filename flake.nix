{
  description = "`information_schema` compatibility tables";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils"; # TODO: pin
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
    # FIXME: eliminate reliance on GitHub
  };

  outputs = { self, flake-utils, nixpkgs }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = (import nixpkgs) { inherit system; };
      in {
        devShell = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go # 1.21
          ];
          buildInputs = with pkgs; [
            # go development
            gopls
            delve
            golangci-lint
            
            sqlite
            podman
            podman-compose
            
            # nix support
            nixpkgs-fmt
            nil
            # goreleaser # TODO: when we have a release

            # general development
            git
            bashInteractive
            lychee
            shellcheck
            fzf
            ripgrep

            # # node development
            # nodejs
            # nodePackages.pnpm
          ];
        };
      }
    );
}
