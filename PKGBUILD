pkgname=goldwarden
pkgver=0.2.7
pkgrel=1
pkgdesc='A feature-packed Bitwarden compatible desktop integration'
arch=('x86_64' 'aarch64')
url="https://github.com/quexten/$pkgname"
license=('MIT')
depends=('libfido2')
makedepends=('go' 'gcc' 'base-devel')
source=("$url/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('57555dab4afd4fc60785e8ad34f41932988b4cd2ce130daaa719625a0e455481')

prepare(){
  cd "$pkgname-$pkgver"
  mkdir -p build/
}

build() {
  cd "$pkgname-$pkgver"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -modcacherw"
  export CGO_ENABLED=1

  go mod tidy
  go build -o build/$pkgname .
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 build/$pkgname "$pkgdir"/usr/bin/$pkgname
  install -Dm644 "$srcdir/$pkgname-$pkgver/resources/com.quexten.goldwarden.policy" "$pkgdir/usr/share/polkit-1/actions/com.quexten.goldwarden.policy"
  chown root:root "$pkgdir/usr/share/polkit-1/actions/com.quexten.goldwarden.policy"
}
