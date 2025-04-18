{ pkgs, lib, config, inputs, ... }:

{
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = [
    pkgs.git
  ];

  dotenv.enable = true;

  # https://devenv.sh/languages/
  # languages.rust.enable = true;
  languages.go.enable = true;
  languages.javascript.enable = true;

  # https://devenv.sh/processes/
  # processes.cargo-watch.exec = "cargo-watch";
  processes.api.exec = "cd api && go run main.go";

  # https://devenv.sh/services/
  # services.postgres.enable = true;
  services.postgres = {
    enable = true;
    initialDatabases = [{
      name = "swipcord";
      user = "swipcord";
      pass = "swipcord";
      }];
    listen_addresses = "0.0.0.0";
  };

  # https://devenv.sh/scripts/
  scripts.hello.exec = ''
    echo hello from $GREET
  '';
  env = {
    DB_NAME = "swipcord";
    DB_USER = "swipcord";
    DB_PASS = "swipcord";
    DB_HOST = "localhost";
    DB_PORT = "5432";
    JWT_SECRET = "secret";
  };
  enterShell = ''
    hello
    git --version
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  # https://devenv.sh/git-hooks/
  # git-hooks.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}
