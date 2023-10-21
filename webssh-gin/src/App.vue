<template>
  <div id="app">
    <div id="xterm"/>
  </div>
</template>

<script>
import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'

export default {
  name: 'WebShell',
  data() {
    return {
      socketURI: 'ws://127.0.0.1:10000/pod/exec/default/nginx/nginx?action=sh'
    }
  },
  mounted() {
    this.initSocket()
  },
  beforeUnmount() {
    this.socket.close()
    this.term && this.term.dispose()
  },
  methods: {
    initTerm() {
      let element = document.getElementById('xterm');
      const term = new Terminal({
        cursorBlink: true,
        fontSize: 14,
        fontFamily: 'monospace',
        fontWeight: 'normal',
        fontWeightBold: 'bold',
        lineHeight: 1,
        letterSpacing: 0,
        cursorStyle: 'underline',  // option 'bar' | 'block' | 'underline' | 'beam'
        scrollback: 1000,
        disableStdin: false,
        convertEol: true,
        screenKeys: false,
        theme: {
          foreground: '#fff',
          background: '#000',
        }
      });
      this.term = term;
      const fitAddon = new FitAddon();
      this.term.loadAddon(fitAddon);
      this.fitAddon = fitAddon;
      term.open(element);
      fitAddon.fit();
      term.focus();
      this.term.onData(data => {
        var msg = {type: "input", input: data}
        this.socket.send(JSON.stringify(msg));
      });
      window.addEventListener('resize',this.resizeTerm);
    },
    getColsAndRows(element) {
      element = element || document.getElementById('xterm');
      return {
        rows: parseInt((element.clientHeight - 0 / 18)),
        cols: 10
      };
    },
    resizeTerm() {
      this.fitAddon.fit();
      this.term.scrollToBottom();
    },
    initSocket() {
      this.socket = new WebSocket(`${this.socketURI}`);
      this.socketOnClose();
      this.socketOnOpen();
      this.socketOnError();
      this.socketOnMessage();
    },
    socketOnOpen() {
      this.socket.onopen = () => {
        this.initTerm()
      }
    },
    socketOnMessage() {
      this.socket.onmessage = (event) => {
        this.term.write(event.data.toString());
    }
  },
  socketOnClose() {
    this.socket.onclose = () => {
      console.log('socket closed');
  }
},
  socketOnError() {
    this.socket.onerror = () => {
      console.log('socket error');
  }
  }
}
}
</script>

<style>
  #xterm {
    padding: 15px 0;
  }
</style>