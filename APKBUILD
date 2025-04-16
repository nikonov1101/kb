# Contributor: xandrow indira <xandrow89@gmail.com> 
# Maintainer: Alex Nikonov <alex@nikonov.tech>
pkgname=kb
pkgver=1.0.3
pkgrel=1
pkgdesc="a simple static site generator"
url="https://github.com/nikonov1101/kb"
arch="all"
license="MIT"
makedepends="go make"
subpackages="$pkgname-doc"
source="$pkgname-$pkgver.tar.gz::https://github.com/nikonov1101/kb/archive/v$pkgver.tar.gz"

builddir="$srcdir/$pkgname-$pkgver/"



build() {
	export GOCACHE="${GOCACHE:-"$srcdir/go-cache"}"
	export GOTMPDIR="${GOTMPDIR:-"$srcdir"}"
	export GOMODCACHE="${GOMODCACHE:-"$srcdir/go"}"
	VERSION="$pkgver-apkbuild" make build
}

check() {
	# Remove and add !check option if there is no check command.
	./kb version
}

package() {
	install -m755 -D kb -t "$pkgdir"/usr/bin/
	install -Dm644 LICENSE "$pkgdir"/usr/share/licenses/$pkgname/LICENSE
}

sha512sums="
adf3ccd2d74aa9c32c2db28e4d17be89dd19e81830ba7a3ffed28045f3d38a2b3f5dbcbe5357d880d73d0dd4fd2eb4910a84a614fe8d7866f9b0ce7c67361690  kb-1.0.3.tar.gz
"
