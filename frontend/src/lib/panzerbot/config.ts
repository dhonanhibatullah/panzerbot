const backendOrigin = process.env.NEXT_PUBLIC_BACKEND_ORIGIN ?? 'localhost:6767';
const backendProto = process.env.NEXT_PUBLIC_BACKEND_PROTO ?? 'http';
const wsProto = backendProto === 'https' ? 'wss' : 'ws';

export const config = {
    httpBase: `${backendProto}://${backendOrigin}`,
    wsBase: `${wsProto}://${backendOrigin}`,
} as const;
