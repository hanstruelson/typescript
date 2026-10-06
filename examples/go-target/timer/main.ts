let count = 0;
setTimeout(() => {
    count++;
    console.log("timer", count);
}, 0);
Promise.resolve().then(() => console.log("promise", count));
console.log("sync", count);
