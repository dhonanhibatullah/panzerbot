#!/usr/bin/env bash
set -euo pipefail

CARD_NAME="USB Audio Device"
ASOUNDRC="$HOME/.asoundrc"

# aplay -l line looks like: "card 4: Device [USB Audio Device], device 0: ..."
# Field 3 is the short hw name used in hw:NAME,0
HW_NAME=$(aplay -l 2>/dev/null | grep -F "[$CARD_NAME]" | head -n1 | awk '{print $3}')

if [[ -z "$HW_NAME" ]]; then
    echo "Error: could not find ALSA card matching '$CARD_NAME'" >&2
    exit 1
fi

cat > "$ASOUNDRC" << EOF
pcm.!default {
    type plug
    slave.pcm "hw:$HW_NAME,0"
}
ctl.!default {
    type hw
    card $HW_NAME
}
EOF

echo "Written $ASOUNDRC:"
cat "$ASOUNDRC"
echo
echo "Testing with speaker-test (Ctrl+C to stop)..."
speaker-test -c 2 -t wav
