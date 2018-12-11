const goLibrary = require('go-library');
const path = require('path');

const options = {
  destinationPath: './bin',
  repo: 'jurekbarth/rpgo',
  version: 'v0.1.5',
  projectname: 'rpgo'
}


goLibrary(options);
