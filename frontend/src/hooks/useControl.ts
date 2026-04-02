'use client';

import { useEffect, useState } from 'react';
import { ControlClient, ControlClientState } from '@/lib/panzerbot/control';

export type UseControlReturn = {
    state: ControlClientState;
    sendMotor: (right: number, left: number) => void;
    sendServo: (pan: number, tilt: number) => void;
    client: ControlClient | null;
};

export function useControl(): UseControlReturn {
    const [client, setClient] = useState<ControlClient | null>(null);
    const [state, setState] = useState<ControlClientState>('closed');

    useEffect(() => {
        const c = new ControlClient({ onStateChange: setState });
        setClient(c);
        c.connect();
        return () => c.destroy();
    }, []);

    const sendMotor = (right: number, left: number) => {
        client?.sendMotor(right, left);
    };

    const sendServo = (pan: number, tilt: number) => {
        client?.sendServo(pan, tilt);
    };

    return { state, sendMotor, sendServo, client };
}
