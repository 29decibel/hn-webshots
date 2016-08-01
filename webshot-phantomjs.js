var page = require('webpage').create();
page.viewportSize = { width: 1920, height: 1080 };

var system = require('system');
var args = system.args;

var url = args[1];
var outputFileName = args[2];

page.open(url, function() {
  page.render(outputFileName,  {format: 'png', quality: '100'});
  phantom.exit();
});
