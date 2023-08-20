pkgname=goldwarden
pkgver=0.1.1
pkgrel=1
pkgdesc='Goldwarden'
arch=('x86_64')
url="https://github.com/quexten/$pkgname"
license=('MIT')
makedepends=('go' 'libfido2' 'gcc' 'wayland' 'libx11' 'libxkbcommon' 'libxkbcommon-x11' 'libxcursor' 'base-devel' 'vulkan-headers')
source=("$url/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('7d38db887437a58758e5f183d4951cf7c4d1b099f37ff6f5e95fb98735634983')

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
  go build -tags autofill -o build/$pkgname .
}

package() {
  cd "$pkgname-$pkgver"
  echo $pkgdir
  install -Dm755 build/$pkgname "$pkgdir"/usr/bin/$pkgname
}