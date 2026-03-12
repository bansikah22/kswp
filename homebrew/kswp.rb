class Kswp < Formula
  desc "A kubectl plugin to scan and sweep unused resources"
  homepage "https://github.com/bansikah22/kswp"
  url "https://github.com/bansikah22/kswp/archive/v0.1.0.tar.gz"
  sha256 "..."
  license "Apache-2.0"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"kswp", "."
  end

  test do
    system "#{bin}/kswp", "--version"
  end
end
