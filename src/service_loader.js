"use strict";

var Docker = require("dockerode");
var logger = require("./logger.js");

module.exports = class ServiceLoader {
  constructor() {
    this.docker = new Docker();
  }

  getAllServices() {
    var self = this;
    return new Promise(function(resolve, reject) {
      self.docker.listServices(function(err, services) {
        if (err) {
          reject(err);
        } else if (services) {
          resolve(services);
        } else {
          reject("No services");
        }
      });
    });
  }

  filterIngressServices(services) {
    return services.filter(function(service) {
      return service.Spec.Labels && service.Spec.Labels.ingress === "true";
    });
  }

  call() {
    logger.info("Starting service sync");
    var self = this;
    return this.getAllServices().then(function(services) {
      var hosts = {};
      self.filterIngressServices(services).forEach(function(service) {
        if (service.Spec.Labels && service.Spec.Labels.dnsname &&
          service.Spec.Labels.serviceport) {
          logger.info("Registering service " + service.Spec.Name + " as " +
            service.Spec.Labels.dnsname);
          hosts[service.Spec.Labels.dnsname] = {
            ServiceName: service.Spec.Name,
            TargetPort: service.Spec.Labels.serviceport
          };
        }
      });
      return hosts;
    }).catch(function(err) {
      logger.error("Failed to retrieve any services: " + err);
      return {};
    });
  }
};
