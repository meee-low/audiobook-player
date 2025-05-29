{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
    buildInputs = with pkgs; [
        go
        pkg-config
        alsa-lib
        sqlc
    ];
}
