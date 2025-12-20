#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
if [[ -z "${version}" ]]; then
  echo "Usage: $0 <version>" >&2
  exit 1
fi

url="https://github.com/steipete/ordercli/archive/refs/tags/v${version}.tar.gz"
tmp="/tmp/ordercli-${version}.tar.gz"

curl -L -o "${tmp}" "${url}" >/dev/null
sha256="$(shasum -a 256 "${tmp}" | awk '{print $1}')"

cat <<EOF
version "${version}"
url "${url}"
sha256 "${sha256}"
EOF
