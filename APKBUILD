# Contributor: xandrow indira <xandrow89@gmail.com>
# Maintainer: Alex Nikonov <alex@nikonov.tech>
pkgname=kb
pkgver=1.0.3
pkgrel=2
pkgdesc="a simple static site generator"
url="https://github.com/nikonov1101/kb"
arch="all"
license="MIT"
makedepends="go make"
# check: no test suite provided
# net: allow "go build" to download dependencies from github
options="!check net"
builddir="$srcdir/$pkgname-$pkgver/"
source="$pkgname-$pkgver.tar.gz::https://github.com/nikonov1101/kb/archive/v$pkgver.tar.gz"


export GOCACHE="${GOCACHE:-"$srcdir/go-cache"}"
export GOTMPDIR="${GOTMPDIR:-"$srcdir"}"
export GOMODCACHE="${GOMODCACHE:-"$srcdir/go"}"


build() {
	# Replace with proper build command(s).
	# Remove if there is no build command.
	VERSION="$pkgver-apkbuild" make build
}


package() {
	install -m755 -D kb -t "$pkgdir"/usr/bin/
}

sha512sums="
adf3ccd2d74aa9c32c2db28e4d17be89dd19e81830ba7a3ffed28045f3d38a2b3f5dbcbe5357d880d73d0dd4fd2eb4910a84a614fe8d7866f9b0ce7c67361690  kb-1.0.3.tar.gz
"
