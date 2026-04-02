'use client';

import { useSoundboard } from '@/hooks/useSoundboard';
import { ControlClient } from '@/lib/panzerbot/control';

type SoundboardProps = {
    client: ControlClient | null;
    className?: string;
};

export function Soundboard({ client, className }: SoundboardProps) {
    const { tracks, loading, error, play, stop } = useSoundboard(client);

    if (loading) return <p className={className}>Loading soundboard...</p>;
    if (error) return <p className={className}>Soundboard error: {error}</p>;

    return (
        <div className={className}>
            <div className="flex flex-wrap gap-2">
                {tracks.map((track) => (
                    <button
                        key={track.index}
                        onClick={() => play(track.index)}
                        className="btn-sound"
                    >
                        {track.name}
                    </button>
                ))}
                {tracks.length > 0 && (
                    <button onClick={stop} className="btn-sound-stop">
                        Stop
                    </button>
                )}
            </div>
        </div>
    );
}
