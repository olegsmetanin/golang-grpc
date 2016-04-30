var PROTO_PATH = __dirname + '/../api/api.proto';

var grpc = require('grpc');
var api_proto = grpc.load(PROTO_PATH).api;

function main() {

  var client = new api_proto.Greeter('localhost:10000',
                                       grpc.credentials.createInsecure());

  var user;
  if (process.argv.length >= 3) {
    user = process.argv[2];
  } else {
    user = 'world';
  }
  client.sayHello({name: user}, function(err, response) {
    if (err) {
      console.log('Error:', err);
    } else {
      console.log('Greeting:', response.message);
    }
  });
}

main();