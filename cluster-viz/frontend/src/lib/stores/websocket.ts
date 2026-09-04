import { writable } from 'svelte/store';
import type { GraphEvent } from '../types/graph';
import { applyEvent } from './graph';

export type ConnectionStatus = 'connecting' | 'connected' | 'disconnected' | 'error';

export const connectionStatus = writable<ConnectionStatus>('disconnected');

let socket: WebSocket | null = null;
let reconnectTimeout: number | null = null;
let reconnectDelay = 1000; // Start at 1 second
const maxReconnectDelay = 30000; // Cap at 30 seconds

export function connectWebSocket() {
  if (socket && (socket.readyState === WebSocket.CONNECTING || socket.readyState === WebSocket.OPEN)) {
    return;
  }

  connectionStatus.set('connecting');

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const wsUrl = `${protocol}//${window.location.host}/api/ws`;

  socket = new WebSocket(wsUrl);

  socket.onopen = () => {
    console.log('WebSocket connected');
    connectionStatus.set('connected');
    reconnectDelay = 1000; // Reset delay on successful connection
  };

  socket.onmessage = (event) => {
    try {
      const graphEvent: GraphEvent = JSON.parse(event.data);
      applyEvent(graphEvent);
    } catch (err) {
      console.error('Failed to parse WebSocket message:', err);
    }
  };

  socket.onerror = (error) => {
    console.error('WebSocket error:', error);
    connectionStatus.set('error');
  };

  socket.onclose = () => {
    console.log('WebSocket closed');
    connectionStatus.set('disconnected');
    socket = null;

    // Reconnect with exponential backoff
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
    }

    reconnectTimeout = window.setTimeout(() => {
      console.log(`Reconnecting in ${reconnectDelay}ms...`);
      connectWebSocket();
      reconnectDelay = Math.min(reconnectDelay * 2, maxReconnectDelay);
    }, reconnectDelay);
  };
}

export function disconnectWebSocket() {
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout);
    reconnectTimeout = null;
  }

  if (socket) {
    socket.close();
    socket = null;
  }

  connectionStatus.set('disconnected');
}
