import {Elysia} from "elysia";
import AwaitEventEmitter from "await-event-emitter"
export const eventEmitter = new AwaitEventEmitter();
export const STYX_PORT = 8844

interface Config {
    subscribedEvents: string[];
    version: string;
    status?: () => string;
}

export default (config: Config) => {
    return new Elysia()
        .get('/_styx/init', () => JSON.stringify({
            // cached on main server for version emitted by "version" endpoint; or if reinitialize endpoint is triggered via deployment (TBD)
            subscribedEvents: config.subscribedEvents
        }))
        .post('/_styx/event', ({ body }) => {
            if (!body) {
                return new Response("Missing body", {status: 400})
            }
            const responseBody = (body as {event: string, data: any})
            if (!responseBody.event || !config.subscribedEvents.includes(responseBody.event)) {
                return new Response("Unsupported event", {status: 400})
            }
            return JSON.stringify({
                // returns false if no listeners handled the event
                success: eventEmitter.emitSync(responseBody.event, responseBody.data)
            })
        })
        .get('/_styx/status', () => {
            let status = 'ready'
            if (config.status !== undefined) {
                status = config.status()
            }
            return JSON.stringify({
                status: status
            });
        })
        .get('/_styx/version', () => JSON.stringify({
            version: config.version // main server will cache config from "init" step until new version is set
        }))
}
