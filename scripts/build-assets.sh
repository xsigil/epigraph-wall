#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TEX_DIR="${ROOT_DIR}/tex"
OUT_DIR="${ROOT_DIR}/images"
SCREEN_RES="${1:-1920x1080}"
TEMPLATE="${TEX_DIR}/template.tex"

mkdir -p "${OUT_DIR}"

TMP_DIR="$(mktemp -d --tmpdir="${XDG_RUNTIME_DIR:-/tmp}" epigraph-build-XXXXXX)"
trap 'rm -rf "${TMP_DIR}"' EXIT

for src in "${TEX_DIR}"/*.tex; do
    filename="$(basename "$src")"
    [[ "$filename" == "template.tex" ]] && continue

    name="${filename%.tex}"
    echo "[*] Processing ${name} with LuaLaTeX..."

    # バックスラッシュを壊さずに __BODY__ の位置にソースを挿入
    target_tex="${TMP_DIR}/${name}.tex"
    while IFS= read -r line || [[ -n "$line" ]]; do
        if [[ "$line" =~ __BODY__ ]]; then
            cat "$src"
        else
            printf '%s\n' "$line"
        fi
    done < "${TEMPLATE}" > "${target_tex}"

    # LuaLaTeX コンパイル (エラー時はログを出力して原因特定できるようにする)
    if ! (cd "${TMP_DIR}" && lualatex -interaction=nonstopmode "${name}.tex" > lualatex.log 2>&1); then
        echo "[-] Error: LuaLaTeX compilation failed for ${name}.tex" >&2
        cat "${TMP_DIR}/lualatex.log" | tail -n 25 >&2
        exit 1
    fi

    # PDF -> PNG (300 DPI)
    pdftoppm -png -r 300 "${TMP_DIR}/${name}.pdf" "${TMP_DIR}/page" >/dev/null 2>&1

    # 白黒反転 (-negate) して画面解像度の黒キャンバス中央に配置
    magick -size "${SCREEN_RES}" xc:"#000000" \
        \( "${TMP_DIR}/page-1.png" -negate \) \
        -gravity center -composite "${OUT_DIR}/${name}.png"

    echo "[+] Generated: images/${name}.png"
    rm -f "${TMP_DIR}"/*
done

echo "[✓] All epigraphs successfully built."
