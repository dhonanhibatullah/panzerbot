'use client';

import { useEffect, useRef } from 'react';

const KEY_MAP: Record<string, { right: number; left: number }> = {
    ArrowUp:    { right:  1, left:  1 },
    ArrowDown:  { right: -1, left: -1 },
    ArrowLeft:  { right: -1, left:  1 },
    ArrowRight: { right:  1, left: -1 },
    w:          { right:  1, left:  1 },
    s:          { right: -1, left: -1 },
    a:          { right: -1, left:  1 },
    d:          { right:  1, left: -1 },
};

export function useControlInput(
    sendMotor: (right: number, left: number) => void,
    speed: number = 0.7,
): void {
    const sendMotorRef = useRef(sendMotor);
    sendMotorRef.current = sendMotor;

    const speedRef = useRef(speed);
    speedRef.current = speed;

    useEffect(() => {
        const pressed = new Set<string>();

        const onKeyDown = (e: KeyboardEvent) => {
            if (pressed.has(e.key)) return;
            const dir = KEY_MAP[e.key];
            if (!dir) return;
            pressed.add(e.key);
            sendMotorRef.current(dir.right * speedRef.current, dir.left * speedRef.current);
        };

        const onKeyUp = (e: KeyboardEvent) => {
            if (!KEY_MAP[e.key]) return;
            pressed.delete(e.key);
            if (pressed.size === 0) sendMotorRef.current(0, 0);
        };

        window.addEventListener('keydown', onKeyDown);
        window.addEventListener('keyup', onKeyUp);
        return () => {
            window.removeEventListener('keydown', onKeyDown);
            window.removeEventListener('keyup', onKeyUp);
        };
    }, []);
}
