class ClipbridgeServer < Formula
  desc "Self-hosted clipboard stack service with an embedded Web UI"
  homepage "https://github.com/cinmou/ClipBridgeServer"
  version "0.10.0"

  if Hardware::CPU.arm?
    url "https://github.com/cinmou/ClipBridgeServer/releases/download/v0.10.0/clipbridge-server-darwin-arm64"
    sha256 "REPLACE_WITH_REAL_SHA256"
  else
    url "https://github.com/cinmou/ClipBridgeServer/releases/download/v0.10.0/clipbridge-server-darwin-amd64"
    sha256 "REPLACE_WITH_REAL_SHA256"
  end

  def install
    bin.install Dir["clipbridge-server-darwin-*"].first => "clipbridge-server"
    (etc/"clipbridge").mkpath
    (var/"clipbridge/data").mkpath
    pkgshare.install "configs/config.example.yaml"
  end

  service do
    run [opt_bin/"clipbridge-server", "-config", etc/"clipbridge/config.yaml"]
    keep_alive true
    working_dir var/"clipbridge"
    log_path var/"log/clipbridge-server.log"
    error_log_path var/"log/clipbridge-server.log"
  end

  test do
    assert_match "ClipBridgeServer", shell_output("#{bin}/clipbridge-server -h 2>&1", 0)
  end
end
