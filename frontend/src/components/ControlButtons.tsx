'use client';

import { useControlButton } from '@/hooks/useControlButton';

type ControlButtonsProps = {
    sendMotor: (right: number, left: number) => void;
    speed?: number;
    className?: string;
};

export function ControlButtons({ sendMotor, speed, className }: ControlButtonsProps) {
    const buttons = useControlButton(sendMotor, speed);

    return (
        <div className={className}>
            <div className="grid grid-cols-3 gap-2 w-fit">
                <div />
                <button
                    className="btn-control"
                    {...buttons.forward}
                >
                    ▲
                </button>
                <div />

                <button
                    className="btn-control"
                    {...buttons.left}
                >
                    ◀
                </button>
                <button
                    className="btn-control"
                    {...buttons.stop}
                >
                    ■
                </button>
                <button
                    className="btn-control"
                    {...buttons.right}
                >
                    ▶
                </button>

                <div />
                <button
                    className="btn-control"
                    {...buttons.backward}
                >
                    ▼
                </button>
                <div />
            </div>
        </div>
    );
}
