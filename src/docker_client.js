"use strict";

var http = require("http");

var logger = require("./logger");

module.exports = class DockerClient {
  constructor() {
    this.socket = "/var/run/docker.sock";
  }

  get(apiPath) {
    var options = {
      path: apiPath,
      socketPath: this.socket
    };
    return new Promise(function(resolve, reject) {
      var body = "";
      http.get(options, function(response) {
        response.on("data", function(data) {
          body += data;
        });

        response.on("end", function() {
          if (response.statusCode === 200) {
            resolve(JSON.parse(body));
          } else {
            reject(body);
          }
        });
      }).on("error", function(e) {
        reject(e);
      }).on("timeout", function(e) {
        reject(e);
      });
    });
  }

  listServices() {
    return this.get("/services").then(function(services) {
      return (services ? services : []);
    }).catch(function(e) {
      logger.error("Failed to load services: " + e);
      throw e;
    });
  }
};
