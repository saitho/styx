import {Elysia} from "elysia";
import * as grpc from "@grpc/grpc-js";
import * as protoLoader from "@grpc/proto-loader";

const PROTO_PATH = __dirname + '/../../proto/service-api.proto';
// Suggested options for similarity to existing grpc.load behavior
const packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
    });
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);

const styxapp = protoDescriptor.styxapp as any;

function getServer() {
    var server = new grpc.Server();
    server.addService(styxapp.EventService.service, {
        Emit: (call: any) => {
            console.log('Emit', call)
            call.end();
        },
        Subscribe: (call: any) => {
            console.log('Subscribe', call)
            call.end();
        },
        Unsubscribe: (call: any) => {
            console.log('Unsubscribe', call)
            call.end();
        },
    });
    server.addService(styxapp.ServiceProvider.service, {
        Ping: (call: any) => {
            console.log('Ping', call)
            call.end();
        },
    });
    return server;
}
var routeServer = getServer();
routeServer.bindAsync('0.0.0.0:50052', grpc.ServerCredentials.createInsecure(), (error, port) => {
    if (error) {
        console.error(error)
        return
    }
    console.log('RPC Server running at port 50052')
});

console.log('Launching Web server at port 8844')
const s = new Elysia()
    .get('/_styx/init', () => JSON.stringify({
        subscribedEvents: ['test:test']
    }))
    .get('/_styx/status', () => JSON.stringify({
        status: 'ready'
    }))
    .listen(8844)
