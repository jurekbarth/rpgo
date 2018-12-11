const goLibrary = require('go-library');
const path = require('path');

const options = {
  destinationPath: './node_modules/.bin',
  repo: 'jurekbarth/rpgo',
  version: 'v0.1.5',
  projectname: 'rpgo'
}


goLibrary(options);
