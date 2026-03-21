let
  # Pin nixpkgs to a specific release for reproducibility
  nixpkgs = builtins.fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/fcd5e1ec8b8a1413ebdad05e84e90869eb34976e.tar.gz";
    sha256 = "sha256:0pcszxn8kplmb29x5h0v6lzalgwm0rwps9n7626j37vh3kbf2jqp";
  };

  pkgs = import nixpkgs { config.allowUnfree = true; };

in
pkgs.mkShell {
  name = "secure-terraform-pipeline-shell";
  pure = true;

  buildInputs =
      [
        pkgs.git
        pkgs.go
      ];

  shellHook = ''
    
    echo "======================================="
  '';
}