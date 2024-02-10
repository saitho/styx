import {Elysia} from "elysia";

const s = new Elysia()
    .get('/_styx/init', () => JSON.stringify({
        subscribedEvents: ['test:test']
    }))
    .get('/_styx/status', () => JSON.stringify({
        status: 'ready'
    }))
    .listen(8844)
