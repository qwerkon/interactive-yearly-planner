#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACKAGE_NAME="${PACKAGE_NAME:-planner}"
VERSION="${VERSION:-}"
ARCH="${ARCH:-}"
MAINTAINER="${MAINTAINER:-Planner Maintainers <maintainers@example.com>}"
DESCRIPTION="PDF planner generator for e-ink devices"

if [[ -z "${VERSION}" ]]; then
  if git -C "${ROOT_DIR}" describe --tags --always --dirty >/dev/null 2>&1; then
    VERSION="$(git -C "${ROOT_DIR}" describe --tags --always --dirty)"
    VERSION="${VERSION#v}"
  else
    VERSION="0.0.0"
  fi
fi
VERSION="${VERSION//-/.}"

if [[ -z "${ARCH}" ]]; then
  ARCH="$(dpkg --print-architecture)"
fi

case "${ARCH}" in
  amd64) GOARCH="amd64" ;;
  arm64) GOARCH="arm64" ;;
  armhf) GOARCH="arm"; GOARM="7" ;;
  *)
    echo "Unsupported Debian architecture: ${ARCH}" >&2
    echo "Set ARCH to one of: amd64, arm64, armhf" >&2
    exit 1
    ;;
esac

if ! command -v go >/dev/null 2>&1; then
  echo "go is required to build ${PACKAGE_NAME}" >&2
  exit 1
fi

if ! command -v dpkg-deb >/dev/null 2>&1; then
  echo "dpkg-deb is required to build the .deb package" >&2
  exit 1
fi

BUILD_ROOT="${ROOT_DIR}/dist/deb"
PKG_ROOT="${BUILD_ROOT}/${PACKAGE_NAME}_${VERSION}_${ARCH}"
BIN_PATH="${PKG_ROOT}/usr/lib/${PACKAGE_NAME}/${PACKAGE_NAME}"
WRAPPER_PATH="${PKG_ROOT}/usr/bin/${PACKAGE_NAME}"
RESOURCE_DIR="${PKG_ROOT}/usr/share/${PACKAGE_NAME}"
CONTROL_DIR="${PKG_ROOT}/DEBIAN"
OUTPUT="${BUILD_ROOT}/${PACKAGE_NAME}_${VERSION}_${ARCH}.deb"

rm -rf "${PKG_ROOT}" "${OUTPUT}"
mkdir -p "$(dirname "${BIN_PATH}")" "$(dirname "${WRAPPER_PATH}")" "${RESOURCE_DIR}" "${CONTROL_DIR}" "${BUILD_ROOT}"

echo "Building ${PACKAGE_NAME} ${VERSION} for linux/${GOARCH}..."
(
  cd "${ROOT_DIR}"
  CGO_ENABLED=0 GOOS=linux GOARCH="${GOARCH}" GOARM="${GOARM:-}" go build -trimpath -ldflags "-s -w" -o "${BIN_PATH}" ./cmd/planner
)

install -m 0755 /dev/stdin "${WRAPPER_PATH}" <<'EOF'
#!/usr/bin/env sh
set -eu

export PLANNER_RESOURCE_DIR="${PLANNER_RESOURCE_DIR:-/usr/share/planner}"
exec /usr/lib/planner/planner "$@"
EOF

cp -a "${ROOT_DIR}/cfg" "${RESOURCE_DIR}/cfg"
cp -a "${ROOT_DIR}/tpls" "${RESOURCE_DIR}/tpls"
cp -a "${ROOT_DIR}/translations" "${RESOURCE_DIR}/translations"
install -m 0644 "${ROOT_DIR}/README.md" "${RESOURCE_DIR}/README.md"
install -m 0644 "${ROOT_DIR}/LICENSE" "${RESOURCE_DIR}/LICENSE"

INSTALLED_SIZE="$(du -sk "${PKG_ROOT}/usr" | cut -f1)"

cat >"${CONTROL_DIR}/control" <<EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: ${MAINTAINER}
Depends: ca-certificates, texlive-xetex, texlive-latex-extra, texlive-lang-polish
Installed-Size: ${INSTALLED_SIZE}
Description: ${DESCRIPTION}
 Builds hyperlinked yearly PDF planners from YAML profiles.
EOF

find "${PKG_ROOT}" -type d -exec chmod 0755 {} +
find "${PKG_ROOT}" -type f -exec chmod 0644 {} +
chmod 0755 "${BIN_PATH}" "${WRAPPER_PATH}"

dpkg-deb --build --root-owner-group "${PKG_ROOT}" "${OUTPUT}"

echo "Created ${OUTPUT}"
echo "Install with: sudo apt install ${OUTPUT}"
