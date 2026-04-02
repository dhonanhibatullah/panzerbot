'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { RTCClient, RTCClientState } from '@/lib/panzerbot/rtc';

export type UseWebRTCReturn = {
    state: RTCClientState;
    remoteStream: MediaStream | null;
    start: () => Promise<void>;
    stop: () => void;
};

export function useWebRTC(): UseWebRTCReturn {
    const clientRef = useRef<RTCClient | null>(null);
    const [state, setState] = useState<RTCClientState>('idle');
    const [remoteStream, setRemoteStream] = useState<MediaStream | null>(null);

    useEffect(() => {
        const client = new RTCClient({
            onStateChange: setState,
            onRemoteStream: setRemoteStream,
        });
        clientRef.current = client;

        return () => {
            client.stop();
            clientRef.current = null;
        };
    }, []);

    const start = useCallback(async () => {
        await clientRef.current?.start();
    }, []);

    const stop = useCallback(() => {
        clientRef.current?.stop();
        setRemoteStream(null);
    }, []);

    return { state, remoteStream, start, stop };
}
