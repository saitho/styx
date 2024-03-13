import {Elysia} from "elysia";
import styx, {eventEmitter} from "./styx.ts";

eventEmitter.addListener('test:test', (data: any) => {
    console.log('Test event!')
    console.log(data)
})

console.log('Launching Web server at port 8844')
const s = new Elysia()
    .use(styx({
        version: "2",
        subscribedEvents: ['test:test']
    }))
    .listen(8844)
