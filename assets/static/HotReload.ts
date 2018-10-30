
import { SocketClient, SocketPayload } from "./SocketClient.js";

interface ReloadCSSPayloadData {
    id: string;
    css: string;
}

function reloadCSS(e: SocketPayload) {
    const { id, css } = <ReloadCSSPayloadData>JSON.parse(e.data);
    let style: HTMLStyleElement = document.querySelector("#" + CSS.escape(id));
    if (!style) {
        // new style element
        style = document.createElement('style');
        style.id = id;
        style.type = 'text/css';
        document.getElementsByTagName('head')[0].appendChild(style);
        style.appendChild(document.createTextNode(css));
    }
    else {
        // replace css within existing style element
        style.childNodes[0].textContent = css;
    }
}

const sc = new SocketClient();
sc.on(e => {
    e.type == "reload-css" && reloadCSS(e);
    e.type == "reload" && window.location.reload();
});
sc.connect();