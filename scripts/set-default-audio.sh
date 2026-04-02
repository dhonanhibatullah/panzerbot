#!/usr/bin/env bash
set -euo pipefail

# This should match the name inside the brackets [ ] in aplay -l
CARD_NAME="USB Audio Device"
ASOUNDRC="$HOME/.asoundrc"

echo "Searching for: $CARD_NAME..."

# 1. Improved extraction logic
# We use -F to treat brackets as literal strings
# We use awk -F'[: ]+' to treat both colons and spaces as delimiters
# Field 2 is usually the Card Number (e.g., '1')
HW_ID=$(aplay -l | grep -F "[$CARD_NAME]" | head -n1 | awk -F'[: ]+' '{print $2}')

if [[ -z "$HW_ID" ]]; then
    echo "Error: Could not find ALSA card matching '$CARD_NAME'" >&2
    echo "Available cards are:"
    aplay -l | grep card
    exit 1
fi

echo "Found $CARD_NAME at card index: $HW_ID"

# 2. Generate .asoundrc
# Using 'plughw' is generally safer for USB cards as it handles sample rate conversions
cat > "$ASOUNDRC" << EOF
pcm.!default {
    type plug
    slave {
        pcm "hw:$HW_ID,0"
    }
}

ctl.!default {
    type hw
    card $HW_ID
}
EOF

echo "Successfully written to $ASOUNDRC"
echo "--------------------------------"
cat "$ASOUNDRC"
echo "--------------------------------"

# 3. Test Audio
echo "Testing audio... (If you hear nothing, check 'alsamixer' to ensure card $HW_ID is unmuted)"
# We explicitly tell speaker-test to use the 'default' we just defined
speaker-test -D default -c 2 -t wav -l 1