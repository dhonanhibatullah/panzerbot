'use client';

import { useEffect, useRef, useState } from 'react';
import { ControlClient, ControlClientState } from '@/lib/panzerbot/control';

export type UseControlReturn = {
    state: ControlClientState;
    sendMotor: (right: number, left: number) => void;
    sendServo: (pan: number, tilt: number) => void;
    client: ControlClient | null;
};

export function useControl(): UseControlReturn {
    const clientRef = useRef<ControlClient | null>(null);
    const [state, setState] = useState<ControlClientState>('closed');

    useEffect(() => {
        const client = new ControlClient({ onStateChange: setState });
        clientRef.current = client;
        client.connect();
        return () => client.destroy();
    }, []);

    const sendMotor = (right: number, left: number) => {
        clientRef.current?.sendMotor(right, left);
    };

    const sendServo = (pan: number, tilt: number) => {
        clientRef.current?.sendServo(pan, tilt);
    };

    return { state, sendMotor, sendServo, client: clientRef.current };
}
