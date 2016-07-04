"use strict";

var DockerClient = require("./docker_client.js");
var logger = require("./logger.js");

module.exports = class ServiceLoader {
  constructor() {
    this.docker = new DockerClient();
  }

  filterIngressServices(services) {
    return services.filter(function(service) {
      return service.Spec.Labels && service.Spec.Labels.ingress === "true";
    });
  }

  call() {
    logger.info("Starting service sync");
    var self = this;
    return this.docker.listServices().then(function(services) {
      var hosts = {};
      self.filterIngressServices(services).forEach(function(service) {
        if (service.Spec.Labels && service.Spec.Labels["ingress.dnsname"] &&
          service.Spec.Labels["ingress.serviceport"]) {
          logger.info("Registering service " + service.Spec.Name + " as " +
            service.Spec.Labels.dnsname);
          hosts[service.Spec.Labels["ingress.dnsname"]] = {
            ServiceName: service.Spec.Name,
            TargetPort: service.Spec.Labels["ingress.serviceport"]
          };
        }
      });
      return hosts;
    }).catch(function(err) {
      logger.error("Failed to build routing table: " + err);
      return {};
    });
  }
};
