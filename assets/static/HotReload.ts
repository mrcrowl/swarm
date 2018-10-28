
import {SocketClient} from "./SocketClient.js";

const sc = new SocketClient();
sc.on(e => e.type == "reload" && window.location.reload());
sc.connect();