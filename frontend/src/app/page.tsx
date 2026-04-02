'use client';

import { useControlInput } from '@/hooks/useControlInput';
import { useControl } from '@/hooks/useControl';
import { useWebRTC } from '@/hooks/useWebRTC';
import { useRef } from 'react';
import { ControlButtons } from '@/components/ControlButtons';
import { Soundboard } from '@/components/Soundboard';
import { VideoFeed } from '@/components/VideoFeed';
import { RTCClientState } from '@/lib/panzerbot/rtc';
import { ControlClientState } from '@/lib/panzerbot/control';

// ── Helpers ───────────────────────────────────────────────────────────────────

function controlDot(state: ControlClientState) {
    if (state === 'open') return 'active';
    if (state === 'error') return 'error';
    if (state === 'connecting') return 'pulse';
    return 'inactive';
}

function rtcDot(state: RTCClientState) {
    if (state === 'connected') return 'active';
    if (state === 'error') return 'error';
    if (state === 'connecting') return 'pulse';
    return 'inactive';
}

function rtcLabel(state: RTCClientState) {
    const map: Record<RTCClientState, string> = {
        idle: 'IDLE',
        connecting: 'CONNECTING',
        connected: 'LIVE',
        disconnected: 'DISCONNECTED',
        error: 'ERROR',
    };
    return map[state];
}

// ── Page ──────────────────────────────────────────────────────────────────────

export default function Home() {
    const { state: ctrlState, sendMotor, sendServo, client } = useControl();
    const { state: rtcState, remoteStream, start, stop } = useWebRTC();
    const servoRef = useRef({ pan: 0, tilt: 0 });

    useControlInput(sendMotor);

    const isLive = rtcState === 'connected';
    const isConnecting = rtcState === 'connecting';

    return (
        <div className="min-h-dvh flex flex-col gap-3 p-3 md:p-4" style={{ background: 'var(--bg)' }}>

            {/* ── Header ── */}
            <header className="flex items-center justify-between panel px-4 py-2">
                <div className="flex items-center gap-3">
                    <span
                        className="text-sm font-bold tracking-[0.2em]"
                        style={{ color: 'var(--accent)', fontFamily: 'var(--font-geist-mono)' }}
                    >
                        PANZERBOT
                    </span>
                    <span className="panel-label">// remote control interface</span>
                </div>

                <div className="flex items-center gap-4">
                    <div className="flex items-center gap-2">
                        <span className={`status-dot ${controlDot(ctrlState)}`} />
                        <span className="panel-label">CTRL</span>
                    </div>
                    <div className="flex items-center gap-2">
                        <span className={`status-dot ${rtcDot(rtcState)}`} />
                        <span className="panel-label">{rtcLabel(rtcState)}</span>
                    </div>
                </div>
            </header>

            {/* ── Main area ── */}
            <div className="flex flex-col lg:flex-row gap-3 flex-1">

                {/* ── Video feed ── */}
                <div className="panel flex-1 flex flex-col overflow-hidden" style={{ minHeight: '280px' }}>
                    <div className="flex items-center justify-between px-3 pt-2 pb-1">
                        <span className="panel-label">CAM FEED</span>
                        <div className="flex items-center gap-2">
                            {isLive && (
                                <span
                                    className="panel-label"
                                    style={{ color: 'var(--green)', animation: 'pulse 1.5s ease-in-out infinite' }}
                                >
                                    ● REC
                                </span>
                            )}
                            <button
                                onClick={isLive || isConnecting ? stop : start}
                                className={`btn-call ${isLive || isConnecting ? 'stop' : 'start'}`}
                                disabled={isConnecting}
                            >
                                {isConnecting ? 'CONNECTING...' : isLive ? 'END CALL' : 'START CALL'}
                            </button>
                        </div>
                    </div>

                    <div
                        className="flex-1 flex items-center justify-center relative"
                        style={{ background: '#030507' }}
                    >
                        {isLive ? (
                            <VideoFeed stream={remoteStream} className="w-full h-full object-cover" />
                        ) : (
                            <div className="flex flex-col items-center gap-2" style={{ color: 'var(--text-dim)' }}>
                                <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.2">
                                    <path d="M15 10l4.553-2.276A1 1 0 0121 8.723v6.554a1 1 0 01-1.447.894L15 14M3 8a2 2 0 012-2h8a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2V8z" />
                                </svg>
                                <span className="panel-label">no signal — press start call</span>
                            </div>
                        )}

                        {/* corner brackets */}
                        {['top-2 left-2', 'top-2 right-2', 'bottom-2 left-2', 'bottom-2 right-2'].map((pos, i) => (
                            <span
                                key={i}
                                className={`absolute ${pos} w-4 h-4 pointer-events-none`}
                                style={{
                                    borderTop: i < 2 ? `1px solid var(--accent-dim)` : 'none',
                                    borderBottom: i >= 2 ? `1px solid var(--accent-dim)` : 'none',
                                    borderLeft: i % 2 === 0 ? `1px solid var(--accent-dim)` : 'none',
                                    borderRight: i % 2 === 1 ? `1px solid var(--accent-dim)` : 'none',
                                }}
                            />
                        ))}
                    </div>
                </div>

                {/* ── Right sidebar ── */}
                <div className="flex flex-col gap-3 lg:w-64">

                    {/* Drive controls */}
                    <div className="panel p-3 flex flex-col gap-3">
                        <span className="panel-label">DRIVE</span>

                        <div className="flex justify-center">
                            <ControlButtons sendMotor={sendMotor} />
                        </div>

                        {/* Keyboard hint */}
                        <div className="flex flex-col gap-1.5" style={{ color: 'var(--text-dim)', fontSize: '0.65rem' }}>
                            <div className="flex items-center gap-1.5">
                                <span className="kbd">W</span><span className="kbd">A</span><span className="kbd">S</span><span className="kbd">D</span>
                                <span style={{ marginLeft: 4 }}>or</span>
                                <span className="kbd">↑</span><span className="kbd">←</span><span className="kbd">↓</span><span className="kbd">→</span>
                            </div>
                            <span>Hold key / button to move. Release to stop.</span>
                        </div>
                    </div>

                    {/* Servo controls */}
                    <div className="panel p-3 flex flex-col gap-2">
                        <span className="panel-label">SERVO / CAMERA</span>
                        <div className="flex flex-col gap-2">
                            <div className="flex flex-col gap-1">
                                <div className="flex justify-between" style={{ fontSize: '0.65rem', color: 'var(--text-dim)' }}>
                                    <span>PAN</span>
                                    <span>← left / right →</span>
                                </div>
                                <input
                                    type="range"
                                    min={-1.5}
                                    max={1.5}
                                    step={0.01}
                                    defaultValue={0}
                                    className="w-full accent-(--accent)"
                                    onChange={(e) => {
                                        servoRef.current.pan = parseFloat(e.target.value);
                                        sendServo(servoRef.current.pan, servoRef.current.tilt);
                                    }}
                                />
                            </div>
                            <div className="flex flex-col gap-1">
                                <div className="flex justify-between" style={{ fontSize: '0.65rem', color: 'var(--text-dim)' }}>
                                    <span>TILT</span>
                                    <span>↑ up / down ↓</span>
                                </div>
                                <input
                                    type="range"
                                    min={-1.5}
                                    max={1.5}
                                    step={0.01}
                                    defaultValue={0}
                                    className="w-full accent-(--accent)"
                                    onChange={(e) => {
                                        servoRef.current.tilt = parseFloat(e.target.value);
                                        sendServo(servoRef.current.pan, servoRef.current.tilt);
                                    }}
                                />
                            </div>
                        </div>
                    </div>

                    {/* Soundboard */}
                    <div className="panel p-3 flex flex-col gap-2 flex-1">
                        <span className="panel-label">SOUNDBOARD</span>
                        <Soundboard client={client} />
                    </div>

                </div>
            </div>

            {/* ── Footer ── */}
            <footer
                className="flex items-center justify-between px-2"
                style={{ fontSize: '0.6rem', color: 'var(--text-dim)', fontFamily: 'var(--font-geist-mono)' }}
            >
                <span>CTRL: {ctrlState.toUpperCase()}</span>
                <span>RTC: {rtcState.toUpperCase()}</span>
            </footer>
        </div>
    );
}
