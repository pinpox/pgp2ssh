{
  description = "Convert GPG/PGP Keys to SSH private keys";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let

      # to work with older version of flakes
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 lastModifiedDate;

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in
    {

      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          pgp2ssh = pkgs.buildGoModule {
            pname = "pgp2ssh";
            inherit version;
            src = ./.;
            vendorHash = "sha256-O4AeSfdJxSGnWwRkNnAQMnOZE+Auy+3BIjncG/PK5EE=";
          };
        });

      # The default package for 'nix build' and 'nix run'
      defaultPackage = forAllSystems (system: self.packages.${system}.pgp2ssh);
    };
}
