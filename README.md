# 📙 GNFinder Plus

A tiny wrapper parses text from document files like `pdf`, `docx` etc. and passes it to [GNFinder](https://github.com/gnames/gnfinder) for extracting names

```sh
go install
go build
```

### Pre-built binary

pre-built binaries by GitHub Actions can be downloaded from [releases](https://github.com/biodiv-platform/gnfinder-plus/releases)

### Usage

#### Web

```sh
# below command starts gnfinder-plus on 3006 port
# different port can be given by -port=3006 etc.
./gnfinder-plus

# in another window
curl http://localhost:3006/parse?file=url
```

#### CLI

```sh
wget https://indiabiodiversity.org/biodiv/content/documents/document-0162468a-7ce7-499e-ac6d-ead2dc273c35/687.pdf
./gnfinder-plus -file=687.pdf
```

### Note

in some cases `pdftotext` binary might be missing please install according to your os

```sh
sudo apt install poppler-utils # debian
```
