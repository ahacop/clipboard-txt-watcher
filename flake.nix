{
  description = "clipboard-txt-watcher";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    let
      # Package builder function that can be used with any pkgs
      mkPackage = pkgs: pkgs.buildGoModule rec {
        pname = "clipboard-txt-watcher";
        baseVersion = builtins.replaceStrings [ "\n" ] [ "" ] (builtins.readFile ./VERSION);
        version = if (self ? shortRev) then "${baseVersion}-${self.shortRev}" else baseVersion;

        src = ./.;

        vendorHash = "sha256-CKFZ/6CMM3C0QJYKNMpIHGXIcCxMopeZi9zxIplTQ10=";

        ldflags = [
          "-s"
          "-w"
          "-X main.version=${version}"
          "-X main.commit=${self.rev or "unknown"}"
        ];

        env.CGO_ENABLED = 0;

        nativeBuildInputs = [ pkgs.makeWrapper ];

        postInstall = ''
          wrapProgram $out/bin/clipboard-txt-watcher \
            --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.wl-clipboard pkgs.xclip ]}
        '';

        meta = with pkgs.lib; {
          description = "Watch a file and sync its contents to the system clipboard";
          homepage = "https://github.com/ahacop/clipboard-txt-watcher";
          license = licenses.gpl3Plus;
          maintainers = [
            {
              name = "Ara Hacopian";
              github = "ahacop";
            }
          ];
        };
      };
    in
    {
      # Overlay for use in NixOS configurations
      overlays.default = final: prev: {
        clipboard-txt-watcher = mkPackage final;
      };

      # Home-manager module
      homeManagerModules.default = { config, lib, pkgs, ... }:
        let
          cfg = config.services.clipboard-txt-watcher;
        in
        {
          options.services.clipboard-txt-watcher = {
            enable = lib.mkEnableOption "clipboard-txt-watcher service";

            watchFile = lib.mkOption {
              type = lib.types.path;
              description = "Path to the file to watch";
            };

            clipboardBackend = lib.mkOption {
              type = lib.types.enum [ "wayland" "x11" ];
              default = "wayland";
              description = "Clipboard backend to use";
            };

            package = lib.mkOption {
              type = lib.types.package;
              default = mkPackage pkgs;
              defaultText = lib.literalExpression "pkgs.clipboard-txt-watcher";
              description = "The clipboard-txt-watcher package to use";
            };
          };

          config = lib.mkIf cfg.enable {
            systemd.user.services.clipboard-txt-watcher = {
              Unit = {
                Description = "Clipboard Text File Watcher";
                After = [ "graphical-session.target" ];
              };

              Service = {
                Type = "simple";
                ExecStart = "${cfg.package}/bin/clipboard-txt-watcher --file ${cfg.watchFile} --backend ${cfg.clipboardBackend}";
                Restart = "on-failure";
                RestartSec = 5;
              };

              Install = {
                WantedBy = [ "default.target" ];
              };
            };
          };
        };
    }
    //
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = mkPackage pkgs;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gofumpt
            golangci-lint
            goreleaser
            wl-clipboard
            xclip
          ];
        };
      }
    );
}
