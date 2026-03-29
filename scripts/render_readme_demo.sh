#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

mkdir -p assets
rm -f assets/pls-demo.webm assets/pls-demo.gif assets/pls-demo-palette.png

ngo=${GO:-go}
"$ngo" build -o bin/pls ./cmd/pls
VHS_NO_SANDBOX=1 vhs demo/readme.tape

ffmpeg -y -i assets/pls-demo.webm -vf "fps=12,scale=1200:-1:flags=lanczos,palettegen" assets/pls-demo-palette.png
ffmpeg -y -i assets/pls-demo.webm -i assets/pls-demo-palette.png -lavfi "fps=12,scale=1200:-1:flags=lanczos[x];[x][1:v]paletteuse" assets/pls-demo.gif
rm -f assets/pls-demo-palette.png

ls -lh assets/pls-demo.webm assets/pls-demo.gif
