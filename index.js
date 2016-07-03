var Docker = require("dockerode");
var http = require("http");
var httpProxy = require('http-proxy');
var Promise = require('bluebird');
var logger = require("./logger.js");

var proxy = httpProxy.createProxyServer({});
var hosts = {};

var docker = new Docker();

function getAllServices(docker) {
  return new Promise(function(resolve, reject) { 
    docker.listServices(function(err, services) {
      if(err) {
        reject(err);
      } else {
        if(services) {
          resolve(services);
        } else {
          reject("No services");
        }
      }
    });
  });
}

function filterIngressServices(services) {
  return services.filter(function(service) {
    return service.Spec.Labels && service.Spec.Labels.ingress == "true";
  });
}

function registerEndpoint(service) {
  service.Endpoint.Ports.forEach(function(endpoint) {
    if(service.Spec.Labels && service.Spec.Labels.dnsname && (service.Spec.Labels.serviceport == endpoint.TargetPort)) {
      logger.info("Registering service " + service.Spec.Name);
      hosts[service.Spec.Labels.dnsname] =  {
        "DNSName": service.Spec.Labels.dnsname,
        "ServiceName": service.Spec.Name,
        "TargetPort": endpoint.TargetPort,
      };
    }
  });
}


function loadServices() {
  logger.info("Starting service sync");
  getAllServices(docker).then(function(services) {
    filterIngressServices(services).forEach(function(service) {
      registerEndpoint(service);
    });
  }).catch(function(err) {
    logger.error("Failed to retrieve services: " + err);
  });
}

loadServices();

setInterval(loadServices, 10000);


http.createServer(function(req, res) {
  logger.info("Processing request");
  hostname = req.headers.host.split(":")[0];
  logger.info("Received request for " + hostname);
  if(hosts[hostname]) {
    logger.info("Proxying to http://" + hosts[hostname].ServiceName + ":" + hosts[hostname].TargetPort);
    proxy.web(req, res, {target: "http://" + hosts[hostname].ServiceName + ":" + hosts[hostname].TargetPort}, function(e) {
      logger.error("Failed to proxy to " + hosts[hostname].ServiceName + ": " + e);
      res.writeHead(503);
      res.end("Service not available");
    });
  } else {
    logger.warn("Failed to find service for " + hostname);
    res.writeHead(404);
    res.end("Service not found");
  }

}).listen(8080);

http.createServer(function(req, res) {
  res.writeHead(200);
  res.end(JSON.stringify(hosts));
}).listen(9090);

logger.info("Router listening on port 8080");
logger.info("Metrics listening on port 9090");
