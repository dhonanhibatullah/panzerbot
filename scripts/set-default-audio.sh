#!/usr/bin/env bash
set -euo pipefail

CARD_NAME="USB Audio Device"
ASOUNDRC="$HOME/.asoundrc"

# Find the ALSA card identifier (e.g. "Device" from "USB Audio Device")
# aplay -l uses the short name in hw:NAME,0 — get it from /proc/asound
CARD_ID=$(aplay -l 2>/dev/null | grep -F "[$CARD_NAME]" | head -n1 | awk '{print $2}' | tr -d ':')

if [[ -z "$CARD_ID" ]]; then
    echo "Error: could not find ALSA card matching '$CARD_NAME'" >&2
    exit 1
fi

# Resolve the short hw name from /proc/asound/cards
HW_NAME=$(grep -F "$CARD_NAME" /proc/asound/cards | awk '{print $2}' | head -n1)

if [[ -z "$HW_NAME" ]]; then
    # Fall back to card number
    HW_NAME="$CARD_ID"
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
