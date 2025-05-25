{
  description = "aerospace-ipc: Socket Client for AeroSpace WM";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, utils, ... }: 
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go

            # To create new subcommands, run:
            # cobra-cli add <subcommand-name>
            cobra-cli

            # To generate the mock for the interfaces, run:
            # mockgen -source=./pkg/cli/cli.go -destination=./pkg/cli/mock/mock_cli.go -package=mock
            mockgen

            # To connect with aerspace socket
            # nc -U /tmp/aerospace.sock < fixture/foobar.json
            # netcat-openbsd
          ];
        };

        packages = {
          default = pkgs.callPackage ./nix/package-default.nix {};
          nightly = pkgs.callPackage ./nix/package-nightly.nix {};
          source = pkgs.callPackage ./nix/package-source.nix {};
        };
    });
}
