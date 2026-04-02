'use client';

import { useCallback, useEffect, useState } from 'react';
import { ControlClient } from '@/lib/panzerbot/control';
import { SoundboardTrack } from '@/lib/panzerbot/types';

export type UseSoundboardReturn = {
    tracks: SoundboardTrack[];
    loading: boolean;
    error: string | null;
    play: (index: number) => Promise<void>;
    stop: () => Promise<void>;
};

export function useSoundboard(client: ControlClient | null): UseSoundboardReturn {
    const [tracks, setTracks] = useState<SoundboardTrack[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!client) return;
        setLoading(true);
        client
            .getSoundboardTracks()
            .then(setTracks)
            .catch((e: unknown) => setError(e instanceof Error ? e.message : String(e)))
            .finally(() => setLoading(false));
    }, [client]);

    const play = useCallback(
        async (index: number) => {
            try {
                await client?.playTrack(index);
            } catch (e: unknown) {
                setError(e instanceof Error ? e.message : String(e));
            }
        },
        [client],
    );

    const stop = useCallback(async () => {
        try {
            await client?.stopTracks();
        } catch (e: unknown) {
            setError(e instanceof Error ? e.message : String(e));
        }
    }, [client]);

    return { tracks, loading, error, play, stop };
}
