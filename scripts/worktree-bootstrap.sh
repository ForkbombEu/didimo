#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2025 Forkbomb BV
#
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INCLUDE_FILE="${ROOT_DIR}/.worktreeinclude"
CLASSIC=0

usage() {
	cat <<'USAGE'
Usage: scripts/worktree-bootstrap.sh [--classic] [--from <primary-path>]

Copies allowlisted gitignored paths from the primary worktree, writes
.env.worktree with unique ports (unless --classic), syncs webapp/.env URLs,
initializes submodules, and ensures .bin tools exist.

Primary is detected via `git worktree list` (first entry) unless --from is set.
USAGE
}

FROM_PATH=""
while [[ $# -gt 0 ]]; do
	case "$1" in
	--classic)
		CLASSIC=1
		shift
		;;
	--from)
		FROM_PATH="${2:-}"
		shift 2
		;;
	-h | --help)
		usage
		exit 0
		;;
	*)
		echo "unknown arg: $1" >&2
		usage >&2
		exit 1
		;;
	esac
done

if [[ -z "${FROM_PATH}" ]]; then
	FROM_PATH="$(git -C "${ROOT_DIR}" worktree list --porcelain | awk '/^worktree / { print $2; exit }')"
fi

if [[ -z "${FROM_PATH}" || ! -d "${FROM_PATH}" ]]; then
	echo "error: could not resolve primary worktree path" >&2
	exit 1
fi

if [[ "$(cd "${FROM_PATH}" && pwd -P)" == "$(cd "${ROOT_DIR}" && pwd -P)" ]]; then
	echo "bootstrap on primary worktree: writing classic ports only"
	CLASSIC=1
fi

copy_path() {
	local rel="$1"
	local src="${FROM_PATH}/${rel}"
	local dst="${ROOT_DIR}/${rel}"
	if [[ ! -e "${src}" ]]; then
		echo "skip missing ${rel}"
		return 0
	fi
	if [[ -e "${dst}" ]]; then
		echo "keep existing ${rel}"
		return 0
	fi
	mkdir -p "$(dirname "${dst}")"
	if [[ -d "${src}" ]]; then
		if cp -cR "${src}" "${dst}" 2>/dev/null; then
			echo "cloned ${rel}"
		else
			cp -R "${src}" "${dst}"
			echo "copied ${rel}"
		fi
	else
		cp "${src}" "${dst}"
		echo "copied ${rel}"
	fi
}

echo "bootstrap from ${FROM_PATH} -> ${ROOT_DIR}"

if [[ -f "${INCLUDE_FILE}" ]]; then
	while IFS= read -r line || [[ -n "${line}" ]]; do
		case "${line}" in
		'' | \#*) continue ;;
		esac
		rel="${line%/}"
		rel="${rel#/}"
		copy_path "${rel}"
	done <"${INCLUDE_FILE}"
else
	echo "warning: ${INCLUDE_FILE} missing; copying .env and webapp/.env only"
	copy_path ".env"
	copy_path "webapp/.env"
fi

if [[ "${CLASSIC}" -eq 1 ]]; then
	"${ROOT_DIR}/scripts/worktree-env.sh" write --classic
else
	"${ROOT_DIR}/scripts/worktree-env.sh" write
fi
"${ROOT_DIR}/scripts/worktree-env.sh" sync-urls

git -C "${ROOT_DIR}" submodule update --init --recursive

if [[ ! -x "${ROOT_DIR}/.bin/stepci-captured-runner" || ! -x "${ROOT_DIR}/.bin/et-tu-cesr" ]]; then
	echo "running make tools"
	make -C "${ROOT_DIR}" tools
fi

echo
echo "bootstrap complete. next: make dev"
"${ROOT_DIR}/scripts/worktree-env.sh" print
