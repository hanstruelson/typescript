// @ts-nocheck
const http2 = require("node:http2");
const server = http2.createServer();

server.on("stream", (stream, headers) => {
	let requestBytes = 0;
	stream.on("error", (error) => {
		console.error("HTTP2_STREAM_ERROR " + (error.code || error.message));
	});
	stream.on("data", (chunk) => {
		requestBytes += chunk.length;
	});
	stream.on("end", () => {
		if (stream.destroyed || stream.closed) return;
		stream.respond({ ":status": 200, "content-type": "text/plain; charset=utf-8" });
		stream.end("bytes=" + requestBytes + " path=" + headers[":path"]);
	});
});

const port = Number(process.env.HTTP2_PORT || "8081");
server.listen(port, "127.0.0.1", () => {
	console.log("HTTP2_READY " + port);
});
