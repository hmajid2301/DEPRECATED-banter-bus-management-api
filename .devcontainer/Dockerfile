# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.163.1/containers/go/.devcontainer/base.Dockerfile

ARG GO_VERSION="1"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${GO_VERSION}

COPY library-scripts/fish-debian.sh /tmp/library-scripts/
COPY postInst.fish postInst.bash /home/vscode/

RUN apt-get update && \
  bash /tmp/library-scripts/fish-debian.sh && \
  curl -fsSL https://starship.rs/install.sh | bash -s -- --yes

COPY ./config.fish /home/vscode/.config/fish/config.fish
COPY ./.bashrc /home/vscode/.bashrc