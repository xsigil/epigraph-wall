#!/usr/bin/env bash
set -euo pipefail

DEST_DIR="$(cd "$(dirname "$0")/../images" && pwd)"
mkdir -p "${DEST_DIR}"

TMP_DIR="$(mktemp -d --tmpdir="${XDG_RUNTIME_DIR:-/tmp}" epigraph-gen-XXXXXX)"
trap 'rm -rf "${TMP_DIR}"' EXIT

SCREEN_RES="${1:-1920x1080}"

build_entry() {
    local name="$1"
    local main_text="$2"
    local sub_text="$3"

    cat << LATEX > "${TMP_DIR}/doc.tex"
\documentclass[varwidth=22cm, border=3cm]{standalone}
\usepackage[utf8]{inputenc}
\usepackage{yfonts}
\begin{document}
\begin{center}
    {\LARGE \textswab{${main_text}}}

    \vspace{0.8cm}

    {\large \textfrak{${sub_text}}}
\end{center}
\end{document}
LATEX

    (cd "${TMP_DIR}" && pdflatex -interaction=nonstopmode doc.tex >/dev/null 2>&1)
    pdftoppm -png -r 300 "${TMP_DIR}/doc.pdf" "${TMP_DIR}/page" >/dev/null 2>&1
    
    magick -size "${SCREEN_RES}" xc:"#000000" \
        \( "${TMP_DIR}/page-1.png" -negate \) \
        -gravity center -composite "${DEST_DIR}/${name}.png"

    echo "[+] Generated: images/${name}.png"
    rm -f "${TMP_DIR}"/*
}

echo "[*] Generating initial epigraphs (${SCREEN_RES})..."

build_entry "01_sovereignty" \
    "Sovereignty begins in the shell." \
    "--- The Fortress Epigraph"

build_entry "02_servitude" \
    "Knowledge without autonomy is nothing but servitude." \
    "--- Sovereign Axiom"

build_entry "03_silentium" \
    "Silentium est aurum. Libertas est omnia." \
    "--- Ancient Proverb"

echo "[✓] Assets generated successfully."
