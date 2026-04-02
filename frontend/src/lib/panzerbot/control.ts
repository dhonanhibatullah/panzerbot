import { config } from './config';
import { ControlCode, ControlFrame, SoundboardTrack } from './types';

export type ControlClientState = 'connecting' | 'open' | 'closed' | 'error';

export type ControlClientOptions = {
    onStateChange?: (state: ControlClientState) => void;
};

export class ControlClient {
    private ws: WebSocket | null = null;
    private state: ControlClientState = 'closed';
    private readonly options: ControlClientOptions;

    constructor(options: ControlClientOptions = {}) {
        this.options = options;
    }

    connect(): void {
        if (this.ws) return;

        this.setState('connecting');
        this.ws = new WebSocket(`${config.wsBase}/v1/peripheral/ws`);

        this.ws.onopen = () => this.setState('open');

        this.ws.onclose = () => {
            this.ws = null;
            this.setState('closed');
        };

        this.ws.onerror = () => {
            this.ws = null;
            this.setState('error');
        };
    }

    destroy(): void {
        this.ws?.close();
        this.ws = null;
    }

    getState(): ControlClientState {
        return this.state;
    }

    sendMotor(right: number, left: number): void {
        this.send({ code: ControlCode.Motor, data: { right, left } });
    }

    sendServo(pan: number, tilt: number): void {
        this.send({ code: ControlCode.Servo, data: { pan, tilt } });
    }

    async getSoundboardTracks(): Promise<SoundboardTrack[]> {
        const res = await fetch(`${config.httpBase}/v1/peripheral/soundboard`);
        if (!res.ok) throw new Error(`Failed to fetch soundboard tracks: ${res.status}`);
        const body: { tracks: string[] } = await res.json();
        return body.tracks.map((name, index) => ({ index, name }));
    }

    async playTrack(index: number): Promise<void> {
        const res = await fetch(`${config.httpBase}/v1/peripheral/soundboard/${index}`, {
            method: 'POST',
        });
        if (!res.ok) throw new Error(`Failed to play track ${index}: ${res.status}`);
    }

    async stopTracks(): Promise<void> {
        const res = await fetch(`${config.httpBase}/v1/peripheral/soundboard/stop`, {
            method: 'POST',
        });
        if (!res.ok) throw new Error(`Failed to stop tracks: ${res.status}`);
    }

    private send(frame: ControlFrame): void {
        if (this.ws?.readyState !== WebSocket.OPEN) return;
        this.ws.send(JSON.stringify(frame));
    }

    private setState(next: ControlClientState): void {
        this.state = next;
        this.options.onStateChange?.(next);
    }
}
