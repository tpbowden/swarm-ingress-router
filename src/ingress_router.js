"use strict";

var ServiceLoader = require("./service_loader.js");
var logger = require("./logger.js");

module.exports = class IngressRouter {
  constructor() {
    this.hosts = {};
    this.serviceLoader = new ServiceLoader();
  }

  handleRequest(req, res) {
    var hostname = req.headers.host.split(":")[0];
  }

  handleMetrics(req, res) {
    res.writeHead(200);
    res.end(JSON.stringify(this.hosts));
  }

  loadServices() {
    var self = this;
    this.serviceLoader.call().then(function(hosts) {
      self.hosts = hosts;
    });
  }
};
