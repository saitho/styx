import {Elysia} from "elysia";
import styx, {eventEmitter, STYX_PORT} from "./styx.ts";

eventEmitter.addListener('test:test', (data: any) => {
    console.log('Test event!')
    console.log(data)
})

console.log('Launching Web server at port ' + STYX_PORT)
new Elysia()
    .use(styx({
        version: "2",
        subscribedEvents: ['test:test']
    }))
    .listen(STYX_PORT)
