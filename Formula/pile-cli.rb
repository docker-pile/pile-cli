class Cli < Formula
    desc "Docker management CLI"
    homepage "https://github.com/docker-pile/pile-cli"
    url "https://github.com/docker-pile/pile-cli/releases/download/v0.0.1/cli-darwin-amd64"
    sha256 "a71d7b90b54da6e751372869a1c2e76d"
    license "MIT"
  
    def install
      bin.install "cli-darwin-amd64" => "pile"
    end
  
    test do
      system "#{bin}/pile", "--help"
    end
  end
  