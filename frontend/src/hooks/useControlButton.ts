'use client';

export function useControlButton(
    sendMotor: (right: number, left: number) => void,
    speed: number = 1.0,
) {
    const makeHandlers = (right: number, left: number) => ({
        onPointerDown: () => sendMotor(right * speed, left * speed),
        onPointerUp: () => sendMotor(0, 0),
        onPointerLeave: () => sendMotor(0, 0),
    });

    return {
        forward: makeHandlers(1, 1),
        backward: makeHandlers(-1, -1),
        left: makeHandlers(-1, 1),
        right: makeHandlers(1, -1),
        stop: {
            onPointerDown: () => sendMotor(0, 0),
        },
    };
}
