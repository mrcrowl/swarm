export interface Listener<T> {
	(event: T): any;
}

export class EventEmitter<T> {
	private listeners: Listener<T>[] = [];
	private listenersOncer: Listener<T>[] = [];

	on = (listener: Listener<T>) => {
		this.listeners.push(listener);
	};

	emit = (event: T) => {
		/** Update any general listeners */
		this.listeners.forEach(listener => listener(event));

		/** Clear the `once` queue */
		this.listenersOncer.forEach(listener => listener(event));
		this.listenersOncer = [];
	};
}

export interface SocketPayload {
    type: string
    data: string
}

export type OnOpenFn = (client: SocketClient) => void;

export class SocketClient {
	url: string;
	emitter: EventEmitter<SocketPayload>;
	client: WebSocket;

	constructor() {
		const port = window.location.port;
		const protocol = location.protocol === "https:" ? "wss://" : "ws://";
		const domain = location.hostname || "localhost";
		this.url = `${protocol}${domain}:${port}/__swarm__/ws`;
		this.emitter = new EventEmitter();
	}
	reconnect() {
		setTimeout(() => this.connect(), 5000)
	}

	on(fn: Listener<SocketPayload>) {
		this.emitter.on(fn);
	}

	connect() {
		console.log("%cConnecting to websocket at " + this.url, "color: #237abe");
		setTimeout(() => {
			this.client = new WebSocket(this.url);
			this.bindEvents();
		}, 0);
    }
    
	send(eventName, data) {
		if (this.client.readyState === 1) {
			this.client.send(JSON.stringify({ event: eventName, data: data || {} }));
		}
	}

	/** Wires up the socket client messages to be emitted on our event emitter */
	private bindEvents() {
		this.client.onopen = event => { console.log("%cConnected", "color: #237abe"); this.client.onclose = (event: CloseEvent) => this.reconnect(); };
		this.client.onerror = (event: any) => console.error(event);
		this.client.onmessage = (event: MessageEvent) => event.data && this.emitter.emit(<SocketPayload>JSON.parse(event.data));
	}
}