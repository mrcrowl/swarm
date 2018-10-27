export interface Listener<T> {
	(event: T): any;
}

export interface Disposable {
	dispose();
}

/** A type safe event emitter */
export class EventEmitter<T> {
	private listeners: Listener<T>[] = [];
	private listenersOncer: Listener<T>[] = [];

	on = (listener: Listener<T>): Disposable => {
		this.listeners.push(listener);
		return {
			dispose: () => this.off(listener),
		};
	};

	once = (listener: Listener<T>): void => {
		this.listenersOncer.push(listener);
	};

	off = (listener: Listener<T>) => {
		var callbackIndex = this.listeners.indexOf(listener);
		if (callbackIndex > -1) this.listeners.splice(callbackIndex, 1);
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
	authSent: boolean;
	emitter: EventEmitter<SocketPayload>;
	client: WebSocket;

	constructor(opts) {
		opts = opts || {};
		const port = opts.port || window.location.port;
		const protocol = location.protocol === "https:" ? "wss://" : "ws://";
		const domain = location.hostname || "localhost";
		this.url = opts.host || `${protocol}${domain}:${port}/ws`;
		if (opts.uri) {
			this.url = opts.uri;
		}
		this.authSent = false;
		this.emitter = new EventEmitter();
	}
	reconnect(fn) {
		setTimeout(() => {
			// this.emitter.emit("reconnect", { message: "Trying to reconnect" });
			this.connect(fn);
		}, 5000);
	}

	on(fn) {
		this.emitter.on(fn);
	}

	connect(fn?: OnOpenFn) {
		console.log("%cConnecting to websocket at " + this.url, "color: #237abe");

		setTimeout(() => {
			this.client = new WebSocket(this.url);
			this.bindEvents(fn);
		}, 0);
    }
    
	close() {
		this.client.close();
    }
    
	send(eventName, data) {
		if (this.client.readyState === 1) {
			this.client.send(JSON.stringify({ event: eventName, data: data || {} }));
		}
	}

	private error(data: { reason: any; message: string }) {
		// this.emitter.emit("error", data);
	}

	/** Wires up the socket client messages to be emitted on our event emitter */
	private bindEvents(fn?: OnOpenFn) {
		this.client.onopen = event => {
			console.log("%cConnected", "color: #237abe");
			if (fn) {
				fn(this);
			}
		};
		this.client.onerror = (event: any) => {
			this.error({ reason: event.reason, message: "Socket error" });
        };
        
		this.client.onclose = event => {
			// this.emitter.emit("close", { message: "Socket closed" });
			if (event.code !== 1011) {
				this.reconnect(fn);
			}
		};
		this.client.onmessage = event => {
			let data = event.data;
			if (data) {
                let payload = <SocketPayload>JSON.parse(data);
                this.emitter.emit(payload)
				// this.emitter.emit(item.type, item.data);
			}
		};
	}
}