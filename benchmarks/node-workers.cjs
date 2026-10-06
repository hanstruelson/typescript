// Uses the same erased compute() body as workers.ts, with real Node worker threads.
const { Worker, isMainThread, parentPort, workerData } = require('node:worker_threads');
const { compute } = require('./worker-body.cjs');

if (isMainThread) {
    Promise.all(Array.from({ length: 8 }, (_, seed) => new Promise((resolve, reject) => {
        const worker = new Worker(__filename, { workerData: seed });
        worker.on('message', resolve);
        worker.on('error', reject);
        worker.on('exit', code => { if (code !== 0) reject(new Error(`Worker exited: ${code}`)); });
    }))).then(values => console.log(values.reduce((a, b) => a + b, 0))).catch(error => {
        console.error(error);
        process.exitCode = 1;
    });
} else {
    parentPort.postMessage(compute(workerData));
}
