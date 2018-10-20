if (!WebAssembly.instantiateStreaming) { // polyfill
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetch("./wasm/main.wasm"), go.importObject).then((result) => {
    mod = result.module;
    inst = result.instance;
    document.getElementById("runButton").disabled = false;

    // Start wasm immediately.
    go.run(inst);
});

async function run() {
    console.clear();

    // This is standard sample code to run wasm from button.
    // It assumes the wasm runs then exits.
    // await go.run(inst);
    // inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
}
