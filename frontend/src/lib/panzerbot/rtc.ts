import { config } from './config';
import { RtcInbound, RtcOutbound } from './types';

export type RTCClientState = 'idle' | 'connecting' | 'connected' | 'disconnected' | 'error';

export type RTCClientOptions = {
    onStateChange?: (state: RTCClientState) => void;
    onRemoteStream?: (stream: MediaStream) => void;
};

export class RTCClient {
    private ws: WebSocket | null = null;
    private pc: RTCPeerConnection | null = null;
    private localStream: MediaStream | null = null;
    private state: RTCClientState = 'idle';
    private readonly options: RTCClientOptions;
    private pendingCandidates: RTCIceCandidateInit[] = [];

    constructor(options: RTCClientOptions = {}) {
        this.options = options;
    }

    async start(): Promise<void> {
        if (this.state !== 'idle' && this.state !== 'disconnected' && this.state !== 'error') return;

        this.setState('connecting');

        try {
            this.localStream = await navigator.mediaDevices.getUserMedia({ audio: true, video: false });
        } catch {
            this.setState('error');
            return;
        }

        this.pc = new RTCPeerConnection({
            iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
        });

        for (const track of this.localStream.getTracks()) {
            this.pc.addTrack(track, this.localStream);
        }

        this.pc.ontrack = (event) => {
            if (event.streams[0]) {
                this.options.onRemoteStream?.(event.streams[0]);
            }
        };

        this.pc.onicecandidate = (event) => {
            if (!event.candidate) return;
            this.wsSend({
                type: 'ice-candidate',
                candidate: event.candidate.candidate,
                sdpMid: event.candidate.sdpMid ?? '',
                sdpMLineIndex: event.candidate.sdpMLineIndex ?? 0,
            });
        };

        this.pc.onconnectionstatechange = () => {
            switch (this.pc?.connectionState) {
                case 'connected':
                    this.setState('connected');
                    break;
                case 'disconnected':
                case 'closed':
                    this.setState('disconnected');
                    break;
                case 'failed':
                    this.setState('error');
                    break;
            }
        };

        this.ws = new WebSocket(`${config.wsBase}/v1/rtc/ws`);

        this.ws.onmessage = (event) => this.handleSignalMessage(event);

        this.ws.onclose = () => {
            if (this.state !== 'error') this.setState('disconnected');
        };

        this.ws.onerror = () => {
            this.setState('error');
        };
    }

    stop(): void {
        this.wsSend({ type: 'close' });
        this.teardown();
    }

    getState(): RTCClientState {
        return this.state;
    }

    private async handleSignalMessage(event: MessageEvent): Promise<void> {
        let msg: RtcInbound;
        try {
            msg = JSON.parse(event.data as string) as RtcInbound;
        } catch {
            return;
        }

        if (msg.type === 'error') {
            this.setState('error');
            this.teardown();
            return;
        }

        if (msg.type === 'offer') {
            if (!this.pc) return;
            await this.pc.setRemoteDescription({ type: 'offer', sdp: msg.sdp });

            for (const c of this.pendingCandidates) {
                await this.pc.addIceCandidate(c).catch(() => {});
            }
            this.pendingCandidates = [];

            const answer = await this.pc.createAnswer();
            await this.pc.setLocalDescription(answer);
            this.wsSend({ type: 'answer', sdp: answer.sdp! });
            return;
        }

        if (msg.type === 'ice-candidate') {
            const init: RTCIceCandidateInit = {
                candidate: msg.candidate,
                sdpMid: msg.sdpMid,
                sdpMLineIndex: msg.sdpMLineIndex,
            };
            if (this.pc?.remoteDescription) {
                await this.pc.addIceCandidate(init).catch(() => {});
            } else {
                this.pendingCandidates.push(init);
            }
        }
    }

    private wsSend(msg: RtcOutbound): void {
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(msg));
        }
    }

    private teardown(): void {
        this.ws?.close();
        this.ws = null;
        this.pc?.close();
        this.pc = null;
        this.localStream?.getTracks().forEach((t) => t.stop());
        this.localStream = null;
        this.pendingCandidates = [];
    }

    private setState(next: RTCClientState): void {
        this.state = next;
        this.options.onStateChange?.(next);
    }
}
