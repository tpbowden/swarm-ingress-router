var winston = require("winston");

var logger = new winston.Logger({
  transports: [
    new (winston.transports.Console)({
      timestamp: function() {
        return new Date().toISOString();
      },
      formatter: function(options) {
        return options.timestamp() + " " + options.level.toUpperCase() + " " +
          (options.message ? options.message : "");
      }
    })
  ]
});

module.exports = logger;
