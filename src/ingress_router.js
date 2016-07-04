"use strict";

var httpProxy = require("http-proxy");
var ServiceLoader = require("./service_loader.js");
var logger = require("./logger.js");

module.exports = class IngressRouter {
  constructor() {
    this.hosts = {};
    this.proxy = httpProxy.createProxyServer({});
    this.serviceLoader = new ServiceLoader();
  }

  handleRequest(req, res) {
    var hostname = req.headers.host.split(":")[0];
    var host = this.hosts[hostname];
    if (!host) {
      logger.warn("Failed to find service for " + hostname);
      res.writeHead(404);
      res.end("Service not found");
      return;
    }

    var serviceAddress = "http://" + host.ServiceName + ":" + host.TargetPort;
    logger.info("Proxying to " + serviceAddress);
    this.proxy.web(req, res, {
      target: serviceAddress
    }, function(e) {
      logger.error("Failed to proxy to " + serviceAddress + " - " + e);
      res.writeHead(503);
      res.end("Service not available");
    });
  }

  handleMetrics(req, res) {
    res.writeHead(200, {"Content-Type": "application/json"});
    res.end(JSON.stringify(this.hosts));
  }

  loadServices() {
    var self = this;
    this.serviceLoader.call().then(function(hosts) {
      self.hosts = hosts;
    });
  }
};
