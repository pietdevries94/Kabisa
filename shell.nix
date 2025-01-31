{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    go
  ];
  # disabling CGO happens automatically when building, but is on by default.
  # To prevent unexpected differences, we dissable it in our shell.
  CGO_ENABLED = 0;
  # to enable debugging, we need to disable some nix-specific hardening
  hardeningDisable = [ "fortify" ];
}
