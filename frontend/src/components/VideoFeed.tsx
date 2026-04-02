'use client';

import { useEffect, useRef } from 'react';

type VideoFeedProps = {
    stream: MediaStream | null;
    className?: string;
};

export function VideoFeed({ stream, className }: VideoFeedProps) {
    const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
        const el = videoRef.current;
        if (!el) return;
        el.srcObject = stream;
        if (stream) {
            el.play().catch(() => {});
        }
    }, [stream]);

    return (
        <video
            ref={videoRef}
            autoPlay
            playsInline
            muted
            className={className}
        />
    );
}
