import {Elysia} from "elysia";

console.log('Launching Web server at port 8844')
const s = new Elysia()
    .get('/_styx/init', () => JSON.stringify({
        subscribedEvents: ['test:test']
    }))
    .get('/_styx/status', () => JSON.stringify({
        status: 'ready'
    }))
    .listen(8844)
