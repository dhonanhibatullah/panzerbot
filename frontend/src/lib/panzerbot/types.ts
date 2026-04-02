// ── Control WebSocket ────────────────────────────────────────────────────────

export enum ControlCode {
    Motor = 'motor',
    Servo = 'servo',
}

export type ControlDataMotor = {
    right: number;
    left: number;
};

export type ControlDataServo = {
    pan: number;
    tilt: number;
};

export type ControlFrame =
    | { code: ControlCode.Motor; data: ControlDataMotor }
    | { code: ControlCode.Servo; data: ControlDataServo };

// ── Soundboard ───────────────────────────────────────────────────────────────

export type SoundboardTrack = {
    index: number;
    name: string;
};

// ── RTC Signalling WebSocket ─────────────────────────────────────────────────

export type RtcInbound =
    | { type: 'offer'; sdp: string }
    | { type: 'ice-candidate'; candidate: string; sdpMid: string; sdpMLineIndex: number }
    | { type: 'error' };

export type RtcOutbound =
    | { type: 'answer'; sdp: string }
    | { type: 'ice-candidate'; candidate: string; sdpMid: string; sdpMLineIndex: number }
    | { type: 'close' };
