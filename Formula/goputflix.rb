# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Goputflix < Formula
  desc "A CLI tool to browse and stream videos from your Put.io account directly to VLC media player"
  homepage "https://github.com/alexraskin/goputflix"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/alexraskin/goputflix/releases/download/v1.0.0/goputflix_Darwin_x86_64.tar.gz"
      sha256 "916277931dbb024b9607468493d7c037d8caba10085d932a07d735f090c07c04"

      def install
        bin.install "goputflix"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/alexraskin/goputflix/releases/download/v1.0.0/goputflix_Darwin_arm64.tar.gz"
      sha256 "aa8a6499503c7f3e58c15b8415742bfc2c264621786a0974db2705e6bc263bd1"

      def install
        bin.install "goputflix"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/alexraskin/goputflix/releases/download/v1.0.0/goputflix_Linux_x86_64.tar.gz"
        sha256 "812bbce9ca29615e4c96d482558cb5ddd13bfb29c2b7c8d205ab6e27f8d5ce64"

        def install
          bin.install "goputflix"
        end
      end
    end
    if Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/alexraskin/goputflix/releases/download/v1.0.0/goputflix_Linux_arm64.tar.gz"
        sha256 "b3194c8745769a62883791df0d37f111c70d6022ccffd87941ec6071ceabbf61"

        def install
          bin.install "goputflix"
        end
      end
    end
  end
end
