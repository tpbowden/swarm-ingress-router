var Docker = require("dockerode");
var http = require("http");
var httpProxy = require('http-proxy');

var proxy = httpProxy.createProxyServer({});


var docker = new Docker();

var hosts = {};

setInterval(function() {
  console.log("syncing services");
  docker.listServices(function(err, services) {
    if(err) {
      console.log(err);
    } else {
      services.filter(function(service) {
        return service.Spec.Labels && service.Spec.Labels.ingress == "true";
      }).forEach(function(service) {
        service.Endpoint.Ports.forEach(function(endpoint) {
          if(service.Spec.Labels && service.Spec.Labels.dnsname && (service.Spec.Labels.serviceport == endpoint.TargetPort)) {
            hosts[service.Spec.Labels.dnsname] =  {
              "DNSName": service.Spec.Labels.dnsname,
              "ServiceName": service.Spec.Name,
              "TargetPort": endpoint.TargetPort,
            };
          }
        });
      });
    }
  });
}, 10000);


http.createServer(function(req, res) {
  hostname = req.headers.host.split(":")[0];
  console.log(hosts);
  if(hosts[hostname]) {
    console.log("Proxying to http://" + hosts[hostname].ServiceName + ":" + hosts[hostname].TargetPort);
    proxy.web(req, res, {target: "http://" + hosts[hostname].ServiceName + ":" + hosts[hostname].TargetPort});
  } else {
    res.writeHead(404);
    res.end("Service not found");
  }

}).listen(8080);
