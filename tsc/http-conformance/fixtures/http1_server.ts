// @ts-nocheck
const http = require("node:http");

const server = http.createServer((request, response) => {
	let requestBytes = 0;
	request.on("data", (chunk) => {
		requestBytes += chunk.length;
	});
	request.on("end", () => {
		response.statusCode = 200;
		response.setHeader("content-type", "text/plain; charset=utf-8");
		response.end(request.method + " " + request.url + " version=" + request.httpVersion + " bytes=" + requestBytes);
	});
});

const port = Number(process.env.HTTP1_PORT || "8080");
server.listen(port, "127.0.0.1", () => {
	console.log("HTTP1_READY " + port);
});
